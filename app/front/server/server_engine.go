package server

import (
	"fmt"
	"market_server/app/front/data_engine"
	"market_server/app/front/front_engine"
	"market_server/app/front/svc"
	"market_server/app/front/ws_server"
	mkconfig "market_server/app/market_aggregate/config"
	"market_server/common/comm"
	config "market_server/common/config"
	"market_server/common/datastruct"

	"github.com/zeromicro/go-zero/core/logx"
)

type ServerEngine struct {
	ctx    *svc.ServiceContext
	Commer *comm.Comm

	RecvDataChan *datastruct.DataChannel
	PubDataChan  *datastruct.DataChannel

	HedgingConfigs    []*mkconfig.HedgingConfig
	NacosClientWorker *config.NacosClient

	FrontEngineWorker *front_engine.FrontEngine
	DataEngineWorker  *data_engine.DataEngine
	WSEngineWorker    *ws_server.WSEngine

	IsTest bool
}

func NewServerEngine(svcCtx *svc.ServiceContext) *ServerEngine {

	s := &ServerEngine{}
	s.ctx = svcCtx
	s.RecvDataChan = datastruct.NewDataChannel()

	s.PubDataChan = datastruct.NewDataChannel()

	s.FrontEngineWorker = front_engine.NewFrontEngine(&svcCtx.Config)
	s.DataEngineWorker = data_engine.NewDataEngine(s.RecvDataChan, svcCtx)
	s.WSEngineWorker = ws_server.NewWSEngine(&svcCtx.Config.WS)

	s.FrontEngineWorker.SetNextWorker(s.DataEngineWorker)
	s.DataEngineWorker.SetNextWorker(s.FrontEngineWorker)
	s.WSEngineWorker.SetNextWorker(s.FrontEngineWorker)

	s.HedgingConfigs = nil
	s.IsTest = false
	s.Commer = comm.NewComm(s.RecvDataChan, s.PubDataChan, svcCtx.Config.Comm)

	return s
}

func (s *ServerEngine) SetTestFlag(value bool) {
	s.IsTest = value
	s.DataEngineWorker.SetTestFlag(value)
	s.FrontEngineWorker.SetTestFlag(value)
}

func (s *ServerEngine) StartNacosClient() {
	logx.Info("****************** StartNacosClient *****************")
	s.NacosClientWorker = config.NewNacosClient(&s.ctx.Config.Nacos)

	logx.Info("Connect Nacos Successfully!")

	HedgeConfigStr, err := s.NacosClientWorker.GetConfigContent("HedgeParams", datastruct.BCTS_GROUP)
	if err != nil {
		logx.Error(err.Error())
	}
	logx.Slow("Requested HedgeConfigStr: " + HedgeConfigStr)
	s.ProcsssHedgeConfigStr(HedgeConfigStr)

	s.NacosClientWorker.ListenConfig("HedgeParams", datastruct.BCTS_GROUP, s.HedgeParamsChanged)

	SymbolConfigStr, err := s.NacosClientWorker.GetConfigContent("SymbolParams", datastruct.BCTS_GROUP)
	if err != nil {
		logx.Error(err.Error())
	}
	logx.Info("Requested SymbolConfigStr: " + SymbolConfigStr)
	s.ProcessSymbolConfigStr(SymbolConfigStr)

	s.NacosClientWorker.ListenConfig("SymbolParams", datastruct.BCTS_GROUP, s.SymbolParamsChanged)
}

func (s *ServerEngine) HedgeParamsChanged(namespace, group, dataId, hedgingContent string) {
	logx.Slow(fmt.Sprintf("HedgeParamsChanged hedgingContent: %s\n", hedgingContent))
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

		if _, ok := symbol_exchange_set[hedge_config.Symbol][datastruct.BCTS_EXCHANGE]; !ok {
			symbol_exchange_set[hedge_config.Symbol][datastruct.BCTS_EXCHANGE] = struct{}{}
			logx.Slow(fmt.Sprintf("New Meta: %s.%s", hedge_config.Symbol, hedge_config.Exchange))
		}
	}

	NewMeta.DepthMeta = symbol_exchange_set
	NewMeta.TradeMeta = symbol_exchange_set
	NewMeta.KlineMeta = symbol_exchange_set

	logx.Info(fmt.Sprintf("HedgeParamsChanged: NewMeta:\n%s \n", NewMeta.String()))

	s.Commer.UpdateMetaData(&NewMeta)
}

func (s *ServerEngine) SymbolParamsChanged(namespace, group, dataId, data string) {
	logx.Infof("SymbolContent: %s\n", data)
	s.ProcessSymbolConfigStr(data)
}

func (s *ServerEngine) ProcessSymbolConfigStr(data string) {
	symbol_configs, err := mkconfig.ParseJsonSymbolConfig(data)

	if err != nil {
		logx.Error(err.Error())
		return
	}

	s.ctx.UpdateSymbolConfigWithSlice(symbol_configs)
}

func (s *ServerEngine) SetTestConfig() {

	symbols := []string{"BTC_USDT", "ETH_USDT"}

	meta_data := datastruct.GetTestMetadata(symbols)
	logx.Infof("\\meta_data: %+v", meta_data)

	// s.Commer.UpdateMetaData(meta_data)
}

func (s *ServerEngine) Start() {
	go s.Commer.Start()
	go s.DataEngineWorker.Start()
	go s.FrontEngineWorker.Start()
	go s.WSEngineWorker.Start()

	if !s.IsTest {
		go s.StartNacosClient()
	}
}
