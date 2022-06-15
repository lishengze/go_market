package front_engine

import (
	"market_server/app/front/net"
	"market_server/common/datastruct"
	"sync"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
)

type DepthPubInfo struct {
	ws_info *net.WSInfo
	data    *datastruct.DepthQuote
}

type TradePubInfo struct {
	ws_info *net.WSInfo
	data    *datastruct.Trade
}

type KlinePubInfo struct {
	ws_info *net.WSInfo
	data    *datastruct.Kline
}

type DepthSubInfo struct {
	mutex sync.Mutex
	Info  map[string]*treemap.Map
}

type KlineSubItem struct {
	ws_info    *treemap.Map
	cache_data *datastruct.Kline
}

type KlineSubInfo struct {
	mutex sync.Mutex
	Info  map[string](map[int]*KlineSubItem)
}

type TradeSubInfo struct {
	mutex sync.Mutex
	Info  map[string]*treemap.Map
}

type SubData struct {
	DepthInfo *DepthSubInfo
	TradeInfo *TradeSubInfo
	KlineInfo *KlineSubInfo
}

func NewSubData() *SubData {
	return nil
}

func (s *SubData) GetDepthPubInfoList(depth *datastruct.DepthQuote) []*DepthPubInfo {
	var rst []*DepthPubInfo

	s.DepthInfo.mutex.Lock()
	defer s.DepthInfo.mutex.Unlock()

	if sub_tree, ok := s.DepthInfo.Info[depth.Symbol]; ok {
		sub_tree_iter := sub_tree.Iterator()
		sub_tree_iter.Begin()
		for sub_tree_iter.Next() {
			rst = append(rst, &DepthPubInfo{
				ws_info: sub_tree_iter.Value().(*net.WSInfo),
				data:    depth,
			})

		}
	}

	return rst
}

func (s *SubData) GetTradePubInfoList(trade *datastruct.Trade) []*TradePubInfo {
	var rst []*TradePubInfo

	s.TradeInfo.mutex.Lock()
	defer s.TradeInfo.mutex.Unlock()

	if sub_tree, ok := s.TradeInfo.Info[trade.Symbol]; ok {
		sub_tree_iter := sub_tree.Iterator()
		sub_tree_iter.Begin()
		for sub_tree_iter.Next() {
			rst = append(rst, &TradePubInfo{
				ws_info: sub_tree_iter.Value().(*net.WSInfo),
				data:    trade,
			})
		}
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

func (s *SubData) GetKlinePubInfoList(kline *datastruct.Kline) []*KlinePubInfo {
	var rst []*KlinePubInfo

	s.KlineInfo.mutex.Lock()
	defer s.KlineInfo.mutex.Unlock()

	if _, ok := s.KlineInfo.Info[kline.Symbol]; !ok {
		return rst
	}

	is_updated := false

	for resolution, sub_info := range s.KlineInfo.Info[kline.Symbol] {
		cache_kline := sub_info.cache_data

		if kline.Time-cache_kline.Time >= int64(resolution) {
			s.KlineInfo.Info[kline.Symbol][int(kline.Resolution)].cache_data = kline
			is_updated = true
		}
	}

	if is_updated {
		sub_tree := s.KlineInfo.Info[kline.Symbol][int(kline.Resolution)].ws_info
		sub_tree_iter := sub_tree.Iterator()
		sub_tree_iter.Begin()
		for sub_tree_iter.Next() {
			rst = append(rst, &KlinePubInfo{
				ws_info: sub_tree_iter.Value().(*net.WSInfo),
				data:    kline,
			})
		}
	}

	return rst
}

func (s *SubData) ProcessKlineHistData(hist_kline *datastruct.RspHistKline) {

	iter := hist_kline.Klines.Iterator()
	iter.Last()
	last_kline := iter.Value().(*datastruct.Kline)

	if _, ok := s.KlineInfo.Info[hist_kline.ReqInfo.Symbol]; !ok {
		s.KlineInfo.Info[hist_kline.ReqInfo.Symbol] = make(map[int]*KlineSubItem)
	}

	if _, ok := s.KlineInfo.Info[hist_kline.ReqInfo.Symbol][int(hist_kline.ReqInfo.Frequency)]; !ok {

		s.KlineInfo.Info[hist_kline.ReqInfo.Symbol][int(hist_kline.ReqInfo.Frequency)].cache_data = last_kline
	}

}

func (s *SubData) SubTrade(symbol string, ws *net.WSInfo) {
	s.TradeInfo.mutex.Lock()
	defer s.TradeInfo.mutex.Unlock()

	if _, ok := s.TradeInfo.Info[symbol]; !ok {
		s.TradeInfo.Info[symbol] = treemap.NewWith(utils.Int64Comparator)
	}

	s.TradeInfo.Info[symbol].Put(ws.ID, ws)
}

func (s *SubData) UnSubTrade(symbol string, ws *net.WSInfo) {
	s.TradeInfo.mutex.Lock()
	defer s.TradeInfo.mutex.Unlock()

	if _, ok := s.TradeInfo.Info[symbol]; !ok {
		return
	}

	s.TradeInfo.Info[symbol].Remove(ws.ID)
}

func (s *SubData) SubDepth(symbol string, ws *net.WSInfo) *datastruct.DepthQuote {
	s.DepthInfo.mutex.Lock()
	defer s.DepthInfo.mutex.Unlock()

	if _, ok := s.DepthInfo.Info[symbol]; !ok {
		s.DepthInfo.Info[symbol] = treemap.NewWith(utils.Int64Comparator)
	}

	s.DepthInfo.Info[symbol].Put(ws.ID, ws)

	return nil
}

func (s *SubData) UnSubDepth(symbol string, ws *net.WSInfo) {
	s.DepthInfo.mutex.Lock()
	defer s.DepthInfo.mutex.Unlock()
	if _, ok := s.DepthInfo.Info[symbol]; !ok {
		return
	}

	s.DepthInfo.Info[symbol].Remove(ws.ID)
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
}
