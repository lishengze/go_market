package ftxapi

import (
	"encoding/json"
	"time"
)

type (
	WsSubscribeMsg struct {
		Op      string `json:"op"`
		Channel string `json:"channel"`
		Market  string `json:"market"`
		//'op': 'subscribe', 'channel': 'trades', 'market': 'BTC-PERP'
	}

	WsUnsubscribeMsg WsSubscribeMsg
)

type StreamMarketTrade struct {
	Channel string `json:"channel"`
	Market  string `json:"market"`
	Type    string `json:"type"`
	Data    []struct {
		Id          int64     `json:"id"`
		Price       float64   `json:"price"`
		Size        float64   `json:"size"`
		Side        string    `json:"side"`
		Liquidation bool      `json:"liquidation"`
		Time        time.Time `json:"time"`
	} `json:"data"`
}

type StreamDepth struct {
	Channel string `json:"channel"`
	Market  string `json:"market"`
	Type    string `json:"type"`
	Data    struct {
		Time     float64         `json:"time"`
		Checksum uint32          `json:"checksum"`
		Bids     [][]json.Number `json:"bids"`
		Asks     [][]json.Number `json:"asks"`
		Action   string          `json:"action"`
	} `json:"data"`
}

type FloatString struct {
	s string
}

type GetMarketRsp struct {
	Success bool `json:"success"`
	Result  []struct {
		Name                  string  `json:"name"`
		BaseCurrency          string  `json:"baseCurrency"`
		QuoteCurrency         string  `json:"quoteCurrency"`
		QuoteVolume24H        float64 `json:"quoteVolume24h"`
		Change1H              float64 `json:"change1h"`
		Change24H             float64 `json:"change24h"`
		ChangeBod             float64 `json:"changeBod"`
		HighLeverageFeeExempt bool    `json:"highLeverageFeeExempt"`
		MinProvideSize        float64 `json:"minProvideSize"`
		Type                  string  `json:"type"`
		Underlying            string  `json:"underlying"`
		Enabled               bool    `json:"enabled"`
		Ask                   float64 `json:"ask"`
		Bid                   float64 `json:"bid"`
		Last                  float64 `json:"last"`
		PostOnly              bool    `json:"postOnly"`
		Price                 float64 `json:"price"`
		PriceIncrement        float64 `json:"priceIncrement"`
		SizeIncrement         float64 `json:"sizeIncrement"`
		Restricted            bool    `json:"restricted"`
		VolumeUsd24H          float64 `json:"volumeUsd24h"`
		LargeOrderThreshold   float64 `json:"largeOrderThreshold"`
	} `json:"result"`
}
