package main

import (
	"fmt"
	"log"
	"market_server/pkg/kafkaclient/mpupb"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	a := mpupb.Trade{
		Timestamp: timestamppb.Now(),
		Exchange:  "BINANCE",
		Symbol:    "BTC_USDT",
		Price:     "1000000",
		Volume:    "1",
	}
	fmt.Println(a)

	msg, err := proto.Marshal(&a)
	if err != nil {
		log.Fatalln(err, "----------1")
	}

	a2 := mpupb.Trade{}
	err = proto.Unmarshal(msg, &a2)
	if err != nil {
		log.Fatalln(err, "-------------2")
	}

	fmt.Println(a2)
}
