package user

import (
	"context"
	"market_server/app/client/api/internal/svc"
	"market_server/app/client/api/internal/types"
	"market_server/app/userCenter/rpc/usercenter"
	"market_server/common/crypto"
	"market_server/common/xerror"
	"net/http"

	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewGetUserInfoLogic(r *http.Request, ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
	}
}

func (l *GetUserInfoLogic) GetUserInfo(req *types.UserInfoReq) (resp *types.UserInfoRsp, err error) {

	userID, err := crypto.ExportUserIDFromHeader(l.r, l.svcCtx.PemFileBase.AESKey)
	if err != nil {
		return nil, errors.Wrapf(xerror.ErrorTryAgain, "GetUserInfo ExportUserIDFromHeader req:%+v, error:%+v", req, err)
	}

	req.UserID = userID

	res, err := l.svcCtx.UserCenterRpc.GetUserInfo(l.ctx, &usercenter.UserInfoReq{
		UserID: req.UserID,
	})

	if err != nil {
		return nil, err
	}

	return &types.UserInfoRsp{
		UserID:   res.UserID,
		UserType: res.UserType,
		Name:     res.Name,
		Email:    res.Email,
		Mobile:   res.Mobile,
		KycLevel: res.KycLevel,
	}, nil
}
