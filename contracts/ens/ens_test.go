// Copyright 2016 The go-ethereum Authors
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

package ens

import (
	"context"
	"math/big"
	"testing"

	"github.com/tomochain/tomochain/accounts/abi/bind"
	"github.com/tomochain/tomochain/accounts/abi/bind/backends"
	"github.com/tomochain/tomochain/contracts/ens/contract"
	"github.com/tomochain/tomochain/core"
	"github.com/tomochain/tomochain/crypto"
)

var (
	key, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	name   = "my name on ENS"
	hash   = crypto.Keccak256Hash([]byte("my content"))
	addr   = crypto.PubkeyToAddress(key.PublicKey)
)

func TestENS(t *testing.T) {
	contractBackend := backends.NewSimulatedBackend(core.GenesisAlloc{addr: {Balance: big.NewInt(10_000_000_000_000_000)}})
	chainID, err := contractBackend.ChainID(context.Background())
	if err != nil {
		t.Fatalf("can't get chainID: %v", err)
	}
	transactOpts, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		t.Fatalf("can't create TransactOpts: %v", err)
	}

	ensAddr, ens, err := DeployENS(transactOpts, contractBackend)
	if err != nil {
		t.Fatalf("can't deploy root registry: %v", err)
	}
	contractBackend.Commit()

	// Set ourself as the owner of the name.
	if _, err := ens.Register(name); err != nil {
		t.Fatalf("can't register: %v", err)
	}
	contractBackend.Commit()

	// Deploy a resolver and make it responsible for the name.
	resolverAddr, _, _, err := contract.DeployPublicResolver(transactOpts, contractBackend, ensAddr)
	if err != nil {
		t.Fatalf("can't deploy resolver: %v", err)
	}
	if _, err := ens.SetResolver(ensNode(name), resolverAddr); err != nil {
		t.Fatalf("can't set resolver: %v", err)
	}
	contractBackend.Commit()

	// Set the content hash for the name.
	if _, err = ens.SetContentHash(name, hash); err != nil {
		t.Fatalf("can't set content hash: %v", err)
	}
	contractBackend.Commit()

	// Try to resolve the name.
	vhost, err := ens.Resolve(name)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if vhost != hash {
		t.Fatalf("resolve error, expected %v, got %v", hash.Hex(), vhost.Hex())
	}
}
