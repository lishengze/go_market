package dingtalk

import "github.com/pkg/errors"

// doc https://open.dingtalk.com/document/robots/custom-robot-access
const (
	// 文本信息
	MSG_TYPE_TEXT = "text"

	// link类型
	MSG_TYPE_LINK = "link"

	// markdown类型
	MSG_TYPE_MARKDOWN = "markdown"

	// ActionCard类型
	MSG_TYPE_ACTION_CARD = "actionCard"

	// ActionCard single类型
	MSG_TYPE_ACTION_CARD_SINGLE = "actionCardSingle"

	// FeedCard类型
	MSG_TYPE_FEED_CARD = "feedCard"
)

type (
	MessageAt struct {
		AtMobiles []string `json:"atMobiles"`
		AtUserIds []string `json:"atUserIds"`
		IsAtAll   bool     `json:"isAtAll"`
	}
	// 文本信息
	TextMessage struct {
		Content string `json:"content"`
	}

	// link类型
	LinkMessage struct {
		Title      string `json:"title"`
		Text       string `json:"text"`
		MessageUrl string `json:"messageUrl"`
		PicUrl     string `json:"picUrl"`
	}

	// markdown类型
	MarkdownMessage struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	}

	// ActionCard Single 类型
	ActionCardSingleMessage struct {
		Title          string `json:"title"`
		Text           string `json:"text"`
		SingleTitle    string `json:"singleTitle"`
		SingleURL      string `json:"singleURL"`
		BtnOrientation string `json:"btnOrientation"`
	}

	ActionCardMessage struct {
		Title          string          `json:"title"`
		Text           string          `json:"text"`
		Btns           []ActionCardBtn `json:"btns"`
		BtnOrientation string          `json:"btnOrientation"`
	}
	ActionCardBtn struct {
		Title     string `json:"title"`
		ActionURL string `json:"actionURL"`
	}

	// FeedCard类型
	FeedCardMessage struct {
		Links []FeedCardLink `json:"links"`
	}

	FeedCardLink struct {
		Title      string `json:"title"`
		MessageUrl string `json:"messageUrl"`
		PicUrl     string `json:"picUrl"`
	}

	// 消息主体
	Message struct {
		MsgType          string                   `json:"msgtype,options=[text,link,markdown,actionCard,feedCard]"`
		At               *MessageAt               `json:"at,omitempty,optional"`
		Text             *TextMessage             `json:"text,omitempty,optional"`
		Link             *LinkMessage             `json:"link,omitempty,optional"`
		Markdown         *MarkdownMessage         `json:"markdown,omitempty,optional"`
		ActionCard       *ActionCardMessage       `json:"actionCard,omitempty,optional"`
		ActionCardSingle *ActionCardSingleMessage `json:"actionCard,omitempty,optional"`
		FeedCard         *FeedCardMessage         `json:"feedCard,omitempty,optional"`
	}

	Response struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
)

func NewMessageRequest(msgType string, message interface{}, at *MessageAt) (*Message, error) {
	switch msgType {
	case MSG_TYPE_TEXT:
		m, ok := message.(*TextMessage)
		if !ok {
			return nil, errors.Errorf("%v is not an object with type *TextMessage", message)
		}
		return &Message{
			MsgType: MSG_TYPE_TEXT,
			At:      at,
			Text:    m,
		}, nil
	case MSG_TYPE_LINK:
		m, ok := message.(*LinkMessage)
		if !ok {
			return nil, errors.Errorf("%v is not an object with type *LinkMessage", message)
		}
		return &Message{
			MsgType: MSG_TYPE_LINK,
			At:      at,
			Link:    m,
		}, nil
	case MSG_TYPE_MARKDOWN:
		m, ok := message.(*MarkdownMessage)
		if !ok {
			return nil, errors.Errorf("%v is not an object with type *MarkdownMessage", message)
		}
		return &Message{
			MsgType:  MSG_TYPE_MARKDOWN,
			At:       at,
			Markdown: m,
		}, nil

	case MSG_TYPE_ACTION_CARD:
		m, ok := message.(*ActionCardMessage)
		if !ok {
			return nil, errors.Errorf("%v is not an object with type *ActionCardMessage", message)
		}
		return &Message{
			MsgType:    MSG_TYPE_ACTION_CARD,
			At:         at,
			ActionCard: m,
		}, nil
	case MSG_TYPE_ACTION_CARD_SINGLE:
		m, ok := message.(*ActionCardSingleMessage)
		if !ok {
			return nil, errors.Errorf("%v is not an object with type *ActionCardSingleMessage", message)
		}
		return &Message{
			MsgType:          MSG_TYPE_ACTION_CARD,
			At:               at,
			ActionCardSingle: m,
		}, nil

	case MSG_TYPE_FEED_CARD:
		m, ok := message.(*FeedCardMessage)
		if !ok {
			return nil, errors.Errorf("%v is not an object with type *FeedCardMessage", message)
		}
		return &Message{
			MsgType:  MSG_TYPE_FEED_CARD,
			At:       at,
			FeedCard: m,
		}, nil
	default:
		return nil, errors.New("invalid msgType")
	}
}
