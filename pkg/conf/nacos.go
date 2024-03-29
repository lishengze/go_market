package config

import (
	"market_aggregate/pkg/util"

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
		util.LOG_ERROR(err.Error())
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

/*
        "symbol_id":"BTC_USD",
        "switch":true,
        "publish_frequency":3000,
        "publish_level":20,
        "price_offset_kind":1,
        "price_offset":0.00001,
        "amount_offset_kind":1,
        "amount_offset":0.001,
        "poll_offset_kind":1,
        "poll_offset":0.001,
        "user":"admin",
        "time":"2022-03-15 08:47:11"
*/

type MarketRiskConfig struct {
	Symbol           string  `json:"symbol_id"`          //币种代码，如 USDT
	Switch           bool    `json:"switch"`             //
	PublishFrequency int     `json:"publish_frequency"`  //发布频率 毫秒
	PublishLevel     int     `json:"publish_level"`      //发布档位
	PriceOffsetKind  int     `json:"price_offset_kind"`  //价格偏移种类
	PriceOffset      float64 `json:"price_offset"`       //价格偏移值
	AmountOffsetKind int     `json:"amount_offset_kind"` //数量偏移种类
	AmountOffset     float64 `json:"amount_offset"`      //数量偏移值
	OTCOffsetKind    int     `json:"poll_offset_kind"`   //OTC 询价偏移种类
	OTCOffset        float64 `json:"poll_offset"`        //OTC 询价偏移值
	User             string  `json:"admin"`              //大额提现阈值
	Time             string  `json:"time"`               //提现手续费算法 1:百比分，2:绝对值。默认为1.
}

/*
HedgeParams
    "platform_id":"B2C2",
    "instrument":"BTC_USDT",
    "switch":true,
    "buy_fund_ratio":0.3,
    "sell_fund_ratio":0.3,
    "price_precision":2,
    "amount_precision":4,
    "sum_precision":2,

    "min_unit":0.0001,
    "min_change_price":0.01,
    "max_margin":1,
    "max_order":100,
    "buy_price_limit":10000,
    "sell_price_limit":5000,
    "max_match_level":5,
    "fee_kind":1,
    "taker_fee":0,
    "maker_fee":0
*/
type HedgingConfig struct {
	Exchange        string  `json:"platform_id"`     //交易平台，如 binance
	Symbol          string  `json:"instrument"`      //品种
	Switch          bool    `json:"switch"`          //充值开关
	BuyFundRatio    float64 `json:"buy_fund_ratio"`  //买盘资金占用比例
	SellFundRatio   float64 `json:"sell_fund_ratio"` //卖盘资金占用比例
	PricePrecision  int     `json:"price_precision"`
	AmountPrecision int     `json:"amount_precision"`
	SumPrecision    int     `json:"sum_precision"`
	MinUnit         float64 `json:"min_unit"`         //最小交易单位
	MinChangePrice  float64 `json:"min_change_price"` //最小变动价
	MaxMargin       float64 `json:"max_margin"`       //
	MaxOrder        float64 `json:"max_order"`        //
	BuyPriceLimit   int     `json:"buy_price_limit"`  //买委托价格限制
	SellPriceLimit  int     `json:"sell_price_limit"` //卖委托价格限制
	MaxMatchLevel   int     `json:"max_match_level"`
	FeeKind         int     `json:"fee_kind"`  //手续费算法，取值1或2，1表示百比分，2表示绝对值。默认为1.
	TakerFee        float64 `json:"taker_fee"` //Taker手续费率
	MakerFee        float64 `json:"maker_fee"` //Maker手续费率
}

/*
        "symbol_id":"BTC_USD",
        "symbol_kind":"现货",
        "bid":"BTC",
        "primary_currency":"BTC",
        "bid_currency":"USD",
        "settle_currency":"USD",
        "switch":true,
        "amount_precision":6,
        "price_precision":6,
        "sum_precision":4,
        "min_unit":0.000001,
        "min_change_price":0.01,
        "spread":0,
        "fee_kind":1,
        "taker_fee":0.1,
        "maker_fee":0.1,

        "min_order":0.0001,
        "max_order":100,
        "min_money":0.01,
        "max_money":10000,
        "buy_price_limit":20,
        "sell_price_limit":20,
        "max_match_level":10,
        "otc_min_order":0.001,
        "otc_max_order":2,
        "otc_min_price":1,
        "otc_max_price":100000,
        "user":"admin",
        "time":"2022-02-14 06:24:48"
*/

type SymbolConfig struct {
	Symbol          string `json:"symbol_id"`        //品种代码，如 BTC_USDT
	SymbolKind      string `json:"symbol_kind"`      //品种类型，如 1-现货、2-期货等
	PrimaryCurrency string `json:"primary_currency"` //基础货币，如 USDT
	BidCurrency     string `json:"bid_currency"`     //报价货币，如 USDT
	SettleCurrency  string `json:"settle_currency"`  //结算货币，如 USDT
	Switch          bool   `json:"switch"`           //交易开关

	AmountPrecision int `json:"amount_precision"` //数量精度
	PricePrecision  int `json:"price_precision"`  //价格精度
	SumPrecision    int `json:"sum_precision"`    //金额精度

	MinUnit              float64 `json:"min_unit"`         //最小交易单位
	MinChangePrice       float64 `json:"min_change_price"` //最小变动价位
	Spread               int     `json:"spread"`           //点差，品种tick值的整数倍
	FeeKind              int     `json:"fee_kind"`         //手续费算法，取值1或2，1表示百比分，2表示绝对值。默认为1.
	TakerFee             float64 `json:"taker_fee"`        //Taker手续费率
	MakerFee             float64 `json:"maker_fee"`        //Maker手续费率
	SingleMinOrderVolume float64 `json:"min_order"`        //单次最小下单量
	SingleMaxOrderVolume float64 `json:"max_order"`        //单次最大下单量
	SingleMinOrderAmount float64 `json:"min_money"`        //单次最小下单金额
	SingleMaxOrderAmount float64 `json:"max_money"`        //单次最大下单金额
	BuyPriceLimit        int     `json:"buy_price_limit"`  //买委托价格限制
	SellPriceLimit       int     `json:"sell_price_limit"` //卖委托价格限制
	MaxMatchLevel        int     `json:"max_match_level"`  //最大成交档位,不得超过20
	OtcMinOrderVolume    float64 `json:"otc_min_order"`    //OTC最小量
	OtcMaxOrderVolume    float64 `json:"otc_max_order"`    //OTC最大量
	OtcMinOrderAmount    float64 `json:"otc_min_price"`    //OTC最小金额
	OtcMaxOrderAmount    float64 `json:"otc_max_price"`    //OTC最大金额
	Operator             string  `json:"user"`             //操作员
	Time                 string  `json:"time"`             //操作时间
}
