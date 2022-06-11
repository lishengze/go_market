package data_engine

import (
	"context"
	"fmt"
	"market_server/app/dataManager/rpc/marketservice"
	"market_server/app/front/config"
	"market_server/common/datastruct"
	"market_server/common/util"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

/*
接受数据的接口;
*/

type AtomData struct {
	price float64
	time  int64
}

type SortedList struct {
	list *[]AtomData
}

type PeriodData struct {
	time_cache_data  *treemap.Map
	price_cache_data *treemap.Map

	Max     float64
	MaxTime int64

	Min     float64
	MinTime int64

	Start     float64
	StartTime int64

	Change float64

	TimeNanos int64
	Count     int

	CurTrade *datastruct.Trade
}

func (p *PeriodData) UpdateWithTrade(trade *datastruct.Trade) {
	p.CurTrade = trade

	p.UpdateMeta()
}

func (p *PeriodData) UpdateWithKline(kline *datastruct.Kline) {
	p.time_cache_data.Put(kline.Time, kline)

	p.EraseOuttimeData()

	p.UpdateMeta()
}

func (p *PeriodData) EraseOuttimeData() {
	type outtime_data struct {
		time  int64
		price float64
	}

	outtime_datalist := []*outtime_data{}

	begin_iter := p.time_cache_data.Iterator()
	if ok := begin_iter.First(); !ok {
		return
	}

	last_iter := p.time_cache_data.Iterator()
	if ok := last_iter.Last(); !ok {
		return
	}

	for begin_iter.Next() {
		if last_iter.Key().(int64)-begin_iter.Key().(int64) > p.TimeNanos {

			outtime_datalist = append(outtime_datalist, &outtime_data{
				time:  begin_iter.Key().(int64),
				price: begin_iter.Value().(*datastruct.Kline).High,
			})
		} else {
			break
		}
	}

	for _, outtime := range outtime_datalist {
		p.time_cache_data.Remove(outtime.time)
		p.price_cache_data.Remove(outtime.price)
	}

}

func (p *PeriodData) InitCacheData(klines *marketservice.HistKlineData) {
	for _, pb_kline := range klines.KlineData {
		kline := marketservice.NewKlineWithPbKline(pb_kline)
		if kline == nil {
			continue
		}

		// if p.MaxTime == 0 || p.Max < kline.High {
		// 	p.Max = kline.High
		// 	p.MaxTime = kline.Time
		// }

		// if p.MinTime == 0 || p.Min > kline.Low {
		// 	p.Min = kline.Low
		// 	p.MinTime = kline.Time
		// }

		// if p.StartTime == 0 || p.StartTime > kline.Time {
		// 	p.StartTime = kline.Time
		// 	p.Start = kline.Open
		// }

		p.time_cache_data.Put(kline.Time, kline)
		p.price_cache_data.Put(kline.High, kline)
	}
}

func (p *PeriodData) UpdateMeta() {

}

func (p *PeriodData) UpdateWithPbKlines(klines *marketservice.HistKlineData) {
	p.InitCacheData(klines)

	p.EraseOuttimeData()

	p.UpdateMeta()
}

type DataEngine struct {
	RecvDataChan *datastruct.DataChannel
	config       *config.Config

	cache_data map[string]*PeriodData // 缓存24小时1分频率的 k 线数据，用来计算24小时的涨跌幅;

	msclient marketservice.MarketService
}

func NewDataEngine(recvDataChan *datastruct.DataChannel, config *config.Config) *DataEngine {

	rst := &DataEngine{
		RecvDataChan: recvDataChan,
		config:       config,
		cache_data:   make(map[string]*PeriodData),
		msclient:     marketservice.NewMarketService(zrpc.MustNewClient(config.RpcConfig)),
	}

	return rst
}

func (a *DataEngine) InitPeriodDara(symbol string) {
	a.cache_data[symbol] = &PeriodData{
		TimeNanos:       24 * 60 * 60 * 1000000000,
		Count:           0,
		MaxTime:         0,
		MinTime:         0,
		time_cache_data: treemap.NewWith(utils.Int64Comparator),
	}

	end_time_nanos := uint64(util.TimeMinuteNanos())
	start_time_nanos := uint64(end_time_nanos - 24*60*1000000000)

	req_hist_info := &marketservice.ReqHishKlineInfo{
		Symbol:    symbol,
		Exchange:  datastruct.BCTS_EXCHANGE,
		StartTime: start_time_nanos,
		EndTime:   end_time_nanos,
		Count:     24 * 60,
		Frequency: 60,
	}

	hist_klines, err := a.msclient.RequestHistKlineData(context.Background(), req_hist_info)

	if err != nil {
		fmt.Printf("err %+v \n", err)
		logx.Errorf("ReqHistKline Err: %+v\n", err)
	}

	fmt.Printf("Rst: %+v \n", hist_klines)

	a.cache_data[symbol].UpdateWithPbKlines(hist_klines)
}

func (a *DataEngine) StartListenRecvdata() {
	logx.Info("[S] DBServer start_listen_recvdata")
	go func() {
		for {
			select {
			case new_depth := <-a.RecvDataChan.DepthChannel:
				a.process_depth(new_depth)
			case new_kline := <-a.RecvDataChan.KlineChannel:
				a.process_kline(new_kline)
			case new_trade := <-a.RecvDataChan.TradeChannel:
				a.process_trade(new_trade)
			}
		}
	}()
	logx.Info("[S] DBServer start_receiver Over!")
}

func (d *DataEngine) process_depth(depth *datastruct.DepthQuote) error {
	return nil
}

func (d *DataEngine) process_kline(depth *datastruct.Kline) error {
	return nil
}

func (d *DataEngine) process_trade(depth *datastruct.Trade) error {
	return nil
}

func (d *DataEngine) publish_depth(*datastruct.DepthQuote) {

}

func (d *DataEngine) publish_trade(*datastruct.Trade) {

}

func (d *DataEngine) publish_kline(*datastruct.Kline) {

}

func (d *DataEngine) SubTrade(symbol string) *datastruct.Trade {
	return nil
}

func (d *DataEngine) UnSubTrade(symbol string) {

}

func (d *DataEngine) SubDepth(symbol string) *datastruct.DepthQuote {
	return nil
}

func (d *DataEngine) UnSubDepth(symbol string) {

}

func (d *DataEngine) SubKline(req_kline_info *datastruct.ReqHistKline) *datastruct.HistKline {
	return nil
}

func (d *DataEngine) UnSubKline(req_kline_info *datastruct.ReqHistKline) {

}
