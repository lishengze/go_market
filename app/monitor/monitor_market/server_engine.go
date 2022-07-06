package monitor_market

import (
	"market_server/common/comm"
	"market_server/common/datastruct"
	"market_server/common/monitorStruct"

	"github.com/zeromicro/go-zero/core/logx"
)

type ServerEngine struct {
	ConfigInfo     *Config
	Commer         *comm.Comm
	WSClientWorker *WSClient

	// OriginalDataChan *datastruct.DataChannel
	KafkaMonitor *MonitorEngine

	// MonitorDataChan chan *monitorStruct.MonitorData
	WSMonotr *MonitorEngine
}

func NewServerEngine(config_info *Config) *ServerEngine {
	MonitorDataChan := make(chan *monitorStruct.MonitorData)
	WSClientWorker := NewWSClient(&config_info.WS, config_info.MonitorMetaInfo.Symbols, MonitorDataChan)

	OriginalDataChan := datastruct.NewDataChannel()
	Commer := comm.NewComm(OriginalDataChan, nil, config_info.Comm)

	KafkaMonitor := NewMonitorEngineWithOrignalDataChannel(OriginalDataChan, &config_info.MonitorConfigInfo, &config_info.DingConfigInfo, "Kafka ")
	WSMonotr := NewMonitorEngineWithMonitorDataChannel(MonitorDataChan, &config_info.MonitorConfigInfo, &config_info.DingConfigInfo, "WS ")
	return &ServerEngine{
		ConfigInfo:     config_info,
		Commer:         Commer,
		WSClientWorker: WSClientWorker,

		KafkaMonitor: KafkaMonitor,
		WSMonotr:     WSMonotr,
	}
}

func (s *ServerEngine) StartCommer() {

	symbol_exchange_set := make(map[string](map[string]struct{}))
	new_meta := datastruct.Metadata{}
	for _, symbol := range s.ConfigInfo.MonitorMetaInfo.Symbols {
		if _, ok := symbol_exchange_set[symbol]; !ok {
			symbol_exchange_set[symbol] = make(map[string]struct{})
		}
		if _, ok := symbol_exchange_set[symbol][datastruct.BCTS_EXCHANGE]; !ok {
			symbol_exchange_set[symbol][datastruct.BCTS_EXCHANGE] = struct{}{}
		}
	}

	new_meta.TradeMeta = symbol_exchange_set
	new_meta.KlineMeta = symbol_exchange_set

	logx.Infof("[I] InitMeta: %v \n", new_meta)

	s.Commer.UpdateMetaData(&new_meta)
}

func (s *ServerEngine) Start() {
	go s.StartCommer()
	go s.WSClientWorker.Start()
}
