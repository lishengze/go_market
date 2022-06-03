package order

import (
	"bcts/app/client/api/internal/svc"
	"bcts/app/client/api/internal/types"
	"bcts/app/order/rpc/order"
	"bcts/common/crypto"
	"bcts/common/xerror"
	"context"
	"github.com/pkg/errors"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtcQuoteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewOtcQuoteLogic(r *http.Request, ctx context.Context, svcCtx *svc.ServiceContext) *OtcQuoteLogic {
	return &OtcQuoteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
	}
}

func (l *OtcQuoteLogic) OtcQuote(req *types.OtcQuoteReq) (resp *types.OtcQuoteRsp, err error) {

	userID, err := crypto.ExportUserIDFromHeader(l.r, l.svcCtx.PemFileBase.AESKey)
	if err != nil {
		return nil, errors.Wrapf(xerror.ErrorTryAgain, "GetUserInfo ExportUserIDFromHeader req:%+v, error:%+v", req, err)
	}
	req.UserID = userID

	logx.Infof("%+v", req)
	params := &order.QuoteReq{
		UserID:    req.UserID,
		Symbol:    req.Symbol,
		Direction: req.Direction,
		QuoteType: req.QuoteType,
		Volume:    req.Volume,
		Amount:    req.Amount,
	}
	result, err := l.svcCtx.OrderRpc.OtcQuote(l.ctx, params)
	if err != nil {
		logx.Error(err)
		return
	}

	resp = &types.OtcQuoteRsp{
		QuoteID: req.QuoteID,
		UserID:  req.UserID,
		Symbol:  req.Symbol,
		Price:   result.Price,
		Amount:  result.Amount,
		Volume:  result.Volume,
	}

	return
}
