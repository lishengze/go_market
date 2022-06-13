package main

import (
	"flag"
	"fmt"
	"market_server/app/market_aggregate/aggregate"
	main_config "market_server/app/market_aggregate/config"
	mkconfig "market_server/app/market_aggregate/config"
	config "market_server/common/config"
	"market_server/common/datastruct"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

type TestEngine struct {
	cfg *main_config.Config

	NacosClientWorker *config.NacosClient

	IsTest bool

	RiskConfigMutex    sync.Mutex
	RiskCtrlConfigMaps aggregate.RiskCtrlConfigMap
	HedgingConfigs     []*mkconfig.HedgingConfig
	MarketRiskConfigs  []*mkconfig.MarketRiskConfig
	SymbolConfigs      []*mkconfig.SymbolConfig
}

func NewServerEngine(cfg *main_config.Config) *TestEngine {

	s := &TestEngine{}
	s.cfg = cfg

	s.HedgingConfigs = nil
	s.MarketRiskConfigs = nil
	s.SymbolConfigs = nil
	s.RiskCtrlConfigMaps = make(map[string]*mkconfig.RiskCtrlConfig)
	s.IsTest = false

	return s
}

func (s *TestEngine) Start() {

	go s.StartNacosClient()
}

func (s *TestEngine) SetTestFlag(value bool) {
	s.IsTest = value
}

func (s *TestEngine) StartNacosClient() {
	fmt.Println("****************** StartNacosClient *****************")
	s.NacosClientWorker = config.NewNacosClient(&s.cfg.Nacos)

	fmt.Println("Connect Nacos Successfully!")

	MarketRiskConfigStr, err := s.NacosClientWorker.GetConfigContent("MarketRisk", datastruct.BCTS_GROUP)
	if err != nil {
		logx.Error(err.Error())
	}
	fmt.Println("Requested MarketRisk: " + MarketRiskConfigStr)
	s.ProcessMarketriskConfigStr(MarketRiskConfigStr)

	HedgeConfigStr, err := s.NacosClientWorker.GetConfigContent("HedgeParams", datastruct.BCTS_GROUP)
	if err != nil {
		logx.Error(err.Error())
	}
	fmt.Println("Requested HedgeConfigStr: " + HedgeConfigStr)
	s.ProcsssHedgeConfigStr(HedgeConfigStr)

	SymbolConfigStr, err := s.NacosClientWorker.GetConfigContent("SymbolParams", datastruct.BCTS_GROUP)
	if err != nil {
		logx.Error(err.Error())
	}
	fmt.Println("Requested SymbolConfigStr: " + SymbolConfigStr)
	s.ProcessSymbolConfigStr(SymbolConfigStr)

	// s.NacosClientWorker.ListenConfig("MarketRisk", datastruct.BCTS_GROUP, s.MarketRiskChanged)

	// s.NacosClientWorker.ListenConfig("HedgeParams", datastruct.BCTS_GROUP, s.HedgeParamsChanged)

	// s.NacosClientWorker.ListenConfig("SymbolParams", datastruct.BCTS_GROUP, s.SymbolParamsChanged)
}

func (s *TestEngine) HedgeParamsChanged(namespace, group, dataId, hedgingContent string) {
	fmt.Printf("HedgeParamsChanged hedgingContent: %s\n", hedgingContent)
	s.ProcsssHedgeConfigStr(hedgingContent)
}

func (s *TestEngine) ProcsssHedgeConfigStr(data string) {
	hedge_configs, err := mkconfig.ParseJsonHedgerConfig(data)
	if err != nil {
		logx.Error(err.Error())
		return
	}

	s.HedgingConfigs = hedge_configs

	symbol_exchange_set := make(map[string](map[string]struct{}))

	NewMeta := datastruct.Metadata{}
	for _, hedge_config := range hedge_configs {
		if _, ok := symbol_exchange_set[hedge_config.Symbol]; !ok {
			symbol_exchange_set[hedge_config.Symbol] = make(map[string]struct{})
		}

		symbol_exchange_set[hedge_config.Symbol][hedge_config.Exchange] = struct{}{}
		fmt.Printf("New Meta: %s.%s", hedge_config.Symbol, hedge_config.Exchange)
	}

	NewMeta.DepthMeta = symbol_exchange_set
	NewMeta.TradeMeta = symbol_exchange_set
	NewMeta.KlineMeta = symbol_exchange_set

	fmt.Printf("HedgeParamsChanged: NewMeta:\n%s \n", NewMeta.String())

	s.UpdateRiskConfigHedgePart(hedge_configs)
}

func (s *TestEngine) MarketRiskChanged(namespace, group, dataId, data string) {
	fmt.Printf("MarketRiskContent: %s\n", data)
	s.ProcessMarketriskConfigStr(data)
}

func (s *TestEngine) ProcessMarketriskConfigStr(data string) {
	market_risk_configs, err := mkconfig.ParseJsonMarketRiskConfig(data)

	if err != nil {
		logx.Error(err.Error())
		return
	}

	s.MarketRiskConfigs = market_risk_configs

	NewAggConfig := mkconfig.AggregateConfig{}
	NewAggConfig.DepthAggregatorConfigMap = make(map[string]mkconfig.AggregateConfigAtom)

	for _, market_risk_config := range market_risk_configs {
		NewAggConfig.DepthAggregatorConfigMap[market_risk_config.Symbol] = mkconfig.AggregateConfigAtom{
			AggregateFreq: time.Duration(market_risk_config.PublishFrequency),
			PublishLevel:  int(market_risk_config.PublishLevel),
			IsPublish:     bool(market_risk_config.Switch)}
	}

	fmt.Printf("MarketRiskChanged:\n%s \n", NewAggConfig.String())

	s.UpdateRiskConfigRiskPart(market_risk_configs)
}

func (s *TestEngine) SymbolParamsChanged(namespace, group, dataId, data string) {
	fmt.Printf("SymbolContent: %s\n", data)
	s.ProcessSymbolConfigStr(data)
}

func (s *TestEngine) ProcessSymbolConfigStr(data string) {
	symbol_configs, err := mkconfig.ParseJsonSymbolConfig(data)

	if err != nil {
		logx.Error(err.Error())
		return
	}

	s.SymbolConfigs = symbol_configs

	s.UpdateRiskConfigSymbolPart(symbol_configs)
}

/*
   "fee_kind":1,
   "taker_fee":0,
*/
func (s *TestEngine) UpdateRiskConfigHedgePart(hedge_configs []*mkconfig.HedgingConfig) {

	for _, hedge_config := range hedge_configs {
		if _, ok := s.RiskCtrlConfigMaps[hedge_config.Symbol]; !ok {
			s.RiskCtrlConfigMaps[hedge_config.Symbol] = &mkconfig.RiskCtrlConfig{}
			s.RiskCtrlConfigMaps[hedge_config.Symbol].HedgeConfigMap = make(map[string]*mkconfig.HedgeConfig)

			fmt.Println("Risk New Symbol: " + hedge_config.Symbol + "\n")
		}

		fmt.Println("Risk Update Symbol: " + hedge_config.Symbol + "\n")

		if s.RiskCtrlConfigMaps[hedge_config.Symbol].HedgeConfigMap == nil {
			s.RiskCtrlConfigMaps[hedge_config.Symbol].HedgeConfigMap = make(map[string]*mkconfig.HedgeConfig)
		}

		s.RiskCtrlConfigMaps[hedge_config.Symbol].HedgeConfigMap[hedge_config.Exchange] = &mkconfig.HedgeConfig{
			FeeKind:  hedge_config.FeeKind,
			FeeValue: hedge_config.TakerFee,
		}
	}

	s.UpdateRiskConfig()

}

/*
        "price_offset_kind":1,
        "price_offset":0.001,
*/
func (s *TestEngine) UpdateRiskConfigRiskPart(market_risk_configs []*mkconfig.MarketRiskConfig) {

	for _, market_risk_config := range market_risk_configs {
		if _, ok := s.RiskCtrlConfigMaps[market_risk_config.Symbol]; !ok {
			s.RiskCtrlConfigMaps[market_risk_config.Symbol] = &mkconfig.RiskCtrlConfig{}
		}

		s.RiskCtrlConfigMaps[market_risk_config.Symbol].PriceBiasKind = market_risk_config.PriceOffsetKind
		s.RiskCtrlConfigMaps[market_risk_config.Symbol].PriceBiasValue = market_risk_config.PriceOffset
		s.RiskCtrlConfigMaps[market_risk_config.Symbol].VolumeBiasKind = market_risk_config.AmountOffsetKind
		s.RiskCtrlConfigMaps[market_risk_config.Symbol].VolumeBiasValue = market_risk_config.AmountOffset
	}

	s.UpdateRiskConfig()
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
func (s *TestEngine) UpdateRiskConfigSymbolPart(SymbolConfigs []*mkconfig.SymbolConfig) {

	s.RiskConfigMutex.Lock()

	for _, symbol_config := range SymbolConfigs {
		if _, ok := s.RiskCtrlConfigMaps[symbol_config.Symbol]; !ok {
			s.RiskCtrlConfigMaps[symbol_config.Symbol] = &mkconfig.RiskCtrlConfig{}
		}

		s.RiskCtrlConfigMaps[symbol_config.Symbol].PricePrecison = uint32(symbol_config.PricePrecision)
		s.RiskCtrlConfigMaps[symbol_config.Symbol].VolumePrecison = uint32(symbol_config.AmountPrecision)
		s.RiskCtrlConfigMaps[symbol_config.Symbol].PriceMinumChange = symbol_config.MinChangePrice
	}

	s.UpdateRiskConfig()

	s.RiskConfigMutex.Unlock()
}

func (s *TestEngine) UpdateRiskConfig() {
	if s.HedgingConfigs != nil && s.MarketRiskConfigs != nil && s.SymbolConfigs != nil {

		fmt.Printf("s.RiskCtrlConfigMaps: \n%s", aggregate.GetRiskCtrlConfigMapString(&s.RiskCtrlConfigMaps))
	}
}

func TestNacos() {
	var configFile = flag.String("f", "client.yaml", "the config file")

	var c main_config.Config
	conf.MustLoad(*configFile, &c)

	logx.MustSetup(c.LogConfig)

	fmt.Printf("Log: %s \n", c.String())

	test := NewServerEngine(&c)

	test.Start()
}

func main() {
	TestNacos()
}
