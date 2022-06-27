package logic

import (
	"context"
	"market_server/app/data_manager/rpc/internal/svc"
	"market_server/app/data_manager/rpc/types/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RequestHistKlineDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRequestHistKlineDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RequestHistKlineDataLogic {
	return &RequestHistKlineDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  服务名字还是用 MarketService, 整个行情系统都要用；
func (l *RequestHistKlineDataLogic) RequestHistKlineData(in *pb.ReqHishKlineInfo) (*pb.HistKlineData, error) {
	// todo: add your logic here and delete this line

	return &pb.HistKlineData{}, nil
}
