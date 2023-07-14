package ethapi

import (
	"context"
	"errors"
	"math/big"

	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/common/hexutil"
	"github.com/tomochain/tomochain/core"
	"github.com/tomochain/tomochain/core/types"
	"github.com/tomochain/tomochain/rlp"
	"github.com/tomochain/tomochain/tomox/tradingstate"
	"github.com/tomochain/tomochain/tomoxlending/lendingstate"
)

// SendOrderRawTransaction will add the signed transaction to the transaction pool.
// The sender is responsible for signing the transaction and using the correct nonce.
func (s *PublicTomoXTransactionPoolAPI) SendOrderRawTransaction(ctx context.Context, encodedTx hexutil.Bytes) (common.Hash, error) {
	tx := new(types.OrderTransaction)
	if err := rlp.DecodeBytes(encodedTx, tx); err != nil {
		return common.Hash{}, err
	}
	return submitOrderTransaction(ctx, s.b, tx)
}

// SendLendingRawTransaction will add the signed transaction to the transaction pool.
// The sender is responsible for signing the transaction and using the correct nonce.
func (s *PublicTomoXTransactionPoolAPI) SendLendingRawTransaction(ctx context.Context, encodedTx hexutil.Bytes) (common.Hash, error) {
	tx := new(types.LendingTransaction)
	if err := rlp.DecodeBytes(encodedTx, tx); err != nil {
		return common.Hash{}, err
	}
	return submitLendingTransaction(ctx, s.b, tx)
}

// GetOrderTxMatchByHash returns the bytes of the transaction for the given hash.
func (s *PublicTomoXTransactionPoolAPI) GetOrderTxMatchByHash(ctx context.Context, hash common.Hash) ([]*tradingstate.OrderItem, error) {
	var tx *types.Transaction
	orders := []*tradingstate.OrderItem{}
	if tx, _, _, _ = core.GetTransaction(s.b.ChainDb(), hash); tx == nil {
		if tx = s.b.GetPoolTransaction(hash); tx == nil {
			return []*tradingstate.OrderItem{}, nil
		}
	}

	batch, err := tradingstate.DecodeTxMatchesBatch(tx.Data())
	if err != nil {
		return []*tradingstate.OrderItem{}, err
	}
	for _, txMatch := range batch.Data {
		order, err := txMatch.DecodeOrder()
		if err != nil {
			return []*tradingstate.OrderItem{}, err
		}
		orders = append(orders, order)
	}
	return orders, nil

}

// GetOrderPoolContent return pending, queued content
func (s *PublicTomoXTransactionPoolAPI) GetOrderPoolContent(ctx context.Context) interface{} {
	pendingOrders := []*tradingstate.OrderItem{}
	queuedOrders := []*tradingstate.OrderItem{}
	pending, queued := s.b.OrderTxPoolContent()

	for _, txs := range pending {
		for _, tx := range txs {
			V, R, S := tx.Signature()
			order := &tradingstate.OrderItem{
				Nonce:           big.NewInt(int64(tx.Nonce())),
				Quantity:        tx.Quantity(),
				Price:           tx.Price(),
				ExchangeAddress: tx.ExchangeAddress(),
				UserAddress:     tx.UserAddress(),
				BaseToken:       tx.BaseToken(),
				QuoteToken:      tx.QuoteToken(),
				Status:          tx.Status(),
				Side:            tx.Side(),
				Type:            tx.Type(),
				Hash:            tx.OrderHash(),
				OrderID:         tx.OrderID(),
				Signature: &tradingstate.Signature{
					V: byte(V.Uint64()),
					R: common.BigToHash(R),
					S: common.BigToHash(S),
				},
			}
			pendingOrders = append(pendingOrders, order)
		}
	}

	for _, txs := range queued {
		for _, tx := range txs {
			V, R, S := tx.Signature()
			order := &tradingstate.OrderItem{
				Nonce:           big.NewInt(int64(tx.Nonce())),
				Quantity:        tx.Quantity(),
				Price:           tx.Price(),
				ExchangeAddress: tx.ExchangeAddress(),
				UserAddress:     tx.UserAddress(),
				BaseToken:       tx.BaseToken(),
				QuoteToken:      tx.QuoteToken(),
				Status:          tx.Status(),
				Side:            tx.Side(),
				Type:            tx.Type(),
				Hash:            tx.OrderHash(),
				OrderID:         tx.OrderID(),
				Signature: &tradingstate.Signature{
					V: byte(V.Uint64()),
					R: common.BigToHash(R),
					S: common.BigToHash(S),
				},
			}
			queuedOrders = append(pendingOrders, order)
		}
	}

	return map[string]interface{}{
		"pending": pendingOrders,
		"queued":  queuedOrders,
	}
}

// GetOrderStats return pending, queued length
func (s *PublicTomoXTransactionPoolAPI) GetOrderStats(ctx context.Context) interface{} {
	pending, queued := s.b.OrderStats()
	return map[string]interface{}{
		"pending": pending,
		"queued":  queued,
	}
}

// OrderMsg struct
type OrderMsg struct {
	AccountNonce    hexutil.Uint64 `json:"nonce"    gencodec:"required"`
	Quantity        hexutil.Big    `json:"quantity,omitempty"`
	Price           hexutil.Big    `json:"price,omitempty"`
	ExchangeAddress common.Address `json:"exchangeAddress,omitempty"`
	UserAddress     common.Address `json:"userAddress,omitempty"`
	BaseToken       common.Address `json:"baseToken,omitempty"`
	QuoteToken      common.Address `json:"quoteToken,omitempty"`
	Status          string         `json:"status,omitempty"`
	Side            string         `json:"side,omitempty"`
	Type            string         `json:"type,omitempty"`
	OrderID         hexutil.Uint64 `json:"orderid,omitempty"`
	// Signature values
	V hexutil.Big `json:"v" gencodec:"required"`
	R hexutil.Big `json:"r" gencodec:"required"`
	S hexutil.Big `json:"s" gencodec:"required"`

	// This is only used when marshaling to JSON.
	Hash common.Hash `json:"hash" rlp:"-"`
}

// LendingMsg api message for lending
type LendingMsg struct {
	AccountNonce    hexutil.Uint64 `json:"nonce"    gencodec:"required"`
	Quantity        hexutil.Big    `json:"quantity,omitempty"`
	RelayerAddress  common.Address `json:"relayerAddress,omitempty"`
	UserAddress     common.Address `json:"userAddress,omitempty"`
	CollateralToken common.Address `json:"collateralToken,omitempty"`
	AutoTopUp       bool           `json:"autoTopUp,omitempty"`
	LendingToken    common.Address `json:"lendingToken,omitempty"`
	Term            hexutil.Uint64 `json:"term,omitempty"`
	Interest        hexutil.Uint64 `json:"interest,omitempty"`
	Status          string         `json:"status,omitempty"`
	Side            string         `json:"side,omitempty"`
	Type            string         `json:"type,omitempty"`
	LendingId       hexutil.Uint64 `json:"lendingId,omitempty"`
	LendingTradeId  hexutil.Uint64 `json:"tradeId,omitempty"`
	ExtraData       string         `json:"extraData,omitempty"`

	// Signature values
	V hexutil.Big `json:"v" gencodec:"required"`
	R hexutil.Big `json:"r" gencodec:"required"`
	S hexutil.Big `json:"s" gencodec:"required"`

	// This is only used when marshaling to JSON.
	Hash common.Hash `json:"hash" rlp:"-"`
}

type PriceVolume struct {
	Price  *big.Int `json:"price,omitempty"`
	Volume *big.Int `json:"volume,omitempty"`
}

type InterestVolume struct {
	Interest *big.Int `json:"interest,omitempty"`
	Volume   *big.Int `json:"volume,omitempty"`
}

// SendOrder will add the signed transaction to the transaction pool.
// The sender is responsible for signing the transaction and using the correct nonce.
func (s *PublicTomoXTransactionPoolAPI) SendOrder(ctx context.Context, msg OrderMsg) (common.Hash, error) {
	tx := types.NewOrderTransaction(uint64(msg.AccountNonce), msg.Quantity.ToInt(), msg.Price.ToInt(), msg.ExchangeAddress, msg.UserAddress, msg.BaseToken, msg.QuoteToken, msg.Status, msg.Side, msg.Type, msg.Hash, uint64(msg.OrderID))
	tx = tx.ImportSignature(msg.V.ToInt(), msg.R.ToInt(), msg.S.ToInt())
	return submitOrderTransaction(ctx, s.b, tx)
}

// SendLending will add the signed transaction to the transaction pool.
// The sender is responsible for signing the transaction and using the correct nonce.
func (s *PublicTomoXTransactionPoolAPI) SendLending(ctx context.Context, msg LendingMsg) (common.Hash, error) {
	tx := types.NewLendingTransaction(uint64(msg.AccountNonce), msg.Quantity.ToInt(), uint64(msg.Interest), uint64(msg.Term), msg.RelayerAddress, msg.UserAddress, msg.LendingToken, msg.CollateralToken, msg.AutoTopUp, msg.Status, msg.Side, msg.Type, msg.Hash, uint64(msg.LendingId), uint64(msg.LendingTradeId), msg.ExtraData)
	tx = tx.ImportSignature(msg.V.ToInt(), msg.R.ToInt(), msg.S.ToInt())
	return submitLendingTransaction(ctx, s.b, tx)
}

// GetOrderCount returns the number of transactions the given address has sent for the given block number
func (s *PublicTomoXTransactionPoolAPI) GetOrderCount(ctx context.Context, addr common.Address) (*hexutil.Uint64, error) {

	nonce, err := s.b.GetOrderNonce(addr.Hash())
	if err != nil {
		return (*hexutil.Uint64)(&nonce), err
	}
	return (*hexutil.Uint64)(&nonce), err
}

func (s *PublicTomoXTransactionPoolAPI) GetBestBid(ctx context.Context, baseToken, quoteToken common.Address) (PriceVolume, error) {

	result := PriceVolume{}
	block := s.b.CurrentBlock()
	if block == nil {
		return result, errors.New("Current block not found")
	}
	tomoxService := s.b.TomoxService()
	if tomoxService == nil {
		return result, errors.New("TomoX service not found")
	}
	author, err := s.b.GetEngine().Author(block.Header())
	if err != nil {
		return result, err
	}
	tomoxState, err := tomoxService.GetTradingState(block, author)
	if err != nil {
		return result, err
	}
	result.Price, result.Volume = tomoxState.GetBestBidPrice(tradingstate.GetTradingOrderBookHash(baseToken, quoteToken))
	if result.Price.Sign() == 0 {
		return result, errors.New("Bid tree not found")
	}
	return result, nil
}

func (s *PublicTomoXTransactionPoolAPI) GetBestAsk(ctx context.Context, baseToken, quoteToken common.Address) (PriceVolume, error) {
	result := PriceVolume{}
	block := s.b.CurrentBlock()
	if block == nil {
		return result, errors.New("Current block not found")
	}
	tomoxService := s.b.TomoxService()
	if tomoxService == nil {
		return result, errors.New("TomoX service not found")
	}
	author, err := s.b.GetEngine().Author(block.Header())
	if err != nil {
		return result, err
	}
	tomoxState, err := tomoxService.GetTradingState(block, author)
	if err != nil {
		return result, err
	}
	result.Price, result.Volume = tomoxState.GetBestAskPrice(tradingstate.GetTradingOrderBookHash(baseToken, quoteToken))
	if result.Price.Sign() == 0 {
		return result, errors.New("Ask tree not found")
	}
	return result, nil
}

func (s *PublicTomoXTransactionPoolAPI) GetBidTree(ctx context.Context, baseToken, quoteToken common.Address) (map[*big.Int]tradingstate.DumpOrderList, error) {
	block := s.b.CurrentBlock()
	if block == nil {
		return nil, errors.New("Current block not found")
	}
	tomoxService := s.b.TomoxService()
	if tomoxService == nil {
		return nil, errors.New("TomoX service not found")
	}
	author, err := s.b.GetEngine().Author(block.Header())
	if err != nil {
		return nil, err
	}
	tomoxState, err := tomoxService.GetTradingState(block, author)
	if err != nil {
		return nil, err
	}
	result, err := tomoxState.DumpBidTrie(tradingstate.GetTradingOrderBookHash(baseToken, quoteToken))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *PublicTomoXTransactionPoolAPI) GetPrice(ctx context.Context, baseToken, quoteToken common.Address) (*big.Int, error) {
	block := s.b.CurrentBlock()
	if block == nil {
		return nil, errors.New("Current block not found")
	}
	tomoxService := s.b.TomoxService()
	if tomoxService == nil {
		return nil, errors.New("TomoX service not found")
	}
	author, err := s.b.GetEngine().Author(block.Header())
	if err != nil {
		return nil, err
	}
	tomoxState, err := tomoxService.GetTradingState(block, author)
	if err != nil {
		return nil, err
	}
	price := tomoxState.GetLastPrice(tradingstate.GetTradingOrderBookHash(baseToken, quoteToken))
	if price == nil || price.Sign() == 0 {
		return common.Big0, errors.New("Order book's price not found")
	}
	return price, nil
}

func (s *PublicTomoXTransactionPoolAPI) GetLastEpochPrice(ctx context.Context, baseToken, quoteToken common.Address) (*big.Int, error) {
	block := s.b.CurrentBlock()
	if block == nil {
		return nil, errors.New("Current block not found")
	}
	tomoxService := s.b.TomoxService()
	if tomoxService == nil {
		return nil, errors.New("TomoX service not found")
	}
	author, err := s.b.GetEngine().Author(block.Header())
	if err != nil {
		return nil, err
	}
	tomoxState, err := tomoxService.GetTradingState(block, author)
	if err != nil {
		return nil, err
	}
	price := tomoxState.GetMediumPriceBeforeEpoch(tradingstate.GetTradingOrderBookHash(baseToken, quoteToken))
	if price == nil || price.Sign() == 0 {
		return common.Big0, errors.New("Order book's price not found")
	}
	return price, nil
}

func (s *PublicTomoXTransactionPoolAPI) GetCurrentEpochPrice(ctx context.Context, baseToken, quoteToken common.Address) (*big.Int, error) {
	block := s.b.CurrentBlock()
	if block == nil {
		return nil, errors.New("Current block not found")
	}
	tomoxService := s.b.TomoxService()
	if tomoxService == nil {
		return nil, errors.New("TomoX service not found")
	}
	author, err := s.b.GetEngine().Author(block.Header())
	if err != nil {
		return nil, err
	}
	tomoxState, err := tomoxService.GetTradingState(block, author)
	if err != nil {
		return nil, err
	}
	price, _ := tomoxState.GetMediumPriceAndTotalAmount(tradingstate.GetTradingOrderBookHash(baseToken, quoteToken))
	if price == nil || price.Sign() == 0 {
		return common.Big0, errors.New("Order book's price not found")
	}
	return price, nil
}

func (s *PublicTomoXTransactionPoolAPI) GetAskTree(ctx context.Context, baseToken, quoteToken common.Address) (map[*big.Int]tradingstate.DumpOrderList, error) {
	block := s.b.CurrentBlock()
	if block == nil {
		return nil, errors.New("Current block not found")
	}
	tomoxService := s.b.TomoxService()
	if tomoxService == nil {
		return nil, errors.New("TomoX service not found")
	}
	author, err := s.b.GetEngine().Author(block.Header())
	if err != nil {
		return nil, err
	}
	tomoxState, err := tomoxService.GetTradingState(block, author)
	if err != nil {
		return nil, err
	}
	result, err := tomoxState.DumpAskTrie(tradingstate.GetTradingOrderBookHash(baseToken, quoteToken))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *PublicTomoXTransactionPoolAPI) GetOrderById(ctx context.Context, baseToken, quoteToken common.Address, orderId uint64) (interface{}, error) {
	block := s.b.CurrentBlock()
	if block == nil {
		return nil, errors.New("Current block not found")
	}
	tomoxService := s.b.TomoxService()
	if tomoxService == nil {
		return nil, errors.New("TomoX service not found")
	}
	author, err := s.b.GetEngine().Author(block.Header())
	if err != nil {
		return nil, err
	}
	tomoxState, err := tomoxService.GetTradingState(block, author)
	if err != nil {
		return nil, err
	}
	orderIdHash := common.BigToHash(new(big.Int).SetUint64(orderId))
	orderitem := tomoxState.GetOrder(tradingstate.GetTradingOrderBookHash(baseToken, quoteToken), orderIdHash)
	if orderitem.Quantity == nil || orderitem.Quantity.Sign() == 0 {
		return nil, errors.New("Order not found")
	}
	return orderitem, nil
}

func (s *PublicTomoXTransactionPoolAPI) GetTradingOrderBookInfo(ctx context.Context, baseToken, quoteToken common.Address) (*tradingstate.DumpOrderBookInfo, error) {
	block := s.b.CurrentBlock()
	if block == nil {
		return nil, errors.New("Current block not found")
	}
	tomoxService := s.b.TomoxService()
	if tomoxService == nil {
		return nil, errors.New("TomoX service not found")
	}
	author, err := s.b.GetEngine().Author(block.Header())
	if err != nil {
		return nil, err
	}
	tomoxState, err := tomoxService.GetTradingState(block, author)
	if err != nil {
		return nil, err
	}
	result, err := tomoxState.DumpOrderBookInfo(tradingstate.GetTradingOrderBookHash(baseToken, quoteToken))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *PublicTomoXTransactionPoolAPI) GetLiquidationPriceTree(ctx context.Context, baseToken, quoteToken common.Address) (map[*big.Int]tradingstate.DumpLendingBook, error) {
	block := s.b.CurrentBlock()
	if block == nil {
		return nil, errors.New("Current block not found")
	}
	tomoxService := s.b.TomoxService()
	if tomoxService == nil {
		return nil, errors.New("TomoX service not found")
	}
	author, err := s.b.GetEngine().Author(block.Header())
	if err != nil {
		return nil, err
	}
	tomoxState, err := tomoxService.GetTradingState(block, author)
	if err != nil {
		return nil, err
	}
	result, err := tomoxState.DumpLiquidationPriceTrie(tradingstate.GetTradingOrderBookHash(baseToken, quoteToken))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *PublicTomoXTransactionPoolAPI) GetInvestingTree(ctx context.Context, lendingToken common.Address, term uint64) (map[*big.Int]lendingstate.DumpOrderList, error) {
	block := s.b.CurrentBlock()
	if block == nil {
		return nil, errors.New("Current block not found")
	}
	lendingService := s.b.LendingService()
	if lendingService == nil {
		return nil, errors.New("TomoX Lending service not found")
	}
	author, err := s.b.GetEngine().Author(block.Header())
	if err != nil {
		return nil, err
	}
	lendingState, err := lendingService.GetLendingState(block, author)
	if err != nil {
		return nil, err
	}
	result, err := lendingState.DumpInvestingTrie(lendingstate.GetLendingOrderBookHash(lendingToken, term))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *PublicTomoXTransactionPoolAPI) GetBorrowingTree(ctx context.Context, lendingToken common.Address, term uint64) (map[*big.Int]lendingstate.DumpOrderList, error) {
	block := s.b.CurrentBlock()
	if block == nil {
		return nil, errors.New("Current block not found")
	}
	lendingService := s.b.LendingService()
	if lendingService == nil {
		return nil, errors.New("TomoX Lending service not found")
	}
	author, err := s.b.GetEngine().Author(block.Header())
	if err != nil {
		return nil, err
	}
	lendingState, err := lendingService.GetLendingState(block, author)
	if err != nil {
		return nil, err
	}
	result, err := lendingState.DumpBorrowingTrie(lendingstate.GetLendingOrderBookHash(lendingToken, term))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *PublicTomoXTransactionPoolAPI) GetLendingOrderBookInfo(tx context.Context, lendingToken common.Address, term uint64) (*lendingstate.DumpOrderBookInfo, error) {
	block := s.b.CurrentBlock()
	if block == nil {
		return nil, errors.New("Current block not found")
	}
	lendingService := s.b.LendingService()
	if lendingService == nil {
		return nil, errors.New("TomoX Lending service not found")
	}
	author, err := s.b.GetEngine().Author(block.Header())
	if err != nil {
		return nil, err
	}
	lendingState, err := lendingService.GetLendingState(block, author)
	if err != nil {
		return nil, err
	}
	result, err := lendingState.DumpOrderBookInfo(lendingstate.GetLendingOrderBookHash(lendingToken, term))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *PublicTomoXTransactionPoolAPI) getLendingOrderTree(ctx context.Context, lendingToken common.Address, term uint64) (map[*big.Int]lendingstate.LendingItem, error) {
	block := s.b.CurrentBlock()
	if block == nil {
		return nil, errors.New("Current block not found")
	}
	lendingService := s.b.LendingService()
	if lendingService == nil {
		return nil, errors.New("TomoX Lending service not found")
	}
	author, err := s.b.GetEngine().Author(block.Header())
	if err != nil {
		return nil, err
	}
	lendingState, err := lendingService.GetLendingState(block, author)
	if err != nil {
		return nil, err
	}
	result, err := lendingState.DumpLendingOrderTrie(lendingstate.GetLendingOrderBookHash(lendingToken, term))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *PublicTomoXTransactionPoolAPI) GetLendingTradeTree(ctx context.Context, lendingToken common.Address, term uint64) (map[*big.Int]lendingstate.LendingTrade, error) {
	block := s.b.CurrentBlock()
	if block == nil {
		return nil, errors.New("Current block not found")
	}
	lendingService := s.b.LendingService()
	if lendingService == nil {
		return nil, errors.New("TomoX Lending service not found")
	}
	author, err := s.b.GetEngine().Author(block.Header())
	if err != nil {
		return nil, err
	}
	lendingState, err := lendingService.GetLendingState(block, author)
	if err != nil {
		return nil, err
	}
	result, err := lendingState.DumpLendingTradeTrie(lendingstate.GetLendingOrderBookHash(lendingToken, term))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *PublicTomoXTransactionPoolAPI) GetLiquidationTimeTree(ctx context.Context, lendingToken common.Address, term uint64) (map[*big.Int]lendingstate.DumpOrderList, error) {
	block := s.b.CurrentBlock()
	if block == nil {
		return nil, errors.New("Current block not found")
	}
	lendingService := s.b.LendingService()
	if lendingService == nil {
		return nil, errors.New("TomoX Lending service not found")
	}
	author, err := s.b.GetEngine().Author(block.Header())
	if err != nil {
		return nil, err
	}
	lendingState, err := lendingService.GetLendingState(block, author)
	if err != nil {
		return nil, err
	}
	result, err := lendingState.DumpLiquidationTimeTrie(lendingstate.GetLendingOrderBookHash(lendingToken, term))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *PublicTomoXTransactionPoolAPI) GetLendingOrderCount(ctx context.Context, addr common.Address) (*hexutil.Uint64, error) {
	block := s.b.CurrentBlock()
	if block == nil {
		return nil, errors.New("Current block not found")
	}
	lendingService := s.b.LendingService()
	if lendingService == nil {
		return nil, errors.New("TomoX Lending service not found")
	}
	author, err := s.b.GetEngine().Author(block.Header())
	if err != nil {
		return nil, err
	}
	lendingState, err := lendingService.GetLendingState(block, author)
	if err != nil {
		return nil, err
	}
	nonce := lendingState.GetNonce(addr.Hash())
	return (*hexutil.Uint64)(&nonce), err
}

func (s *PublicTomoXTransactionPoolAPI) GetBestInvesting(ctx context.Context, lendingToken common.Address, term uint64) (InterestVolume, error) {
	result := InterestVolume{}
	block := s.b.CurrentBlock()
	if block == nil {
		return result, errors.New("Current block not found")
	}
	lendingService := s.b.LendingService()
	if lendingService == nil {
		return result, errors.New("TomoX Lending service not found")
	}
	author, err := s.b.GetEngine().Author(block.Header())
	if err != nil {
		return result, err
	}
	lendingState, err := lendingService.GetLendingState(block, author)
	if err != nil {
		return result, err
	}
	result.Interest, result.Volume = lendingState.GetBestInvestingRate(lendingstate.GetLendingOrderBookHash(lendingToken, term))
	return result, nil
}

func (s *PublicTomoXTransactionPoolAPI) GetBestBorrowing(ctx context.Context, lendingToken common.Address, term uint64) (InterestVolume, error) {
	result := InterestVolume{}
	block := s.b.CurrentBlock()
	if block == nil {
		return result, errors.New("Current block not found")
	}
	lendingService := s.b.LendingService()
	if lendingService == nil {
		return result, errors.New("TomoX Lending service not found")
	}
	author, err := s.b.GetEngine().Author(block.Header())
	if err != nil {
		return result, err
	}
	lendingState, err := lendingService.GetLendingState(block, author)
	if err != nil {
		return result, err
	}
	result.Interest, result.Volume = lendingState.GetBestBorrowRate(lendingstate.GetLendingOrderBookHash(lendingToken, term))
	return result, nil
}

func (s *PublicTomoXTransactionPoolAPI) GetBids(ctx context.Context, baseToken, quoteToken common.Address) (map[*big.Int]*big.Int, error) {
	block := s.b.CurrentBlock()
	if block == nil {
		return nil, errors.New("Current block not found")
	}
	tomoxService := s.b.TomoxService()
	if tomoxService == nil {
		return nil, errors.New("TomoX service not found")
	}
	author, err := s.b.GetEngine().Author(block.Header())
	if err != nil {
		return nil, err
	}
	tomoxState, err := tomoxService.GetTradingState(block, author)
	if err != nil {
		return nil, err
	}
	result, err := tomoxState.GetBids(tradingstate.GetTradingOrderBookHash(baseToken, quoteToken))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *PublicTomoXTransactionPoolAPI) GetAsks(ctx context.Context, baseToken, quoteToken common.Address) (map[*big.Int]*big.Int, error) {
	block := s.b.CurrentBlock()
	if block == nil {
		return nil, errors.New("Current block not found")
	}
	tomoxService := s.b.TomoxService()
	if tomoxService == nil {
		return nil, errors.New("TomoX service not found")
	}
	author, err := s.b.GetEngine().Author(block.Header())
	if err != nil {
		return nil, err
	}
	tomoxState, err := tomoxService.GetTradingState(block, author)
	if err != nil {
		return nil, err
	}
	result, err := tomoxState.GetAsks(tradingstate.GetTradingOrderBookHash(baseToken, quoteToken))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *PublicTomoXTransactionPoolAPI) GetInvests(ctx context.Context, lendingToken common.Address, term uint64) (map[*big.Int]*big.Int, error) {
	block := s.b.CurrentBlock()
	if block == nil {
		return nil, errors.New("Current block not found")
	}
	lendingService := s.b.LendingService()
	if lendingService == nil {
		return nil, errors.New("TomoX Lending service not found")
	}
	author, err := s.b.GetEngine().Author(block.Header())
	if err != nil {
		return nil, err
	}
	lendingState, err := lendingService.GetLendingState(block, author)
	if err != nil {
		return nil, err
	}
	result, err := lendingState.GetInvestings(lendingstate.GetLendingOrderBookHash(lendingToken, term))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *PublicTomoXTransactionPoolAPI) GetBorrows(ctx context.Context, lendingToken common.Address, term uint64) (map[*big.Int]*big.Int, error) {
	block := s.b.CurrentBlock()
	if block == nil {
		return nil, errors.New("Current block not found")
	}
	lendingService := s.b.LendingService()
	if lendingService == nil {
		return nil, errors.New("TomoX Lending service not found")
	}
	author, err := s.b.GetEngine().Author(block.Header())
	if err != nil {
		return nil, err
	}
	lendingState, err := lendingService.GetLendingState(block, author)
	if err != nil {
		return nil, err
	}
	result, err := lendingState.GetBorrowings(lendingstate.GetLendingOrderBookHash(lendingToken, term))
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetLendingTxMatchByHash returns lendingItems which have been processed at tx of the given txhash
func (s *PublicTomoXTransactionPoolAPI) GetLendingTxMatchByHash(ctx context.Context, hash common.Hash) ([]*lendingstate.LendingItem, error) {
	var tx *types.Transaction
	if tx, _, _, _ = core.GetTransaction(s.b.ChainDb(), hash); tx == nil {
		if tx = s.b.GetPoolTransaction(hash); tx == nil {
			return []*lendingstate.LendingItem{}, nil
		}
	}

	batch, err := lendingstate.DecodeTxLendingBatch(tx.Data())
	if err != nil {
		return []*lendingstate.LendingItem{}, err
	}
	return batch.Data, nil
}

// GetLiquidatedTradesByTxHash returns trades which closed by TomoX protocol at the tx of the give hash
func (s *PublicTomoXTransactionPoolAPI) GetLiquidatedTradesByTxHash(ctx context.Context, hash common.Hash) (lendingstate.FinalizedResult, error) {
	var tx *types.Transaction
	if tx, _, _, _ = core.GetTransaction(s.b.ChainDb(), hash); tx == nil {
		if tx = s.b.GetPoolTransaction(hash); tx == nil {
			return lendingstate.FinalizedResult{}, nil
		}
	}

	finalizedResult, err := lendingstate.DecodeFinalizedResult(tx.Data())
	if err != nil {
		return lendingstate.FinalizedResult{}, err
	}
	finalizedResult.TxHash = hash
	return finalizedResult, nil
}

func (s *PublicTomoXTransactionPoolAPI) GetLendingOrderById(ctx context.Context, lendingToken common.Address, term uint64, orderId uint64) (lendingstate.LendingItem, error) {
	lendingItem := lendingstate.LendingItem{}
	block := s.b.CurrentBlock()
	if block == nil {
		return lendingItem, errors.New("Current block not found")
	}
	lendingService := s.b.LendingService()
	if lendingService == nil {
		return lendingItem, errors.New("TomoX Lending service not found")
	}
	author, err := s.b.GetEngine().Author(block.Header())
	if err != nil {
		return lendingItem, err
	}
	lendingState, err := lendingService.GetLendingState(block, author)
	if err != nil {
		return lendingItem, err
	}
	lendingOrderBook := lendingstate.GetLendingOrderBookHash(lendingToken, term)
	orderIdHash := common.BigToHash(new(big.Int).SetUint64(orderId))
	lendingItem = lendingState.GetLendingOrder(lendingOrderBook, orderIdHash)
	if lendingItem.LendingId != orderId {
		return lendingItem, errors.New("Lending Item not found")
	}
	return lendingItem, nil
}

func (s *PublicTomoXTransactionPoolAPI) GetLendingTradeById(ctx context.Context, lendingToken common.Address, term uint64, tradeId uint64) (lendingstate.LendingTrade, error) {
	lendingItem := lendingstate.LendingTrade{}
	block := s.b.CurrentBlock()
	if block == nil {
		return lendingItem, errors.New("Current block not found")
	}
	lendingService := s.b.LendingService()
	if lendingService == nil {
		return lendingItem, errors.New("TomoX Lending service not found")
	}
	author, err := s.b.GetEngine().Author(block.Header())
	if err != nil {
		return lendingItem, err
	}
	lendingState, err := lendingService.GetLendingState(block, author)
	if err != nil {
		return lendingItem, err
	}
	lendingOrderBook := lendingstate.GetLendingOrderBookHash(lendingToken, term)
	tradeIdHash := common.BigToHash(new(big.Int).SetUint64(tradeId))
	lendingItem = lendingState.GetLendingTrade(lendingOrderBook, tradeIdHash)
	if lendingItem.TradeId != tradeId {
		return lendingItem, errors.New("Lending Item not found")
	}
	return lendingItem, nil
}
