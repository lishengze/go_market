package comm

import (
	"fmt"
	"market_aggregate/pkg/datastruct"
	"market_aggregate/pkg/riskctrl"
	"strconv"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
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

func SetDepthTreeMap(src *treemap.Map, proto_depth []*PriceVolume) {
	for _, value := range proto_depth {
		price, err := strconv.ParseFloat(value.Price, 64)
		if err != nil {
			riskctrl.LOG_ERROR(err.Error())
			continue
		}
		volume, err := strconv.ParseFloat(value.Volume, 64)
		if err != nil {
			riskctrl.LOG_ERROR(err.Error())
			continue
		}
		src.Put(price, volume)
	}
}

func (p *ProtobufSerializer) DecodeDepth(raw_msg []byte) (*datastruct.DepthQuote, error) {
	proto_depth := Depth{}
	err := proto.Unmarshal(raw_msg, &proto_depth)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	asks := treemap.NewWith(utils.Float64Comparator)
	bids := treemap.NewWith(utils.Float64Comparator)

	SetDepthTreeMap(asks, proto_depth.Asks)
	SetDepthTreeMap(bids, proto_depth.Bids)

	// 	// for item

	return &datastruct.DepthQuote{
		Exchange: proto_depth.Exchange,
		Symbol:   proto_depth.Symbol,
		Time:     int64(proto_depth.Timestamp.Nanos),
		Asks:     asks,
		Bids:     bids,
	}, nil
}

func (p *ProtobufSerializer) DecodeKline(raw_msg []byte) (*datastruct.Kline, error) {
	return nil, nil
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
