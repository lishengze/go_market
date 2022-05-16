package main

import (
	"fmt"
	"market_aggregate/pkg/comm"
	"market_aggregate/pkg/conf"
	"market_aggregate/pkg/datastruct"
	"market_aggregate/pkg/kafka"
	"market_aggregate/pkg/util"
	"time"
)

func ListenRecvData(RecvDataChan *datastruct.DataChannel) {
	util.LOG_INFO("ListenRecvData")
	for {
		select {
		case local_depth := <-RecvDataChan.DepthChannel:
			fmt.Printf("local_depth: %+v\n\n", local_depth)
		case local_kline := <-RecvDataChan.KlineChannel:
			fmt.Printf("local_kline: %+v\n\n", local_kline)
		case local_trade := <-RecvDataChan.TradeChannel:
			fmt.Printf("local_trade: %+v\n\n", local_trade)
		}
	}
}

func TestKafka() {
	Config := &conf.Config{
		IP:            "43.154.179.47:9117",
		NetServerType: "KAFKA",
		SerialType:    "PROTOBUF",
	}

	symbol_set := make(map[string](map[string]struct{}))
	exchange_set := make(map[string]struct{})
	exchange_set["FTX"] = struct{}{}
	symbol_set["BTC_USDT"] = exchange_set

	MetaData := datastruct.Metadata{}

	// MetaData.DepthMeta = symbol_set

	MetaData.TradeMeta = symbol_set

	Serializer := &comm.ProtobufSerializer{}

	kafka_server := kafka.KafkaServer{}

	RecvDataChan := &datastruct.DataChannel{
		DepthChannel: make(chan *datastruct.DepthQuote),
		KlineChannel: make(chan *datastruct.Kline),
		TradeChannel: make(chan *datastruct.Trade),
	}

	PubDataChan := &datastruct.DataChannel{
		DepthChannel: make(chan *datastruct.DepthQuote),
		KlineChannel: make(chan *datastruct.Kline),
		TradeChannel: make(chan *datastruct.Trade),
	}

	go ListenRecvData(RecvDataChan)

	kafka_server.Init(Config, Serializer, RecvDataChan, PubDataChan, MetaData)

	kafka_server.IsTest = false

	kafka_server.Start()

	time.Sleep(time.Hour)
}

func main() {
	fmt.Println("---- Test Kafka ----")

	// aggregate.TestInnerDepth()

	// aggregate.TestImport()

	// aggregate.TestWorker()

	// TestDepthChannel()

	// TestTreeMap()

	// comm.TestSeKline()

	TestKafka()
}
