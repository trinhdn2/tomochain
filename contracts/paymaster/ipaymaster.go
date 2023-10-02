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
	Bin: "0x608060405234801561001057600080fd5b506103ed806100206000396000f3fe6080604052600436106100345760003560e01c806362f2938614610039578063d16a2f951461004e578063f0fdf83414610078575b600080fd5b61004c610047366004610215565b6100a6565b005b61006161005c3660046102bb565b61012f565b60405161006f9291906102e8565b60405180910390f35b34801561008457600080fd5b50610098610093366004610347565b6101dc565b60405190815260200161006f565b600080546001810182558180526000805160206103988339815191528101558390036101085760405162461bcd60e51b815260206004820152600c60248201526b0cadae0e8f240e8f090c2e6d60a31b60448201526064015b60405180910390fd5b50506000805460018101825590805260008051602061039883398151915281015550505050565b60008054600181018255818052600080516020610398833981519152810155606083820361018e5760405162461bcd60e51b815260206004820152600c60248201526b0cadae0e8f240e8f090c2e6d60a31b60448201526064016100ff565b6000805460018101825581805260008051602061039883398151915281015554604080516020810192909252016040516020818303038152906040526101d390610360565b91509250929050565b600081815481106101ec57600080fd5b600091825260209091200154905081565b60006020828403121561020f57600080fd5b50919050565b60008060008060008060a0878903121561022e57600080fd5b863567ffffffffffffffff8082111561024657600080fd5b818901915089601f83011261025a57600080fd5b81358181111561026957600080fd5b8a602082850101111561027b57600080fd5b602092830198509650610292918a915089016101fd565b9350604087013592506102a888606089016101fd565b9150608087013590509295509295509295565b600080604083850312156102ce57600080fd5b823591506102df84602085016101fd565b90509250929050565b63ffffffff60e01b8316815260006020604081840152835180604085015260005b8181101561032557858101830151858201606001528201610309565b506000606082860101526060601f19601f830116850101925050509392505050565b60006020828403121561035957600080fd5b5035919050565b805160208201516001600160e01b0319808216929190600483101561038f5780818460040360031b1b83161693505b50505091905056fe290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563a2646970667358221220165af7dae870bb8785dfd78360b524bd9fdf0ce426b2d031d47a0099783ef58e64736f6c63430008150033",
}

// PaymasterABI is the input ABI used to generate the binding from.
// Deprecated: Use PaymasterMetaData.ABI instead.
var PaymasterABI = PaymasterMetaData.ABI

// PaymasterBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use PaymasterMetaData.Bin instead.
var PaymasterBin = PaymasterMetaData.Bin

// DeployPaymaster deploys a new Ethereum contract, binding an instance of Paymaster to it.
func DeployPaymaster(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Paymaster, error) {
	parsed, err := PaymasterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PaymasterBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Paymaster{PaymasterCaller: PaymasterCaller{contract: contract}, PaymasterTransactor: PaymasterTransactor{contract: contract}, PaymasterFilterer: PaymasterFilterer{contract: contract}}, nil
}

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
