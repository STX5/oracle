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

// ResponseMetaData contains all meta data concerning the Response contract.
var ResponseMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"getValue\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"value\",\"type\":\"string\"}],\"name\":\"setValue\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610766806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c8063209652551461003b578063510fbf8d14610059575b600080fd5b610043610075565b6040516100509190610223565b60405180910390f35b610073600480360381019061006e91906103ec565b610143565b005b60606000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002080546100c090610477565b80601f01602080910402602001604051908101604052809291908181526020018280546100ec90610477565b80156101395780601f1061010e57610100808354040283529160200191610139565b820191906000526020600020905b81548152906001019060200180831161011c57829003601f168201915b5050505050905090565b806000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020908161018e919061065e565b505050565b600081519050919050565b600082825260208201905092915050565b60005b838110156101cd5780820151818401526020810190506101b2565b60008484015250505050565b6000601f19601f8301169050919050565b60006101f582610193565b6101ff818561019e565b935061020f8185602086016101af565b610218816101d9565b840191505092915050565b6000602082019050818103600083015261023d81846101ea565b905092915050565b6000604051905090565b600080fd5b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061028482610259565b9050919050565b61029481610279565b811461029f57600080fd5b50565b6000813590506102b18161028b565b92915050565b600080fd5b600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6102f9826101d9565b810181811067ffffffffffffffff82111715610318576103176102c1565b5b80604052505050565b600061032b610245565b905061033782826102f0565b919050565b600067ffffffffffffffff821115610357576103566102c1565b5b610360826101d9565b9050602081019050919050565b82818337600083830152505050565b600061038f61038a8461033c565b610321565b9050828152602081018484840111156103ab576103aa6102bc565b5b6103b684828561036d565b509392505050565b600082601f8301126103d3576103d26102b7565b5b81356103e384826020860161037c565b91505092915050565b600080604083850312156104035761040261024f565b5b6000610411858286016102a2565b925050602083013567ffffffffffffffff81111561043257610431610254565b5b61043e858286016103be565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000600282049050600182168061048f57607f821691505b6020821081036104a2576104a1610448565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b60006008830261050a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff826104cd565b61051486836104cd565b95508019841693508086168417925050509392505050565b6000819050919050565b6000819050919050565b600061055b6105566105518461052c565b610536565b61052c565b9050919050565b6000819050919050565b61057583610540565b61058961058182610562565b8484546104da565b825550505050565b600090565b61059e610591565b6105a981848461056c565b505050565b5b818110156105cd576105c2600082610596565b6001810190506105af565b5050565b601f821115610612576105e3816104a8565b6105ec846104bd565b810160208510156105fb578190505b61060f610607856104bd565b8301826105ae565b50505b505050565b600082821c905092915050565b600061063560001984600802610617565b1980831691505092915050565b600061064e8383610624565b9150826002028217905092915050565b61066782610193565b67ffffffffffffffff8111156106805761067f6102c1565b5b61068a8254610477565b6106958282856105d1565b600060209050601f8311600181146106c857600084156106b6578287015190505b6106c08582610642565b865550610728565b601f1984166106d6866104a8565b60005b828110156106fe578489015182556001820191506020850194506020810190506106d9565b8683101561071b5784890151610717601f891682610624565b8355505b6001600288020188555050505b50505050505056fea2646970667358221220eef28ff65f908abbe38a5e0a6d49aafa1b172d2b352c93d70ca9cc51b4b6d1fb64736f6c63430008120033",
}

// ResponseABI is the input ABI used to generate the binding from.
// Deprecated: Use ResponseMetaData.ABI instead.
var ResponseABI = ResponseMetaData.ABI

// ResponseBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ResponseMetaData.Bin instead.
var ResponseBin = ResponseMetaData.Bin

// DeployResponse deploys a new Ethereum contract, binding an instance of Response to it.
func DeployResponse(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Response, error) {
	parsed, err := ResponseMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ResponseBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Response{ResponseCaller: ResponseCaller{contract: contract}, ResponseTransactor: ResponseTransactor{contract: contract}, ResponseFilterer: ResponseFilterer{contract: contract}}, nil
}

// Response is an auto generated Go binding around an Ethereum contract.
type Response struct {
	ResponseCaller     // Read-only binding to the contract
	ResponseTransactor // Write-only binding to the contract
	ResponseFilterer   // Log filterer for contract events
}

// ResponseCaller is an auto generated read-only Go binding around an Ethereum contract.
type ResponseCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ResponseTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ResponseTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ResponseFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ResponseFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ResponseSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ResponseSession struct {
	Contract     *Response         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ResponseCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ResponseCallerSession struct {
	Contract *ResponseCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ResponseTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ResponseTransactorSession struct {
	Contract     *ResponseTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ResponseRaw is an auto generated low-level Go binding around an Ethereum contract.
type ResponseRaw struct {
	Contract *Response // Generic contract binding to access the raw methods on
}

// ResponseCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ResponseCallerRaw struct {
	Contract *ResponseCaller // Generic read-only contract binding to access the raw methods on
}

// ResponseTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ResponseTransactorRaw struct {
	Contract *ResponseTransactor // Generic write-only contract binding to access the raw methods on
}

// NewResponse creates a new instance of Response, bound to a specific deployed contract.
func NewResponse(address common.Address, backend bind.ContractBackend) (*Response, error) {
	contract, err := bindResponse(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Response{ResponseCaller: ResponseCaller{contract: contract}, ResponseTransactor: ResponseTransactor{contract: contract}, ResponseFilterer: ResponseFilterer{contract: contract}}, nil
}

// NewResponseCaller creates a new read-only instance of Response, bound to a specific deployed contract.
func NewResponseCaller(address common.Address, caller bind.ContractCaller) (*ResponseCaller, error) {
	contract, err := bindResponse(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ResponseCaller{contract: contract}, nil
}

// NewResponseTransactor creates a new write-only instance of Response, bound to a specific deployed contract.
func NewResponseTransactor(address common.Address, transactor bind.ContractTransactor) (*ResponseTransactor, error) {
	contract, err := bindResponse(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ResponseTransactor{contract: contract}, nil
}

// NewResponseFilterer creates a new log filterer instance of Response, bound to a specific deployed contract.
func NewResponseFilterer(address common.Address, filterer bind.ContractFilterer) (*ResponseFilterer, error) {
	contract, err := bindResponse(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ResponseFilterer{contract: contract}, nil
}

// bindResponse binds a generic wrapper to an already deployed contract.
func bindResponse(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ResponseMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Response *ResponseRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Response.Contract.ResponseCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Response *ResponseRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Response.Contract.ResponseTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Response *ResponseRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Response.Contract.ResponseTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Response *ResponseCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Response.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Response *ResponseTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Response.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Response *ResponseTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Response.Contract.contract.Transact(opts, method, params...)
}

// GetValue is a free data retrieval call binding the contract method 0x20965255.
//
// Solidity: function getValue() view returns(string)
func (_Response *ResponseCaller) GetValue(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Response.contract.Call(opts, &out, "getValue")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetValue is a free data retrieval call binding the contract method 0x20965255.
//
// Solidity: function getValue() view returns(string)
func (_Response *ResponseSession) GetValue() (string, error) {
	return _Response.Contract.GetValue(&_Response.CallOpts)
}

// GetValue is a free data retrieval call binding the contract method 0x20965255.
//
// Solidity: function getValue() view returns(string)
func (_Response *ResponseCallerSession) GetValue() (string, error) {
	return _Response.Contract.GetValue(&_Response.CallOpts)
}

// SetValue is a paid mutator transaction binding the contract method 0x510fbf8d.
//
// Solidity: function setValue(address from, string value) returns()
func (_Response *ResponseTransactor) SetValue(opts *bind.TransactOpts, from common.Address, value string) (*types.Transaction, error) {
	return _Response.contract.Transact(opts, "setValue", from, value)
}

// SetValue is a paid mutator transaction binding the contract method 0x510fbf8d.
//
// Solidity: function setValue(address from, string value) returns()
func (_Response *ResponseSession) SetValue(from common.Address, value string) (*types.Transaction, error) {
	return _Response.Contract.SetValue(&_Response.TransactOpts, from, value)
}

// SetValue is a paid mutator transaction binding the contract method 0x510fbf8d.
//
// Solidity: function setValue(address from, string value) returns()
func (_Response *ResponseTransactorSession) SetValue(from common.Address, value string) (*types.Transaction, error) {
	return _Response.Contract.SetValue(&_Response.TransactOpts, from, value)
}
