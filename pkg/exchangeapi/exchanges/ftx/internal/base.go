package internal

import (
	"exterior-interactor/pkg/exchangeapi/exmodel"
	"exterior-interactor/pkg/exchangeapi/extools"
	"exterior-interactor/pkg/httptools"
)

type FtxBase struct {
	*FtxSigner
	extools.ExBase
}

func NewFtxBase(config exmodel.AccountConfig) (extools.ExBase, error) {
	signer := &FtxSigner{AccountConfig: config}

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

	return &FtxBase{
		FtxSigner: signer,
		ExBase: extools.NewExBase(signer, extools.ExBaseConfig{
			AccountConfig: config,
			Exchange:      exmodel.FTX,
			HttpClient:    client,
		}, requestInterceptor{}),
	}, nil

}
