package extools

import "exterior-interactor/pkg/exchangeapi/exmodel"

type (
	SymbolManager interface {
		PullAllSymbol() error
		GetSymbol(stdFormat exmodel.StdSymbol) (*exmodel.Symbol, error)
		GetAllSymbol() []*exmodel.Symbol
		Convert(exFormat string, apiType exmodel.ApiType) (*exmodel.Symbol, error)
	}

	MarketTradeManager interface {
		Sub(symbols ...exmodel.StdSymbol)
		OutputCh() <-chan *exmodel.StreamMarketTrade
	}

	DepthManager interface {
		Sub(symbols ...exmodel.StdSymbol)
		OutputCh() <-chan *exmodel.StreamDepth
	}

	TradeManager interface {
		PlaceOrder(req exmodel.PlaceOrderReq) (*exmodel.PlaceOrderRsp, error)
		CancelOrder(req exmodel.CancelOrderReq) (*exmodel.CancelOrderRsp, error)
		QueryOrder(req exmodel.QueryOrderReq) (*exmodel.Order, error)
		QueryTrades(req exmodel.QueryTradeReq) ([]*exmodel.Trade, error)
		OutputUpdateCh() <-chan *exmodel.OrderTradesUpdate
		Close()
	}

	WalletManager interface {
		QueryBalance(req exmodel.QueryBalanceReq) (*exmodel.QueryBalanceRsp, error)
	}
)
