package monitor_market

import (
	"market_server/common/datastruct"
	"market_server/common/dingtalk"
	"market_server/common/monitorStruct"

	"github.com/zeromicro/go-zero/core/logx"
)

type MonitorEngine struct {
	MonitorDataChan  chan *monitorStruct.MonitorData
	OriginalDataChan *datastruct.DataChannel

	MonitorMarketDataWorker *monitorStruct.MonitorMarketData
	MonitorChan             *monitorStruct.MonitorChannel

	DingClient *dingtalk.Client
	MetaInfo   string

	RateParam    float64
	InitDeadLine int64
	CheckSecs    int64
}

func NewMonitorEngineWithOrignalDataChannel(original_data_chan *datastruct.DataChannel, monitor_config *monitorStruct.MonitorConfig,
	ding_config *DingConfig, meta_info string) *MonitorEngine {
	monitor_chan := monitorStruct.NewMonitorChannel()
	dingtalk := dingtalk.NewClient(ding_config.Token, ding_config.Secret)
	return &MonitorEngine{
		MonitorDataChan:  nil,
		OriginalDataChan: original_data_chan,

		MonitorChan:  monitor_chan,
		RateParam:    monitor_config.RateParam,
		InitDeadLine: monitor_config.InitDeadLine,
		CheckSecs:    monitor_config.CheckSecs,
		DingClient:   dingtalk,

		MetaInfo:                meta_info,
		MonitorMarketDataWorker: monitorStruct.NewMonitorMarketData(meta_info, monitor_config, monitor_chan),
	}
}

func NewMonitorEngineWithMonitorDataChannel(monitor_data_chan chan *monitorStruct.MonitorData, monitor_config *monitorStruct.MonitorConfig,
	ding_config *DingConfig, meta_info string) *MonitorEngine {
	monitor_chan := monitorStruct.NewMonitorChannel()
	dingtalk := dingtalk.NewClient(ding_config.Token, ding_config.Secret)

	return &MonitorEngine{
		MonitorDataChan:  monitor_data_chan,
		OriginalDataChan: nil,

		MonitorChan:  monitor_chan,
		RateParam:    monitor_config.RateParam,
		InitDeadLine: monitor_config.InitDeadLine,
		CheckSecs:    monitor_config.CheckSecs,
		DingClient:   dingtalk,

		MetaInfo:                meta_info,
		MonitorMarketDataWorker: monitorStruct.NewMonitorMarketData(meta_info, monitor_config, monitor_chan),
	}
}

func (k *MonitorEngine) Start() {
	logx.Infof("MonitorEngine Start!")

	k.StartListenRecvdata()
}

func (k *MonitorEngine) StartListenRecvdata() {
	logx.Info("[S] MonitorEngine start_listen_recvdata")

	if k.MonitorDataChan != nil {
		go func() {
			for {
				select {
				case data := <-k.MonitorDataChan:
					if data.DataType == datastruct.DEPTH_TYPE {
						go k.process_depth(data.Symbol)
					} else if data.DataType == datastruct.TRADE_TYPE {
						go k.process_trade(data.Symbol)
					} else if data.DataType == datastruct.KLINE_TYPE {
						go k.process_kline(data.Symbol)
					}
				}
			}
		}()
	} else if k.OriginalDataChan != nil {
		go func() {
			for {
				select {
				case depth := <-k.OriginalDataChan.DepthChannel:
					go k.process_depth(depth.Exchange + "_" + depth.Symbol)
				case trade := <-k.OriginalDataChan.TradeChannel:
					go k.process_trade(trade.Exchange + "_" + trade.Symbol)
				case kline := <-k.OriginalDataChan.KlineChannel:
					go k.process_kline(kline.Exchange + "_" + kline.Symbol)
				}
			}
		}()
	}

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

func catch_exp() {
	errMsg := recover()
	if errMsg != nil {
		logx.Errorf("errMsg: %+v \n", errMsg)
		logx.Infof("errMsg: %+v \n", errMsg)
	}
}

func (k *MonitorEngine) process_depth(symbol string) error {

	defer catch_exp()

	logx.Slowf("%s Update depth %s", k.MetaInfo, symbol)

	k.MonitorMarketDataWorker.UpdateDepth(symbol)

	return nil
}

func (k *MonitorEngine) process_kline(symbol string) error {
	defer catch_exp()

	logx.Slowf("%s Update kline %s", k.MetaInfo, symbol)

	k.MonitorMarketDataWorker.UpdateKline(symbol)

	return nil
}

func (k *MonitorEngine) process_trade(symbol string) error {
	defer catch_exp()

	logx.Slowf("%s Update Trade %s", k.MetaInfo, symbol)

	k.MonitorMarketDataWorker.UpdateTrade(symbol)

	return nil
}

func (k *MonitorEngine) process_invalid_depth(montior_atom *monitorStruct.MonitorAtom) error {
	logx.Info(montior_atom.InvalidInfo)
	k.DingClient.SendMessage(k.MetaInfo + "\n" + montior_atom.InvalidInfo)
	return nil
}

func (k *MonitorEngine) process_invalid_trade(montior_atom *monitorStruct.MonitorAtom) error {
	logx.Info(montior_atom.InvalidInfo)
	k.DingClient.SendMessage(k.MetaInfo + "\n" + montior_atom.InvalidInfo)
	return nil
}

func (k *MonitorEngine) process_invalid_kline(montior_atom *monitorStruct.MonitorAtom) error {
	logx.Info(montior_atom.InvalidInfo)
	k.DingClient.SendMessage(k.MetaInfo + "\n" + montior_atom.InvalidInfo)
	return nil
}
