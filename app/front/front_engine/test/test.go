package main

import (
	"encoding/json"
	"fmt"
	"market_server/app/front/front_engine"
	"market_server/common/datastruct"

	"github.com/zeromicro/go-zero/core/logx"
)

func InitLogx() {

	LogConfig := logx.LogConf{
		Compress:            true,
		KeepDays:            0,
		Level:               "info",
		Mode:                "file",
		Path:                "./log",
		ServiceName:         "client",
		StackCooldownMillis: 100,
		TimeFormat:          "2006-01-02 15:04:05",
	}

	logx.MustSetup(LogConfig)
}

func TestPubSymbolData() {
	symbol_list := []string{"BTC_USDT", "ETH_USDT"}

	fmt.Printf("Original Data: %+v \n", symbol_list)

	rst := front_engine.NewSymbolListMsg(symbol_list)

	fmt.Printf("\nJson string: %s \n", string(rst))

	var trans_data front_engine.PubSymbolistJson
	err := json.Unmarshal(rst, &trans_data)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("transed Data: %+v \n", trans_data)
	}
}

func TestPubDepthData() {
	original_depth := datastruct.GetTestDepth()

	fmt.Printf("Original Data: %+v \n", original_depth)

	rst := front_engine.NewDepthJsonMsg(original_depth)

	fmt.Printf("\nJson string: %s \n", string(rst))

	var trans_depth front_engine.PubDepthJson
	err := json.Unmarshal(rst, &trans_depth)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("transed Data: %+v \n", trans_depth)
	}

}

func TestPubHistKlineData() {
	req_kline_info := &datastruct.ReqHistKline{
		Symbol:    "BTC_USDT",
		Exchange:  "FTX",
		StartTime: 0,
		EndTime:   0,
		Count:     2,
		Frequency: 60,
	}
	original_histkline := datastruct.GetTestHistKline(req_kline_info)

	fmt.Printf("Original Data: %+v \n", original_histkline)

	rst := front_engine.NewHistKlineJsonMsg(&datastruct.RspHistKline{
		ReqInfo: req_kline_info,
		Klines:  original_histkline,
	})

	fmt.Printf("\nJson string: %s \n", string(rst))

	var trans_data front_engine.PubKlineJson
	err := json.Unmarshal(rst, &trans_data)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("\n transed Data: %+v \n", trans_data)
	}

}

func TestPubKlineData() {
	original_kline := datastruct.GetTestKline()

	fmt.Printf("Original Data: %+v \n", original_kline)

	rst := front_engine.NewKlineUpdateJsonMsg(original_kline)

	fmt.Printf("\nJson string: %s \n", string(rst))

	var trans_kline front_engine.PubKlineJson
	err := json.Unmarshal(rst, &trans_kline)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("\ntransed Data: %+v \n", trans_kline)
	}

}

func TestPubTradeData() {
	trade := datastruct.GetTestTrade()

	fmt.Printf("Original Data: %+v \n", trade)

	rst := front_engine.NewTradeJsonMsg(trade, nil, 0.0)

	fmt.Printf("\nJson string: %s \n", string(rst))

	var trans_trade front_engine.PubTradeJson
	err := json.Unmarshal(rst, &trans_trade)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("\ntransed Data: %+v \n", trans_trade)
	}

}

func main() {
	InitLogx()

	TestPubSymbolData()

	// TestPubTradeData()

	// TestPubDepthData()

	// TestPubKlineData()

	// TestPubHistKlineData()
}
