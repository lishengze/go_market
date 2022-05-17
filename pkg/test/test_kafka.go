package main

import (
	"fmt"
	"market_aggregate/pkg/comm"
	config "market_aggregate/pkg/conf"
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

func StartPubData(PubDataChan *datastruct.DataChannel) {
	timer := time.Tick(3 * time.Second)

	// index := 0
	for {
		select {
		case <-timer:
			depth_quote := datastruct.GetTestDepth()
			depth_quote.Exchange = datastruct.BCTS_EXCHANGE
			depth_quote.Symbol = "BTC_USDT"
			PubDataChan.DepthChannel <- depth_quote

			local_kline := datastruct.GetTestKline()
			local_kline.Exchange = datastruct.BCTS_EXCHANGE
			local_kline.Symbol = "BTC_USDT"
			PubDataChan.KlineChannel <- local_kline

			local_trade := datastruct.GetTestTrade()
			local_trade.Exchange = datastruct.BCTS_EXCHANGE
			local_trade.Symbol = "BTC_USDT"
			PubDataChan.TradeChannel <- local_trade
		}
	}
}

func TestKafka() {
	config.NATIVE_CONFIG_INIT("client.yaml")

	util.LOG_INFO(fmt.Sprintf("CONFIG: %+v", *config.NATIVE_CONFIG()))

	// return

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

	kafka_server.Init(Serializer, RecvDataChan, PubDataChan)

	kafka_server.IsTest = false

	kafka_server.Start()

	go StartPubData(PubDataChan)
	// kafka_server.UpdateMetaData(MetaData)

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

	// aggregate.TestAggregator()

	// aggregate.TestAddWorker()

	// config.TestConf()
}
