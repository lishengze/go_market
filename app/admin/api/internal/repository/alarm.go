package repository

import (
	"context"
	"fmt"
	"market_server/app/admin/api/internal/svc"
	"market_server/common/dingtalk"
	"market_server/common/nacosAdapter"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmRepository struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAlarmRepository(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmRepository {
	return &AlarmRepository{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

var (
	BrockConf *nacosAdapter.BrokerConf
	ClientX   *nacosAdapter.Client
)

func (l *AlarmRepository) DingAlertMsg(alertContent string) (bool, error) {
	l.ConProcessLoad()
	if l.svcCtx.Config.Environment.Env != "" {
		alertContent = strings.Replace(alertContent, "env", l.svcCtx.Config.Environment.Env, 1)
	}
	if l.svcCtx.Config.Environment.Env == "local" {
		return true, nil
	}
	dingClient := dingtalk.NewClient(BrockConf.DingDingAskTalk.AccessToken, BrockConf.DingDingAskTalk.Secret)
	_, err := dingClient.SendMessage(alertContent)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return true, nil
}

func (l *AlarmRepository) DingPmsAlertMsg(alertContent string) (bool, error) {
	l.ConProcessLoad()
	if l.svcCtx.Config.Environment.Env != "" {
		alertContent = alertContent + "     env: " + l.svcCtx.Config.Environment.Env
	}
	dingClient := dingtalk.NewClient(BrockConf.DingDingPms.AccessToken, BrockConf.DingDingPms.Secret)
	_, err := dingClient.SendMessage(alertContent)
	if err != nil {
		fmt.Printf("DingPmsAlertMsg: error(%+v)", err)
		return false, err
	}
	return true, nil
}

func (l *AlarmRepository) ConProcessLoad() {
	if BrockConf == nil {
		BrockConf = l.GetConfNaAdapterMs()
	}
}

func (l *AlarmRepository) GetConfNaAdapterMs() *nacosAdapter.BrokerConf {
	parameters := l.svcCtx.Parameters
	brockConf, err := parameters.GetBrokeConf()
	if err != nil {
		fmt.Printf("DingAlertMsg: get config err(%+v)", err)
	}
	return brockConf
}
