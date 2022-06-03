package logic

import (
	"context"

	"bcts/app/dataService/rpc/internal/svc"
	"bcts/app/dataService/rpc/types/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteNacosConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteNacosConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteNacosConfigLogic {
	return &DeleteNacosConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除参数
func (l *DeleteNacosConfigLogic) DeleteNacosConfig(in *pb.DeleteNacosConfig) (*pb.DSEmpty, error) {
	// todo: add your logic here and delete this line

	return &pb.DSEmpty{}, nil
}
