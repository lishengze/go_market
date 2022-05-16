package kafka

import (
	"context"
	"fmt"
	"market_aggregate/pkg/conf"
	"market_aggregate/pkg/datastruct"
	"market_aggregate/pkg/protostruct"
	"market_aggregate/pkg/util"
	"sync"

	"github.com/Shopify/sarama"
	"google.golang.org/protobuf/proto"
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

	Config   *conf.Config
	MetaData datastruct.Metadata

	Serializer datastruct.SerializerI

	RecvDataChan *datastruct.DataChannel
	PubDataChan  *datastruct.DataChannel

	ConsumeSet map[string](*ConsumeItem)

	PublishMutex sync.Mutex

	IsTest bool
}

// Init(*conf.Config, SerializerI, *DataChannel)
func (k *KafkaServer) Init(config *conf.Config, serializer datastruct.SerializerI,
	recv_data_chan *datastruct.DataChannel,
	pub_data_chan *datastruct.DataChannel) error {

	k.Serializer = serializer
	k.Config = config
	k.RecvDataChan = recv_data_chan
	k.PubDataChan = pub_data_chan

	util.LOG_INFO("KafkaServer.Init")

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
	util.LOG_INFO("KafkaServer.InitKafkaApi")

	var err error

	k.Producer, err = sarama.NewSyncProducer([]string{k.Config.IP}, nil)
	if err != nil {
		util.LOG_ERROR(err.Error())
		return err
	}

	k.Broker = sarama.NewBroker(k.Config.IP)
	broker_config := sarama.NewConfig()
	err = k.Broker.Open(broker_config)
	if err != nil {
		util.LOG_ERROR(err.Error())
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

	util.LOG_INFO("KafkaServer.InitListenPubChan")

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

func (k *KafkaServer) Start() {
	if k.IsTest {
		k.start_test()
		return
	}

	k.start_consume()
}

func (k *KafkaServer) start_consume() {
	for _, consume_item := range k.ConsumeSet {
		go k.ConsumeSingleTopic(consume_item)
	}
}

func (k *KafkaServer) start_test() {

}

func (k *KafkaServer) UpdateMetaData(meta_data datastruct.Metadata) {
	NewConsumeSet := GetConsumeSet(k.MetaData)
	util.LOG_INFO(fmt.Sprintf("UpdatedTopicSet: %+v", k.ConsumeSet))

	for new_topic, consume_item := range NewConsumeSet {
		if _, ok := k.ConsumeSet[new_topic]; ok == false {
			util.LOG_INFO("Start Consume Topic: " + new_topic)
			go k.ConsumeSingleTopic(consume_item)
			k.ConsumeSet[new_topic] = consume_item
		}
	}

	for old_topic, consume_item := range k.ConsumeSet {
		if _, ok := NewConsumeSet[old_topic]; ok == false {
			util.LOG_INFO("Stop Consume Topic: " + old_topic)
			consume_item.CancelFunc()
			delete(k.ConsumeSet, old_topic)
		}
	}
}

func (k *KafkaServer) ConsumeSingleTopic(consume_item *ConsumeItem) {
	consumer, err := sarama.NewConsumer([]string{k.Config.IP}, nil)

	util.LOG_INFO("ConsumeSingleTopic: " + consume_item.Topic)

	if err != nil {
		util.LOG_ERROR(err.Error())
		return
	}

	for {
		select {
		case <-consume_item.Ctx.Done():
			util.LOG_INFO(consume_item.Topic + " listen Over!")
			return
		default:
			k.ConsumeAtom(consume_item.Topic, consumer)
		}
	}
}

func (k *KafkaServer) ConsumeAtom(topic string, consumer sarama.Consumer) {
	util.LOG_INFO("ConsumeAtom: " + topic)
	partitionList, err := consumer.Partitions(topic) // 根据topic取到所有的分区
	if err != nil {
		util.LOG_ERROR(err.Error())
		return
	}

	for partition := range partitionList {
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)

		if err != nil {
			util.LOG_ERROR(err.Error())
			continue
		}
		defer pc.AsyncClose()

		for msg := range pc.Messages() {

			// util.LOG_INFO(msg.Topic)

			topic_type := GetTopicType(msg.Topic)

			// util.LOG_INFO(topic_type)

			switch topic_type {
			case DEPTH_TYPE:
				go k.ProcessDepthBytes(msg.Value)
			case KLINE_TYPE:
				go k.ProcessKlineBytes(msg.Value)
			case TRADE_TYPE:
				go k.ProcessTradeBytes(msg.Value)
			default:
				util.LOG_ERROR("Unknown Topic " + topic_type)
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
		util.LOG_ERROR(err.Error())
		return err
	}
	return nil
}

func (k *KafkaServer) PublishDepth(local_depth *datastruct.DepthQuote) error {
	serialize_str, err := k.Serializer.EncodeDepth(local_depth)

	if err != nil {
		util.LOG_ERROR(err.Error())
		return err
	}

	topic := GetDepthTopic(local_depth.Symbol, local_depth.Exchange)

	return k.PublishMsg(topic, serialize_str)
}

func (k *KafkaServer) PublishKline(local_kline *datastruct.Kline) error {
	serialize_str, err := k.Serializer.EncodeKline(local_kline)

	if err != nil {
		util.LOG_ERROR(err.Error())
		return err
	}

	topic := GetKlineTopic(local_kline.Symbol, local_kline.Exchange)

	return k.PublishMsg(topic, serialize_str)
}

func (k *KafkaServer) PublishTrade(local_trade *datastruct.Trade) error {

	serialize_str, err := k.Serializer.EncodeTrade(local_trade)

	if err != nil {
		util.LOG_ERROR(err.Error())
		return err
	}

	topic := GetTradeTopic(local_trade.Symbol, local_trade.Exchange)

	return k.PublishMsg(topic, serialize_str)
}

func (k *KafkaServer) ProcessDepthBytes(depth_bytes []byte) error {

	local_depth, err := k.Serializer.DecodeDepth(depth_bytes)

	if err != nil {
		util.LOG_ERROR(err.Error())
		return err
	}

	k.SendRecvedDepth(local_depth)

	return nil
}

func (k *KafkaServer) ProcessKlineBytes(kline_bytes []byte) error {
	local_kline, err := k.Serializer.DecodeKline(kline_bytes)

	if err != nil {
		util.LOG_ERROR(err.Error())
		return err
	}

	k.SendRecvedKline(local_kline)

	return nil
}

func (k *KafkaServer) ProcessTradeBytes(trade_bytes []byte) error {
	local_trade, err := k.Serializer.DecodeTrade(trade_bytes)

	if err != nil {
		util.LOG_ERROR(err.Error())
		return err
	}

	k.SendRecvedTrade(local_trade)

	return nil
}

func (k *KafkaServer) SendRecvedDepth(depth *datastruct.DepthQuote) {
	k.RecvDataChan.DepthChannel <- depth
}

func (k *KafkaServer) SendRecvedKline(kline *datastruct.Kline) {
	k.RecvDataChan.KlineChannel <- kline
}

func (k *KafkaServer) SendRecvedTrade(trade *datastruct.Trade) {
	k.RecvDataChan.TradeChannel <- trade
}

func TestConsumer() {
	fmt.Println("------- TestConsumer-------")
	consumer, err := sarama.NewConsumer([]string{"43.154.179.47:9117"}, nil)

	topic := "TRADEx-BTC_USDT"

	partitionList, err := consumer.Partitions(topic) // 根据topic取到所有的分区
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}
	fmt.Println("partitionList: ", partitionList)
	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return
		}
		defer pc.AsyncClose()
		// 异步从每个分区消费信息
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("%+v \n", msg)

				trade := protostruct.Trade{}
				// //data:=mpupb.Kline{}
				// //data:=mpupb.Trade{}
				err := proto.Unmarshal(msg.Value, &trade)
				if err != nil {
					fmt.Println(err)
					continue
				}

				fmt.Printf("%+v \n", trade)

				// fmt.Println(
				// 	data.Timestamp.AsTime().Format(time.RFC3339Nano),
				// 	data.MpuTimestamp.AsTime().Format(time.RFC3339Nano),
				// 	data.Symbol,
				// 	//data.String(),
				// )
				//fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v", msg.Partition, msg.Offset, msg.Key, msg.Value)
			}
		}(pc)
	}

	select {}

}
