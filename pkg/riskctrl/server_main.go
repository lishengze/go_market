package riskctrl

import (
	"market_aggregate/pkg/comm"
	"market_aggregate/pkg/conf"
	"market_aggregate/pkg/datastruct"
)

type ServerEngine struct {
	AggregateWorkr Aggregator
	Riskworker     RiskWorkerManager
	Commer         comm.Comm

	RecvDataChan *datastruct.DataChannel
	PubDataChan  *datastruct.DataChannel

	MetaData  datastruct.Metadata
	AggConfig conf.AggregateConfig
}

func (s *ServerEngine) Init(config *conf.Config, meta_data datastruct.Metadata) {
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

	s.AggConfig = conf.AggregateConfig{
		DepthAggregatorMillsecs: 5,
	}

	s.Commer = comm.Comm{}
	s.Commer.Init(config, s.RecvDataChan, s.PubDataChan, meta_data)

	s.AggregateWorkr = Aggregator{}
	s.AggregateWorkr.Init(s.RecvDataChan, s.PubDataChan, s.AggConfig)

}

func (*ServerEngine) Start() {

}