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
)
