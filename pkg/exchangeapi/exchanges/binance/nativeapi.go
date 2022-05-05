package binance

import (
	"exterior-interactor/pkg/exchangeapi/exchanges/binance/api/binancespot"
	"exterior-interactor/pkg/exchangeapi/exchanges/binance/internal"
	"exterior-interactor/pkg/exchangeapi/exmodel"
	"exterior-interactor/pkg/httptools"
)

const (
	Name = exmodel.BINANCE
)

type (
	// NativeApi 交易所原生接口
	NativeApi struct {
		SpotApi *binancespot.Api
	}
)

func NewNativeApi(config exmodel.AccountConfig) *NativeApi {
	base := internal.NewBinanceBase(config)
	return &NativeApi{
		SpotApi: binancespot.NewApi(base),
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
