package data_engine

import (
	"market_server/common/datastruct"
	"market_server/common/monitorStruct"
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
)

type KafkaMonitor struct {
	RecvDataChan *datastruct.DataChannel

	depth_cache_map       map[string](*monitorStruct.MonitorAtom)
	depth_cache_map_mutex sync.Mutex

	trade_cache_map       map[string](*monitorStruct.MonitorAtom)
	trade_cache_map_mutex sync.Mutex

	kline_cache_map       map[string](*monitorStruct.MonitorAtom)
	kline_cache_map_mutex sync.Mutex
}

func NewKafkaMonitor(recvDataChan *datastruct.DataChannel) *KafkaMonitor {
	return &KafkaMonitor{
		RecvDataChan:    recvDataChan,
		depth_cache_map: make(map[string]*monitorStruct.MonitorAtom),
		trade_cache_map: make(map[string]*monitorStruct.MonitorAtom),
		kline_cache_map: make(map[string]*monitorStruct.MonitorAtom),
	}
}

func (d *KafkaMonitor) Start() {
	logx.Infof("KafkaMonitor Start!")

	d.StartListenRecvdata()
}

func (a *KafkaMonitor) StartListenRecvdata() {
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

func catch_depth_exp(depth *datastruct.DepthQuote) {
	errMsg := recover()
	if errMsg != nil {
		// fmt.Println("This is catch_exp func")
		logx.Errorf("catch_exp depth:  %+v\n", depth.String(3))
		logx.Errorf("errMsg: %+v \n", errMsg)

		logx.Infof("catch_exp depth:  %+v\n", depth.String(3))
		logx.Infof("errMsg: %+v \n", errMsg)
		// fmt.Println(errMsg)
	}
}

func (d *KafkaMonitor) process_depth(depth *datastruct.DepthQuote) error {

	defer catch_depth_exp(depth)

	return nil
}

func catch_kline_exp(kline *datastruct.Kline) {
	errMsg := recover()
	if errMsg != nil {
		// fmt.Println("This is catch_exp func")
		logx.Errorf("catch_exp kline:  %+v\n", kline.String())
		logx.Errorf("errMsg: %+v \n", errMsg)

		logx.Infof("catch_exp kline:  %+v\n", kline.String())
		logx.Infof("errMsg: %+v \n", errMsg)
		// fmt.Println(errMsg)
	}
}

func (d *KafkaMonitor) process_kline(kline *datastruct.Kline) error {
	defer catch_kline_exp(kline)
	// kline.Time = kline.Time / datastruct.NANO_PER_SECS

	// logx.Statf("Rcv kline: %s", kline.String())

	// d.PublishChangeinfo(d.cache_period_data[kline.Symbol].GetChangeInfo(), nil)

	return nil
}

func catch_trade_exp(trade *datastruct.Trade) {
	errMsg := recover()
	if errMsg != nil {
		// fmt.Println("This is catch_exp func")
		logx.Errorf("catch_exp trade:  %+v\n", trade.String())
		logx.Errorf("errMsg: %+v \n", errMsg)

		logx.Infof("catch_exp trade:  %+v\n", trade.String())
		logx.Infof("errMsg: %+v \n", errMsg)
		// fmt.Println(errMsg)
	}
}

func (d *KafkaMonitor) process_trade(trade *datastruct.Trade) error {
	defer catch_trade_exp(trade)

	return nil
}
