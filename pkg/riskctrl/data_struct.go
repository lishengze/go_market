package riskctrl

import (
	"encoding/json"

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

func (d *DepthQuote) Init() {
	d.Asks = treemap.NewWith(utils.Float64Comparator)
	d.Bids = treemap.NewWith(utils.Float64Comparator)
}

func (depth_quote *DepthQuote) String(len int) string {

	res, _ := json.Marshal(*depth_quote)

	return string(res)
}
