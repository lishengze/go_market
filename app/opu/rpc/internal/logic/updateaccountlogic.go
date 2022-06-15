package logic

import (
	"context"

	"exterior-interactor/app/opu/rpc/internal/svc"
	"exterior-interactor/app/opu/rpc/opupb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAccountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAccountLogic {
	return &UpdateAccountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateAccountLogic) UpdateAccount(in *opupb.UpdateAccountReq) (*opupb.EmptyRsp, error) {
	return l.svcCtx.OPU.UpdateAccount(in)
}
