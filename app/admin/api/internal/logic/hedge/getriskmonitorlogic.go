package hedge

import (
	"bcts/app/admin/api/internal/svc"
	"bcts/app/admin/api/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRiskMonitorLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRiskMonitorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRiskMonitorLogic {
	return &GetRiskMonitorLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRiskMonitorLogic) GetRiskMonitor(req *types.ReqQueryHedgeRisk) (resp *types.QueryHedgeRiskReply, err error) {
	// todo: add your logic here and delete this line

	return
}
