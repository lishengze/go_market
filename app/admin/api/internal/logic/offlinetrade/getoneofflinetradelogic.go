package offlinetrade

import (
	"context"
	"market_server/app/admin/api/internal/svc"
	"market_server/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOneOfflineTradeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOneOfflineTradeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOneOfflineTradeLogic {
	return &GetOneOfflineTradeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOneOfflineTradeLogic) GetOneOfflineTrade(req *types.OfflineTradeReq) (resp *types.OfflineTradeInput, err error) {
	record, err := l.svcCtx.OfflineTradeInputModel.FindOne(l.ctx, req.ID)
	if err != nil {
		return nil, err
	}

	resp = ConvertModel(record)
	return
}
