# Notify

集成邮件、短信、钉钉等消息发送

## Usage

### 邮件
```go
import "bcts/internal/notify"

notify := notify.GetNotify("email")
data := &EmailApiPostData{
    Channel:  "trading",
    Token:    "xxxxxx",
    From:     "from@email",
    FromName: "from@email",
    To:       []string{"to@email"},
    Subject:  "test",
    Content:  "test",
}

resp, err := notify.Send(context.Background(), mockServer.URL, 60, data)
```

### 短信
```go
import "bcts/internal/notify"

notify := notify.GetNotify("sms")
data := &SmsApiPostData{
    Channel:  "trading",
    Token:    "xxxxxx",
    From:     "from phone num",
    To:       "to phone num",
    Content:  "test",
}

resp, err := notify.Send(context.Background(), mockServer.URL, 60, data)
```

### 钉钉群机器人

```go
// 接入指引：https://open.dingtalk.com/document/robots/custom-robot-access
import "bcts/internal/notify"

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

res, err := notify.Send(context.Background(), endpoint, 30, msg)
res.(*dingtalk.Response)
```
