package main

import (
	"exterior-interactor/pkg/exchangeapi/exchanges/ftx"
	"exterior-interactor/pkg/exchangeapi/exmodel"
	"exterior-interactor/pkg/exchangeapi/extools"
	"fmt"
)

func main() {
	api := ftx.NewNativeApiWithProxy(exmodel.EmptyAccountConfig, "http://localhost:1081")
	s := ftx.NewSymbolManager(api)

	m := ftx.NewMarketTradeManager(s, api)

	k := extools.NewKlineGenerator()

	m.Sub("BTC_USDT")

	for item:=range m.OutputCh(){
		fmt.Printf("%+v \n",item)
	}

	go func() {
		for item := range k.OutputCh() {
			fmt.Printf("%+v \n", *item)
		}
	}()

	for item := range m.OutputCh() {
		k.InputMarketTrade(item)
	}

}