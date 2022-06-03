package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"market_server/app/pms/rpc/internal/kafka/types"
	"market_server/app/pms/rpc/internal/svc"
	"math"

	"github.com/zeromicro/go-zero/core/logx"
)

type TradeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTradeLogic(ctx context.Context, svcCtx *svc.ServiceContext) LogicHandle {
	return LogicHandle(&TradeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	})
}

type MarketDetail struct {
	MarketPrice float64
	Timestamp   uint64
}

func (l *TradeLogic) Handle(msg *types.TradeData) error {
	//数据存入redis
	marketPrice := float64(msg.Price.Value) / math.Pow10(int(msg.Price.Precise))
	redisKey := "marketPrice:" + msg.Symbol
	fmt.Println(msg.Symbol, msg.Time, msg.Price, marketPrice)
	value := MarketDetail{
		marketPrice,
		msg.Time,
	}
	valueStr, _ := json.Marshal(value)
	err := l.svcCtx.RedisConn.Set(redisKey, fmt.Sprintf("%s", valueStr))
	if err != nil {
		return err
	}
	return nil
}
