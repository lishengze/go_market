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

func (s *SubData) GetKlinePubInfoListWithTrade(trade *datastruct.Trade) []*KlinePubInfo {
	defer util.CatchExp("GetKlinePubInfoListWithTrade")

	// logx.Slowf("0")
	if trade == nil {
		logx.Errorf("trade is nil")
		return nil
	}
	// logx.Slowf("1")

	var rst []*KlinePubInfo

	s.KlineInfo.mutex.Lock()
	defer s.KlineInfo.mutex.Unlock()

	// logx.Slowf("2")

	if _, ok := s.KlineInfo.Info[trade.Symbol]; !ok {
		return rst
	}

	// logx.Slowf("3")

	for resolution, sub_info := range s.KlineInfo.Info[trade.Symbol] {
		// logx.Slowf("4")
		if sub_info.cache_data == nil {
			logx.Errorf("Hisk Kline %s, %d , cache_data empty", trade.Symbol, resolution)
			continue
		}

		var pub_kline *datastruct.Kline

		if datastruct.IsNewKlineStartTime(trade.Time, int64(resolution)) {

			tmp_kline := &datastruct.Kline{
				Exchange:   trade.Exchange,
				Symbol:     trade.Symbol,
				Time:       trade.Time,
				Open:       trade.Price,
				High:       trade.Price,
				Low:        trade.Price,
				Close:      trade.Price,
				Volume:     trade.Volume,
				Resolution: resolution,
			}

			logx.Slowf("New Kline With: \nTrade %s;\nkline: %s \n", trade.String(), tmp_kline.FullString())

			s.KlineInfo.Info[trade.Symbol][resolution].cache_data = datastruct.NewKlineWithKline(tmp_kline)
			s.KlineInfo.Info[trade.Symbol][resolution].cache_data.Resolution = resolution

			pub_kline = datastruct.NewKlineWithKline(tmp_kline)
		} else {
			cache_kline := sub_info.cache_data

			NextKlineTime := cache_kline.Time + int64(resolution)*datastruct.NANO_PER_SECS

			if resolution == 60 {
				NextKlineTime = NextKlineTime + int64(resolution)*datastruct.NANO_PER_SECS
			}

			if trade.Time <= cache_kline.Time {
				logx.Errorf("Trade.Time %s, earlier than CachedKlineTime: %s", util.TimeStrFromInt(trade.Time), util.TimeStrFromInt(cache_kline.Time))
				continue
			}

			if trade.Time > NextKlineTime {
				logx.Errorf("Trade.Time %s, later than NextKlineTime: %s", util.TimeStrFromInt(trade.Time), util.TimeStrFromInt(NextKlineTime))
				continue
			}

			cache_kline.Close = trade.Price
			if cache_kline.Low > trade.Price {
				cache_kline.Low = trade.Price

			}
			if cache_kline.High < trade.Price {
				cache_kline.High = trade.Price
			}

			cache_kline.Volume = cache_kline.Volume + trade.Volume

			pub_kline = datastruct.NewKlineWithKline(cache_kline)
		}

		if pub_kline != nil {
			cur_pub_list := s.GetKlinePubInfoListAtom(sub_info, pub_kline)
			rst = append(rst, cur_pub_list...)
		}
	}

	return rst
}

func (s *SubData) GetKlinePubInfoListAtom(sub_info *KlineSubItem, pub_kline *datastruct.Kline) []*KlinePubInfo {
	defer util.CatchExp("GetKlinePubInfoListAtom")
	var rst []*KlinePubInfo

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

	return rst
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
	defer util.CatchExp("GetKlinePubInfoList")

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
			logx.Errorf("NewKLineTime: %d, CachedKlineTime: %d", kline.Time, cache_kline.Time)
			continue
		}

		var pub_kline *datastruct.Kline
		pub_kline = nil

		if datastruct.IsOldKlineEnd(kline, int64(resolution)) {
			logx.Slowf("Old Kline End: rsl:%d, %s", resolution, kline.String())

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
			// byte_data := NewKlineUpdateJsonMsg(pub_kline)
			// sub_tree := sub_info.ws_info
			// sub_tree_iter := sub_tree.Iterator()

			// for sub_tree_iter.Begin(); sub_tree_iter.Next(); {
			// 	rst = append(rst, &KlinePubInfo{
			// 		ws_info:    sub_tree_iter.Value().(*net.WSInfo),
			// 		data:       byte_data,
			// 		Symbol:     pub_kline.Symbol,
			// 		Resolution: pub_kline.Resolution,
			// 	})
			// }

			s.KlineInfo.Info[kline.Symbol][resolution].cache_data = datastruct.NewKlineWithKline(pub_kline)

		} else if datastruct.IsNewKlineStart(kline, int64(resolution)) {
			logx.Slowf("New Kline Start: rsl:%d, %s", resolution, kline.String())

			s.KlineInfo.Info[kline.Symbol][resolution].cache_data = datastruct.NewKlineWithKline(kline)
			s.KlineInfo.Info[kline.Symbol][resolution].cache_data.Resolution = resolution

			pub_kline = datastruct.NewKlineWithKline(kline)

		} else {
			cache_kline.Close = kline.Close
			cache_kline.Low = util.MinFloat64(cache_kline.Low, kline.Low)
			cache_kline.High = util.MaxFloat64(cache_kline.High, kline.High)
			cache_kline.Volume = cache_kline.Volume + kline.Volume

			logx.Slowf("Cached Kline:%d, %s", resolution, kline.String())

			pub_kline = datastruct.NewKlineWithKline(cache_kline)
		}

		if pub_kline != nil {
			cur_pub_list := s.GetKlinePubInfoListAtom(sub_info, pub_kline)
			rst = append(rst, cur_pub_list...)
		}
	}

	return rst
}

func (s *SubData) ProcessKlineHistData(hist_kline *datastruct.RspHistKline) {
	defer util.CatchExp("GetKlinePubInfoList")

	logx.Slowf("[SD] HistK, rsl: %d, %s", hist_kline.ReqInfo.Frequency, datastruct.HistKlineSimpleTime(hist_kline.Klines))

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

	if hist_kline.Klines.Size() == int(hist_kline.ReqInfo.Count)+1 {
		iter.First()
		first_kline := iter.Value().(*datastruct.Kline)
		hist_kline.Klines.Remove(first_kline.Time)
	}

	// if !hist_kline.IsLastComplete {
	// logx.Infof("hist_kline.Klines.Size: %d, hist_kline.ReqInfo.Count: %d, last kline %+v, is not complete kline data, leave it in cache, wait for next round!",
	// 	hist_kline.Klines.Size(), int(hist_kline.ReqInfo.Count), last_kline)
	// hist_kline.Klines.Remove(last_kline.Time)
	// }

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
	logx.Infof("UnSubTrade Remove %s : %+v", symbol, ws)
}

func (s *SubData) SubKline(req_kline_info *datastruct.ReqHistKline, ws *net.WSInfo) {
	defer util.CatchExp("SubKline")

	s.KlineInfo.mutex.Lock()
	defer s.KlineInfo.mutex.Unlock()

	if _, ok := s.KlineInfo.Info[req_kline_info.Symbol]; !ok {

		s.KlineInfo.Info[req_kline_info.Symbol] = make(map[int]*KlineSubItem)

		s.KlineInfo.Info[req_kline_info.Symbol][int(req_kline_info.Frequency)] = &KlineSubItem{
			ws_info:    treemap.NewWith(utils.Int64Comparator),
			cache_data: nil,
		}
	}

	if _, ok := s.KlineInfo.Info[req_kline_info.Symbol][int(req_kline_info.Frequency)]; !ok {
		s.KlineInfo.Info[req_kline_info.Symbol][int(req_kline_info.Frequency)] = &KlineSubItem{
			ws_info:    treemap.NewWith(utils.Int64Comparator),
			cache_data: nil,
		}
	}

	s.KlineInfo.Info[req_kline_info.Symbol][int(req_kline_info.Frequency)].ws_info.Put(ws.ID, ws)

	logx.Infof("SubKline After Sub %s, %s,%s", req_kline_info.String(), ws.String(), s.KlineInfo.String())
}

func (s *SubData) UnSubKline(req_kline_info *datastruct.ReqHistKline, ws *net.WSInfo) {
	defer util.CatchExp("UnSubKline")

	s.KlineInfo.mutex.Lock()
	defer s.KlineInfo.mutex.Unlock()

	if _, ok := s.KlineInfo.Info[req_kline_info.Symbol]; !ok {
		return
	} else if _, ok := s.KlineInfo.Info[req_kline_info.Symbol][int(req_kline_info.Frequency)]; !ok {
		return
	}

	if s.KlineInfo.Info[req_kline_info.Symbol][int(req_kline_info.Frequency)].ws_info != nil {
		s.KlineInfo.Info[req_kline_info.Symbol][int(req_kline_info.Frequency)].ws_info.Remove(ws.ID)
		logx.Infof("UnSubKline Remove %s, %d : %+v", req_kline_info.Symbol, int(req_kline_info.Frequency), ws)
	} else {
		logx.Infof("KlineInfo, %s.%d is already nil!", req_kline_info.Symbol, req_kline_info.Frequency)
	}

}
