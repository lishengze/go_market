package logic

import (
	"context"
	"market_server/common/nacosAdapter"
	"market_server/common/xerror"

	"github.com/pkg/errors"

	"market_server/app/dataService/rpc/internal/svc"
	"market_server/app/dataService/rpc/types/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSymbolLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSymbolLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSymbolLogic {
	return &GetSymbolLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取nacos参数
func (l *GetSymbolLogic) GetSymbol(in *pb.GetSymbolReq) (*pb.GetSymbolRsp, error) {

	l.Logger.Infof("DataService GetSymbolm in:%+v", in)

	var symbolList []*pb.SymbolInfo

	//获取一个
	if len(in.Symbol) > 0 {
		if r, ok := l.svcCtx.SymbolInfoMap.Load(in.Symbol); ok {
			res := r.(*nacosAdapter.Symbol)

			symbolList = append(symbolList, &pb.SymbolInfo{
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
			})

			return &pb.GetSymbolRsp{
				Count: 1,
				List:  symbolList,
			}, nil

		}
		return nil, errors.Wrapf(xerror.ErrorDSSymbol, "DataService symbol:%s", in.Symbol)
	}

	//获取全部

	l.svcCtx.SymbolInfoMap.Range(func(key, value interface{}) bool {
		res := value.(*nacosAdapter.Symbol)
		symbolList = append(symbolList, &pb.SymbolInfo{
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
		})
		return true
	})

	return &pb.GetSymbolRsp{
		Count: int64(len(symbolList)),
		List:  symbolList,
	}, nil
}
