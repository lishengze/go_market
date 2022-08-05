package main

import (
	"context"
	"fmt"
	"market_server/app/data_manager/rpc/marketservice"
	"market_server/common/datastruct"
	"market_server/common/util"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

func test_req_hist_kline(ctx context.Context, msclient marketservice.MarketService) {
	req_hist_info := &marketservice.ReqHishKlineInfo{
		Symbol:    "BTC_USDT",
		Exchange:  datastruct.BCTS_EXCHANGE,
		StartTime: 1654297007842658763,
		EndTime:   1654297013959596689,
		Count:     5,
		Frequency: 60,
	}

	rst, err := msclient.RequestHistKlineData(ctx, req_hist_info)

	if err != nil {
		fmt.Printf("err %+v \n", err)
	}

	fmt.Printf("Rst: %+v \n", rst)
}

func test_req_trade(ctx context.Context, msclient marketservice.MarketService) {
	in := &marketservice.ReqTradeInfo{
		Symbol:   "BTC_USDT",
		Exchange: datastruct.BCTS_EXCHANGE,
		Time:     1654373953705096962,
	}

	rst, err := msclient.RequestTradeData(ctx, in)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("ReqTrade Rst: %+v \n", rst)
}

func test_rpc() {
	zconfig := zrpc.RpcClientConf{}
	conf.MustLoad("front.yaml", &zconfig)

	zclient := zrpc.MustNewClient(zconfig)
	msclient := marketservice.NewMarketService(zclient)

	ctx := context.Background()

	test_req_hist_kline(ctx, msclient)

	test_req_trade(ctx, msclient)
}

type TestRpc struct {
	KlineCache datastruct.KlineCache
	ZClient    zrpc.Client
	MSClient   marketservice.MarketService
	Ctx        context.Context
}

func NewTestRpc() *TestRpc {

	util.InitTestLogx()

	cacheConfig := &datastruct.CacheConfig{
		Count: 1000,
	}

	zconfig := zrpc.RpcClientConf{}
	conf.MustLoad("front.yaml", &zconfig)
	zclient := zrpc.MustNewClient(zconfig)

	return &TestRpc{
		KlineCache: *datastruct.NewKlineCache(cacheConfig),
		ZClient:    zclient,
		MSClient:   marketservice.NewMarketService(zclient),
		Ctx:        context.Background(),
	}
}

func (t *TestRpc) Start() {
	t.TestKline()
	t.TestTrade()
}

func (t *TestRpc) TestKline() {
	resolution := 5 * datastruct.NANO_PER_MIN
	symbol := "BTC_USDT"

	req_hist_info := &marketservice.ReqHishKlineInfo{
		Symbol:    symbol,
		Exchange:  datastruct.BCTS_EXCHANGE,
		StartTime: 1654297007842658763,
		EndTime:   1654297013959596689,
		Count:     30,
		Frequency: uint32(resolution),
	}

	rst, err := t.MSClient.RequestHistKlineData(t.Ctx, req_hist_info)

	if err != nil {
		fmt.Printf("err %+v \n", err)
	}

	klines := marketservice.TransPbKlines(rst.KlineData)

	logx.Infof("ori data_count: %d", len(klines))
	for _, kline := range klines {
		logx.Infof(kline.FullString())
	}

	t.KlineCache.ReleaseInputKlines(klines, symbol, resolution)

	logx.Slowf("KlineCache: %s", t.KlineCache.String(symbol, resolution))
}

func (t *TestRpc) TestTrade() {
	in := &marketservice.ReqTradeInfo{
		Symbol:   "BTC_USDT",
		Exchange: datastruct.BCTS_EXCHANGE,
		Time:     1654373953705096962,
	}

	rst, err := t.MSClient.RequestTradeData(t.Ctx, in)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("ReqTrade Rst: %+v \n", rst)
}

func main() {
	fmt.Println("----- Test Rpc -------")

	test_rpc()
}
