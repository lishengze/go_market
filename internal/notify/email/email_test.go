package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"market_server/internal/notify"
	nethttp "net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/mapping"
)

var mockServer = httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
	//w.Write()
	var mockEmailApiServerResp = EmailApiResponse{
		//Error:   EmailApiError{Message: "sssss", Code: "xxxx", Details: []string{}},
		Success: true,
		Data: &EmailApiRespData{
			TemplateId:      "d-4c2f4e8997f14ec6b828e978e204a690",
			Status:          "QUEUED",
			TemplateVersion: "1",
			DateCreated:     "2022-05-18T10:24:11.569047",
			Locale:          "en",
		},
	}
	//`{"data":{"requestChannel":"trading","from":"Xpert.noreply@hashkey.com","fromName":"Hashkey Xpert QA","to":["hucuiliang@hashkeytech.com"],"name":"hub.openapi","subject":"Hashkey Xpert：Deposit","locale":"en","templateId":"d-4c2f4e8997f14ec6b828e978e204a690","templateVersion":"1","messageId":"","status":"QUEUED""dateCreated":"2022-05-18T10:24:11.569047"},"success":true}`

	var val map[string]map[string]interface{}

	val, err := mapping.Marshal(mockEmailApiServerResp)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v", val)
	fmt.Println()
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(val["json"]); err != nil {
		fmt.Println(err)
		return
	}

	w.WriteHeader(nethttp.StatusOK)
	w.Write(buf.Bytes())
}))

type ErrEmailApiPostData struct {
	Channel   string `header:"X-Channel"`
	Token     string `header:"X-Token"`
	From      string `json:"from"`
	FromName  string `json:"fromName"`
	To        string `json:"to"`
	Subject   string `json:"subject"`
	Content   string `json:"content"`
	MessageId string `json:"messageId"`
}

func TestEmailNotify_Send(t *testing.T) {
	notity := notify.GetNotify("email")
	data := &EmailApiPostData{
		Channel:  "trading",
		Token:    "xxxxxx",
		From:     "dev@hashkeyxpert.com",
		FromName: "Hashkey Xpert",
		To:       []string{"hucuiliang@hashkeytech.com"},
		Subject:  "test",
		Content:  "test",
	}

	res, err := notity.Send(context.Background(), mockServer.URL, 60, data)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", res)
}

func TestEmailNotify_ErrEmailApiPostData(t *testing.T) {
	notity := notify.GetNotify("email")
	data := &ErrEmailApiPostData{
		Channel:  "trading",
		Token:    "xxxxxx",
		From:     "dev@hashkeyxpert.com",
		FromName: "Hashkey Xpert",
		To:       "hucuiliang@hashkeytech.com",
		Subject:  "test",
		Content:  "test",
	}

	_, err := notity.Send(context.Background(), mockServer.URL, 60, data)
	expected := errors.Errorf("%v is not an object with type *email.EmailApiPostData", data)
	if !reflect.DeepEqual(err, expected) {
		t.Logf("Expected %v, but got: %v", err, expected)
	}
}
