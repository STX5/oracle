// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package request

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

// RequestMetaData contains all meta data concerning the Request contract.
var RequestMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"value\",\"type\":\"string\"}],\"name\":\"RequestEvent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"url\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"pattern\",\"type\":\"string\"}],\"name\":\"request\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506104e6806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c806301d511f114610030575b600080fd5b61004a600480360381019061004591906102c1565b61004c565b005b60006040518060400160405280600981526020017f277b2275726c223a220000000000000000000000000000000000000000000000815250905060006040518060400160405280600d81526020017f222c227061747465726e223a2200000000000000000000000000000000000000815250905060006040518060400160405280600381526020017f227d2700000000000000000000000000000000000000000000000000000000008152509050600083868487856040516020016101159594939291906103aa565b60405160208183030381529060405290507f6a000d5e93e3e328bdb9bc8b8bb220bd5f92d74c38c3d1154dfac49a8293fa663382604051610157929190610480565b60405180910390a1505050505050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6101ce82610185565b810181811067ffffffffffffffff821117156101ed576101ec610196565b5b80604052505050565b6000610200610167565b905061020c82826101c5565b919050565b600067ffffffffffffffff82111561022c5761022b610196565b5b61023582610185565b9050602081019050919050565b82818337600083830152505050565b600061026461025f84610211565b6101f6565b9050828152602081018484840111156102805761027f610180565b5b61028b848285610242565b509392505050565b600082601f8301126102a8576102a761017b565b5b81356102b8848260208601610251565b91505092915050565b600080604083850312156102d8576102d7610171565b5b600083013567ffffffffffffffff8111156102f6576102f5610176565b5b61030285828601610293565b925050602083013567ffffffffffffffff81111561032357610322610176565b5b61032f85828601610293565b9150509250929050565b600081519050919050565b600081905092915050565b60005b8381101561036d578082015181840152602081019050610352565b60008484015250505050565b600061038482610339565b61038e8185610344565b935061039e81856020860161034f565b80840191505092915050565b60006103b68288610379565b91506103c28287610379565b91506103ce8286610379565b91506103da8285610379565b91506103e68284610379565b91508190509695505050505050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610420826103f5565b9050919050565b61043081610415565b82525050565b600082825260208201905092915050565b600061045282610339565b61045c8185610436565b935061046c81856020860161034f565b61047581610185565b840191505092915050565b60006040820190506104956000830185610427565b81810360208301526104a78184610447565b9050939250505056fea2646970667358221220832148ace8dbed75c8095b1dadb726ed0fba018f9cead9415f1b7781c4dec86964736f6c63430008120033",
}

// RequestABI is the input ABI used to generate the binding from.
// Deprecated: Use RequestMetaData.ABI instead.
var RequestABI = RequestMetaData.ABI

// RequestBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use RequestMetaData.Bin instead.
var RequestBin = RequestMetaData.Bin

// DeployRequest deploys a new Ethereum contract, binding an instance of Request to it.
func DeployRequest(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Request, error) {
	parsed, err := RequestMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RequestBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Request{RequestCaller: RequestCaller{contract: contract}, RequestTransactor: RequestTransactor{contract: contract}, RequestFilterer: RequestFilterer{contract: contract}}, nil
}

// Request is an auto generated Go binding around an Ethereum contract.
type Request struct {
	RequestCaller     // Read-only binding to the contract
	RequestTransactor // Write-only binding to the contract
	RequestFilterer   // Log filterer for contract events
}

// RequestCaller is an auto generated read-only Go binding around an Ethereum contract.
type RequestCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RequestTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RequestTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RequestFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RequestFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RequestSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RequestSession struct {
	Contract     *Request          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RequestCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RequestCallerSession struct {
	Contract *RequestCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// RequestTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RequestTransactorSession struct {
	Contract     *RequestTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// RequestRaw is an auto generated low-level Go binding around an Ethereum contract.
type RequestRaw struct {
	Contract *Request // Generic contract binding to access the raw methods on
}

// RequestCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RequestCallerRaw struct {
	Contract *RequestCaller // Generic read-only contract binding to access the raw methods on
}

// RequestTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RequestTransactorRaw struct {
	Contract *RequestTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRequest creates a new instance of Request, bound to a specific deployed contract.
func NewRequest(address common.Address, backend bind.ContractBackend) (*Request, error) {
	contract, err := bindRequest(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Request{RequestCaller: RequestCaller{contract: contract}, RequestTransactor: RequestTransactor{contract: contract}, RequestFilterer: RequestFilterer{contract: contract}}, nil
}

// NewRequestCaller creates a new read-only instance of Request, bound to a specific deployed contract.
func NewRequestCaller(address common.Address, caller bind.ContractCaller) (*RequestCaller, error) {
	contract, err := bindRequest(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RequestCaller{contract: contract}, nil
}

// NewRequestTransactor creates a new write-only instance of Request, bound to a specific deployed contract.
func NewRequestTransactor(address common.Address, transactor bind.ContractTransactor) (*RequestTransactor, error) {
	contract, err := bindRequest(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RequestTransactor{contract: contract}, nil
}

// NewRequestFilterer creates a new log filterer instance of Request, bound to a specific deployed contract.
func NewRequestFilterer(address common.Address, filterer bind.ContractFilterer) (*RequestFilterer, error) {
	contract, err := bindRequest(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RequestFilterer{contract: contract}, nil
}

// bindRequest binds a generic wrapper to an already deployed contract.
func bindRequest(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RequestMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Request *RequestRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Request.Contract.RequestCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Request *RequestRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Request.Contract.RequestTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Request *RequestRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Request.Contract.RequestTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Request *RequestCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Request.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Request *RequestTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Request.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Request *RequestTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Request.Contract.contract.Transact(opts, method, params...)
}

// Request is a paid mutator transaction binding the contract method 0x01d511f1.
//
// Solidity: function request(string url, string pattern) returns()
func (_Request *RequestTransactor) Request(opts *bind.TransactOpts, url string, pattern string) (*types.Transaction, error) {
	return _Request.contract.Transact(opts, "request", url, pattern)
}

// Request is a paid mutator transaction binding the contract method 0x01d511f1.
//
// Solidity: function request(string url, string pattern) returns()
func (_Request *RequestSession) Request(url string, pattern string) (*types.Transaction, error) {
	return _Request.Contract.Request(&_Request.TransactOpts, url, pattern)
}

// Request is a paid mutator transaction binding the contract method 0x01d511f1.
//
// Solidity: function request(string url, string pattern) returns()
func (_Request *RequestTransactorSession) Request(url string, pattern string) (*types.Transaction, error) {
	return _Request.Contract.Request(&_Request.TransactOpts, url, pattern)
}

// RequestRequestEventIterator is returned from FilterRequestEvent and is used to iterate over the raw logs and unpacked data for RequestEvent events raised by the Request contract.
type RequestRequestEventIterator struct {
	Event *RequestRequestEvent // Event containing the contract specifics and raw log

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
func (it *RequestRequestEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RequestRequestEvent)
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
		it.Event = new(RequestRequestEvent)
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
func (it *RequestRequestEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RequestRequestEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RequestRequestEvent represents a RequestEvent event raised by the Request contract.
type RequestRequestEvent struct {
	From  common.Address
	Value string
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterRequestEvent is a free log retrieval operation binding the contract event 0x6a000d5e93e3e328bdb9bc8b8bb220bd5f92d74c38c3d1154dfac49a8293fa66.
//
// Solidity: event RequestEvent(address from, string value)
func (_Request *RequestFilterer) FilterRequestEvent(opts *bind.FilterOpts) (*RequestRequestEventIterator, error) {

	logs, sub, err := _Request.contract.FilterLogs(opts, "RequestEvent")
	if err != nil {
		return nil, err
	}
	return &RequestRequestEventIterator{contract: _Request.contract, event: "RequestEvent", logs: logs, sub: sub}, nil
}

// WatchRequestEvent is a free log subscription operation binding the contract event 0x6a000d5e93e3e328bdb9bc8b8bb220bd5f92d74c38c3d1154dfac49a8293fa66.
//
// Solidity: event RequestEvent(address from, string value)
func (_Request *RequestFilterer) WatchRequestEvent(opts *bind.WatchOpts, sink chan<- *RequestRequestEvent) (event.Subscription, error) {

	logs, sub, err := _Request.contract.WatchLogs(opts, "RequestEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RequestRequestEvent)
				if err := _Request.contract.UnpackLog(event, "RequestEvent", log); err != nil {
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

// ParseRequestEvent is a log parse operation binding the contract event 0x6a000d5e93e3e328bdb9bc8b8bb220bd5f92d74c38c3d1154dfac49a8293fa66.
//
// Solidity: event RequestEvent(address from, string value)
func (_Request *RequestFilterer) ParseRequestEvent(log types.Log) (*RequestRequestEvent, error) {
	event := new(RequestRequestEvent)
	if err := _Request.contract.UnpackLog(event, "RequestEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
