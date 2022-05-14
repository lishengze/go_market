package comm

import (
	"market_aggregate/pkg/conf"
	"market_aggregate/pkg/datastruct"
)

type NetServerI interface {
	Init(*conf.Config)
	PublishDepth(*datastruct.DepthQuote)
	PublishKline(*datastruct.Kline)
	PublishTrade(*datastruct.Trade)
}
