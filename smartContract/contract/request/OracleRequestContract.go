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

// ContractMetaData contains all meta data concerning the Contract contract.
var ContractMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"taskId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"value\",\"type\":\"string\"}],\"name\":\"SentEvent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"taskId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"value\",\"type\":\"string\"}],\"name\":\"request\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610398806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c806301d511f114610030575b600080fd5b61004a600480360381019061004591906101e5565b61004c565b005b7fe2a195737141626111612ba29c267b27f85c0b0b308764e9a6385294e1ae112433838360405161007f9392919061031d565b60405180910390a15050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6100f2826100a9565b810181811067ffffffffffffffff82111715610111576101106100ba565b5b80604052505050565b600061012461008b565b905061013082826100e9565b919050565b600067ffffffffffffffff8211156101505761014f6100ba565b5b610159826100a9565b9050602081019050919050565b82818337600083830152505050565b600061018861018384610135565b61011a565b9050828152602081018484840111156101a4576101a36100a4565b5b6101af848285610166565b509392505050565b600082601f8301126101cc576101cb61009f565b5b81356101dc848260208601610175565b91505092915050565b600080604083850312156101fc576101fb610095565b5b600083013567ffffffffffffffff81111561021a5761021961009a565b5b610226858286016101b7565b925050602083013567ffffffffffffffff8111156102475761024661009a565b5b610253858286016101b7565b9150509250929050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006102888261025d565b9050919050565b6102988161027d565b82525050565b600081519050919050565b600082825260208201905092915050565b60005b838110156102d85780820151818401526020810190506102bd565b60008484015250505050565b60006102ef8261029e565b6102f981856102a9565b93506103098185602086016102ba565b610312816100a9565b840191505092915050565b6000606082019050610332600083018661028f565b818103602083015261034481856102e4565b9050818103604083015261035881846102e4565b905094935050505056fea2646970667358221220aa8552c2d80923702fd65d62cd7fcb99cd9baf81df9dee12071807474a20722a64736f6c63430008120033",
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

// Request is a paid mutator transaction binding the contract method 0x01d511f1.
//
// Solidity: function request(string taskId, string value) returns()
func (_Contract *ContractTransactor) Request(opts *bind.TransactOpts, taskId string, value string) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "request", taskId, value)
}

// Request is a paid mutator transaction binding the contract method 0x01d511f1.
//
// Solidity: function request(string taskId, string value) returns()
func (_Contract *ContractSession) Request(taskId string, value string) (*types.Transaction, error) {
	return _Contract.Contract.Request(&_Contract.TransactOpts, taskId, value)
}

// Request is a paid mutator transaction binding the contract method 0x01d511f1.
//
// Solidity: function request(string taskId, string value) returns()
func (_Contract *ContractTransactorSession) Request(taskId string, value string) (*types.Transaction, error) {
	return _Contract.Contract.Request(&_Contract.TransactOpts, taskId, value)
}

// ContractSentEventIterator is returned from FilterSentEvent and is used to iterate over the raw logs and unpacked data for SentEvent events raised by the Contract contract.
type ContractSentEventIterator struct {
	Event *ContractSentEvent // Event containing the contract specifics and raw log

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
func (it *ContractSentEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractSentEvent)
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
		it.Event = new(ContractSentEvent)
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
func (it *ContractSentEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractSentEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractSentEvent represents a SentEvent event raised by the Contract contract.
type ContractSentEvent struct {
	Sender common.Address
	TaskId string
	Value  string
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterSentEvent is a free log retrieval operation binding the contract event 0xe2a195737141626111612ba29c267b27f85c0b0b308764e9a6385294e1ae1124.
//
// Solidity: event SentEvent(address sender, string taskId, string value)
func (_Contract *ContractFilterer) FilterSentEvent(opts *bind.FilterOpts) (*ContractSentEventIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "SentEvent")
	if err != nil {
		return nil, err
	}
	return &ContractSentEventIterator{contract: _Contract.contract, event: "SentEvent", logs: logs, sub: sub}, nil
}

// WatchSentEvent is a free log subscription operation binding the contract event 0xe2a195737141626111612ba29c267b27f85c0b0b308764e9a6385294e1ae1124.
//
// Solidity: event SentEvent(address sender, string taskId, string value)
func (_Contract *ContractFilterer) WatchSentEvent(opts *bind.WatchOpts, sink chan<- *ContractSentEvent) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "SentEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractSentEvent)
				if err := _Contract.contract.UnpackLog(event, "SentEvent", log); err != nil {
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

// ParseSentEvent is a log parse operation binding the contract event 0xe2a195737141626111612ba29c267b27f85c0b0b308764e9a6385294e1ae1124.
//
// Solidity: event SentEvent(address sender, string taskId, string value)
func (_Contract *ContractFilterer) ParseSentEvent(log types.Log) (*ContractSentEvent, error) {
	event := new(ContractSentEvent)
	if err := _Contract.contract.UnpackLog(event, "SentEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
