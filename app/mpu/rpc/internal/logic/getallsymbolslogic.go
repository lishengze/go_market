package logic

import (
	"context"

	"exterior-interactor/app/mpu/rpc/internal/svc"
	"exterior-interactor/app/mpu/rpc/mpupb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllSymbolsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllSymbolsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllSymbolsLogic {
	return &GetAllSymbolsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllSymbolsLogic) GetAllSymbols(in *mpupb.EmptyReq) (*mpupb.GetAllSymbolsRsp, error) {
	// todo: add your logic here and delete this line

	return &mpupb.GetAllSymbolsRsp{}, nil
}
