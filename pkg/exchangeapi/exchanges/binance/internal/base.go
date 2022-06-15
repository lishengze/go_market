package internal

import (
	"exterior-interactor/pkg/exchangeapi/exmodel"
	"exterior-interactor/pkg/exchangeapi/extools"
)

type BinanceBase struct {
	*BinanceSigner
	extools.ExBase
}

func NewBinanceBase(config exmodel.AccountConfig) extools.ExBase {
	signer := &BinanceSigner{AccountConfig: config}
	return &BinanceBase{
		BinanceSigner: signer,
		ExBase: extools.NewExBase(signer, extools.ExBaseConfig{
			AccountConfig: config,
			Exchange:      exmodel.BINANCE,
			HttpClient:    Client,
		}, requestInterceptor{}),
	}
}
