package dingtalk

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	nethttp "net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/jsonx"

	"bcts/internal/client/http"
	"bcts/internal/notify"

	"github.com/pkg/errors"
)

const Name = "ding_talk"

var dingTalkUrl url.URL = url.URL{
	Host:   "oapi.dingtalk.com",
	Scheme: "https",
	Path:   "robot/send",
}

func init() {
	notify.RegisterNotify(dingTalkNotify{})
}

type dingTalkNotify struct {
}

// Todo
func (dingTalkNotify) Send(ctx context.Context, endpoint string, timeout int64, data interface{}) (interface{}, error) {
	if _, ok := data.(*Message); !ok {
		return nil, errors.Errorf("%v is not an object with type *dingtalk.Message", data)
	}

	if ctx == nil {
		ctx = context.Background()
	}

	t := 30 * time.Millisecond
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
	var val Response
	err = jsonx.Unmarshal(bodyData, &val)

	return &val, errors.WithStack(err)
}

func (dingTalkNotify) Name() string {
	return Name
}

func GetRobotSendUrl(accessToken, secret string) (string, error) {
	dtu := dingTalkUrl
	value := url.Values{}
	value.Set("access_token", accessToken)

	if len(secret) == 0 {
		dtu.RawQuery = value.Encode()
		return dtu.String(), nil
	}
	timestamp := strconv.FormatInt(time.Now().Unix()*1000, 10)
	sign, err := sign(timestamp, secret)
	if err != nil {
		dtu.RawQuery = value.Encode()
		return dtu.String(), nil
	}

	value.Set("timestamp", timestamp)
	value.Set("sign", sign)
	dtu.RawQuery = value.Encode()
	return dtu.String(), nil
}

func sign(timestamp string, secret string) (string, error) {
	stringToSign := fmt.Sprintf("%s\n%s", timestamp, secret)
	h := hmac.New(sha256.New, []byte(secret))
	if _, err := io.WriteString(h, stringToSign); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}
