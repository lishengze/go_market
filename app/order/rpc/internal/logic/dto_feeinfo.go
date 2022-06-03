package logic

import "github.com/shopspring/decimal"

type FeeInfo struct {
	Symbol                            string `json:"symbol"`                  //币种代码，如 BTC_USDT
	SymbolKind                        int    `json:"symbol_kind"`             //品种类型，如 现货、期货等
	Underlying                        string `json:"underlying"`              //标的，如 BTC
	PrimaryCurrency                   string `json:"primary_currency"`        //基础货币，如 USDT
	BidCurrency                       string `json:"bid_currency"`            //报价货币，如 USDT
	SettleCurrency                    string `json:"settle_currency"`         //结算货币，如 USDT
	FeeKind                           int    `json:"fee_kind"`                //手续费算法，取值1或2，1表示百比分，2表示绝对值。默认为1.
	TakerFee                          string `json:"taker_fee"`               //Taker手续费率
	MakerFee                          string `json:"maker_fee"`               //Maker手续费率
	BuyCommissionSwitch               bool   `json:"buy_commission_switch"`   //买委托限制开关,true表示禁止
	SellCommissionSwitch              bool   `json:"sell_commission_switch"`  //卖委托限制开关,true表示禁止
	Switch                            bool   `json:"switch"`                  //交易开关
	VolumePrecision                   int    `json:"volume_precision"`        //数量精度
	PricePrecision                    int    `json:"price_precision"`         //价格精度
	AmountPrecision                   int    `json:"amount_precision"`        //金额精度
	MinUnit                           string `json:"min_unit"`                //最小交易单位
	MinChangePrice                    string `json:"min_change_price"`        //最小变动价位
	Spread                            string `json:"spread"`                  //点差，品种tick值的整数倍
	SingleMinOrderVolume              string `json:"single_min_order_volume"` //单次最小下单量
	SingleMaxOrderVolume              string `json:"single_max_order_volume"` //单次最大下单量
	SingleMinOrderAmount              string `json:"single_min_order_amount"` //单次最小下单金额
	SingleMaxOrderAmount              string `json:"single_max_order_amount"` //单次最大下单金额
	BuyPriceLimit                     string `json:"buy_price_limit"`         //买委托价格限制
	SellPriceLimit                    string `json:"sell_price_limit"`        //卖委托价格限制
	MaxMatchGear                      int    `json:"max_match_gear"`          //最大成交档位,不得超过20
	OtcMinOrderVolume                 string `json:"otc_min_order_volume"`    //OTC最小量
	OtcMaxOrderVolume                 string `json:"otc_max_order_volume"`    //OTC最大量
	OtcMinOrderAmount                 string `json:"otc_min_order_amount"`    //OTC最小金额
	OtcMaxOrderAmount                 string `json:"otc_max_order_amount"`    //OTC最大金额
	DailyTradingLimitSwitch           bool
	DailyBuyingLimitPerClientWeekday  decimal.Decimal //每人每日买交易限额
	DailyBuyingLimitPerClientWeekend  decimal.Decimal //每人每日买交易限额
	DailySellingLimitPerClientWeekday decimal.Decimal //每人每日卖交易限额
	DailySellingLimitPerClientWeekend decimal.Decimal //每人每日卖交易限额
}

type ReqQueryFee struct {
	RequestID string `form:"request_id" json:"request_id"`
	Symbol    string `form:"symbol" json:"symbol"` //币对
}

type RspQueryFee []*FeeInfo
