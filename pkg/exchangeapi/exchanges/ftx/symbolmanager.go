package ftx

import (
	"exterior-interactor/pkg/exchangeapi/exmodel"
	"exterior-interactor/pkg/exchangeapi/extools"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
	"sync"
)

type symbolManager struct {
	Std sync.Map // key: std symbol format, value: *exmodel.Symbol
	Ex  sync.Map // key: ex symbol format, value: *exmodel.Symbol
	api *NativeApi
}

func NewSymbolManager(api *NativeApi) extools.SymbolManager {
	mgr := &symbolManager{
		Std: sync.Map{},
		Ex:  sync.Map{},
		api: api,
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
	}, func() error {
		return o.pullFutureSymbol()
	})
}

func (o *symbolManager) pullSpotSymbol() error {
	rsp, err := o.api.Api.GetMarket()
	if err != nil {
		return err
	}
	for _, r := range rsp.Result {
		if r.Type != "spot" {
			continue
		}

		var (
			stdFormat = exmodel.StdSymbol(fmt.Sprintf("%s_%s", r.BaseCurrency, r.QuoteCurrency))
			exFormat  = r.Name
		)

		symbol := &exmodel.Symbol{
			Exchange:     Name,
			ExFormat:     exFormat,
			StdSymbol:    stdFormat,
			BaseCurrency: exmodel.Currency(r.BaseCurrency),
			Type:         exmodel.SymbolTypeSpot,
			VolumeScale:  fmt.Sprint(r.SizeIncrement),
			PriceScale:   fmt.Sprint(r.PriceIncrement),
			MinVolume:    fmt.Sprint(r.MinProvideSize),
			ContractSize: "1",
		}

		//fmt.Println(symbol.ExFormat,"_____",symbol.PriceScale,"_______",symbol.VolumeScale)
		o.Std.Store(stdFormat, symbol)
		o.Ex.Store(exFormat, symbol)
	}
	return nil
}

func (o *symbolManager) pullFutureSymbol() error {
	// todo
	return nil
}

func (o *symbolManager) GetSymbol(stdFormat exmodel.StdSymbol) (*exmodel.Symbol, error) {
	res, ok := o.Std.Load(stdFormat)
	if !ok {
		err := o.PullAllSymbol()
		if err != nil {
			logx.Errorf("pullAllSymbol err:%v", err)
			return nil, fmt.Errorf("cat not find symbol:%s in %s. ", stdFormat, Name)
		}

		res2, ok2 := o.Std.Load(stdFormat)
		if !ok2 {
			return nil, fmt.Errorf("cat not find symbol:%s in %s. ", stdFormat, Name)
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
	case exmodel.ApiTypeUnified:
		res, ok := o.Ex.Load(exFormat)
		if !ok {
			err := o.pullSpotSymbol()
			if err != nil {
				logx.Errorf("pullSpotSymbol err:%v", err)
				return nil, fmt.Errorf("cat not parse symbol:%s in %s . ", exFormat, Name)
			}
			res2, ok2 := o.Ex.Load(exFormat)
			if !ok2 {
				return nil, fmt.Errorf("cat not parse symbol:%s in %s . ", exFormat, Name)
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
