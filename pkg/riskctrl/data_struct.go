package riskctrl

import (
	"fmt"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
)

const (
	BCTS_EXCHANGE = "_bcts_"
)

type TestData struct {
	Name string
}

type TSymbol string
type TExchange string

// type RFloat float64
type TPrice float64
type TVolume float64

type InnerDepth struct {
	Volume         float64
	ExchangeVolume map[string]float64
}

func (src *InnerDepth) Add(other *InnerDepth) {
	if src == other {
		return
	}

	src.Volume += other.Volume

	for exchange, volume := range other.ExchangeVolume {
		src.ExchangeVolume[exchange] += volume
	}
}

type DepthQuote struct {
	Exchange string
	Symbol   string
	Time     int64
	Asks     *treemap.Map
	Bids     *treemap.Map
}

type Kline struct {
	Exchange   string
	Symbol     string
	Time       int64
	Open       float64
	High       float64
	Low        float64
	Close      float64
	Volume     float64
	Resolution int
}

type Trade struct {
	Exchange string
	Symbol   string
	Time     int64
	Price    float64
	Volume   float64
}

func (t *Trade) String() string {
	res := fmt.Sprintf("%s.%s, %+v, p: %f v: %f \n", t.Exchange, t.Symbol, time.Unix(int64(t.Time), 0), t.Price, t.Volume)
	return res
}

func (k *Kline) String() string {
	res := fmt.Sprintf("%s.%s, %+v, o: %f, h: %f, l: %f, c: %f, v: %f\n",
		k.Exchange, k.Symbol, time.Unix(int64(k.Time), 0),
		k.Open, k.High, k.Low, k.Close, k.Volume)
	return res
}

func (d *DepthQuote) String(len int) string {

	res := fmt.Sprintf("%s.%s, %v\nAsks: %s\nBids: %s \n", d.Exchange, d.Symbol, time.Unix(int64(d.Time), 0), d.Asks.String(), d.Bids.String())

	return string(res)
}

func NewTrade(src *Trade) *Trade {
	if src != nil {
		rst := &Trade{
			Exchange: src.Exchange,
			Symbol:   src.Symbol,
			Time:     src.Time,
			Price:    src.Price,
			Volume:   src.Volume,
		}
		return rst
	} else {
		rst := &Trade{}
		return rst
	}
}

func NewKline(src *Kline) *Kline {
	if src != nil {
		rst := &Kline{
			Exchange:   src.Exchange,
			Symbol:     src.Symbol,
			Time:       src.Time,
			Open:       src.Open,
			High:       src.High,
			Low:        src.Low,
			Close:      src.Close,
			Volume:     src.Volume,
			Resolution: src.Resolution,
		}
		return rst
	} else {
		rst := &Kline{}
		return rst
	}
}

func InitKlineByTrade(src *Kline, trade *Trade) {
	src.Exchange = BCTS_EXCHANGE
	src.Symbol = trade.Symbol
	src.Time = trade.Time
	src.Resolution = 60
	src.Open = trade.Price
	src.High = trade.Price
	src.Low = trade.Price
	src.Close = trade.Price
	src.Volume = trade.Volume
}

func NewDepth(src *DepthQuote) *DepthQuote {
	if src != nil {
		rst := &DepthQuote{
			Exchange: src.Exchange,
			Symbol:   src.Symbol,
			Time:     src.Time,
			Asks:     treemap.NewWith(utils.Float64Comparator),
			Bids:     treemap.NewWith(utils.Float64Comparator),
		}

		ask_iter := src.Asks.Iterator()
		for ask_iter.Begin(); ask_iter.Next(); {
			rst.Asks.Put(ask_iter.Key(), ask_iter.Value())
		}

		bid_iter := src.Bids.Iterator()
		for bid_iter.Begin(); bid_iter.Next(); {
			rst.Bids.Put(bid_iter.Key(), bid_iter.Value())
		}

		return rst
	} else {
		rst := &DepthQuote{
			Asks: treemap.NewWith(utils.Float64Comparator),
			Bids: treemap.NewWith(utils.Float64Comparator),
		}
		return rst
	}
}

type DataChannel struct {
	TradeChannel chan *Trade
	KlineChannel chan *Kline
	DepthChannel chan *DepthQuote
}

func (d *DepthQuote) Init() {
	d.Asks = treemap.NewWith(utils.Float64Comparator)
	d.Bids = treemap.NewWith(utils.Float64Comparator)
}
