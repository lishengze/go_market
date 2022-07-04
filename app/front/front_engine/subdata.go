package front_engine

import (
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

func (s *SubData) GetTradePubInfoList(trade *datastruct.Trade, change_info *datastruct.ChangeInfo) []*TradePubInfo {
	var rst []*TradePubInfo

	s.TradeInfo.mutex.Lock()
	defer s.TradeInfo.mutex.Unlock()

	// if len(s.TradeInfo.Info) > 0 {
	// 	logx.Statf("CurTradeSubInfo: %s", s.TradeInfo.String())
	// }

	byte_data := NewTradeJsonMsg(trade, change_info)
	if sub_tree, ok := s.TradeInfo.Info[trade.Symbol]; ok {
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

func (s *SubData) UpdateKlineCacheData(kline *datastruct.Kline) {
	if _, ok := s.KlineInfo.Info[kline.Symbol]; !ok {
		s.KlineInfo.Info[kline.Symbol] = make(map[int]*KlineSubItem)
	}

	if _, ok := s.KlineInfo.Info[kline.Symbol][int(kline.Resolution)]; !ok {
		s.KlineInfo.Info[kline.Symbol][int(kline.Resolution)].cache_data = kline
	} else {
		cache_kline := s.KlineInfo.Info[kline.Symbol][int(kline.Resolution)].cache_data

		if kline.Time-cache_kline.Time >= int64(kline.Resolution) {
			s.KlineInfo.Info[kline.Symbol][int(kline.Resolution)].cache_data = kline
		}
	}
}

/*
判断是否是之前的 k 线数据结束;
逻辑描述:
	1. 原始的K 线数据，是由trade数据聚合而成，比如15分频的 9:45 的数据，是由 [9:45, 10:00) 的数据聚合而成;
	2. 若是 用15分频的数据聚合60分频, 9:45 的原始K线数据来时，时间已经到了10:00, 需要立刻结合当前的累计的高开低收，然后发布数据;
	3. 所以对于新来的原始k 线，有三种状态: 旧的聚合区间的最后一个数据， 新的聚合区间的第一个数据, 旧的聚合区间的中间的数据，
	4. 旧的聚合区间的最后一个数据： (kline.resolution + kline.time) % cur_resolution == 0;
		4.1 需要判断是否需要进行累计计算-不同聚合：kline.resolution == cur_resolution, 可以直接发布;
		4.2 否则: 与当前累计数据一起计算高低收， 发布;
	5. 新的聚合区间的第一个数据: kline.resolution % cur_resolution == 0;
		5.1 已经过滤掉了不用聚合直接发布的情况;
		5.2 将当前 cache 的数据更新为当前的数据;
	6. 旧的聚合区间的中间的数据
		根据规则更新 cache 的 高低收;
*/

func (s *SubData) GetKlinePubInfoList(kline *datastruct.Kline) []*KlinePubInfo {
	var rst []*KlinePubInfo

	s.KlineInfo.mutex.Lock()
	defer s.KlineInfo.mutex.Unlock()

	if _, ok := s.KlineInfo.Info[kline.Symbol]; !ok {
		return rst
	}

	for resolution, sub_info := range s.KlineInfo.Info[kline.Symbol] {
		if sub_info.cache_data == nil {
			logx.Errorf("Hisk Kline %s, %d , cache_data empty", kline.Symbol, resolution)
			continue
		}

		cache_kline := sub_info.cache_data

		if kline.Time <= cache_kline.Time {
			logx.Errorf("NewKLineTime: %d, CachedKlineTime: %d")
			continue
		}

		if datastruct.IsOldKlineEnd(kline, int64(resolution)) {
			logx.Slowf("Old Kline End: %s", kline.String())

			var pub_kline *datastruct.Kline
			if kline.Resolution != resolution {
				cache_kline.Close = kline.Close
				cache_kline.Low = util.MinFloat64(cache_kline.Low, kline.Low)
				cache_kline.High = util.MaxFloat64(cache_kline.High, kline.High)
				cache_kline.Volume += kline.Volume

				pub_kline = datastruct.NewKlineWithKline(cache_kline)
			} else {
				pub_kline = datastruct.NewKlineWithKline(kline)
			}

			// logx.Statf("CurKlineInfo %s", s.KlineInfo.String())
			byte_data := NewKlineUpdateJsonMsg(pub_kline)
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

			s.KlineInfo.Info[kline.Symbol][resolution].cache_data = datastruct.NewKlineWithKline(pub_kline)

		} else if datastruct.IsNewKlineStart(kline, int64(resolution)) {
			logx.Slowf("New Kline Start: %s", kline.String())

			s.KlineInfo.Info[kline.Symbol][resolution].cache_data = datastruct.NewKlineWithKline(kline)
			s.KlineInfo.Info[kline.Symbol][resolution].cache_data.Resolution = resolution
		} else {
			cache_kline.Close = kline.Close
			cache_kline.Low = util.MinFloat64(cache_kline.Low, kline.Low)
			cache_kline.High = util.MaxFloat64(cache_kline.High, kline.High)

			logx.Slowf("Cached Kline: %s", kline.String())
		}
	}

	return rst
}

func (s *SubData) ProcessKlineHistData(hist_kline *datastruct.RspHistKline) {

	logx.Slowf("SubData: Hist: %s", datastruct.HistKlineString(hist_kline.Klines))

	iter := hist_kline.Klines.Iterator()
	iter.Last()
	last_kline := iter.Value().(*datastruct.Kline)

	if _, ok := s.KlineInfo.Info[hist_kline.ReqInfo.Symbol]; !ok {
		s.KlineInfo.Info[hist_kline.ReqInfo.Symbol] = make(map[int]*KlineSubItem)
	}

	if _, ok := s.KlineInfo.Info[hist_kline.ReqInfo.Symbol][int(hist_kline.ReqInfo.Frequency)]; !ok {

		s.KlineInfo.Info[hist_kline.ReqInfo.Symbol][int(hist_kline.ReqInfo.Frequency)].cache_data = datastruct.NewKlineWithKline(last_kline)
		logx.Slowf("[Store] %s ", last_kline.String())
	} else if s.KlineInfo.Info[hist_kline.ReqInfo.Symbol][int(hist_kline.ReqInfo.Frequency)].cache_data == nil {

		s.KlineInfo.Info[hist_kline.ReqInfo.Symbol][int(hist_kline.ReqInfo.Frequency)].cache_data = datastruct.NewKlineWithKline(last_kline)
		logx.Slowf("[Store] %s ", last_kline.String())
	}

	if hist_kline.Klines.Size() > int(hist_kline.ReqInfo.Count) {
		logx.Infof("hist_kline.Klines.Size: %d, hist_kline.ReqInfo.Count: %d, last kline %+v, is not complete kline data, leave it in cache, wait for next round!",
			hist_kline.Klines.Size(), int(hist_kline.ReqInfo.Count), last_kline)
		hist_kline.Klines.Remove(last_kline.Time)
	}

}

func (s *SubData) SubSymbol(ws *net.WSInfo) {
	s.SymbolInfo.mutex.Lock()
	defer s.SymbolInfo.mutex.Unlock()

	if _, ok := s.SymbolInfo.ws_info.Get(ws.ID); !ok {
		s.SymbolInfo.ws_info.Put(ws.ID, ws)
	}

	logx.Infof("After Sub %s, %s", ws.String(), s.SymbolInfo.String())
}

func (s *SubData) UnSubSymbol(ws *net.WSInfo) {
	s.SymbolInfo.mutex.Lock()
	defer s.SymbolInfo.mutex.Unlock()

	if _, ok := s.SymbolInfo.ws_info.Get(ws.ID); !ok {
		return
	}

	s.SymbolInfo.ws_info.Remove(ws.ID)

	logx.Infof("SymbolInfo Remove %s : %+v", ws)
}

func (s *SubData) SubDepth(symbol string, ws *net.WSInfo) *datastruct.DepthQuote {
	s.DepthInfo.mutex.Lock()
	defer s.DepthInfo.mutex.Unlock()

	if _, ok := s.DepthInfo.Info[symbol]; !ok {
		s.DepthInfo.Info[symbol] = treemap.NewWith(utils.Int64Comparator)
	}

	s.DepthInfo.Info[symbol].Put(ws.ID, ws)

	logx.Infof("After Sub %s, %s, %s", symbol, ws.String(), s.DepthInfo.String())

	return nil
}

func (s *SubData) UnSubDepth(symbol string, ws *net.WSInfo) {
	s.DepthInfo.mutex.Lock()
	defer s.DepthInfo.mutex.Unlock()
	if _, ok := s.DepthInfo.Info[symbol]; !ok {
		return
	}

	s.DepthInfo.Info[symbol].Remove(ws.ID)

	logx.Infof("Depth Remove %s : %+v", symbol, ws)
}

func (s *SubData) SubTrade(symbol string, ws *net.WSInfo) {
	s.TradeInfo.mutex.Lock()
	defer s.TradeInfo.mutex.Unlock()

	if _, ok := s.TradeInfo.Info[symbol]; !ok {
		s.TradeInfo.Info[symbol] = treemap.NewWith(utils.Int64Comparator)
	}

	s.TradeInfo.Info[symbol].Put(ws.ID, ws)

	logx.Infof("After Sub %s, %s, %s", symbol, ws.String(), s.TradeInfo.String())
}

func (s *SubData) UnSubTrade(symbol string, ws *net.WSInfo) {
	s.TradeInfo.mutex.Lock()
	defer s.TradeInfo.mutex.Unlock()

	if _, ok := s.TradeInfo.Info[symbol]; !ok {
		return
	}

	s.TradeInfo.Info[symbol].Remove(ws.ID)
	logx.Infof("Trade Remove %s : %+v", symbol, ws)
}

func (s *SubData) SubKline(req_kline_info *datastruct.ReqHistKline, ws *net.WSInfo) {
	s.KlineInfo.mutex.Lock()
	defer s.KlineInfo.mutex.Unlock()

	if _, ok := s.KlineInfo.Info[req_kline_info.Symbol]; !ok {

		s.KlineInfo.Info[req_kline_info.Symbol] = make(map[int]*KlineSubItem)

		s.KlineInfo.Info[req_kline_info.Symbol][int(req_kline_info.Frequency)] = &KlineSubItem{
			ws_info:    treemap.NewWith(utils.Int64Comparator),
			cache_data: nil,
		}
	}

	s.KlineInfo.Info[req_kline_info.Symbol][int(req_kline_info.Frequency)].ws_info.Put(ws.ID, ws)

	logx.Infof("After Sub %s, %s,%s", req_kline_info.String(), ws.String(), s.KlineInfo.String())
}

func (s *SubData) UnSubKline(req_kline_info *datastruct.ReqHistKline, ws *net.WSInfo) {
	s.KlineInfo.mutex.Lock()
	defer s.KlineInfo.mutex.Unlock()

	if _, ok := s.KlineInfo.Info[req_kline_info.Symbol]; !ok {
		return
	} else if _, ok := s.KlineInfo.Info[req_kline_info.Symbol][int(req_kline_info.Frequency)]; !ok {
		return
	}

	s.KlineInfo.Info[req_kline_info.Symbol][int(req_kline_info.Frequency)].ws_info.Remove(ws.ID)

	logx.Infof("KLine Remove %s, %d : %+v", req_kline_info.Symbol, int(req_kline_info.Frequency), ws)
}
