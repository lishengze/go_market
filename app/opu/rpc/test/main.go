package main

import (
	"context"
	"exterior-interactor/app/opu/rpc/opu"
	"fmt"
	"github.com/zeromicro/go-zero/zrpc"
	"log"
)

const (
	opuAddress   = "152.32.254.76:8606"
	kafkaAddress = "152.32.254.76:9117"
)

func main() {
	testRegisterAccount()
}

func getOpu(address string) opu.Opu {
	client, err := zrpc.NewClientWithTarget(address)
	if err != nil {
		log.Fatalln(err)
	}
	return opu.NewOpu(client)
}

func testRegisterAccount() {
	o := getOpu(opuAddress)
	rsp, err := o.RegisterAccount(context.Background(), &opu.RegisterAccountReq{
		Alias:          "FTX_MCA_OTC_TRADING",
		Key:            "1IVOM9S2EELFcOZGJ3RLIg3X_ETRwOM4sDH9a_D5",
		Secret:         "A3US9cW3biFIO9xQYfRdTFpD4UX0gNCcTyzY2F5W",
		Passphrase:     "",
		Exchange:       "FTX",
		SubAccountName: "Xpert RFQ",
	})

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(rsp)
}

func testGetBalance() {
	o := getOpu(opuAddress)

	rsp, err := o.QueryBalance(context.Background(), &opu.QueryBalanceReq{
		AccountAlias: "FTX_MCA_OTC_TRADING",
	})

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(rsp)
}
