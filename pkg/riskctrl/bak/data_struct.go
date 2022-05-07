package riskctrlbak

import (
	"encoding/json"
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
	Volume         TVolume
	ExchangeVolume map[TExchange]TVolume
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
	Exchange TExchange             `json:"Exchange"`
	Symbol   TSymbol               `json:"Symbol"`
	Time     uint64                `json:"Time"`
	Asks     map[TPrice]InnerDepth `json:"Asks"`
	Bids     map[TPrice]InnerDepth `json:"Bids"`
}

func (depth_quote *DepthQuote) String(len int) string {

	res, _ := json.Marshal(*depth_quote)

	return string(res)
}
