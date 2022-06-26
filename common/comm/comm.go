package comm

import (
	"market_server/common/config"
	"market_server/common/datastruct"
	"market_server/common/kafka"

	"github.com/zeromicro/go-zero/core/logx"
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
	cfg config.CommConfig) *Comm {
	logx.Infof("NewComm, Config: %+v\n", cfg)
	c := &Comm{}

	if cfg.SerialType == COMM_PROTOBUF {
		c.Serializer = &ProtobufSerializer{}
	}

	if cfg.NetServerType == COMM_KAFKA {
		// c.NetServer = &kafka.KafkaServer{}
		// c.NetServer.InitKafka(c.Serializer, recv_chan, pub_chan, cfg.KafkaConfig)
		var err error
		c.NetServer, err = kafka.NewKafka(c.Serializer, recv_chan, pub_chan, cfg.KafkaConfig)

		if err != nil {
			logx.Errorf("NewKafka Error: %+v", err)
		}
	}

	return c
}

func (c *Comm) Init(
	recv_chan *datastruct.DataChannel,
	pub_chan *datastruct.DataChannel,
	cfg config.CommConfig) error {

	if cfg.SerialType == COMM_PROTOBUF {
		c.Serializer = &ProtobufSerializer{}
	}

	if cfg.NetServerType == COMM_KAFKA {
		var err error
		c.NetServer, err = kafka.NewKafka(c.Serializer, recv_chan, pub_chan, cfg.KafkaConfig)
		logx.Errorf("NewKafka Error: %+v", err)

		// c.NetServer.InitKafka(c.Serializer, recv_chan, pub_chan, cfg.KafkaConfig)
	}

	// c.DataChan = data_chan
	return nil
}

func (c *Comm) Start() {
	if c.NetServer != nil {
		c.NetServer.Start()
	} else {
		logx.Error("c.NetServer is nil")
	}
}

func (c *Comm) UpdateMetaData(meta *datastruct.Metadata) {
	if c.NetServer != nil {
		logx.Info("Commer Update MetaData")
		go c.NetServer.UpdateMetaData(meta)
	} else {
		logx.Error("c.NetServer is nil")
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
