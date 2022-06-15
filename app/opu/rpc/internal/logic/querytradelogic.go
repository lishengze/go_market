package logic

import (
	"context"

	"exterior-interactor/app/opu/rpc/internal/svc"
	"exterior-interactor/app/opu/rpc/opupb"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryTradeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryTradeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryTradeLogic {
	return &QueryTradeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QueryTradeLogic) QueryTrade(in *opupb.QueryTradeReq) (*opupb.QueryTradeRsp, error) {
	return l.svcCtx.OPU.QueryTrade(in)
}
