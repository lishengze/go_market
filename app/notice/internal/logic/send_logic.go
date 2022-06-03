package logic

import (
	"bcts/internal/notify"
	"bcts/internal/notify/email"
	"bcts/internal/notify/sms"
	"context"
	"market_server/app/notice/internal/config"
	"market_server/common/xerror"
	"strings"

	v1 "market_server/app/notice/api/v1"
	"market_server/app/notice/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendLogic {
	return &SendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendLogic) Send(in *v1.NoticeSendReq) (*v1.NoticeSendResp, error) {
	tm, ok := l.svcCtx.Config.TemplateConf.List[in.BizType.String()]
	if !ok || len(in.To) == 0 || in.RegisterType == v1.RegisterTypeOption_RegisterTypeUnknown {
		return nil, xerror.ErrorParamError
	}
	if len(in.To) > 100 {
		return nil, xerror.ErrNotifyTooManyRecipients
	}
	// 替换内容模板的变量
	if len(in.Amount) > 0 {
		tm.Content = strings.Replace(tm.Content, config.TEMPLATE_VARS_AMOUNT, in.Amount, -1)
	}

	if len(in.Currency) > 0 {
		tm.Content = strings.Replace(tm.Content, config.TEMPLATE_VARS_CURRENCY, in.Currency, -1)
	}

	if len(in.ChargeNo) > 0 {
		tm.Content = strings.Replace(tm.Content, config.TEMPLATE_VARS_CHARGENO, in.ChargeNo, -1)
	}

	switch in.RegisterType {
	case v1.RegisterTypeOption_RegisterWithMobile:
		return l.smsNotice(tm, in.To[1])
	case v1.RegisterTypeOption_RegisterWithEmail:
		return l.emailNotice(tm, in.To)
	}

	return nil, xerror.ErrorTryAgain
}

func (l *SendLogic) emailNotice(tm config.TemplateItem, to []string) (*v1.NoticeSendResp, error) {
	nofify := notify.GetNotify("email")
	// 构造邮件发送的数据
	data := &email.EmailApiPostData{
		Channel:  l.svcCtx.Config.EmailServerConf.Channel,
		Token:    l.svcCtx.Config.EmailServerConf.Token,
		From:     l.svcCtx.Config.EmailServerConf.From,
		FromName: l.svcCtx.Config.EmailServerConf.FromName,
		To:       to,
		Subject:  tm.Subject,
		Content:  strings.Replace(l.svcCtx.Config.TemplateConf.EmailHtmlTemplate, config.TEMPLATE_VARS_CONTENT, tm.Content, -1),
	}

	resp, err := nofify.Send(l.ctx, l.svcCtx.Config.EmailServerConf.Endpoint, l.svcCtx.Config.EmailServerConf.TimeOut, data)
	if err != nil {
		return nil, xerror.ErrNotifySendFail
	}

	out, ok := resp.(*email.EmailApiResponse)
	if !ok {
		return nil, xerror.ErrNotifyRespParseFail
	}

	if !out.Success && out.Error != nil {
		return nil, xerror.ErrNotifySendFail.WithMessage(out.Error.Code, out.Error.Message)
	}

	return &v1.NoticeSendResp{
		Success:     out.Success,
		Status:      out.Data.Status,
		DateCreated: out.Data.DateCreated,
	}, nil
}

func (l *SendLogic) smsNotice(tm config.TemplateItem, to string) (*v1.NoticeSendResp, error) {
	nofify := notify.GetNotify("sms")
	// 构造短信发送的数据
	data := &sms.SmsApiPostData{
		Channel: l.svcCtx.Config.SmsServerConf.Channel,
		Token:   l.svcCtx.Config.SmsServerConf.Token,
		To:      to,
		Content: tm.Content,
	}

	resp, err := nofify.Send(l.ctx, l.svcCtx.Config.SmsServerConf.Endpoint, l.svcCtx.Config.SmsServerConf.TimeOut, data)
	if err != nil {
		return nil, xerror.ErrNotifySendFail
	}

	out, ok := resp.(*sms.SmsApiResponse)
	if !ok {
		return nil, xerror.ErrNotifyRespParseFail
	}

	if !out.Success && out.Error != nil {
		return nil, xerror.NewCodeError(xerror.ErrNotifySendFail.ErrCode(), out.Error.Code, out.Error.Message)
	}

	return &v1.NoticeSendResp{
		Success:     out.Success,
		Status:      out.Data.Status,
		DateCreated: out.Data.DateCreated,
	}, nil
}
