package binance

import (
	"exterior-interactor/pkg/exchangeapi/exchanges/binance/api/binancespot"
	"exterior-interactor/pkg/exchangeapi/exchanges/binance/internal"
	"exterior-interactor/pkg/exchangeapi/exmodel"
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

func NewNativeApi(config exmodel.AccountConfig) (*NativeApi, error) {
	base, err := internal.NewBinanceBase(config)
	if err != nil {
		return nil, err
	}
	return &NativeApi{
		SpotApi: binancespot.NewApi(base),
	}, nil
}
