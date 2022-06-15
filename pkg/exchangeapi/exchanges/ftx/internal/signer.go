package internal

import (
	"exterior-interactor/pkg/exchangeapi/exmodel"
	"exterior-interactor/pkg/exchangeapi/extools"
	"exterior-interactor/pkg/httptools"
	"fmt"
	"time"
)

var _ extools.Signer = new(FtxSigner)

type FtxSigner struct {
	exmodel.AccountConfig
}

func (o FtxSigner) Sign() func(params *httptools.IntegralParam) error {
	return func(params *httptools.IntegralParam) error {
		if o.AccountConfig == exmodel.EmptyAccountConfig {
			return fmt.Errorf("empty account config, can't sign ")
		}

		params.Url.RawQuery = params.Param.Encode()

		ts := fmt.Sprint(time.Now().UnixMilli())
		bodyStr, err := params.JsonBody.TrimmedString()
		if err != nil {
			return err
		}

		var pathWithParam string

		if params.Param.Encode() == "" {
			pathWithParam = params.Url.Path
		} else {
			pathWithParam = fmt.Sprintf("%s?%s", params.Url.Path, params.Param.Encode())
		}
		signData := fmt.Sprintf("%s%s%s%s", ts, params.HttpMethod, pathWithParam, bodyStr)
		//logx.Infof("signData: %s", signData)
		signature, err := extools.GetParamHmacSHA256Sign(o.Secret, signData)
		if err != nil {
			return err
		}
		params.Header.Set("FTX-KEY", o.AccountConfig.Key)
		params.Header.Set("FTX-SIGN", signature)
		params.Header.Set("FTX-TS", ts)

		if o.AccountConfig.SubAccountName != "" {
			params.Header.Set("FTX-SUBACCOUNT", o.AccountConfig.SubAccountName)
		}

		return nil
	}
}
