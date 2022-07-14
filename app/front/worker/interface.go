package worker

import (
	"market_server/app/front/net"
	"market_server/common/datastruct"
)

type WorkerI interface {
	PublishSymbol(symbol_list []string, ws *net.WSInfo)
	PublishDepth(*datastruct.DepthQuote, *net.WSInfo)
	PublishTrade(*datastruct.RspTrade, *net.WSInfo)
	PublishKline(*datastruct.Kline, *net.WSInfo)
	PublishChangeinfo(*datastruct.ChangeInfo, *net.WSInfo)
	PublishHistKline(kline *datastruct.RspHistKline, ws *net.WSInfo)

	SubSymbol(ws *net.WSInfo)

	SubTrade(symbol string, ws *net.WSInfo) (string, bool)
	UnSubTrade(symbol string, ws *net.WSInfo)

	SubDepth(symbol string, ws *net.WSInfo) (string, bool)
	UnSubDepth(symbol string, ws *net.WSInfo)

	SubKline(req_kline_info *datastruct.ReqHistKline, ws *net.WSInfo) (string, bool)
	UnSubKline(req_kline_info *datastruct.ReqHistKline, ws *net.WSInfo)
}
