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

// Transaction is an auto generated low-level Go binding around an user-defined struct.
type Transaction struct {
	From common.Address
}

// PaymasterMetaData contains all meta data concerning the Paymaster contract.
var PaymasterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"a\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_context\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"internalType\":\"structTransaction\",\"name\":\"_transaction\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"_txHash\",\"type\":\"bytes32\"},{\"internalType\":\"enumExecutionResult\",\"name\":\"_txResult\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"_maxRefundedGas\",\"type\":\"uint256\"}],\"name\":\"postTransaction\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_txHash\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"internalType\":\"structTransaction\",\"name\":\"_transaction\",\"type\":\"tuple\"}],\"name\":\"validateAndPayForPaymasterTransaction\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"magic\",\"type\":\"bytes4\"},{\"internalType\":\"bytes\",\"name\":\"context\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x608060405234801561000f575f80fd5b506103de8061001d5f395ff3fe608060405260043610610036575f3560e01c806392979f5a14610041578063d16a2f9514610056578063f0fdf83414610080575f80fd5b3661003d57005b5f80fd5b61005461004f36600461020f565b6100ad565b005b6100696100643660046102b3565b610132565b6040516100779291906102de565b60405180910390f35b34801561008b575f80fd5b5061009f61009a36600461033a565b6101db565b604051908152602001610077565b5f80546001810182558180525f8051602061038983398151915281015583900361010d5760405162461bcd60e51b815260206004820152600c60248201526b0cadae0e8f240e8f090c2e6d60a31b60448201526064015b60405180910390fd5b50505f80546001810182559080525f8051602061038983398151915281015550505050565b5f80546001810182558180525f80516020610389833981519152810155606083820361018f5760405162461bcd60e51b815260206004820152600c60248201526b0cadae0e8f240e8f090c2e6d60a31b6044820152606401610104565b5f80546001810182558180525f8051602061038983398151915281015554604080516020810192909252016040516020818303038152906040526101d290610351565b91509250929050565b5f81815481106101e9575f80fd5b5f91825260209091200154905081565b5f60208284031215610209575f80fd5b50919050565b5f805f805f8060a08789031215610224575f80fd5b863567ffffffffffffffff8082111561023b575f80fd5b818901915089601f83011261024e575f80fd5b81358181111561025c575f80fd5b8a602082850101111561026d575f80fd5b602092830198509650610284918a915089016101f9565b93506040870135925060608701356002811061029e575f80fd5b80925050608087013590509295509295509295565b5f80604083850312156102c4575f80fd5b823591506102d584602085016101f9565b90509250929050565b63ffffffff60e01b831681525f602060408184015283518060408501525f5b81811015610319578581018301518582016060015282016102fd565b505f606082860101526060601f19601f830116850101925050509392505050565b5f6020828403121561034a575f80fd5b5035919050565b805160208201516001600160e01b031980821692919060048310156103805780818460040360031b1b83161693505b50505091905056fe290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563a26469706673582212201c5aff1ac1a4b9cebce8a30143926ce5b238497f156a952bfaff7af647a5a86064736f6c63430008150033",
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

// PostTransaction is a paid mutator transaction binding the contract method 0x92979f5a.
//
// Solidity: function postTransaction(bytes _context, (address) _transaction, bytes32 _txHash, uint8 _txResult, uint256 _maxRefundedGas) payable returns()
func (_Paymaster *PaymasterTransactor) PostTransaction(opts *bind.TransactOpts, _context []byte, _transaction Transaction, _txHash [32]byte, _txResult uint8, _maxRefundedGas *big.Int) (*types.Transaction, error) {
	return _Paymaster.contract.Transact(opts, "postTransaction", _context, _transaction, _txHash, _txResult, _maxRefundedGas)
}

// PostTransaction is a paid mutator transaction binding the contract method 0x92979f5a.
//
// Solidity: function postTransaction(bytes _context, (address) _transaction, bytes32 _txHash, uint8 _txResult, uint256 _maxRefundedGas) payable returns()
func (_Paymaster *PaymasterSession) PostTransaction(_context []byte, _transaction Transaction, _txHash [32]byte, _txResult uint8, _maxRefundedGas *big.Int) (*types.Transaction, error) {
	return _Paymaster.Contract.PostTransaction(&_Paymaster.TransactOpts, _context, _transaction, _txHash, _txResult, _maxRefundedGas)
}

// PostTransaction is a paid mutator transaction binding the contract method 0x92979f5a.
//
// Solidity: function postTransaction(bytes _context, (address) _transaction, bytes32 _txHash, uint8 _txResult, uint256 _maxRefundedGas) payable returns()
func (_Paymaster *PaymasterTransactorSession) PostTransaction(_context []byte, _transaction Transaction, _txHash [32]byte, _txResult uint8, _maxRefundedGas *big.Int) (*types.Transaction, error) {
	return _Paymaster.Contract.PostTransaction(&_Paymaster.TransactOpts, _context, _transaction, _txHash, _txResult, _maxRefundedGas)
}

// ValidateAndPayForPaymasterTransaction is a paid mutator transaction binding the contract method 0xd16a2f95.
//
// Solidity: function validateAndPayForPaymasterTransaction(bytes32 _txHash, (address) _transaction) payable returns(bytes4 magic, bytes context)
func (_Paymaster *PaymasterTransactor) ValidateAndPayForPaymasterTransaction(opts *bind.TransactOpts, _txHash [32]byte, _transaction Transaction) (*types.Transaction, error) {
	return _Paymaster.contract.Transact(opts, "validateAndPayForPaymasterTransaction", _txHash, _transaction)
}

// ValidateAndPayForPaymasterTransaction is a paid mutator transaction binding the contract method 0xd16a2f95.
//
// Solidity: function validateAndPayForPaymasterTransaction(bytes32 _txHash, (address) _transaction) payable returns(bytes4 magic, bytes context)
func (_Paymaster *PaymasterSession) ValidateAndPayForPaymasterTransaction(_txHash [32]byte, _transaction Transaction) (*types.Transaction, error) {
	return _Paymaster.Contract.ValidateAndPayForPaymasterTransaction(&_Paymaster.TransactOpts, _txHash, _transaction)
}

// ValidateAndPayForPaymasterTransaction is a paid mutator transaction binding the contract method 0xd16a2f95.
//
// Solidity: function validateAndPayForPaymasterTransaction(bytes32 _txHash, (address) _transaction) payable returns(bytes4 magic, bytes context)
func (_Paymaster *PaymasterTransactorSession) ValidateAndPayForPaymasterTransaction(_txHash [32]byte, _transaction Transaction) (*types.Transaction, error) {
	return _Paymaster.Contract.ValidateAndPayForPaymasterTransaction(&_Paymaster.TransactOpts, _txHash, _transaction)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Paymaster *PaymasterTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Paymaster.contract.Transact(opts, "receive", nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Paymaster *PaymasterSession) Receive() (*types.Transaction, error) {
	return _Paymaster.Contract.Receive(&_Paymaster.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Paymaster *PaymasterTransactorSession) Receive() (*types.Transaction, error) {
	return _Paymaster.Contract.Receive(&_Paymaster.TransactOpts)
}
