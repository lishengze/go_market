package sms

import (
	"context"
	"crypto/tls"
	"io/ioutil"
	nethttp "net/http"
	"time"

	"github.com/zeromicro/go-zero/core/jsonx"

	"github.com/pkg/errors"

	"bcts/internal/client/http"
	"bcts/internal/notify"
)

const Name = "sms"

func init() {
	notify.RegisterNotify(smsNotify{})
}

type (
	SmsApiPostData struct {
		Channel string `header:"X-Channel"`
		Token   string `header:"X-Token"`
		To      string `json:"to"`
		Content string `json:"content"`
	}
	SmsApiResponse struct {
		Error   *SmsApiError    `json:"error,optional"`
		Success bool            `json:"success"`
		Data    *SmsApiRespData `json:"data,optional"`
	}
	SmsApiError struct {
		Code       string   `json:"code"`
		Message    string   `json:"message"`
		StatusCode string   `json:"statusCode"`
		Details    []string `json:"details"`
	}

	SmsApiRespData struct {
		Name        string `json:"name"`
		Status      string `json:"status"`
		DateCreated string `json:"dateCreated"`
		Locale      string `json:"locale"`
	}
)

type smsNotify struct {
}

func (smsNotify) Send(ctx context.Context, endpoint string, timeout int64, data interface{}) (interface{}, error) {
	if _, ok := data.(*SmsApiPostData); !ok {
		return nil, errors.Errorf("%v is not an object with type *email.SmsApiPostData", data)
	}

	t := 600 * time.Millisecond
	if timeout > 0 {
		t = time.Duration(timeout) * time.Millisecond
	}

	// 实例化httpclient
	client, err := http.NewClient(
		ctx,
		endpoint,
		http.WithTimeout(t),
		http.WithTransport(&nethttp.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}),
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// 请求接口
	resp, err := client.Do(nethttp.MethodPost, data)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer resp.Body.Close()

	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// 解析body数据
	var val SmsApiResponse
	err = jsonx.Unmarshal(bodyData, &val)

	return &val, errors.WithStack(err)
}

func (smsNotify) Name() string {
	return Name
}
