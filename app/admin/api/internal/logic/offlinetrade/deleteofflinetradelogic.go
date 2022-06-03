package offlinetrade

import (
	"bcts/app/admin/api/internal/svc"
	"bcts/app/admin/api/internal/types"
	"bcts/common/middleware"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteOfflineTradeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteOfflineTradeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteOfflineTradeLogic {
	return &DeleteOfflineTradeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteOfflineTradeLogic) BeforeDelete(req *types.DeleteOfflineTradeReq) (err error) {
	req.Operator = l.ctx.Value(middleware.JWT_USER_NAME).(string)

	base := NewBaseLogic(l.ctx, l.svcCtx)
	validators := []validator{
		base.validateFirstStatus,
		base.validateOperator,
	}

	for _, validator := range validators {
		err = validator(&types.OfflineTradeReq{
			ID:            req.ID,
			Operator:      req.Operator,
			TwoFactorCode: req.TwoFactorCode,
		})
		if err != nil {
			return err
		}
	}
	l.Logger.Info("BeforeCreate")
	return
}

func (l *DeleteOfflineTradeLogic) DeleteOfflineTrade(req *types.DeleteOfflineTradeReq) (resp *types.OfflineTradeInput, err error) {
	if err = l.BeforeDelete(req); err != nil {
		return
	}

	err = l.svcCtx.OfflineTradeInputModel.Delete(l.ctx, nil, req.ID)
	return
}
