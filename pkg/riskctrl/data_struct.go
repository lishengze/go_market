package riskctrl

import (
	"fmt"

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
	Time     uint64
	Asks     *treemap.Map
	Bids     *treemap.Map
}

type Kline struct {
	Exchange string
	Symbol   string
	Time     uint64
	Open     float64
	High     float64
	Low      float64
	Close    float64
	Volume   float64
}

type Trade struct {
	Exchange string
	Symbol   string
	Time     uint64
	Price    float64
	Volume   float64
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
			Exchange: src.Exchange,
			Symbol:   src.Symbol,
			Time:     src.Time,
			Open:     src.Open,
			High:     src.High,
			Low:      src.Low,
			Close:    src.Close,
			Volume:   src.Volume,
		}
		return rst
	} else {
		rst := &Kline{}
		return rst
	}
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

func (d *DepthQuote) String(len int) string {

	res := fmt.Sprintf("%s.%s, %d\nAsks: %s\nBids: %s \n", d.Exchange, d.Symbol, d.Time, d.Asks.String(), d.Bids.String())

	return string(res)
}
