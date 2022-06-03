package model

import "strings"

type OfflineTradeInputStatus int

const (
	OfflineTradeInputStatusFirst OfflineTradeInputStatus = iota + 1
	OfflineTradeInputStatusSecond
)

func (o OfflineTradeInputStatus) INT() int {
	switch o {
	case OfflineTradeInputStatusFirst:
		return 1
	case OfflineTradeInputStatusSecond:
		return 2
	}
	return 0
}

type (
	TradeType   int
	TradeStatus int
	SourceType  int
)

const (
	TradeTypeInner   TradeType = iota // 0 内部成交
	TradeTypeOuter                    // 1 外部成交
	TradeTypeOffline                  //2 线下
)

const (
	TradeStatusNone    TradeStatus = iota // 未对冲 "none"
	TradeStatusHedging                    // 对冲处理中 "hedging"
	TradeStatusDone                       // 已完成 "done"
)

func (o TradeType) INT() int {
	switch o {
	case TradeTypeInner:
		return 0
	case TradeTypeOuter:
		return 1
	case TradeTypeOffline:
		return 2
	}
	return 0
}

func IsTradeStatusValid(status string) bool {
	if status != "" {
		statusAfter := strings.Split(status, ",")
		for _, v := range statusAfter {
			if v == "0" || v == "1" || v == "2" {
				return true
			} else {
				return false
			}
		}
	}
	return false
}

// MemberTransaction

type TransactionType int

const (
	TransTypeRecharge       TransactionType = iota // 0 充值 "RECHARGE"
	TransTypeWithdraw                              // 1 提现 "WITHDRAW"
	TransTypeWithdrawFee                           // 2 提现手续费 "WITHDRAW Fee"
	TransTypeTransfer                              // 3 转账 "TRANSFER"
	TransTypeExchange                              // 4 币币交易 "EXCHANGE"
	TransTypeFee                                   // 5 交易手续费 "Fee"
	TransTypeExchangeTrade                         // 6 币币交易成交 "EXCHANGE"
	TransTypeExchangeCancel                        // 7 币币交易取消 "EXCHANGE"
	TransTypeBuy                                   // 8 买入 "BUY"
	TransTypeSell                                  // 9 卖出 "SELL"
	TransTypeDividend                              // 10 分红 "DIVIDEND"
	TransTypeTradeFix                              // 11 交易修正
	TransTypeFeeFix                                //12 交易手续费修正
)
