package front_engine

import (
	"fmt"
	"market_server/app/front/net"
	"market_server/common/datastruct"
	"market_server/common/util"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

type SubData struct {
	SymbolInfo *SymbolSubInfo
	DepthInfo  *DepthSubInfo
	TradeInfo  *TradeSubInfo
	KlineInfo  *KlineSubInfo
}

func NewSubData() *SubData {
	return &SubData{
		SymbolInfo: NewSymbolSubInfo(),
		DepthInfo:  NewDepthSubInfo(),
		TradeInfo:  NewTradeSubInfo(),
		KlineInfo:  NewKlineSubInfo(),
	}
}

func (s *SubData) GetSymbolPubInfoList(symbollist []string) []*SymbolPubInfo {
	var rst []*SymbolPubInfo

	s.SymbolInfo.mutex.Lock()
	defer s.SymbolInfo.mutex.Unlock()

	iter := s.SymbolInfo.ws_info.Iterator()

	byte_data := NewSymbolListMsg(symbollist)
	for iter.Begin(); iter.Next(); {
		rst = append(rst, &SymbolPubInfo{
			ws_info: iter.Value().(*net.WSInfo),
			data:    byte_data,
		})
	}

	return rst
}

func (s *SubData) GetDepthPubInfoList(depth *datastruct.DepthQuote) []*DepthPubInfo {
	var rst []*DepthPubInfo

	s.DepthInfo.mutex.Lock()
	defer s.DepthInfo.mutex.Unlock()

	//
	// if len(s.DepthInfo.Info) > 0 {
	// 	logx.Slowf("CurDepthSubInfo: %s", s.DepthInfo.String())
	// }

	byte_data := NewDepthJsonMsg(depth)
	if sub_tree, ok := s.DepthInfo.Info[depth.Symbol]; ok {
		sub_tree_iter := sub_tree.Iterator()
		sub_tree_iter.Begin()
		for sub_tree_iter.Next() {
			rst = append(rst, &DepthPubInfo{
				ws_info: sub_tree_iter.Value().(*net.WSInfo),
				data:    byte_data,
			})

		}
	} else {
		// logx.Errorf("Depth Symobl: %s Not Subed!", depth.Symbol)
	}

	return rst
}

func (s *SubData) GetTradePubInfoList(trade *datastruct.RspTrade) []*TradePubInfo {
	var rst []*TradePubInfo

	s.TradeInfo.mutex.Lock()
	defer s.TradeInfo.mutex.Unlock()

	// if len(s.TradeInfo.Info) > 0 {
	// 	logx.Statf("CurTradeSubInfo: %s", s.TradeInfo.String())
	// }

	byte_data := NewTradeJsonMsg(trade)
	if sub_tree, ok := s.TradeInfo.Info[trade.TradeData.Symbol]; ok {
		sub_tree_iter := sub_tree.Iterator()
		sub_tree_iter.Begin()
		for sub_tree_iter.Next() {
			rst = append(rst, &TradePubInfo{
				ws_info: sub_tree_iter.Value().(*net.WSInfo),
				data:    byte_data,
			})
		}
	} else {
		// logx.Errorf("Trade Symobl: %s Not Subed!", trade.Symbol)
	}

	return rst
}

func (s *SubData) GetKlinePubInfoListAtom(sub_info *KlineSubItem, pub_kline *datastruct.Kline) []*KlinePubInfo {
	defer util.CatchExp("GetKlinePubInfoListAtom")
	var rst []*KlinePubInfo

	byte_data := NewKlineUpdateJsonMsg(pub_kline)

	// if pub_kline.Resolution > datastruct.NANO_PER_SECS {
	// 	pub_kline.Resolution = pub_kline.Resolution / datastruct.NANO_PER_SECS
	// }

	sub_tree := sub_info.ws_info
	sub_tree_iter := sub_tree.Iterator()

	for sub_tree_iter.Begin(); sub_tree_iter.Next(); {
		rst = append(rst, &KlinePubInfo{
			ws_info:    sub_tree_iter.Value().(*net.WSInfo),
			data:       byte_data,
			Symbol:     pub_kline.Symbol,
			Resolution: pub_kline.Resolution,
		})
	}

	return rst
}

func (s *SubData) GetKlinePubInfoList(kline *datastruct.Kline) []*KlinePubInfo {
	defer util.CatchExp("GetKlinePubInfoList")

	var rst []*KlinePubInfo

	s.KlineInfo.mutex.Lock()
	defer s.KlineInfo.mutex.Unlock()

	if _, ok := s.KlineInfo.Info[kline.Symbol]; !ok {
		return rst
	}

	for _, sub_info := range s.KlineInfo.Info[kline.Symbol] {
		cur_pub_list := s.GetKlinePubInfoListAtom(sub_info, kline)
		rst = append(rst, cur_pub_list...)
	}

	return rst
}

func (s *SubData) SubSymbol(ws *net.WSInfo) {
	defer util.CatchExp("SubSymbol")

	s.SymbolInfo.mutex.Lock()
	defer s.SymbolInfo.mutex.Unlock()

	if _, ok := s.SymbolInfo.ws_info.Get(ws.ID); !ok {
		s.SymbolInfo.ws_info.Put(ws.ID, ws)
	}

	logx.Infof("SubSymbol %s, %s", ws.String(), s.SymbolInfo.String())

	iter := s.SymbolInfo.ws_info.Iterator()
	invalid_sub_list := make([]int64, s.SymbolInfo.ws_info.Size())
	for iter.Begin(); iter.Next(); {
		ws_info := iter.Value().(*net.WSInfo)

		if !ws_info.IsAlive() {
			invalid_sub_list = append(invalid_sub_list, ws_info.ID)
		}
	}

	for _, id := range invalid_sub_list {
		s.SymbolInfo.ws_info.Remove(id)
		logx.Infof("SymbolInfo Remove: %d", id)
	}
}

func (s *SubData) UnSubSymbol(ws *net.WSInfo) {
	defer util.CatchExp("UnSubSymbol")

	s.SymbolInfo.mutex.Lock()
	defer s.SymbolInfo.mutex.Unlock()

	if _, ok := s.SymbolInfo.ws_info.Get(ws.ID); !ok {
		return
	}

	s.SymbolInfo.ws_info.Remove(ws.ID)

	logx.Infof("SymbolInfo Remove : %d", ws.ID)
}

func (s *SubData) SubDepth(symbol string, ws *net.WSInfo) *datastruct.DepthQuote {
	defer util.CatchExp("SubDepth")

	s.DepthInfo.mutex.Lock()
	defer s.DepthInfo.mutex.Unlock()

	if _, ok := s.DepthInfo.Info[symbol]; !ok {
		s.DepthInfo.Info[symbol] = treemap.NewWith(utils.Int64Comparator)
	}

	s.DepthInfo.Info[symbol].Put(ws.ID, ws)

	logx.Infof("SubDepth %s, %s, %s", symbol, ws.String(), s.DepthInfo.String())

	return nil
}

func (s *SubData) UnSubDepth(symbol string, ws *net.WSInfo) {
	defer util.CatchExp("UnSubDepth")

	s.DepthInfo.mutex.Lock()
	defer s.DepthInfo.mutex.Unlock()
	if _, ok := s.DepthInfo.Info[symbol]; !ok {
		return
	}

	s.DepthInfo.Info[symbol].Remove(ws.ID)

	if s.DepthInfo.Info[symbol].Size() == 0 {
		delete(s.DepthInfo.Info, symbol)
	}

	logx.Infof("UnSubDepth %s : %+v", symbol, ws)
}

func (s *SubData) SubTrade(symbol string, ws *net.WSInfo) {
	defer util.CatchExp("SubTrade")

	s.TradeInfo.mutex.Lock()
	defer s.TradeInfo.mutex.Unlock()

	if _, ok := s.TradeInfo.Info[symbol]; !ok {
		s.TradeInfo.Info[symbol] = treemap.NewWith(utils.Int64Comparator)
	}

	s.TradeInfo.Info[symbol].Put(ws.ID, ws)

	logx.Infof("SubTrade %s, %s, %s", symbol, ws.String(), s.TradeInfo.String())
}

func (s *SubData) UnSubTrade(symbol string, ws *net.WSInfo) {
	defer util.CatchExp("UnSubTrade")

	s.TradeInfo.mutex.Lock()
	defer s.TradeInfo.mutex.Unlock()

	if _, ok := s.TradeInfo.Info[symbol]; !ok {
		return
	}

	s.TradeInfo.Info[symbol].Remove(ws.ID)

	if s.TradeInfo.Info[symbol].Size() == 0 {
		delete(s.TradeInfo.Info, symbol)
	}

	logx.Infof("UnSubTrade Remove %s : %+v", symbol, ws)
}

func (s *SubData) SetSubHistKlineFlag(req_kline_info *datastruct.ReqHistKline, ws *net.WSInfo) {
	defer util.CatchExp("SubKline")

	s.KlineInfo.mutex.Lock()
	defer s.KlineInfo.mutex.Unlock()

	if _, ok := s.KlineInfo.Info[req_kline_info.Symbol]; !ok {

		logx.Errorf("SetSubHistKlineFlag symbol %s not subed!\n", req_kline_info.Symbol)
		return
	}

	if _, ok := s.KlineInfo.Info[req_kline_info.Symbol][req_kline_info.Frequency]; !ok {
		logx.Errorf("SetSubHistKlineFlag symbol: %s, resolution: %d not subed!\n",
			req_kline_info.Symbol, req_kline_info.Frequency)
		return
	}

	value, ok := s.KlineInfo.Info[req_kline_info.Symbol][req_kline_info.Frequency].ws_info.Get(ws.ID)

	if ok {
		fmt.Println(value)
	} else {
		logx.Errorf("SetSubHistKlineFlag symbol: %s, resolution: %d, ws: %d, not subed!\n",
			req_kline_info.Symbol, req_kline_info.Frequency, ws.ID)
		return
	}

	logx.Infof("SubKline After Sub %s, %s\n", req_kline_info.String(), ws.String())
}

func (s *SubData) SubKline(req_kline_info *datastruct.ReqHistKline, ws *net.WSInfo) {
	defer util.CatchExp("SubKline")

	s.KlineInfo.mutex.Lock()
	defer s.KlineInfo.mutex.Unlock()

	if _, ok := s.KlineInfo.Info[req_kline_info.Symbol]; !ok {

		s.KlineInfo.Info[req_kline_info.Symbol] = make(map[uint64]*KlineSubItem)

		s.KlineInfo.Info[req_kline_info.Symbol][req_kline_info.Frequency] = &KlineSubItem{
			ws_info: treemap.NewWith(utils.Int64Comparator),
		}
	}

	if _, ok := s.KlineInfo.Info[req_kline_info.Symbol][req_kline_info.Frequency]; !ok {
		s.KlineInfo.Info[req_kline_info.Symbol][req_kline_info.Frequency] = &KlineSubItem{
			ws_info: treemap.NewWith(utils.Int64Comparator),
		}
	}

	s.KlineInfo.Info[req_kline_info.Symbol][req_kline_info.Frequency].ws_info.Put(ws.ID, ws)

	info := "SubInfo : \n"
	for symbol, data := range s.KlineInfo.Info {
		for resolution, ws_info := range data {
			info = info + fmt.Sprintf("%s.%d: %s \n", symbol, resolution, ws_info.String())
		}
	}

	logx.Infof("\n After SubKline After Sub %s, %d\n%s", req_kline_info.String(), req_kline_info.Frequency, info)
}

func (s *SubData) UnSubKline(req_kline_info *datastruct.ReqHistKline, ws *net.WSInfo) {
	defer util.CatchExp("UnSubKline")

	s.KlineInfo.mutex.Lock()
	defer s.KlineInfo.mutex.Unlock()

	// logx.Infof("WSEngine UnSubInfo: %s,%d", symbol, resolution)

	if _, ok := s.KlineInfo.Info[req_kline_info.Symbol]; !ok {
		return
	} else if _, ok := s.KlineInfo.Info[req_kline_info.Symbol][req_kline_info.Frequency]; !ok {
		return
	}

	if s.KlineInfo.Info[req_kline_info.Symbol][req_kline_info.Frequency].ws_info != nil {
		s.KlineInfo.Info[req_kline_info.Symbol][req_kline_info.Frequency].ws_info.Remove(ws.ID)

		if s.KlineInfo.Info[req_kline_info.Symbol][req_kline_info.Frequency].ws_info.Size() == 0 {
			delete(s.KlineInfo.Info[req_kline_info.Symbol], req_kline_info.Frequency)

			if len(s.KlineInfo.Info[req_kline_info.Symbol]) == 0 {
				delete(s.KlineInfo.Info, req_kline_info.Symbol)
			}
		}

		info := "SubInfo : \n"
		for symbol, data := range s.KlineInfo.Info {
			for resolution, ws_info := range data {
				info = info + fmt.Sprintf("%s.%d: %s \n", symbol, resolution, ws_info.String())
			}
		}

		logx.Infof("\n After UnSubKline Remove %s, %d\n%s", req_kline_info.Symbol, int(req_kline_info.Frequency), info)
	} else {
		logx.Infof("KlineInfo, %s.%d is already nil!", req_kline_info.Symbol, req_kline_info.Frequency)
	}

}
