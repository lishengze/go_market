package binance

import (
	"exterior-interactor/pkg/exchangeapi/exmodel"
	"exterior-interactor/pkg/exchangeapi/extools"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
	"sync"
)

const (
	priceFilterStr       = "PRICE_FILTER"
	loftSizeFilterStr    = "LOT_SIZE"
	minNotionalFilterStr = "MIN_NOTIONAL"

	contractTrading = "TRADING"
)

const (
	// PERPETUAL 永续合约
	//CURRENT_QUARTER 当季交割合约
	//NEXT_QUARTER 次季交割合约
	contractTypePerp        = "PERPETUAL"
	contractTypeCurQuarter  = "CURRENT_QUARTER"
	contractTypeNextQuarter = "NEXT_QUARTER"
)

type symbolManager struct {
	Std  sync.Map // key: std symbol format, value: *exmodel.Symbol
	Spot sync.Map // key: ex symbol format, value: *exmodel.Symbol
	api  *NativeApi
}

func NewSymbolManager(api *NativeApi) extools.SymbolManager {
	mgr := &symbolManager{
		Std:  sync.Map{},
		Spot: sync.Map{},
		api:  api,
	}
	err := mgr.PullAllSymbol()
	if err != nil {
		panic(err)
	}
	return mgr
}

func (o *symbolManager) PullAllSymbol() error {
	return mr.Finish(func() error {
		return o.pullSpotSymbol()
	})
}

func (o *symbolManager) pullSpotSymbol() error {
	info, err := o.api.SpotApi.GetExchangeInfo()
	if err != nil {
		return err
	}
	for _, r := range info.Symbols {
		if r.Status != contractTrading {
			//logx.Infof("binance symbol %s status:%s ignored. ", r.Symbol, r.Status)
			continue
		}

		var (
			qtyScale, pxScale string
			minVolume         string
			stdFormat         = exmodel.StdSymbol(fmt.Sprintf("%s_%s", r.BaseAsset, r.QuoteAsset))
			exFormat          = r.Symbol
		)

		for _, filter := range r.Filters {
			if filter.FilterType == priceFilterStr {
				pxScale = filter.TickSize
			}
			if filter.FilterType == loftSizeFilterStr {
				qtyScale = filter.StepSize
				minVolume = filter.MinQty
			}
		}

		symbol := &exmodel.Symbol{
			Exchange:     exmodel.BINANCE,
			ExFormat:     exFormat,
			StdSymbol:    stdFormat,
			BaseCurrency: exmodel.Currency(r.BaseAsset),
			Type:         exmodel.SymbolTypeSpot,
			VolumeScale:  qtyScale,
			PriceScale:   pxScale,
			MinVolume:    minVolume,
			ContractSize: "1",
		}

		o.Std.Store(stdFormat, symbol)
		o.Spot.Store(exFormat, symbol)
	}
	return nil
}

func (o *symbolManager) GetSymbol(stdFormat exmodel.StdSymbol) (*exmodel.Symbol, error) {
	res, ok := o.Std.Load(stdFormat)
	if !ok {
		err := o.PullAllSymbol()
		if err != nil {
			logx.Errorf("pullAllSymbol err:%v", err)
			return nil, fmt.Errorf("cat not find symbol:%s in binance. ", stdFormat)
		}

		res2, ok2 := o.Std.Load(stdFormat)
		if !ok2 {
			return nil, fmt.Errorf("cat not find symbol:%s in binance. ", stdFormat)
		}

		symbolCopy := *res2.(*exmodel.Symbol)
		return &symbolCopy, nil
	}

	symbolCopy := *res.(*exmodel.Symbol)
	return &symbolCopy, nil
}

func (o *symbolManager) GetAllSymbol() []*exmodel.Symbol {
	res := make([]*exmodel.Symbol, 0)
	o.Std.Range(func(key, value interface{}) bool {
		symbolCopy := *value.(*exmodel.Symbol)
		res = append(res, &symbolCopy)
		return true
	})
	return res
}

func (o *symbolManager) Convert(exFormat string, apiType exmodel.ApiType) (*exmodel.Symbol, error) {
	switch apiType {
	case exmodel.ApiTypeSpot:
		res, ok := o.Spot.Load(exFormat)
		if !ok {
			err := o.pullSpotSymbol()
			if err != nil {
				logx.Errorf("pullSpotSymbol err:%v", err)
				return nil, fmt.Errorf("cat not parse symbol:%s in binance. ", exFormat)
			}
			res2, ok2 := o.Spot.Load(exFormat)
			if !ok2 {
				return nil, fmt.Errorf("cat not parse symbol:%s in binance. ", exFormat)
			}
			symbolCopy := *res2.(*exmodel.Symbol)
			return &symbolCopy, nil
		}

		symbolCopy := *res.(*exmodel.Symbol)
		return &symbolCopy, nil
	default:
		return nil, fmt.Errorf("not suppot apiType:%s", apiType)
	}

}
