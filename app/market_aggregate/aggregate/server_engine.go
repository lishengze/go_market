package aggregate

import (
	"fmt"
	mkconfig "market_server/app/market_aggregate/config"
	"market_server/app/market_aggregate/svc"
	"market_server/common/comm"
	config "market_server/common/config"
	"market_server/common/datastruct"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type ServerEngine struct {
	ctx             *svc.ServiceContext
	AggregateWorker *Aggregator
	Riskworker      *RiskWorkerManager
	Commer          *comm.Comm

	RecvDataChan *datastruct.DataChannel
	PubDataChan  *datastruct.DataChannel

	RiskConfigMutex    sync.Mutex
	RiskCtrlConfigMaps RiskCtrlConfigMap
	HedgingConfigs     []*mkconfig.HedgingConfig
	MarketRiskConfigs  []*mkconfig.MarketRiskConfig
	SymbolConfigs      []*mkconfig.SymbolConfig

	NacosClientWorker *config.NacosClient

	IsTest bool
}

func NewServerEngine(svcCtx *svc.ServiceContext) *ServerEngine {

	s := &ServerEngine{}
	s.ctx = svcCtx
	s.RecvDataChan = datastruct.NewDataChannel()

	s.PubDataChan = datastruct.NewDataChannel()

	s.HedgingConfigs = nil
	s.MarketRiskConfigs = nil
	s.SymbolConfigs = nil
	s.RiskCtrlConfigMaps = make(map[string]*mkconfig.RiskCtrlConfig)
	s.IsTest = false
	s.Commer = comm.NewComm(s.RecvDataChan, s.PubDataChan, svcCtx.Config.Comm)

	s.Riskworker = NewRiskWorkerManager(&s.ctx.Config)

	s.AggregateWorker = NewAggregator(s.RecvDataChan, s.PubDataChan, s.Riskworker, &s.ctx.Config)
	return s
}

func (s *ServerEngine) Start() {

	s.Commer.Start()

	s.AggregateWorker.Start()

	if !s.IsTest {
		go s.StartNacosClient()
	}
}

func (s *ServerEngine) SetTestFlag(value bool) {
	s.IsTest = value
}

func (s *ServerEngine) StartNacosClient() {
	logx.Info("****************** StartNacosClient *****************")
	s.NacosClientWorker = config.NewNacosClient(&s.ctx.Config.Nacos)

	logx.Info("Connect Nacos Successfully!")

	MarketRiskConfigStr, err := s.NacosClientWorker.GetConfigContent("MarketRisk", datastruct.BCTS_GROUP)
	if err != nil {
		logx.Error(err.Error())
	}
	logx.Info("Requested MarketRisk: " + MarketRiskConfigStr)
	s.ProcessMarketriskConfigStr(MarketRiskConfigStr)

	HedgeConfigStr, err := s.NacosClientWorker.GetConfigContent("HedgeParams", datastruct.BCTS_GROUP)
	if err != nil {
		logx.Error(err.Error())
	}
	logx.Info("Requested HedgeConfigStr: " + HedgeConfigStr)
	s.ProcsssHedgeConfigStr(HedgeConfigStr)

	SymbolConfigStr, err := s.NacosClientWorker.GetConfigContent("SymbolParams", datastruct.BCTS_GROUP)
	if err != nil {
		logx.Error(err.Error())
	}
	logx.Info("Requested SymbolConfigStr: " + SymbolConfigStr)
	s.ProcessSymbolConfigStr(SymbolConfigStr)

	s.NacosClientWorker.ListenConfig("MarketRisk", datastruct.BCTS_GROUP, s.MarketRiskChanged)

	s.NacosClientWorker.ListenConfig("HedgeParams", datastruct.BCTS_GROUP, s.HedgeParamsChanged)

	s.NacosClientWorker.ListenConfig("SymbolParams", datastruct.BCTS_GROUP, s.SymbolParamsChanged)
}

func (s *ServerEngine) HedgeParamsChanged(namespace, group, dataId, hedgingContent string) {
	logx.Info(fmt.Sprintf("HedgeParamsChanged hedgingContent: %s\n", hedgingContent))
	s.ProcsssHedgeConfigStr(hedgingContent)
}

func (s *ServerEngine) ProcsssHedgeConfigStr(data string) {
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
		logx.Info(fmt.Sprintf("New Meta: %s.%s", hedge_config.Symbol, hedge_config.Exchange))
	}

	NewMeta.DepthMeta = symbol_exchange_set
	NewMeta.TradeMeta = symbol_exchange_set
	NewMeta.KlineMeta = symbol_exchange_set

	logx.Info(fmt.Sprintf("HedgeParamsChanged: NewMeta:\n%s \n", NewMeta.String()))

	s.Commer.UpdateMetaData(&NewMeta)

	s.UpdateRiskConfigHedgePart(hedge_configs)
}

func (s *ServerEngine) MarketRiskChanged(namespace, group, dataId, data string) {
	logx.Info(fmt.Sprintf("MarketRiskContent: %s\n", data))
	s.ProcessMarketriskConfigStr(data)
}

func (s *ServerEngine) ProcessMarketriskConfigStr(data string) {
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

	logx.Info(fmt.Sprintf("MarketRiskChanged:\n%s \n", NewAggConfig.String()))

	s.AggregateWorker.UpdateConfig(NewAggConfig)

	s.UpdateRiskConfigRiskPart(market_risk_configs)
}

func (s *ServerEngine) SymbolParamsChanged(namespace, group, dataId, data string) {
	logx.Info(fmt.Sprintf("SymbolContent: %s\n", data))
	s.ProcessSymbolConfigStr(data)
}

func (s *ServerEngine) ProcessSymbolConfigStr(data string) {
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
func (s *ServerEngine) UpdateRiskConfigHedgePart(hedge_configs []*mkconfig.HedgingConfig) {

	s.RiskConfigMutex.Lock()

	for _, hedge_config := range hedge_configs {
		if _, ok := s.RiskCtrlConfigMaps[hedge_config.Symbol]; !ok {
			s.RiskCtrlConfigMaps[hedge_config.Symbol] = &mkconfig.RiskCtrlConfig{}
			s.RiskCtrlConfigMaps[hedge_config.Symbol].HedgeConfigMap = make(map[string]*mkconfig.HedgeConfig)

			logx.Info("Risk New Symbol: " + hedge_config.Symbol + "\n")
		}

		logx.Info("Risk Update Symbol: " + hedge_config.Symbol + "\n")

		if s.RiskCtrlConfigMaps[hedge_config.Symbol].HedgeConfigMap == nil {
			s.RiskCtrlConfigMaps[hedge_config.Symbol].HedgeConfigMap = make(map[string]*mkconfig.HedgeConfig)
		}

		s.RiskCtrlConfigMaps[hedge_config.Symbol].HedgeConfigMap[hedge_config.Exchange] = &mkconfig.HedgeConfig{
			FeeKind:  hedge_config.FeeKind,
			FeeValue: hedge_config.TakerFee,
		}
	}

	s.UpdateRiskConfig()

	s.RiskConfigMutex.Unlock()
}

/*
        "price_offset_kind":1,
        "price_offset":0.001,
*/
func (s *ServerEngine) UpdateRiskConfigRiskPart(market_risk_configs []*mkconfig.MarketRiskConfig) {

	s.RiskConfigMutex.Lock()

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
func (s *ServerEngine) UpdateRiskConfigSymbolPart(SymbolConfigs []*mkconfig.SymbolConfig) {

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

func (s *ServerEngine) UpdateRiskConfig() {
	if s.HedgingConfigs != nil && s.MarketRiskConfigs != nil && s.SymbolConfigs != nil {

		logx.Info(fmt.Sprintf("s.RiskCtrlConfigMaps: \n%s", GetRiskCtrlConfigMapString(&s.RiskCtrlConfigMaps)))

		s.Riskworker.UpdateConfig(&s.RiskCtrlConfigMaps)
	}
}

func (s *ServerEngine) SetTestConfig() {
	risk_config := GetTestRiskConfig()
	logx.Infof("\\nrisk_config: %+v", risk_config)

	s.Riskworker.UpdateConfig(&risk_config)

	AggConfig := GetTestAggConfig()
	logx.Infof("\\nAggConfig: %+v", AggConfig)

	s.AggregateWorker.UpdateConfig(AggConfig)

	meta_data := datastruct.GetTestMetadata("BTC_USDT")
	logx.Infof("\\nmeta_data: %+v", meta_data)

	s.Commer.UpdateMetaData(meta_data)
}

func (s *ServerEngine) TestKafkaCancelListen() {

	meta_data := datastruct.GetTestMetadata("ETH_USDT")
	logx.Infof("meta_data: %+v", meta_data)

	s.Commer.UpdateMetaData(meta_data)
}
