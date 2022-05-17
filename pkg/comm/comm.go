package comm

import (
	config "market_aggregate/pkg/conf"
	"market_aggregate/pkg/datastruct"
	"market_aggregate/pkg/kafka"
)

type CommI interface {
	Init(*datastruct.DataChannel, *datastruct.DataChannel)

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

func (c *Comm) Init(
	recv_chan *datastruct.DataChannel,
	pub_chan *datastruct.DataChannel) error {

	if config.NATIVE_CONFIG().SerialType == COMM_PROTOBUF {
		c.Serializer = &ProtobufSerializer{}
	}

	if config.NATIVE_CONFIG().NetServerType == COMM_KAFKA {
		c.NetServer = &kafka.KafkaServer{}
		c.NetServer.Init(c.Serializer, recv_chan, pub_chan)
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
