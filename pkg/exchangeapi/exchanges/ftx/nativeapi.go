package ftx

import (
	"exterior-interactor/pkg/exchangeapi/exchanges/ftx/ftxapi"
	"exterior-interactor/pkg/exchangeapi/exchanges/ftx/internal"
	"exterior-interactor/pkg/exchangeapi/exmodel"
)

const (
	Name = exmodel.FTX
)

// NativeApi 原生接口
type NativeApi struct {
	*ftxapi.Api
}

func NewNativeApi(config exmodel.AccountConfig) (*NativeApi, error) {
	base, err := internal.NewFtxBase(config)
	if err != nil {
		return nil, err
	}
	return &NativeApi{
		Api: ftxapi.NewApi(base),
	}, nil
}
