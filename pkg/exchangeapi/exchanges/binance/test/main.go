package main

import (
	"exterior-interactor/pkg/exchangeapi/exchanges/binance"
	"exterior-interactor/pkg/exchangeapi/exmodel"
	"fmt"
)

func main() {
	api, _ := binance.NewNativeApi(exmodel.AccountConfig{Proxy: "http://localhost:1081"})

	s := binance.NewSymbolManager(api)
	d := binance.NewDepthManager(s, api)
	//t := binance.NewMarketTradeManager(s, api)
	//k := extools.NewKlineGenerator()

	d.Sub("BTC_USDT")
	//d.Sub("SAND_USDT", "BTC_USDT")
	//t.Sub("BTC_USDT")
	//t.Sub("SAND_USDT", "BTC_USDT")

	//go func() {
	//	for item := range t.OutputCh() {
	//		k.InputMarketTrade(item)
	//	}
	//}()
	//
	//for item := range k.OutputCh() {
	//	fmt.Printf("%+v \n", *item)
	//}

	for item := range d.OutputCh() {
		fmt.Println(">>>>>>>>>>>>", item.Symbol.StdSymbol)
		fmt.Println("ask", item.Asks)
		fmt.Println("bid", item.Bids)
	}
}
