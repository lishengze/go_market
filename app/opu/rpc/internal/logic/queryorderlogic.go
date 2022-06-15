package logic

import (
	"context"

	"exterior-interactor/app/opu/rpc/internal/svc"
	"exterior-interactor/app/opu/rpc/opupb"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryOrderLogic {
	return &QueryOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QueryOrderLogic) QueryOrder(in *opupb.QueryOrderReq) (*opupb.QueryOrderRsp, error) {
	return l.svcCtx.OPU.QueryOrder(in)
}
