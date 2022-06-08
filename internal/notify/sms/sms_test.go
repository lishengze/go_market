package sms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/zeromicro/go-zero/core/mapping"

	"market_server/internal/notify"
)

var mockServer = httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
	//w.Write()
	var mockEmailApiServerResp = SmsApiResponse{
		//Error:   SmsApiError{Message: "sssss", Code: "xxxx", Details: []string{}},
		Success: true,
		Data: &SmsApiRespData{
			Name:        "hub.openapi",
			Status:      "xxxx",
			DateCreated: "xxxx",
			Locale:      "en",
		},
	}
	var val map[string]map[string]interface{}
	val, err := mapping.Marshal(mockEmailApiServerResp)
	if err != nil {
		fmt.Println(err)
		return
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(val["json"]); err != nil {
		fmt.Println(err)
		return
	}

	w.WriteHeader(nethttp.StatusOK)
	w.Write(buf.Bytes())
}))

func TestSmsNotify_Send(t *testing.T) {
	notity := notify.GetNotify("sms")
	data := &SmsApiPostData{
		Channel: "trading",
		Token:   "xxxxxx",
		To:      "15221560954",
		Content: "test",
	}

	resp, err := notity.Send(context.Background(), mockServer.URL, 60, data)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%+v", resp)
}
