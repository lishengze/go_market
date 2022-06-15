package logic

import (
	"context"

	"exterior-interactor/app/opu/rpc/internal/svc"
	"exterior-interactor/app/opu/rpc/opupb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSymbolLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSymbolLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSymbolLogic {
	return &GetSymbolLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSymbolLogic) GetSymbol(in *opupb.GetSymbolReq) (*opupb.GetSymbolRsp, error) {
	return l.svcCtx.OPU.GetSymbol(in)
}
