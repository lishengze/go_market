package main

import (
	"exterior-interactor/pkg/exchangeapi/exchanges/ftx"
	"exterior-interactor/pkg/exchangeapi/exchanges/ftx/ftxapi"
	"exterior-interactor/pkg/exchangeapi/exmodel"
	"fmt"
	"log"
)

func main() {
	api, _ := ftx.NewNativeApi(exmodel.AccountConfig{
		Proxy:          "http://localhost:1081",
		Alias:          "FTX_MCA_OTC_TRADING",
		Key:            "1IVOM9S2EELFcOZGJ3RLIg3X_ETRwOM4sDH9a_D5",
		Secret:         "A3US9cW3biFIO9xQYfRdTFpD4UX0gNCcTyzY2F5W",
		SubAccountName: "Xpert RFQ",
	})
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

	rsp, err := api.QueryOrderByClientOrderId("1538848152356917248")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%+v\n", *rsp)

	return

	d, err := api.GetStreamDepth()
	if err != nil {
		log.Fatalln(err)
	}
	d.Sub("BTC-PERP")

	for item := range d.ReadCh() {
		fmt.Println("asks:", item.(*ftxapi.StreamDepth).Data.Asks)
		fmt.Println("bids:", item.(*ftxapi.StreamDepth).Data.Bids)
	}
}
