package aggregate

import (
	"fmt"
	"market_aggregate/pkg/comm"
	config "market_aggregate/pkg/conf"
	"market_aggregate/pkg/datastruct"
	"market_aggregate/pkg/util"
)

type ServerEngine struct {
	AggregateWorkr *Aggregator
	Riskworker     *RiskWorkerManager
	Commer         comm.Comm

	RecvDataChan *datastruct.DataChannel
	PubDataChan  *datastruct.DataChannel

	MetaData  datastruct.Metadata
	AggConfig config.AggregateConfig
}

func (s *ServerEngine) Init() {
	s.RecvDataChan = &datastruct.DataChannel{
		DepthChannel: make(chan *datastruct.DepthQuote),
		KlineChannel: make(chan *datastruct.Kline),
		TradeChannel: make(chan *datastruct.Trade),
	}

	s.PubDataChan = &datastruct.DataChannel{
		DepthChannel: make(chan *datastruct.DepthQuote),
		KlineChannel: make(chan *datastruct.Kline),
		TradeChannel: make(chan *datastruct.Trade),
	}

	s.Commer = comm.Comm{}
	s.Commer.Init(s.RecvDataChan, s.PubDataChan)

	s.Riskworker = &RiskWorkerManager{}
	s.Riskworker.Init()

	s.AggregateWorkr = &Aggregator{}
	s.AggregateWorkr.Init(s.RecvDataChan, s.PubDataChan, s.Riskworker)

}

func (s *ServerEngine) Start() {
	risk_config := GetTestRiskConfig()
	s.Riskworker.UpdateConfig(&risk_config)

	AggConfig := GetTestAggConfig()
	s.AggregateWorkr.UpdateConfig(AggConfig)

	s.Commer.Start()
	s.AggregateWorkr.Start()
}

func TestServerEngine() {
	config.NATIVE_CONFIG_INIT("client.yaml")

	util.LOG_INFO(fmt.Sprintf("CONFIG: %+v", *config.NATIVE_CONFIG()))

	server_engine := new(ServerEngine)
	server_engine.Init()
	server_engine.Start()

	risk_config := GetTestRiskConfig()
	util.LOG_INFO(fmt.Sprintf("risk_config: %+v", risk_config))

	server_engine.Riskworker.UpdateConfig(&risk_config)

	AggConfig := GetTestAggConfig()
	util.LOG_INFO(fmt.Sprintf("AggConfig: %+v", AggConfig))

	server_engine.AggregateWorkr.UpdateConfig(AggConfig)

	meta_data := datastruct.GetTestMetadata()
	util.LOG_INFO(fmt.Sprintf("meta_data: %+v", meta_data))

	server_engine.Commer.UpdateMetaData(meta_data)

	select {}
}
