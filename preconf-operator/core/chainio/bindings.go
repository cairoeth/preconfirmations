package chainio

import (
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/logging"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	regcoord "github.com/Layr-Labs/eigensdk-go/contracts/bindings/RegistryCoordinator"
	erc20mock "github.com/cairoeth/preconfirmations/contracts/bindings/ERC20Mock"
	preconfchallengemanager "github.com/cairoeth/preconfirmations/contracts/bindings/PreconfChallengeManager"
	preconfservicemanager "github.com/cairoeth/preconfirmations/contracts/bindings/PreconfServiceManager"
)

type AvsManagersBindings struct {
	ChallengeManager *preconfchallengemanager.ContractPreconfChallengeManager
	ServiceManager   *preconfservicemanager.ContractPreconfServiceManager
	ethClient        eth.EthClient
	logger           logging.Logger
}

func NewAvsManagersBindings(registryCoordinatorAddr, operatorStateRetrieverAddr common.Address, ethclient eth.EthClient, logger logging.Logger) (*AvsManagersBindings, error) {
	contractRegistryCoordinator, err := regcoord.NewContractRegistryCoordinator(registryCoordinatorAddr, ethclient)
	if err != nil {
		return nil, err
	}
	serviceManagerAddr, err := contractRegistryCoordinator.ServiceManager(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	contractServiceManager, err := preconfservicemanager.NewContractPreconfServiceManager(serviceManagerAddr, ethclient)
	if err != nil {
		logger.Error("Failed to fetch IServiceManager contract", "err", err)
		return nil, err
	}

	taskManagerAddr, err := contractServiceManager.ChallengeManager(&bind.CallOpts{})
	if err != nil {
		logger.Error("Failed to fetch ChallengeManager address", "err", err)
		return nil, err
	}
	contractChallengeManager, err := preconfchallengemanager.NewContractPreconfChallengeManager(taskManagerAddr, ethclient)
	if err != nil {
		logger.Error("Failed to fetch IChallengeManager contract", "err", err)
		return nil, err
	}

	return &AvsManagersBindings{
		ServiceManager:   contractServiceManager,
		ChallengeManager: contractChallengeManager,
		ethClient:        ethclient,
		logger:           logger,
	}, nil
}

func (b *AvsManagersBindings) GetErc20Mock(tokenAddr common.Address) (*erc20mock.ContractERC20Mock, error) {
	contractErc20Mock, err := erc20mock.NewContractERC20Mock(tokenAddr, b.ethClient)
	if err != nil {
		b.logger.Error("Failed to fetch ERC20Mock contract", "err", err)
		return nil, err
	}
	return contractErc20Mock, nil
}
