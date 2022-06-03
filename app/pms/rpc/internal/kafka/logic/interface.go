package logic

import (
	"market_server/app/pms/rpc/internal/kafka/types"
)

type LogicHandle interface {
	Handle(msg *types.TradeData) error
}
