package kafka

import (
	"fmt"
	"market_aggregate/pkg/conf"
	"market_aggregate/pkg/datastruct"
	"market_aggregate/pkg/protostruct"
	"market_aggregate/pkg/util"

	"github.com/Shopify/sarama"
	"google.golang.org/protobuf/proto"
)

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

type KafkaServer struct {
	Consumer sarama.Consumer
	Producer sarama.SyncProducer
	Broker   *sarama.Broker

	Serializer   datastruct.SerializerI
	RecvDataChan *datastruct.DataChannel
}

// Init(*conf.Config, SerializerI, *DataChannel)
func (k *KafkaServer) Init(config *conf.Config, serializer datastruct.SerializerI, recv_data_chan *datastruct.DataChannel) error {
	var err error

	k.Serializer = serializer
	k.RecvDataChan = recv_data_chan

	k.Consumer, err = sarama.NewConsumer([]string{config.IP}, nil)
	if err != nil {
		util.LOG_ERROR(err.Error())
		return err
	}

	k.Producer, err = sarama.NewSyncProducer([]string{config.IP}, nil)
	if err != nil {
		util.LOG_ERROR(err.Error())
		return err
	}

	k.Broker = sarama.NewBroker(config.IP)
	broker_config := sarama.NewConfig()
	err = k.Broker.Open(broker_config)
	if err != nil {
		util.LOG_ERROR(err.Error())
		return err
	}

	return nil
}

func (k *KafkaServer) PublishDepth(*datastruct.DepthQuote) {

}

func (k *KafkaServer) PublishKline(*datastruct.Kline) {

}

func (k *KafkaServer) PublishTrade(*datastruct.Trade) {

}

func (k *KafkaServer) SendRecvedDepth(*datastruct.DepthQuote) {

}

func (k *KafkaServer) SendRecvedKline(*datastruct.Kline) {

}

func (k *KafkaServer) SendRedvedTrade(*datastruct.Trade) {

}
