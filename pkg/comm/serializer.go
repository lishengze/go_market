package comm

import (
	"fmt"
	"market_aggregate/app/datastruct"
	"market_aggregate/app/protostruct"
	"market_aggregate/app/util"
	"strconv"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
	"google.golang.org/protobuf/proto"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type ProtobufSerializer struct {
}

func SetProtoDepth(dst *[]*protostruct.PriceVolume, src *treemap.Map) {
	iter := src.Iterator()
	for iter.Begin(); iter.Next(); {
		cur_pricevolume := protostruct.PriceVolume{}
		cur_pricevolume.Price = strconv.FormatFloat(iter.Key().(float64), 'f', -1, 64)
		cur_pricevolume.Volume = strconv.FormatFloat(iter.Value().(*datastruct.InnerDepth).Volume, 'f', -1, 64)
		*dst = append(*dst, &cur_pricevolume)
	}
}

func (p *ProtobufSerializer) EncodeDepth(local_depth *datastruct.DepthQuote) ([]byte, error) {
	proto_depth := protostruct.Depth{}

	proto_depth.Exchange = local_depth.Exchange
	proto_depth.Symbol = local_depth.Symbol
	proto_depth.Timestamp = &timestamppb.Timestamp{Seconds: local_depth.Time / datastruct.NANO_PER_SECS, Nanos: int32(local_depth.Time % datastruct.NANO_PER_SECS)}
	proto_depth.MpuTimestamp = &timestamppb.Timestamp{Seconds: local_depth.Time / datastruct.NANO_PER_SECS, Nanos: int32(local_depth.Time % datastruct.NANO_PER_SECS)}

	SetProtoDepth(&proto_depth.Asks, local_depth.Asks)
	SetProtoDepth(&proto_depth.Bids, local_depth.Bids)

	// fmt.Printf("ProtoDepth: %+v \n", proto_depth)

	msg, err := proto.Marshal(&proto_depth)

	// proto_depth.
	return msg, err
}

func (p *ProtobufSerializer) EncodeKline(local_kline *datastruct.Kline) ([]byte, error) {

	proto_kline := protostruct.Kline{}
	proto_kline.Exchange = local_kline.Exchange
	proto_kline.Symbol = local_kline.Symbol
	proto_kline.Timestamp = &timestamppb.Timestamp{Seconds: local_kline.Time / datastruct.NANO_PER_SECS, Nanos: int32(local_kline.Time % datastruct.NANO_PER_SECS)}
	proto_kline.Resolution = uint32(local_kline.Resolution)

	proto_kline.Open = strconv.FormatFloat(local_kline.Open, 'f', -1, 64)
	proto_kline.High = strconv.FormatFloat(local_kline.High, 'f', -1, 64)
	proto_kline.Low = strconv.FormatFloat(local_kline.Low, 'f', -1, 64)
	proto_kline.Close = strconv.FormatFloat(local_kline.Close, 'f', -1, 64)

	proto_kline.Volume = strconv.FormatFloat(local_kline.Volume, 'f', -1, 64)

	// fmt.Printf("proto_kline: %+v \n", proto_kline)

	msg, err := proto.Marshal(&proto_kline)

	return msg, err
}

func (p *ProtobufSerializer) EncodeTrade(local_trade *datastruct.Trade) ([]byte, error) {
	proto_trade := protostruct.Trade{}

	proto_trade.Exchange = local_trade.Exchange
	proto_trade.Symbol = local_trade.Symbol
	proto_trade.Timestamp = &timestamppb.Timestamp{Seconds: local_trade.Time / datastruct.NANO_PER_SECS, Nanos: int32(local_trade.Time % datastruct.NANO_PER_SECS)}
	proto_trade.Price = strconv.FormatFloat(local_trade.Price, 'f', -1, 64)
	proto_trade.Volume = strconv.FormatFloat(local_trade.Volume, 'f', -1, 64)

	msg, err := proto.Marshal(&proto_trade)

	// fmt.Printf("\nproto_trade: %+v \n", proto_trade)

	return msg, err
}

func SetDepthTreeMap(src *treemap.Map, proto_depth []*protostruct.PriceVolume, exchange string) {
	for _, value := range proto_depth {
		price, err := strconv.ParseFloat(value.Price, 64)
		if err != nil {
			util.LOG_ERROR(err.Error())
			continue
		}
		volume, err := strconv.ParseFloat(value.Volume, 64)
		if err != nil {
			util.LOG_ERROR(err.Error())
			continue
		}

		inner_depth := datastruct.InnerDepth{
			Volume:         volume,
			ExchangeVolume: make(map[string]float64),
		}
		inner_depth.ExchangeVolume[exchange] = volume

		src.Put(price, &inner_depth)
	}
}

func (p *ProtobufSerializer) DecodeDepth(raw_msg []byte) (*datastruct.DepthQuote, error) {
	proto_depth := protostruct.Depth{}
	err := proto.Unmarshal(raw_msg, &proto_depth)
	if err != nil {
		util.LOG_ERROR(err.Error())
		return nil, err
	}

	asks := treemap.NewWith(utils.Float64Comparator)
	bids := treemap.NewWith(utils.Float64Comparator)

	SetDepthTreeMap(asks, proto_depth.Asks, proto_depth.Exchange)
	SetDepthTreeMap(bids, proto_depth.Bids, proto_depth.Exchange)

	return &datastruct.DepthQuote{
		Exchange: proto_depth.Exchange,
		Symbol:   proto_depth.Symbol,
		Time:     int64(proto_depth.Timestamp.Seconds*datastruct.NANO_PER_SECS + int64(proto_depth.Timestamp.Nanos)),
		Asks:     asks,
		Bids:     bids,
	}, nil
}

func (p *ProtobufSerializer) DecodeKline(raw_msg []byte) (*datastruct.Kline, error) {
	proto_kline := protostruct.Kline{}
	err := proto.Unmarshal(raw_msg, &proto_kline)
	if err != nil {
		util.LOG_ERROR(err.Error())
		return nil, err
	}

	open, err := strconv.ParseFloat(proto_kline.Open, 64)
	if err != nil {
		util.LOG_ERROR(err.Error())
		return nil, err
	}

	high, err := strconv.ParseFloat(proto_kline.High, 64)
	if err != nil {
		util.LOG_ERROR(err.Error())
		return nil, err
	}

	low, err := strconv.ParseFloat(proto_kline.Low, 64)
	if err != nil {
		util.LOG_ERROR(err.Error())
		return nil, err
	}

	close, err := strconv.ParseFloat(proto_kline.Close, 64)
	if err != nil {
		util.LOG_ERROR(err.Error())
		return nil, err
	}

	volume, err := strconv.ParseFloat(proto_kline.Volume, 64)
	if err != nil {
		util.LOG_ERROR(err.Error())
		return nil, err
	}

	return &datastruct.Kline{
		Exchange:   proto_kline.Exchange,
		Symbol:     proto_kline.Symbol,
		Time:       int64(proto_kline.Timestamp.Seconds*datastruct.NANO_PER_SECS + int64(proto_kline.Timestamp.Nanos)),
		Resolution: int(proto_kline.Resolution),
		Open:       open,
		High:       high,
		Low:        low,
		Close:      close,
		Volume:     volume,
	}, nil
}

func (p *ProtobufSerializer) DecodeTrade(raw_msg []byte) (*datastruct.Trade, error) {
	proto_trade := protostruct.Trade{}
	err := proto.Unmarshal(raw_msg, &proto_trade)
	if err != nil {
		util.LOG_ERROR(err.Error())
		return nil, err
	}

	price, err := strconv.ParseFloat(proto_trade.Price, 64)
	if err != nil {
		util.LOG_ERROR(err.Error())
		return nil, err
	}

	volume, err := strconv.ParseFloat(proto_trade.Volume, 64)
	if err != nil {
		util.LOG_ERROR(err.Error())
		return nil, err
	}

	// fmt.Printf("proto_trade.Timestamp: %+v\n\n", proto_trade.Timestamp)

	return &datastruct.Trade{
		Exchange: proto_trade.Exchange,
		Symbol:   proto_trade.Symbol,
		Price:    price,
		Volume:   volume,
		Time:     int64(proto_trade.Timestamp.Seconds*datastruct.NANO_PER_SECS + int64(proto_trade.Timestamp.Nanos)),
	}, nil
}

func TestSeDepth() {
	test_depth := datastruct.GetTestDepth()
	fmt.Printf("OriginalDepth: %+v\n", test_depth)

	PS := ProtobufSerializer{}

	bytes, _ := PS.EncodeDepth(test_depth)

	reco_depth, _ := PS.DecodeDepth(bytes)

	fmt.Printf("RecoDepth: %+v\n", reco_depth)

}

func TestSeTrade() {
	original_trade := datastruct.GetTestTrade()

	fmt.Printf("\nOriginalTrade: %+v\n", original_trade)

	PS := ProtobufSerializer{}

	bytes, _ := PS.EncodeTrade(original_trade)

	reco_trade, _ := PS.DecodeTrade(bytes)

	fmt.Printf("\nRecoTrade: %+v\n", reco_trade)

}

func TestSeKline() {
	original_kline := datastruct.GetTestKline()

	fmt.Printf("\nOriginalKline: %+v\n", original_kline)

	PS := ProtobufSerializer{}

	bytes, _ := PS.EncodeKline(original_kline)

	reco_kline, _ := PS.DecodeKline(bytes)

	fmt.Printf("\nRecoKline: %+v\n", reco_kline)
}
