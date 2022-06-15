package logic

import (
	"context"

	"exterior-interactor/app/opu/rpc/internal/svc"
	"exterior-interactor/app/opu/rpc/opupb"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryBalanceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryBalanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryBalanceLogic {
	return &QueryBalanceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QueryBalanceLogic) QueryBalance(in *opupb.QueryBalanceReq) (*opupb.QueryBalanceRsp, error) {
	return l.svcCtx.OPU.QueryBalance(in)
}
