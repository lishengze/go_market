package exmodel

const (
	SymbolTypeSpot         SymbolType = "SPOT"
	SymbolTypeCoinPerp     SymbolType = "COIN_PERP"
	SymbolTypeCoinDelivery SymbolType = "COIN_DELIVERY"
	SymbolTypeUsdtPerp     SymbolType = "USDT_PERP"
	SymbolTypeUsdtDelivery SymbolType = "USDT_DELIVERY"
)

type (
	/*
		StdSymbol e.g.
		BTC_USDT
		BTC_USD_220325
		BTC_USDT_PERP
	*/
	StdSymbol string

	SymbolType string

	Symbol struct {
		Exchange      Exchange
		Type          SymbolType
		StdSymbol     StdSymbol
		BaseCurrency  Currency
		QuoteCurrency Currency
		ApiType       ApiType
		ExFormat      string
		VolumeScale   string // 下单最小数量间隔
		PriceScale    string // 下单最小价格间隔
		MinVolume     string // 最小 量
		ContractSize  string // 合约面值
	}
)

func ConvertStrsToStdSymbols(s []string) []StdSymbol {
	var symbols = make([]StdSymbol, 0)

	for _, symbol := range s {
		symbols = append(symbols, StdSymbol(symbol))
	}
	return symbols
}

func (o StdSymbol) String() string {
	return string(o)
}

func (o ApiType) String() string {
	return string(o)
}

func (o SymbolType) String() string {
	return string(o)
}
