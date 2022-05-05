package main

import (
	"exterior-interactor/pkg/exchangeapi/exchanges/ftx"
	"exterior-interactor/pkg/exchangeapi/exchanges/ftx/ftxapi"
	"exterior-interactor/pkg/exchangeapi/exmodel"
	"fmt"
	"log"
)

func main() {
	api := ftx.NewNativeApiWithProxy(exmodel.EmptyAccountConfig, "http://localhost:1081")
	//rsp, err := api.GetMarket()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//fmt.Printf("%+v \n",rsp.Result[0])
	//return
	//
	//m, err := api.GetStreamMarketTrade()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//m.Sub("BTC-PERP")
	//
	//for item := range m.ReadCh() {
	//	fmt.Println(item)
	//}

	d, err := api.GetStreamDepth()
	if err != nil {
		log.Fatalln(err)
	}
	d.Sub("BTC-PERP")

	for item := range d.ReadCh() {
		fmt.Println("asks:",item.(*ftxapi.StreamDepth).Data.Asks)
		fmt.Println("bids:",item.(*ftxapi.StreamDepth).Data.Bids)
	}
}
