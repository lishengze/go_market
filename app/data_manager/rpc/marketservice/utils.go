package marketservice

import (
	"market_server/app/data_manager/rpc/types/pb"
	"market_server/common/datastruct"
	"strconv"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

func NewKlineWithPbKline(pb_kline *Kline) *datastruct.Kline {
	open, err := strconv.ParseFloat(pb_kline.Open, 64)
	if err != nil {
		logx.Error(err.Error())
		return nil
	}

	high, err := strconv.ParseFloat(pb_kline.High, 64)
	if err != nil {
		logx.Error(err.Error())
		return nil
	}

	low, err := strconv.ParseFloat(pb_kline.Low, 64)
	if err != nil {
		logx.Error(err.Error())
		return nil
	}

	close, err := strconv.ParseFloat(pb_kline.Close, 64)
	if err != nil {
		logx.Error(err.Error())
		return nil
	}

	volume, err := strconv.ParseFloat(pb_kline.Volume, 64)
	if err != nil {
		logx.Error(err.Error())
		return nil
	}

	return &datastruct.Kline{
		Exchange:   pb_kline.Exchange,
		Symbol:     pb_kline.Symbol,
		Time:       int64(pb_kline.Timestamp.Seconds*datastruct.NANO_PER_SECS + int64(pb_kline.Timestamp.Nanos)),
		Resolution: int(pb_kline.Resolution),
		Open:       open,
		High:       high,
		Low:        low,
		Close:      close,
		Volume:     volume,
	}
}

func NewKlineWithPbTrade(pb_trade *Trade) *datastruct.Trade {
	price, err := strconv.ParseFloat(pb_trade.Price, 64)
	if err != nil {
		logx.Error(err.Error())
		return nil
	}

	volume, err := strconv.ParseFloat(pb_trade.Volume, 64)
	if err != nil {
		logx.Error(err.Error())
		return nil
	}

	return &datastruct.Trade{
		Exchange: pb_trade.Exchange,
		Symbol:   pb_trade.Symbol,
		Price:    price,
		Volume:   volume,
		Time:     int64(pb_trade.Timestamp.Seconds*datastruct.NANO_PER_SECS + int64(pb_trade.Timestamp.Nanos)),
	}
}

func SetDepthTreeMap(src *treemap.Map, proto_depth []*pb.PriceVolume, exchange string) {
	for _, value := range proto_depth {
		price, err := strconv.ParseFloat(value.Price, 64)
		if err != nil {
			logx.Error(err.Error())
			continue
		}
		volume, err := strconv.ParseFloat(value.Volume, 64)
		if err != nil {
			logx.Error(err.Error())
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

func NewKlineWithPbDepth(pb_depth *Depth) *datastruct.DepthQuote {
	asks := treemap.NewWith(utils.Float64Comparator)
	bids := treemap.NewWith(utils.Float64Comparator)

	SetDepthTreeMap(asks, pb_depth.Asks, pb_depth.Exchange)
	SetDepthTreeMap(bids, pb_depth.Bids, pb_depth.Exchange)

	return &datastruct.DepthQuote{
		Exchange: pb_depth.Exchange,
		Symbol:   pb_depth.Symbol,
		Time:     int64(pb_depth.Timestamp.Seconds*datastruct.NANO_PER_SECS + int64(pb_depth.Timestamp.Nanos)),
		Asks:     asks,
		Bids:     bids,
	}
}
