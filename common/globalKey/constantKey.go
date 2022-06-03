package globalKey

/**
global constant key
*/

//软删除
var DelStateNo int64 = 0  //未删除
var DelStateYes int64 = 1 //已删除

//时间格式化模版
var DateTimeFormatTplStandardDateTime = "Y-m-d H:i:s"
var DateTimeFormatTplStandardDate = "Y-m-d"
var DateTimeFormatTplStandardTime = "H:i:s"

//Direction
const (
	DirectionBuy  = 0
	DirectionSell = 1
)

//手续费
const (
	FeeKindPercentage = 1 //百分比
	FeeKindAbsolute   = 2 //绝对值
)

//询价方式
const (
	QuoteTypeVolume = 1
	QuoteTypeAmount = 2
)

//OrderType
const (
	OrderTypeOTC   = 1 //OTC交易
	OrderTypeAggre = 2 //聚合交易
)

//OrderMode
const (
	OrderModeVolume = 1
	OrderModeAmount = 2
)

//OrderPriceType
const (
	OrderPriceTypeLimit  = 1
	OrderPriceTypeMarket = 2
)

//MarketMaker
const (
	MarketTaker = 1
	MarketMaker = 2
)

const (
	OrderStatusSent            = 0 //已发送
	OrderStatusFilled          = 1 //全部成交
	OrderStatusPartialInQueue  = 2 //部分成交还在队列
	OrderStatusPartialCanceled = 3 //部分成交已经撤单
	OrderStatusCanceled        = 4 //撤单
)
