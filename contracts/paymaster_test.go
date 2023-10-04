package contracts

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tomochain/tomochain/accounts/abi/bind"
	"github.com/tomochain/tomochain/accounts/abi/bind/backends"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/contracts/paymaster"
	"github.com/tomochain/tomochain/core"
	"github.com/tomochain/tomochain/core/types"
	"github.com/tomochain/tomochain/crypto"
	"github.com/tomochain/tomochain/params"
)

var (
	pKey            *ecdsa.PrivateKey
	addr            common.Address
	pmAddress       common.Address
	pmInstance      *paymaster.Paymaster
	contractBackend *backends.SimulatedBackend
	transactOpts    *bind.TransactOpts

	ctx     = context.Background()
	chainID = big.NewInt(1337)
	nonce   = uint64(1)
)

func init() {
	pKey, _ = crypto.GenerateKey()
	addr = crypto.PubkeyToAddress(pKey.PublicKey)
	contractBackend = backends.NewSimulatedBackend(core.GenesisAlloc{addr: {Balance: big.NewInt(1_000_000_000_000_000_000)}})
	transactOpts = bind.NewKeyedTransactor(pKey)
	var err error
	pmAddress, _, pmInstance, err = paymaster.DeployPaymaster(transactOpts, contractBackend)
	if err != nil {
		panic(err)
	}
	contractBackend.Commit()
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: common.MinGasPrice,
		Gas:      params.TxGas + 1000,
		To:       &pmAddress,
		Value:    big.NewInt(1_000_000_000_000_000),
		Data:     nil,
	})
	nonce++
	signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, pKey)
	if err != nil {
		panic(err)
	}
	err = contractBackend.SendTransaction(ctx, signedTx)
	if err != nil {
		panic(err)
	}
	contractBackend.Commit()
}

func TestSuccessPaymasterTx(t *testing.T) {
	addrBalance, err := contractBackend.BalanceAt(ctx, addr, nil)
	pmBalance, err := contractBackend.BalanceAt(ctx, pmAddress, nil)
	tx := types.NewTx(&types.PaymasterTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasPrice:  common.MinGasPrice,
		Gas:       200000,
		To:        &addr,
		Value:     nil,
		Data:      nil,
		PmPayload: append(pmAddress.Bytes(), common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000001")...),
	})
	nonce++
	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), pKey)
	if err != nil {
		t.Error(err)
	}
	err = contractBackend.SendTransaction(ctx, signedTx)
	if err != nil {
		t.Error(err)
	}
	contractBackend.Commit()
	newAddrBalance, err := contractBackend.BalanceAt(ctx, addr, nil)
	newPmBalance, err := contractBackend.BalanceAt(ctx, pmAddress, nil)
	assert.Less(t, newPmBalance.Uint64(), pmBalance.Uint64(), "Balance of the paymaster contract must be decreased")
	assert.Exactly(t, addrBalance.Uint64(), newAddrBalance.Uint64(), "Balance of the from address must not be changed")
}
