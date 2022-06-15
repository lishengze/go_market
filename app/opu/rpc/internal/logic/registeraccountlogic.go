package logic

import (
	"context"

	"exterior-interactor/app/opu/rpc/internal/svc"
	"exterior-interactor/app/opu/rpc/opupb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterAccountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterAccountLogic {
	return &RegisterAccountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterAccountLogic) RegisterAccount(in *opupb.RegisterAccountReq) (*opupb.RegisterAccountRsp, error) {
	return l.svcCtx.OPU.RegisterAccount(in)
}
