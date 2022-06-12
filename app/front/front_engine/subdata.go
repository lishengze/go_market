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

func (s *SubData) GetKlinePubInfoList(kline *datastruct.Kline) []*KlinePubInfo {
	var rst []*KlinePubInfo

	// if symbol_sub_info, ok := s.KlineInfo.Info[kline.Symbol]; ok {
	// 	for frequency, sub_info := range symbol_sub_info {

	// 	}
	// }
	return rst
}

func (s *SubData) ProcessKlineHistData(hist_kline *datastruct.RspHistKline) {
	if _, ok := s.KlineInfo.Info[hist_kline.ReqInfo.Symbol]; !ok {
		s.KlineInfo.Info[hist_kline.ReqInfo.Symbol] = make(map[int]*KlineSubItem)
	}

	if _, ok := s.KlineInfo.Info[hist_kline.ReqInfo.Symbol][int(hist_kline.ReqInfo.Frequency)]; !ok {

	}

}
