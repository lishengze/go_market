package logic

import (
	"context"
	"market_server/app/order/model"
	"market_server/common/xerror"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"

	"market_server/app/order/rpc/internal/svc"
	"market_server/app/order/rpc/types/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOTCOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOTCOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOTCOrderLogic {
	return &GetOTCOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetOTCOrderLogic) GetOTCOrder(in *pb.OTCOrderQueryReq) (*pb.OrderQueryRsp, error) {
	whereBuilder := l.svcCtx.OrderModel.RowBuilder()
	countBuilder := l.svcCtx.OrderModel.CountBuilder("id")
	if in.UserID > 0 {
		whereBuilder = whereBuilder.Where(squirrel.Eq{"user_id": in.UserID})
		countBuilder = countBuilder.Where(squirrel.Eq{"user_id": in.UserID})
	}
	count, err := l.svcCtx.OrderModel.FindCount(l.ctx, countBuilder)
	if err != nil {
		return nil, errors.Wrapf(xerror.ErrorDB, "Order ExchangeOrders err: %+v , in :%+v", err, in)
	}
	whereBuilder = whereBuilder.Offset(uint64(in.PageSize * (in.Page - 1))).Limit(uint64(in.PageSize))
	res, err := l.svcCtx.OrderModel.FindAll(l.ctx, whereBuilder, "")
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerror.ErrorDB, "Order ExchangeOrders err: %+v , in :%+v", err, in)
	}
	var orders []*pb.ExchangeOrder
	for _, v := range res {
		order := &pb.ExchangeOrder{
			Id:              v.Id,
			OrderLocalId:    v.OrderLocalId,
			UserId:          v.UserId,
			OrderType:       v.OrderType,
			OrderMode:       v.OrderMode,
			OrderPriceType:  v.OrderPriceType,
			Symbol:          v.Symbol,
			Direction:       v.Direction,
			OrderMaker:      v.OrderMaker,
			Volume:          v.Volume.String(),
			Amount:          v.Amount.String(),
			Price:           v.Price.String(),
			OrderStatus:     v.OrderStatus,
			TradeVolume:     v.TradeVolume.String(),
			TradeAmount:     v.TradeAmount.String(),
			FeeKind:         v.FeeKind,
			FeeRate:         v.FeeRate.String(),
			OrderCreateTime: v.OrderCreateTime.Format("2006-01-02 15:04:05"),
			OrderModifyTime: v.OrderModifyTime.Format("2006-01-02 15:04:05"),
		}
		orders = append(orders, order)
	}
	result := &pb.OrderQueryRsp{
		Count:  count,
		Orders: orders,
	}

	return result, nil
}
