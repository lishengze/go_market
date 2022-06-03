package logic

import (
	"bcts/internal/notify"
	"context"
	"market_server/common/xerror"
	"strconv"

	v1 "market_server/app/notice/api/v1"
	"market_server/app/notice/internal/svc"

	"bcts/internal/notify/dingtalk"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlertLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlertLogic {
	return &AlertLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AlertLogic) Alert(in *v1.AlertSendReq) (*v1.AlertSendResp, error) {
	robotConf, ok := l.svcCtx.Config.DingTalkRobotConf[in.Group.String()]
	if !ok {
		return nil, xerror.ErrorParamError.WithMessage("undefined robot conf", "机器人配置不存在")
	}

	m, msgtype, err := l.buildMessage(in)
	if err != nil {
		return nil, err
	}

	at := &dingtalk.MessageAt{
		AtMobiles: in.At,
		IsAtAll:   in.IsAtAll,
	}

	msg, err := dingtalk.NewMessageRequest(msgtype, m, at)
	if err != nil {
		return nil, err
	}

	endpoint, err := dingtalk.GetRobotSendUrl(robotConf.AccessToken, robotConf.Secret)
	if err != nil {
		return nil, err
	}

	notify := notify.GetNotify("ding_talk")
	res, err := notify.Send(context.Background(), endpoint, 30, msg)
	if err != nil {
		return nil, err
	}

	return &v1.AlertSendResp{
		Errcode: int32(res.(*dingtalk.Response).ErrCode),
		Errmsg:  res.(*dingtalk.Response).ErrMsg,
	}, nil
}

func (l *AlertLogic) buildMessage(in *v1.AlertSendReq) (interface{}, string, error) {
	switch in.MsgType {
	case v1.AlertMsgTypeOPtions_MsgTypeText:
		if len(in.Text.Content) == 0 {
			return nil, dingtalk.MSG_TYPE_TEXT, xerror.ErrorParamError
		}
		return dingtalk.TextMessage{
			Content: in.Text.Content,
		}, dingtalk.MSG_TYPE_TEXT, nil
	case v1.AlertMsgTypeOPtions_MsgTypeLink:
		if len(in.Link.Title) == 0 || len(in.Link.Text) == 0 || len(in.Link.MessageUrl) == 0 {
			return nil, dingtalk.MSG_TYPE_LINK, xerror.ErrorParamError
		}
		return dingtalk.LinkMessage{
			Title:      in.Link.Title,
			Text:       in.Link.Text,
			MessageUrl: in.Link.MessageUrl,
			PicUrl:     in.Link.PicUrl,
		}, dingtalk.MSG_TYPE_LINK, nil
	case v1.AlertMsgTypeOPtions_MsgTypeMarkdown:
		if len(in.Markdown.Title) == 0 || len(in.Markdown.Text) == 0 {
			return nil, dingtalk.MSG_TYPE_MARKDOWN, xerror.ErrorParamError
		}
		return dingtalk.MarkdownMessage{
			Title: in.Markdown.Title,
			Text:  in.Markdown.Text,
		}, dingtalk.MSG_TYPE_MARKDOWN, nil
	case v1.AlertMsgTypeOPtions_MsgTypeActionCardSingle:
		if len(in.ActionCardSingle.Title) == 0 ||
			len(in.ActionCardSingle.Text) == 0 ||
			len(in.ActionCardSingle.SingleTitle) == 0 ||
			len(in.ActionCardSingle.SingleUrl) == 0 {
			return nil, dingtalk.MSG_TYPE_ACTION_CARD_SINGLE, xerror.ErrorParamError
		}
		return dingtalk.ActionCardSingleMessage{
			Title:          in.ActionCardSingle.Title,
			Text:           in.ActionCardSingle.Text,
			SingleTitle:    in.ActionCardSingle.SingleTitle,
			SingleURL:      in.ActionCardSingle.SingleUrl,
			BtnOrientation: strconv.Itoa(int(in.ActionCardSingle.BtnOrientation)),
		}, dingtalk.MSG_TYPE_ACTION_CARD_SINGLE, nil
	case v1.AlertMsgTypeOPtions_MsgTypeActionCard:
		if len(in.ActionCard.Title) == 0 ||
			len(in.ActionCard.Text) == 0 ||
			len(in.ActionCard.Btns) == 0 {
			return nil, dingtalk.MSG_TYPE_ACTION_CARD, xerror.ErrorParamError
		}
		btns := make([]dingtalk.ActionCardBtn, 0, len(in.ActionCard.Btns))
		for _, btn := range in.ActionCard.Btns {
			btns = append(btns, dingtalk.ActionCardBtn{
				Title:     btn.Title,
				ActionURL: btn.ActionUrl,
			})
		}

		return dingtalk.ActionCardMessage{
			Title:          in.ActionCard.Title,
			Text:           in.ActionCard.Text,
			Btns:           btns,
			BtnOrientation: strconv.Itoa(int(in.ActionCard.BtnOrientation)),
		}, dingtalk.MSG_TYPE_ACTION_CARD, nil
	case v1.AlertMsgTypeOPtions_MsgTypeFeedCard:
		if len(in.FeedCard.Links) == 0 {
			return nil, dingtalk.MSG_TYPE_FEED_CARD, xerror.ErrorParamError
		}
		links := make([]dingtalk.FeedCardLink, 0, len(in.FeedCard.Links))
		for _, link := range in.FeedCard.Links {
			links = append(links, dingtalk.FeedCardLink{
				Title:      link.Title,
				MessageUrl: link.MessageUrl,
				PicUrl:     link.PicUrl,
			})
		}
	}

	return nil, dingtalk.MSG_TYPE_FEED_CARD, xerror.ErrorParamError
}
