package withdraw

import (
	"context"

	"bcts/app/hedge/cmd/api/internal/svc"
	"bcts/app/hedge/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWithdrawInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetWithdrawInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWithdrawInfoLogic {
	return &GetWithdrawInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWithdrawInfoLogic) GetWithdrawInfo(req *types.WithdrawReq) (resp *types.WithdrawReply, err error) {
	// todo: add your logic here and delete this line

	return
}
