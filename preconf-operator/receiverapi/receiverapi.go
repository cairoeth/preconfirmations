// Package receiverapi contains the logic to receive preconfirmation callbacks.
package receiverapi

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

type ReceiveAPI struct {
	ipPortAddr string
	logger     logging.Logger
	client     eth.EthClient
}

type ReceiveResponse struct {
	jsonrpc string
	id      int
}

type ReceiveTx struct {
	Tx *hexutil.Bytes `json:"tx,omitempty"`
}

type JSONRPCRequest struct {
	ID      interface{} `json:"id"`
	Method  string      `json:"method"`
	Params  []ReceiveTx `json:"params"`
	Version string      `json:"jsonrpc,omitempty"`
}

func NewReceiveAPI(IPPortAddr string, logger logging.Logger, client eth.EthClient) *ReceiveAPI {
	receiveAPI := &ReceiveAPI{
		ipPortAddr: IPPortAddr,
		logger:     logger,
		client:     client,
	}
	return receiveAPI
}

// Start starts the receiver api server in a goroutine
func (api *ReceiveAPI) Start() <-chan error {
	api.logger.Infof("Starting receiver api server at address %v", api.ipPortAddr)

	mux := http.NewServeMux()
	httpServer := http.Server{
		Addr:    api.ipPortAddr,
		Handler: mux,
	}

	mux.HandleFunc("/receive", api.receive)
	errChan := run(api.logger, &httpServer)
	return errChan
}

func (api *ReceiveAPI) receive(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var txs JSONRPCRequest
	err := decoder.Decode(&txs)
	if err != nil {
		api.logger.Error("could not read request body", "err", err)
	}
	api.logger.Info("processed receive", "txs", txs)

	for _, tx := range txs.Params {
		api.logger.Info("tx received", "tx", tx.Tx)

		var txBytes types.Transaction
		err := txBytes.UnmarshalBinary(*tx.Tx)
		if err != nil {
			api.logger.Error("could not decode tx into bytes", "err", err)
		}

		err = api.client.SendTransaction(context.Background(), &txBytes)
		if err != nil {
			api.logger.Error("could not send tx to rpc", "err", err)
		}
	}

	response := &ReceiveResponse{"2.0", 1}
	err = jsonResponse(w, response)
	if err != nil {
		api.logger.Error("Error in receive endpoint", "err", err)
	}
}

func run(logger logging.Logger, httpServer *http.Server) <-chan error {
	errChan := make(chan error, 1)
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		<-ctx.Done()

		logger.Info("shutdown signal received")

		defer func() {
			stop()
			close(errChan)
		}()

		if err := httpServer.Shutdown(context.Background()); err != nil {
			errChan <- err
		}
		logger.Info("shutdown completed")
	}()

	go func() {
		logger.Info("receiver api server running", "addr", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()

	return errChan
}

func jsonResponse(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	return err
}
