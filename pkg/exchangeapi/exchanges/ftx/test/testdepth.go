package main

import (
	"exterior-interactor/pkg/exchangeapi/exchanges/ftx"
	"exterior-interactor/pkg/exchangeapi/exmodel"
)

func main() {
	api, _ := ftx.NewNativeApi(exmodel.AccountConfig{Proxy: "http://localhost:1081"})
	s := ftx.NewSymbolManager(api)

	d := ftx.NewDepthManager(s, api)

	d.Sub("XRP_USD")
	d.Sub("ETH_USD")
	d.Sub("BTC_USD")
	d.Sub("BTC_USDT")
	d.Sub("ETH_USDT")

	//for item:=range d.OutputCh(){
	//	fmt.Printf("symbol %+v \n",item.Symbol.StdSymbol)
	//	fmt.Printf("ask %+v \n",item.Asks)
	//	fmt.Printf("bid %+v \n",item.Bids)
	//}

	for range d.OutputCh() {

	}

}
