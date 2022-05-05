package logic

import (
	"context"

	"exterior-interactor/app/mpu/rpc/internal/svc"
	"exterior-interactor/app/mpu/rpc/mpupb"

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

func (l *GetSymbolLogic) GetSymbol(in *mpupb.GetSymbolReq) (*mpupb.GetSymbolRsp, error) {
	// todo: add your logic here and delete this line

	return &mpupb.GetSymbolRsp{}, nil
}
