package front_engine

import (
	"market_server/common/datastruct"
)

type FrontEngine struct {
}

func (d *FrontEngine) publish_depth(*datastruct.DepthQuote) {

}

func (d *FrontEngine) publish_trade(*datastruct.Trade) {

}

func (d *FrontEngine) publish_kline(*datastruct.Kline) {

}

func (d *FrontEngine) publish_changeinfo(*datastruct.ChangeInfo) {

}

func (d *FrontEngine) SubTrade(symbol string) *datastruct.Trade {
	return nil
}

func (d *FrontEngine) UnSubTrade(symbol string) {

}

func (d *FrontEngine) SubDepth(symbol string) *datastruct.DepthQuote {
	return nil
}

func (d *FrontEngine) UnSubDepth(symbol string) {

}

func (d *FrontEngine) SubKline(req_kline_info *datastruct.ReqHistKline) *datastruct.HistKline {
	return nil
}

func (d *FrontEngine) UnSubKline(req_kline_info *datastruct.ReqHistKline) {

}
