package comm

import (
	"market_aggregate/pkg/conf"
	"market_aggregate/pkg/datastruct"
	"market_aggregate/pkg/kafka"
)

type CommI interface {
	Init(*conf.Config, *datastruct.DataChannel)

	PublishDepth(*datastruct.DepthQuote)
	PublishKline(*datastruct.Kline)
	PublishTrade(*datastruct.Trade)

	// SendDepth(*datastruct.DepthQuote)
	// SendKline(*datastruct.Kline)
	// SendTrade(*datastruct.Trade)
}

const (
	COMM_KAFKA = "KAFKA"
	COMM_REDIS = "REDIS"
	COMM_GRPC  = "GRPC"
)

const (
	COMM_PROTOBUF = "PROTOBUF"
	COMM_JSON     = "JSON"
	COMM_FLATBUF  = "FLATBUF"
)

type Comm struct {
	RecvCommType   string // Kafka, Redis, Grpc
	PubCommType    string //
	SerializerType string // Protobuf, Json,

	NetServer  datastruct.NetServerI  //
	Serializer datastruct.SerializerI //
}

func (c *Comm) Init(config *conf.Config, recv_chan *datastruct.DataChannel) error {

	if config.SerialType == COMM_PROTOBUF {
		c.Serializer = &ProtobufSerializer{}
	}

	if config.NetServerType == COMM_KAFKA {
		c.NetServer = &kafka.KafkaServer{
			RecvDataChan: recv_chan,
			Serializer:   c.Serializer,
			Config:       config,
		}
		c.NetServer.Init(config, c.Serializer, recv_chan)
	}

	// c.DataChan = data_chan
	return nil
}

func (c *Comm) PublishDepth(depth *datastruct.DepthQuote) {
	c.NetServer.PublishDepth(depth)
}

func (c *Comm) PublishKline(kline *datastruct.Kline) {
	c.NetServer.PublishKline(kline)
}

func (c *Comm) PublishTrade(trade *datastruct.Trade) {
	c.NetServer.PublishTrade(trade)
}

// func (c *Comm) SendDepth(depth *datastruct.DepthQuote) {
// 	c.DataChan.SendDepth(depth)
// }

// func (c *Comm) SendKline(kline *datastruct.Kline) {
// 	c.DataChan.SendKline(kline)
// }

// func (c *Comm) SendTrade(trade *datastruct.Trade) {
// 	c.DataChan.SendTrade(trade)
// }
