package riskctrl

import (
	"fmt"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
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
	open     float64
	high     float64
	low      float64
	close    float64
	volume   float64
}

func (d *DepthQuote) Init() {
	d.Asks = treemap.NewWith(utils.Float64Comparator)
	d.Bids = treemap.NewWith(utils.Float64Comparator)
}

func (d *DepthQuote) String(len int) string {

	res := fmt.Sprintf("%s.%s, %d\nAsks: %s\nBids: %s \n", d.Exchange, d.Symbol, d.Time, d.Asks.String(), d.Bids.String())

	return string(res)
}
