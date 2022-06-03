package email

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

const Name = "email"

func init() {
	notify.RegisterNotify(emailNotify{})
}

type emailNotify struct {
}

type (
	EmailApiPostData struct {
		Channel   string   `header:"X-Channel"`
		Token     string   `header:"X-Token"`
		From      string   `json:"from"`
		FromName  string   `json:"fromName"`
		To        []string `json:"to"`
		Subject   string   `json:"subject"`
		Content   string   `json:"content"`
		MessageId string   `json:"messageId"`
	}
	EmailApiResponse struct {
		Error   *EmailApiError    `json:"error,optional"`
		Success bool              `json:"success"`
		Data    *EmailApiRespData `json:"data,optional"`
	}
	EmailApiError struct {
		Code    string   `json:"code"`
		Message string   `json:"message"`
		Details []string `json:"details"`
	}

	EmailApiRespData struct {
		TemplateId      string `json:"templateId"`
		TemplateVersion string `json:"templateVersion"`
		Status          string `json:"status"`
		DateCreated     string `json:"dateCreated"`
		Locale          string `json:"locale"`
	}
)

func (emailNotify) Send(ctx context.Context, endpoint string, timeout int64, data interface{}) (interface{}, error) {
	if _, ok := data.(*EmailApiPostData); !ok {
		return nil, errors.Errorf("%v is not an object with type *email.EmailApiPostData", data)
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
	var val EmailApiResponse
	err = jsonx.Unmarshal(bodyData, &val)

	return &val, errors.WithStack(err)
}

func (emailNotify) Name() string {
	return Name
}
