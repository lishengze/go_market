package config

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type ServerConfig struct {
	IpAddr string
	Port   uint64
}

type ClientConfig struct {
	NamespaceId         string
	TimeoutMs           uint64
	NotLoadCacheAtStart bool
	LogDir              string
	CacheDir            string
	LogLevel            string
}

type NacosConfig struct {
	*ServerConfig
	*ClientConfig
}

// type NacosConfig struct {
// 	IpAddr              string
// 	Port                int32
// 	NamespaceId         string
// 	TimeoutMs           int
// 	NotLoadCacheAtStart bool
// 	LogDir              string
// 	CacheDir            string
// 	RotateTime          string
// 	MaxAge              int32
// 	LogLevel            string
// }

type NacosClient struct {
	iClient config_client.IConfigClient
}

func NewNacosClient(c *NacosConfig) *NacosClient {
	sc := []constant.ServerConfig{
		{
			IpAddr: c.IpAddr,
			Port:   c.Port,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         c.NamespaceId,
		TimeoutMs:           c.TimeoutMs,
		NotLoadCacheAtStart: c.NotLoadCacheAtStart,
		LogDir:              c.LogDir,
		CacheDir:            c.CacheDir,
		LogLevel:            c.LogLevel,
	}

	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		panic(err)
	}
	return &NacosClient{
		iClient: client,
	}
}

func (c *NacosClient) GetConfigContent(dataId string, group string) (string, error) {
	content, err := c.iClient.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	if err != nil {
		return "", err
	}
	return content, err
}

func (c *NacosClient) PublishConfig(dataId string, group string, content string) error {
	_, err := c.iClient.PublishConfig(vo.ConfigParam{
		DataId:  dataId,
		Group:   group,
		Content: content,
	})
	return err
}

func (c *NacosClient) DeleteConfig(dataId string, group string) error {
	_, err := c.iClient.DeleteConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	return err
}

func (c *NacosClient) ListenConfig(dataId string, group string, f func(namespace, group, dataId, data string)) error {
	err := c.iClient.ListenConfig(vo.ConfigParam{
		DataId:   dataId,
		Group:    group,
		OnChange: f,
	})
	return err
}

type CurrencyConfig struct {
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

type HedgingConfig struct {
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

/*
 */
// std::string SymbolId;           // `json:"symbol_id"` //品种代码，如 BTC_USDT
// std::string SymbolKind;         //`json:"symbol_kind"` //品种类型，如 现货、期货等
// std::string Bid;                // `json:"bid"` //标的，如 BTC
// std::string PrimaryCurrency;    //  `json:"primary_currency"` //基础货币，如 USDT
// std::string BidCurrency;       // `json:"bid_currency"` //报价货币，如 USDT
// std::string SettleCurrency;     //   `json:"settle_currency"` //结算货币，如 USDT
// bool Switch;                  // `json:"switch"` //交易开关
// int AmountPrecision;         //`json:"amount_precision"` //数量精度
// int PricePrecision;       //json:"price_precision"` //价格精度
// int SumPrecision;         //`json:"sum_precision"` //金额精度
// double MinUnit;            //`json:"min_unit"` //最小交易单位
// double MinChangePrice;    //`json:"min_change_price"` //最小变动价位
// double Spread;            //`json:"spread"` //点差，品种tick值的整数倍
// int    FeeKind;        //`json:"fee_kind"` //手续费算法，取值1或2，1表示百比分，2表示绝对值。默认为1.
// double TakerFee;      //`json:"taker_fee"` //Taker手续费率
// double MakerFee;      //`json:"maker_fee"` //Maker手续费率
// double MinOrder;      //`json:"min_order"` //单次最小下单量
// double MaxOrder;    //`json:"max_order"` //单次最大下单量
// double MinMoney;    //      `json:"min_money"` //单次最小下单金额
// double MaxMoney;     // float64    `json:"max_money"` //单次最大下单金额
// double BuyPriceLimit;    //float64    `json:"buy_price_limit"` //买委托价格限制
// double SellPriceLimit;  //float64    `json:"sell_price_limit"` //卖委托价格限制
// int MaxMatchLevel;         //`json:"max_match_level"` //最大成交档位,不得超过20
// double OtcMinOrder; //`json:"otc_min_order"` //OTC最小量
// double OtcMaxOrder;    //`json:"otc_max_order"` //OTC最大量
// double OtcMinPrice;     //`json:"otc_min_price"` //OTC最小金额
// double OtcMaxPrice;     //`json:"otc_max_price"` //OTC最大金额
// std::string User;       //`json:"user"` //操作员
// std::string Time;          //`json:"time"` //操作时间

type SymbolConfig struct {
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
