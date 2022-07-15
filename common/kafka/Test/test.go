package main

import (
	"market_server/common/comm"
	"market_server/common/config"
	"market_server/common/datastruct"
	"market_server/common/util"

	"market_server/common/kafka"

	"github.com/zeromicro/go-zero/core/logx"
)

func GetTestMetadata() *datastruct.Metadata {
	symbol_set := make(map[string](map[string]struct{}))
	exchange_set := make(map[string]struct{})
	exchange_set["FTX"] = struct{}{}
	symbol := "BTC_USDT"
	symbol_set[symbol] = exchange_set

	MetaData := datastruct.Metadata{}

	// MetaData.DepthMeta = symbol_set
	MetaData.TradeMeta = symbol_set

	return &MetaData
}

func GetInitMeta() *datastruct.Metadata {
	init_symbol_list := []string{"BTC_USDT"}

	symbol_exchange_set := make(map[string](map[string]struct{}))
	new_meta := datastruct.Metadata{}
	for _, symbol := range init_symbol_list {
		if _, ok := symbol_exchange_set[symbol]; !ok {
			symbol_exchange_set[symbol] = make(map[string]struct{})
		}
		if _, ok := symbol_exchange_set[symbol][datastruct.BCTS_EXCHANGE]; !ok {
			symbol_exchange_set[symbol][datastruct.BCTS_EXCHANGE] = struct{}{}
		}
	}

	// new_meta.TradeMeta = symbol_exchange_set
	// new_meta.KlineMeta = symbol_exchange_set

	new_meta.DepthMeta = symbol_exchange_set

	logx.Infof("[I] InitMeta: %v \n", new_meta)

	return &new_meta
}

func test_kafka() {

	util.InitTestLogx()

	config := config.CommConfig{
		KafkaConfig: config.KafkaConfig{
			IP: "10.8.147.73:8117",
		},
		NetServerType: "KAFKA",
		SerialType:    "PROTOBUF",
	}
	RecvDataChan := datastruct.NewDataChannel()

	PubDataChan := datastruct.NewDataChannel()

	Serializer := &comm.ProtobufSerializer{}
	kafka_server, err := kafka.NewKafka(Serializer, RecvDataChan, PubDataChan, config.KafkaConfig)
	if err != nil {
		logx.Errorf("NewKafka Error: %+v", err)
	}

	kafka_server.UpdateMetaData(GetInitMeta())
	kafka_server.Start()

	// logx.Infof("CreateTopic: %+v", kafka_server.CreatedTopics())

	// topic := "golang_test"
	// if kafka_server.CreateTopic(topic) {
	// 	logx.Infof("Create %s Successesfully!", topic)
	// } else {
	// 	logx.Infof("Create %s Failed!", topic)
	// }

	// kafka_server.UpdateCreateTopics()
	// logx.Infof("CreateTopic: %+v", kafka_server.CreatedTopics())

	select {}
}

func StartListenRecvdata(RecvDataChan *datastruct.DataChannel) {
	logx.Info("[S] DBServer start_listen_recvdata")
	go func() {
		for {
			select {
			case new_depth := <-RecvDataChan.DepthChannel:
				store_depth(new_depth)
			case new_kline := <-RecvDataChan.KlineChannel:
				store_kline(new_kline)
			case new_trade := <-RecvDataChan.TradeChannel:
				store_trade(new_trade)
			}
		}
	}()
	logx.Info("[S] DBServer start_receiver Over!")
}

func store_kline(kline *datastruct.Kline) error {

	return nil
}

func store_trade(trade *datastruct.Trade) error {

	return nil
}

func store_depth(depth *datastruct.DepthQuote) error {
	logx.Info(depth.Bids)
	return nil
}

func main() {
	test_kafka()
}
