package exmodel

type ApiType string // 表示 交易所 的 接口类型

const (
	ApiTypeUnified    = "UNIFIED" // 统一接口，不区分类型
	ApiTypeSpot       = "SPOT"
	ApiTypeCoinFuture = "COIN_FUTURE"
	ApiTypeUsdtFuture = "USDT_FUTURE"
)
