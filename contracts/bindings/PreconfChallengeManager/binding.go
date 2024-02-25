// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractPreconfChallengeManager

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

// PreconfChallengeManagerPreconf is an auto generated low-level Go binding around an user-defined struct.
type PreconfChallengeManagerPreconf struct {
	Hashes      [][32]byte
	BlockTarget *big.Int
}

// ContractPreconfChallengeManagerMetaData contains all meta data concerning the ContractPreconfChallengeManager contract.
var ContractPreconfChallengeManagerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"CHALLENGE_WINDOW_BLOCK\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"challenges\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_serviceManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_prover\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"prover\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"raiseAndResolveChallenge\",\"inputs\":[{\"name\":\"preconf\",\"type\":\"tuple\",\"internalType\":\"structPreconfChallengeManager.Preconf\",\"components\":[{\"name\":\"hashes\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"blockTarget\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"preconfSigned\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"proof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"serviceManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Challenged\",\"inputs\":[{\"name\":\"preconf\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structPreconfChallengeManager.Preconf\",\"components\":[{\"name\":\"hashes\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"blockTarget\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"challenger\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AlreadyChallenged\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EmptyBundle\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotReached\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"WindowExpired\",\"inputs\":[]}]",
	Bin: "0x608060405234801561001057600080fd5b5061001961001e565b6100de565b600054610100900460ff161561008a5760405162461bcd60e51b815260206004820152602760248201527f496e697469616c697a61626c653a20636f6e747261637420697320696e697469604482015266616c697a696e6760c81b606482015260840160405180910390fd5b60005460ff90811610156100dc576000805460ff191660ff9081179091556040519081527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b565b611046806100ed6000396000f3fe608060405234801561001057600080fd5b50600436106100935760003560e01c8063715018a611610066578063715018a61461010d5780638da5cb5b14610115578063c0c53b8b14610126578063c1e69b6614610139578063f2fde38b1461016c57600080fd5b80631b4f683f1461009857806332a8f30f146100ad5780633998fdd3146100dd57806348e30cdf146100f0575b600080fd5b6100ab6100a6366004610b73565b61017f565b005b609a546100c0906001600160a01b031681565b6040516001600160a01b0390911681526020015b60405180910390f35b6099546100c0906001600160a01b031681565b6100f8606481565b60405163ffffffff90911681526020016100d4565b6100ab610400565b6033546001600160a01b03166100c0565b6100ab610134366004610c25565b610414565b61015c610147366004610c70565b609b6020526000908152604090205460ff1681565b60405190151581526020016100d4565b6100ab61017a366004610c89565b6105c2565b84602001354310156101a457604051635203092960e11b815260040160405180910390fd5b6101b360646020870135610cad565b4311156101d357604051633ddfdb7f60e11b815260040160405180910390fd5b6101dd8580610cd3565b151590506101fe57604051631563113f60e21b815260040160405180910390fd5b6000856040516020016102119190610da9565b60408051601f1981840301815291815281516020928301206000818152609b90935291205490915060ff161561025a5760405163f1082a9360e01b815260040160405180910390fd5b609a5460405163e19108c360e01b81526000916001600160a01b03169063e19108c39061028f90879087908690600401610dbc565b6000604051808303816000875af11580156102ae573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526102d69190810190610e64565b9050600081604001518060200190518101906102f29190610f63565b90506001811515146103a35760006103408489898080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061063b92505050565b609954604051630e323b9960e21b81526001600160a01b0380841660048301529293509116906338c8ee6490602401600060405180830381600087803b15801561038957600080fd5b505af115801561039d573d6000803e3d6000fd5b50505050505b6000838152609b602052604090819020805460ff19166001179055517f9955dcf204285cdb270c001137aafffc8da100b0f75a21ace093b1bb92b49dc0906103ee908a903390610f85565b60405180910390a15050505050505050565b61040861065f565b61041260006106b9565b565b600054610100900460ff16158080156104345750600054600160ff909116105b8061044e5750303b15801561044e575060005460ff166001145b6104b65760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b60648201526084015b60405180910390fd5b6000805460ff1916600117905580156104d9576000805461ff0019166101001790555b6104e161070b565b61053d6040518060400160405280601781526020017f507265636f6e664368616c6c656e67654d616e61676572000000000000000000815250604051806040016040528060058152602001640302e312e360dc1b81525061073a565b610546846106b9565b609980546001600160a01b038086166001600160a01b031992831617909255609a80549285169290911691909117905580156105bc576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b50505050565b6105ca61065f565b6001600160a01b03811661062f5760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b60648201526084016104ad565b610638816106b9565b50565b600080600061064a858561076f565b91509150610657816107df565b509392505050565b6033546001600160a01b031633146104125760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016104ad565b603380546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b600054610100900460ff166107325760405162461bcd60e51b81526004016104ad90610faf565b61041261099a565b600054610100900460ff166107615760405162461bcd60e51b81526004016104ad90610faf565b61076b82826109ca565b5050565b6000808251604114156107a65760208301516040840151606085015160001a61079a87828585610a0b565b945094505050506107d8565b8251604014156107d057602083015160408401516107c5868383610af8565b9350935050506107d8565b506000905060025b9250929050565b60008160048111156107f3576107f3610ffa565b14156107fc5750565b600181600481111561081057610810610ffa565b141561085e5760405162461bcd60e51b815260206004820152601860248201527f45434453413a20696e76616c6964207369676e6174757265000000000000000060448201526064016104ad565b600281600481111561087257610872610ffa565b14156108c05760405162461bcd60e51b815260206004820152601f60248201527f45434453413a20696e76616c6964207369676e6174757265206c656e6774680060448201526064016104ad565b60038160048111156108d4576108d4610ffa565b141561092d5760405162461bcd60e51b815260206004820152602260248201527f45434453413a20696e76616c6964207369676e6174757265202773272076616c604482015261756560f01b60648201526084016104ad565b600481600481111561094157610941610ffa565b14156106385760405162461bcd60e51b815260206004820152602260248201527f45434453413a20696e76616c6964207369676e6174757265202776272076616c604482015261756560f01b60648201526084016104ad565b600054610100900460ff166109c15760405162461bcd60e51b81526004016104ad90610faf565b610412336106b9565b600054610100900460ff166109f15760405162461bcd60e51b81526004016104ad90610faf565b815160209283012081519190920120606591909155606655565b6000807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a0831115610a425750600090506003610aef565b8460ff16601b14158015610a5a57508460ff16601c14155b15610a6b5750600090506004610aef565b6040805160008082526020820180845289905260ff881692820192909252606081018690526080810185905260019060a0016020604051602081039080840390855afa158015610abf573d6000803e3d6000fd5b5050604051601f1901519150506001600160a01b038116610ae857600060019250925050610aef565b9150600090505b94509492505050565b6000806001600160ff1b03831681610b1560ff86901c601b610cad565b9050610b2387828885610a0b565b935093505050935093915050565b60008083601f840112610b4357600080fd5b50813567ffffffffffffffff811115610b5b57600080fd5b6020830191508360208285010111156107d857600080fd5b600080600080600060608688031215610b8b57600080fd5b853567ffffffffffffffff80821115610ba357600080fd5b908701906040828a031215610bb757600080fd5b90955060208701359080821115610bcd57600080fd5b610bd989838a01610b31565b90965094506040880135915080821115610bf257600080fd5b50610bff88828901610b31565b969995985093965092949392505050565b6001600160a01b038116811461063857600080fd5b600080600060608486031215610c3a57600080fd5b8335610c4581610c10565b92506020840135610c5581610c10565b91506040840135610c6581610c10565b809150509250925092565b600060208284031215610c8257600080fd5b5035919050565b600060208284031215610c9b57600080fd5b8135610ca681610c10565b9392505050565b60008219821115610cce57634e487b7160e01b600052601160045260246000fd5b500190565b6000808335601e19843603018112610cea57600080fd5b83018035915067ffffffffffffffff821115610d0557600080fd5b6020019150600581901b36038213156107d857600080fd5b60008135601e19833603018112610d3357600080fd5b8201803567ffffffffffffffff811115610d4c57600080fd5b8060051b803603851315610d5f57600080fd5b604080875286018290526001600160fb1b03821115610d7d57600080fd5b806020840160608801376060818701019250505060008152602083013560208501528091505092915050565b602081526000610ca66020830184610d1d565b6040815282604082015282846060830137600060608483018101919091529115156020820152601f909201601f191690910101919050565b634e487b7160e01b600052604160045260246000fd5b6040516060810167ffffffffffffffff81118282101715610e2d57610e2d610df4565b60405290565b604051601f8201601f1916810167ffffffffffffffff81118282101715610e5c57610e5c610df4565b604052919050565b60006020808385031215610e7757600080fd5b825167ffffffffffffffff80821115610e8f57600080fd5b9084019060608287031215610ea357600080fd5b610eab610e0a565b8251610eb681610c10565b81528284015184820152604083015182811115610ed257600080fd5b80840193505086601f840112610ee757600080fd5b825182811115610ef957610ef9610df4565b610f0b601f8201601f19168601610e33565b92508083528785828601011115610f2157600080fd5b60005b81811015610f3f578481018601518482018701528501610f24565b81811115610f505760008683860101525b5050604081019190915295945050505050565b600060208284031215610f7557600080fd5b81518015158114610ca657600080fd5b604081526000610f986040830185610d1d565b905060018060a01b03831660208301529392505050565b6020808252602b908201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960408201526a6e697469616c697a696e6760a81b606082015260800190565b634e487b7160e01b600052602160045260246000fdfea2646970667358221220340a989c59206123e766cbdf70865ce707f29516fb9c10a9eef25ca6cb7687e564736f6c634300080c0033",
}

// ContractPreconfChallengeManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractPreconfChallengeManagerMetaData.ABI instead.
var ContractPreconfChallengeManagerABI = ContractPreconfChallengeManagerMetaData.ABI

// ContractPreconfChallengeManagerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractPreconfChallengeManagerMetaData.Bin instead.
var ContractPreconfChallengeManagerBin = ContractPreconfChallengeManagerMetaData.Bin

// DeployContractPreconfChallengeManager deploys a new Ethereum contract, binding an instance of ContractPreconfChallengeManager to it.
func DeployContractPreconfChallengeManager(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ContractPreconfChallengeManager, error) {
	parsed, err := ContractPreconfChallengeManagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractPreconfChallengeManagerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractPreconfChallengeManager{ContractPreconfChallengeManagerCaller: ContractPreconfChallengeManagerCaller{contract: contract}, ContractPreconfChallengeManagerTransactor: ContractPreconfChallengeManagerTransactor{contract: contract}, ContractPreconfChallengeManagerFilterer: ContractPreconfChallengeManagerFilterer{contract: contract}}, nil
}

// ContractPreconfChallengeManager is an auto generated Go binding around an Ethereum contract.
type ContractPreconfChallengeManager struct {
	ContractPreconfChallengeManagerCaller     // Read-only binding to the contract
	ContractPreconfChallengeManagerTransactor // Write-only binding to the contract
	ContractPreconfChallengeManagerFilterer   // Log filterer for contract events
}

// ContractPreconfChallengeManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractPreconfChallengeManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractPreconfChallengeManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractPreconfChallengeManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractPreconfChallengeManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractPreconfChallengeManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractPreconfChallengeManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractPreconfChallengeManagerSession struct {
	Contract     *ContractPreconfChallengeManager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                    // Call options to use throughout this session
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// ContractPreconfChallengeManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractPreconfChallengeManagerCallerSession struct {
	Contract *ContractPreconfChallengeManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                          // Call options to use throughout this session
}

// ContractPreconfChallengeManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractPreconfChallengeManagerTransactorSession struct {
	Contract     *ContractPreconfChallengeManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                          // Transaction auth options to use throughout this session
}

// ContractPreconfChallengeManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractPreconfChallengeManagerRaw struct {
	Contract *ContractPreconfChallengeManager // Generic contract binding to access the raw methods on
}

// ContractPreconfChallengeManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractPreconfChallengeManagerCallerRaw struct {
	Contract *ContractPreconfChallengeManagerCaller // Generic read-only contract binding to access the raw methods on
}

// ContractPreconfChallengeManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractPreconfChallengeManagerTransactorRaw struct {
	Contract *ContractPreconfChallengeManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractPreconfChallengeManager creates a new instance of ContractPreconfChallengeManager, bound to a specific deployed contract.
func NewContractPreconfChallengeManager(address common.Address, backend bind.ContractBackend) (*ContractPreconfChallengeManager, error) {
	contract, err := bindContractPreconfChallengeManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractPreconfChallengeManager{ContractPreconfChallengeManagerCaller: ContractPreconfChallengeManagerCaller{contract: contract}, ContractPreconfChallengeManagerTransactor: ContractPreconfChallengeManagerTransactor{contract: contract}, ContractPreconfChallengeManagerFilterer: ContractPreconfChallengeManagerFilterer{contract: contract}}, nil
}

// NewContractPreconfChallengeManagerCaller creates a new read-only instance of ContractPreconfChallengeManager, bound to a specific deployed contract.
func NewContractPreconfChallengeManagerCaller(address common.Address, caller bind.ContractCaller) (*ContractPreconfChallengeManagerCaller, error) {
	contract, err := bindContractPreconfChallengeManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractPreconfChallengeManagerCaller{contract: contract}, nil
}

// NewContractPreconfChallengeManagerTransactor creates a new write-only instance of ContractPreconfChallengeManager, bound to a specific deployed contract.
func NewContractPreconfChallengeManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractPreconfChallengeManagerTransactor, error) {
	contract, err := bindContractPreconfChallengeManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractPreconfChallengeManagerTransactor{contract: contract}, nil
}

// NewContractPreconfChallengeManagerFilterer creates a new log filterer instance of ContractPreconfChallengeManager, bound to a specific deployed contract.
func NewContractPreconfChallengeManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractPreconfChallengeManagerFilterer, error) {
	contract, err := bindContractPreconfChallengeManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractPreconfChallengeManagerFilterer{contract: contract}, nil
}

// bindContractPreconfChallengeManager binds a generic wrapper to an already deployed contract.
func bindContractPreconfChallengeManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractPreconfChallengeManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractPreconfChallengeManager.Contract.ContractPreconfChallengeManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractPreconfChallengeManager.Contract.ContractPreconfChallengeManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractPreconfChallengeManager.Contract.ContractPreconfChallengeManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractPreconfChallengeManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractPreconfChallengeManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractPreconfChallengeManager.Contract.contract.Transact(opts, method, params...)
}

// CHALLENGEWINDOWBLOCK is a free data retrieval call binding the contract method 0x48e30cdf.
//
// Solidity: function CHALLENGE_WINDOW_BLOCK() view returns(uint32)
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerCaller) CHALLENGEWINDOWBLOCK(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _ContractPreconfChallengeManager.contract.Call(opts, &out, "CHALLENGE_WINDOW_BLOCK")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// CHALLENGEWINDOWBLOCK is a free data retrieval call binding the contract method 0x48e30cdf.
//
// Solidity: function CHALLENGE_WINDOW_BLOCK() view returns(uint32)
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerSession) CHALLENGEWINDOWBLOCK() (uint32, error) {
	return _ContractPreconfChallengeManager.Contract.CHALLENGEWINDOWBLOCK(&_ContractPreconfChallengeManager.CallOpts)
}

// CHALLENGEWINDOWBLOCK is a free data retrieval call binding the contract method 0x48e30cdf.
//
// Solidity: function CHALLENGE_WINDOW_BLOCK() view returns(uint32)
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerCallerSession) CHALLENGEWINDOWBLOCK() (uint32, error) {
	return _ContractPreconfChallengeManager.Contract.CHALLENGEWINDOWBLOCK(&_ContractPreconfChallengeManager.CallOpts)
}

// Challenges is a free data retrieval call binding the contract method 0xc1e69b66.
//
// Solidity: function challenges(bytes32 ) view returns(bool)
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerCaller) Challenges(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _ContractPreconfChallengeManager.contract.Call(opts, &out, "challenges", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Challenges is a free data retrieval call binding the contract method 0xc1e69b66.
//
// Solidity: function challenges(bytes32 ) view returns(bool)
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerSession) Challenges(arg0 [32]byte) (bool, error) {
	return _ContractPreconfChallengeManager.Contract.Challenges(&_ContractPreconfChallengeManager.CallOpts, arg0)
}

// Challenges is a free data retrieval call binding the contract method 0xc1e69b66.
//
// Solidity: function challenges(bytes32 ) view returns(bool)
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerCallerSession) Challenges(arg0 [32]byte) (bool, error) {
	return _ContractPreconfChallengeManager.Contract.Challenges(&_ContractPreconfChallengeManager.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractPreconfChallengeManager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerSession) Owner() (common.Address, error) {
	return _ContractPreconfChallengeManager.Contract.Owner(&_ContractPreconfChallengeManager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerCallerSession) Owner() (common.Address, error) {
	return _ContractPreconfChallengeManager.Contract.Owner(&_ContractPreconfChallengeManager.CallOpts)
}

// Prover is a free data retrieval call binding the contract method 0x32a8f30f.
//
// Solidity: function prover() view returns(address)
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerCaller) Prover(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractPreconfChallengeManager.contract.Call(opts, &out, "prover")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Prover is a free data retrieval call binding the contract method 0x32a8f30f.
//
// Solidity: function prover() view returns(address)
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerSession) Prover() (common.Address, error) {
	return _ContractPreconfChallengeManager.Contract.Prover(&_ContractPreconfChallengeManager.CallOpts)
}

// Prover is a free data retrieval call binding the contract method 0x32a8f30f.
//
// Solidity: function prover() view returns(address)
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerCallerSession) Prover() (common.Address, error) {
	return _ContractPreconfChallengeManager.Contract.Prover(&_ContractPreconfChallengeManager.CallOpts)
}

// ServiceManager is a free data retrieval call binding the contract method 0x3998fdd3.
//
// Solidity: function serviceManager() view returns(address)
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerCaller) ServiceManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractPreconfChallengeManager.contract.Call(opts, &out, "serviceManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ServiceManager is a free data retrieval call binding the contract method 0x3998fdd3.
//
// Solidity: function serviceManager() view returns(address)
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerSession) ServiceManager() (common.Address, error) {
	return _ContractPreconfChallengeManager.Contract.ServiceManager(&_ContractPreconfChallengeManager.CallOpts)
}

// ServiceManager is a free data retrieval call binding the contract method 0x3998fdd3.
//
// Solidity: function serviceManager() view returns(address)
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerCallerSession) ServiceManager() (common.Address, error) {
	return _ContractPreconfChallengeManager.Contract.ServiceManager(&_ContractPreconfChallengeManager.CallOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address owner, address _serviceManager, address _prover) returns()
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerTransactor) Initialize(opts *bind.TransactOpts, owner common.Address, _serviceManager common.Address, _prover common.Address) (*types.Transaction, error) {
	return _ContractPreconfChallengeManager.contract.Transact(opts, "initialize", owner, _serviceManager, _prover)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address owner, address _serviceManager, address _prover) returns()
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerSession) Initialize(owner common.Address, _serviceManager common.Address, _prover common.Address) (*types.Transaction, error) {
	return _ContractPreconfChallengeManager.Contract.Initialize(&_ContractPreconfChallengeManager.TransactOpts, owner, _serviceManager, _prover)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address owner, address _serviceManager, address _prover) returns()
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerTransactorSession) Initialize(owner common.Address, _serviceManager common.Address, _prover common.Address) (*types.Transaction, error) {
	return _ContractPreconfChallengeManager.Contract.Initialize(&_ContractPreconfChallengeManager.TransactOpts, owner, _serviceManager, _prover)
}

// RaiseAndResolveChallenge is a paid mutator transaction binding the contract method 0x1b4f683f.
//
// Solidity: function raiseAndResolveChallenge((bytes32[],uint256) preconf, bytes preconfSigned, bytes proof) returns()
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerTransactor) RaiseAndResolveChallenge(opts *bind.TransactOpts, preconf PreconfChallengeManagerPreconf, preconfSigned []byte, proof []byte) (*types.Transaction, error) {
	return _ContractPreconfChallengeManager.contract.Transact(opts, "raiseAndResolveChallenge", preconf, preconfSigned, proof)
}

// RaiseAndResolveChallenge is a paid mutator transaction binding the contract method 0x1b4f683f.
//
// Solidity: function raiseAndResolveChallenge((bytes32[],uint256) preconf, bytes preconfSigned, bytes proof) returns()
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerSession) RaiseAndResolveChallenge(preconf PreconfChallengeManagerPreconf, preconfSigned []byte, proof []byte) (*types.Transaction, error) {
	return _ContractPreconfChallengeManager.Contract.RaiseAndResolveChallenge(&_ContractPreconfChallengeManager.TransactOpts, preconf, preconfSigned, proof)
}

// RaiseAndResolveChallenge is a paid mutator transaction binding the contract method 0x1b4f683f.
//
// Solidity: function raiseAndResolveChallenge((bytes32[],uint256) preconf, bytes preconfSigned, bytes proof) returns()
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerTransactorSession) RaiseAndResolveChallenge(preconf PreconfChallengeManagerPreconf, preconfSigned []byte, proof []byte) (*types.Transaction, error) {
	return _ContractPreconfChallengeManager.Contract.RaiseAndResolveChallenge(&_ContractPreconfChallengeManager.TransactOpts, preconf, preconfSigned, proof)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractPreconfChallengeManager.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerSession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractPreconfChallengeManager.Contract.RenounceOwnership(&_ContractPreconfChallengeManager.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractPreconfChallengeManager.Contract.RenounceOwnership(&_ContractPreconfChallengeManager.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ContractPreconfChallengeManager.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractPreconfChallengeManager.Contract.TransferOwnership(&_ContractPreconfChallengeManager.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractPreconfChallengeManager.Contract.TransferOwnership(&_ContractPreconfChallengeManager.TransactOpts, newOwner)
}

// ContractPreconfChallengeManagerChallengedIterator is returned from FilterChallenged and is used to iterate over the raw logs and unpacked data for Challenged events raised by the ContractPreconfChallengeManager contract.
type ContractPreconfChallengeManagerChallengedIterator struct {
	Event *ContractPreconfChallengeManagerChallenged // Event containing the contract specifics and raw log

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
func (it *ContractPreconfChallengeManagerChallengedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPreconfChallengeManagerChallenged)
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
		it.Event = new(ContractPreconfChallengeManagerChallenged)
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
func (it *ContractPreconfChallengeManagerChallengedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPreconfChallengeManagerChallengedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPreconfChallengeManagerChallenged represents a Challenged event raised by the ContractPreconfChallengeManager contract.
type ContractPreconfChallengeManagerChallenged struct {
	Preconf    PreconfChallengeManagerPreconf
	Challenger common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterChallenged is a free log retrieval operation binding the contract event 0x9955dcf204285cdb270c001137aafffc8da100b0f75a21ace093b1bb92b49dc0.
//
// Solidity: event Challenged((bytes32[],uint256) preconf, address challenger)
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerFilterer) FilterChallenged(opts *bind.FilterOpts) (*ContractPreconfChallengeManagerChallengedIterator, error) {

	logs, sub, err := _ContractPreconfChallengeManager.contract.FilterLogs(opts, "Challenged")
	if err != nil {
		return nil, err
	}
	return &ContractPreconfChallengeManagerChallengedIterator{contract: _ContractPreconfChallengeManager.contract, event: "Challenged", logs: logs, sub: sub}, nil
}

// WatchChallenged is a free log subscription operation binding the contract event 0x9955dcf204285cdb270c001137aafffc8da100b0f75a21ace093b1bb92b49dc0.
//
// Solidity: event Challenged((bytes32[],uint256) preconf, address challenger)
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerFilterer) WatchChallenged(opts *bind.WatchOpts, sink chan<- *ContractPreconfChallengeManagerChallenged) (event.Subscription, error) {

	logs, sub, err := _ContractPreconfChallengeManager.contract.WatchLogs(opts, "Challenged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPreconfChallengeManagerChallenged)
				if err := _ContractPreconfChallengeManager.contract.UnpackLog(event, "Challenged", log); err != nil {
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

// ParseChallenged is a log parse operation binding the contract event 0x9955dcf204285cdb270c001137aafffc8da100b0f75a21ace093b1bb92b49dc0.
//
// Solidity: event Challenged((bytes32[],uint256) preconf, address challenger)
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerFilterer) ParseChallenged(log types.Log) (*ContractPreconfChallengeManagerChallenged, error) {
	event := new(ContractPreconfChallengeManagerChallenged)
	if err := _ContractPreconfChallengeManager.contract.UnpackLog(event, "Challenged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractPreconfChallengeManagerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ContractPreconfChallengeManager contract.
type ContractPreconfChallengeManagerInitializedIterator struct {
	Event *ContractPreconfChallengeManagerInitialized // Event containing the contract specifics and raw log

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
func (it *ContractPreconfChallengeManagerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPreconfChallengeManagerInitialized)
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
		it.Event = new(ContractPreconfChallengeManagerInitialized)
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
func (it *ContractPreconfChallengeManagerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPreconfChallengeManagerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPreconfChallengeManagerInitialized represents a Initialized event raised by the ContractPreconfChallengeManager contract.
type ContractPreconfChallengeManagerInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerFilterer) FilterInitialized(opts *bind.FilterOpts) (*ContractPreconfChallengeManagerInitializedIterator, error) {

	logs, sub, err := _ContractPreconfChallengeManager.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ContractPreconfChallengeManagerInitializedIterator{contract: _ContractPreconfChallengeManager.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ContractPreconfChallengeManagerInitialized) (event.Subscription, error) {

	logs, sub, err := _ContractPreconfChallengeManager.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPreconfChallengeManagerInitialized)
				if err := _ContractPreconfChallengeManager.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerFilterer) ParseInitialized(log types.Log) (*ContractPreconfChallengeManagerInitialized, error) {
	event := new(ContractPreconfChallengeManagerInitialized)
	if err := _ContractPreconfChallengeManager.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractPreconfChallengeManagerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ContractPreconfChallengeManager contract.
type ContractPreconfChallengeManagerOwnershipTransferredIterator struct {
	Event *ContractPreconfChallengeManagerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ContractPreconfChallengeManagerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPreconfChallengeManagerOwnershipTransferred)
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
		it.Event = new(ContractPreconfChallengeManagerOwnershipTransferred)
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
func (it *ContractPreconfChallengeManagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPreconfChallengeManagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPreconfChallengeManagerOwnershipTransferred represents a OwnershipTransferred event raised by the ContractPreconfChallengeManager contract.
type ContractPreconfChallengeManagerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ContractPreconfChallengeManagerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractPreconfChallengeManager.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ContractPreconfChallengeManagerOwnershipTransferredIterator{contract: _ContractPreconfChallengeManager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ContractPreconfChallengeManagerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractPreconfChallengeManager.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPreconfChallengeManagerOwnershipTransferred)
				if err := _ContractPreconfChallengeManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ContractPreconfChallengeManager *ContractPreconfChallengeManagerFilterer) ParseOwnershipTransferred(log types.Log) (*ContractPreconfChallengeManagerOwnershipTransferred, error) {
	event := new(ContractPreconfChallengeManagerOwnershipTransferred)
	if err := _ContractPreconfChallengeManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
