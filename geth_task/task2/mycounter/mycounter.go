// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mycounter

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

// MycounterMetaData contains all meta data concerning the Mycounter contract.
var MycounterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_num\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"number\",\"type\":\"uint256\"}],\"name\":\"IncreNumber\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"increment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"number\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newNumber\",\"type\":\"uint256\"}],\"name\":\"setNumber\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b506040516102b73803806102b78339818101604052810190602e9190606b565b805f81905550506091565b5f5ffd5b5f819050919050565b604d81603d565b81146056575f5ffd5b50565b5f815190506065816046565b92915050565b5f60208284031215607d57607c6039565b5b5f6088848285016059565b91505092915050565b6102198061009e5f395ff3fe608060405234801561000f575f5ffd5b506004361061003f575f3560e01c80633fb5c1cb146100435780638381f58a1461005f578063d09de08a1461007d575b5f5ffd5b61005d6004803603810190610058919061011c565b610087565b005b610067610090565b6040516100749190610156565b60405180910390f35b610085610095565b005b805f8190555050565b5f5481565b5f5f8154809291906100a69061019c565b91905055507f998d26fb0e10fa8363195144582a6825ce8755d3b6c6b1d9df0ac74c937d44945f546040516100db9190610156565b60405180910390a1565b5f5ffd5b5f819050919050565b6100fb816100e9565b8114610105575f5ffd5b50565b5f81359050610116816100f2565b92915050565b5f60208284031215610131576101306100e5565b5b5f61013e84828501610108565b91505092915050565b610150816100e9565b82525050565b5f6020820190506101695f830184610147565b92915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f6101a6826100e9565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036101d8576101d761016f565b5b60018201905091905056fea264697066735822122074c2dd9d9c89b02f8a8e883937def0f6bfcb7342e8432d702c69fdcbe9126dc864736f6c634300081d0033",
}

// MycounterABI is the input ABI used to generate the binding from.
// Deprecated: Use MycounterMetaData.ABI instead.
var MycounterABI = MycounterMetaData.ABI

// MycounterBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MycounterMetaData.Bin instead.
var MycounterBin = MycounterMetaData.Bin

// DeployMycounter deploys a new Ethereum contract, binding an instance of Mycounter to it.
func DeployMycounter(auth *bind.TransactOpts, backend bind.ContractBackend, _num *big.Int) (common.Address, *types.Transaction, *Mycounter, error) {
	parsed, err := MycounterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MycounterBin), backend, _num)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Mycounter{MycounterCaller: MycounterCaller{contract: contract}, MycounterTransactor: MycounterTransactor{contract: contract}, MycounterFilterer: MycounterFilterer{contract: contract}}, nil
}

// Mycounter is an auto generated Go binding around an Ethereum contract.
type Mycounter struct {
	MycounterCaller     // Read-only binding to the contract
	MycounterTransactor // Write-only binding to the contract
	MycounterFilterer   // Log filterer for contract events
}

// MycounterCaller is an auto generated read-only Go binding around an Ethereum contract.
type MycounterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MycounterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MycounterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MycounterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MycounterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MycounterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MycounterSession struct {
	Contract     *Mycounter        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MycounterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MycounterCallerSession struct {
	Contract *MycounterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// MycounterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MycounterTransactorSession struct {
	Contract     *MycounterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// MycounterRaw is an auto generated low-level Go binding around an Ethereum contract.
type MycounterRaw struct {
	Contract *Mycounter // Generic contract binding to access the raw methods on
}

// MycounterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MycounterCallerRaw struct {
	Contract *MycounterCaller // Generic read-only contract binding to access the raw methods on
}

// MycounterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MycounterTransactorRaw struct {
	Contract *MycounterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMycounter creates a new instance of Mycounter, bound to a specific deployed contract.
func NewMycounter(address common.Address, backend bind.ContractBackend) (*Mycounter, error) {
	contract, err := bindMycounter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Mycounter{MycounterCaller: MycounterCaller{contract: contract}, MycounterTransactor: MycounterTransactor{contract: contract}, MycounterFilterer: MycounterFilterer{contract: contract}}, nil
}

// NewMycounterCaller creates a new read-only instance of Mycounter, bound to a specific deployed contract.
func NewMycounterCaller(address common.Address, caller bind.ContractCaller) (*MycounterCaller, error) {
	contract, err := bindMycounter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MycounterCaller{contract: contract}, nil
}

// NewMycounterTransactor creates a new write-only instance of Mycounter, bound to a specific deployed contract.
func NewMycounterTransactor(address common.Address, transactor bind.ContractTransactor) (*MycounterTransactor, error) {
	contract, err := bindMycounter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MycounterTransactor{contract: contract}, nil
}

// NewMycounterFilterer creates a new log filterer instance of Mycounter, bound to a specific deployed contract.
func NewMycounterFilterer(address common.Address, filterer bind.ContractFilterer) (*MycounterFilterer, error) {
	contract, err := bindMycounter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MycounterFilterer{contract: contract}, nil
}

// bindMycounter binds a generic wrapper to an already deployed contract.
func bindMycounter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MycounterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Mycounter *MycounterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Mycounter.Contract.MycounterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Mycounter *MycounterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mycounter.Contract.MycounterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Mycounter *MycounterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Mycounter.Contract.MycounterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Mycounter *MycounterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Mycounter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Mycounter *MycounterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mycounter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Mycounter *MycounterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Mycounter.Contract.contract.Transact(opts, method, params...)
}

// Number is a free data retrieval call binding the contract method 0x8381f58a.
//
// Solidity: function number() view returns(uint256)
func (_Mycounter *MycounterCaller) Number(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Mycounter.contract.Call(opts, &out, "number")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Number is a free data retrieval call binding the contract method 0x8381f58a.
//
// Solidity: function number() view returns(uint256)
func (_Mycounter *MycounterSession) Number() (*big.Int, error) {
	return _Mycounter.Contract.Number(&_Mycounter.CallOpts)
}

// Number is a free data retrieval call binding the contract method 0x8381f58a.
//
// Solidity: function number() view returns(uint256)
func (_Mycounter *MycounterCallerSession) Number() (*big.Int, error) {
	return _Mycounter.Contract.Number(&_Mycounter.CallOpts)
}

// Increment is a paid mutator transaction binding the contract method 0xd09de08a.
//
// Solidity: function increment() returns()
func (_Mycounter *MycounterTransactor) Increment(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mycounter.contract.Transact(opts, "increment")
}

// Increment is a paid mutator transaction binding the contract method 0xd09de08a.
//
// Solidity: function increment() returns()
func (_Mycounter *MycounterSession) Increment() (*types.Transaction, error) {
	return _Mycounter.Contract.Increment(&_Mycounter.TransactOpts)
}

// Increment is a paid mutator transaction binding the contract method 0xd09de08a.
//
// Solidity: function increment() returns()
func (_Mycounter *MycounterTransactorSession) Increment() (*types.Transaction, error) {
	return _Mycounter.Contract.Increment(&_Mycounter.TransactOpts)
}

// SetNumber is a paid mutator transaction binding the contract method 0x3fb5c1cb.
//
// Solidity: function setNumber(uint256 newNumber) returns()
func (_Mycounter *MycounterTransactor) SetNumber(opts *bind.TransactOpts, newNumber *big.Int) (*types.Transaction, error) {
	return _Mycounter.contract.Transact(opts, "setNumber", newNumber)
}

// SetNumber is a paid mutator transaction binding the contract method 0x3fb5c1cb.
//
// Solidity: function setNumber(uint256 newNumber) returns()
func (_Mycounter *MycounterSession) SetNumber(newNumber *big.Int) (*types.Transaction, error) {
	return _Mycounter.Contract.SetNumber(&_Mycounter.TransactOpts, newNumber)
}

// SetNumber is a paid mutator transaction binding the contract method 0x3fb5c1cb.
//
// Solidity: function setNumber(uint256 newNumber) returns()
func (_Mycounter *MycounterTransactorSession) SetNumber(newNumber *big.Int) (*types.Transaction, error) {
	return _Mycounter.Contract.SetNumber(&_Mycounter.TransactOpts, newNumber)
}

// MycounterIncreNumberIterator is returned from FilterIncreNumber and is used to iterate over the raw logs and unpacked data for IncreNumber events raised by the Mycounter contract.
type MycounterIncreNumberIterator struct {
	Event *MycounterIncreNumber // Event containing the contract specifics and raw log

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
func (it *MycounterIncreNumberIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MycounterIncreNumber)
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
		it.Event = new(MycounterIncreNumber)
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
func (it *MycounterIncreNumberIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MycounterIncreNumberIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MycounterIncreNumber represents a IncreNumber event raised by the Mycounter contract.
type MycounterIncreNumber struct {
	Number *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterIncreNumber is a free log retrieval operation binding the contract event 0x998d26fb0e10fa8363195144582a6825ce8755d3b6c6b1d9df0ac74c937d4494.
//
// Solidity: event IncreNumber(uint256 number)
func (_Mycounter *MycounterFilterer) FilterIncreNumber(opts *bind.FilterOpts) (*MycounterIncreNumberIterator, error) {

	logs, sub, err := _Mycounter.contract.FilterLogs(opts, "IncreNumber")
	if err != nil {
		return nil, err
	}
	return &MycounterIncreNumberIterator{contract: _Mycounter.contract, event: "IncreNumber", logs: logs, sub: sub}, nil
}

// WatchIncreNumber is a free log subscription operation binding the contract event 0x998d26fb0e10fa8363195144582a6825ce8755d3b6c6b1d9df0ac74c937d4494.
//
// Solidity: event IncreNumber(uint256 number)
func (_Mycounter *MycounterFilterer) WatchIncreNumber(opts *bind.WatchOpts, sink chan<- *MycounterIncreNumber) (event.Subscription, error) {

	logs, sub, err := _Mycounter.contract.WatchLogs(opts, "IncreNumber")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MycounterIncreNumber)
				if err := _Mycounter.contract.UnpackLog(event, "IncreNumber", log); err != nil {
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

// ParseIncreNumber is a log parse operation binding the contract event 0x998d26fb0e10fa8363195144582a6825ce8755d3b6c6b1d9df0ac74c937d4494.
//
// Solidity: event IncreNumber(uint256 number)
func (_Mycounter *MycounterFilterer) ParseIncreNumber(log types.Log) (*MycounterIncreNumber, error) {
	event := new(MycounterIncreNumber)
	if err := _Mycounter.contract.UnpackLog(event, "IncreNumber", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
