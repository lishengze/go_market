package extools

import (
	"encoding/json"
	"exterior-interactor/pkg/exchangeapi/exmodel"
	"exterior-interactor/pkg/httptools"
	"github.com/zeromicro/go-zero/core/logx"
	"net/url"
)

type (
	ExBaseConfig struct {
		exmodel.AccountConfig
		exmodel.Exchange
		*httptools.HttpClient
	}

	Signer interface {
		Sign() func(params *httptools.IntegralParam) error
	}

	ExBase interface {
		Config() ExBaseConfig
		Request(req, res interface{}, meta Meta, fns ...func(params *httptools.IntegralParam) error) error
	}

	defaultExBase struct {
		interceptors []RequestInterceptor
		Signer
		ExBaseConfig
	}
)

func NewExBase(signer Signer, config ExBaseConfig, interceptors ...RequestInterceptor) ExBase {
	return &defaultExBase{
		interceptors: interceptors,
		ExBaseConfig: config,
		Signer:       signer,
	}
}

func (o *defaultExBase) Config() ExBaseConfig {
	return o.ExBaseConfig
}

func (o *defaultExBase) Request(req, res interface{}, meta Meta, fns ...func(params *httptools.IntegralParam) error) error {
	request := httptools.NewRequestWithHttpClient(o.ExBaseConfig.HttpClient).SetReq(req).
		SetParams(func(params *httptools.IntegralParam) error {
			params.HttpMethod = meta.HttpMethod()
			u, err := url.Parse(meta.Url())
			params.Url = u
			return err
		}).SetParams(fns...)

	if meta.NeedSign() {
		request.SetParams(o.Sign())
	}

	for _, interceptor := range o.interceptors {
		err := interceptor.BeforeRequest(meta, request)
		if err != nil {
			return err
		}
	}

	rsp, err := request.Request()
	if err != nil {
		return err
	}

	for _, interceptor := range o.interceptors {
		err = interceptor.AfterRequest(meta, request, rsp)
		if err != nil {
			return err
		}
	}

	bytes, err := httptools.DecodeResponseToBytes(rsp)
	if err != nil {
		return err
	}

	//fmt.Println(string(bytes))

	err = json.Unmarshal(bytes, &res)
	if err != nil {
		logx.Errorf("err:%s bytes: %s", err, string(bytes))
		return err
	}

	return nil
}
