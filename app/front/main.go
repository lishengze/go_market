package main

import (
	"flag"
	"fmt"
	fconfig "market_server/app/front/config"
	"market_server/app/front/server"
	"market_server/app/front/svc"
	"market_server/common/datastruct"
	"market_server/common/util"
	"os"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

func TestEngine() {
	flag.Parse()
	env := os.Args[1]

	fmt.Printf("env: %+v \n", env)
	var configFile = flag.String("f", "etc/"+env+"/client.yaml", "the config file")

	is_test := false
	if len(os.Args) > 2 {
		is_test = true
	}

	fmt.Println(*configFile)

	var c fconfig.Config
	conf.MustLoad(*configFile, &c)

	fmt.Printf("config: %+v \n", c)
	fmt.Printf("is_test: %+v \n", is_test)

	logx.MustSetup(c.LogConfig)

	logx.Infof("config: %+v \n", c)

	// return

	ctx := svc.NewServiceContext(c)

	svr := server.NewServerEngine(ctx)

	svr.SetTestFlag(is_test)

	svr.Start()

	if is_test {
		StartTest(svr)
	}

	select {}
}

func StartTest(svr *server.ServerEngine) {
	test_pub_map := make(map[string]struct{})
	test_pub_map[datastruct.KLINE_TYPE] = struct{}{}
	test_pub_map[datastruct.TRADE_TYPE] = struct{}{}
	test_pub_map[datastruct.DEPTH_TYPE] = struct{}{}

	logx.Statf("test_pub_map: %+v \n", test_pub_map)

	go StartPublishTestData(svr.RecvDataChan, test_pub_map)

	time.Sleep(time.Second * 4)

	test_sub_map := make(map[string]struct{})
	// test_sub_map[datastruct.KLINE_TYPE] = struct{}{}
	// test_sub_map[datastruct.TRADE_TYPE] = struct{}{}
	// test_sub_map[datastruct.DEPTH_TYPE] = struct{}{}

	logx.Statf("test_sub_map: %+v \n", test_sub_map)
	go StartSubTest(svr, test_sub_map)
}

func StartSubTest(svr *server.ServerEngine, test_map map[string]struct{}) {

	svr.FrontEngineWorker.TestSub(test_map)
}

func PubDepthTestData(depth_chan chan *datastruct.DepthQuote) {
	timer := time.Tick(3 * time.Second)

	// symbols := []string{"BTC_USDT", "ETH_USDT", "USDT_USD"}

	symbols := []string{"BTC_USDT"}

	for {
		select {
		case <-timer:
			tmp_depth := datastruct.GetTestDepthMultiSymbols(symbols, datastruct.BCTS_EXCHANGE)
			logx.Slowf("tmp_depth: %s \n", tmp_depth.String(3))
			depth_chan <- tmp_depth
		}
	}
}

func PubTradeTestData(trade_chan chan *datastruct.Trade) {
	timer := time.Tick(3 * time.Second)

	// symbols := []string{"BTC_USDT", "ETH_USDT", "USDT_USD"}

	symbols := []string{"BTC_USDT"}

	for {
		select {
		case <-timer:
			tmp_trade := datastruct.GetTestTradeMultiSymbols(symbols, datastruct.BCTS_EXCHANGE)
			logx.Slowf("tmp_trade: %s \n", tmp_trade.String())
			trade_chan <- tmp_trade
		}
	}
}

func PubKlineTestData(kline_chan chan *datastruct.Kline) {
	timer := time.Tick(3 * time.Second)

	// symbols := []string{"BTC_USDT", "ETH_USDT", "USDT_USD"}

	symbols := []string{"BTC_USDT"}

	cur_time := util.TimeMinuteNanos()

	for {
		select {
		case <-timer:
			tmp_kline := datastruct.GetTestKlineMultiSymbols(symbols, datastruct.BCTS_EXCHANGE, cur_time)
			// logx.Slowf("tmp_kline: %s \n", tmp_kline.String())
			kline_chan <- tmp_kline
			cur_time += datastruct.SECS_PER_MIN * datastruct.NANO_PER_SECS
		}
	}
}

func StartPublishTestData(data_chan *datastruct.DataChannel, test_map map[string]struct{}) {
	if _, ok := test_map[datastruct.DEPTH_TYPE]; ok {
		go PubDepthTestData(data_chan.DepthChannel)
	}

	if _, ok := test_map[datastruct.TRADE_TYPE]; ok {
		go PubTradeTestData(data_chan.TradeChannel)
	}

	if _, ok := test_map[datastruct.KLINE_TYPE]; ok {
		go PubKlineTestData(data_chan.KlineChannel)
	}
}

func main() {
	fmt.Println("--------------- This is Main -----------")
	TestEngine()
}
