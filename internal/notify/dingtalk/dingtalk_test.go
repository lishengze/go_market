package dingtalk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/zeromicro/go-zero/core/mapping"

	"bcts/internal/notify"
)

func TestGetRobotSendUrl(t *testing.T) {
	//notify := notify2.GetNotify("dingTalk")
	endpoint, err := GetRobotSendUrl("db898ab53bb81ceb95635034a6c326a62e0c56dce511fce45dd3441ef297fd54", "SEC440dc8873cca42bcfe5b57754e7b4477f1358551713366007e3998cebb8cb828")
	if err != nil {
		t.Error(err)
	}
	t.Log(endpoint)
}

var mockServer = httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
	//w.Write()
	var mockServerResp = Response{
		ErrCode: 0,
		ErrMsg:  "ok",
	}
	var val map[string]map[string]interface{}
	val, err := mapping.Marshal(mockServerResp)
	if err != nil {
		fmt.Println(err)
		return
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(val); err != nil {
		fmt.Println(err)
		return
	}

	w.WriteHeader(nethttp.StatusOK)
	w.Write(buf.Bytes())
}))

func TestDingTalkNotify_Send(t *testing.T) {
	endpoint, err := GetRobotSendUrl("db898ab53bb81ceb95635034a6c326a62e0c56dce511fce45dd3441ef297fd54", "SEC440dc8873cca42bcfe5b57754e7b4477f1358551713366007e3998cebb8cb828")
	if err != nil {
		t.Error(err)
	}
	notify := notify.GetNotify("ding_talk")
	m := &LinkMessage{
		Text:       "这个即将发布的新版本，创始人xx称它为红树林。而在此之前，每当面临重大升级，产品经理们都会取一个应景的代号，这一次，为什么是红树林",
		Title:      "时代的火车向前开test",
		PicUrl:     "",
		MessageUrl: "https://www.dingtalk.com/s?__biz=MzA4NjMwMTA2Ng==&mid=2650316842&idx=1&sn=60da3ea2b29f1dcc43a7c8e4a7c97a16&scene=2&srcid=09189AnRJEdIiWVaKltFzNTw&from=timeline&isappinstalled=0&key=&ascene=2&uin=&devicetype=android-23&version=26031933&nettype=WIFI",
	}
	msg, err := NewMessageRequest(MSG_TYPE_LINK, m, nil)

	if err != nil {
		t.Error(err)
	}
	// 访问mockserver
	endpoint = mockServer.URL

	res, err := notify.Send(context.Background(), endpoint, 30, msg)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%+v", res)
}
