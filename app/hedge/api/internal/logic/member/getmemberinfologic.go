package member

import (
	"context"

	"market_server/app/hedge/api/internal/svc"
	"market_server/app/hedge/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMemberInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMemberInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMemberInfoLogic {
	return &GetMemberInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMemberInfoLogic) GetMemberInfo(req *types.MemberReq) (resp *types.MemberReply, err error) {
	// todo: add your logic here and delete this line

	return
}
