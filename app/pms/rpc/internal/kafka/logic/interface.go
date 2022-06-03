package logic

import (
	"bcts/app/pms/rpc/internal/kafka/types"
)

type LogicHandle interface {
	Handle(msg *types.TradeData) error
}
