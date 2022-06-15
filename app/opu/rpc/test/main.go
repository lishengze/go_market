package main

import (
	"context"
	"exterior-interactor/app/opu/rpc/opu"
	"exterior-interactor/app/opu/rpc/opupb"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/protobuf/proto"
	"log"
)

const (
	opuAddress   = "152.32.254.76:8606"
	kafkaAddress = "152.32.254.76:9117"

	topic = "ORDER.FTX_MCA_OTC_TRADING"
)

var o = getOpu(opuAddress)

func main() {
	subUpdate()

	//testRegisterAccount()
	//testGetBalance()
	//testPlaceOrder()
	testCancelOrder()

	select {}
}

func getOpu(address string) opu.Opu {
	client, err := zrpc.NewClientWithTarget(address)
	if err != nil {
		log.Fatalln(err)
	}
	return opu.NewOpu(client)
}

func testPlaceOrder() {
	rsp, err := o.PlaceOrder(context.Background(), &opu.PlaceOrderReq{
		AccountId:     "",
		AccountAlias:  "FTX_MCA_OTC_TRADING",
		ClientOrderId: "12345",
		StdSymbol:     "BTC_USDT",
		Volume:        "0.01",
		Price:         "12000",
		Type:          "LIMIT",
		Side:          "BUY",
	})

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(rsp)
}

func testCancelOrder() {
	rsp, err := o.CancelOrder(context.Background(), &opu.CancelOrderReq{
		AccountId:     "",
		AccountAlias:  "FTX_MCA_OTC_TRADING",
		ClientOrderId: "12345",
	})

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(rsp)
}

func testRegisterAccount() {
	rsp, err := o.RegisterAccount(context.Background(), &opu.RegisterAccountReq{
		Alias:          "FTX_MCA_OTC_TRADING",
		Key:            "1IVOM9S2EELFcOZGJ3RLIg3X_ETRwOM4sDH9a_D5",
		Secret:         "A3US9cW3biFIO9xQYfRdTFpD4UX0gNCcTyzY2F5W",
		Passphrase:     "",
		Exchange:       "FTX",
		SubAccountName: "Xpert RFQ",
	})

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(rsp)
}

func testGetBalance() {
	rsp, err := o.QueryBalance(context.Background(), &opu.QueryBalanceReq{
		AccountAlias: "FTX_MCA_OTC_TRADING",
	})

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(rsp)
}

func subUpdate() {
	// kafka consumer
	consumer, err := sarama.NewConsumer([]string{kafkaAddress}, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return
	}
	partitionList, err := consumer.Partitions(topic) // 根据topic取到所有的分区
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}
	//fmt.Println(partitionList)
	for partition := range partitionList {
		// 遍历所有的分区
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
				data := opupb.OrderTradesUpdate{}
				err := proto.Unmarshal(msg.Value, &data)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println(
					data.String(),
				)
			}
		}(pc)
	}
}
