package front_engine

import (
	"market_server/app/front/config"
	"market_server/app/front/net"
	"market_server/app/front/worker"
	"market_server/common/datastruct"
	"market_server/common/util"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

type FrontEngine struct {
	sub_data    *SubData
	next_worker worker.WorkerI
	config      *config.Config

	IsTest bool
}

func NewFrontEngine(config *config.Config) *FrontEngine {

	rst := &FrontEngine{
		config:   config,
		sub_data: NewSubData(),
	}
	return rst
}

func (f *FrontEngine) SetTestFlag(value bool) {
	f.IsTest = value
}

func (f *FrontEngine) SetNextWorker(next_worker worker.WorkerI) {
	f.next_worker = next_worker
}

func (f *FrontEngine) Start() {
	logx.Infof("FrontEngine Start!")
}

// func catch_exp(msg []byte, ws *net.WSInfo) {
// 	errMsg := recover()
// 	if errMsg != nil {
// 		fmt.Println("This is catch_exp func")
// 		logx.Errorf("catch_exp OriginalMsg: %s, WSInfo: %+v\n", msg, *ws)
// 		logx.Errorf("errMsg: %+v \n", errMsg)
// 		fmt.Println(errMsg)
// 	}

// }

func (f *FrontEngine) PublishSymbol(symbol_list []string, ws *net.WSInfo) {
	logx.Infof("PubSymbolist: %+v", symbol_list)

	if ws != nil {
		byte_data := NewSymbolListMsg(symbol_list)

		if ws.IsAlive() {
			err := ws.SendMsg(1, byte_data)
			if err != nil {
				logx.Errorf("PublishSymbol err: %+v \n", err)
			}
		} else {
			logx.Infof("ws:%+v is not alive", ws)
			f.sub_data.UnSubSymbol(ws)
		}

	} else {
		symbol_pub_list := f.sub_data.GetSymbolPubInfoList(symbol_list)

		for _, info := range symbol_pub_list {
			logx.Statf("symbol_pub_info: %s \n", info.String())
			if info.ws_info.IsAlive() {
				err := info.ws_info.SendMsg(1, info.data)
				if err != nil {
					logx.Errorf("PublishSymbol err: %+v \n", err)
				}
			} else {
				logx.Infof("ws:%+v is not alive", ws)
				f.sub_data.UnSubSymbol(ws)
			}
		}
	}

}

func (f *FrontEngine) PublishDepth(depth *datastruct.DepthQuote, ws *net.WSInfo) {
	defer func(depth *datastruct.DepthQuote, ws *net.WSInfo) {
		errMsg := recover()
		if errMsg != nil {
			logx.Errorf("PublishDepth depth: %+v, ws_info: %+v\n", depth, ws)
			logx.Errorf("errMsg: %+v \n", errMsg)
		}
	}(depth, ws)

	if ws != nil {
		if ws.IsAlive() {
			byte_data := NewDepthJsonMsg(depth)
			err := ws.SendMsg(websocket.TextMessage, byte_data)

			if err != nil {
				logx.Errorf("PublishDepth err: %+v \n", err)
			}
		} else {
			logx.Infof("ws:%+v is not alive", ws)
			f.sub_data.UnSubDepth(depth.Symbol, ws)
		}
	} else {
		depth_pub_list := f.sub_data.GetDepthPubInfoList(depth)

		for _, info := range depth_pub_list {
			// logx.Slowf("depth_pub_info: %s \n", info.String())
			if info.ws_info.IsAlive() {
				err := info.ws_info.SendMsg(1, info.data)
				if err != nil {
					logx.Errorf("PublishDepth err: %+v \n", err)
				}
			} else {
				logx.Infof("ws:%+v is not alive", info.ws_info)
				f.sub_data.UnSubDepth(depth.Symbol, info.ws_info)
			}
		}
	}

}

func (f *FrontEngine) PublishTrade(trade *datastruct.Trade, change_info *datastruct.ChangeInfo, ws *net.WSInfo) {

	defer func(trade *datastruct.Trade, change_info *datastruct.ChangeInfo, ws *net.WSInfo) {
		errMsg := recover()
		if errMsg != nil {
			// fmt.Println("This is catch_exp func")
			logx.Errorf("PublishTrade trade: %+v, change_info: %+v, ws_info: %+v\n",
				trade, change_info, ws)
			logx.Errorf("errMsg: %+v \n", errMsg)
			// fmt.Println(errMsg)
		}
	}(trade, change_info, ws)

	logx.Slowf("PublishTrade:%+v, %+v\n", trade, change_info)

	if ws != nil {
		if ws.IsAlive() {
			// logx.Infof("ws:%+v is not alive")
			byte_data := NewTradeJsonMsg(trade, change_info)
			err := ws.SendMsg(websocket.TextMessage, byte_data)

			if err != nil {
				logx.Errorf("PublishDepth err: %+v \n", err)
			}
		} else {
			f.sub_data.UnSubTrade(trade.Symbol, ws)
			logx.Infof("ws:%+v is not alive", ws)
		}
	} else {
		trade_pub_list := f.sub_data.GetTradePubInfoList(trade, change_info)
		// logx.Info("After GetTradePubInfoList")

		for _, info := range trade_pub_list {
			logx.Slowf("trade_pub_info: %s \n", info.String())
			if info.ws_info.IsAlive() {
				err := info.ws_info.SendMsg(websocket.TextMessage, info.data)
				if err != nil {
					logx.Errorf("PublishTrade err: %+v \n", err)
				}
			} else {
				logx.Infof("ws:%+v is not alive", info.ws_info)
				f.sub_data.UnSubTrade(trade.Symbol, info.ws_info)
			}
		}
	}

}

func (f *FrontEngine) PublishKline(kline *datastruct.Kline, ws *net.WSInfo) {
	defer func(kline *datastruct.Kline, ws *net.WSInfo) {
		errMsg := recover()
		if errMsg != nil {
			// fmt.Println("This is catch_exp func")
			logx.Errorf("PublishUpdateKline kline: %+v, ws_info: %+v\n", kline, ws)
			logx.Errorf("errMsg: %+v \n", errMsg)
			// fmt.Println(errMsg)
		}
	}(kline, ws)

	if ws != nil {
		cur_req := &datastruct.ReqHistKline{
			Symbol:    kline.Symbol,
			Frequency: uint32(kline.Resolution),
		}

		if ws.IsAlive() {

			logx.Infof("[PubKline]: %s", kline.String())
			byte_data := NewKlineUpdateJsonMsg(kline)
			err := ws.SendMsg(websocket.TextMessage, byte_data)

			if err != nil {
				logx.Errorf("PublishDepth err: %+v \n", err)
			}
		} else {
			logx.Infof("ws:%+v is not alive", ws)
			f.sub_data.UnSubKline(cur_req, ws)
		}
	} else {

		kline_pub_list := f.sub_data.GetKlinePubInfoList(kline)
		if nil != kline_pub_list {
			for _, info := range kline_pub_list {
				cur_req := &datastruct.ReqHistKline{
					Symbol:    info.Symbol,
					Frequency: uint32(info.Resolution),
				}

				// logx.Infof("kline_pub_info: %s \n", info.String())
				if info.ws_info.IsAlive() {
					err := info.ws_info.SendMsg(1, info.data)
					if err != nil {
						logx.Errorf("PublishKline err: %+v \n", err)
					}
				} else {
					logx.Infof("ws:%d is not alive", info.ws_info.ID)
					f.sub_data.UnSubKline(cur_req, info.ws_info)
				}
			}
		}

	}

}

func (f *FrontEngine) PublishHistKline(klines *datastruct.RspHistKline, ws *net.WSInfo) {
	defer func(klines *datastruct.RspHistKline, ws *net.WSInfo) {
		errMsg := recover()
		if errMsg != nil {
			// fmt.Println("This is catch_exp func")
			logx.Errorf("PublishHistKline klines: %+v, ws_info: %+v\n", klines, ws)
			logx.Errorf("errMsg: %+v \n", errMsg)
			// fmt.Println(errMsg)
		}
	}(klines, ws)

	f.sub_data.ProcessKlineHistData(klines)

	if ws != nil {
		logx.Infof("PublishHistKline: %s", klines.SimpleTimeList())

		byte_data := NewHistKlineJsonMsg(klines)
		if ws.IsAlive() {
			err := ws.SendMsg(1, byte_data)
			if err != nil {
				logx.Errorf("PublishHistKline err: %+v \n", err)
			}
		} else {
			logx.Infof("ws:%+v is not alive", ws)
			f.sub_data.UnSubKline(klines.ReqInfo, ws)
		}
	} else {
		logx.Errorf("PublishHistKline ws is null\n")
	}
	// publish his kline to client;
}

func (f *FrontEngine) PublishChangeinfo(change_info *datastruct.ChangeInfo, ws *net.WSInfo) {

}

func (f *FrontEngine) SubSymbol(ws *net.WSInfo) {
	f.sub_data.SubSymbol(ws)
	f.next_worker.SubSymbol(ws)
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
	f.next_worker.SubKline(req_kline_info, ws)
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
			Frequency: datastruct.SECS_PER_MIN * 5,
		}

		logx.Statf("Sub Kline: %+v, ws_info: %+v\n", req_kline_info, ws_info)

		f.SubKline(&req_kline_info, ws_info)
	}

}
