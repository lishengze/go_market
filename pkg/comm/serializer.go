package comm

import (
	"fmt"
	"market_aggregate/pkg/datastruct"
	"strconv"

	"google.golang.org/protobuf/proto"
)

type SerializerI interface {
	EncodeDepth(*datastruct.DepthQuote) (string, error)
	EncodeKline(*datastruct.Kline) (string, error)
	EncodeTrade(*datastruct.Trade) (string, error)

	DecodeDepth([]byte) (*datastruct.DepthQuote, error)
	DecodeKline([]byte) (*datastruct.Kline, error)
	DecodeTrade([]byte) (*datastruct.Trade, error)
}

type ProtobufSerializer struct {
}

func (p *ProtobufSerializer) EncodeDepth(*datastruct.DepthQuote) (string, error) {
	return "", nil
}

func (p *ProtobufSerializer) EncodeKline(*datastruct.Kline) (string, error) {
	return "", nil
}

func (p *ProtobufSerializer) EncodeTrade(*datastruct.Trade) (string, error) {
	return "", nil
}

func (p *ProtobufSerializer) DecodeDepth(raw_msg []byte) (*datastruct.DepthQuote, error) {

}

func (p *ProtobufSerializer) DecodeKline(raw_msg []byte) (*datastruct.Kline, error) {

}

func (p *ProtobufSerializer) DecodeTrade(raw_msg []byte) (*datastruct.Trade, error) {
	proto_trade := Trade{}
	err := proto.Unmarshal(raw_msg, &proto_trade)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	price, err := strconv.ParseFloat(proto_trade.Price, 64)
	if err != nil {
		return nil, err
	}

	volume, err := strconv.ParseFloat(proto_trade.Volume, 64)
	if err != nil {
		return nil, err
	}

	return &datastruct.Trade{
		Exchange: proto_trade.Exchange,
		Symbol:   proto_trade.Symbol,
		Price:    price,
		Volume:   volume,
		Time:     int64(proto_trade.Timestamp.Nanos),
	}, nil

}
