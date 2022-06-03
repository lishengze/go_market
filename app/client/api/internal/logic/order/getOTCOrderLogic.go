package order

import (
	"context"
	"market_server/app/order/rpc/order"
	"net/http"
	"strconv"

	"market_server/app/client/api/internal/svc"
	"market_server/app/client/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOTCOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewGetOTCOrderLogic(r *http.Request, ctx context.Context, svcCtx *svc.ServiceContext) *GetOTCOrderLogic {
	return &GetOTCOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
	}
}

func (l *GetOTCOrderLogic) GetOTCOrder(req *types.OTCOrderQueryReq) (resp *types.OrderQueryRsp, err error) {
	memberId := int64(32292864) //TODO: 从当前登录用户session中获取user_id
	logx.Infof("memberId:%d", memberId)
	logx.Infof("%+v", req)

	params := &order.OTCOrderQueryReq{
		UserID:   memberId,
		Page:     int64(req.Page),
		PageSize: int64(req.PageSize),
	}
	result, err := l.svcCtx.OrderRpc.GetOTCOrder(l.ctx, params)
	if err != nil {
		logx.Error(err)
		return
	}
	var orders []*types.OrderOrderItem
	for _, v := range result.Orders {
		otcOrder := &types.OrderOrderItem{
			OrderLocalId:    v.OrderLocalId,
			UserID:          strconv.FormatInt(v.UserId, 10),
			OrderType:       int(v.OrderType),
			OrderPriceType:  int(v.OrderPriceType),
			Symbol:          v.Symbol,
			Direction:       int(v.Direction),
			Price:           v.Price,
			Volume:          v.Volume,
			Amount:          v.Amount,
			OrderStatus:     int(v.OrderStatus),
			OrderMaker:      int(v.OrderMaker),
			OrderCreateTime: v.OrderCreateTime,
			OrderModifyTime: v.OrderModifyTime,
		}
		orders = append(orders, otcOrder)
	}
	resp = &types.OrderQueryRsp{
		Count:  strconv.FormatInt(result.Count, 10),
		Orders: orders,
	}

	return
}
