// Copyright 2015 The go-ethereum Authors
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
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/consensus"
	"github.com/tomochain/tomochain/core"
	"github.com/tomochain/tomochain/core/types"
	"github.com/tomochain/tomochain/eth/downloader"
	"github.com/tomochain/tomochain/eth/fetcher"
	"github.com/tomochain/tomochain/eth/protocols/eth"
	"github.com/tomochain/tomochain/ethdb"
	"github.com/tomochain/tomochain/event"
	"github.com/tomochain/tomochain/log"
	"github.com/tomochain/tomochain/p2p"
	"github.com/tomochain/tomochain/params"
)

const (
	softResponseLimit = 2 * 1024 * 1024 // Target maximum size of returned blocks, headers or node data.
	estHeaderRlpSize  = 500             // Approximate size of an RLP encoded block header

	// txChanSize is the size of channel listening to TxPreEvent.
	// The number is referenced from the size of tx pool.
	txChanSize = 4096
)

var (
	daoChallengeTimeout = 15 * time.Second // Time allowance for a node to reply to the DAO handshake challenge
)

// errIncompatibleConfig is returned if the requested protocols and configs are
// not compatible (low protocol version restrictions and high requirements).
var errIncompatibleConfig = errors.New("incompatible configuration")

func errResp(code errCode, format string, v ...interface{}) error {
	return fmt.Errorf("%v - %v", code, fmt.Sprintf(format, v...))
}

type handler struct {
	networkId uint64

	fastSync  uint32 // Flag whether fast sync is enabled (gets disabled if we already have blocks)
	acceptTxs uint32 // Flag whether we're considered synchronised (enables transaction processing)

	database    ethdb.Database
	txpool      txPool
	orderpool   orderPool
	lendingpool lendingPool
	chainconfig *params.ChainConfig
	maxPeers    int

	downloader   *downloader.Downloader
	chain        *core.BlockChain
	txFetcher    *fetcher.TxFetcher
	blockFetcher *fetcher.BlockFetcher
	peers        *peerSet

	eventMux      *event.TypeMux
	txCh          chan core.TxPreEvent
	orderTxCh     chan core.OrderTxPreEvent
	lendingTxCh   chan core.LendingTxPreEvent
	txSub         event.Subscription
	orderTxSub    event.Subscription
	lendingTxSub  event.Subscription
	minedBlockSub *event.TypeMuxSubscription

	// channels for fetcher, syncer, txsyncLoop
	newPeerCh   chan *eth.Peer
	txsyncCh    chan *txsync
	quitSync    chan struct{}
	noMorePeers chan struct{}

	// wait group is used for graceful shutdowns during downloading
	// and processing
	chainSync *chainSyncer
	wg        sync.WaitGroup
	peerWG    sync.WaitGroup

	handlerStartCh chan struct{}
	handlerDoneCh  chan struct{}
}

// NewProtocolManagerEx add order pool to protocol
func NewProtocolManagerEx(config *params.ChainConfig, mode downloader.SyncMode, networkID uint64, mux *event.TypeMux, txpool txPool, orderpool orderPool, lendingpool lendingPool, engine consensus.Engine, blockchain *core.BlockChain, chaindb ethdb.Database) (*handler, error) {
	protocol, err := NewProtocolManager(config, mode, networkID, mux, txpool, engine, blockchain, chaindb)
	if err != nil {
		return nil, err
	}
	protocol.addOrderPoolProtocol(orderpool)
	protocol.addLendingPoolProtocol(lendingpool)
	return protocol, nil
}

// NewProtocolManager returns a new ethereum sub protocol manager. The Ethereum sub protocol manages peers capable
// with the ethereum network.
func NewProtocolManager(config *params.ChainConfig, mode downloader.SyncMode, networkID uint64, mux *event.TypeMux, txpool txPool, engine consensus.Engine, blockchain *core.BlockChain, chaindb ethdb.Database) (*handler, error) {
	// Create the protocol handler with the base fields
	h := &handler{
		networkId:      networkID,
		eventMux:       mux,
		database:       chaindb,
		txpool:         txpool,
		chain:          blockchain,
		chainconfig:    config,
		peers:          newPeerSet(),
		newPeerCh:      make(chan *eth.Peer),
		noMorePeers:    make(chan struct{}),
		txsyncCh:       make(chan *txsync),
		quitSync:       make(chan struct{}),
		orderpool:      nil,
		lendingpool:    nil,
		orderTxSub:     nil,
		lendingTxSub:   nil,
		handlerDoneCh:  make(chan struct{}),
		handlerStartCh: make(chan struct{}),
	}
	// Figure out whether to allow fast sync or not
	if mode == downloader.FastSync && blockchain.CurrentBlock().NumberU64() > 0 {
		log.Warn("Blockchain not empty, fast sync disabled")
		mode = downloader.FullSync
	}
	if mode == downloader.FastSync {
		h.fastSync = uint32(1)
	}
	// Construct the different synchronisation mechanisms
	h.downloader = downloader.New(mode, chaindb, h.eventMux, blockchain, nil, h.removePeer)

	validator := func(header *types.Header) error {
		return engine.VerifyHeader(blockchain, header, true)
	}
	heighter := func() uint64 {
		return blockchain.CurrentBlock().NumberU64()
	}
	inserter := func(blocks types.Blocks) (int, error) {
		// If fast sync is running, deny importing weird blocks. This is a problematic
		// clause when starting up a new network, because fast-syncing miners might not
		// accept each others' blocks until a restart. Unfortunately we haven't figured
		// out a way yet where nodes can decide unilaterally whether the network is new
		// or not. This should be fixed if we figure out a solution.
		if atomic.LoadUint32(&h.fastSync) == 1 {
			log.Warn("Fast syncing, discarded propagated block", "number", blocks[0].Number(), "hash", blocks[0].Hash())
			return 0, nil
		}
		n, err := h.chain.InsertChain(blocks)
		if err == nil {
			atomic.StoreUint32(&h.acceptTxs, 1) // Mark initial sync done on any fetcher import
		}
		return n, err
	}
	h.blockFetcher = fetcher.NewBlockFetcher(false, nil, blockchain.GetBlockByHash, validator, h.BroadcastBlock, heighter, nil, inserter, h.removePeer)
	fetchTx := func(peer string, hashes []common.Hash) error {
		p := h.peers.peer(peer)
		if p == nil {
			return errors.New("unknown peer")
		}
		return p.RequestTxs(hashes)
	}
	addTxs := func(txs []*types.Transaction) []error {
		return h.txpool.Add(txs, false, false)
	}
	h.txFetcher = fetcher.NewTxFetcher(h.txpool.Has, addTxs, fetchTx, h.removePeer)
	h.chainSync = newChainSyncer(h)
	return h, nil
}

// protoTracker tracks the number of active protocol handlers.
func (h *handler) protoTracker() {
	defer h.wg.Done()
	var active int
	for {
		select {
		case <-h.handlerStartCh:
			active++
		case <-h.handlerDoneCh:
			active--
		case <-h.quitSync:
			// Wait for all active handlers to finish.
			for ; active > 0; active-- {
				<-h.handlerDoneCh
			}
			return
		}
	}
}

// incHandlers signals to increment the number of active handlers if not
// quitting.
func (h *handler) incHandlers() bool {
	select {
	case h.handlerStartCh <- struct{}{}:
		return true
	case <-h.quitSync:
		return false
	}
}

// decHandlers signals to decrement the number of active handlers.
func (h *handler) decHandlers() {
	h.handlerDoneCh <- struct{}{}
}

// runEthPeer
func (h *handler) runEthPeer(peer *eth.Peer, handler eth.Handler) error {
	if !h.chainSync.handlePeerEvent(peer) {
		return p2p.DiscQuitting
	}
	h.peerWG.Add(1)
	defer h.peerWG.Done()

	// Execute the Ethereum handshake
	var (
		genesis = h.chain.Genesis()
		head    = h.chain.CurrentHeader()
		hash    = head.Hash()
		number  = head.Number.Uint64()
		td      = h.chain.GetTd(hash, number)
	)
	if err := peer.Handshake(h.networkId, td, hash, genesis.Hash()); err != nil {
		peer.Log().Debug("Ethereum handshake failed", "err", err)
		return err
	}
	// Ignore maxPeers if this is a trusted peer
	if h.peers.Len() >= h.maxPeers && !peer.Peer.Info().Network.Trusted {
		return p2p.DiscTooManyPeers
	}
	peer.Log().Debug("Ethereum peer connected", "name", peer.Name())

	// Register the peer locally
	if err := h.peers.registerEthPeer(peer); err != nil {
		peer.Log().Error("Ethereum peer registration failed", "err", err)
		return err
	}
	defer h.removePeer(peer.ID())

	p := h.peers.peer(peer.ID())
	if p == nil {
		return errors.New("peer dropped during handling")
	}
	// Register the peer in the downloader. If the downloader considers it banned, we disconnect
	if err := h.downloader.RegisterPeer(peer.ID(), peer.Version(), peer); err != nil {
		return err
	}
	h.chainSync.handlePeerEvent(peer)
	// Propagate existing transactions. new transactions appearing
	// after this will be sent via broadcasts.
	h.syncTransactions(peer)
	// Handle incoming messages until the connection is torn down
	return handler(peer)
}

func (h *handler) addOrderPoolProtocol(orderpool orderPool) {
	h.orderpool = orderpool
}

func (h *handler) addLendingPoolProtocol(lendingpool lendingPool) {
	h.lendingpool = lendingpool
}

// removePeer requests disconnection of a peer.
func (h *handler) removePeer(id string) {
	peer := h.peers.peer(id)
	if peer != nil {
		peer.Peer.Disconnect(p2p.DiscUselessPeer)
	}
}

func (h *handler) Start(maxPeers int) {
	h.maxPeers = maxPeers

	// broadcast transactions
	h.txCh = make(chan core.TxPreEvent, txChanSize)
	h.txSub = h.txpool.SubscribeTxPreEvent(h.txCh)

	h.orderTxCh = make(chan core.OrderTxPreEvent, txChanSize)
	if h.orderpool != nil {
		h.orderTxSub = h.orderpool.SubscribeTxPreEvent(h.orderTxCh)
	}
	h.lendingTxCh = make(chan core.LendingTxPreEvent, txChanSize)
	if h.lendingpool != nil {
		h.lendingTxSub = h.lendingpool.SubscribeTxPreEvent(h.lendingTxCh)
	}
	go h.txBroadcastLoop()
	go h.orderTxBroadcastLoop()
	go h.lendingTxBroadcastLoop()
	// broadcast mined blocks
	h.minedBlockSub = h.eventMux.Subscribe(core.NewMinedBlockEvent{})
	go h.minedBroadcastLoop()

	// start sync handlers
	h.wg.Add(1)
	go h.chainSync.loop()

	h.wg.Add(1)
	go h.txsyncLoop()

	// start peer handler tracker
	h.wg.Add(1)
	go h.protoTracker()
}

func (h *handler) Stop() {
	log.Info("Stopping Ethereum protocol")

	h.txSub.Unsubscribe() // quits txBroadcastLoop
	if h.orderTxSub != nil {
		h.orderTxSub.Unsubscribe()
	}
	if h.lendingTxSub != nil {
		h.lendingTxSub.Unsubscribe()
	}
	h.minedBlockSub.Unsubscribe() // quits blockBroadcastLoop

	// Quit the sync loop.
	// After this send has completed, no new peers will be accepted.
	h.noMorePeers <- struct{}{}

	// Quit fetcher, txsyncLoop.
	close(h.quitSync)

	// Disconnect existing sessions.
	// This also closes the gate for any new registrations on the peer set.
	// sessions which are already established but not added to pm.peers yet
	// will exit when they try to register.
	h.peers.close()

	// Wait for all peer handler goroutines and the loops to come down.
	h.wg.Wait()

	log.Info("Ethereum protocol stopped")
}

/*
// handle is the callback invoked to manage the life cycle of an eth peer. When
// this function terminates, the peer is disconnected.
func (h *handler) handle(p *eth.Peer) error {
	// Ignore maxPeers if this is a trusted peer
	if h.peers.Len() >= h.maxPeers && !p.Peer.Info().Network.Trusted {
		return p2p.DiscTooManyPeers
	}
	p.Log().Debug("Ethereum peer connected", "name", p.Name())

	// Execute the Ethereum handshake
	var (
		genesis = h.chain.Genesis()
		head    = h.chain.CurrentHeader()
		hash    = head.Hash()
		number  = head.Number.Uint64()
		td      = h.chain.GetTd(hash, number)
	)
	if err := p.Handshake(h.networkId, td, hash, genesis.Hash()); err != nil {
		p.Log().Debug("Ethereum handshake failed", "err", err)
		return err
	}
	if rw, ok := p.rw.(*meteredMsgReadWriter); ok {
		rw.Init(p.version)
	}
	// Register the peer locally
	err := h.peers.Register(p)
	if err != nil && err != p2p.ErrAddPairPeer {
		p.Log().Error("Ethereum peer registration failed", "err", err)
		return err
	}
	defer h.removePeer(p.id)
	if err != p2p.ErrAddPairPeer {
		// Register the peer in the downloader. If the downloader considers it banned, we disconnect
		if err := h.downloader.RegisterPeer(p.id, uint(p.version), p); err != nil {
			return err
		}
		// Propagate existing transactions. new transactions appearing
		// after this will be sent via broadcasts.
		h.syncTransactions(p)

		// If we're DAO hard-fork aware, validate any remote peer with regard to the hard-fork
		if daoBlock := h.chainconfig.DAOForkBlock; daoBlock != nil {
			// Request the peer's DAO fork header for extra-data validation
			if err := p.RequestHeadersByNumber(daoBlock.Uint64(), 1, 0, false); err != nil {
				return err
			}
			// Start a timer to disconnect if the peer doesn't reply in time
			p.forkDrop = time.AfterFunc(daoChallengeTimeout, func() {
				p.Log().Debug("Timed out DAO fork-check, dropping")
				h.removePeer(p.id)
			})
			// Make sure it's cleaned up if the peer dies off
			defer func() {
				if p.forkDrop != nil {
					p.forkDrop.Stop()
					p.forkDrop = nil
				}
			}()
		}
	}
	// main loop. handle incoming messages.
	for {
		if err := h.handleMsg(p); err != nil {
			p.Log().Debug("Ethereum message handling failed", "err", err)
			return err
		}
	}
}

// handleMsg is invoked whenever an inbound message is received from a remote
// peer. The remote connection is torn down upon returning any error.
func (h *handler) handleMsg(p *peer) error {
	// Read the next message from the remote peer, and ensure it's fully consumed
	msg, err := p.rw.ReadMsg()
	if err != nil {
		return err
	}
	if msg.Size > ProtocolMaxMsgSize {
		return errResp(ErrMsgTooLarge, "%v > %v", msg.Size, ProtocolMaxMsgSize)
	}
	defer msg.Discard()

	// Handle the message depending on its contents
	switch {
	case msg.Code == StatusMsg:
		// Status messages should never arrive after the handshake
		return errResp(ErrExtraStatusMsg, "uncontrolled status message")

		// Block header query, collect the requested headers and reply
	case msg.Code == GetBlockHeadersMsg:
		// Decode the complex header query
		var query getBlockHeadersData
		if err := msg.Decode(&query); err != nil {
			return errResp(ErrDecode, "%v: %v", msg, err)
		}
		hashMode := query.Origin.Hash != (common.Hash{})

		// Gather headers until the fetch or network limits is reached
		var (
			bytes   common.StorageSize
			headers []*types.Header
			unknown bool
		)
		for !unknown && len(headers) < int(query.Amount) && bytes < softResponseLimit && len(headers) < downloader.MaxHeaderFetch {
			// Retrieve the next header satisfying the query
			var origin *types.Header
			if hashMode {
				origin = h.chain.GetHeaderByHash(query.Origin.Hash)
			} else {
				origin = h.chain.GetHeaderByNumber(query.Origin.Number)
			}
			if origin == nil {
				break
			}
			number := origin.Number.Uint64()
			headers = append(headers, origin)
			bytes += estHeaderRlpSize

			// Advance to the next header of the query
			switch {
			case query.Origin.Hash != (common.Hash{}) && query.Reverse:
				// Hash based traversal towards the genesis block
				for i := 0; i < int(query.Skip)+1; i++ {
					if header := h.chain.GetHeader(query.Origin.Hash, number); header != nil {
						query.Origin.Hash = header.ParentHash
						number--
					} else {
						unknown = true
						break
					}
				}
			case query.Origin.Hash != (common.Hash{}) && !query.Reverse:
				// Hash based traversal towards the leaf block
				var (
					current = origin.Number.Uint64()
					next    = current + query.Skip + 1
				)
				if next <= current {
					infos, _ := json.MarshalIndent(p.Peer.Info(), "", "  ")
					p.Log().Warn("GetBlockHeaders skip overflow attack", "current", current, "skip", query.Skip, "next", next, "attacker", infos)
					unknown = true
				} else {
					if header := h.chain.GetHeaderByNumber(next); header != nil {
						if h.chain.GetBlockHashesFromHash(header.Hash(), query.Skip+1)[query.Skip] == query.Origin.Hash {
							query.Origin.Hash = header.Hash()
						} else {
							unknown = true
						}
					} else {
						unknown = true
					}
				}
			case query.Reverse:
				// Number based traversal towards the genesis block
				if query.Origin.Number >= query.Skip+1 {
					query.Origin.Number -= query.Skip + 1
				} else {
					unknown = true
				}

			case !query.Reverse:
				// Number based traversal towards the leaf block
				query.Origin.Number += query.Skip + 1
			}
		}
		return p.SendBlockHeaders(headers)

	case msg.Code == BlockHeadersMsg:
		// A batch of headers arrived to one of our previous requests
		var headers []*types.Header
		if err := msg.Decode(&headers); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		// If no headers were received, but we're expending a DAO fork check, maybe it's that
		if len(headers) == 0 && p.forkDrop != nil {
			// Possibly an empty reply to the fork header checks, sanity check TDs
			verifyDAO := true

			// If we already have a DAO header, we can check the peer's TD against it. If
			// the peer's ahead of this, it too must have a reply to the DAO check
			if daoHeader := h.chain.GetHeaderByNumber(h.chainconfig.DAOForkBlock.Uint64()); daoHeader != nil {
				if _, td := p.Head(); td.Cmp(h.chain.GetTd(daoHeader.Hash(), daoHeader.Number.Uint64())) >= 0 {
					verifyDAO = false
				}
			}
			// If we're seemingly on the same chain, disable the drop timer
			if verifyDAO {
				p.Log().Debug("Seems to be on the same side of the DAO fork")
				p.forkDrop.Stop()
				p.forkDrop = nil
				return nil
			}
		}
		// Filter out any explicitly requested headers, deliver the rest to the downloader
		filter := len(headers) == 1
		if filter {
			// If it's a potential DAO fork check, validate against the rules
			if p.forkDrop != nil && h.chainconfig.DAOForkBlock.Cmp(headers[0].Number) == 0 {
				// Disable the fork drop timer
				p.forkDrop.Stop()
				p.forkDrop = nil

				// Validate the header and either drop the peer or continue
				if err := misc.VerifyDAOHeaderExtraData(h.chainconfig, headers[0]); err != nil {
					p.Log().Debug("Verified to be on the other side of the DAO fork, dropping")
					return err
				}
				p.Log().Debug("Verified to be on the same side of the DAO fork")
				return nil
			}
			// Irrelevant of the fork checks, send the header to the fetcher just in case
			headers = h.fetcher.FilterHeaders(p.id, headers, time.Now())
		}
		if len(headers) > 0 || !filter {
			err := h.downloader.DeliverHeaders(p.id, headers)
			if err != nil {
				log.Debug("Failed to deliver headers", "err", err)
			}
		}

	case msg.Code == GetBlockBodiesMsg:
		// Decode the retrieval message
		msgStream := rlp.NewStream(msg.Payload, uint64(msg.Size))
		if _, err := msgStream.List(); err != nil {
			return err
		}
		// Gather blocks until the fetch or network limits is reached
		var (
			hash   common.Hash
			bytes  int
			bodies []rlp.RawValue
		)
		for bytes < softResponseLimit && len(bodies) < downloader.MaxBlockFetch {
			// Retrieve the hash of the next block
			if err := msgStream.Decode(&hash); err == rlp.EOL {
				break
			} else if err != nil {
				return errResp(ErrDecode, "msg %v: %v", msg, err)
			}
			// Retrieve the requested block body, stopping if enough was found
			if data := h.chain.GetBodyRLP(hash); len(data) != 0 {
				bodies = append(bodies, data)
				bytes += len(data)
			}
		}
		return p.SendBlockBodiesRLP(bodies)

	case msg.Code == BlockBodiesMsg:
		// A batch of block bodies arrived to one of our previous requests
		var request blockBodiesData
		if err := msg.Decode(&request); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		// Deliver them all to the downloader for queuing
		trasactions := make([][]*types.Transaction, len(request))
		uncles := make([][]*types.Header, len(request))

		for i, body := range request {
			trasactions[i] = body.Transactions
			uncles[i] = body.Uncles
		}
		// Filter out any explicitly requested bodies, deliver the rest to the downloader
		filter := len(trasactions) > 0 || len(uncles) > 0
		if filter {
			trasactions, uncles = h.fetcher.FilterBodies(p.id, trasactions, uncles, time.Now())
		}
		if len(trasactions) > 0 || len(uncles) > 0 || !filter {
			err := h.downloader.DeliverBodies(p.id, trasactions, uncles)
			if err != nil {
				log.Debug("Failed to deliver bodies", "err", err)
			}
		}

	case p.version >= eth63 && msg.Code == GetNodeDataMsg:
		// Decode the retrieval message
		msgStream := rlp.NewStream(msg.Payload, uint64(msg.Size))
		if _, err := msgStream.List(); err != nil {
			return err
		}
		// Gather state data until the fetch or network limits is reached
		var (
			hash  common.Hash
			bytes int
			data  [][]byte
		)
		for bytes < softResponseLimit && len(data) < downloader.MaxStateFetch {
			// Retrieve the hash of the next state entry
			if err := msgStream.Decode(&hash); err == rlp.EOL {
				break
			} else if err != nil {
				return errResp(ErrDecode, "msg %v: %v", msg, err)
			}
			// Retrieve the requested state entry, stopping if enough was found
			if entry, err := h.chain.TrieNode(hash); err == nil {
				data = append(data, entry)
				bytes += len(entry)
			}
		}
		return p.SendNodeData(data)

	case p.version >= eth63 && msg.Code == NodeDataMsg:
		// A batch of node state data arrived to one of our previous requests
		var data [][]byte
		if err := msg.Decode(&data); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		// Deliver all to the downloader
		if err := h.downloader.DeliverNodeData(p.id, data); err != nil {
			log.Debug("Failed to deliver node state data", "err", err)
		}

	case p.version >= eth63 && msg.Code == GetReceiptsMsg:
		// Decode the retrieval message
		msgStream := rlp.NewStream(msg.Payload, uint64(msg.Size))
		if _, err := msgStream.List(); err != nil {
			return err
		}
		// Gather state data until the fetch or network limits is reached
		var (
			hash     common.Hash
			bytes    int
			receipts []rlp.RawValue
		)
		for bytes < softResponseLimit && len(receipts) < downloader.MaxReceiptFetch {
			// Retrieve the hash of the next block
			if err := msgStream.Decode(&hash); err == rlp.EOL {
				break
			} else if err != nil {
				return errResp(ErrDecode, "msg %v: %v", msg, err)
			}
			// Retrieve the requested block's receipts, skipping if unknown to us
			results := h.chain.GetReceiptsByHash(hash)
			if results == nil {
				if header := h.chain.GetHeaderByHash(hash); header == nil || header.ReceiptHash != types.EmptyRootHash {
					continue
				}
			}
			// If known, encode and queue for response packet
			if encoded, err := rlp.EncodeToBytes(results); err != nil {
				log.Error("Failed to encode receipt", "err", err)
			} else {
				receipts = append(receipts, encoded)
				bytes += len(encoded)
			}
		}
		return p.SendReceiptsRLP(receipts)

	case p.version >= eth63 && msg.Code == ReceiptsMsg:
		// A batch of receipts arrived to one of our previous requests
		var receipts [][]*types.Receipt
		if err := msg.Decode(&receipts); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		// Deliver all to the downloader
		if err := h.downloader.DeliverReceipts(p.id, receipts); err != nil {
			log.Debug("Failed to deliver receipts", "err", err)
		}

	case msg.Code == NewBlockHashesMsg:
		var announces newBlockHashesData
		if err := msg.Decode(&announces); err != nil {
			return errResp(ErrDecode, "%v: %v", msg, err)
		}
		// Mark the hashes as present at the remote node
		for _, block := range announces {
			p.MarkBlock(block.Hash)
		}
		// Schedule all the unknown hashes for retrieval
		unknown := make(newBlockHashesData, 0, len(announces))
		for _, block := range announces {
			if !h.chain.HasBlock(block.Hash, block.Number) {
				unknown = append(unknown, block)
			}
		}
		for _, block := range unknown {
			h.fetcher.Notify(p.id, block.Hash, block.Number, time.Now(), p.RequestOneHeader, p.RequestBodies)
		}

	case msg.Code == NewBlockMsg:
		// Retrieve and decode the propagated block
		var request newBlockData
		if err := msg.Decode(&request); err != nil {
			return errResp(ErrDecode, "%v: %v", msg, err)
		}
		request.Block.ReceivedAt = msg.ReceivedAt
		request.Block.ReceivedFrom = p

		// Mark the peer as owning the block and schedule it for import
		p.MarkBlock(request.Block.Hash())
		h.fetcher.Enqueue(p.id, request.Block)

		// Assuming the block is importable by the peer, but possibly not yet done so,
		// calculate the head hash and TD that the peer truly must have.
		var (
			trueHead = request.Block.ParentHash()
			trueTD   = new(big.Int).Sub(request.TD, request.Block.Difficulty())
		)
		// Update the peers total difficulty if better than the previous
		if _, td := p.Head(); trueTD.Cmp(td) > 0 {
			p.SetHead(trueHead, trueTD)

			// Schedule a sync if above ours. Note, this will not fire a sync for a gap of
			// a singe block (as the true TD is below the propagated block), however this
			// scenario should easily be covered by the fetcher.
			currentBlock := h.chain.CurrentBlock()
			if trueTD.Cmp(h.chain.GetTd(currentBlock.Hash(), currentBlock.NumberU64())) > 0 {
				go h.synchronise(p)
			}
		}

	case msg.Code == TxMsg:
		// Transactions arrived, make sure we have a valid and fresh chain to handle them
		if atomic.LoadUint32(&h.acceptTxs) == 0 {
			break
		}
		// Transactions can be processed, parse all of them and deliver to the pool
		var txs []*types.Transaction
		if err := msg.Decode(&txs); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		var unkownTxs []*types.Transaction

		for i, tx := range txs {
			// Validate and mark the remote transaction
			if tx == nil {
				return errResp(ErrDecode, "transaction %d is nil", i)
			}
			p.MarkTransaction(tx.Hash())
			exist, _ := h.knownTxs.ContainsOrAdd(tx.Hash(), true)
			if !exist {
				unkownTxs = append(unkownTxs, tx)
			} else {
				log.Trace("Discard known tx", "hash", tx.Hash(), "nonce", tx.Nonce(), "to", tx.To())
			}

		}
		h.txpool.AddRemotes(txs)

	case msg.Code == OrderTxMsg:
		// Transactions arrived, make sure we have a valid and fresh chain to handle them
		if atomic.LoadUint32(&h.acceptTxs) == 0 {
			break
		}
		// Transactions can be processed, parse all of them and deliver to the pool
		var txs []*types.OrderTransaction
		if err := msg.Decode(&txs); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		var unkownOrderTxs []*types.OrderTransaction

		for i, tx := range txs {
			// Validate and mark the remote transaction
			if tx == nil {
				return errResp(ErrDecode, "transaction %d is nil", i)
			}
			p.MarkOrderTransaction(tx.Hash())
			exist, _ := h.knowOrderTxs.ContainsOrAdd(tx.Hash(), true)
			if !exist {
				unkownOrderTxs = append(unkownOrderTxs, tx)
			} else {
				log.Trace("Discard known tx", "hash", tx.Hash(), "nonce", tx.Nonce())
			}

		}

		if h.orderpool != nil {
			h.orderpool.AddRemotes(txs)
		}

	case msg.Code == LendingTxMsg:
		// Transactions arrived, make sure we have a valid and fresh chain to handle them
		if atomic.LoadUint32(&h.acceptTxs) == 0 {
			break
		}
		// Transactions can be processed, parse all of them and deliver to the pool
		var txs []*types.LendingTransaction
		if err := msg.Decode(&txs); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		var unkownLendingTxs []*types.LendingTransaction

		for i, tx := range txs {
			// Validate and mark the remote transaction
			if tx == nil {
				return errResp(ErrDecode, "transaction %d is nil", i)
			}
			p.MarkLendingTransaction(tx.Hash())
			exist, _ := h.knowLendingTxs.ContainsOrAdd(tx.Hash(), true)
			if !exist {
				unkownLendingTxs = append(unkownLendingTxs, tx)
			} else {
				log.Trace("Discard known tx", "hash", tx.Hash(), "nonce", tx.Nonce())
			}

		}

		if h.lendingpool != nil {
			h.lendingpool.AddRemotes(txs)
		}

	default:
		return errResp(ErrInvalidMsgCode, "%v", msg.Code)
	}
	return nil
}
*/

// BroadcastBlock will either propagate a block to a subset of its peers, or
// will only announce its availability (depending on what's requested).
func (h *handler) BroadcastBlock(block *types.Block, propagate bool) {
	hash := block.Hash()
	peers := h.peers.peersWithoutBlock(hash)

	// If propagation is requested, send to a subset of the peer
	if propagate {
		// Calculate the TD of the block (it's not imported yet, so block.Td is not valid)
		var td *big.Int
		if parent := h.chain.GetBlock(block.ParentHash(), block.NumberU64()-1); parent != nil {
			td = new(big.Int).Add(block.Difficulty(), h.chain.GetTd(block.ParentHash(), block.NumberU64()-1))
		} else {
			log.Error("Propagating dangling block", "number", block.Number(), "hash", hash)
			return
		}
		// Send the block to a subset of our peers
		for _, peer := range peers {
			peer.SendNewBlock(block, td)
		}
		log.Trace("Propagated block", "hash", hash, "recipients", len(peers), "duration", common.PrettyDuration(time.Since(block.ReceivedAt)))
		return
	}
	// Otherwise if the block is indeed in out own chain, announce it
	if h.chain.HasBlock(hash, block.NumberU64()) {
		for _, peer := range peers {
			peer.SendNewBlockHashes([]common.Hash{hash}, []uint64{block.NumberU64()})
		}
		log.Trace("Announced block", "hash", hash, "recipients", len(peers), "duration", common.PrettyDuration(time.Since(block.ReceivedAt)))
	}
}

// BroadcastTx will propagate a transaction to all peers which are not known to
// already have the given transaction.
func (h *handler) BroadcastTx(hash common.Hash, tx *types.Transaction) {
	// Broadcast transaction to a batch of peers not knowing about it
	peers := h.peers.peersWithoutTransaction(hash)
	//FIXME include this again: peers = peers[:int(math.Sqrt(float64(len(peers))))]
	for _, peer := range peers {
		peer.SendTransactions(types.Transactions{tx})
	}
	log.Trace("Broadcast transaction", "hash", hash, "recipients", len(peers))
}

// OrderBroadcastTx will propagate a transaction to all peers which are not known to
// already have the given transaction.
func (h *handler) OrderBroadcastTx(hash common.Hash, tx *types.OrderTransaction) {
	// Broadcast transaction to a batch of peers not knowing about it
	peers := h.peers.OrderPeersWithoutTx(hash)
	//FIXME include this again: peers = peers[:int(math.Sqrt(float64(len(peers))))]
	for _, peer := range peers {
		peer.SendOrderTransactions(types.OrderTransactions{tx})
	}
	log.Trace("Broadcast order transaction", "hash", hash, "recipients", len(peers))
}

// LendingBroadcastTx will propagate a transaction to all peers which are not known to
// already have the given transaction.
func (h *handler) LendingBroadcastTx(hash common.Hash, tx *types.LendingTransaction) {
	// Broadcast transaction to a batch of peers not knowing about it
	peers := h.peers.LendingPeersWithoutTx(hash)
	//FIXME include this again: peers = peers[:int(math.Sqrt(float64(len(peers))))]
	for _, peer := range peers {
		peer.SendLendingTransactions(types.LendingTransactions{tx})
	}
	log.Trace("Broadcast lending transaction", "hash", hash, "recipients", len(peers))
}

// minedBroadcastLoop broadcast loop
func (h *handler) minedBroadcastLoop() {
	// automatically stops if unsubscribe
	for obj := range h.minedBlockSub.Chan() {
		switch ev := obj.Data.(type) {
		case core.NewMinedBlockEvent:
			h.BroadcastBlock(ev.Block, true) // First propagate block to peers
			//self.BroadcastBlock(ev.Block, false) // Only then announce to the rest
		}
	}
}

func (h *handler) txBroadcastLoop() {
	for {
		select {
		case event := <-h.txCh:
			h.BroadcastTx(event.Tx.Hash(), event.Tx)

			// Err() channel will be closed when unsubscribing.
		case <-h.txSub.Err():
			return
		}
	}
}

// orderTxBroadcastLoop broadcast order
func (h *handler) orderTxBroadcastLoop() {
	if h.orderTxSub == nil {
		return
	}
	for {
		select {
		case event := <-h.orderTxCh:
			h.OrderBroadcastTx(event.Tx.Hash(), event.Tx)

			// Err() channel will be closed when unsubscribing.
		case <-h.orderTxSub.Err():
			return
		}
	}
}

// lendingTxBroadcastLoop broadcast order
func (h *handler) lendingTxBroadcastLoop() {
	if h.lendingTxSub == nil {
		return
	}
	for {
		select {
		case event := <-h.lendingTxCh:
			h.LendingBroadcastTx(event.Tx.Hash(), event.Tx)

			// Err() channel will be closed when unsubscribing.
		case <-h.lendingTxSub.Err():
			return
		}
	}
}

// NodeInfo represents a short summary of the Ethereum sub-protocol metadata
// known about the host peer.
type NodeInfo struct {
	Network    uint64              `json:"network"`    // Ethereum network ID (1=Frontier, 2=Morden, Ropsten=3, Rinkeby=4)
	Difficulty *big.Int            `json:"difficulty"` // Total difficulty of the host's blockchain
	Genesis    common.Hash         `json:"genesis"`    // SHA3 hash of the host's genesis block
	Config     *params.ChainConfig `json:"config"`     // Chain configuration for the fork rules
	Head       common.Hash         `json:"head"`       // SHA3 hash of the host's best owned block
}

// NodeInfo retrieves some protocol metadata about the running host node.
func (h *handler) NodeInfo() *NodeInfo {
	currentBlock := h.chain.CurrentBlock()
	return &NodeInfo{
		Network:    h.networkId,
		Difficulty: h.chain.GetTd(currentBlock.Hash(), currentBlock.NumberU64()),
		Genesis:    h.chain.Genesis().Hash(),
		Config:     h.chain.Config(),
		Head:       currentBlock.Hash(),
	}
}
