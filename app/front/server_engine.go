package main

import (
	"flag"
	"fmt"
	fconfig "market_server/app/front/config"
	"market_server/app/front/data_engine"
	"market_server/app/front/front_engine"
	"market_server/app/front/svc"
	mkconfig "market_server/app/market_aggregate/config"
	"market_server/common/comm"
	config "market_server/common/config"
	"market_server/common/datastruct"
	"os"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

type ServerEngine struct {
	ctx    *svc.ServiceContext
	Commer *comm.Comm

	RecvDataChan *datastruct.DataChannel
	PubDataChan  *datastruct.DataChannel

	HedgingConfigs    []*mkconfig.HedgingConfig
	NacosClientWorker *config.NacosClient

	FrontEngine *front_engine.FrontEngine
	DataEngine  *data_engine.DataEngine

	IsTest bool
}

func NewServerEngine(svcCtx *svc.ServiceContext) *ServerEngine {

	s := &ServerEngine{}
	s.ctx = svcCtx
	s.RecvDataChan = datastruct.NewDataChannel()

	s.PubDataChan = datastruct.NewDataChannel()

	s.HedgingConfigs = nil
	s.IsTest = false
	s.Commer = comm.NewComm(s.RecvDataChan, s.PubDataChan, svcCtx.Config.Comm)

	return s
}

func (s *ServerEngine) SetTestFlag(value bool) {
	s.IsTest = value
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

		symbol_exchange_set[hedge_config.Symbol][hedge_config.Exchange] = struct{}{}
		logx.Slow(fmt.Sprintf("New Meta: %s.%s", hedge_config.Symbol, hedge_config.Exchange))
	}

	NewMeta.DepthMeta = symbol_exchange_set
	NewMeta.TradeMeta = symbol_exchange_set
	NewMeta.KlineMeta = symbol_exchange_set

	logx.Info(fmt.Sprintf("HedgeParamsChanged: NewMeta:\n%s \n", NewMeta.String()))

	s.Commer.UpdateMetaData(&NewMeta)
}

func (s *ServerEngine) SetTestConfig() {

	symbols := []string{"BTC_USDT", "ETH_USDT"}

	meta_data := datastruct.GetTestMetadata(symbols)
	logx.Infof("\\nmeta_data: %+v", meta_data)

	s.Commer.UpdateMetaData(meta_data)
}

func (s *ServerEngine) Start() {
	s.Commer.Start()

	if !s.IsTest {
		go s.StartNacosClient()
	} else {
		StartPublishTestData()
	}
}

func StartPublishTestData(data_chan *datastruct.DataChannel) {
	timer := time.Tick(1 * time.Second)

	// index := 0
	for {
		select {
		case <-timer:
			// depth_quote := datastruct.GetTestDepth()
			// index++
			// RecvDataChan.DepthChannel <- depth_quote
			data_chan.TradeChannel <- datastruct.GetTestTrade()
		}
	}
}

func TestEngine() {
	flag.Parse()

	env := "local"

	for _, v := range os.Args {
		env = v
	}

	fmt.Printf("env: %+v \n", env)
	var configFile = flag.String("f", "etc/"+env+"/marketData.yaml", "the config file")

	fmt.Println(*configFile)

	var c fconfig.Config
	conf.MustLoad(*configFile, &c)

	fmt.Printf("config: %+v \n", c)

	logx.MustSetup(c.LogConfig)

	logx.Infof("config: %+v \n", c)

	ctx := svc.NewServiceContext(c)
	svr := NewServerEngine(ctx)

	svr.SetTestFlag(true)
	svr.Start()

	time.Sleep(time.Second * 3)

	select {}
}

func StartSubTest() {

}
