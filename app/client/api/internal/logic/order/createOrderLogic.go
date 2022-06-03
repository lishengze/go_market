package order

import (
	"context"
	"market_server/app/client/api/internal/svc"
	"market_server/app/client/api/internal/types"
	"market_server/app/order/rpc/order"
	"market_server/common/crypto"
	"market_server/common/xerror"
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewCreateOrderLogic(r *http.Request, ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
	}
}

// 用户报单
func (l *CreateOrderLogic) CreateOrder(req *types.CreateOrderReq) (resp *types.CreateOrderRsp, err error) {

	userID, err := crypto.ExportUserIDFromHeader(l.r, l.svcCtx.PemFileBase.AESKey)
	if err != nil {
		return nil, errors.Wrapf(xerror.ErrorTryAgain, "GetUserInfo ExportUserIDFromHeader req:%+v, error:%+v", req, err)
	}

	req.UserID = userID

	param := &order.CreateOrderReq{
		RequestID:      req.RequestID,
		UserID:         req.UserID,
		Symbol:         req.Symbol,
		OrderMode:      req.OrderMode,
		OrderType:      req.OrderType,
		OrderPriceType: req.OrderPriceType,
		Direction:      req.Direction,
		Volume:         req.Volume,
		Amount:         req.Amount,
		Price:          req.Price,
	}

	res, err := l.svcCtx.OrderRpc.CreateOrder(l.ctx, param)
	if err != nil {
		return nil, err
	}

	var rsp types.CreateOrderRsp
	copier.Copy(&rsp, res)

	return &rsp, nil
}
