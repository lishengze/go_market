package front_engine

import (
	"market_server/common/datastruct"
)

type WSInfo struct {
}

type DepthPubInfo struct {
	ws_info *WSInfo
	data    *datastruct.DepthQuote
}

type TradePubInfo struct {
	ws_info *WSInfo
	data    *datastruct.Trade
}

type KlinePubInfo struct {
	ws_info *WSInfo
	data    *datastruct.Kline
}

type DepthSubInfo struct {
	Info map[string][]*WSInfo
}

type KlineSubItem struct {
	ws_info    []*WSInfo
	cache_data *datastruct.Kline
}

type KlineSubInfo struct {
	Info map[string](map[int]*KlineSubItem)
}

type TradeSubInfo struct {
	Info map[string][]*WSInfo
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

	if sub_list, ok := s.DepthInfo.Info[depth.Symbol]; ok {
		for _, sub_info := range sub_list {
			rst = append(rst, &DepthPubInfo{
				ws_info: sub_info,
				data:    depth,
			})
		}
	}

	return rst
}

func (s *SubData) GetTradePubInfoList(trade *datastruct.Trade) []*TradePubInfo {
	var rst []*TradePubInfo

	if sub_list, ok := s.TradeInfo.Info[trade.Symbol]; ok {
		for _, sub_info := range sub_list {
			rst = append(rst, &TradePubInfo{
				ws_info: sub_info,
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

		if kline.Time-cache_kline.Time >= int64(kline.Resolution*1000000000) {
			s.KlineInfo.Info[kline.Symbol][int(kline.Resolution)].cache_data = kline
		}
	}
}

func (s *SubData) GetKlinePubInfoList(kline *datastruct.Kline) []*KlinePubInfo {
	var rst []*KlinePubInfo

	if _, ok := s.KlineInfo.Info[kline.Symbol]; !ok {
		return rst
		// s.KlineInfo.Info[kline.Symbol] = make(map[int]*KlineSubItem)
	}

	is_updated := false

	for resolution, sub_info := range s.KlineInfo.Info[kline.Symbol] {
		cache_kline := sub_info.cache_data

		if kline.Time-cache_kline.Time >= int64(resolution*1000000000) {
			s.KlineInfo.Info[kline.Symbol][int(kline.Resolution)].cache_data = kline
			is_updated = true
		}
	}

	if is_updated {
		ws_info_list := s.KlineInfo.Info[kline.Symbol][int(kline.Resolution)].ws_info
		for _, ws_info := range ws_info_list {
			rst = append(rst, &KlinePubInfo{
				ws_info: ws_info,
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
	} else {
		cache_kline := s.KlineInfo.Info[hist_kline.ReqInfo.Symbol][int(hist_kline.ReqInfo.Frequency)].cache_data

		if last_kline.Time-cache_kline.Time >= int64(hist_kline.ReqInfo.Frequency*1000000000) {
			s.KlineInfo.Info[hist_kline.ReqInfo.Symbol][int(hist_kline.ReqInfo.Frequency)].cache_data = last_kline
		}
	}

}
