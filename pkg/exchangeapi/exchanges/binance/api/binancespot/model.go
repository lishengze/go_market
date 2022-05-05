package binancespot

type GetExchangeInfoRsp struct {
	Timezone   string `json:"timezone"`
	ServerTime int64  `json:"serverTime"`
	RateLimits []struct {
	} `json:"rateLimits"`
	ExchangeFilters []interface{} `json:"exchangeFilters"`
	Symbols         []struct {
		Symbol                 string   `json:"symbol"`
		Status                 string   `json:"status"`
		BaseAsset              string   `json:"baseAsset"`
		BaseAssetPrecision     int64    `json:"baseAssetPrecision"`
		QuoteAsset             string   `json:"quoteAsset"`
		QuotePrecision         int64    `json:"quotePrecision"`
		QuoteAssetPrecision    int64    `json:"quoteAssetPrecision"`
		OrderTypes             []string `json:"orderTypes"`
		IcebergAllowed         bool     `json:"icebergAllowed"`
		OcoAllowed             bool     `json:"ocoAllowed"`
		IsSpotTradingAllowed   bool     `json:"isSpotTradingAllowed"`
		IsMarginTradingAllowed bool     `json:"isMarginTradingAllowed"`
		Filters                []struct {
			FilterType        string `json:"filterType"`
			MaxPrice          string `json:"maxPrice,omitempty"`
			MinPrice          string `json:"minPrice,omitempty"`
			TickSize          string `json:"tickSize,omitempty"`
			MaxQty            string `json:"maxQty,omitempty"`
			MinQty            string `json:"minQty,omitempty"`
			StepSize          string `json:"stepSize,omitempty"`
			Limit             int64  `json:"limit,omitempty"`
			MinNotional       string `json:"minNotional,omitempty"`
			MultiplierUp      string `json:"multiplierUp,omitempty"`
			MultiplierDown    string `json:"multiplierDown,omitempty"`
			MultiplierDecimal string `json:"multiplierDecimal,omitempty"`
		} `json:"filters"`
		Permissions []string `json:"permissions"`
	} `json:"symbols"`
}

type (
	GetDepthReq struct {
		Symbol string `param:"symbol,required"`
		Limit  int64  `param:"limit,default=100"`
	}

	GetDepthRsp struct {
		LastUpdateId int64       `json:"lastUpdateId"`
		Bids         [][2]string `json:"bids"`
		Asks         [][2]string `json:"asks"`
	}
)

// WsMarketStreamSendMsg 发送给 行情 ws 信道中的消息格式
type WsMarketStreamSendMsg struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	Id     int32    `json:"id"`
}

type WsDiffDepth struct {
	Stream string `json:"stream"`
	Data   struct {
		E       string      `json:"e"`
		E1      int64       `json:"E"`
		S       string      `json:"s"`
		StartId int64       `json:"U"`
		EndId   int64       `json:"u"`
		B       [][2]string `json:"b"`
		A       [][2]string `json:"a"`
	} `json:"data"`
}

type WsKline struct {
	Stream string `json:"stream"`
	Data   struct {
		E  string `json:"e"`
		E1 int64  `json:"E"`
		S  string `json:"s"`
		K  struct {
			T  int64  `json:"t"`
			T1 int64  `json:"T"`
			S  string `json:"s"`
			I  string `json:"i"`
			F  int64  `json:"f"`
			L1 int64  `json:"L"`
			O  string `json:"o"`
			C  string `json:"c"`
			H  string `json:"h"`
			L  string `json:"l"`
			V  string `json:"v"`
			N  int64  `json:"n"`
			X  bool   `json:"x"`
			Q  string `json:"q"`
			V1 string `json:"V"`
			Q1 string `json:"Q"`
			B  string `json:"B"`
		} `json:"k"`
	} `json:"data"`
}

type StreamMarketTrade struct {
	Stream string `json:"stream"`
	Data   struct {
		E  string `json:"e"`
		E1 int64  `json:"E"`
		S  string `json:"s"`
		T  int64  `json:"t"`
		P  string `json:"p"`
		Q  string `json:"q"`
		B  int64  `json:"b"`
		A  int64  `json:"a"`
		T1 int64  `json:"T"`
		M  bool   `json:"m"`
		M1 bool   `json:"M"`
	} `json:"data"`
}
