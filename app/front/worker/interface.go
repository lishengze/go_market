package worker

import (
	"market_server/common/datastruct"
)

type WorkerI interface {
	PublishDepth(*datastruct.DepthQuote)
	PublishTrade(*datastruct.Trade)
	PublishKline(*datastruct.Kline)
	PublishChangeinfo(*datastruct.ChangeInfo)
	PublishHistKline(kline *datastruct.HistKline)

	SubTrade(symbol string) *datastruct.Trade
	UnSubTrade(symbol string)

	SubDepth(symbol string) *datastruct.DepthQuote
	UnSubDepth(symbol string)

	SubKline(req_kline_info *datastruct.ReqHistKline) *datastruct.HistKline
	UnSubKline(req_kline_info *datastruct.ReqHistKline)
}
