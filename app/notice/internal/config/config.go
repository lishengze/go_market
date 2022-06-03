package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

const (
	TEMPLATE_VARS_AMOUNT   = "{{AMOUNT}}"
	TEMPLATE_VARS_CURRENCY = "{{CURRENCY}}"
	TEMPLATE_VARS_CONTENT  = "{{CONTEXT}}"
	TEMPLATE_VARS_CHARGENO = "{{CHARGENO}}"
)

type Config struct {
	zrpc.RpcServerConf
	EmailServerConf   EmailServerConf
	SmsServerConf     SmsServerConf
	TemplateConf      TemplateConf
	DingTalkRobotConf DingTalkRobotConf
}

type EmailServerConf struct {
	Endpoint string
	From     string
	FromName string
	Channel  string
	Token    string
	// setting 0 means no timeout
	TimeOut int64 `json:",default=2000"`
}

type SmsServerConf struct {
	Endpoint string
	Channel  string
	Token    string
	TimeOut  int64 `json:",default=2000"`
}

type TemplateConf struct {
	List              map[string]TemplateItem
	EmailHtmlTemplate string
}

type TemplateItem struct {
	Subject string
	Content string
}

// 订单群助手机器人配置
type DingTalkRobotConf map[string]RobotConf

type RobotConf struct {
	AccessToken string
	Secret      string
}
