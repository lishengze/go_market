package kafka

import (
	"context"
	"fmt"
	"market_server/common/config"
	"market_server/common/datastruct"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/zeromicro/go-zero/core/logx"
)

type ConsumeItem struct {
	Topic      string
	Ctx        context.Context
	CancelFunc context.CancelFunc
}

type KafkaServer struct {
	// Consumer sarama.Consumer

	Producer sarama.SyncProducer
	Broker   *sarama.Broker

	MetaData datastruct.Metadata

	Serializer datastruct.SerializerI

	RecvDataChan *datastruct.DataChannel
	PubDataChan  *datastruct.DataChannel

	ConsumeSet map[string](*ConsumeItem)

	PublishMutex sync.Mutex

	config config.KafkaConfig
	IsTest bool
}

// Init(*config.Config, SerializerI, *DataChannel)
func (k *KafkaServer) InitKafka(serializer datastruct.SerializerI,
	recv_data_chan *datastruct.DataChannel,
	pub_data_chan *datastruct.DataChannel,
	config config.KafkaConfig) error {

	k.Serializer = serializer
	k.RecvDataChan = recv_data_chan
	k.PubDataChan = pub_data_chan

	k.config = config

	logx.Infof("KafkaServer.Init, config: %+v", k.config)

	var err error

	err = k.InitKafkaApi()
	if err != nil {
		return err
	}

	err = k.InitListenPubChan()
	if err != nil {
		return err
	}

	return nil
}

func (k *KafkaServer) InitKafkaApi() error {
	logx.Info("KafkaServer.InitKafkaApi")

	var err error

	k.Producer, err = sarama.NewSyncProducer([]string{k.config.IP}, nil)
	if err != nil {
		logx.Error(err.Error())
		return err
	}

	k.Broker = sarama.NewBroker(k.config.IP)
	broker_config := sarama.NewConfig()
	err = k.Broker.Open(broker_config)
	if err != nil {
		logx.Error(err.Error())
		return err
	}

	return nil
}

func (k *KafkaServer) InitListenPubChan() error {
	// k.PubDataChan = &datastruct.DataChannel{
	// 	TradeChannel: make(chan *datastruct.Trade),
	// 	KlineChannel: make(chan *datastruct.Kline),
	// 	DepthChannel: make(chan *datastruct.DepthQuote),
	// }

	logx.Info("KafkaServer.InitListenPubChan")

	go func() {
		for {
			select {
			case local_depth := <-k.PubDataChan.DepthChannel:
				go k.PublishDepth(local_depth)
			case local_kline := <-k.PubDataChan.KlineChannel:
				go k.PublishKline(local_kline)
			case local_trade := <-k.PubDataChan.TradeChannel:
				go k.PublishTrade(local_trade)
			}
		}
	}()
	return nil
}

// Start Consume Topic
func (k *KafkaServer) Start() {
	if k.IsTest {
		k.start_test()
		return
	}

	k.start_consume()
}

func (k *KafkaServer) start_consume() {
	if len(k.ConsumeSet) == 0 {
		logx.Info("ConsumeSet is Empty!")
		return
	}

	for _, consume_item := range k.ConsumeSet {
		go k.ConsumeSingleTopic(consume_item)
	}
}

func (k *KafkaServer) start_test() {

}

func (k *KafkaServer) UpdateMetaData(meta_data *datastruct.Metadata) {
	logx.Info("UpdateMetaData: " + fmt.Sprintf("%+v", meta_data))

	NewConsumeSet := GetConsumeSet(*meta_data)
	logx.Info(fmt.Sprintf("NewTopicSet: %+v", NewConsumeSet))

	if k.ConsumeSet == nil {
		k.ConsumeSet = make(map[string](*ConsumeItem))
	}

	for new_topic, consume_item := range NewConsumeSet {
		if _, ok := k.ConsumeSet[new_topic]; !ok {
			logx.Info("Start Consume Topic: " + new_topic)
			go k.ConsumeSingleTopic(consume_item)
			k.ConsumeSet[new_topic] = consume_item
		}
	}

	for old_topic, consume_item := range k.ConsumeSet {
		if _, ok := NewConsumeSet[old_topic]; !ok {
			logx.Info("Stop Consume Topic: " + old_topic)
			consume_item.CancelFunc()
			delete(k.ConsumeSet, old_topic)
		}
	}
}

func (k *KafkaServer) ConsumeSingleTopic(consume_item *ConsumeItem) {
	consumer, err := sarama.NewConsumer([]string{k.config.IP}, nil)

	logx.Info("ConsumeSingleTopic: " + consume_item.Topic)

	if err != nil {
		logx.Error(err.Error())
		return
	}

	for {
		select {
		case <-consume_item.Ctx.Done():
			logx.Info(consume_item.Topic + " listen Over!")
			return
		default:
			k.ConsumeAtom(consume_item.Topic, consumer)
		}
	}
}

func (k *KafkaServer) ConsumeAtom(topic string, consumer sarama.Consumer) {
	logx.Info("ConsumeAtom: " + topic)
	partitionList, err := consumer.Partitions(topic) // 根据topic取到所有的分区
	if err != nil {
		logx.Error(err.Error())
		return
	}

	for partition := range partitionList {
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)

		if err != nil {
			logx.Error(err.Error())
			continue
		}
		defer pc.AsyncClose()

		for msg := range pc.Messages() {

			// logx.Info(msg.Topic)

			topic_type := GetTopicType(msg.Topic)

			// logx.Info(topic_type)

			switch topic_type {
			case DEPTH_TYPE:
				go k.ProcessDepthBytes(msg.Value)
			case KLINE_TYPE:
				go k.ProcessKlineBytes(msg.Value)
			case TRADE_TYPE:
				go k.ProcessTradeBytes(msg.Value)
			default:
				logx.Error("Unknown Topic " + topic_type)
			}
		}
	}
}

func (k *KafkaServer) PublishMsg(topic string, origin_bytes []byte) error {
	defer k.PublishMutex.Unlock()

	msgs := []*sarama.ProducerMessage{{
		Topic: topic,
		Value: sarama.ByteEncoder(origin_bytes),
	}}

	k.PublishMutex.Lock()
	err := k.Producer.SendMessages(msgs)

	if err != nil {
		logx.Error(err.Error())
		return err
	}
	return nil
}

func (k *KafkaServer) PublishDepth(local_depth *datastruct.DepthQuote) error {
	logx.Info(fmt.Sprintf("Pub Depth %+v", local_depth.String(3)))
	serialize_str, err := k.Serializer.EncodeDepth(local_depth)

	if err != nil {
		logx.Error(err.Error())
		return err
	}

	topic := GetDepthTopic(local_depth.Symbol, local_depth.Exchange)

	return k.PublishMsg(topic, serialize_str)
}

func (k *KafkaServer) PublishKline(local_kline *datastruct.Kline) error {
	logx.Info(fmt.Sprintf("Pub kline %+v", local_kline))

	serialize_str, err := k.Serializer.EncodeKline(local_kline)

	if err != nil {
		logx.Error(err.Error())
		return err
	}

	topic := GetKlineTopic(local_kline.Symbol, local_kline.Exchange)

	return k.PublishMsg(topic, serialize_str)
}

func (k *KafkaServer) PublishTrade(local_trade *datastruct.Trade) error {
	logx.Info(fmt.Sprintf("Pub Trade %+v", local_trade))

	serialize_str, err := k.Serializer.EncodeTrade(local_trade)

	if err != nil {
		logx.Error(err.Error())
		return err
	}

	topic := GetTradeTopic(local_trade.Symbol, local_trade.Exchange)

	return k.PublishMsg(topic, serialize_str)
}

func (k *KafkaServer) ProcessDepthBytes(depth_bytes []byte) error {

	local_depth, err := k.Serializer.DecodeDepth(depth_bytes)

	if err != nil {
		logx.Error(err.Error())
		return err
	}

	k.SendRecvedDepth(local_depth)

	return nil
}

func (k *KafkaServer) ProcessKlineBytes(kline_bytes []byte) error {
	local_kline, err := k.Serializer.DecodeKline(kline_bytes)

	if err != nil {
		logx.Error(err.Error())
		return err
	}

	k.SendRecvedKline(local_kline)

	return nil
}

func (k *KafkaServer) ProcessTradeBytes(trade_bytes []byte) error {
	local_trade, err := k.Serializer.DecodeTrade(trade_bytes)

	if err != nil {
		logx.Error(err.Error())
		return err
	}

	k.SendRecvedTrade(local_trade)

	return nil
}

func (k *KafkaServer) SendRecvedDepth(depth *datastruct.DepthQuote) {
	logx.Slowf("[kafka] Rcv Depth: %s \n", depth.String(3))
	k.RecvDataChan.DepthChannel <- depth
}

func (k *KafkaServer) SendRecvedKline(kline *datastruct.Kline) {
	logx.Slowf("[kafka] Rcv kline: %s \n", kline.String())
	k.RecvDataChan.KlineChannel <- kline
}

func (k *KafkaServer) SendRecvedTrade(trade *datastruct.Trade) {
	logx.Slowf("[kafka] Rcv Trade: %s \n", trade.String())
	k.RecvDataChan.TradeChannel <- trade
}
