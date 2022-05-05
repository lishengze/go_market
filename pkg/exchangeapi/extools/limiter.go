package extools

import (
	"exterior-interactor/pkg/httptools"
)

type LimiterManager interface {
	ThroughLimiters(request httptools.Request, meta Meta) error
	//RejectAllRequest(duration time.Time) // 交易所发出警告时（收到 http 418 429 码时） 调用此方法
}
