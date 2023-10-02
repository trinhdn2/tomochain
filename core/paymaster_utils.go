package core

import (
	"fmt"
	"strings"

	"github.com/tomochain/tomochain/accounts/abi"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/contracts/paymaster"
	"github.com/tomochain/tomochain/core/vm"
)

var IPaymasterABI abi.ABI

func init() {
	var err error
	IPaymasterABI, err = abi.JSON(strings.NewReader(paymaster.PaymasterMetaData.ABI))
	if err != nil {
		panic(fmt.Sprintf("Error reading abi: %v", err))
	}
}

func validateAndPayForPaymaster(originMsg *Message, evm *vm.EVM, tx *paymaster.IPaymasterTransaction,
	txHash common.Hash) ([4]byte, []byte, uint64, error) {
	payload, err := IPaymasterABI.Pack("validateAndPayForPaymasterTransaction", txHash, tx)
	if err != nil {
		return [4]byte{}, nil, 0, err
	}
	ret, usedGas, err := constructAndApplySmcCallMsg(originMsg, evm, payload)
	if err != nil {
		return [4]byte{}, nil, 0, err
	}
	// unpack result
	var validateResult struct {
		magic   [4]byte
		context []byte
	}
	err = IPaymasterABI.Unpack(&validateResult, "validateAndPayForPaymasterTransaction", ret)
	if err != nil {
		return [4]byte{}, nil, 0, err
	}
	return validateResult.magic, validateResult.context, usedGas, nil
}

func postTransaction(originMsg *Message, evm *vm.EVM, tx *paymaster.IPaymasterTransaction, txHash common.Hash,
	context []byte, maxRefundedGas uint64, txResult *paymaster.IPaymasterExecutionResult) (uint64, error) {
	payload, err := IPaymasterABI.Pack("postTransaction", context, txHash, tx, txResult, maxRefundedGas)
	if err != nil {
		return 0, err
	}
	_, usedGas, err := constructAndApplySmcCallMsg(originMsg, evm, payload)
	if err != nil {
		return 0, err
	}
	return usedGas, nil
}

func constructAndApplySmcCallMsg(originMsg *Message, evm *vm.EVM, data []byte) ([]byte, uint64, error) {
	validateMsg := Message{
		To:                &originMsg.PmAddress,
		From:              originMsg.From,
		Nonce:             0,
		Value:             nil,
		GasLimit:          originMsg.GasLimit,
		GasPrice:          originMsg.GasPrice,
		Data:              data,
		PmAddress:         originMsg.PmAddress,
		PmPayload:         originMsg.PmPayload,
		BalanceTokenFee:   nil,
		SkipAccountChecks: true,
	}
	return apply(evm, validateMsg)
}

// apply the sub message on top of current EVM, returns the byte result, used gas and VM error
func apply(evm *vm.EVM, msg Message) ([]byte, uint64, error) {
	sender := vm.AccountRef(msg.From)
	ret, leftOverGas, vmErr := evm.Call(sender, *msg.To, msg.Data, msg.GasLimit, msg.Value)
	return ret, msg.GasLimit - leftOverGas, vmErr
}
