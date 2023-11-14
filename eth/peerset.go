// Copyright 2020 The go-ethereum Authors
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

package eth

import (
	"errors"
	"math/big"
	"sync"
	"time"

	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/eth/protocols/eth"
	"github.com/tomochain/tomochain/event"
	"github.com/tomochain/tomochain/p2p"
)

var (
	// errPeerSetClosed is returned if a peer is attempted to be added or removed
	// from the peer set after it has been terminated.
	errPeerSetClosed = errors.New("peerset closed")

	// errPeerAlreadyRegistered is returned if a peer is attempted to be added
	// to the peer set, but one with the same id already exists.
	errPeerAlreadyRegistered = errors.New("peer already registered")

	// errPeerNotRegistered is returned if a peer is attempted to be removed from
	// a peer set, but no peer with the given id exists.
	errPeerNotRegistered = errors.New("peer not registered")

	// ethConnectTimeout is the `snap` timeout for `eth` to connect too.
	ethConnectTimeout = 3 * time.Second
)

// peerSet represents the collection of active peers currently participating in
// the `eth` or `snap` protocols.
type peerSet struct {
	ethPeers map[string]*ethPeer // Peers connected on the `eth` protocol

	ethJoinFeed  event.Feed // Events when an `eth` ethPeer successfully joins
	ethDropFeed  event.Feed // Events when an `eth` ethPeer gets dropped
	snapJoinFeed event.Feed // Events when a `snap` ethPeer joins on both `eth` and `snap`
	snapDropFeed event.Feed // Events when a `snap` ethPeer gets dropped (only if fully joined)

	scope event.SubscriptionScope // Subscription group to unsubscribe everyone at once

	lock   sync.RWMutex
	closed bool
}

// newPeerSet creates a new peer set to track the active participants.
func newPeerSet() *peerSet {
	return &peerSet{
		ethPeers: make(map[string]*ethPeer),
	}
}

// subscribeEthJoin registers a subscription for peers joining (and completing
// the handshake) on the `eth` protocol.
func (ps *peerSet) subscribeEthJoin(ch chan<- *eth.Peer) event.Subscription {
	return ps.scope.Track(ps.ethJoinFeed.Subscribe(ch))
}

// subscribeEthDrop registers a subscription for peers being dropped from the
// `eth` protocol.
func (ps *peerSet) subscribeEthDrop(ch chan<- *eth.Peer) event.Subscription {
	return ps.scope.Track(ps.ethDropFeed.Subscribe(ch))
}

// registerEthPeer injects a new `eth` peer into the working set, or returns an
// error if the peer is already known. The peer is announced on the `eth` join
// feed and if it completes a pending `snap` peer, also on that feed.
func (ps *peerSet) registerEthPeer(peer *eth.Peer) error {
	ps.lock.Lock()
	if ps.closed {
		ps.lock.Unlock()
		return errPeerSetClosed
	}
	id := peer.ID()
	if _, ok := ps.ethPeers[id]; ok {
		ps.lock.Unlock()
		return errPeerAlreadyRegistered
	}
	ps.ethPeers[id] = &ethPeer{Peer: peer}

	return nil
}

// unregisterEthPeer removes a remote peer from the active set, disabling any further
// actions to/from that particular entity. The drop is announced on the `eth` drop
// feed and also on the `snap` feed if the eth/snap duality was broken just now.
func (ps *peerSet) unregisterEthPeer(id string) error {
	ps.lock.Lock()
	eth, ok := ps.ethPeers[id]
	if !ok {
		ps.lock.Unlock()
		return errPeerNotRegistered
	}
	delete(ps.ethPeers, id)
	ps.ethDropFeed.Send(eth)
	return nil
}

// peer retrieves the registered `eth` peer with the given id.
func (ps *peerSet) ethPeer(id string) *ethPeer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	return ps.ethPeers[id]
}

// ethPeersWithoutBlock retrieves a list of `eth` peers that do not have a given
// block in their set of known hashes so it might be propagated to them.
func (ps *peerSet) ethPeersWithoutBlock(hash common.Hash) []*ethPeer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]*ethPeer, 0, len(ps.ethPeers))
	for _, p := range ps.ethPeers {
		if !p.KnownBlock(hash) {
			list = append(list, p)
		}
	}
	return list
}

// ethPeersWithoutTransacion retrieves a list of `eth` peers that do not have a
// given transaction in their set of known hashes.
func (ps *peerSet) ethPeersWithoutTransacion(hash common.Hash) []*ethPeer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]*ethPeer, 0, len(ps.ethPeers))
	for _, p := range ps.ethPeers {
		if !p.KnownTransaction(hash) {
			list = append(list, p)
		}
	}
	return list
}

// Len returns if the current number of `eth` peers in the set. Since the `snap`
// peers are tied to the existence of an `eth` connection, that will always be a
// subset of `eth`.
func (ps *peerSet) Len() int {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	return len(ps.ethPeers)
}

// ethPeerWithHighestTD retrieves the known peer with the currently highest total
// difficulty.
func (ps *peerSet) ethPeerWithHighestTD() *eth.Peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	var (
		bestPeer *eth.Peer
		bestTd   *big.Int
	)
	for _, p := range ps.ethPeers {
		if _, td := p.Head(); bestPeer == nil || td.Cmp(bestTd) > 0 {
			bestPeer, bestTd = p.Peer, td
		}
	}
	return bestPeer
}

// close disconnects all peers.
func (ps *peerSet) close() {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	for _, p := range ps.ethPeers {
		p.Disconnect(p2p.DiscQuitting)
	}
	ps.closed = true
}
