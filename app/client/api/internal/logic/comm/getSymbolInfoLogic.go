package comm

import (
	"context"
	"net/http"

	"github.com/jinzhu/copier"

	"market_server/app/client/api/internal/svc"
	"market_server/app/client/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSymbolInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewGetSymbolInfoLogic(r *http.Request, ctx context.Context, svcCtx *svc.ServiceContext) *GetSymbolInfoLogic {
	return &GetSymbolInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
	}
}

func (l *GetSymbolInfoLogic) GetSymbolInfo(req *types.SymbolInfoReq) (resp *types.SymbolInfoRsp, err error) {
	res, err := l.svcCtx.NacosClient.GetSymbol(req.Symbol)
	if err != nil {
		return nil, err
	}

	var symbolList []*types.SymbolInfo
	if len(res) > 0 {
		for _, v := range res {
			var symbol types.SymbolInfo
			copier.Copy(&symbol, v)
			symbolList = append(symbolList, &symbol)
		}
	}

	return &types.SymbolInfoRsp{
		Count: len(symbolList),
		List:  symbolList,
	}, nil
}
