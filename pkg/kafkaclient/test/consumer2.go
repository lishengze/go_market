package main

//const topic = "my-topic"
//const topic = "DEPTH.BTC_USDT.BINANCE"
//const topic = "KLINE.BTC_USDT.BINANCE"
//const topic = "TRADE.BTC_USDT.BINANCE"

const topic = "DEPTH.BTC_USDT._bcts_"

//const topic = "DEPTH.ETH_USDT.FTX"

//const topic = "KLINE.BTC_USDT.FTX"
//const topic = "TRADE.BTC_USDT.FTX"

// func main() {
// 	// kafka consumer
// 	consumer, err := sarama.NewConsumer([]string{"152.32.254.76:9117"}, nil)
// 	if err != nil {
// 		fmt.Printf("fail to start consumer, err:%v\n", err)
// 		return
// 	}
// 	partitionList, err := consumer.Partitions(topic) // 根据topic取到所有的分区
// 	if err != nil {
// 		fmt.Printf("fail to get list of partition:err%v\n", err)
// 		return
// 	}
// 	fmt.Println(partitionList)
// 	for partition := range partitionList { // 遍历所有的分区
// 		// 针对每个分区创建一个对应的分区消费者
// 		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
// 		if err != nil {
// 			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
// 			return
// 		}
// 		defer pc.AsyncClose()
// 		// 异步从每个分区消费信息
// 		go func(sarama.PartitionConsumer) {
// 			for msg := range pc.Messages() {
// 				data := mpupb.Depth{}
// 				//data:=mpupb.Kline{}
// 				//data:=mpupb.Trade{}
// 				err := proto.Unmarshal(msg.Value, &data)
// 				if err != nil {
// 					fmt.Println(err)
// 					continue
// 				}
// 				fmt.Println(
// 					data.Timestamp.AsTime().Format(time.RFC3339Nano),
// 					data.MpuTimestamp.AsTime().Format(time.RFC3339Nano),
// 					data.Symbol,
// 					data.String(),
// 				)
// 				time.Sleep(time.Second)
// 				//fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v", msg.Partition, msg.Offset, msg.Key, msg.Value)
// 			}
// 		}(pc)
// 	}

// 	select {}
// }
