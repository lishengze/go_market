package monitor_market

import (
	"market_server/common/datastruct"
	"market_server/common/monitorStruct"

	"github.com/zeromicro/go-zero/core/logx"
)

type KafkaMonitor struct {
	RecvDataChan            *datastruct.DataChannel
	MonitorMarketDataWorker *monitorStruct.MonitorMarketData
	MonitorChan             *monitorStruct.MonitorChannel

	RateParam    float64
	InitDeadLine int64
	CheckSecs    int64
}

func NewKafkaMonitor(recvDataChan *datastruct.DataChannel) *KafkaMonitor {
	monitor_chan := monitorStruct.NewMonitorChannel()
	rate_param := 1.8
	init_dead_line := datastruct.NANO_PER_MIN * 10
	check_secs := 15 * datastruct.SECS_PER_MIN

	return &KafkaMonitor{
		RecvDataChan: recvDataChan,
		MonitorChan:  monitor_chan,
		RateParam:    1.8,
		InitDeadLine: int64(init_dead_line),
		CheckSecs:    int64(check_secs),

		MonitorMarketDataWorker: monitorStruct.NewMonitorMarketData(rate_param, int64(init_dead_line),
			int64(check_secs), monitor_chan),
	}
}

func (k *KafkaMonitor) Start() {
	logx.Infof("KafkaMonitor Start!")

	k.StartListenRecvdata()
}

func (k *KafkaMonitor) StartListenRecvdata() {
	logx.Info("[S] KafkaMonitor start_listen_recvdata")
	go func() {
		for {
			select {
			case new_depth := <-k.RecvDataChan.DepthChannel:
				go k.process_depth(new_depth)
			case new_kline := <-k.RecvDataChan.KlineChannel:
				go k.process_kline(new_kline)
			case new_trade := <-k.RecvDataChan.TradeChannel:
				go k.process_trade(new_trade)
			}
		}
	}()
	logx.Info("[S] DBServer start_receiver Over!")
}

func (k *KafkaMonitor) StartListenInvalidData() {
	logx.Info("[S] KafkaMonitor StartListenInvalidData")
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

func (k *KafkaMonitor) process_depth(depth *datastruct.DepthQuote) error {

	defer catch_depth_exp(depth)

	k.MonitorMarketDataWorker.UpdateDepth(depth.Symbol)

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

func (k *KafkaMonitor) process_kline(kline *datastruct.Kline) error {
	defer catch_kline_exp(kline)

	k.MonitorMarketDataWorker.UpdateKline(kline.Symbol)

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

func (k *KafkaMonitor) process_trade(trade *datastruct.Trade) error {
	defer catch_trade_exp(trade)

	k.MonitorMarketDataWorker.UpdateTrade(trade.Symbol)

	return nil
}

func (k *KafkaMonitor) process_invalid_depth(montior_atom *monitorStruct.MonitorAtom) error {

	return nil
}

func (k *KafkaMonitor) process_invalid_trade(montior_atom *monitorStruct.MonitorAtom) error {

	return nil
}

func (k *KafkaMonitor) process_invalid_kline(montior_atom *monitorStruct.MonitorAtom) error {

	return nil
}
