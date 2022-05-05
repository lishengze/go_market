package ftx

import (
	"exterior-interactor/pkg/exchangeapi/exchanges/ftx/ftxapi"
	"exterior-interactor/pkg/exchangeapi/exchanges/ftx/internal"
	"exterior-interactor/pkg/exchangeapi/exmodel"
	"exterior-interactor/pkg/httptools"
)

const (
	Name = exmodel.FTX
)

// NativeApi 原生接口
type NativeApi struct {
	*ftxapi.Api
}

func NewNativeApi(config exmodel.AccountConfig) *NativeApi {
	base := internal.NewFtxBase(config)
	return &NativeApi{
		Api: ftxapi.NewApi(base),
	}
}

func NewNativeApiWithProxy(config exmodel.AccountConfig, proxy string) *NativeApi {
	if proxy == "" {
		return NewNativeApi(config)
	}

	c, err := httptools.NewHttpClientWithProxy(proxy)
	if err != nil {
		panic(err)
	}
	internal.Client = c

	return NewNativeApi(config)
}
