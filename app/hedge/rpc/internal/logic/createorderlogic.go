package logic

import (
	"context"

	"bcts/app/hedge/cmd/rpc/internal/svc"
	"bcts/app/hedge/cmd/rpc/types/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateOrderLogic) CreateOrder(in *pb.OrderReq) (*pb.OrderRsp, error) {
	// todo: add your logic here and delete this line

	return &pb.OrderRsp{}, nil
}
