package aggregate

import (
	"fmt"
	"market_aggregate/pkg/comm"
	config "market_aggregate/pkg/conf"
	"market_aggregate/pkg/datastruct"
	"market_aggregate/pkg/util"
	"sync"
	"time"
)

type ServerEngine struct {
	AggregateWorker *Aggregator
	Riskworker      *RiskWorkerManager
	Commer          comm.Comm

	RecvDataChan *datastruct.DataChannel
	PubDataChan  *datastruct.DataChannel

	// MetaData   datastruct.Metadata
	// AggConfig  config.AggregateConfig
	// RiskConfig RiskCtrlConfigMap

	RiskConfigMutex    sync.Mutex
	RiskCtrlConfigMaps RiskCtrlConfigMap
	HedgingConfigs     []*config.HedgingConfig
	MarketRiskConfigs  []*config.MarketRiskConfig
	SymbolConfigs      []*config.SymbolConfig

	NacosClientWorker *config.NacosClient
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

	s.HedgingConfigs = nil
	s.MarketRiskConfigs = nil
	s.SymbolConfigs = nil
	s.RiskCtrlConfigMaps = make(map[string]*config.RiskCtrlConfig)

	s.InitConfig()
	go s.StartNacosClient()

	// s.Commer = comm.Comm{}
	// s.Commer.Init(s.RecvDataChan, s.PubDataChan)

	// s.Riskworker = &RiskWorkerManager{}
	// s.Riskworker.Init()

	// s.AggregateWorker = &Aggregator{}
	// s.AggregateWorker.Init(s.RecvDataChan, s.PubDataChan, s.Riskworker)

	// go s.StartNacosClient()
}

func (s *ServerEngine) Start() {

	// s.Commer.Start()
	// s.AggregateWorker.Start()
	// go s.StartNacosClient()
}

func (s *ServerEngine) InitConfig() {
	config.NATIVE_CONFIG_INIT("client.yaml")
	util.LOG_INFO(fmt.Sprintf("CONFIG: %+v", *config.NATIVE_CONFIG()))
	util.LOG_INFO(fmt.Sprintf("NacoIP: %s:%d", config.NATIVE_CONFIG().Nacos.IpAddr, config.NATIVE_CONFIG().Nacos.Port))

}

func (s *ServerEngine) StartNacosClient() {
	util.LOG_INFO("****************** StartNacosClient *****************")
	s.NacosClientWorker = config.NewNacosClient(&config.NATIVE_CONFIG().Nacos)

	util.LOG_INFO("Connect Nacos Successfully!")

	s.NacosClientWorker.ListenConfig("MarketRisk", datastruct.BCTS_GROUP, s.MarketRiskChanged)

	s.NacosClientWorker.ListenConfig("HedgeParams", datastruct.BCTS_GROUP, s.HedgeParamsChanged)

	s.NacosClientWorker.ListenConfig("SymbolParams", datastruct.BCTS_GROUP, s.SymbolParamsChanged)
}

func (s *ServerEngine) HedgeParamsChanged(namespace, group, dataId, hedgingContent string) {
	util.LOG_INFO(fmt.Sprintf("hedgingContent: %s\n", hedgingContent))
	hedge_configs, err := config.ParseJsonHedgerConfig(hedgingContent)
	if err != nil {
		util.LOG_ERROR(err.Error())
		return
	}

	s.HedgingConfigs = hedge_configs

	symbol_exchange_set := make(map[string](map[string]struct{}))

	NewMeta := datastruct.Metadata{}
	for _, hedge_config := range hedge_configs {
		if _, ok := symbol_exchange_set[hedge_config.Symbol]; ok == false {
			symbol_exchange_set[hedge_config.Symbol] = make(map[string]struct{})
		}

		symbol_exchange_set[hedge_config.Symbol][hedge_config.Exchange] = struct{}{}
	}

	NewMeta.DepthMeta = symbol_exchange_set
	NewMeta.TradeMeta = symbol_exchange_set
	NewMeta.KlineMeta = symbol_exchange_set

	// s.Commer.UpdateMetaData(&NewMeta)

	util.LOG_INFO(fmt.Sprintf("HedgeParamsChanged: NewMeta:%+v \n", NewMeta))

	s.UpdateRiskConfigHedgePart(hedge_configs)
}

func (s *ServerEngine) MarketRiskChanged(namespace, group, dataId, data string) {
	util.LOG_INFO(fmt.Sprintf("MarketRiskContent: %s\n", data))

	market_risk_configs, err := config.ParseJsonMarketRiskConfig(data)

	if err != nil {
		util.LOG_ERROR(err.Error())
		return
	}

	s.MarketRiskConfigs = market_risk_configs

	NewAggConfig := config.AggregateConfig{}
	NewAggConfig.DepthAggregatorConfigMap = make(map[string]config.AggregateConfigAtom)

	for _, market_risk_config := range market_risk_configs {
		NewAggConfig.DepthAggregatorConfigMap[market_risk_config.Symbol] = config.AggregateConfigAtom{
			AggregateFreq: time.Duration(market_risk_config.PublishFrequency),
			PublishLevel:  int(market_risk_config.PublishLevel),
			IsPublish:     bool(market_risk_config.Switch)}
	}

	util.LOG_INFO(fmt.Sprintf("MarketRiskChanged:%+v \n", NewAggConfig))

	// s.AggregateWorker.UpdateConfig(NewAggConfig)

	s.UpdateRiskConfigRiskPart(market_risk_configs)
}

func (s *ServerEngine) SymbolParamsChanged(namespace, group, dataId, data string) {
	util.LOG_INFO(fmt.Sprintf("SymbolContent: %s\n", data))

	symbol_configs, err := config.ParseJsonSymbolConfig(data)

	if err != nil {
		util.LOG_ERROR(err.Error())
		return
	}

	s.SymbolConfigs = symbol_configs

	s.UpdateRiskConfigSymbolPart(symbol_configs)
}

/*
   "fee_kind":1,
   "taker_fee":0,
*/
func (s *ServerEngine) UpdateRiskConfigHedgePart(hedge_configs []*config.HedgingConfig) {

	s.RiskConfigMutex.Lock()

	for _, hedge_config := range hedge_configs {
		if _, ok := s.RiskCtrlConfigMaps[hedge_config.Symbol]; ok == false {
			s.RiskCtrlConfigMaps[hedge_config.Symbol] = &config.RiskCtrlConfig{}
			s.RiskCtrlConfigMaps[hedge_config.Symbol].HedgeConfigMap = make(map[string]*config.HedgeConfig)
		}

		s.RiskCtrlConfigMaps[hedge_config.Symbol].HedgeConfigMap[hedge_config.Exchange] = &config.HedgeConfig{
			FeeKind:  hedge_config.FeeKind,
			FeeValue: hedge_config.TakerFee,
		}
	}

	if s.HedgingConfigs != nil && s.MarketRiskConfigs != nil && s.SymbolConfigs != nil {
		s.Riskworker.UpdateConfig(&s.RiskCtrlConfigMaps)
	}

	s.RiskConfigMutex.Unlock()
}

/*
        "price_offset_kind":1,
        "price_offset":0.001,
*/
func (s *ServerEngine) UpdateRiskConfigRiskPart(market_risk_configs []*config.MarketRiskConfig) {

	s.RiskConfigMutex.Lock()

	for _, market_risk_config := range market_risk_configs {
		if _, ok := s.RiskCtrlConfigMaps[market_risk_config.Symbol]; ok == false {
			s.RiskCtrlConfigMaps[market_risk_config.Symbol] = &config.RiskCtrlConfig{}
		}

		s.RiskCtrlConfigMaps[market_risk_config.Symbol].PriceBiasKind = market_risk_config.PriceOffsetKind
		s.RiskCtrlConfigMaps[market_risk_config.Symbol].PriceBiasValue = market_risk_config.PriceOffset
		s.RiskCtrlConfigMaps[market_risk_config.Symbol].VolumeBiasKind = market_risk_config.AmountOffsetKind
		s.RiskCtrlConfigMaps[market_risk_config.Symbol].VolumeBiasValue = market_risk_config.AmountOffset
	}

	s.UpdateRiskConfig()

	s.RiskConfigMutex.Unlock()
}

/*
        "amount_precision":4,
        "price_precision":2,
        "sum_precision":4,
	PricePrecison  uint32
	VolumePrecison uint32

	PriceBiasValue float64
	PriceBiasKind  int

	VolumeBiasValue float64
	VolumeBiasKind  int

	PriceMinumChange float64

*/
func (s *ServerEngine) UpdateRiskConfigSymbolPart(SymbolConfigs []*config.SymbolConfig) {

	s.RiskConfigMutex.Lock()

	for _, symbol_config := range SymbolConfigs {
		if _, ok := s.RiskCtrlConfigMaps[symbol_config.Symbol]; ok == false {
			s.RiskCtrlConfigMaps[symbol_config.Symbol] = &config.RiskCtrlConfig{}
		}

		s.RiskCtrlConfigMaps[symbol_config.Symbol].PricePrecison = uint32(symbol_config.PricePrecision)
		s.RiskCtrlConfigMaps[symbol_config.Symbol].VolumePrecison = uint32(symbol_config.AmountPrecision)
		s.RiskCtrlConfigMaps[symbol_config.Symbol].PriceMinumChange = symbol_config.MinChangePrice
	}

	s.UpdateRiskConfig()

	s.RiskConfigMutex.Unlock()
}

func (s *ServerEngine) UpdateRiskConfig() {
	if s.HedgingConfigs != nil && s.MarketRiskConfigs != nil && s.SymbolConfigs != nil {
		util.LOG_INFO(fmt.Sprintf("%+v", s.RiskCtrlConfigMaps))
		// s.Riskworker.UpdateConfig(&s.RiskCtrlConfigMaps)
	}
}

func TestServerEngine() {

	server_engine := new(ServerEngine)
	server_engine.Init()
	server_engine.Start()

	// risk_config := GetTestRiskConfig()
	// util.LOG_INFO(fmt.Sprintf("risk_config: %+v", risk_config))

	// server_engine.Riskworker.UpdateConfig(&risk_config)

	// AggConfig := GetTestAggConfig()
	// util.LOG_INFO(fmt.Sprintf("AggConfig: %+v", AggConfig))

	// server_engine.AggregateWorker.UpdateConfig(AggConfig)

	// meta_data := datastruct.GetTestMetadata()
	// util.LOG_INFO(fmt.Sprintf("meta_data: %+v", meta_data))

	// server_engine.Commer.UpdateMetaData(meta_data)

	select {}
}
