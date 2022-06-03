package logic

import (
	"context"
	"market_server/app/admin/rpc/internal/svc"
	"market_server/app/admin/rpc/types/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWebsocketUrlLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetWebsocketUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWebsocketUrlLogic {
	return &GetWebsocketUrlLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetWebsocketUrlLogic) GetWebsocketUrl(in *pb.GetWebsocketUrlReq) (*pb.GetWebsocketUrlRsp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetWebsocketUrlRsp{}, nil
}
