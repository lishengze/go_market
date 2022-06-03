package user

import (
	"bcts/app/userCenter/rpc/usercenter"
	"bcts/common/crypto"
	"bcts/common/xerror"
	"context"
	"github.com/pkg/errors"
	"net/http"

	"bcts/app/client/api/internal/svc"
	"bcts/app/client/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAccountInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewGetAccountInfoLogic(r *http.Request, ctx context.Context, svcCtx *svc.ServiceContext) *GetAccountInfoLogic {
	return &GetAccountInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
	}
}

func (l *GetAccountInfoLogic) GetAccountInfo(req *types.AccountInfoReq) (resp *types.AccountInfoRsp, err error) {

	userID, err := crypto.ExportUserIDFromHeader(l.r, l.svcCtx.PemFileBase.AESKey)
	if err != nil {
		return nil, errors.Wrapf(xerror.ErrorTryAgain, "GetUserInfo ExportUserIDFromHeader req:%+v, error:%+v", req, err)
	}

	req.UserID = userID

	res, err := l.svcCtx.UserCenterRpc.GetAccountInfo(l.ctx, &usercenter.AccountInfoReq{
		UserID:   userID,
		Currency: req.Currency,
	})

	if err != nil {
		return nil, err
	}

	var data = make([]*types.AccountInfo, 0)
	for _, v := range res.AccountInfos {
		data = append(data, &types.AccountInfo{
			UserID:    v.UserID,
			Currency:  v.Currency,
			Available: v.Available,
			Frozen:    v.Frozen,
			Balance:   v.Balance,
			CostPrice: v.CostPrice,
		})
	}

	return &types.AccountInfoRsp{
		Count: len(data),
		List:  data,
	}, nil
}
