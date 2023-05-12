// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package response

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

// ContractMetaData contains all meta data concerning the Contract contract.
var ContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"taskId\",\"type\":\"string\"}],\"name\":\"read\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"taskResultMap\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"taskId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"taskResult\",\"type\":\"string\"}],\"name\":\"write\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610885806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80631e7c019714610046578063616ffe8314610062578063f748780c14610092575b600080fd5b610060600480360381019061005b91906103b2565b6100c2565b005b61007c6004803603810190610077919061042a565b6100f2565b60405161008991906104f2565b60405180910390f35b6100ac60048036038101906100a7919061042a565b6101a2565b6040516100b991906104f2565b60405180910390f35b806000836040516100d39190610550565b908152602001604051809103902090816100ed919061077d565b505050565b60606000826040516101049190610550565b9081526020016040518091039020805461011d90610596565b80601f016020809104026020016040519081016040528092919081815260200182805461014990610596565b80156101965780601f1061016b57610100808354040283529160200191610196565b820191906000526020600020905b81548152906001019060200180831161017957829003601f168201915b50505050509050919050565b60008180516020810182018051848252602083016020850120818352809550505050505060009150905080546101d790610596565b80601f016020809104026020016040519081016040528092919081815260200182805461020390610596565b80156102505780601f1061022557610100808354040283529160200191610250565b820191906000526020600020905b81548152906001019060200180831161023357829003601f168201915b505050505081565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6102bf82610276565b810181811067ffffffffffffffff821117156102de576102dd610287565b5b80604052505050565b60006102f1610258565b90506102fd82826102b6565b919050565b600067ffffffffffffffff82111561031d5761031c610287565b5b61032682610276565b9050602081019050919050565b82818337600083830152505050565b600061035561035084610302565b6102e7565b90508281526020810184848401111561037157610370610271565b5b61037c848285610333565b509392505050565b600082601f8301126103995761039861026c565b5b81356103a9848260208601610342565b91505092915050565b600080604083850312156103c9576103c8610262565b5b600083013567ffffffffffffffff8111156103e7576103e6610267565b5b6103f385828601610384565b925050602083013567ffffffffffffffff81111561041457610413610267565b5b61042085828601610384565b9150509250929050565b6000602082840312156104405761043f610262565b5b600082013567ffffffffffffffff81111561045e5761045d610267565b5b61046a84828501610384565b91505092915050565b600081519050919050565b600082825260208201905092915050565b60005b838110156104ad578082015181840152602081019050610492565b60008484015250505050565b60006104c482610473565b6104ce818561047e565b93506104de81856020860161048f565b6104e781610276565b840191505092915050565b6000602082019050818103600083015261050c81846104b9565b905092915050565b600081905092915050565b600061052a82610473565b6105348185610514565b935061054481856020860161048f565b80840191505092915050565b600061055c828461051f565b915081905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b600060028204905060018216806105ae57607f821691505b6020821081036105c1576105c0610567565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b6000600883026106297fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff826105ec565b61063386836105ec565b95508019841693508086168417925050509392505050565b6000819050919050565b6000819050919050565b600061067a6106756106708461064b565b610655565b61064b565b9050919050565b6000819050919050565b6106948361065f565b6106a86106a082610681565b8484546105f9565b825550505050565b600090565b6106bd6106b0565b6106c881848461068b565b505050565b5b818110156106ec576106e16000826106b5565b6001810190506106ce565b5050565b601f82111561073157610702816105c7565b61070b846105dc565b8101602085101561071a578190505b61072e610726856105dc565b8301826106cd565b50505b505050565b600082821c905092915050565b600061075460001984600802610736565b1980831691505092915050565b600061076d8383610743565b9150826002028217905092915050565b61078682610473565b67ffffffffffffffff81111561079f5761079e610287565b5b6107a98254610596565b6107b48282856106f0565b600060209050601f8311600181146107e757600084156107d5578287015190505b6107df8582610761565b865550610847565b601f1984166107f5866105c7565b60005b8281101561081d578489015182556001820191506020850194506020810190506107f8565b8683101561083a5784890151610836601f891682610743565b8355505b6001600288020188555050505b50505050505056fea2646970667358221220ea5722fd20e0bedaee85ed8a66d624068a7ea611dd6628cbc861f72d580ed10564736f6c63430008120033",
}

// ContractABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractMetaData.ABI instead.
var ContractABI = ContractMetaData.ABI

// ContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractMetaData.Bin instead.
var ContractBin = ContractMetaData.Bin

// DeployContract deploys a new Ethereum contract, binding an instance of Contract to it.
func DeployContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Contract, error) {
	parsed, err := ContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// Contract is an auto generated Go binding around an Ethereum contract.
type Contract struct {
	ContractCaller     // Read-only binding to the contract
	ContractTransactor // Write-only binding to the contract
	ContractFilterer   // Log filterer for contract events
}

// ContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractSession struct {
	Contract     *Contract         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractCallerSession struct {
	Contract *ContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractTransactorSession struct {
	Contract     *ContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractRaw struct {
	Contract *Contract // Generic contract binding to access the raw methods on
}

// ContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractCallerRaw struct {
	Contract *ContractCaller // Generic read-only contract binding to access the raw methods on
}

// ContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractTransactorRaw struct {
	Contract *ContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContract creates a new instance of Contract, bound to a specific deployed contract.
func NewContract(address common.Address, backend bind.ContractBackend) (*Contract, error) {
	contract, err := bindContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// NewContractCaller creates a new read-only instance of Contract, bound to a specific deployed contract.
func NewContractCaller(address common.Address, caller bind.ContractCaller) (*ContractCaller, error) {
	contract, err := bindContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractCaller{contract: contract}, nil
}

// NewContractTransactor creates a new write-only instance of Contract, bound to a specific deployed contract.
func NewContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractTransactor, error) {
	contract, err := bindContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractTransactor{contract: contract}, nil
}

// NewContractFilterer creates a new log filterer instance of Contract, bound to a specific deployed contract.
func NewContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractFilterer, error) {
	contract, err := bindContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractFilterer{contract: contract}, nil
}

// bindContract binds a generic wrapper to an already deployed contract.
func bindContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.ContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transact(opts, method, params...)
}

// Read is a free data retrieval call binding the contract method 0x616ffe83.
//
// Solidity: function read(string taskId) view returns(string)
func (_Contract *ContractCaller) Read(opts *bind.CallOpts, taskId string) (string, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "read", taskId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Read is a free data retrieval call binding the contract method 0x616ffe83.
//
// Solidity: function read(string taskId) view returns(string)
func (_Contract *ContractSession) Read(taskId string) (string, error) {
	return _Contract.Contract.Read(&_Contract.CallOpts, taskId)
}

// Read is a free data retrieval call binding the contract method 0x616ffe83.
//
// Solidity: function read(string taskId) view returns(string)
func (_Contract *ContractCallerSession) Read(taskId string) (string, error) {
	return _Contract.Contract.Read(&_Contract.CallOpts, taskId)
}

// TaskResultMap is a free data retrieval call binding the contract method 0xf748780c.
//
// Solidity: function taskResultMap(string ) view returns(string)
func (_Contract *ContractCaller) TaskResultMap(opts *bind.CallOpts, arg0 string) (string, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "taskResultMap", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TaskResultMap is a free data retrieval call binding the contract method 0xf748780c.
//
// Solidity: function taskResultMap(string ) view returns(string)
func (_Contract *ContractSession) TaskResultMap(arg0 string) (string, error) {
	return _Contract.Contract.TaskResultMap(&_Contract.CallOpts, arg0)
}

// TaskResultMap is a free data retrieval call binding the contract method 0xf748780c.
//
// Solidity: function taskResultMap(string ) view returns(string)
func (_Contract *ContractCallerSession) TaskResultMap(arg0 string) (string, error) {
	return _Contract.Contract.TaskResultMap(&_Contract.CallOpts, arg0)
}

// Write is a paid mutator transaction binding the contract method 0x1e7c0197.
//
// Solidity: function write(string taskId, string taskResult) returns()
func (_Contract *ContractTransactor) Write(opts *bind.TransactOpts, taskId string, taskResult string) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "write", taskId, taskResult)
}

// Write is a paid mutator transaction binding the contract method 0x1e7c0197.
//
// Solidity: function write(string taskId, string taskResult) returns()
func (_Contract *ContractSession) Write(taskId string, taskResult string) (*types.Transaction, error) {
	return _Contract.Contract.Write(&_Contract.TransactOpts, taskId, taskResult)
}

// Write is a paid mutator transaction binding the contract method 0x1e7c0197.
//
// Solidity: function write(string taskId, string taskResult) returns()
func (_Contract *ContractTransactorSession) Write(taskId string, taskResult string) (*types.Transaction, error) {
	return _Contract.Contract.Write(&_Contract.TransactOpts, taskId, taskResult)
}
