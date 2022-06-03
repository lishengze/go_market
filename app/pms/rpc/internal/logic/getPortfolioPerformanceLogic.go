package logic

import (
	"context"

	"bcts/app/pms/rpc/internal/svc"
	"bcts/app/pms/rpc/types/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPortfolioPerformanceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPortfolioPerformanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPortfolioPerformanceLogic {
	return &GetPortfolioPerformanceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPortfolioPerformanceLogic) GetPortfolioPerformance(in *pb.PortfolioPerformanceReq) (*pb.PortfolioPerformanceRsp, error) {
	// todo: add your logic here and delete this line

	return &pb.PortfolioPerformanceRsp{}, nil
}
