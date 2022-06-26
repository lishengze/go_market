package main

import (
	"market_server/common/comm"
	"market_server/common/config"
	"market_server/common/datastruct"
	"market_server/common/util"
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

	Commer := comm.NewComm(RecvDataChan, PubDataChan, config)

	Commer.Start()

	test_meta := GetTestMetadata()

	Commer.UpdateMetaData(test_meta)
}

func main() {
	// kafka.TestConsumer()
}
