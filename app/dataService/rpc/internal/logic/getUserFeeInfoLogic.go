package logic

import (
	"bcts/common/nacosAdapter"
	"bcts/common/utils"
	"bcts/common/xerror"
	"context"
	"github.com/pkg/errors"

	"bcts/app/dataService/rpc/internal/svc"
	"bcts/app/dataService/rpc/types/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserFeeInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserFeeInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserFeeInfoLogic {
	return &GetUserFeeInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取用户交易手续费
func (l *GetUserFeeInfoLogic) GetUserFeeInfo(in *pb.GetUserFeeInfoReq) (*pb.GetUserFeeInfoRsp, error) {
	//1. symbol不能为空
	if len(in.Symbol) == 0 {
		return nil, errors.Wrapf(xerror.ErrorSymbolIsEmpty, "in:%+v", in)
	}

	r, ok := l.svcCtx.SymbolInfoMap.Load(in.Symbol)
	if !ok {
		return nil, errors.Wrapf(xerror.ErrorDSSymbol, "symbol not found, in:%+v", in)
	}

	res := r.(*nacosAdapter.Symbol)
	targetCurrency, baseCurrency := utils.SeparateSymbol(in.Symbol)

	//2. 用户组参数

	symbolInfo := &pb.SymbolInfo{
		SymbolID:             res.SymbolID,
		SymbolKind:           int64(res.SymbolKind),
		Underlying:           res.Underlying,
		PrimaryCurrency:      res.Underlying,
		BidCurrency:          res.BidCurrency,
		SettleCurrency:       res.SettleCurrency,
		Switch:               res.Switch,
		VolumePrecision:      int64(res.VolumePrecision),
		PricePrecision:       int64(res.VolumePrecision),
		AmountPrecision:      int64(res.AmountPrecision),
		MinUnit:              res.MinUnit,
		MinChangePrice:       res.MinChangePrice,
		Spread:               res.Spread,
		FeeKind:              int64(res.FeeKind),
		TakerFee:             res.TakerFee,
		MakerFee:             res.MakerFee,
		SingleMaxOrderAmount: res.SingleMinOrderAmount,
		SingleMinOrderAmount: res.SingleMinOrderAmount,
		SingleMinOrderVolume: res.SingleMinOrderVolume,
		SingleMaxOrderVolume: res.SingleMaxOrderVolume,
		BuyPriceLimit:        res.BuyPriceLimit,
		SellPriceLimit:       res.SellPriceLimit,
		MaxMatchGear:         int64(res.MaxMatchGear),
		OtcMinOrderVolume:    res.OtcMinOrderVolume,
		OtcMaxOrderVolume:    res.OtcMaxOrderVolume,
		OtcMinOrderAmount:    res.OtcMinOrderAmount,
		OtcMaxOrderAmount:    res.OtcMaxOrderAmount,
	}

	return &pb.GetUserFeeInfoRsp{
		UserID:         in.UserID,
		Symbol:         in.Symbol,
		BaseCurrency:   baseCurrency,
		TargetCurrency: targetCurrency,
		FeeKind:        int32(res.FeeKind),
		TakerFee:       res.TakerFee,
		MakerFee:       res.MakerFee,
		SymbolInfo:     symbolInfo,
	}, nil

}
