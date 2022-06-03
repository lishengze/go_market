package logic

import "github.com/shopspring/decimal"

// 交易所行情实体
type MarketDepth struct {
	//Direction    int             // 行情方向
	Price  decimal.Decimal // 价格
	Volume decimal.Decimal // 总数量

	TradedVolume decimal.Decimal // 成交量
	//Turnover     decimal.Decimal // 成交额
	//Fee          decimal.Decimal // 手续费
}

func (m *MarketDepth) GetRestVolume() decimal.Decimal {
	return m.Volume.Sub(m.TradedVolume)
}
