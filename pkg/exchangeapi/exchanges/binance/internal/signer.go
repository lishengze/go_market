package internal

import (
	"exterior-interactor/pkg/exchangeapi/exmodel"
	"exterior-interactor/pkg/exchangeapi/extools"
	"exterior-interactor/pkg/httptools"
	"fmt"
	"time"
)

var _ extools.Signer = new(BinanceSigner)

type BinanceSigner struct {
	exmodel.AccountConfig
}

func (o BinanceSigner) Sign() func(params *httptools.IntegralParam) error {
	return func(params *httptools.IntegralParam) error {
		if o.AccountConfig==exmodel.EmptyAccountConfig{
			return fmt.Errorf("empty account config, can't sign ")
		}
		params.Header.Set("X-MBX-APIKEY", o.AccountConfig.Key)
		params.Param.Set("recvWindow", "10000")
		params.Param.Set("timestamp", fmt.Sprint(time.Now().UnixMilli()))
		payload := params.Param.Encode()

		sign, err := extools.GetParamHmacSHA256Sign(o.AccountConfig.Secret, payload)
		if err != nil {
			return err
		}
		params.Param.Set("signature", sign)

		return nil
	}
}
