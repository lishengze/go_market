package main

import (
	"exterior-interactor/app/mpu/rpc/mpupb"
	"fmt"
	"google.golang.org/protobuf/proto"
	"time"

	"github.com/Shopify/sarama"
)

//const topic = "my-topic"
const topic = "DEPTH.BTC_USDT.BINANCE"

func main() {
	// kafka consumer
	consumer, err := sarama.NewConsumer([]string{"152.32.254.76:9117"}, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return
	}
	partitionList, err := consumer.Partitions(topic) // 根据topic取到所有的分区
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}
	fmt.Println(partitionList)
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
				depth:=mpupb.Depth{}
				err:=proto.Unmarshal(msg.Value,&depth)
				if err!=nil{
					fmt.Println(err)
					continue
				}
				fmt.Println(depth.Timestamp.AsTime().Format(time.RFC3339Nano))
				//fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v", msg.Partition, msg.Offset, msg.Key, msg.Value)
			}
		}(pc)
	}

	select {}
}
