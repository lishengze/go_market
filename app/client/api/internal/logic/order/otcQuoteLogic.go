package order

import (
	"context"
	"market_server/app/client/api/internal/svc"
	"market_server/app/client/api/internal/types"
	"market_server/app/order/rpc/order"
	"market_server/common/crypto"
	"market_server/common/xerror"
	"net/http"

	"github.com/pkg/errors"

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
