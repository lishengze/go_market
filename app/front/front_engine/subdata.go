package front_engine

import (
	"market_server/app/front/net"
	"market_server/common/datastruct"
	"market_server/common/util"
	"sync"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

type DepthPubInfo struct {
	ws_info *net.WSInfo
	data    *datastruct.DepthQuote
}

type SymbolPubInfo struct {
	ws_info *net.WSInfo
	data    []string
}

type TradePubInfo struct {
	ws_info *net.WSInfo
	data    *datastruct.Trade
}

type KlinePubInfo struct {
	ws_info *net.WSInfo
	data    *datastruct.Kline
}

type SymbolSubInfo struct {
	mutex   sync.Mutex
	ws_info *treemap.Map
}

func NewSymbolSubInfo() *SymbolSubInfo {
	return &SymbolSubInfo{
		ws_info: treemap.NewWith(utils.Int64Comparator),
	}
}

type DepthSubInfo struct {
	mutex sync.Mutex
	Info  map[string]*treemap.Map
}

func NewDepthSubInfo() *DepthSubInfo {
	return &DepthSubInfo{
		Info: make(map[string]*treemap.Map),
	}
}

type KlineSubItem struct {
	ws_info    *treemap.Map
	cache_data *datastruct.Kline
}

func NewKlineWithKline() *KlineSubItem {
	return &KlineSubItem{
		ws_info: treemap.NewWith(utils.Int64Comparator),
	}
}

type KlineSubInfo struct {
	mutex sync.Mutex
	Info  map[string](map[int]*KlineSubItem)
}

func NewKlineSubInfo() *KlineSubInfo {
	return &KlineSubInfo{
		Info: make(map[string](map[int]*KlineSubItem)),
	}
}

type TradeSubInfo struct {
	mutex sync.Mutex
	Info  map[string]*treemap.Map
}

func NewTradeSubInfo() *TradeSubInfo {
	return &TradeSubInfo{
		Info: make(map[string]*treemap.Map),
	}
}

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

	for iter.Begin(); iter.Next(); {
		rst = append(rst, &SymbolPubInfo{
			ws_info: iter.Value().(*net.WSInfo),
			data:    symbollist,
		})
	}

	return rst
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
		cache_kline := sub_info.cache_data

		if kline.Time <= cache_kline.Time {
			logx.Errorf("NewKLineTime: %d, CachedKlineTime: %d")
			continue
		}

		if datastruct.IsOldKlineEnd(kline, int64(resolution)) {
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

			sub_tree := s.KlineInfo.Info[kline.Symbol][int(kline.Resolution)].ws_info
			sub_tree_iter := sub_tree.Iterator()
			for sub_tree_iter.Begin(); sub_tree_iter.Next(); {
				rst = append(rst, &KlinePubInfo{
					ws_info: sub_tree_iter.Value().(*net.WSInfo),
					data:    pub_kline,
				})
			}
			s.KlineInfo.Info[kline.Symbol][resolution].cache_data = datastruct.NewKlineWithKline(pub_kline)

		} else if datastruct.IsNewKlineStart(kline, int64(resolution)) {
			s.KlineInfo.Info[kline.Symbol][int(kline.Resolution)].cache_data = datastruct.NewKlineWithKline(kline)
			s.KlineInfo.Info[kline.Symbol][int(kline.Resolution)].cache_data.Resolution = resolution
		} else {
			cache_kline.Close = kline.Close
			cache_kline.Low = util.MinFloat64(cache_kline.Low, kline.Low)
			cache_kline.High = util.MaxFloat64(cache_kline.High, kline.High)
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

func (s *SubData) SubSymbol(ws *net.WSInfo) {
	s.SymbolInfo.mutex.Lock()
	defer s.SymbolInfo.mutex.Unlock()

	if _, ok := s.SymbolInfo.ws_info.Get(ws.ID); !ok {
		s.SymbolInfo.ws_info.Put(ws.ID, ws)
	}
}

func (s *SubData) UnSubSymbol(ws *net.WSInfo) {
	s.SymbolInfo.mutex.Lock()
	defer s.SymbolInfo.mutex.Unlock()

	if _, ok := s.SymbolInfo.ws_info.Get(ws.ID); !ok {
		return
	}

	s.SymbolInfo.ws_info.Remove(ws.ID)
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
