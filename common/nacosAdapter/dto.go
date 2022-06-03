package nacosAdapter

type Currency struct {
	CurrencyID     string `json:"currency_id"`      //币种代码，如 USDT
	Kind           int    `json:"kind"`             //币种类型，1:法币 2: 公链数字货币 3：稳定币
	ChineseName    string `json:"chinese_name"`     //中文名称，如 泰达币
	EnglishName    string `json:"english_name"`     //英文名称，如 Tether USD
	MinUnit        string `json:"min_unit"`         //币种最小单位
	DepositSwitch  bool   `json:"deposit_switch"`   //充值开关
	WithdrawSwitch bool   `json:"withdraw_switch"`  //提现开关
	MinWithdraw    string `json:"min_withdraw"`     //单笔最小提现金额
	MaxWithdraw    string `json:"max_withdraw"`     //单笔最大提现金额
	MaxDayWithdraw string `json:"max_day_withdraw"` //当日最大提现金额
	Threshold      string `json:"threshold"`        //大额提现阈值
	FeeKind        int    `json:"fee_kind"`         //提现手续费算法 1:百比分，2:绝对值。默认为1.
	Fee            string `json:"fee"`              //提现手续费率
	Operator       string `json:"operator"`         //操作员
	Time           string `json:"time"`             //操作时间
}

type Hedging struct {
	PlatformID           string `json:"platform_id"`             //交易平台，如 binance
	Symbol               string `json:"symbol"`                  //品种
	Switch               bool   `json:"switch"`                  //充值开关
	BuyFundRatio         string `json:"buy_fund_ratio"`          //买盘资金占用比例
	SellFundRatio        string `json:"sell_fund_ratio"`         //卖盘资金占用比例
	MinUnit              string `json:"min_unit"`                //最小交易单位
	MinChangePrice       string `json:"min_change_price"`        //最小变动价
	MaxLeverage          int    `json:"max_leverage"`            //最大杠杆倍数
	SingleMaxOrderVolume string `json:"single_max_order_volume"` //单次最大下单量
	BuyPriceLimit        string `json:"buy_price_limit"`         //买委托价格限制
	SellPriceLimit       string `json:"sell_price_limit"`        //卖委托价格限制
	MaxMatchGear         int    `json:"max_match_gear"`          //最大成交档位
	FeeKind              int    `json:"fee_kind"`                //手续费算法，取值1或2，1表示百比分，2表示绝对值。默认为1.
	TakerFee             string `json:"taker_fee"`               //Taker手续费率
	MakerFee             string `json:"maker_fee"`               //Maker手续费率
}

type Symbol struct {
	SymbolID             string `json:"symbol_id"`               //品种代码，如 BTC_USDT
	SymbolKind           int    `json:"symbol_kind"`             //品种类型，如 1-现货、2-期货等
	Underlying           string `json:"underlying"`              //标的，如 BTC
	PrimaryCurrency      string `json:"primary_currency"`        //基础货币，如 USDT
	BidCurrency          string `json:"bid_currency"`            //报价货币，如 USDT
	SettleCurrency       string `json:"settle_currency"`         //结算货币，如 USDT
	Switch               bool   `json:"switch"`                  //交易开关
	VolumePrecision      int    `json:"volume_precision"`        //数量精度
	PricePrecision       int    `json:"price_precision"`         //价格精度
	AmountPrecision      int    `json:"amount_precision"`        //金额精度
	MinUnit              string `json:"min_unit"`                //最小交易单位
	MinChangePrice       string `json:"min_change_price"`        //最小变动价位
	Spread               string `json:"spread"`                  //点差，品种tick值的整数倍
	FeeKind              int    `json:"fee_kind"`                //手续费算法，取值1或2，1表示百比分，2表示绝对值。默认为1.
	TakerFee             string `json:"taker_fee"`               //Taker手续费率
	MakerFee             string `json:"maker_fee"`               //Maker手续费率
	SingleMinOrderVolume string `json:"single_min_order_volume"` //单次最小下单量
	SingleMaxOrderVolume string `json:"single_max_order_volume"` //单次最大下单量
	SingleMinOrderAmount string `json:"single_min_order_amount"` //单次最小下单金额
	SingleMaxOrderAmount string `json:"single_max_order_amount"` //单次最大下单金额
	BuyPriceLimit        string `json:"buy_price_limit"`         //买委托价格限制
	SellPriceLimit       string `json:"sell_price_limit"`        //卖委托价格限制
	MaxMatchGear         int    `json:"max_match_gear"`          //最大成交档位,不得超过20
	OtcMinOrderVolume    string `json:"otc_min_order_volume"`    //OTC最小量
	OtcMaxOrderVolume    string `json:"otc_max_order_volume"`    //OTC最大量
	OtcMinOrderAmount    string `json:"otc_min_order_amount"`    //OTC最小金额
	OtcMaxOrderAmount    string `json:"otc_max_order_amount"`    //OTC最大金额
	Operator             string `json:"operator"`                //操作员
	Time                 string `json:"time"`                    //操作时间
}

type TradeSwitch struct {
	TradeSwitch int `json:"trade_switch"`
}

type GroupTradeParam struct {
	GroupCode                         string `json:"group_code"`                             //客户组，如 WXBroker-VIP1
	SymbolId                          string `json:"symbol_id"`                              //币种代码，如 BTC_USDT
	Spread                            string `json:"spread"`                                 //点差，品种tick值的整数倍
	FeeKind                           int    `json:"fee_kind"`                               //手续费算法，取值1或2，1表示百比分，2表示绝对值。默认为1.
	TakerFee                          string `json:"taker_fee"`                              //Taker手续费率
	MakerFee                          string `json:"maker_fee"`                              //Maker手续费率
	Operator                          string `json:"operator"`                               //操作员
	Time                              string `json:"time"`                                   //操作时间
	DailyBuyingLimitPerClientWeekday  string `json:"daily_buying_limit_per_client_weekday"`  //每人每日买交易限额
	DailyBuyingLimitPerClientWeekend  string `json:"daily_buying_limit_per_client_weekend"`  //每人每日买交易限额
	DailySellingLimitPerClientWeekday string `json:"daily_selling_limit_per_client_weekday"` //每人每日卖交易限额
	DailySellingLimitPerClientWeekend string `json:"daily_selling_limit_per_client_weekend"` //每人每日卖交易限额
}

type GroupTradeAuth struct {
	GroupCode            string `json:"group_code"`             //客户组，如 WXBroker-VIP1
	SymbolId             string `json:"symbol_id"`              //币种代码，如 BTC_USDT
	BuyCommissionSwitch  bool   `json:"buy_commission_switch"`  //买委托限制开关,true表示禁止
	SellCommissionSwitch bool   `json:"sell_commission_switch"` //卖委托限制开关,true表示禁止
	User                 string `json:"user"`                   //操作员
	Time                 string `json:"time"`                   //操作时间
}

type BrokerConf struct {
	Mysql           Mysql
	DingDingTalk    DingDingTalk
	DingDingAskTalk DingDingAskTalk
	DingDingPms     DingDingPms
}

type Mysql struct {
	Addr string
}

type DingDingTalk struct {
	AccessToken string
	Secret      string
}

type DingDingAskTalk struct {
	AccessToken string
	Secret      string
}

type DingDingPms struct {
	AccessToken string
	Secret      string
}
