package logic

import (
	"context"
	"exterior-interactor/app/idsrv/rpc/idsrvpb"
	"exterior-interactor/pkg/xsnowflake"

	"exterior-interactor/app/idsrv/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetIdLogic {
	return &GetIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetIdLogic) GetId(in *idsrvpb.EmptyReq) (*idsrvpb.GetIdRsp, error) {
	return &idsrvpb.GetIdRsp{
		Id: xsnowflake.GenerateId().String(),
	}, nil
}
