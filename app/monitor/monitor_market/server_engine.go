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

	KafkaMonitor *MonitorEngine
	WSMonotr     *MonitorEngine
}

func NewServerEngine(config_info *Config) *ServerEngine {
	rst := &ServerEngine{
		ConfigInfo:     config_info,
		Commer:         nil,
		WSClientWorker: nil,

		KafkaMonitor: nil,
		WSMonotr:     nil,
	}

	if config_info.MonitorObjectConfig.WSMonitor {
		MonitorDataChan := make(chan *monitorStruct.MonitorData)
		rst.WSClientWorker = NewWSClient(&config_info.WS, config_info.MonitorMetaInfo.Symbols, MonitorDataChan)
		rst.WSMonotr = NewMonitorEngineWithMonitorDataChannel(MonitorDataChan, &config_info.MonitorConfigInfo, &config_info.DingConfigInfo, "WS ")
	}

	if config_info.MonitorObjectConfig.KafkaMonitor {
		OriginalDataChan := datastruct.NewDataChannel()
		rst.Commer = comm.NewComm(OriginalDataChan, nil, config_info.Comm)
		rst.KafkaMonitor = NewMonitorEngineWithOrignalDataChannel(OriginalDataChan, &config_info.MonitorConfigInfo, &config_info.DingConfigInfo, "Kafka ")
	}

	return rst
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
	new_meta.DepthMeta = symbol_exchange_set

	logx.Infof("[I] InitMeta: %v \n", new_meta)

	s.Commer.Start()

	s.Commer.UpdateMetaData(&new_meta)
}

func (s *ServerEngine) Start() {
	if s.ConfigInfo.MonitorObjectConfig.KafkaMonitor {
		go s.StartCommer()
		go s.KafkaMonitor.Start()
	}

	if s.ConfigInfo.MonitorObjectConfig.WSMonitor {
		go s.WSClientWorker.Start()
		go s.WSMonotr.Start()
	}
}
