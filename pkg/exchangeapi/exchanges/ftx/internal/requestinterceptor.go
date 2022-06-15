package internal

import (
	"exterior-interactor/pkg/exchangeapi/extools"
	"exterior-interactor/pkg/httptools"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
)

type requestInterceptor struct{}

func (o requestInterceptor) BeforeRequest(meta extools.Meta, request httptools.Request) error {
	if request.GetIntegralParam().Header.Get(httptools.ContentType) == "" {
		request.SetParams(func(params *httptools.IntegralParam) error {
			if len(params.JsonBody) != 0 {
				params.Header.Set(httptools.ContentType, httptools.ContentTypeJson)
				return nil
			}
			params.Header.Set(httptools.ContentType, httptools.ContentTypeForm)
			return nil
		})
	}
	return DefaultLimiterManager.ThroughLimiters(request, meta)
}

func (o requestInterceptor) AfterRequest(meta extools.Meta, request httptools.Request, rsp *http.Response) error {
	switch rsp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusTooManyRequests:
	case http.StatusTeapot:
	}

	bytes, err := httptools.DecodeResponseToBytes(rsp)
	if err != nil {
		logx.Error(err, string(bytes))
	}

	logx.Errorf("reqInfo:%s http code:%v, headers:%v, content:%s",
		request.ReqInfo(), rsp.StatusCode, rsp.Header, string(bytes))

	return fmt.Errorf("httpcode:%v, body:%s", rsp.StatusCode, string(bytes))
}
