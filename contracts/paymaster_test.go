package contracts

import (
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/tomochain/tomochain/accounts/abi/bind"
	"github.com/tomochain/tomochain/accounts/abi/bind/backends"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/contracts/paymaster"
	"github.com/tomochain/tomochain/core"
	"github.com/tomochain/tomochain/crypto"
)

var (
	pKey            *ecdsa.PrivateKey
	addr            common.Address
	pmAddress       common.Address
	pmInstance      *paymaster.Paymaster
	contractBackend *backends.SimulatedBackend
	transactOpts    *bind.TransactOpts
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
}

func TestAccumulateGasUsed(t *testing.T) {

}
