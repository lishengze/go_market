package comm

import (
	"market_aggregate/pkg/conf"
	"market_aggregate/pkg/datastruct"
)

type DataRecvI interface {
	SendDepth(*datastruct.DepthQuote)
	SendKline(*datastruct.Kline)
	SendTrade(*datastruct.Trade)
}

type CommI interface {
	Init(*conf.Config, DataRecvI)

	PublishDepth(*datastruct.DepthQuote)
	PublishKline(*datastruct.Kline)
	PublishTrade(*datastruct.Trade)

	SendDepth(*datastruct.DepthQuote)
	SendKline(*datastruct.Kline)
	SendTrade(*datastruct.Trade)
}

const (
	COMM_KAFKA = "KAFKA"
	COMM_REDIS = "REDIS"
	COMM_GRPC  = "GRPC"
)

type Comm struct {
	RecvCommType   string // Kafka, Redis, Grpc
	PubCommType    string //
	SerializerType string // Protobuf, Json,

	NetServer   NetServerI  //
	SerializerI SerializerI //

	DataChan DataRecvI
}

func (c *Comm) Init(config *conf.Config, data_chan DataRecvI) {
	c.DataChan = data_chan
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

func (c *Comm) SendDepth(depth *datastruct.DepthQuote) {
	c.DataChan.SendDepth(depth)
}

func (c *Comm) SendKline(kline *datastruct.Kline) {
	c.DataChan.SendKline(kline)
}

func (c *Comm) SendTrade(trade *datastruct.Trade) {
	c.DataChan.SendTrade(trade)
}
