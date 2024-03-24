package operator

// OUTDATED
// This file contains cli functions for registering an operator with the AVS and printing status
// However, all of this functionality has been moved to the plugin/ package
// we are just waiting for eigenlayer-cli to be open sourced so we can completely get rid of this registration functionality in the operator

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	eigenSdkTypes "github.com/Layr-Labs/eigensdk-go/types"
)

func (o *Operator) registerOperatorOnStartup(
	operatorEcdsaPrivateKey *ecdsa.PrivateKey,
	mockTokenStrategyAddr common.Address,
) {
	err := o.RegisterOperatorWithEigenlayer()
	if err != nil {
		// This error might only be that the operator was already registered with eigenlayer, so we don't want to fatal
		o.logger.Error("Error registering operator with eigenlayer", "err", err)
	} else {
		o.logger.Infof("Registered operator with eigenlayer")
	}

	// TODO(samlaf): shouldn't hardcode number here
	amount := big.NewInt(1000)
	err = o.DepositIntoStrategy(mockTokenStrategyAddr, amount)
	if err != nil {
		o.logger.Fatal("Error depositing into strategy", "err", err)
	}
	o.logger.Infof("Deposited %s into strategy %s", amount, mockTokenStrategyAddr)

	err = o.RegisterOperatorWithAvs(operatorEcdsaPrivateKey)
	if err != nil {
		// This error might only be that the operator was already registered in AVS, so we don't want to fatal
		o.logger.Error("Error registering operator with avs", "err", err)
	}
	o.logger.Infof("Registered operator with avs")
}

func (o *Operator) RegisterOperatorWithEigenlayer() error {
	op := eigenSdkTypes.Operator{
		Address:                 o.operatorAddr.String(),
		EarningsReceiverAddress: o.operatorAddr.String(),
	}
	_, err := o.eigenlayerWriter.RegisterAsOperator(context.Background(), op)
	if err != nil {
		o.logger.Errorf("Error registering operator with eigenlayer")
		return err
	}
	return nil
}

func (o *Operator) DepositIntoStrategy(strategyAddr common.Address, amount *big.Int) error {
	_, tokenAddr, err := o.eigenlayerReader.GetStrategyAndUnderlyingToken(&bind.CallOpts{}, strategyAddr)
	if err != nil {
		o.logger.Error("Failed to fetch strategy contract", "err", err)
		return err
	}

	o.logger.Info("ERC20Mock contract", "address", tokenAddr)

	contractErc20Mock, err := o.avsReader.GetErc20Mock(context.Background(), tokenAddr)
	if err != nil {
		o.logger.Error("Failed to fetch ERC20Mock contract", "err", err)
		return err
	}

	txOpts, _ := o.avsWriter.TxMgr.GetNoSendTxOpts()
	tx, err := contractErc20Mock.Mint(txOpts, o.operatorAddr, amount)
	if err != nil {
		o.logger.Error("Error assembling Mint tx", "err", err)
		return err
	}
	_, err = o.avsWriter.TxMgr.Send(context.Background(), tx)
	if err != nil {
		o.logger.Errorf("Error submitting Mint tx")
		return err
	}

	_, err = o.eigenlayerWriter.DepositERC20IntoStrategy(context.Background(), strategyAddr, amount)
	if err != nil {
		o.logger.Errorf("Error depositing into strategy", "err", err)
		return err
	}
	return nil
}

// RegisterOperatorWithAvs Registration specific functions
func (o *Operator) RegisterOperatorWithAvs(
	operatorEcdsaKeyPair *ecdsa.PrivateKey,
) error {
	// hardcode these things for now
	quorumNumbers := []byte{0}
	socket := "Not Needed"
	operatorToAvsRegistrationSigSalt := [32]byte{123}
	curBlockNum, err := o.ethClient.BlockNumber(context.Background())
	if err != nil {
		o.logger.Errorf("Unable to get current block number")
		return err
	}
	curBlock, err := o.ethClient.BlockByNumber(context.Background(), big.NewInt(int64(curBlockNum)))
	if err != nil {
		o.logger.Errorf("Unable to get current block")
		return err
	}
	sigValidForSeconds := int64(1_000_000)
	operatorToAvsRegistrationSigExpiry := big.NewInt(int64(curBlock.Time()) + sigValidForSeconds)
	_, err = o.avsWriter.RegisterOperatorInQuorumWithAVSRegistryCoordinator(
		context.Background(),
		operatorEcdsaKeyPair, operatorToAvsRegistrationSigSalt, operatorToAvsRegistrationSigExpiry,
		o.blsKeypair, quorumNumbers, socket,
	)
	if err != nil {
		o.logger.Errorf("Unable to register operator with avs registry coordinator")
		return err
	}
	o.logger.Infof("Registered operator with avs registry coordinator.")

	return nil
}
