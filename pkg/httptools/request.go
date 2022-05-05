package httptools

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	ContentType = "Content-Type"
	UserAgent   = "User-Agent"
)

const (
	ContentTypeForm = "application/x-www-form-urlencoded"
	ContentTypeJson = "application/json"

	DefaultUserAgent = "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36"
)

var _ Request = new(request)

type (
	Request interface {
		SetReq(req interface{}) Request // 如果调用此方法，需要首先被调用
		SetUrl(url string) Request
		SetProxy(proxy string) Request
		SetHeader(k, v string) Request
		SetParams(fns ...func(p *IntegralParam) error) Request
		SetHttpMethod(method string) Request
		Request() (*http.Response, error)
		GetIntegralParam() *IntegralParam
		ReqInfo() string
	}

	request struct {
		err error
		//proxy      string
		httpClient *HttpClient
		*IntegralParam
	}
)

func NewRequest() Request {
	return &request{
		httpClient:    NewHttpClient(),
		IntegralParam: NewIntegralParam(),
	}
}

func NewRequestWithHttpClient(client *HttpClient) Request {
	return &request{
		httpClient:    client,
		IntegralParam: NewIntegralParam(),
	}
}

func (o *request) SetUrl(s string) Request {
	o.IntegralParam.Url, o.err = url.Parse(s)
	return o
}

func (o *request) GetIntegralParam() *IntegralParam {
	return o.IntegralParam
}

func (o *request) SetProxy(proxy string) Request {
	uProxy, err := url.Parse(proxy)
	if err != nil {
		o.err = fmt.Errorf("proxy url err:%v ", err)
		return o
	}

	tr := &http.Transport{
		MaxIdleConnsPerHost: 200,
		MaxIdleConns:        200,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		Proxy:               http.ProxyURL(uProxy),
	}

	o.httpClient.Transport = tr
	return o
}

func (o *request) SetHeader(k, v string) Request {
	o.IntegralParam.Header.Set(k, v)
	return o
}

func (o *request) SetReq(req interface{}) Request {
	ip, err := ParseReqParam(req)
	if err != nil {
		o.err = err
	} else {
		o.IntegralParam = ip
	}
	return o
}

func (o *request) SetParams(fns ...func(p *IntegralParam) error) Request {
	for _, fn := range fns {
		err := fn(o.IntegralParam)
		if err != nil {
			o.err = err
			return o
		}
	}
	return o
}

func (o *request) SetHttpMethod(m string) Request {
	o.IntegralParam.HttpMethod = m
	return o
}

func (o *request) ReqInfo() string {
	return fmt.Sprintf("[%s] url:%s, proxy:%s, header:%v, form:%v, body:%v",
		o.HttpMethod, o.Url.String(), o.httpClient.Proxy, o.Header, o.Form, o.JsonBody)
}

func (o *request) Request() (*http.Response, error) {
	if o.err != nil {
		return nil, o.err
	}

	if o.IntegralParam.Header.Get(UserAgent) == "" {
		o.IntegralParam.Header.Set(UserAgent, DefaultUserAgent)
	}

	//如果参数中有中文参数,这个方法会进行URLEncode
	o.Url.RawQuery = o.IntegralParam.Param.Encode()
	path := o.Url.String()

	var httpRequest *http.Request

	switch o.IntegralParam.Header.Get(ContentType) {
	case ContentTypeForm:
		req, err := http.NewRequest(o.HttpMethod, path, strings.NewReader(o.IntegralParam.Form.Encode()))
		if err != nil {
			return nil, err
		}
		httpRequest = req

	case ContentTypeJson:
		bodyBytes, err := json.Marshal(o.IntegralParam.JsonBody)
		if err != nil {
			return nil, err
		}
		req, err := http.NewRequest(o.HttpMethod, path, bytes.NewReader(bodyBytes))
		if err != nil {
			return nil, err
		}
		httpRequest = req
	case "":
		return nil, fmt.Errorf("not set content-type. ")
	default:
		return nil, fmt.Errorf("not support content-type:%s. ", o.IntegralParam.Header.Get(ContentType))
	}

	httpRequest.Close = true
	httpRequest.Header = o.IntegralParam.Header

	startTime := time.Now()
	rsp, err := o.httpClient.Do(httpRequest)
	duration := time.Now().Sub(startTime)

	if err != nil {
		return rsp, err
	}

	logx.WithDuration(duration).Info(o.ReqInfo())

	return rsp, nil
}

func DecodeResponseToBytes(rsp *http.Response) ([]byte, error) {
	defer func() {
		_ = rsp.Body.Close()
	}()
	bodyData, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	return bodyData, nil
}
