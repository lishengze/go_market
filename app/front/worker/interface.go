package data_worker

import "market_server/common/datastruct"

type WorkerI interface {
	publish_depth(*datastruct.DepthQuote)
	publish_trade(*datastruct.Trade)
	publish_kline(*datastruct.Kline)

	SubTrade(symbol string) *datastruct.Trade
	UnSubTrade(symbol string)

	SubDepth(symbol string) *datastruct.DepthQuote
	UnSubDepth(symbol string)

	SubKline(req_kline_info *datastruct.ReqHistKline) *datastruct.HistKline
	UnSubKline(req_kline_info *datastruct.ReqHistKline)
}
