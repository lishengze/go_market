package internal

import (
	"exterior-interactor/pkg/exchangeapi/exmodel"
	"exterior-interactor/pkg/exchangeapi/extools"
	"exterior-interactor/pkg/httptools"
)

type BinanceBase struct {
	*BinanceSigner
	extools.ExBase
}

func NewBinanceBase(config exmodel.AccountConfig) (extools.ExBase, error) {
	signer := &BinanceSigner{AccountConfig: config}

	var client *httptools.HttpClient
	if config.Proxy != "" {
		c, err := httptools.NewHttpClientWithProxy(config.Proxy)
		if err != nil {
			return nil, err
		}
		client = c
	} else {
		client = httptools.NewHttpClient()
	}

	return &BinanceBase{
		BinanceSigner: signer,
		ExBase: extools.NewExBase(signer, extools.ExBaseConfig{
			AccountConfig: config,
			Exchange:      exmodel.BINANCE,
			HttpClient:    client,
		}, requestInterceptor{}),
	}, nil
}
