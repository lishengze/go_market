package notify

import (
	"context"
	"strings"
)

type (
	Notify interface {
		Send(ctx context.Context, endpoint string, timeout int64, data interface{}) (interface{}, error)
		// Name returns the name of the Notify implementation. The returned string
		Name() string
	}
)

var registeredNotifies = make(map[string]Notify)

// RegisterNotify registers the provided Notify for use with all Transport clients and
// servers.
func RegisterNotify(notify Notify) {
	if notify == nil {
		panic("cannot register a nil Notify")
	}
	if notify.Name() == "" {
		panic("cannot register Notify with empty string result for Name()")
	}
	notifyType := strings.ToLower(notify.Name())
	registeredNotifies[notifyType] = notify
}

// GetNotify gets a registered Notify by notify-type, or nil if no Notify is
// registered for the notify-type.
// The notify-type is expected to be lowercase.
func GetNotify(notifyType string) Notify {
	return registeredNotifies[notifyType]
}
