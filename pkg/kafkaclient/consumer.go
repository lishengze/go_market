package kafkaclient

import (
	"encoding/json"
	"fmt"
	"market_server/pkg/kafkaclient/mpupb"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"google.golang.org/protobuf/proto"
)

var DepthDataList = make(map[string]mpupb.Depth, 2000)

type KafkaDepthPeer struct {
	TopicSymbolList        []string
	Address                string
	Port                   string
	latestTimestampMap     map[string]int64
	latestTimestampMapLock sync.RWMutex //读写锁
	latestMpuTimestamp     map[string]int64
	latestMpuTimestampLock sync.RWMutex //读写锁
}

func NewKafkaClient(symbolList []string, Address, Port string) *KafkaDepthPeer {
	return &KafkaDepthPeer{
		TopicSymbolList:        symbolList,
		Address:                Address,
		Port:                   Port,
		latestTimestampMap:     make(map[string]int64),
		latestTimestampMapLock: sync.RWMutex{},
		latestMpuTimestamp:     make(map[string]int64),
		latestMpuTimestampLock: sync.RWMutex{},
	}
}

func (kc *KafkaDepthPeer) GetDepthInfo() map[string]mpupb.Depth {
	return DepthDataList
}

func (kc *KafkaDepthPeer) FetchDepthWorkWather() {
	config := sarama.NewConfig() //使用默认配置，断开连接后恢复时会自动重连
	config.Net.DialTimeout = 10 * time.Second
	config.Net.ReadTimeout = 10 * time.Second
	config.Net.WriteTimeout = 10 * time.Second
	consumer, err := sarama.NewConsumer([]string{kc.Address + ":" + kc.Port}, config)
	if err != nil {
		logx.Errorf("fail to start consumer, err:%v\n", err)
		return
	}
	for _, symbol := range kc.TopicSymbolList {
		symbolTopic := "DEPTH." + symbol + "._bcts_"

		partitionList, err := consumer.Partitions(symbolTopic) // 根据topic取到所有的分区
		if err != nil {
			logx.Errorf("fail to get list of partition:err%v\n", err)
			return
		}
		for partition := range partitionList { // 遍历所有的分区
			// 针对每个分区创建一个对应的分区消费者
			pc, err := consumer.ConsumePartition(symbolTopic, int32(partition), sarama.OffsetNewest)
			if err != nil {
				logx.Errorf("failed to start consumer for partition %d,err:%v\n", partition, err)
				return
			}
			defer pc.AsyncClose()
			// 异步从每个分区消费信息
			threading.GoSafe(func() {
				var i int
				for msg := range pc.Messages() {
					data := mpupb.Depth{}
					err := proto.Unmarshal(msg.Value, &data)
					if err != nil {
						logx.Error(err)
						continue
					}
					if i == 0 { //测试用，只打印一次日志
						//logx.Infof("kafka receive message:%+v", data)
						bytes, _ := json.Marshal(data)
						fmt.Printf("kafka receive message:%+v\n", string(bytes))
						i++
					}
					DepthDataList[symbol] = data
					kc.setLatestTimestamp(data.Symbol, data.Timestamp.AsTime().Unix())
					kc.setLatestTimestampMap(data.Symbol, data.MpuTimestamp.AsTime().Unix())
				}
			})
		}
	}
	select {}
}

func (kc *KafkaDepthPeer) setLatestTimestamp(symbol string, t int64) {
	kc.latestTimestampMapLock.Lock()
	defer kc.latestTimestampMapLock.Unlock()
	kc.latestTimestampMap[symbol] = t
	//logx.Infof("setLatestTimestamp t=%d", t)
}

func (kc *KafkaDepthPeer) GetLatestMpuTimestamp(symbol string) int64 {
	kc.latestTimestampMapLock.RLock()
	defer kc.latestTimestampMapLock.RUnlock()
	if t, ok := kc.latestTimestampMap[symbol]; ok {
		return t
	}
	return 0
}

func (kc *KafkaDepthPeer) setLatestTimestampMap(symbol string, t int64) {
	kc.latestMpuTimestampLock.Lock()
	defer kc.latestMpuTimestampLock.Unlock()
	kc.latestMpuTimestamp[symbol] = t
	//logx.Infof("setLatestTimestampMap t=%d", t)
}

func (kc *KafkaDepthPeer) GetLatestTimestampMap(symbol string) int64 {
	kc.latestMpuTimestampLock.RLock()
	defer kc.latestMpuTimestampLock.RUnlock()
	if t, ok := kc.latestMpuTimestamp[symbol]; ok {
		return t
	}
	return 0
}
