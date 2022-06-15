package ftxapi

import (
	"encoding/json"
	"time"
)

type (
	WsSubscribeMsg struct {
		Op      string `json:"op"`
		Channel string `json:"channel"`
		Market  string `json:"market,omitempty"`
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

type GetBalanceRsp struct {
	Success bool `json:"success"`
	Result  []struct {
		Coin                   string  `json:"coin"`
		Free                   float64 `json:"free"`
		SpotBorrow             float64 `json:"spotBorrow"`
		Total                  float64 `json:"total"`
		UsdValue               float64 `json:"usdValue"`
		AvailableWithoutBorrow float64 `json:"availableWithoutBorrow"`
	} `json:"result"`
}

type PlaceOrderReq struct {
	Market            string `json:"market"`
	Side              string `json:"side,options=buy|sell"`
	Price             string `json:"price"`
	Type              string `json:"type,options=limit|market"`
	Size              string `json:"size"`
	ClientId          string `json:"clientId,omitempty"`
	ReduceOnly        bool   `json:"reduceOnly,omitempty"`
	Ioc               bool   `json:"ioc,omitempty"`
	PostOnly          bool   `json:"postOnly,omitempty"`
	RejectOnPriceBand bool   `json:"rejectOnPriceBand,omitempty"`
	RejectAfterTs     string `json:"rejectAfterTs,omitempty"`
}

type PlaceOrderRsp struct {
	Success bool `json:"success"`
	Result  struct {
		CreatedAt     time.Time `json:"createdAt"`
		FilledSize    float64   `json:"filledSize"`
		Future        string    `json:"future"`
		Id            int       `json:"id"`
		Market        string    `json:"market"`
		Price         float64   `json:"price"`
		RemainingSize float64   `json:"remainingSize"`
		Side          string    `json:"side"`
		Size          float64   `json:"size"`
		Status        string    `json:"status"`
		Type          string    `json:"type"`
		ReduceOnly    bool      `json:"reduceOnly"`
		Ioc           bool      `json:"ioc"`
		PostOnly      bool      `json:"postOnly"`
		ClientId      string    `json:"clientId"`
	} `json:"result"`
}

type QueryOrderRsp struct {
	Success bool `json:"success"`
	Result  struct {
		CreatedAt     time.Time   `json:"createdAt"`
		FilledSize    float64     `json:"filledSize"`
		Future        string      `json:"future"`
		Id            int         `json:"id"`
		Market        string      `json:"market"`
		Price         float64     `json:"price"`
		AvgFillPrice  float64     `json:"avgFillPrice"`
		RemainingSize float64     `json:"remainingSize"`
		Side          string      `json:"side"`
		Size          float64     `json:"size"`
		Status        string      `json:"status"`
		Type          string      `json:"type"`
		ReduceOnly    bool        `json:"reduceOnly"`
		Ioc           bool        `json:"ioc"`
		PostOnly      bool        `json:"postOnly"`
		ClientId      string      `json:"clientId"`
		Liquidation   interface{} `json:"liquidation"`
	} `json:"result"`
}

type CancelOrderRsp struct {
	Success bool   `json:"success"`
	Result  string `json:"result"`
}

type QueryTradesReq struct {
	market    string `param:"market,omitempty"`
	StartTime int64  `param:"start_time,omitempty"`
	EndTime   int64  `param:"end_time,omitempty"`
	Order     string `param:"order,omitempty"`
	OrderId   string `param:"orderId,omitempty"`
}

type QueryTradesRsp struct {
	Success bool `json:"success"`
	Result  []struct {
		Fee           float64   `json:"fee"`
		FeeCurrency   string    `json:"feeCurrency"`
		FeeRate       float64   `json:"feeRate"`
		Future        string    `json:"future"`
		Id            int       `json:"id"`
		Liquidity     string    `json:"liquidity"`
		Market        string    `json:"market"`
		BaseCurrency  string    `json:"baseCurrency"`
		QuoteCurrency string    `json:"quoteCurrency"`
		OrderId       int       `json:"orderId"`
		TradeId       int       `json:"tradeId"`
		Price         float64   `json:"price"`
		Side          string    `json:"side"`
		Size          float64   `json:"size"`
		Time          time.Time `json:"time"`
		Type          string    `json:"type"`
	} `json:"result"`
}

type WsFills struct {
	Channel string `json:"channel"`
	Data    struct {
		Fee       float64   `json:"fee"`
		FeeRate   float64   `json:"feeRate"`
		Future    string    `json:"future"`
		Id        int       `json:"id"`
		Liquidity string    `json:"liquidity"`
		Market    string    `json:"market"`
		OrderId   int       `json:"orderId"`
		TradeId   int       `json:"tradeId"`
		Price     float64   `json:"price"`
		Side      string    `json:"side"`
		Size      float64   `json:"size"`
		Time      time.Time `json:"time"`
		Type      string    `json:"type"`
	} `json:"data"`
	Type string `json:"type"`
}

type WsOrders struct {
	Channel string `json:"channel"`
	Data    struct {
		Id            int       `json:"id"`
		ClientId      string    `json:"clientId"`
		Market        string    `json:"market"`
		Type          string    `json:"type"`
		Side          string    `json:"side"`
		Size          float64   `json:"size"`
		Price         float64   `json:"price"`
		ReduceOnly    bool      `json:"reduceOnly"`
		Ioc           bool      `json:"ioc"`
		PostOnly      bool      `json:"postOnly"`
		Status        string    `json:"status"`
		FilledSize    float64   `json:"filledSize"`
		RemainingSize float64   `json:"remainingSize"`
		AvgFillPrice  float64   `json:"avgFillPrice"`
		CreatedAt     time.Time `json:"createdAt"`
	} `json:"data"`
	Type string `json:"type"`
}

type WsLoginReq struct {
	WsLoginArgs `json:"args"`
	Op          string `json:"op"`
}

type WsLoginArgs struct {
	Key        string `json:"key"`
	Sign       string `json:"sign"`
	Time       int64  `json:"time"`
	SubAccount string `json:"subaccount,omitempty"`
}
