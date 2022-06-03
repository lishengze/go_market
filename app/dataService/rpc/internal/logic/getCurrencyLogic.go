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

type GetCurrencyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCurrencyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCurrencyLogic {
	return &GetCurrencyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取币种参数
func (l *GetCurrencyLogic) GetCurrency(in *pb.GetCurrencyReq) (*pb.GetCurrencyRsp, error) {
	l.Logger.Infof("DataService GetCurrency in:%+v", in)

	var currencyList []*pb.CurrencyInfo
	//获取一个
	if len(in.Currency) > 0 {
		if r, ok := l.svcCtx.CurrencyInfoMap.Load(in.Currency); ok {
			res := r.(*nacosAdapter.Currency)

			currencyList = append(currencyList, &pb.CurrencyInfo{
				CurrencyID:     res.CurrencyID,
				Kind:           int64(res.Kind),
				ChineseName:    res.ChineseName,
				EnglishName:    res.EnglishName,
				MinUnit:        res.MinUnit,
				DepositSwitch:  res.DepositSwitch,
				WithdrawSwitch: res.WithdrawSwitch,
				MinWithdraw:    res.MinWithdraw,
				MaxWithdraw:    res.MaxWithdraw,
				MaxDayWithdraw: res.MaxDayWithdraw,
				Threshold:      res.Threshold,
				FeeKind:        int64(res.FeeKind),
				Fee:            res.Fee,
				Operator:       res.Operator,
				Time:           res.Time,
			})

			return &pb.GetCurrencyRsp{
				Count: 1,
				List:  currencyList,
			}, nil

		}
		return nil, errors.Wrapf(xerror.ErrorDSCurrency, "DataService currency:%s", in.Currency)
	}

	//获取全部
	l.svcCtx.CurrencyInfoMap.Range(func(key, value interface{}) bool {
		res := value.(*nacosAdapter.Currency)

		currencyList = append(currencyList, &pb.CurrencyInfo{
			CurrencyID:     res.CurrencyID,
			Kind:           int64(res.Kind),
			ChineseName:    res.ChineseName,
			EnglishName:    res.EnglishName,
			MinUnit:        res.MinUnit,
			DepositSwitch:  res.DepositSwitch,
			WithdrawSwitch: res.WithdrawSwitch,
			MinWithdraw:    res.MinWithdraw,
			MaxWithdraw:    res.MaxWithdraw,
			MaxDayWithdraw: res.MaxDayWithdraw,
			Threshold:      res.Threshold,
			FeeKind:        int64(res.FeeKind),
			Fee:            res.Fee,
			Operator:       res.Operator,
			Time:           res.Time,
		})
		return true
	})

	return &pb.GetCurrencyRsp{
		Count: int64(len(currencyList)),
		List:  currencyList,
	}, nil
}
