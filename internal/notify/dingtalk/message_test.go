package dingtalk

import (
	"reflect"
	"testing"

	"github.com/zeromicro/go-zero/core/mapping"
)

func TestNewMessageRequest(t *testing.T) {
	m := &TextMessage{
		Content: "test",
	}
	msg, err := NewMessageRequest(MSG_TYPE_TEXT, m, &MessageAt{
		AtMobiles: []string{"15221560954"},
	})

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(msg.Text, m) {
		t.Errorf("expected message to be %+v, got %+v ", m, msg.Text)
	}
	val, err := mapping.Marshal(msg)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", val)
}
