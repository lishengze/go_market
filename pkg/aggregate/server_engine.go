package aggregate

import (
	"market_aggregate/pkg/comm"
	config "market_aggregate/pkg/conf"
	"market_aggregate/pkg/datastruct"
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

	s.AggregateWorkr.Start()
}
