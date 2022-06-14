package worker

import (
	"market_server/app/front/net"
	"market_server/common/datastruct"
)

type WorkerI interface {
	PublishDepth(*datastruct.DepthQuote, *net.WSInfo)
	PublishTrade(*datastruct.Trade, *net.WSInfo)
	PublishKline(*datastruct.Kline, *net.WSInfo)
	PublishChangeinfo(*datastruct.ChangeInfo, *net.WSInfo)
	PublishHistKline(kline *datastruct.RspHistKline, ws *net.WSInfo)

	SubTrade(symbol string, ws *net.WSInfo)
	// UnSubTrade(symbol string)

	SubDepth(symbol string, ws *net.WSInfo)
	// UnSubDepth(symbol string)

	SubKline(req_kline_info *datastruct.ReqHistKline, ws *net.WSInfo)
	// UnSubKline(req_kline_info *datastruct.ReqHistKline)
}
