package main

import (
	"exterior-interactor/pkg/exchangeapi/exchanges/binance"
	"exterior-interactor/pkg/exchangeapi/exchanges/binance/api/binancespot"
	"exterior-interactor/pkg/exchangeapi/exmodel"
	"fmt"
	"log"
	"time"
)

func main() {
	api, _ := binance.NewNativeApi(exmodel.AccountConfig{Proxy: "http://localhost:1081"})
	k, err := api.SpotApi.GetStreamKline()
	if err != nil {
		log.Fatalln(err)
	}

	k.Sub("btcusdt@kline_1m")

	for item := range k.ReadCh() {
		data := item.(*binancespot.WsKline)
		fmt.Println("*************")
		fmt.Println("startTime", time.UnixMilli(data.Data.K.T))
		fmt.Println("startTime", time.UnixMilli(data.Data.K.T1))
		fmt.Println("open", data.Data.K.O)
		fmt.Println("high", data.Data.K.H)
		fmt.Println("low", data.Data.K.L)
		fmt.Println("close", data.Data.K.C)
		fmt.Println("volume", data.Data.K.V)
		fmt.Println("value", data.Data.K.Q)
		fmt.Println("tradeAmount", data.Data.K.N)
		fmt.Println("startTradeId", data.Data.K.F)
		fmt.Println("endTradeId", data.Data.K.L1)
	}
}
