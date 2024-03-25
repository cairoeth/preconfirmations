// Package operator AVS operator logic.
package operator

import (
	"context"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/Layr-Labs/incredible-squaring-avs/metrics"
	"github.com/cairoeth/preconfirmations/preconf-operator/core/chainio"
	"github.com/cairoeth/preconfirmations/preconf-operator/receiverapi"
	"github.com/cairoeth/preconfirmations/preconf-operator/sse"
	"github.com/cairoeth/preconfirmations/preconf-operator/types"

	"github.com/Layr-Labs/eigensdk-go/chainio/clients"
	sdkelcontracts "github.com/Layr-Labs/eigensdk-go/chainio/clients/elcontracts"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/chainio/txmgr"
	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
	sdkecdsa "github.com/Layr-Labs/eigensdk-go/crypto/ecdsa"
	"github.com/Layr-Labs/eigensdk-go/logging"
	sdkmetrics "github.com/Layr-Labs/eigensdk-go/metrics"
	"github.com/Layr-Labs/eigensdk-go/metrics/collectors/economic"
	rpccalls "github.com/Layr-Labs/eigensdk-go/metrics/collectors/rpc_calls"
	"github.com/Layr-Labs/eigensdk-go/nodeapi"
	"github.com/Layr-Labs/eigensdk-go/signerv2"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
	"github.com/cairoeth/preconfirmations/preconf-share/preconshare"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ybbus/jsonrpc/v3"
)

const (
	AvsName = "preconfirmations"
	SemVer  = "0.0.1"
)

type Operator struct {
	config    types.NodeConfig
	logger    logging.Logger
	ethClient eth.EthClient
	// TODO(samlaf): remove both avsWriter and eigenlayerWrite from operator
	// they are only used for registration, so we should make a special registration package
	// this way, auditing this operator code makes it obvious that operators don't need to
	// write to the chain during the course of their normal operations
	// writing to the chain should be done via the cli only
	metricsReg       *prometheus.Registry
	metrics          metrics.Metrics
	nodeAPI          *nodeapi.NodeApi
	receiveAPI       *receiverapi.ReceiveAPI
	avsWriter        *chainio.AvsWriter
	avsReader        chainio.AvsReaderer
	eigenlayerReader sdkelcontracts.ELReader
	eigenlayerWriter sdkelcontracts.ELWriter
	blsKeypair       *bls.KeyPair
	operatorID       bls.OperatorId
	operatorAddr     common.Address
	// ip address of aggregator
	aggregatorServerIPPortAddr string
	// rpc client to send signed task responses to aggregator
	aggregatorRPCClient AggregatorRPCClienter
	// needed when opting in to avs (allow this service manager contract to slash operator)
	credibleSquaringServiceManagerAddr common.Address
}

// NewOperatorFromConfig the config in core (which is shared with aggregator and challenger)
func NewOperatorFromConfig(c types.NodeConfig) (*Operator, error) {
	var logLevel logging.LogLevel
	if c.Production {
		logLevel = logging.Production
	} else {
		logLevel = logging.Development
	}
	logger, err := logging.NewZapLogger(logLevel)
	if err != nil {
		return nil, err
	}
	reg := prometheus.NewRegistry()
	eigenMetrics := sdkmetrics.NewEigenMetrics(AvsName, c.EigenMetricsIPPortAddress, reg, logger)
	avsAndEigenMetrics := metrics.NewAvsAndEigenMetrics(AvsName, eigenMetrics, reg)

	// Setup Node Api
	nodeAPI := nodeapi.NewNodeApi(AvsName, SemVer, c.NodeAPIIPPortAddress, logger)

	var ethRPCClient eth.EthClient
	if c.EnableMetrics {
		rpcCallsCollector := rpccalls.NewCollector(AvsName, reg)
		ethRPCClient, err = eth.NewInstrumentedClient(c.EthRPCURL, rpcCallsCollector)
		if err != nil {
			logger.Errorf("Cannot create http ethclient", "err", err)
			return nil, err
		}
	} else {
		ethRPCClient, err = eth.NewClient(c.EthRPCURL)
		if err != nil {
			logger.Errorf("Cannot create http ethclient", "err", err)
			return nil, err
		}
	}

	// Setup Receive Api
	receiveAPI := receiverapi.NewReceiveAPI("localhost:8000", logger, ethRPCClient)

	blsKeyPassword, ok := os.LookupEnv("OPERATOR_BLS_KEY_PASSWORD")
	if !ok {
		logger.Warnf("OPERATOR_BLS_KEY_PASSWORD env var not set. using empty string")
	}
	blsKeyPair, err := bls.ReadPrivateKeyFromFile(c.BlsPrivateKeyStorePath, blsKeyPassword)
	if err != nil {
		logger.Errorf("Cannot parse bls private key", "err", err)
		return nil, err
	}
	// TODO(samlaf): should we add the chainID to the config instead?
	// this way we can prevent creating a signer that signs on mainnet by mistake
	// if the config says chainID=5, then we can only create a goerli signer
	chainID, err := ethRPCClient.ChainID(context.Background())
	if err != nil {
		logger.Error("Cannot get chainID", "err", err)
		return nil, err
	}

	ecdsaKeyPassword, ok := os.LookupEnv("OPERATOR_ECDSA_KEY_PASSWORD")
	if !ok {
		logger.Warnf("OPERATOR_ECDSA_KEY_PASSWORD env var not set. using empty string")
	}

	signerV2, _, err := signerv2.SignerFromConfig(signerv2.Config{
		KeystorePath: c.EcdsaPrivateKeyStorePath,
		Password:     ecdsaKeyPassword,
	}, chainID)
	if err != nil {
		panic(err)
	}
	chainioConfig := clients.BuildAllConfig{
		EthHttpUrl:                 c.EthRPCURL,
		EthWsUrl:                   c.EthRPCURL,
		RegistryCoordinatorAddr:    c.AVSRegistryCoordinatorAddress,
		OperatorStateRetrieverAddr: c.OperatorStateRetrieverAddress,
		AvsName:                    AvsName,
		PromMetricsIpPortAddress:   c.EigenMetricsIPPortAddress,
	}

	sdkClients, err := clients.BuildAll(chainioConfig, common.HexToAddress(c.OperatorAddress), signerV2, logger)
	if err != nil {
		panic(err)
	}

	txMgr := txmgr.NewSimpleTxManager(ethRPCClient, logger, signerV2, common.HexToAddress(c.OperatorAddress))

	avsWriter, err := chainio.BuildAvsWriter(
		txMgr, common.HexToAddress(c.AVSRegistryCoordinatorAddress),
		common.HexToAddress(c.OperatorStateRetrieverAddress), ethRPCClient, logger,
	)
	if err != nil {
		logger.Error("Cannot create AvsWriter", "err", err)
		return nil, err
	}

	avsReader, err := chainio.BuildAvsReader(
		common.HexToAddress(c.AVSRegistryCoordinatorAddress),
		common.HexToAddress(c.OperatorStateRetrieverAddress),
		ethRPCClient, logger)
	if err != nil {
		logger.Error("Cannot create AvsReader", "err", err)
		return nil, err
	}

	// We must register the economic metrics separately because they are exported metrics (from jsonrpc or subgraph calls)
	// and not instrumented metrics: see https://prometheus.io/docs/instrumenting/writing_clientlibs/#overall-structure
	quorumNames := map[sdktypes.QuorumNum]string{
		0: "quorum0",
	}
	economicMetricsCollector := economic.NewCollector(
		sdkClients.ElChainReader, sdkClients.AvsRegistryChainReader,
		AvsName, logger, common.HexToAddress(c.OperatorAddress), quorumNames)
	reg.MustRegister(economicMetricsCollector)

	aggregatorRPCClient, err := NewAggregatorRPCClient(c.AggregatorServerIPPortAddress, logger, avsAndEigenMetrics)
	if err != nil {
		logger.Error("Cannot create AggregatorRPCClient. Is aggregator running?", "err", err)
		return nil, err
	}

	operator := &Operator{
		config:                             c,
		logger:                             logger,
		metricsReg:                         reg,
		metrics:                            avsAndEigenMetrics,
		nodeAPI:                            nodeAPI,
		receiveAPI:                         receiveAPI,
		ethClient:                          ethRPCClient,
		avsWriter:                          avsWriter,
		avsReader:                          avsReader,
		eigenlayerReader:                   sdkClients.ElChainReader,
		eigenlayerWriter:                   sdkClients.ElChainWriter,
		blsKeypair:                         blsKeyPair,
		operatorAddr:                       common.HexToAddress(c.OperatorAddress),
		aggregatorServerIPPortAddr:         c.AggregatorServerIPPortAddress,
		aggregatorRPCClient:                aggregatorRPCClient,
		credibleSquaringServiceManagerAddr: common.HexToAddress(c.AVSRegistryCoordinatorAddress),
		operatorID:                         [32]byte{0}, // this is set below

	}

	if c.RegisterOperatorOnStartup {
		operatorEcdsaPrivateKey, err := sdkecdsa.ReadKey(
			c.EcdsaPrivateKeyStorePath,
			ecdsaKeyPassword,
		)
		if err != nil {
			return nil, err
		}
		operator.registerOperatorOnStartup(operatorEcdsaPrivateKey, common.HexToAddress(c.TokenStrategyAddr))
	}

	// OperatorId is set in contract during registration so we get it after registering operator.
	operatorID, err := sdkClients.AvsRegistryChainReader.GetOperatorId(&bind.CallOpts{}, operator.operatorAddr)
	if err != nil {
		logger.Error("Cannot get operator id", "err", err)
		return nil, err
	}
	operator.operatorID = operatorID
	logger.Info("Operator info",
		"operatorID", operatorID,
		"operatorAddr", c.OperatorAddress,
		"operatorG1Pubkey", operator.blsKeypair.GetPubKeyG1(),
		"operatorG2Pubkey", operator.blsKeypair.GetPubKeyG2(),
	)

	return operator, nil
}

func (o *Operator) Start(ctx context.Context) error {
	operatorIsRegistered, err := o.avsReader.IsOperatorRegistered(&bind.CallOpts{}, o.operatorAddr)
	if err != nil {
		o.logger.Error("Error checking if operator is registered", "err", err)
		return err
	}
	if !operatorIsRegistered {
		// We bubble the error all the way up instead of using logger.Fatal because logger.Fatal prints a huge stack trace
		// that hides the actual error message. This error msg is more explicit and doesn't require showing a stack trace to the user.
		return fmt.Errorf("operator is not registered. Registering operator using the operator-cli before starting operator")
	}

	o.logger.Infof("Starting operator.")

	if o.config.EnableNodeAPI {
		o.nodeAPI.Start()
	}

	o.receiveAPI.Start()

	var metricsErrChan <-chan error
	if o.config.EnableMetrics {
		metricsErrChan = o.metrics.Start(ctx, o.metricsReg)
	} else {
		metricsErrChan = make(chan error, 1)
	}

	client := sse.New("localhost:6379")
	eventChan := make(chan sse.Event)

	sub, err := client.Subscribe(eventChan)
	if err != nil {
		o.logger.Error("Cannot subscribe to preconf-share endpoint", "err", err)
	}

	o.logger.Infof("Listening to preconf-share stream..")

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-metricsErrChan:
			// TODO(samlaf); we should also register the service as unhealthy in the node api
			// https://eigen.nethermind.io/docs/spec/api/
			o.logger.Fatal("Error in metrics server", "err", err)
		case event := <-eventChan:
			if event.Error != nil {
				o.logger.Error("Error occurred from preconf-share stream", "err", event.Error)
				sub.Stop()
			}

			ecdsaKeyPassword, _ := os.LookupEnv("OPERATOR_ECDSA_KEY_PASSWORD")

			privateKey, err := sdkecdsa.ReadKey(
				o.config.EcdsaPrivateKeyStorePath,
				ecdsaKeyPassword,
			)
			if err != nil {
				o.logger.Error("Error getting the private key", "err", err)
				sub.Stop()
			}

			// TODO: update this
			dataToSign := "{request: wwdwdw}"

			// keccak256 hash of the data
			dataBytes := []byte(dataToSign)
			hashData := crypto.Keccak256Hash(dataBytes)

			signatureBytes, err := crypto.Sign(hashData.Bytes(), privateKey)
			if err != nil {
				o.logger.Error("Error occurred signing preconfirmation", "err", err)
				sub.Stop()
			}

			response := preconshare.ConfirmRequestArgs{
				Version: "v0.1",
				Preconf: preconshare.ConfirmPreconf{
					Hash:  event.Data.Hash,
					Block: 510,
				},
				Signature: (*hexutil.Bytes)(&signatureBytes),
				Endpoint:  "http://localhost:8000/receive",
			}

			o.logger.Infof("Sending preconfirmation for hash", event.Data.Hash.Hex())

			// Send signed preconfirmation to preconf-share
			rpcClient := jsonrpc.NewClient("http://localhost:8080")
			_, err = rpcClient.Call(ctx, "preconf_confirmRequest", []*preconshare.ConfirmRequestArgs{&response})
			if err != nil {
				o.logger.Error("Failed to send signed preconfirmation", "err", err)
				sub.Stop()
			}
		}
	}
}
