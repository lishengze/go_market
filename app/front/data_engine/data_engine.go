package data_engine

import (
	"context"
	"fmt"
	"market_server/app/dataManager/rpc/marketservice"
	"market_server/app/front/config"
	"market_server/app/front/worker"
	"market_server/common/datastruct"
	"market_server/common/util"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

type DataEngine struct {
	RecvDataChan *datastruct.DataChannel
	config       *config.Config

	cache_kline_data map[string]*PeriodData // 缓存24小时1分频率的 k 线数据，用来计算24小时的涨跌幅;

	msclient marketservice.MarketService

	next_worker worker.WorkerI
}

func NewDataEngine(recvDataChan *datastruct.DataChannel, config *config.Config) *DataEngine {

	rst := &DataEngine{
		RecvDataChan:     recvDataChan,
		config:           config,
		cache_kline_data: make(map[string]*PeriodData),
		msclient:         marketservice.NewMarketService(zrpc.MustNewClient(config.RpcConfig)),
	}

	return rst
}

func (d *DataEngine) UpdateMeta(symbols []string) {
	for _, symbol := range symbols {
		if _, ok := d.cache_kline_data[symbol]; !ok {
			d.InitPeriodDara(symbol)
		}
	}
}

func (a *DataEngine) InitPeriodDara(symbol string) {
	a.cache_kline_data[symbol] = &PeriodData{
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

	a.cache_kline_data[symbol].UpdateWithPbKlines(hist_klines)
}

func (a *DataEngine) StartListenRecvdata() {
	logx.Info("[S] DBServer start_listen_recvdata")
	go func() {
		for {
			select {
			case new_depth := <-a.RecvDataChan.DepthChannel:
				go a.process_depth(new_depth)
			case new_kline := <-a.RecvDataChan.KlineChannel:
				go a.process_kline(new_kline)
			case new_trade := <-a.RecvDataChan.TradeChannel:
				go a.process_trade(new_trade)
			}
		}
	}()
	logx.Info("[S] DBServer start_receiver Over!")
}

func (d *DataEngine) process_depth(depth *datastruct.DepthQuote) error {
	d.publish_depth(depth)
	return nil
}

func (d *DataEngine) process_kline(kline *datastruct.Kline) error {

	if _, ok := d.cache_kline_data[kline.Symbol]; !ok {
		d.InitPeriodDara(kline.Symbol)
	}

	d.cache_kline_data[kline.Symbol].UpdateWithKline(kline)

	d.publish_kline(kline)

	d.publish_changeinfo(d.cache_kline_data[kline.Symbol].GetChangeInfo())

	return nil
}

func (d *DataEngine) process_trade(trade *datastruct.Trade) error {

	d.cache_kline_data[trade.Symbol].UpdateWithTrade(trade)

	d.publish_changeinfo(d.cache_kline_data[trade.Symbol].GetChangeInfo())

	d.publish_trade(trade)
	return nil
}

func (d *DataEngine) publish_depth(depth *datastruct.DepthQuote) {
	d.next_worker.publish_depth(depth)
}

func (d *DataEngine) publish_trade(trade *datastruct.Trade) {
	d.next_worker.publish_trade(trade)
}

func (d *DataEngine) publish_kline(kline *datastruct.Kline) {
	d.next_worker.publish_kline(kline)
}

func (d *DataEngine) publish_hist_kline(kline []*datastruct.Kline) {
	// d.publish_kline(kline)
}

func (d *DataEngine) publish_changeinfo(change_info *datastruct.ChangeInfo) {
	d.publish_changeinfo(change_info)
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
	req_hist_info := &marketservice.ReqHishKlineInfo{
		Symbol:    req_kline_info.symbol,
		Exchange:  req_kline_info.exchange,
		StartTime: req_kline_info.start_time,
		EndTime:   req_kline_info.end_time,
		Count:     req_kline_info.count,
		Frequency: req_kline_info.frequency,
	}

	hist_klines, err := d.msclient.RequestHistKlineData(context.Background(), req_hist_info)

	if err != nil {
		return nil
	}

	for _, pb_kline := range hist_klines.KlineData {
		kline := marketservice.NewKlineWithPbKline(pb_kline)
		if kline == nil {
			continue
		}
	}
}

func (d *DataEngine) UnSubKline(req_kline_info *datastruct.ReqHistKline) {

}
