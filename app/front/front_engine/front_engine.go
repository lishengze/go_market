package front_engine

import (
	"market_server/common/datastruct"
)

type FrontEngine struct {
	sub_data *SubData
}

func (f *FrontEngine) PublishDepth(*datastruct.DepthQuote) {

}

func (f *FrontEngine) PublishTrade(*datastruct.Trade) {

}

func (f *FrontEngine) PublishKline(*datastruct.Kline) {

}

func (f *FrontEngine) PublishChangeinfo(*datastruct.ChangeInfo) {

}

func (f *FrontEngine) PublishHistKline(klines *datastruct.HistKline) {
	// d.publish_kline(kline)
	f.sub_data.ProcessKlineHistData(klines)

}

func (f *FrontEngine) SubTrade(symbol string) *datastruct.Trade {
	return nil
}

func (f *FrontEngine) UnSubTrade(symbol string) {

}

func (f *FrontEngine) SubDepth(symbol string) *datastruct.DepthQuote {
	return nil
}

func (f *FrontEngine) UnSubDepth(symbol string) {

}

func (f *FrontEngine) SubKline(req_kline_info *datastruct.ReqHistKline) *datastruct.HistKline {
	return nil
}

func (f *FrontEngine) UnSubKline(req_kline_info *datastruct.ReqHistKline) {

}
