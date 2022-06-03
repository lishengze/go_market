package user

import (
	"bcts/app/client/api/internal/svc"
	"bcts/app/client/api/internal/types"
	"bcts/common/crypto"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
)

type AddTestLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewAddTestLogic(r *http.Request, ctx context.Context, svcCtx *svc.ServiceContext) *AddTestLogic {
	return &AddTestLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
	}
}

func (l *AddTestLogic) AddTest(req *types.TestAddReq) (resp *types.TestAddRsp, err error) {

	userID, err := crypto.ExportUserIDFromHeader(l.r, l.svcCtx.PemFileBase.AESKey)

	l.Logger.Infof("req: %+v, userid:%s", req, userID)

	return &types.TestAddRsp{
		ID: "1234567",
	}, nil
}
