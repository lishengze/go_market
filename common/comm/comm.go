package comm

import (
	"market_server/common/datastruct"
	"market_server/common/kafka"
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

func NewComm(recv_chan *datastruct.DataChannel,
	pub_chan *datastruct.DataChannel,
	serial_type string,
	net_server_type string) *Comm {
	c := &Comm{}

	if serial_type == COMM_PROTOBUF {
		c.Serializer = &ProtobufSerializer{}
	}

	if net_server_type == COMM_KAFKA {
		c.NetServer = &kafka.KafkaServer{}
		c.NetServer.Init(c.Serializer, recv_chan, pub_chan)
	}

	return c
}

func (c *Comm) Init(
	recv_chan *datastruct.DataChannel,
	pub_chan *datastruct.DataChannel,
	serial_type string,
	net_server_type string) error {

	if serial_type == COMM_PROTOBUF {
		c.Serializer = &ProtobufSerializer{}
	}

	if net_server_type == COMM_KAFKA {
		c.NetServer = &kafka.KafkaServer{}
		c.NetServer.Init(c.Serializer, recv_chan, pub_chan)
	}

	// c.DataChan = data_chan
	return nil
}

func (c *Comm) Start() {
	if c.NetServer != nil {
		c.NetServer.Start()
	}
}

func (c *Comm) UpdateMetaData(meta *datastruct.Metadata) {
	if c.NetServer != nil {
		c.NetServer.UpdateMetaData(meta)
	}
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
