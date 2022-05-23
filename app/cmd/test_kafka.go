package main

import (
	"fmt"
	config "market_aggregate/app/conf"
	"market_aggregate/app/datastruct"
	"market_aggregate/pkg/comm"
	"market_aggregate/pkg/kafka"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

func ListenRecvData(RecvDataChan *datastruct.DataChannel) {
	logx.Info("ListenRecvData")
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

	logx.Info(fmt.Sprintf("CONFIG: %+v", *config.NATIVE_CONFIG()))

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

	time.Sleep(time.Hour)
}

// func main() {
// 	fmt.Println("---- Test Kafka ----")

// 	// aggregate.TestInnerDepth()

// 	// aggregate.TestImport()

// 	// aggregate.TestWorker()

// 	// TestDepthChannel()

// 	// TestTreeMap()

// 	// comm.TestSeKline()

// 	// TestKafka()

// 	// aggregate.TestAggregator()

// 	// aggregate.TestAddWorker()

// 	// config.TestConf()

// 	aggregate.TestServerEngine()

// 	// util.TestTimeStr()

// 	// ulog.TestLog()

// 	// aggregate.TestLog()
// }
