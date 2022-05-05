package internal

import (
	"exterior-interactor/pkg/exchangeapi/exmodel"
	"exterior-interactor/pkg/exchangeapi/extools"
)

type FtxBase struct {
	*FtxSigner
	extools.ExBase
}

func NewFtxBase(config exmodel.AccountConfig) extools.ExBase {
	signer := &FtxSigner{AccountConfig: config}
	return &FtxBase{
		FtxSigner: signer,
		ExBase: extools.NewExBase(signer, extools.ExBaseConfig{
			Exchange:   exmodel.FTX,
			HttpClient: Client,
		}, requestInterceptor{}),
	}
}
