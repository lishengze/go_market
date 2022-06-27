package logic

import (
	"context"
	"market_server/app/data_manager/rpc/internal/svc"
	"market_server/app/data_manager/rpc/types/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RequestTradeDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRequestTradeDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RequestTradeDataLogic {
	return &RequestTradeDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RequestTradeDataLogic) RequestTradeData(in *pb.ReqTradeInfo) (*pb.Trade, error) {
	// todo: add your logic here and delete this line

	return &pb.Trade{}, nil
}
