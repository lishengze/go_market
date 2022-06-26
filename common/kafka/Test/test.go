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

func test_kafka() {

	util.InitTestLogx()

	config := config.CommConfig{
		KafkaConfig: config.KafkaConfig{
			IP: "43.154.179.47:9117",
		},
		NetServerType: "KAFKA",
		SerialType:    "PROTOBUF",
	}
	RecvDataChan := datastruct.NewDataChannel()

	PubDataChan := datastruct.NewDataChannel()

	// Commer := comm.NewComm(RecvDataChan, PubDataChan, config)

	// Commer.Start()

	// test_meta := GetTestMetadata()

	// Commer.UpdateMetaData(test_meta)

	Serializer := &comm.ProtobufSerializer{}
	kafka_server, err := kafka.NewKafka(Serializer, RecvDataChan, PubDataChan, config.KafkaConfig)
	if err != nil {
		logx.Errorf("NewKafka Error: %+v", err)
	}
	kafka_server.Start()

	logx.Infof("CreateTopic: %+v", kafka_server.CreatedTopics())

	topic := "golang_test"
	if kafka_server.CreateTopic(topic) {
		logx.Infof("Create %s Successesfully!", topic)
	} else {
		logx.Infof("Create %s Failed!", topic)
	}

	kafka_server.UpdateCreateTopics()
	logx.Infof("CreateTopic: %+v", kafka_server.CreatedTopics())

	select {}
}

func main() {
	test_kafka()
}
