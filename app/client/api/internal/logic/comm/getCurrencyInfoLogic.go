package comm

import (
	"bcts/app/client/api/internal/svc"
	"bcts/app/client/api/internal/types"
	"context"
	"github.com/jinzhu/copier"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCurrencyInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewGetCurrencyInfoLogic(r *http.Request, ctx context.Context, svcCtx *svc.ServiceContext) *GetCurrencyInfoLogic {
	return &GetCurrencyInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCurrencyInfoLogic) GetCurrencyInfo(req *types.CurrencyInfoReq) (resp *types.CurrencyInfoRsp, err error) {

	res, err := l.svcCtx.NacosClient.GetCurrency(req.Currency)
	if err != nil {
		return nil, err
	}

	var currencyList []*types.CurrencyInfo
	if len(res) > 0 {
		for _, v := range res {
			var currency types.CurrencyInfo
			copier.Copy(&currency, v)
			currencyList = append(currencyList, &currency)
		}
	}

	return &types.CurrencyInfoRsp{
		Count: len(currencyList),
		List:  currencyList,
	}, nil
}
