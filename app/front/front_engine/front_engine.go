package front_engine

import (
	"market_server/app/front/net"
	"market_server/app/front/worker"
	"market_server/common/datastruct"
)

type FrontEngine struct {
	sub_data *SubData

	next_worker worker.WorkerI
}

func (f *FrontEngine) PublishDepth(depth *datastruct.DepthQuote, ws *net.WSInfo) {
	if ws != nil {

	} else {
		// depth_pub_list := f.sub_data.GetDepthPubInfoList(depth)
	}

}

func (f *FrontEngine) PublishTrade(trade *datastruct.Trade, ws *net.WSInfo) {
	if ws != nil {

	} else {
		// trade_pub_list := f.sub_data.GetTradePubInfoList(trade)
	}

}

func (f *FrontEngine) PublishKline(kline *datastruct.Kline, ws *net.WSInfo) {
	if ws != nil {

	} else {
		// kline_pub_list := f.sub_data.GetKlinePubInfoList(kline)
	}

}

func (f *FrontEngine) PublishChangeinfo(change_info *datastruct.ChangeInfo, ws *net.WSInfo) {

}

func (f *FrontEngine) PublishHistKline(klines *datastruct.RspHistKline, ws *net.WSInfo) {
	// d.publish_kline(kline)
	f.sub_data.ProcessKlineHistData(klines)

	// publish his kline to client;
}

func (f *FrontEngine) SubTrade(symbol string, ws *net.WSInfo) {
	f.sub_data.SubTrade(symbol, ws)
	f.next_worker.SubTrade(symbol, ws)
}

func (f *FrontEngine) UnSubTrade(symbol string, ws *net.WSInfo) {
	f.sub_data.UnSubTrade(symbol, ws)
}

func (f *FrontEngine) SubDepth(symbol string, ws *net.WSInfo) {
	f.sub_data.SubDepth(symbol, ws)
	f.next_worker.SubDepth(symbol, ws)
}

func (f *FrontEngine) UnSubDepth(symbol string, ws *net.WSInfo) {
	f.sub_data.UnSubDepth(symbol, ws)
}

func (f *FrontEngine) SubKline(req_kline_info *datastruct.ReqHistKline, ws *net.WSInfo) {
	f.sub_data.SubKline(req_kline_info, ws)
}

func (f *FrontEngine) UnSubKline(req_kline_info *datastruct.ReqHistKline, ws *net.WSInfo) {
	f.sub_data.UnSubKline(req_kline_info, ws)
}
