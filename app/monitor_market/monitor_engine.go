package monitor_market

import (
	"market_server/common/datastruct"
	"market_server/common/dingtalk"
	"market_server/common/monitorStruct"

	"github.com/zeromicro/go-zero/core/logx"
)

type MonitorEngine struct {
	RecvDataChan            chan *monitorStruct.MonitorData
	MonitorMarketDataWorker *monitorStruct.MonitorMarketData
	MonitorChan             *monitorStruct.MonitorChannel
	DingClient              *dingtalk.Client

	RateParam    float64
	InitDeadLine int64
	CheckSecs    int64
}

func NewMonitorEngine(recvDataChan chan *monitorStruct.MonitorData, monitor_config *monitorStruct.MonitorConfig, ding_config *DingConfig) *MonitorEngine {
	monitor_chan := monitorStruct.NewMonitorChannel()
	dingtalk := dingtalk.NewClient(ding_config.token, ding_config.secret)

	return &MonitorEngine{
		RecvDataChan: recvDataChan,
		MonitorChan:  monitor_chan,
		RateParam:    monitor_config.RateParam,
		InitDeadLine: monitor_config.InitDeadLine,
		CheckSecs:    monitor_config.CheckSecs,
		DingClient:   dingtalk,

		MonitorMarketDataWorker: monitorStruct.NewMonitorMarketData(monitor_config, monitor_chan),
	}
}

func (k *MonitorEngine) Start() {
	logx.Infof("MonitorEngine Start!")

	k.StartListenRecvdata()
}

func (k *MonitorEngine) StartListenRecvdata() {
	logx.Info("[S] MonitorEngine start_listen_recvdata")
	go func() {
		for {
			select {
			case data := <-k.RecvDataChan:
				if data.DataType == datastruct.DEPTH_TYPE {
					k.process_depth(data.Symbol)
				} else if data.DataType == datastruct.TRADE_TYPE {
					k.process_trade(data.Symbol)
				} else if data.DataType == datastruct.KLINE_TYPE {
					k.process_kline(data.Symbol)
				}
			}
		}
	}()
	logx.Info("[S] DBServer start_receiver Over!")
}

func (k *MonitorEngine) StartListenInvalidData() {
	logx.Info("[S] MonitorEngine StartListenInvalidData")
	go func() {
		for {
			select {
			case invalid_depth := <-k.MonitorChan.DepthChan:
				go k.process_invalid_depth(invalid_depth)
			case invalid_trade := <-k.MonitorChan.TradeChan:
				go k.process_invalid_trade(invalid_trade)
			case invalid_kline := <-k.MonitorChan.KlineChan:
				go k.process_invalid_kline(invalid_kline)
			}
		}
	}()
	logx.Info("[S] DBServer start_receiver Over!")
}

func catch_depth_exp(depth *datastruct.DepthQuote) {
	errMsg := recover()
	if errMsg != nil {
		logx.Errorf("catch_exp depth:  %+v\n", depth.String(3))
		logx.Errorf("errMsg: %+v \n", errMsg)

		logx.Infof("catch_exp depth:  %+v\n", depth.String(3))
		logx.Infof("errMsg: %+v \n", errMsg)
	}
}

func (k *MonitorEngine) process_depth(symbol string) error {

	// defer catch_depth_exp(depth)

	k.MonitorMarketDataWorker.UpdateDepth(symbol)

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

func (k *MonitorEngine) process_kline(symbol string) error {
	// defer catch_kline_exp(kline)

	k.MonitorMarketDataWorker.UpdateKline(symbol)

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

func (k *MonitorEngine) process_trade(symbol string) error {
	// defer catch_trade_exp(trade)

	k.MonitorMarketDataWorker.UpdateTrade(symbol)

	return nil
}

func (k *MonitorEngine) process_invalid_depth(montior_atom *monitorStruct.MonitorAtom) error {

	return nil
}

func (k *MonitorEngine) process_invalid_trade(montior_atom *monitorStruct.MonitorAtom) error {

	return nil
}

func (k *MonitorEngine) process_invalid_kline(montior_atom *monitorStruct.MonitorAtom) error {

	return nil
}
