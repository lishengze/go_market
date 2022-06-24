package logic

import (
	"context"

	"exterior-interactor/app/opu/rpc/internal/svc"
	"exterior-interactor/app/opu/rpc/opupb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelAllOpenOrdersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCancelAllOpenOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelAllOpenOrdersLogic {
	return &CancelAllOpenOrdersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CancelAllOpenOrdersLogic) CancelAllOpenOrders(in *opupb.CancelAllOpenOrdersReq) (*opupb.EmptyRsp, error) {
	return l.svcCtx.OPU.CancelAllOpenOrders(in)
}
