package core

import (
	"exterior-interactor/app/opu/model"
	"exterior-interactor/pkg/exchangeapi/exchanges/ftx"
	"exterior-interactor/pkg/exchangeapi/exmodel"
	"exterior-interactor/pkg/exchangeapi/extools"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type symbolManager struct {
	*svcCtx
	extools.SymbolManager
	updateInterval time.Duration
}

func newSymbolManager(svcCtx *svcCtx, exchange, proxy string) *symbolManager {
	var sm extools.SymbolManager

	switch exchange {
	case exmodel.FTX.String():
		api, err := ftx.NewNativeApi(exmodel.AccountConfig{Proxy: proxy})
		if err != nil {
			panic(err)
		}
		sm = ftx.NewSymbolManager(api)
	default:
		panic(fmt.Sprintf("not support exchange:%s", exchange))
	}
	res := &symbolManager{
		svcCtx:         svcCtx,
		SymbolManager:  sm,
		updateInterval: time.Minute * 10,
	}
	go res.run()
	return res
}

func (o *symbolManager) run() {
	for {
		o.pullSymbols()
		o.updateSymbols()
		time.Sleep(o.updateInterval)
	}
}

// pullSymbols 更新所有symbol
func (o *symbolManager) pullSymbols() {
	err := o.SymbolManager.PullAllSymbol()
	if err != nil {
		logx.Errorf("ftx PullAllSymbol err:%v", err)
	}
}

// updateSymbols 更新所有symbol
func (o *symbolManager) updateSymbols() {
	var symbols []*model.Symbol

	res := o.SymbolManager.GetAllSymbol()
	for _, s := range res {
		symbol, err := o.svcCtx.SymbolModel.FindOneByStdSymbolExchange(s.StdSymbol.String(), s.Exchange.String())
		switch err {
		case model.ErrNotFound:
			symbols = append(symbols, &model.Symbol{
				Id:            o.svcCtx.IdSrv.MustGetId(),
				Tp:            s.Type.String(),
				ApiType:       s.ApiType.String(),
				StdSymbol:     s.StdSymbol.String(),
				ExFormat:      s.ExFormat,
				Exchange:      s.Exchange.String(),
				BaseCurrency:  s.BaseCurrency.String(),
				QuoteCurrency: s.QuoteCurrency.String(),
				VolumeScale:   s.VolumeScale,
				PriceScale:    s.PriceScale,
				MinVolume:     s.MinVolume,
				ContractSize:  s.ContractSize,
			})
		case nil:
			// 如果有不同，则更新
			if symbol.VolumeScale != s.VolumeScale || symbol.PriceScale != s.PriceScale ||
				symbol.MinVolume != s.MinVolume || symbol.ContractSize != s.ContractSize {
				err := o.SymbolModel.Update(symbol, func() {
					symbol.VolumeScale = s.VolumeScale
					symbol.PriceScale = s.PriceScale
					symbol.MinVolume = s.MinVolume
					symbol.ContractSize = s.ContractSize
				})
				if err != nil {
					logx.Errorf("update symbol err:%v, symbol:%+v", err, *symbol)
				}
			}
		default:
			logx.Error(err)
		}
	}

	err := o.svcCtx.SymbolModel.BulkInsert(symbols)
	if err != nil {
		logx.Errorf("BulkInsert err:%v", err)
	}
}
