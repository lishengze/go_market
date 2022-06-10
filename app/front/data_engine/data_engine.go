package data_engine

import (
	"market_server/app/front/config"
	"market_server/common/datastruct"

	"github.com/zeromicro/go-zero/core/logx"
)

/*
接受数据的接口;
*/

type PeriodData struct {
	cache_data []*datastruct.Kline

	Max     float64
	MaxTime int64

	Min     float64
	MinTime int64

	Start  float64
	Change float64

	TimeMinutes int
	Count       int
}

func (o *PeriodData) UpdateWithTrade(trade datastruct.Trade) {

}

func (o *PeriodData) UpdateWithKline(kline datastruct.Kline) {

}

type DataEngine struct {
	RecvDataChan *datastruct.DataChannel
	cache_config *config.CacheConfig

	cache_data map[string]*PeriodData // 缓存24小时1分频率的 k 线数据，用来计算24小时的涨跌幅;
}

func NewDataEngine(recvDataChan *datastruct.DataChannel, config *config.CacheConfig) *DataEngine {
	rst := &DataEngine{
		RecvDataChan: recvDataChan,
		cache_config: config,
		cache_data:   make(map[string]*PeriodData),
	}

	return rst
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
