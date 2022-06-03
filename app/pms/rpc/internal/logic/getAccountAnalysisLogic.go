package logic

import (
	"context"

	"market_server/app/pms/rpc/internal/svc"
	"market_server/app/pms/rpc/types/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAccountAnalysisLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAccountAnalysisLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAccountAnalysisLogic {
	return &GetAccountAnalysisLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  资产列表
func (l *GetAccountAnalysisLogic) GetAccountAnalysis(in *pb.AccountAnalysisReq) (*pb.AccountAnalysisRsp, error) {
	// todo: add your logic here and delete this line

	return &pb.AccountAnalysisRsp{}, nil
}
