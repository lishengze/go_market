package exmodel

import (
	"strings"
	"time"
)

const (
	OrderSideBuy     OrderSide = "BUY"
	OrderSideSell    OrderSide = "SELL"
	OrderSideUnknown OrderSide = "UNKNOWN"
)

const (
	OrderTypeLimit   OrderType = "LIMIT"
	OrderTypeMarket  OrderType = "MARKET"
	OrderTypeUnknown OrderType = "UNKNOWN"
)

const (
	OrderStatusUnknown    OrderStatus = "UNKNOWN"
	OrderStatusPending    OrderStatus = "PENDING"
	OrderStatusRejected   OrderStatus = "REJECTED" // 终态
	OrderStatusSent       OrderStatus = "SENT"
	OrderStatusPartial    OrderStatus = "PARTIAL"
	OrderStatusCancelling OrderStatus = "CANCELLING"
	OrderStatusCancelled  OrderStatus = "CANCELLED" // 终态
	OrderStatusFilled     OrderStatus = "FILLED"    // 终态 完全成交
	//OrderStatusForcedLiquidate OrderStatus = "FORCED_LIQUIDATE" // 终态 强制平仓
)

const (
	OrderUpdate  OrderTradesUpdateType = "ORDER_UPDATE"
	TradesUpdate OrderTradesUpdateType = "TRADES_UPDATE"
)

const (
	LiquidityTaker   Liquidity = "TAKER"
	LiquidityMaker   Liquidity = "MAKER"
	LiquidityUnknown Liquidity = "UNKNOWN"
)

type (
	OrderSide             string
	OrderType             string
	TimeInForce           string
	OrderStatus           string
	OrderTradesUpdateType string
	Liquidity             string

	Order struct {
		OrderId        string // 交易所的 OrderId
		ClientOrderId  string
		Volume         string
		Price          string
		FilledVolume   string // 已成交总量
		Exchange       Exchange
		SymbolExFormat string
		ApiType
		Side   OrderSide
		Type   OrderType
		Status OrderStatus
	}

	Trade struct {
		TradeId string // 交易所的 TradeId
		OrderId string // 交易所的 OrderId
		//ClientOrderId string
		Exchange       Exchange
		SymbolExFormat string
		ApiType
		Liquidity   Liquidity
		Volume      string
		Price       string
		Fee         string
		FeeCurrency Currency
		TradeTime   time.Time
	}

	PlaceOrderReq struct {
		SymbolExFormat string
		ApiType
		ClientOrderId string
		OrderType     OrderType
		Side          OrderSide
		Volume        string
		Price         string
	}

	PlaceOrderRsp struct {
		OrderId       string // 交易所的 OrderId
		ClientOrderId string
		FilledVolume  string
		Status        OrderStatus
	}

	CancelOrderReq struct {
		OrderId        string // 交易所的 OrderId
		ClientOrderId  string
		SymbolExFormat string
		ApiType
	}

	CancelOrderRsp struct {
	}

	QueryOrderReq struct {
		OrderId        string // 交易所的 OrderId
		ClientOrderId  string
		SymbolExFormat string
		ApiType
	}

	QueryTradeReq struct {
		OrderId        string // 交易所的 OrderId
		SymbolExFormat string
		ApiType
		StartTime time.Time
		EndTime   time.Time
	}

	// OrderTradesUpdate 订单更新或成交更新 推送
	OrderTradesUpdate struct {
		Type          OrderTradesUpdateType
		OrderId       string // 交易所的 OrderId
		ClientOrderId string
		*OrderUpdateInfo
		*TradesUpdateInfo
	}

	OrderUpdateInfo struct {
		OrderStatus  OrderStatus
		FilledVolume string    // 已成交总量
		UpdateTime   time.Time // 订单更新时间
	}

	TradesUpdateInfo struct {
		TradeId     string // 交易所的 TradeId
		Price       string
		Volume      string
		Fee         string
		FeeCurrency Currency
		Liquidity   Liquidity
		TradeTime   time.Time // 成交时间
	}
)

func (o TimeInForce) String() string {
	return string(o)
}

func (o OrderStatus) String() string {
	return string(o)
}

func (o OrderTradesUpdateType) String() string {
	return string(o)
}

// IsClosed 判断订单是否终结
func (o OrderStatus) IsClosed() bool {
	switch o {
	case OrderStatusRejected, OrderStatusFilled, OrderStatusCancelled:
		return true
	default:
		return false
	}
}

func (o OrderType) String() string {
	return string(o)
}

func (o OrderType) Lower() string {
	return strings.ToLower(string(o))
}

func (o OrderType) Upper() string {
	return strings.ToUpper(string(o))
}

func (o OrderSide) String() string {
	return string(o)
}

func (o OrderSide) Upper() string {
	return strings.ToUpper(string(o))
}

func (o OrderSide) Lower() string {
	return strings.ToLower(string(o))
}

func (o Liquidity) String() string {
	return string(o)
}

func (o Liquidity) Upper() string {
	return strings.ToUpper(string(o))
}

func (o Liquidity) Lower() string {
	return strings.ToLower(string(o))
}
