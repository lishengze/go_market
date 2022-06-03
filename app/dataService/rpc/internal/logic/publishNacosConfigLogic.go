package logic

import (
	"context"

	"bcts/app/dataService/rpc/internal/svc"
	"bcts/app/dataService/rpc/types/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishNacosConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishNacosConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishNacosConfigLogic {
	return &PublishNacosConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 发布参数
func (l *PublishNacosConfigLogic) PublishNacosConfig(in *pb.PublishNacosConfig) (*pb.DSEmpty, error) {
	// todo: add your logic here and delete this line

	return &pb.DSEmpty{}, nil
}
