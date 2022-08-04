package test

import (
	"fmt"
	"market_server/common/comm"
	"market_server/common/config"
	"market_server/common/datastruct"
	"market_server/common/kafka"
	"market_server/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type TestKlineCache struct {
	RecvDataChan *datastruct.DataChannel
	PubDataChan  *datastruct.DataChannel
	KlineCache   *datastruct.KlineCache
	KafkaServer  *kafka.KafkaServer
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
	new_meta.KlineMeta = symbol_exchange_set
	// new_meta.DepthMeta = symbol_exchange_set

	logx.Infof("[I] InitMeta: %v \n", new_meta)

	return &new_meta
}

func NewTestKlineCache() *TestKlineCache {
	util.InitTestLogx()

	config := config.CommConfig{
		KafkaConfig: config.KafkaConfig{
			IP: "10.8.147.73:8117",
		},
		NetServerType: "KAFKA",
		SerialType:    "PROTOBUF",
	}
	recvDataChan := datastruct.NewDataChannel()

	pubDataChan := datastruct.NewDataChannel()

	serializer := &comm.ProtobufSerializer{}
	kafkaServer, err := kafka.NewKafka(serializer, recvDataChan, pubDataChan, config.KafkaConfig)

	cacheConfig := &datastruct.CacheConfig{
		Count: 1500,
	}
	klineCache := datastruct.NewKlineCache(cacheConfig)

	if err != nil {
		logx.Errorf("NewKafka Error: %+v", err)
	}

	return &TestKlineCache{
		RecvDataChan: recvDataChan,
		PubDataChan:  pubDataChan,
		KafkaServer:  kafkaServer,
		KlineCache:   klineCache,
	}

	// kafka_server.UpdateMetaData(GetInitMeta())
	// kafka_server.Start()
	// StartListenRecvdata(RecvDataChan)
}

func (t *TestKlineCache) StartListenRecvdata() {
	logx.Info("[S] TestKlineCache start_listen_recvdata")
	go func() {
		for {
			select {
			case new_depth := <-t.RecvDataChan.DepthChannel:
				t.process_depth(new_depth)
			case new_kline := <-t.RecvDataChan.KlineChannel:
				go t.process_kline(new_kline)
			case new_trade := <-t.RecvDataChan.TradeChannel:
				t.process_trade(new_trade)
			}
		}
	}()
	logx.Info("[S] TestKlineCache start_receiver Over!")
}

func (t *TestKlineCache) process_kline(kline *datastruct.Kline) error {
	defer util.CatchExp(fmt.Sprintf("ProcessKline %s", kline.FullString()))

	if kline == nil {
		return fmt.Errorf("kline is Nil")
	}

	resolution := 5 * datastruct.NANO_PER_MIN
	t.KlineCache.UpdateWithKline(kline, resolution)

	// if kline.IsHistory() {
	// 	logx.Slowf("[HK] %s", kline.FullString())
	// } else {
	// 	logx.Slowf("[RK] %s", kline.FullString())

	// 	trade := datastruct.NewTradeWithRealTimeKline(kline)
	// 	logx.Slowf("[RT] %s", trade.String())
	// }

	return nil
}

func (t *TestKlineCache) process_trade(trade *datastruct.Trade) error {

	return nil
}

func (t *TestKlineCache) process_depth(depth *datastruct.DepthQuote) error {
	// logx.Info(depth.Bids)
	return nil
}

func (t *TestKlineCache) Start() {
	logx.Infof("TestKlineCache Start!")

	t.KafkaServer.UpdateMetaData(GetInitMeta())
	t.KafkaServer.Start()
	t.StartListenRecvdata()
}

func TestKlineCacheMain() {
	t := NewTestKlineCache()
	t.Start()

	select {}
}
