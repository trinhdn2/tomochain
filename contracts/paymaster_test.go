package contracts

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"testing"

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
		Nonce:    1,
		GasPrice: common.MinGasPrice,
		Gas:      params.TxGas + 1000,
		To:       &pmAddress,
		Value:    big.NewInt(1_000_000_000_000_000),
		Data:     nil,
	})
	signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, pKey)
	if err != nil {
		panic(err)
	}
	err = contractBackend.SendTransaction(ctx, signedTx)
	if err != nil {
		panic(err)
	}
	contractBackend.Commit()
	fmt.Println("addr:", addr.Hex(), "\nPmAddress:", pmAddress.Hex())
}

func TestSuccessPaymasterTx(t *testing.T) {
	balance, err := contractBackend.BalanceAt(ctx, addr, nil)
	fmt.Println(balance.String())
	tx := types.NewTx(&types.PaymasterTx{
		ChainID:   chainID,
		Nonce:     2,
		GasPrice:  common.MinGasPrice,
		Gas:       150000,
		To:        &addr,
		Value:     nil,
		Data:      nil,
		PmPayload: append(pmAddress.Bytes(), common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000001")...),
	})
	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), pKey)
	if err != nil {
		t.Error(err)
	}
	err = contractBackend.SendTransaction(ctx, signedTx)
	if err != nil {
		t.Error(err)
	}
	contractBackend.Commit()
	receipt, err := contractBackend.TransactionReceipt(ctx, signedTx.Hash())
	fmt.Println(receipt)
	balance, err = contractBackend.BalanceAt(ctx, addr, nil)
	fmt.Println(balance.String())
}

//0x
//d16a2f95
//0000000000000000000000000000000000000000000000000000000000000000
//000000000000000000000000f7e6258432cda2b44b013d6b67ced090ec4bf78f
//
//0xd16a2f950000000000000000000000000000000000000000000000000000000000000000000000000000000000000000f7e6258432cda2b44b013d6b67ced090ec4bf78f
