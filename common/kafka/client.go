package kafka

import (
	"context"
	"fmt"
	"market_server/common/config"
	"market_server/common/datastruct"
	"sync"
	"time"

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
	Client   sarama.Client

	MetaData datastruct.Metadata

	Serializer datastruct.SerializerI

	RecvDataChan *datastruct.DataChannel
	PubDataChan  *datastruct.DataChannel

	ConsumeSet map[string](*ConsumeItem)

	PublishMutex sync.Mutex

	config config.KafkaConfig
	IsTest bool

	consume_lock sync.Mutex

	consume_topics       map[string]struct{}
	consume_topics_mutex sync.Mutex

	created_topics map[string]struct{}

	statistic_secs     int
	rcv_statistic_info sync.Map
	pub_statistic_info sync.Map

	statistic_start time.Time
}

func NewKafka(serializer datastruct.SerializerI,
	recv_data_chan *datastruct.DataChannel,
	pub_data_chan *datastruct.DataChannel,
	config config.KafkaConfig) (*KafkaServer, error) {
	server := &KafkaServer{
		Serializer:     serializer,
		RecvDataChan:   recv_data_chan,
		PubDataChan:    pub_data_chan,
		config:         config,
		consume_topics: make(map[string]struct{}),
		statistic_secs: 10,
		ConsumeSet:     make(map[string](*ConsumeItem)),
		created_topics: make(map[string]struct{}),
	}

	err := server.InitApi()

	return server, err
}

func (k *KafkaServer) StatisticTimeTaskMain() {
	logx.Info("---- StatisticTimeTask Start!")
	duration := time.Duration((time.Duration)(k.statistic_secs) * time.Second)
	timer := time.Tick(duration)

	k.statistic_start = time.Now()
	for {
		select {
		case <-timer:
			k.UpdateStatisticInfo()
		}
	}
}

func (k *KafkaServer) OutputRcvInfo(key, value interface{}) bool {
	if value.(int) != 0 {
		logx.Statf("[rcv] %s : %d ", key, value)
		k.rcv_statistic_info.Store(key, 0)
	}

	return true
}

func (k *KafkaServer) OutputPubInfo(key, value interface{}) bool {
	if value.(int) != 0 {
		logx.Statf("[pub] %s : %d ", key, value)
		k.pub_statistic_info.Store(key, 0)
	}
	return true
}

func (k *KafkaServer) UpdateStatisticInfo() {

	logx.Statf("kafka Statistic Start: %+v \n", k.statistic_start)

	k.rcv_statistic_info.Range(k.OutputRcvInfo)

	k.pub_statistic_info.Range(k.OutputPubInfo)

	k.statistic_start = time.Now()

	logx.Statf("kafka Statistic End: %+v \n", k.statistic_start)
}

func (k *KafkaServer) CreatedTopics() map[string]struct{} {
	return k.created_topics
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
	k.consume_topics = make(map[string]struct{})
	k.statistic_secs = 10

	logx.Infof("KafkaServer.Init, config: %+v\n", k.config)

	var err error

	err = k.InitApi()
	if err != nil {
		return err
	} else {
		logx.Infof("InitApi Err: %+v", err)
	}

	k.UpdateCreateTopics()
	k.StartListenPubChan()

	return nil
}

func (k *KafkaServer) InitApi() error {
	logx.Info("KafkaServer.InitApi")

	var err error

	k.Producer, err = sarama.NewSyncProducer([]string{k.config.IP}, nil)
	if err != nil {
		logx.Error(err.Error())
		return err
	}

	k.Client, err = sarama.NewClient([]string{k.config.IP}, nil)
	if err != nil {
		logx.Errorf("NewClient Failed %s ", err.Error())
		return err
	}

	k.Broker = sarama.NewBroker(k.config.IP)
	broker_config := sarama.NewConfig()
	err = k.Broker.Open(broker_config)
	if err != nil {
		logx.Errorf("Broker Open Error: %s", err.Error())
		return err
	}
	ok, err := k.Broker.Connected()
	if !ok {
		logx.Errorf("Broker Connect Error: %s", err.Error())
		return err
	}

	return nil
}

func (k *KafkaServer) UpdateCreateTopics() error {
	logx.Info("UpdateCreateTopics")

	online_topics, err := k.Client.Topics()
	if err != nil {
		logx.Errorf("Get Online Topics Failed %s ", err.Error())
		return err
	}
	logx.Infof("Online Topics: %+v", online_topics)

	for _, topic := range online_topics {
		if _, ok := k.created_topics[topic]; !ok {
			k.created_topics[topic] = struct{}{}
		}
	}
	return nil
}

func (k *KafkaServer) StartListenPubChan() {
	logx.Info("KafkaServer.StartListenPubChan")

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
}

func (k *KafkaServer) IsTopicConsumed(topic string) bool {

	k.consume_topics_mutex.Lock()
	_, ok := k.consume_topics[topic]
	k.consume_topics_mutex.Unlock()

	return ok
}

func (k *KafkaServer) AddConsumeTopic(topic string) {
	k.consume_topics_mutex.Lock()
	if _, ok := k.consume_topics[topic]; !ok {
		k.consume_topics[topic] = struct{}{}

		logx.Infof("Add Consume Topic: %s", topic)
	}
	k.consume_topics_mutex.Unlock()
}

func (k *KafkaServer) DelConsumeTopic(topic string) {
	k.consume_topics_mutex.Lock()
	if _, ok := k.consume_topics[topic]; !ok {
		delete(k.consume_topics, topic)

		logx.Infof("Del Consume Topic: %s", topic)
	}
	k.consume_topics_mutex.Unlock()
}

// Start Consume Topic
func (k *KafkaServer) Start() error {

	k.UpdateCreateTopics()
	k.StartListenPubChan()

	go k.StatisticTimeTaskMain()
	k.start_consume()
	return nil
}

func (k *KafkaServer) start_consume() {
	k.consume_lock.Lock()

	defer k.consume_lock.Unlock()

	logx.Infof("CurrConsumeSet: %+v", k.ConsumeSet)

	if len(k.ConsumeSet) == 0 {
		logx.Info("ConsumeSet is Empty!")
	} else {
		for _, consume_item := range k.ConsumeSet {
			go k.ConsumeSingleTopic(consume_item)
		}
	}
}

func (k *KafkaServer) UpdateMetaData(meta_data *datastruct.Metadata) {
	k.consume_lock.Lock()

	logx.Infof("UpdateMetaData: %+v", meta_data)

	NewConsumeSet := GetConsumeSet(*meta_data)
	logx.Infof("NewTopicSet: %+v", NewConsumeSet)

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

	defer k.consume_lock.Unlock()
}

func (k *KafkaServer) ConsumeSingleTopic(consume_item *ConsumeItem) {
	consumer, err := sarama.NewConsumer([]string{k.config.IP}, nil)

	if err != nil {
		logx.Error(err.Error())
		return
	}

	if k.IsTopicConsumed(consume_item.Topic) {
		logx.Infof("Topic %s already consumed!", consume_item.Topic)
		return
	}
	k.AddConsumeTopic(consume_item.Topic)

	partitionList, err := consumer.Partitions(consume_item.Topic) // 根据topic取到所有的分区
	if err != nil {
		logx.Error(err.Error())
		return
	}

	logx.Infof("Conume Topic: %s ", consume_item.Topic)

	for {
		for partition := range partitionList {
			pc, err := consumer.ConsumePartition(consume_item.Topic, int32(partition), sarama.OffsetNewest)

			if err != nil {
				logx.Error(err.Error())
				continue
			}
			defer pc.AsyncClose()

			for msg := range pc.Messages() {
				topic_type := GetTopicType(msg.Topic)

				if value, ok := k.rcv_statistic_info.Load(msg.Topic); ok {
					k.rcv_statistic_info.Store(msg.Topic, value.(int)+1)
				} else {
					k.rcv_statistic_info.Store(msg.Topic, 1)
				}

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

				select {
				case <-consume_item.Ctx.Done():
					logx.Info(consume_item.Topic + " listen Over!")
					return
				default:
					// time.Sleep(time.Second)
				}
			}
		}

		logx.Infof("Wait For Connect!")

		time.Sleep(time.Second * 3)
	}

}

func (k *KafkaServer) ConsumeAtom(topic string, consumer sarama.Consumer) {

}

func (k *KafkaServer) CheckTopic(topic string) bool {
	_, ok := k.created_topics[topic]

	if !ok {
		logx.Infof("Topic: %+v not exists, need to be created!")
	}
	return ok
}

func (k *KafkaServer) CreateTopic(topic string) bool {

	create_detail_map := make(map[string]*sarama.TopicDetail)

	create_detail_map[topic] = &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	create_req := &sarama.CreateTopicsRequest{
		TopicDetails: create_detail_map,
		Timeout:      time.Second * 15,
	}

	create_respons, err := k.Broker.CreateTopics(create_req)

	if err != nil {
		logx.Errorf("CreateTopics Error: %+v", err)
		return false
	} else {
		logx.Infof("create_respons: %+v", create_respons)
	}

	k.UpdateCreateTopics()

	_, ok := k.created_topics[topic]

	return ok
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

func (k *KafkaServer) PublishMsg(topic string, origin_bytes []byte) error {

	if !k.CheckTopic(topic) && !k.CreateTopic(topic) {
		return fmt.Errorf("%s is not created and create it failed!", topic)
	}

	if value, ok := k.pub_statistic_info.Load(topic); ok {
		k.pub_statistic_info.Store(topic, value.(int)+1)
	} else {
		k.pub_statistic_info.Store(topic, 1)
	}

	msgs := []*sarama.ProducerMessage{{
		Topic: topic,
		Value: sarama.ByteEncoder(origin_bytes),
	}}

	k.PublishMutex.Lock()
	err := k.Producer.SendMessages(msgs)

	if err != nil {
		logx.Errorf("[kafka] SendMessages: %s ", err.Error())
		logx.Infof("[kafka] SendMessages: %s ", err.Error())
		k.PublishMutex.Unlock()
		return err
	}
	k.PublishMutex.Unlock()
	return nil
}

func (k *KafkaServer) PublishDepth(local_depth *datastruct.DepthQuote) error {
	// logx.Slow(fmt.Sprintf("Pub Depth %+v", local_depth.String(3)))
	serialize_str, err := k.Serializer.EncodeDepth(local_depth)

	if err != nil {
		logx.Error(err.Error())
		return err
	}

	topic := GetDepthTopic(local_depth.Symbol, local_depth.Exchange)

	return k.PublishMsg(topic, serialize_str)
}

func (k *KafkaServer) PublishKline(local_kline *datastruct.Kline) error {
	// logx.Slow(fmt.Sprintf("Pub kline %+v", local_kline))

	serialize_str, err := k.Serializer.EncodeKline(local_kline)

	if err != nil {
		logx.Error(err.Error())
		return err
	}

	topic := GetKlineTopic(local_kline.Symbol, local_kline.Exchange)

	return k.PublishMsg(topic, serialize_str)
}

func (k *KafkaServer) PublishTrade(local_trade *datastruct.Trade) error {
	logx.Slowf("Pub Trade %+v", local_trade)

	serialize_str, err := k.Serializer.EncodeTrade(local_trade)

	if err != nil {
		logx.Error(err.Error())
		return err
	}

	topic := GetTradeTopic(local_trade.Symbol, local_trade.Exchange)

	return k.PublishMsg(topic, serialize_str)
}

func (k *KafkaServer) SendRecvedDepth(depth *datastruct.DepthQuote) {
	// logx.Slowf("[kafka] Rcv Depth: %s \n", depth.String(3))
	k.RecvDataChan.DepthChannel <- depth
}

func (k *KafkaServer) SendRecvedKline(kline *datastruct.Kline) {
	// logx.Slowf("[kafka] Rcv kline: %s \n", kline.String())
	k.RecvDataChan.KlineChannel <- kline
}

func (k *KafkaServer) SendRecvedTrade(trade *datastruct.Trade) {
	// logx.Slowf("[kafka] Rcv Trade: %s \n", trade.String())
	k.RecvDataChan.TradeChannel <- trade
}
