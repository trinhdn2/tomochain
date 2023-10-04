// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"errors"
	"fmt"
	"math"
	"math/big"

	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/contracts/paymaster"
	"github.com/tomochain/tomochain/core/types"
	"github.com/tomochain/tomochain/core/vm"
	"github.com/tomochain/tomochain/params"
)

var (
	errInsufficientBalanceForGas = errors.New("insufficient balance to pay for gas")
	invalidMagic                 = "fffffffff"
)

/*
The State Transitioning Model

A state transition is a change made when a transaction is applied to the current world state
The state transitioning model does all all the necessary work to work out a valid new state root.

1) Nonce handling
2) Pre pay gas
3) Create a new state object if the recipient is \0*32
4) Value transfer
== If contract creation ==

	4a) Attempt to run transaction data
	4b) If valid, use result as code for the new state object

== end ==
5) Run Script section
6) Derive new state root
*/
type StateTransition struct {
	gp         *GasPool
	msg        *Message
	gas        uint64
	gasPrice   *big.Int
	initialGas uint64
	value      *big.Int
	data       []byte
	state      vm.StateDB
	evm        *vm.EVM

	// Paymaster fields
	magic   [4]byte
	context []byte
}

// A Message contains the data derived from a single transaction that is relevant to state
// processing.
type Message struct {
	To              *common.Address
	From            common.Address
	Nonce           uint64
	Value           *big.Int
	GasLimit        uint64
	GasPrice        *big.Int
	Data            []byte
	PmAddress       common.Address
	PmPayload       []byte
	BalanceTokenFee *big.Int

	// When SkipAccountChecks is true, the message nonce is not checked against the
	// account nonce in state. It also disables checking that the sender is an EOA.
	// This field will be set to true for operations like RPC eth_call.
	SkipAccountChecks bool
}

// message no matter the execution itself is successful or not.
type ExecutionResult struct {
	UsedGas    uint64 // Total used gas but include the refunded gas
	Err        error  // Any error encountered during the execution(listed in core/vm/errors.go)
	ReturnData []byte // Returned data from evm(function result or data supplied with revert opcode)
}

// Unwrap returns the internal evm error which allows us for further
// analysis outside.
func (result *ExecutionResult) Unwrap() error {
	return result.Err
}

// Failed returns the indicator whether the execution is successful or not
func (result *ExecutionResult) Failed() bool { return result.Err != nil }

// Return is a helper function to help caller distinguish between revert reason
// and function return. Return returns the data after execution if no error occurs.
func (result *ExecutionResult) Return() []byte {
	if result.Err != nil {
		return nil
	}
	return common.CopyBytes(result.ReturnData)
}

// Revert returns the concrete revert reason if the execution is aborted by `REVERT`
// opcode. Note the reason can be nil if no data supplied with revert opcode.
func (result *ExecutionResult) Revert() []byte {
	if result.Err != vm.ErrExecutionReverted {
		return nil
	}
	return common.CopyBytes(result.ReturnData)
}

// IntrinsicGas computes the 'intrinsic gas' for a message with the given data.
func IntrinsicGas(data []byte, contractCreation, homestead bool) (uint64, error) {
	// Set the starting gas for the raw transaction
	var gas uint64
	if contractCreation && homestead {
		gas = params.TxGasContractCreation
	} else {
		gas = params.TxGas
	}
	// Bump the required gas by the amount of transactional data
	if len(data) > 0 {
		// Zero and non-zero bytes are priced differently
		var nz uint64
		for _, byt := range data {
			if byt != 0 {
				nz++
			}
		}
		// Make sure we don't exceed uint64 for all data combinations
		if (math.MaxUint64-gas)/params.TxDataNonZeroGas < nz {
			return 0, ErrGasUintOverflow
		}
		gas += nz * params.TxDataNonZeroGas

		z := uint64(len(data)) - nz
		if (math.MaxUint64-gas)/params.TxDataZeroGas < z {
			return 0, ErrGasUintOverflow
		}
		gas += z * params.TxDataZeroGas
	}
	return gas, nil
}

// NewStateTransition initialises and returns a new state transition object.
func NewStateTransition(evm *vm.EVM, msg *Message, gp *GasPool) *StateTransition {
	return &StateTransition{
		gp:       gp,
		evm:      evm,
		msg:      msg,
		gasPrice: msg.GasPrice,
		value:    msg.Value,
		data:     msg.Data,
		state:    evm.StateDB,
	}
}

// TransactionToMessage converts a transaction into a Message.
func TransactionToMessage(tx *types.Transaction, s types.Signer, balanceFee *big.Int, number *big.Int) (*Message, error) {
	from, err := types.Sender(s, tx)
	msg := &Message{
		From:              from,
		Nonce:             tx.Nonce(),
		GasLimit:          tx.Gas(),
		GasPrice:          new(big.Int).Set(tx.GasPrice()),
		To:                tx.To(),
		Value:             tx.Value(),
		Data:              tx.Data(),
		PmAddress:         from,
		PmPayload:         tx.PmPayload(),
		SkipAccountChecks: false,
		BalanceTokenFee:   balanceFee,
	}
	if len(msg.PmPayload) >= 20 {
		msg.PmAddress = common.BytesToAddress(msg.PmPayload[:20]) // the first 20 bytes of PmPayload is the address of Paymaster contract
	}
	if balanceFee != nil {
		if number.Cmp(common.TIPTRC21Fee) > 0 {
			msg.GasPrice = common.TRC21GasPrice
		} else {
			msg.GasPrice = common.TRC21GasPriceBefore
		}
	}
	return msg, err
}

// ApplyMessage computes the new state by applying the given message
// against the old state within the environment.
//
// ApplyMessage returns the bytes returned by any EVM execution (if it took place),
// the gas used (which includes gas refunds) and an error if it failed. An error always
// indicates a core error meaning that the message would always fail for that particular
// state and would never be accepted within a block.
func ApplyMessage(evm *vm.EVM, msg *Message, gp *GasPool, owner common.Address) (*ExecutionResult, error) {
	return NewStateTransition(evm, msg, gp).TransitionDb(owner)
}

func (st *StateTransition) from() vm.AccountRef {
	f := st.msg.From
	if !st.state.Exist(f) {
		st.state.CreateAccount(f)
	}
	return vm.AccountRef(f)
}

func (st *StateTransition) pmAddress() vm.AccountRef {
	f := st.msg.PmAddress
	if !st.state.Exist(f) {
		st.state.CreateAccount(f)
	}
	return vm.AccountRef(f)
}

func (st *StateTransition) balanceTokenFee() *big.Int {
	return st.msg.BalanceTokenFee
}

func (st *StateTransition) to() vm.AccountRef {
	if st.msg == nil {
		return vm.AccountRef{}
	}
	to := st.msg.To
	if to == nil {
		return vm.AccountRef{} // contract creation
	}

	reference := vm.AccountRef(*to)
	if !st.state.Exist(*to) {
		st.state.CreateAccount(*to)
	}
	return reference
}

func (st *StateTransition) buyGas() error {
	var (
		state           = st.state
		balanceTokenFee = st.balanceTokenFee()
		from            = st.pmAddress()
	)
	mgval := new(big.Int).Mul(new(big.Int).SetUint64(st.msg.GasLimit), st.gasPrice)
	if balanceTokenFee == nil {
		if state.GetBalance(from.Address()).Cmp(mgval) < 0 {
			return errInsufficientBalanceForGas
		}
	} else if balanceTokenFee.Cmp(mgval) < 0 {
		return errInsufficientBalanceForGas
	}
	if err := st.gp.SubGas(st.msg.GasLimit); err != nil {
		return err
	}
	st.gas = st.msg.GasLimit - st.gas // gas left is gasLimit minus gas used by validateAndPayForPaymaster

	st.initialGas = st.msg.GasLimit
	if balanceTokenFee == nil {
		state.SubBalance(from.Address(), mgval)
	}
	return nil
}

func (st *StateTransition) preCheck() error {
	// Only check transactions that are not fake
	msg := st.msg
	if !msg.SkipAccountChecks {
		// Make sure this transaction's nonce is correct.
		stNonce := st.state.GetNonce(msg.From)
		if msgNonce := msg.Nonce; stNonce < msgNonce {
			return fmt.Errorf("%w: address %v, tx: %d state: %d", ErrNonceTooHigh,
				msg.From.Hex(), msgNonce, stNonce)
		} else if stNonce > msgNonce {
			return fmt.Errorf("%w: address %v, tx: %d state: %d", ErrNonceTooLow,
				msg.From.Hex(), msgNonce, stNonce)
		} else if stNonce+1 < stNonce {
			return fmt.Errorf("%w: address %v, nonce: %d", ErrNonceMax,
				msg.From.Hex(), stNonce)
		}
		// Make sure the sender is an EOA
		codeHash := st.state.GetCodeHash(msg.From)
		if codeHash != (common.Hash{}) && codeHash != types.EmptyCodeHash {
			return fmt.Errorf("%w: address %v, codehash: %s", ErrSenderNoEOA,
				msg.From.Hex(), codeHash)
		}
	}

	return st.buyGas()
}

// TransitionDb will transition the state by applying the current message and
// returning the evm execution result with following fields.
//
//   - used gas:
//     total gas used (including gas being refunded)
//   - returndata:
//     the returned data from evm
//   - concrete execution error:
//     various **EVM** error which aborts the execution,
//     e.g. ErrOutOfGas, ErrExecutionReverted
//
// However if any consensus issue encountered, return the error directly with
// nil evm execution result.
func (st *StateTransition) TransitionDb(owner common.Address) (*ExecutionResult, error) {
	// First check this message satisfies all consensus rules before
	// applying the message. The rules include these clauses
	//
	// 1. the nonce of the message caller is correct
	// 2. caller has enough balance to cover transaction fee(gaslimit * gasprice)
	// 3. the amount of gas required is available in the block
	// 4. the purchased gas is enough to cover intrinsic usage
	// 5. there is no overflow when calculating intrinsic gas
	// 6. caller has enough balance to cover asset transfer for **topmost** call

	// paymaster validate and pay
	var (
		pmMagic   [4]byte
		pmContext []byte
		pmGasUsed uint64
	)
	if len(st.msg.PmPayload) > 0 {
		magic, context, gasUsed, err := validateAndPayForPaymaster(st.msg, st.evm, &paymaster.IPaymasterTransaction{From: st.msg.From}, common.BytesToHash(st.msg.PmPayload[20:]))
		st.gas += gasUsed // mark gas used by validating
		pmContext = context
		pmMagic = magic
		fmt.Printf("@@@@@@@@@@@@@@ validate: %v, %v, %d, %v\n", magic, context, gasUsed, err)
		if isValidMagic(magic) && err == nil {
			if len(st.msg.PmPayload) >= 20 {
				st.msg.PmAddress = common.BytesToAddress(st.msg.PmPayload[:20]) // the first 20 bytes of PmPayload is the address of Paymaster contract
			} else {
				copy(magic[:], []byte(invalidMagic))
				fmt.Println("@@@@@@@@@@@@@ magic", magic)
			}
		}
	}

	// Check clauses 1-3, buy gas if everything is correct
	if err := st.preCheck(); err != nil {
		return nil, err
	}
	msg := st.msg
	sender := st.from() // err checked in preCheck

	homestead := st.evm.ChainConfig().IsHomestead(st.evm.BlockNumber)
	contractCreation := msg.To == nil

	// Check clauses 4-5, subtract intrinsic gas if everything is correct
	gas, err := IntrinsicGas(st.data, contractCreation, homestead)
	if err != nil {
		return nil, err
	}
	if st.gas < gas {
		return nil, ErrIntrinsicGas
	}
	st.gas -= gas

	// check clause 6
	if msg.Value.Sign() > 0 && !st.evm.CanTransfer(st.state, msg.From, msg.Value) {
		return nil, ErrInsufficientFundsForTransfer
	}

	var (
		ret   []byte
		vmerr error
	)
	// for debugging purpose
	// TODO: clean it after fixing the issue https://github.com/tomochain/tomochain/issues/401
	nonce := uint64(1)
	if contractCreation {
		ret, _, st.gas, vmerr = st.evm.Create(sender, st.data, st.gas, st.value)
	} else {
		// Increment the nonce for the next transaction
		nonce = st.state.GetNonce(sender.Address()) + 1
		st.state.SetNonce(sender.Address(), nonce)
		ret, st.gas, vmerr = st.evm.Call(sender, st.to().Address(), st.data, st.gas, st.value)
	}

	// paymaster post transaction
	if len(st.msg.PmPayload) > 0 && isValidMagic(pmMagic) {
		pmGasUsed, err = postTransaction(st.msg, st.evm, &paymaster.IPaymasterTransaction{From: st.msg.From}, common.BytesToHash(st.msg.PmPayload[20:]),
			pmContext, 0, &paymaster.IPaymasterExecutionResult{Success: true})
		fmt.Printf("@@@@@@@@@@@@@@ postTransaction: %d %v\n", pmGasUsed, err)
		// not enough for postTransaction execution
		if st.gas-pmGasUsed > st.gas {
			err = ErrPostTransactionOutOfGas
		}
	}
	st.refundGas()

	if st.evm.BlockNumber.Cmp(common.TIPTRC21Fee) > 0 {
		if (owner != common.Address{}) {
			st.state.AddBalance(owner, new(big.Int).Mul(new(big.Int).SetUint64(st.gasUsed()), st.gasPrice))
		}
	} else {
		st.state.AddBalance(st.evm.Coinbase, new(big.Int).Mul(new(big.Int).SetUint64(st.gasUsed()), st.gasPrice))
	}

	return &ExecutionResult{
		UsedGas:    st.gasUsed(),
		Err:        vmerr,
		ReturnData: ret,
	}, err
}

func (st *StateTransition) refundGas() {
	// Apply refund counter, capped to half of the used gas.
	refund := st.gasUsed() / 2
	if refund > st.state.GetRefund() {
		refund = st.state.GetRefund()
	}
	st.gas += refund

	balanceTokenFee := st.balanceTokenFee()
	if balanceTokenFee == nil {
		from := st.pmAddress()
		// Return ETH for remaining gas, exchanged at the original rate.
		remaining := new(big.Int).Mul(new(big.Int).SetUint64(st.gas), st.gasPrice)
		st.state.AddBalance(from.Address(), remaining)
	}
	// Also return remaining gas to the block gas counter so it is
	// available for the next transaction.
	st.gp.AddGas(st.gas)
}

// gasUsed returns the amount of gas used up by the state transition.
func (st *StateTransition) gasUsed() uint64 {
	return st.initialGas - st.gas
}

func isValidMagic(magic [4]byte) bool {
	return string(magic[:]) != invalidMagic
}
