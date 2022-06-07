package aggregate

import (
	"fmt"
	"market_server/app/market_aggregate/config"
	mkconfig "market_server/app/market_aggregate/config"
	"market_server/common/datastruct"
	"market_server/common/util"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

func GetTestRiskConfig() RiskCtrlConfigMap {
	rst := RiskCtrlConfigMap{
		"BTC_USDT": {
			HedgeConfigMap: map[string]*config.HedgeConfig{"FTX": &config.HedgeConfig{FeeKind: 1, FeeValue: 0.1},
				"OKEX":  {FeeKind: 1, FeeValue: 0.2},
				"HUOBI": {FeeKind: 1, FeeValue: 0.3}},
			PricePrecison:    2,
			VolumePrecison:   3,
			PriceBiasValue:   0.1,
			PriceBiasKind:    1,
			VolumeBiasValue:  0.1,
			VolumeBiasKind:   1,
			PriceMinumChange: 1.0,
		},
		"ETH_USDT": {
			HedgeConfigMap:   map[string]*config.HedgeConfig{"FTX": {FeeKind: 1, FeeValue: 0.1}, "OKEX": {FeeKind: 1, FeeValue: 0.2}, "HUOBI": {FeeKind: 1, FeeValue: 0.3}},
			PricePrecison:    2,
			VolumePrecison:   3,
			PriceBiasValue:   0.1,
			PriceBiasKind:    1,
			VolumeBiasValue:  0.1,
			VolumeBiasKind:   1,
			PriceMinumChange: 1.0,
		},
		"DOT_USDT": {
			HedgeConfigMap:   map[string]*config.HedgeConfig{"FTX": {FeeKind: 1, FeeValue: 0.1}, "OKEX": {FeeKind: 1, FeeValue: 0.2}, "HUOBI": {FeeKind: 1, FeeValue: 0.3}},
			PricePrecison:    2,
			VolumePrecison:   3,
			PriceBiasValue:   0.1,
			PriceBiasKind:    1,
			VolumeBiasValue:  0.1,
			VolumeBiasKind:   1,
			PriceMinumChange: 1.0,
		},
	}

	return rst
}

func TestWorker() {
	depth_quote := datastruct.GetTestDepth()
	config := GetTestRiskConfig()

	risk_worker_manager := RiskWorkerManager{}
	// risk_worker_manager.Init()

	risk_worker_manager.UpdateConfig(&config)

	// logx.Infof("depth_quote: %v\n", depth_quote)
	// logx.Infof("config: %v\n", config)

	risk_worker_manager.Execute(depth_quote)
}

func TestInnerDepth() {
	// a := datastruct.InnerDepth{0, make(map[string]float64)}

	// var e1 = map[string]float64{
	// 	"FTX": 1.1,
	// }

	// e := map[string]float64{
	// 	"FTX": 1.1,
	// }

	// a := datastruct.InnerDepth{0, map[string]float64{"FTX": 1.1}}
}

func test_json() {
	depth_quote := datastruct.GetTestDepth()
	fmt.Println(depth_quote.String(3))
}

func TestImport() {
	data := datastruct.TestData{
		Name: "Tom",
	}

	fmt.Println(data)
}

func TestAddWorker() {
	// risk_work := RiskWorkerManager{}
	// risk_work.Init()
}

// func main() {
// 	test_worker()

// 	// test_get_sorted_keys()

// 	// test_inner_depth()

// 	// test_json()
// }

func StartRecvDataChannel(RecvDataChan *datastruct.DataChannel) {
	timer := time.Tick(1 * time.Second)

	// index := 0
	for {
		select {
		case <-timer:
			// depth_quote := datastruct.GetTestDepth()
			// index++
			// RecvDataChan.DepthChannel <- depth_quote
			RecvDataChan.TradeChannel <- datastruct.GetTestTrade()
		}
	}
}

func StartPubDataChannel(PubDataChan *datastruct.DataChannel) {
	for {
		select {
		case depth := <-PubDataChan.DepthChannel:
			logx.Info(fmt.Sprintf("\n~~~~~~~~Processed Depth %+v \n", depth))
		case trade := <-PubDataChan.TradeChannel:
			logx.Info(fmt.Sprintf("\n~~~~~~~~Processed Trade %+v \n", trade))
		case kline := <-PubDataChan.KlineChannel:
			logx.Info(fmt.Sprintf("\n~~~~~~~~Processed Kline %+v \n", kline))
		}
	}
}

func GetTestAggConfig() mkconfig.AggregateConfig {
	return mkconfig.AggregateConfig{
		DepthAggregatorConfigMap: map[string]mkconfig.AggregateConfigAtom{"BTC_USDT": mkconfig.AggregateConfigAtom{4000, 30, true}, "ETH_USDT": mkconfig.AggregateConfigAtom{6000, 40, true}},
	}
}

func TestAggregator() {
	aggregator := Aggregator{}

	RecvDataChan := &datastruct.DataChannel{
		DepthChannel: make(chan *datastruct.DepthQuote),
		KlineChannel: make(chan *datastruct.Kline),
		TradeChannel: make(chan *datastruct.Trade),
	}

	PubDataChan := &datastruct.DataChannel{
		DepthChannel: make(chan *datastruct.DepthQuote),
		KlineChannel: make(chan *datastruct.Kline),
		TradeChannel: make(chan *datastruct.Trade),
	}

	risk_config := GetTestRiskConfig()
	risk_worker := &RiskWorkerManager{}
	// risk_worker.Init()
	risk_worker.UpdateConfig(&risk_config)

	AggConfig := GetTestAggConfig()
	// aggregator.Init(RecvDataChan, PubDataChan, risk_worker)
	aggregator.UpdateConfig(AggConfig)
	aggregator.Start()

	logx.Info(fmt.Sprintf("\n------- risk_worker.RiskConfig: %+v\n\n", risk_worker.RiskConfig))

	go StartRecvDataChannel(RecvDataChan)

	go StartPubDataChannel(PubDataChan)

	time.Sleep(time.Hour)
}

func write_log1() {
	for {
		logx.Info("[1] " + util.CurrTimeString())
		time.Sleep(time.Second * 1)
	}
}

func write_log2() {
	for {
		// Info("[2] " + util.CurrTimeString())

		logx.Infof("f[2] %s ", util.CurrTimeString())
		time.Sleep(time.Second * 1)
	}
}

func write_log3() {
	for {
		logx.Errorf("f[3] %s", util.CurrTimeString())
		time.Sleep(time.Second * 1)
	}
}

func write_log4() {
	for {
		logx.Slowf("f[4] %s", util.CurrTimeString())
		time.Sleep(time.Second * 1)
	}
}

func write_log5() {
	for {
		logx.Severef("f[5] %s", util.CurrTimeString())
		time.Sleep(time.Second * 1)
	}
}

func write_log6() {
	for {
		logx.Statf("f[6] %s", util.CurrTimeString())
		time.Sleep(time.Second * 1)
	}
}

func TestLog() {
	go write_log1()

	go write_log2()

	go write_log3()

	go write_log4()

	go write_log5()

	go write_log6()

	select {}
}
