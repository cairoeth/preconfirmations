// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractPreconfServiceManager

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// ISignatureUtilsSignatureWithSaltAndExpiry is an auto generated low-level Go binding around an user-defined struct.
type ISignatureUtilsSignatureWithSaltAndExpiry struct {
	Signature []byte
	Salt      [32]byte
	Expiry    *big.Int
}

// ContractPreconfServiceManagerMetaData contains all meta data concerning the ContractPreconfServiceManager contract.
var ContractPreconfServiceManagerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_avsDirectory\",\"type\":\"address\",\"internalType\":\"contractIAVSDirectory\"},{\"name\":\"_registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"},{\"name\":\"_stakeRegistry\",\"type\":\"address\",\"internalType\":\"contractIStakeRegistry\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"avsDirectory\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"challengeManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractPreconfChallengeManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deregisterOperatorFromAVS\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"freezeOperator\",\"inputs\":[{\"name\":\"operatorAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getOperatorRestakedStrategies\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRestakeableStrategies\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_pauserRegistry\",\"type\":\"address\",\"internalType\":\"contractIPauserRegistry\"},{\"name\":\"_initialPausedStatus\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_challengeManager\",\"type\":\"address\",\"internalType\":\"contractPreconfChallengeManager\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pauseAll\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pauserRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIPauserRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerOperatorToAVS\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operatorSignature\",\"type\":\"tuple\",\"internalType\":\"structISignatureUtils.SignatureWithSaltAndExpiry\",\"components\":[{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"expiry\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMetadataURI\",\"inputs\":[{\"name\":\"_metadataURI\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPauserRegistry\",\"inputs\":[{\"name\":\"newPauserRegistry\",\"type\":\"address\",\"internalType\":\"contractIPauserRegistry\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PauserRegistrySet\",\"inputs\":[{\"name\":\"pauserRegistry\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIPauserRegistry\"},{\"name\":\"newPauserRegistry\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIPauserRegistry\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"NotChallengeManager\",\"inputs\":[]}]",
	Bin: "0x60e06040523480156200001157600080fd5b5060405162001f6838038062001f6883398101604081905262000034916200014c565b6001600160a01b0380841660c052808316608052811660a0528282826200005a62000071565b506200006891505062000071565b505050620001a0565b600054610100900460ff1615620000de5760405162461bcd60e51b815260206004820152602760248201527f496e697469616c697a61626c653a20636f6e747261637420697320696e697469604482015266616c697a696e6760c81b606482015260840160405180910390fd5b60005460ff908116101562000131576000805460ff191660ff9081179091556040519081527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b565b6001600160a01b03811681146200014957600080fd5b50565b6000806000606084860312156200016257600080fd5b83516200016f8162000133565b6020850151909350620001828162000133565b6040850151909250620001958162000133565b809150509250925092565b60805160a05160c051611d36620002326000396000818161021201528181610bf401528181610cb60152610d8a0152600081816106be0152818161081a015281816108b101528181610e810152818161100501526110a40152600081816104e901528181610578015281816105f801528181610c6201528181610d2e01528181610dbf0152610f600152611d366000f3fe608060405234801561001057600080fd5b50600436106101215760003560e01c80636b3aa72e116100ad5780639926ee7d116100715780639926ee7d14610275578063a364f4da14610288578063e481af9d1461029b578063f2fde38b146102a3578063fabc1cbc146102b657600080fd5b80636b3aa72e14610210578063715018a614610236578063750521f51461023e578063886f1195146102515780638da5cb5b1461026457600080fd5b8063358394d8116100f4578063358394d81461019e57806338c8ee64146101b1578063595c6a67146101c45780635ac86ab7146101cc5780635c975abb146101ff57600080fd5b8063023a96fe1461012657806310d67a2f14610156578063136439dd1461016b57806333cfb7b71461017e575b600080fd5b60c954610139906001600160a01b031681565b6040516001600160a01b0390911681526020015b60405180910390f35b6101696101643660046116eb565b6102c9565b005b61016961017936600461170f565b610385565b61019161018c3660046116eb565b6104c4565b60405161014d9190611728565b6101696101ac366004611775565b610994565b6101696101bf3660046116eb565b610acf565b610169610afa565b6101ef6101da3660046117d7565b609854600160ff9092169190911b9081161490565b604051901515815260200161014d565b60985460405190815260200161014d565b7f0000000000000000000000000000000000000000000000000000000000000000610139565b610169610bc1565b61016961024c3660046118a9565b610bd5565b609754610139906001600160a01b031681565b6033546001600160a01b0316610139565b6101696102833660046118fa565b610c57565b6101696102963660046116eb565b610d23565b610191610db9565b6101696102b13660046116eb565b611183565b6101696102c436600461170f565b6111f9565b609760009054906101000a90046001600160a01b03166001600160a01b031663eab66d7a6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561031c573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061034091906119a5565b6001600160a01b0316336001600160a01b0316146103795760405162461bcd60e51b8152600401610370906119c2565b60405180910390fd5b61038281611355565b50565b60975460405163237dfb4760e11b81523360048201526001600160a01b03909116906346fbf68e90602401602060405180830381865afa1580156103cd573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103f19190611a0c565b61040d5760405162461bcd60e51b815260040161037090611a2e565b609854818116146104865760405162461bcd60e51b815260206004820152603860248201527f5061757361626c652e70617573653a20696e76616c696420617474656d70742060448201527f746f20756e70617573652066756e6374696f6e616c69747900000000000000006064820152608401610370565b609881905560405181815233907fab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d906020015b60405180910390a250565b6040516309aa152760e11b81526001600160a01b0382811660048301526060916000917f000000000000000000000000000000000000000000000000000000000000000016906313542a4e90602401602060405180830381865afa158015610530573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105549190611a76565b60405163871ef04960e01b8152600481018290529091506000906001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063871ef04990602401602060405180830381865afa1580156105bf573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105e39190611a8f565b90506001600160c01b038116158061067d57507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316639aa1653d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610654573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106789190611ab8565b60ff16155b1561069957505060408051600081526020810190915292915050565b60006106ad826001600160c01b031661144c565b90506000805b8251811015610783577f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316633ca5a5f58483815181106106fd576106fd611ad5565b01602001516040516001600160e01b031960e084901b16815260f89190911c6004820152602401602060405180830381865afa158015610741573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107659190611a76565b61076f9083611b01565b91508061077b81611b19565b9150506106b3565b5060008167ffffffffffffffff81111561079f5761079f6117f4565b6040519080825280602002602001820160405280156107c8578160200160208202803683370190505b5090506000805b84518110156109875760008582815181106107ec576107ec611ad5565b0160200151604051633ca5a5f560e01b815260f89190911c6004820181905291506000906001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690633ca5a5f590602401602060405180830381865afa158015610861573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108859190611a76565b905060005b81811015610971576040516356e4026d60e11b815260ff84166004820152602481018290527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063adc804da906044016040805180830381865afa1580156108ff573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109239190611b34565b6000015186868151811061093957610939611ad5565b6001600160a01b03909216602092830291909101909101528461095b81611b19565b955050808061096990611b19565b91505061088a565b505050808061097f90611b19565b9150506107cf565b5090979650505050505050565b600054610100900460ff16158080156109b45750600054600160ff909116105b806109ce5750303b1580156109ce575060005460ff166001145b610a315760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b6064820152608401610370565b6000805460ff191660011790558015610a54576000805461ff0019166101001790555b610a5e858561150f565b610a67836115f9565b60c980546001600160a01b0319166001600160a01b0384161790558015610ac8576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b5050505050565b60c9546001600160a01b03163314610382576040516301d6da7b60e41b815260040160405180910390fd5b60975460405163237dfb4760e11b81523360048201526001600160a01b03909116906346fbf68e90602401602060405180830381865afa158015610b42573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b669190611a0c565b610b825760405162461bcd60e51b815260040161037090611a2e565b600019609881905560405190815233907fab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d9060200160405180910390a2565b610bc961164b565b610bd360006115f9565b565b610bdd61164b565b60405163a98fb35560e01b81526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063a98fb35590610c29908490600401611bf1565b600060405180830381600087803b158015610c4357600080fd5b505af1158015610ac8573d6000803e3d6000fd5b336001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614610c9f5760405162461bcd60e51b815260040161037090611c04565b604051639926ee7d60e01b81526001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690639926ee7d90610ced9085908590600401611c7c565b600060405180830381600087803b158015610d0757600080fd5b505af1158015610d1b573d6000803e3d6000fd5b505050505050565b336001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614610d6b5760405162461bcd60e51b815260040161037090611c04565b6040516351b27a6d60e11b81526001600160a01b0382811660048301527f0000000000000000000000000000000000000000000000000000000000000000169063a364f4da90602401610c29565b606060007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316639aa1653d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610e1b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e3f9190611ab8565b60ff16905080610e5d57505060408051600081526020810190915290565b6000805b82811015610f1257604051633ca5a5f560e01b815260ff821660048201527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031690633ca5a5f590602401602060405180830381865afa158015610ed0573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ef49190611a76565b610efe9083611b01565b915080610f0a81611b19565b915050610e61565b5060008167ffffffffffffffff811115610f2e57610f2e6117f4565b604051908082528060200260200182016040528015610f57578160200160208202803683370190505b5090506000805b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316639aa1653d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610fbc573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610fe09190611ab8565b60ff1681101561117957604051633ca5a5f560e01b815260ff821660048201526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031690633ca5a5f590602401602060405180830381865afa158015611054573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906110789190611a76565b905060005b81811015611164576040516356e4026d60e11b815260ff84166004820152602481018290527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063adc804da906044016040805180830381865afa1580156110f2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906111169190611b34565b6000015185858151811061112c5761112c611ad5565b6001600160a01b03909216602092830291909101909101528361114e81611b19565b945050808061115c90611b19565b91505061107d565b5050808061117190611b19565b915050610f5e565b5090949350505050565b61118b61164b565b6001600160a01b0381166111f05760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152608401610370565b610382816115f9565b609760009054906101000a90046001600160a01b03166001600160a01b031663eab66d7a6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561124c573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061127091906119a5565b6001600160a01b0316336001600160a01b0316146112a05760405162461bcd60e51b8152600401610370906119c2565b60985419811960985419161461131e5760405162461bcd60e51b815260206004820152603860248201527f5061757361626c652e756e70617573653a20696e76616c696420617474656d7060448201527f7420746f2070617573652066756e6374696f6e616c69747900000000000000006064820152608401610370565b609881905560405181815233907f3582d1828e26bf56bd801502bc021ac0bc8afb57c826e4986b45593c8fad389c906020016104b9565b6001600160a01b0381166113e35760405162461bcd60e51b815260206004820152604960248201527f5061757361626c652e5f73657450617573657252656769737472793a206e657760448201527f50617573657252656769737472792063616e6e6f7420626520746865207a65726064820152686f206164647265737360b81b608482015260a401610370565b609754604080516001600160a01b03928316815291831660208301527f6e9fcd539896fca60e8b0f01dd580233e48a6b0f7df013b89ba7f565869acdb6910160405180910390a1609780546001600160a01b0319166001600160a01b0392909216919091179055565b606060008061145a846116a5565b61ffff1667ffffffffffffffff811115611476576114766117f4565b6040519080825280601f01601f1916602001820160405280156114a0576020820181803683370190505b5090506000805b8251821080156114b8575061010081105b15611179576001811b9350858416156114ff578060f81b8383815181106114e1576114e1611ad5565b60200101906001600160f81b031916908160001a9053508160010191505b61150881611b19565b90506114a7565b6097546001600160a01b031615801561153057506001600160a01b03821615155b6115b25760405162461bcd60e51b815260206004820152604760248201527f5061757361626c652e5f696e697469616c697a655061757365723a205f696e6960448201527f7469616c697a6550617573657228292063616e206f6e6c792062652063616c6c6064820152666564206f6e636560c81b608482015260a401610370565b609881905560405181815233907fab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d9060200160405180910390a26115f582611355565b5050565b603380546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b6033546001600160a01b03163314610bd35760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610370565b6000805b82156116d0576116ba600184611cc7565b90921691806116c881611cde565b9150506116a9565b92915050565b6001600160a01b038116811461038257600080fd5b6000602082840312156116fd57600080fd5b8135611708816116d6565b9392505050565b60006020828403121561172157600080fd5b5035919050565b6020808252825182820181905260009190848201906040850190845b818110156117695783516001600160a01b031683529284019291840191600101611744565b50909695505050505050565b6000806000806080858703121561178b57600080fd5b8435611796816116d6565b93506020850135925060408501356117ad816116d6565b915060608501356117bd816116d6565b939692955090935050565b60ff8116811461038257600080fd5b6000602082840312156117e957600080fd5b8135611708816117c8565b634e487b7160e01b600052604160045260246000fd5b6040516060810167ffffffffffffffff8111828210171561182d5761182d6117f4565b60405290565b600067ffffffffffffffff8084111561184e5761184e6117f4565b604051601f8501601f19908116603f01168101908282118183101715611876576118766117f4565b8160405280935085815286868601111561188f57600080fd5b858560208301376000602087830101525050509392505050565b6000602082840312156118bb57600080fd5b813567ffffffffffffffff8111156118d257600080fd5b8201601f810184136118e357600080fd5b6118f284823560208401611833565b949350505050565b6000806040838503121561190d57600080fd5b8235611918816116d6565b9150602083013567ffffffffffffffff8082111561193557600080fd5b908401906060828703121561194957600080fd5b61195161180a565b82358281111561196057600080fd5b83019150601f8201871361197357600080fd5b61198287833560208501611833565b815260208301356020820152604083013560408201528093505050509250929050565b6000602082840312156119b757600080fd5b8151611708816116d6565b6020808252602a908201527f6d73672e73656e646572206973206e6f74207065726d697373696f6e6564206160408201526939903ab73830bab9b2b960b11b606082015260800190565b600060208284031215611a1e57600080fd5b8151801515811461170857600080fd5b60208082526028908201527f6d73672e73656e646572206973206e6f74207065726d697373696f6e6564206160408201526739903830bab9b2b960c11b606082015260800190565b600060208284031215611a8857600080fd5b5051919050565b600060208284031215611aa157600080fd5b81516001600160c01b038116811461170857600080fd5b600060208284031215611aca57600080fd5b8151611708816117c8565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b60008219821115611b1457611b14611aeb565b500190565b6000600019821415611b2d57611b2d611aeb565b5060010190565b600060408284031215611b4657600080fd5b6040516040810181811067ffffffffffffffff82111715611b6957611b696117f4565b6040528251611b77816116d6565b815260208301516bffffffffffffffffffffffff81168114611b9857600080fd5b60208201529392505050565b6000815180845260005b81811015611bca57602081850181015186830182015201611bae565b81811115611bdc576000602083870101525b50601f01601f19169290920160200192915050565b6020815260006117086020830184611ba4565b60208082526052908201527f536572766963654d616e61676572426173652e6f6e6c7952656769737472794360408201527f6f6f7264696e61746f723a2063616c6c6572206973206e6f742074686520726560608201527133b4b9ba393c9031b7b7b93234b730ba37b960711b608082015260a00190565b60018060a01b0383168152604060208201526000825160606040840152611ca660a0840182611ba4565b90506020840151606084015260408401516080840152809150509392505050565b600082821015611cd957611cd9611aeb565b500390565b600061ffff80831681811415611cf657611cf6611aeb565b600101939250505056fea2646970667358221220cf95b333f0b94dda6934370be3efe83f632adf3ea65eb3fa7c9c430b81c64b6e64736f6c634300080c0033",
}

// ContractPreconfServiceManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractPreconfServiceManagerMetaData.ABI instead.
var ContractPreconfServiceManagerABI = ContractPreconfServiceManagerMetaData.ABI

// ContractPreconfServiceManagerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractPreconfServiceManagerMetaData.Bin instead.
var ContractPreconfServiceManagerBin = ContractPreconfServiceManagerMetaData.Bin

// DeployContractPreconfServiceManager deploys a new Ethereum contract, binding an instance of ContractPreconfServiceManager to it.
func DeployContractPreconfServiceManager(auth *bind.TransactOpts, backend bind.ContractBackend, _avsDirectory common.Address, _registryCoordinator common.Address, _stakeRegistry common.Address) (common.Address, *types.Transaction, *ContractPreconfServiceManager, error) {
	parsed, err := ContractPreconfServiceManagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractPreconfServiceManagerBin), backend, _avsDirectory, _registryCoordinator, _stakeRegistry)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractPreconfServiceManager{ContractPreconfServiceManagerCaller: ContractPreconfServiceManagerCaller{contract: contract}, ContractPreconfServiceManagerTransactor: ContractPreconfServiceManagerTransactor{contract: contract}, ContractPreconfServiceManagerFilterer: ContractPreconfServiceManagerFilterer{contract: contract}}, nil
}

// ContractPreconfServiceManager is an auto generated Go binding around an Ethereum contract.
type ContractPreconfServiceManager struct {
	ContractPreconfServiceManagerCaller     // Read-only binding to the contract
	ContractPreconfServiceManagerTransactor // Write-only binding to the contract
	ContractPreconfServiceManagerFilterer   // Log filterer for contract events
}

// ContractPreconfServiceManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractPreconfServiceManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractPreconfServiceManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractPreconfServiceManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractPreconfServiceManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractPreconfServiceManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractPreconfServiceManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractPreconfServiceManagerSession struct {
	Contract     *ContractPreconfServiceManager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                  // Call options to use throughout this session
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// ContractPreconfServiceManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractPreconfServiceManagerCallerSession struct {
	Contract *ContractPreconfServiceManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                        // Call options to use throughout this session
}

// ContractPreconfServiceManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractPreconfServiceManagerTransactorSession struct {
	Contract     *ContractPreconfServiceManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                        // Transaction auth options to use throughout this session
}

// ContractPreconfServiceManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractPreconfServiceManagerRaw struct {
	Contract *ContractPreconfServiceManager // Generic contract binding to access the raw methods on
}

// ContractPreconfServiceManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractPreconfServiceManagerCallerRaw struct {
	Contract *ContractPreconfServiceManagerCaller // Generic read-only contract binding to access the raw methods on
}

// ContractPreconfServiceManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractPreconfServiceManagerTransactorRaw struct {
	Contract *ContractPreconfServiceManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractPreconfServiceManager creates a new instance of ContractPreconfServiceManager, bound to a specific deployed contract.
func NewContractPreconfServiceManager(address common.Address, backend bind.ContractBackend) (*ContractPreconfServiceManager, error) {
	contract, err := bindContractPreconfServiceManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractPreconfServiceManager{ContractPreconfServiceManagerCaller: ContractPreconfServiceManagerCaller{contract: contract}, ContractPreconfServiceManagerTransactor: ContractPreconfServiceManagerTransactor{contract: contract}, ContractPreconfServiceManagerFilterer: ContractPreconfServiceManagerFilterer{contract: contract}}, nil
}

// NewContractPreconfServiceManagerCaller creates a new read-only instance of ContractPreconfServiceManager, bound to a specific deployed contract.
func NewContractPreconfServiceManagerCaller(address common.Address, caller bind.ContractCaller) (*ContractPreconfServiceManagerCaller, error) {
	contract, err := bindContractPreconfServiceManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractPreconfServiceManagerCaller{contract: contract}, nil
}

// NewContractPreconfServiceManagerTransactor creates a new write-only instance of ContractPreconfServiceManager, bound to a specific deployed contract.
func NewContractPreconfServiceManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractPreconfServiceManagerTransactor, error) {
	contract, err := bindContractPreconfServiceManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractPreconfServiceManagerTransactor{contract: contract}, nil
}

// NewContractPreconfServiceManagerFilterer creates a new log filterer instance of ContractPreconfServiceManager, bound to a specific deployed contract.
func NewContractPreconfServiceManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractPreconfServiceManagerFilterer, error) {
	contract, err := bindContractPreconfServiceManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractPreconfServiceManagerFilterer{contract: contract}, nil
}

// bindContractPreconfServiceManager binds a generic wrapper to an already deployed contract.
func bindContractPreconfServiceManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractPreconfServiceManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractPreconfServiceManager.Contract.ContractPreconfServiceManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractPreconfServiceManager.Contract.ContractPreconfServiceManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractPreconfServiceManager.Contract.ContractPreconfServiceManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractPreconfServiceManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractPreconfServiceManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractPreconfServiceManager.Contract.contract.Transact(opts, method, params...)
}

// AvsDirectory is a free data retrieval call binding the contract method 0x6b3aa72e.
//
// Solidity: function avsDirectory() view returns(address)
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerCaller) AvsDirectory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractPreconfServiceManager.contract.Call(opts, &out, "avsDirectory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AvsDirectory is a free data retrieval call binding the contract method 0x6b3aa72e.
//
// Solidity: function avsDirectory() view returns(address)
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerSession) AvsDirectory() (common.Address, error) {
	return _ContractPreconfServiceManager.Contract.AvsDirectory(&_ContractPreconfServiceManager.CallOpts)
}

// AvsDirectory is a free data retrieval call binding the contract method 0x6b3aa72e.
//
// Solidity: function avsDirectory() view returns(address)
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerCallerSession) AvsDirectory() (common.Address, error) {
	return _ContractPreconfServiceManager.Contract.AvsDirectory(&_ContractPreconfServiceManager.CallOpts)
}

// ChallengeManager is a free data retrieval call binding the contract method 0x023a96fe.
//
// Solidity: function challengeManager() view returns(address)
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerCaller) ChallengeManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractPreconfServiceManager.contract.Call(opts, &out, "challengeManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ChallengeManager is a free data retrieval call binding the contract method 0x023a96fe.
//
// Solidity: function challengeManager() view returns(address)
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerSession) ChallengeManager() (common.Address, error) {
	return _ContractPreconfServiceManager.Contract.ChallengeManager(&_ContractPreconfServiceManager.CallOpts)
}

// ChallengeManager is a free data retrieval call binding the contract method 0x023a96fe.
//
// Solidity: function challengeManager() view returns(address)
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerCallerSession) ChallengeManager() (common.Address, error) {
	return _ContractPreconfServiceManager.Contract.ChallengeManager(&_ContractPreconfServiceManager.CallOpts)
}

// GetOperatorRestakedStrategies is a free data retrieval call binding the contract method 0x33cfb7b7.
//
// Solidity: function getOperatorRestakedStrategies(address operator) view returns(address[])
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerCaller) GetOperatorRestakedStrategies(opts *bind.CallOpts, operator common.Address) ([]common.Address, error) {
	var out []interface{}
	err := _ContractPreconfServiceManager.contract.Call(opts, &out, "getOperatorRestakedStrategies", operator)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetOperatorRestakedStrategies is a free data retrieval call binding the contract method 0x33cfb7b7.
//
// Solidity: function getOperatorRestakedStrategies(address operator) view returns(address[])
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerSession) GetOperatorRestakedStrategies(operator common.Address) ([]common.Address, error) {
	return _ContractPreconfServiceManager.Contract.GetOperatorRestakedStrategies(&_ContractPreconfServiceManager.CallOpts, operator)
}

// GetOperatorRestakedStrategies is a free data retrieval call binding the contract method 0x33cfb7b7.
//
// Solidity: function getOperatorRestakedStrategies(address operator) view returns(address[])
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerCallerSession) GetOperatorRestakedStrategies(operator common.Address) ([]common.Address, error) {
	return _ContractPreconfServiceManager.Contract.GetOperatorRestakedStrategies(&_ContractPreconfServiceManager.CallOpts, operator)
}

// GetRestakeableStrategies is a free data retrieval call binding the contract method 0xe481af9d.
//
// Solidity: function getRestakeableStrategies() view returns(address[])
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerCaller) GetRestakeableStrategies(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _ContractPreconfServiceManager.contract.Call(opts, &out, "getRestakeableStrategies")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetRestakeableStrategies is a free data retrieval call binding the contract method 0xe481af9d.
//
// Solidity: function getRestakeableStrategies() view returns(address[])
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerSession) GetRestakeableStrategies() ([]common.Address, error) {
	return _ContractPreconfServiceManager.Contract.GetRestakeableStrategies(&_ContractPreconfServiceManager.CallOpts)
}

// GetRestakeableStrategies is a free data retrieval call binding the contract method 0xe481af9d.
//
// Solidity: function getRestakeableStrategies() view returns(address[])
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerCallerSession) GetRestakeableStrategies() ([]common.Address, error) {
	return _ContractPreconfServiceManager.Contract.GetRestakeableStrategies(&_ContractPreconfServiceManager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractPreconfServiceManager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerSession) Owner() (common.Address, error) {
	return _ContractPreconfServiceManager.Contract.Owner(&_ContractPreconfServiceManager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerCallerSession) Owner() (common.Address, error) {
	return _ContractPreconfServiceManager.Contract.Owner(&_ContractPreconfServiceManager.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5ac86ab7.
//
// Solidity: function paused(uint8 index) view returns(bool)
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerCaller) Paused(opts *bind.CallOpts, index uint8) (bool, error) {
	var out []interface{}
	err := _ContractPreconfServiceManager.contract.Call(opts, &out, "paused", index)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5ac86ab7.
//
// Solidity: function paused(uint8 index) view returns(bool)
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerSession) Paused(index uint8) (bool, error) {
	return _ContractPreconfServiceManager.Contract.Paused(&_ContractPreconfServiceManager.CallOpts, index)
}

// Paused is a free data retrieval call binding the contract method 0x5ac86ab7.
//
// Solidity: function paused(uint8 index) view returns(bool)
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerCallerSession) Paused(index uint8) (bool, error) {
	return _ContractPreconfServiceManager.Contract.Paused(&_ContractPreconfServiceManager.CallOpts, index)
}

// Paused0 is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(uint256)
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerCaller) Paused0(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ContractPreconfServiceManager.contract.Call(opts, &out, "paused0")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Paused0 is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(uint256)
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerSession) Paused0() (*big.Int, error) {
	return _ContractPreconfServiceManager.Contract.Paused0(&_ContractPreconfServiceManager.CallOpts)
}

// Paused0 is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(uint256)
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerCallerSession) Paused0() (*big.Int, error) {
	return _ContractPreconfServiceManager.Contract.Paused0(&_ContractPreconfServiceManager.CallOpts)
}

// PauserRegistry is a free data retrieval call binding the contract method 0x886f1195.
//
// Solidity: function pauserRegistry() view returns(address)
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerCaller) PauserRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractPreconfServiceManager.contract.Call(opts, &out, "pauserRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PauserRegistry is a free data retrieval call binding the contract method 0x886f1195.
//
// Solidity: function pauserRegistry() view returns(address)
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerSession) PauserRegistry() (common.Address, error) {
	return _ContractPreconfServiceManager.Contract.PauserRegistry(&_ContractPreconfServiceManager.CallOpts)
}

// PauserRegistry is a free data retrieval call binding the contract method 0x886f1195.
//
// Solidity: function pauserRegistry() view returns(address)
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerCallerSession) PauserRegistry() (common.Address, error) {
	return _ContractPreconfServiceManager.Contract.PauserRegistry(&_ContractPreconfServiceManager.CallOpts)
}

// DeregisterOperatorFromAVS is a paid mutator transaction binding the contract method 0xa364f4da.
//
// Solidity: function deregisterOperatorFromAVS(address operator) returns()
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerTransactor) DeregisterOperatorFromAVS(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _ContractPreconfServiceManager.contract.Transact(opts, "deregisterOperatorFromAVS", operator)
}

// DeregisterOperatorFromAVS is a paid mutator transaction binding the contract method 0xa364f4da.
//
// Solidity: function deregisterOperatorFromAVS(address operator) returns()
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerSession) DeregisterOperatorFromAVS(operator common.Address) (*types.Transaction, error) {
	return _ContractPreconfServiceManager.Contract.DeregisterOperatorFromAVS(&_ContractPreconfServiceManager.TransactOpts, operator)
}

// DeregisterOperatorFromAVS is a paid mutator transaction binding the contract method 0xa364f4da.
//
// Solidity: function deregisterOperatorFromAVS(address operator) returns()
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerTransactorSession) DeregisterOperatorFromAVS(operator common.Address) (*types.Transaction, error) {
	return _ContractPreconfServiceManager.Contract.DeregisterOperatorFromAVS(&_ContractPreconfServiceManager.TransactOpts, operator)
}

// FreezeOperator is a paid mutator transaction binding the contract method 0x38c8ee64.
//
// Solidity: function freezeOperator(address operatorAddr) returns()
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerTransactor) FreezeOperator(opts *bind.TransactOpts, operatorAddr common.Address) (*types.Transaction, error) {
	return _ContractPreconfServiceManager.contract.Transact(opts, "freezeOperator", operatorAddr)
}

// FreezeOperator is a paid mutator transaction binding the contract method 0x38c8ee64.
//
// Solidity: function freezeOperator(address operatorAddr) returns()
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerSession) FreezeOperator(operatorAddr common.Address) (*types.Transaction, error) {
	return _ContractPreconfServiceManager.Contract.FreezeOperator(&_ContractPreconfServiceManager.TransactOpts, operatorAddr)
}

// FreezeOperator is a paid mutator transaction binding the contract method 0x38c8ee64.
//
// Solidity: function freezeOperator(address operatorAddr) returns()
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerTransactorSession) FreezeOperator(operatorAddr common.Address) (*types.Transaction, error) {
	return _ContractPreconfServiceManager.Contract.FreezeOperator(&_ContractPreconfServiceManager.TransactOpts, operatorAddr)
}

// Initialize is a paid mutator transaction binding the contract method 0x358394d8.
//
// Solidity: function initialize(address _pauserRegistry, uint256 _initialPausedStatus, address _initialOwner, address _challengeManager) returns()
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerTransactor) Initialize(opts *bind.TransactOpts, _pauserRegistry common.Address, _initialPausedStatus *big.Int, _initialOwner common.Address, _challengeManager common.Address) (*types.Transaction, error) {
	return _ContractPreconfServiceManager.contract.Transact(opts, "initialize", _pauserRegistry, _initialPausedStatus, _initialOwner, _challengeManager)
}

// Initialize is a paid mutator transaction binding the contract method 0x358394d8.
//
// Solidity: function initialize(address _pauserRegistry, uint256 _initialPausedStatus, address _initialOwner, address _challengeManager) returns()
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerSession) Initialize(_pauserRegistry common.Address, _initialPausedStatus *big.Int, _initialOwner common.Address, _challengeManager common.Address) (*types.Transaction, error) {
	return _ContractPreconfServiceManager.Contract.Initialize(&_ContractPreconfServiceManager.TransactOpts, _pauserRegistry, _initialPausedStatus, _initialOwner, _challengeManager)
}

// Initialize is a paid mutator transaction binding the contract method 0x358394d8.
//
// Solidity: function initialize(address _pauserRegistry, uint256 _initialPausedStatus, address _initialOwner, address _challengeManager) returns()
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerTransactorSession) Initialize(_pauserRegistry common.Address, _initialPausedStatus *big.Int, _initialOwner common.Address, _challengeManager common.Address) (*types.Transaction, error) {
	return _ContractPreconfServiceManager.Contract.Initialize(&_ContractPreconfServiceManager.TransactOpts, _pauserRegistry, _initialPausedStatus, _initialOwner, _challengeManager)
}

// Pause is a paid mutator transaction binding the contract method 0x136439dd.
//
// Solidity: function pause(uint256 newPausedStatus) returns()
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerTransactor) Pause(opts *bind.TransactOpts, newPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractPreconfServiceManager.contract.Transact(opts, "pause", newPausedStatus)
}

// Pause is a paid mutator transaction binding the contract method 0x136439dd.
//
// Solidity: function pause(uint256 newPausedStatus) returns()
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerSession) Pause(newPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractPreconfServiceManager.Contract.Pause(&_ContractPreconfServiceManager.TransactOpts, newPausedStatus)
}

// Pause is a paid mutator transaction binding the contract method 0x136439dd.
//
// Solidity: function pause(uint256 newPausedStatus) returns()
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerTransactorSession) Pause(newPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractPreconfServiceManager.Contract.Pause(&_ContractPreconfServiceManager.TransactOpts, newPausedStatus)
}

// PauseAll is a paid mutator transaction binding the contract method 0x595c6a67.
//
// Solidity: function pauseAll() returns()
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerTransactor) PauseAll(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractPreconfServiceManager.contract.Transact(opts, "pauseAll")
}

// PauseAll is a paid mutator transaction binding the contract method 0x595c6a67.
//
// Solidity: function pauseAll() returns()
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerSession) PauseAll() (*types.Transaction, error) {
	return _ContractPreconfServiceManager.Contract.PauseAll(&_ContractPreconfServiceManager.TransactOpts)
}

// PauseAll is a paid mutator transaction binding the contract method 0x595c6a67.
//
// Solidity: function pauseAll() returns()
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerTransactorSession) PauseAll() (*types.Transaction, error) {
	return _ContractPreconfServiceManager.Contract.PauseAll(&_ContractPreconfServiceManager.TransactOpts)
}

// RegisterOperatorToAVS is a paid mutator transaction binding the contract method 0x9926ee7d.
//
// Solidity: function registerOperatorToAVS(address operator, (bytes,bytes32,uint256) operatorSignature) returns()
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerTransactor) RegisterOperatorToAVS(opts *bind.TransactOpts, operator common.Address, operatorSignature ISignatureUtilsSignatureWithSaltAndExpiry) (*types.Transaction, error) {
	return _ContractPreconfServiceManager.contract.Transact(opts, "registerOperatorToAVS", operator, operatorSignature)
}

// RegisterOperatorToAVS is a paid mutator transaction binding the contract method 0x9926ee7d.
//
// Solidity: function registerOperatorToAVS(address operator, (bytes,bytes32,uint256) operatorSignature) returns()
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerSession) RegisterOperatorToAVS(operator common.Address, operatorSignature ISignatureUtilsSignatureWithSaltAndExpiry) (*types.Transaction, error) {
	return _ContractPreconfServiceManager.Contract.RegisterOperatorToAVS(&_ContractPreconfServiceManager.TransactOpts, operator, operatorSignature)
}

// RegisterOperatorToAVS is a paid mutator transaction binding the contract method 0x9926ee7d.
//
// Solidity: function registerOperatorToAVS(address operator, (bytes,bytes32,uint256) operatorSignature) returns()
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerTransactorSession) RegisterOperatorToAVS(operator common.Address, operatorSignature ISignatureUtilsSignatureWithSaltAndExpiry) (*types.Transaction, error) {
	return _ContractPreconfServiceManager.Contract.RegisterOperatorToAVS(&_ContractPreconfServiceManager.TransactOpts, operator, operatorSignature)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractPreconfServiceManager.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerSession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractPreconfServiceManager.Contract.RenounceOwnership(&_ContractPreconfServiceManager.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractPreconfServiceManager.Contract.RenounceOwnership(&_ContractPreconfServiceManager.TransactOpts)
}

// SetMetadataURI is a paid mutator transaction binding the contract method 0x750521f5.
//
// Solidity: function setMetadataURI(string _metadataURI) returns()
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerTransactor) SetMetadataURI(opts *bind.TransactOpts, _metadataURI string) (*types.Transaction, error) {
	return _ContractPreconfServiceManager.contract.Transact(opts, "setMetadataURI", _metadataURI)
}

// SetMetadataURI is a paid mutator transaction binding the contract method 0x750521f5.
//
// Solidity: function setMetadataURI(string _metadataURI) returns()
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerSession) SetMetadataURI(_metadataURI string) (*types.Transaction, error) {
	return _ContractPreconfServiceManager.Contract.SetMetadataURI(&_ContractPreconfServiceManager.TransactOpts, _metadataURI)
}

// SetMetadataURI is a paid mutator transaction binding the contract method 0x750521f5.
//
// Solidity: function setMetadataURI(string _metadataURI) returns()
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerTransactorSession) SetMetadataURI(_metadataURI string) (*types.Transaction, error) {
	return _ContractPreconfServiceManager.Contract.SetMetadataURI(&_ContractPreconfServiceManager.TransactOpts, _metadataURI)
}

// SetPauserRegistry is a paid mutator transaction binding the contract method 0x10d67a2f.
//
// Solidity: function setPauserRegistry(address newPauserRegistry) returns()
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerTransactor) SetPauserRegistry(opts *bind.TransactOpts, newPauserRegistry common.Address) (*types.Transaction, error) {
	return _ContractPreconfServiceManager.contract.Transact(opts, "setPauserRegistry", newPauserRegistry)
}

// SetPauserRegistry is a paid mutator transaction binding the contract method 0x10d67a2f.
//
// Solidity: function setPauserRegistry(address newPauserRegistry) returns()
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerSession) SetPauserRegistry(newPauserRegistry common.Address) (*types.Transaction, error) {
	return _ContractPreconfServiceManager.Contract.SetPauserRegistry(&_ContractPreconfServiceManager.TransactOpts, newPauserRegistry)
}

// SetPauserRegistry is a paid mutator transaction binding the contract method 0x10d67a2f.
//
// Solidity: function setPauserRegistry(address newPauserRegistry) returns()
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerTransactorSession) SetPauserRegistry(newPauserRegistry common.Address) (*types.Transaction, error) {
	return _ContractPreconfServiceManager.Contract.SetPauserRegistry(&_ContractPreconfServiceManager.TransactOpts, newPauserRegistry)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ContractPreconfServiceManager.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractPreconfServiceManager.Contract.TransferOwnership(&_ContractPreconfServiceManager.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractPreconfServiceManager.Contract.TransferOwnership(&_ContractPreconfServiceManager.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0xfabc1cbc.
//
// Solidity: function unpause(uint256 newPausedStatus) returns()
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerTransactor) Unpause(opts *bind.TransactOpts, newPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractPreconfServiceManager.contract.Transact(opts, "unpause", newPausedStatus)
}

// Unpause is a paid mutator transaction binding the contract method 0xfabc1cbc.
//
// Solidity: function unpause(uint256 newPausedStatus) returns()
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerSession) Unpause(newPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractPreconfServiceManager.Contract.Unpause(&_ContractPreconfServiceManager.TransactOpts, newPausedStatus)
}

// Unpause is a paid mutator transaction binding the contract method 0xfabc1cbc.
//
// Solidity: function unpause(uint256 newPausedStatus) returns()
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerTransactorSession) Unpause(newPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractPreconfServiceManager.Contract.Unpause(&_ContractPreconfServiceManager.TransactOpts, newPausedStatus)
}

// ContractPreconfServiceManagerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ContractPreconfServiceManager contract.
type ContractPreconfServiceManagerInitializedIterator struct {
	Event *ContractPreconfServiceManagerInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractPreconfServiceManagerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPreconfServiceManagerInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractPreconfServiceManagerInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractPreconfServiceManagerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPreconfServiceManagerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPreconfServiceManagerInitialized represents a Initialized event raised by the ContractPreconfServiceManager contract.
type ContractPreconfServiceManagerInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerFilterer) FilterInitialized(opts *bind.FilterOpts) (*ContractPreconfServiceManagerInitializedIterator, error) {

	logs, sub, err := _ContractPreconfServiceManager.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ContractPreconfServiceManagerInitializedIterator{contract: _ContractPreconfServiceManager.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ContractPreconfServiceManagerInitialized) (event.Subscription, error) {

	logs, sub, err := _ContractPreconfServiceManager.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPreconfServiceManagerInitialized)
				if err := _ContractPreconfServiceManager.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerFilterer) ParseInitialized(log types.Log) (*ContractPreconfServiceManagerInitialized, error) {
	event := new(ContractPreconfServiceManagerInitialized)
	if err := _ContractPreconfServiceManager.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractPreconfServiceManagerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ContractPreconfServiceManager contract.
type ContractPreconfServiceManagerOwnershipTransferredIterator struct {
	Event *ContractPreconfServiceManagerOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractPreconfServiceManagerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPreconfServiceManagerOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractPreconfServiceManagerOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractPreconfServiceManagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPreconfServiceManagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPreconfServiceManagerOwnershipTransferred represents a OwnershipTransferred event raised by the ContractPreconfServiceManager contract.
type ContractPreconfServiceManagerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ContractPreconfServiceManagerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractPreconfServiceManager.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ContractPreconfServiceManagerOwnershipTransferredIterator{contract: _ContractPreconfServiceManager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ContractPreconfServiceManagerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractPreconfServiceManager.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPreconfServiceManagerOwnershipTransferred)
				if err := _ContractPreconfServiceManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerFilterer) ParseOwnershipTransferred(log types.Log) (*ContractPreconfServiceManagerOwnershipTransferred, error) {
	event := new(ContractPreconfServiceManagerOwnershipTransferred)
	if err := _ContractPreconfServiceManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractPreconfServiceManagerPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the ContractPreconfServiceManager contract.
type ContractPreconfServiceManagerPausedIterator struct {
	Event *ContractPreconfServiceManagerPaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractPreconfServiceManagerPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPreconfServiceManagerPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractPreconfServiceManagerPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractPreconfServiceManagerPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPreconfServiceManagerPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPreconfServiceManagerPaused represents a Paused event raised by the ContractPreconfServiceManager contract.
type ContractPreconfServiceManagerPaused struct {
	Account         common.Address
	NewPausedStatus *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0xab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d.
//
// Solidity: event Paused(address indexed account, uint256 newPausedStatus)
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerFilterer) FilterPaused(opts *bind.FilterOpts, account []common.Address) (*ContractPreconfServiceManagerPausedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ContractPreconfServiceManager.contract.FilterLogs(opts, "Paused", accountRule)
	if err != nil {
		return nil, err
	}
	return &ContractPreconfServiceManagerPausedIterator{contract: _ContractPreconfServiceManager.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0xab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d.
//
// Solidity: event Paused(address indexed account, uint256 newPausedStatus)
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *ContractPreconfServiceManagerPaused, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ContractPreconfServiceManager.contract.WatchLogs(opts, "Paused", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPreconfServiceManagerPaused)
				if err := _ContractPreconfServiceManager.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePaused is a log parse operation binding the contract event 0xab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d.
//
// Solidity: event Paused(address indexed account, uint256 newPausedStatus)
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerFilterer) ParsePaused(log types.Log) (*ContractPreconfServiceManagerPaused, error) {
	event := new(ContractPreconfServiceManagerPaused)
	if err := _ContractPreconfServiceManager.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractPreconfServiceManagerPauserRegistrySetIterator is returned from FilterPauserRegistrySet and is used to iterate over the raw logs and unpacked data for PauserRegistrySet events raised by the ContractPreconfServiceManager contract.
type ContractPreconfServiceManagerPauserRegistrySetIterator struct {
	Event *ContractPreconfServiceManagerPauserRegistrySet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractPreconfServiceManagerPauserRegistrySetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPreconfServiceManagerPauserRegistrySet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractPreconfServiceManagerPauserRegistrySet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractPreconfServiceManagerPauserRegistrySetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPreconfServiceManagerPauserRegistrySetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPreconfServiceManagerPauserRegistrySet represents a PauserRegistrySet event raised by the ContractPreconfServiceManager contract.
type ContractPreconfServiceManagerPauserRegistrySet struct {
	PauserRegistry    common.Address
	NewPauserRegistry common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterPauserRegistrySet is a free log retrieval operation binding the contract event 0x6e9fcd539896fca60e8b0f01dd580233e48a6b0f7df013b89ba7f565869acdb6.
//
// Solidity: event PauserRegistrySet(address pauserRegistry, address newPauserRegistry)
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerFilterer) FilterPauserRegistrySet(opts *bind.FilterOpts) (*ContractPreconfServiceManagerPauserRegistrySetIterator, error) {

	logs, sub, err := _ContractPreconfServiceManager.contract.FilterLogs(opts, "PauserRegistrySet")
	if err != nil {
		return nil, err
	}
	return &ContractPreconfServiceManagerPauserRegistrySetIterator{contract: _ContractPreconfServiceManager.contract, event: "PauserRegistrySet", logs: logs, sub: sub}, nil
}

// WatchPauserRegistrySet is a free log subscription operation binding the contract event 0x6e9fcd539896fca60e8b0f01dd580233e48a6b0f7df013b89ba7f565869acdb6.
//
// Solidity: event PauserRegistrySet(address pauserRegistry, address newPauserRegistry)
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerFilterer) WatchPauserRegistrySet(opts *bind.WatchOpts, sink chan<- *ContractPreconfServiceManagerPauserRegistrySet) (event.Subscription, error) {

	logs, sub, err := _ContractPreconfServiceManager.contract.WatchLogs(opts, "PauserRegistrySet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPreconfServiceManagerPauserRegistrySet)
				if err := _ContractPreconfServiceManager.contract.UnpackLog(event, "PauserRegistrySet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePauserRegistrySet is a log parse operation binding the contract event 0x6e9fcd539896fca60e8b0f01dd580233e48a6b0f7df013b89ba7f565869acdb6.
//
// Solidity: event PauserRegistrySet(address pauserRegistry, address newPauserRegistry)
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerFilterer) ParsePauserRegistrySet(log types.Log) (*ContractPreconfServiceManagerPauserRegistrySet, error) {
	event := new(ContractPreconfServiceManagerPauserRegistrySet)
	if err := _ContractPreconfServiceManager.contract.UnpackLog(event, "PauserRegistrySet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractPreconfServiceManagerUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the ContractPreconfServiceManager contract.
type ContractPreconfServiceManagerUnpausedIterator struct {
	Event *ContractPreconfServiceManagerUnpaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractPreconfServiceManagerUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPreconfServiceManagerUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractPreconfServiceManagerUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractPreconfServiceManagerUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPreconfServiceManagerUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPreconfServiceManagerUnpaused represents a Unpaused event raised by the ContractPreconfServiceManager contract.
type ContractPreconfServiceManagerUnpaused struct {
	Account         common.Address
	NewPausedStatus *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x3582d1828e26bf56bd801502bc021ac0bc8afb57c826e4986b45593c8fad389c.
//
// Solidity: event Unpaused(address indexed account, uint256 newPausedStatus)
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerFilterer) FilterUnpaused(opts *bind.FilterOpts, account []common.Address) (*ContractPreconfServiceManagerUnpausedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ContractPreconfServiceManager.contract.FilterLogs(opts, "Unpaused", accountRule)
	if err != nil {
		return nil, err
	}
	return &ContractPreconfServiceManagerUnpausedIterator{contract: _ContractPreconfServiceManager.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x3582d1828e26bf56bd801502bc021ac0bc8afb57c826e4986b45593c8fad389c.
//
// Solidity: event Unpaused(address indexed account, uint256 newPausedStatus)
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *ContractPreconfServiceManagerUnpaused, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ContractPreconfServiceManager.contract.WatchLogs(opts, "Unpaused", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPreconfServiceManagerUnpaused)
				if err := _ContractPreconfServiceManager.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUnpaused is a log parse operation binding the contract event 0x3582d1828e26bf56bd801502bc021ac0bc8afb57c826e4986b45593c8fad389c.
//
// Solidity: event Unpaused(address indexed account, uint256 newPausedStatus)
func (_ContractPreconfServiceManager *ContractPreconfServiceManagerFilterer) ParseUnpaused(log types.Log) (*ContractPreconfServiceManagerUnpaused, error) {
	event := new(ContractPreconfServiceManagerUnpaused)
	if err := _ContractPreconfServiceManager.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
