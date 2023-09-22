// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package paymaster

import (
	"errors"
	"math/big"
	"strings"

	"github.com/tomochain/tomochain"
	"github.com/tomochain/tomochain/accounts/abi/bind"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/core/types"
	"github.com/tomochain/tomochain/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = tomochain.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// IPaymasterExecutionResult is an auto generated low-level Go binding around an user-defined struct.
type IPaymasterExecutionResult struct {
	Success bool
}

// IPaymasterTransaction is an auto generated low-level Go binding around an user-defined struct.
type IPaymasterTransaction struct {
	From common.Address
}

// PaymasterMetaData contains all meta data concerning the Paymaster contract.
var PaymasterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_context\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"internalType\":\"structIPaymaster.Transaction\",\"name\":\"_transaction\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"_txHash\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"internalType\":\"structIPaymaster.ExecutionResult\",\"name\":\"_txResult\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"_maxRefundedGas\",\"type\":\"uint256\"}],\"name\":\"postTransaction\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_txHash\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"internalType\":\"structIPaymaster.Transaction\",\"name\":\"_transaction\",\"type\":\"tuple\"}],\"name\":\"validateAndPayForPaymasterTransaction\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"magic\",\"type\":\"bytes4\"},{\"internalType\":\"bytes\",\"name\":\"context\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
}

// PaymasterABI is the input ABI used to generate the binding from.
// Deprecated: Use PaymasterMetaData.ABI instead.
var PaymasterABI = PaymasterMetaData.ABI

// Paymaster is an auto generated Go binding around an Ethereum contract.
type Paymaster struct {
	PaymasterCaller     // Read-only binding to the contract
	PaymasterTransactor // Write-only binding to the contract
	PaymasterFilterer   // Log filterer for contract events
}

// PaymasterCaller is an auto generated read-only Go binding around an Ethereum contract.
type PaymasterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PaymasterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PaymasterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PaymasterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PaymasterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PaymasterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PaymasterSession struct {
	Contract     *Paymaster        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PaymasterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PaymasterCallerSession struct {
	Contract *PaymasterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// PaymasterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PaymasterTransactorSession struct {
	Contract     *PaymasterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// PaymasterRaw is an auto generated low-level Go binding around an Ethereum contract.
type PaymasterRaw struct {
	Contract *Paymaster // Generic contract binding to access the raw methods on
}

// PaymasterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PaymasterCallerRaw struct {
	Contract *PaymasterCaller // Generic read-only contract binding to access the raw methods on
}

// PaymasterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PaymasterTransactorRaw struct {
	Contract *PaymasterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPaymaster creates a new instance of Paymaster, bound to a specific deployed contract.
func NewPaymaster(address common.Address, backend bind.ContractBackend) (*Paymaster, error) {
	contract, err := bindPaymaster(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Paymaster{PaymasterCaller: PaymasterCaller{contract: contract}, PaymasterTransactor: PaymasterTransactor{contract: contract}, PaymasterFilterer: PaymasterFilterer{contract: contract}}, nil
}

// NewPaymasterCaller creates a new read-only instance of Paymaster, bound to a specific deployed contract.
func NewPaymasterCaller(address common.Address, caller bind.ContractCaller) (*PaymasterCaller, error) {
	contract, err := bindPaymaster(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PaymasterCaller{contract: contract}, nil
}

// NewPaymasterTransactor creates a new write-only instance of Paymaster, bound to a specific deployed contract.
func NewPaymasterTransactor(address common.Address, transactor bind.ContractTransactor) (*PaymasterTransactor, error) {
	contract, err := bindPaymaster(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PaymasterTransactor{contract: contract}, nil
}

// NewPaymasterFilterer creates a new log filterer instance of Paymaster, bound to a specific deployed contract.
func NewPaymasterFilterer(address common.Address, filterer bind.ContractFilterer) (*PaymasterFilterer, error) {
	contract, err := bindPaymaster(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PaymasterFilterer{contract: contract}, nil
}

// bindPaymaster binds a generic wrapper to an already deployed contract.
func bindPaymaster(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PaymasterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Paymaster *PaymasterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Paymaster.Contract.PaymasterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Paymaster *PaymasterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Paymaster.Contract.PaymasterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Paymaster *PaymasterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Paymaster.Contract.PaymasterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Paymaster *PaymasterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Paymaster.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Paymaster *PaymasterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Paymaster.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Paymaster *PaymasterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Paymaster.Contract.contract.Transact(opts, method, params...)
}

// PostTransaction is a paid mutator transaction binding the contract method 0x62f29386.
//
// Solidity: function postTransaction(bytes _context, (address) _transaction, bytes32 _txHash, (bool) _txResult, uint256 _maxRefundedGas) payable returns()
func (_Paymaster *PaymasterTransactor) PostTransaction(opts *bind.TransactOpts, _context []byte, _transaction IPaymasterTransaction, _txHash [32]byte, _txResult IPaymasterExecutionResult, _maxRefundedGas *big.Int) (*types.Transaction, error) {
	return _Paymaster.contract.Transact(opts, "postTransaction", _context, _transaction, _txHash, _txResult, _maxRefundedGas)
}

// PostTransaction is a paid mutator transaction binding the contract method 0x62f29386.
//
// Solidity: function postTransaction(bytes _context, (address) _transaction, bytes32 _txHash, (bool) _txResult, uint256 _maxRefundedGas) payable returns()
func (_Paymaster *PaymasterSession) PostTransaction(_context []byte, _transaction IPaymasterTransaction, _txHash [32]byte, _txResult IPaymasterExecutionResult, _maxRefundedGas *big.Int) (*types.Transaction, error) {
	return _Paymaster.Contract.PostTransaction(&_Paymaster.TransactOpts, _context, _transaction, _txHash, _txResult, _maxRefundedGas)
}

// PostTransaction is a paid mutator transaction binding the contract method 0x62f29386.
//
// Solidity: function postTransaction(bytes _context, (address) _transaction, bytes32 _txHash, (bool) _txResult, uint256 _maxRefundedGas) payable returns()
func (_Paymaster *PaymasterTransactorSession) PostTransaction(_context []byte, _transaction IPaymasterTransaction, _txHash [32]byte, _txResult IPaymasterExecutionResult, _maxRefundedGas *big.Int) (*types.Transaction, error) {
	return _Paymaster.Contract.PostTransaction(&_Paymaster.TransactOpts, _context, _transaction, _txHash, _txResult, _maxRefundedGas)
}

// ValidateAndPayForPaymasterTransaction is a paid mutator transaction binding the contract method 0xd16a2f95.
//
// Solidity: function validateAndPayForPaymasterTransaction(bytes32 _txHash, (address) _transaction) payable returns(bytes4 magic, bytes context)
func (_Paymaster *PaymasterTransactor) ValidateAndPayForPaymasterTransaction(opts *bind.TransactOpts, _txHash [32]byte, _transaction IPaymasterTransaction) (*types.Transaction, error) {
	return _Paymaster.contract.Transact(opts, "validateAndPayForPaymasterTransaction", _txHash, _transaction)
}

// ValidateAndPayForPaymasterTransaction is a paid mutator transaction binding the contract method 0xd16a2f95.
//
// Solidity: function validateAndPayForPaymasterTransaction(bytes32 _txHash, (address) _transaction) payable returns(bytes4 magic, bytes context)
func (_Paymaster *PaymasterSession) ValidateAndPayForPaymasterTransaction(_txHash [32]byte, _transaction IPaymasterTransaction) (*types.Transaction, error) {
	return _Paymaster.Contract.ValidateAndPayForPaymasterTransaction(&_Paymaster.TransactOpts, _txHash, _transaction)
}

// ValidateAndPayForPaymasterTransaction is a paid mutator transaction binding the contract method 0xd16a2f95.
//
// Solidity: function validateAndPayForPaymasterTransaction(bytes32 _txHash, (address) _transaction) payable returns(bytes4 magic, bytes context)
func (_Paymaster *PaymasterTransactorSession) ValidateAndPayForPaymasterTransaction(_txHash [32]byte, _transaction IPaymasterTransaction) (*types.Transaction, error) {
	return _Paymaster.Contract.ValidateAndPayForPaymasterTransaction(&_Paymaster.TransactOpts, _txHash, _transaction)
}
