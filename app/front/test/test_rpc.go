package main

import (
	"context"
	"fmt"
	"market_server/app/dataManager/rpc/marketservice"
	"market_server/common/datastruct"

	"github.com/zeromicro/go-zero/core/conf"
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

func main() {
	fmt.Println("----- Test Rpc -------")

	test_rpc()
}