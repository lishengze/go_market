package user

import (
	"bcts/app/client/api/internal/svc"
	"bcts/app/client/api/internal/types"
	"bcts/app/userCenter/rpc/usercenter"
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserSettlementInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewGetUserSettlementInfoLogic(r *http.Request, ctx context.Context, svcCtx *svc.ServiceContext) *GetUserSettlementInfoLogic {
	return &GetUserSettlementInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
	}
}

func (l *GetUserSettlementInfoLogic) GetUserSettlementInfo(req *types.UserSettlementInfoReq) (resp *types.UserSettlementInfoRsp, err error) {
	_, _ = l.svcCtx.UserCenterRpc.GetFiatSettlementInfo(l.ctx, &usercenter.FiatSettlementInfoReq{
		UserID: req.UserID,
	})

	if err != nil {
		return nil, err
	}
	return nil, nil

	//return &types.UserSettlementInfoRsp{
	//	UserID:             res.UserID,
	//	BankName:           res.BankName,
	//	Swift:              res.Swift,
	//	RouteCode:          res.RouteCode,
	//	BeneficiaryName:    res.BeneficiaryName,
	//	BeneficiaryAccount: res.BeneficiaryAccount,
	//	BeneficiaryAddress: res.BeneficiaryAddress,
	//}, nil
}
