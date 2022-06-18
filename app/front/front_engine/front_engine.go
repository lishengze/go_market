package front_engine

import (
	"market_server/app/front/config"
	"market_server/app/front/net"
	"market_server/app/front/worker"
	"market_server/common/datastruct"
	"market_server/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type FrontEngine struct {
	sub_data    *SubData
	next_worker worker.WorkerI
	config      *config.Config
}

func NewFrontEngine(config *config.Config) *FrontEngine {

	rst := &FrontEngine{
		config:   config,
		sub_data: NewSubData(),
	}
	return rst
}

func (f *FrontEngine) SetNextWorker(next_worker worker.WorkerI) {
	f.next_worker = next_worker
}

func (f *FrontEngine) Start() {
	logx.Infof("FrontEngine Start!")
}

func (f *FrontEngine) PublishSymbol(symbol_list []string, ws *net.WSInfo) {
	if ws != nil {

	} else {
		symbol_pub_list := f.sub_data.GetSymbolPubInfoList(symbol_list)

		logx.Statf("symbol_pub_list: %+v \n", symbol_pub_list)
	}

}

func (f *FrontEngine) PublishDepth(depth *datastruct.DepthQuote, ws *net.WSInfo) {
	if ws != nil {

	} else {
		depth_pub_list := f.sub_data.GetDepthPubInfoList(depth)

		logx.Statf("depth_pub_list: %+v \n", depth_pub_list)
	}

}

func (f *FrontEngine) PublishTrade(trade *datastruct.Trade, ws *net.WSInfo) {
	if ws != nil {

	} else {
		trade_pub_list := f.sub_data.GetTradePubInfoList(trade)

		logx.Statf("trade_pub_list: %+v \n", trade_pub_list)
	}

}

func (f *FrontEngine) PublishKline(kline *datastruct.Kline, ws *net.WSInfo) {
	if ws != nil {

	} else {
		kline_pub_list := f.sub_data.GetKlinePubInfoList(kline)

		logx.Statf("kline_pub_list: %+v \n", kline_pub_list)
	}

}

func (f *FrontEngine) PublishChangeinfo(change_info *datastruct.ChangeInfo, ws *net.WSInfo) {

}

func (f *FrontEngine) PublishHistKline(klines *datastruct.RspHistKline, ws *net.WSInfo) {
	// d.publish_kline(kline)
	f.sub_data.ProcessKlineHistData(klines)

	if ws != nil {
		logx.Statf("rsp klines: %+v \n", klines)
	}
	// publish his kline to client;
}

func (f *FrontEngine) SubSymbol(ws *net.WSInfo) {
	f.sub_data.SubSymbol(ws)
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

func (f *FrontEngine) TestSub(test_map map[string]struct{}) {
	ws_info := &net.WSInfo{
		ID: util.UTCNanoTime(),
	}

	symbol := "BTC_USDT"

	f.SubSymbol(ws_info)

	if _, ok := test_map[datastruct.DEPTH_TYPE]; ok {
		logx.Statf("Sub Depth: %s, ws_info: %+v\n", symbol, ws_info)
		f.SubDepth(symbol, ws_info)
	}

	if _, ok := test_map[datastruct.TRADE_TYPE]; ok {
		logx.Statf("Sub Trade: %s, ws_info: %+v\n", symbol, ws_info)
		f.SubTrade(symbol, ws_info)
	}

	if _, ok := test_map[datastruct.KLINE_TYPE]; ok {
		req_kline_info := datastruct.ReqHistKline{
			Symbol:    symbol,
			Exchange:  datastruct.BCTS_EXCHANGE,
			StartTime: 0,
			EndTime:   0,
			Count:     1000,
			Frequency: datastruct.SECS_PER_MIN,
		}

		logx.Statf("Sub Kline: %+v, ws_info: %+v\n", req_kline_info, ws_info)

		f.SubKline(&req_kline_info, ws_info)
	}

}
