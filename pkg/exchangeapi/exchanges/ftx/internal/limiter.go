package internal

import (
	"context"
	"exterior-interactor/pkg/exchangeapi/extools"
	"exterior-interactor/pkg/httptools"
	"fmt"
	"golang.org/x/time/rate"
	"time"
)

var DefaultLimiterManager = newDefaultLimiterManager()

type (
	limiterManager struct {
		*rate.Limiter
	}
)

func newDefaultLimiterManager() extools.LimiterManager {
	/*
		Rate limits
		Please do not send more than 30 requests per second: doing so will result in HTTP 429 errors.
		We strongly recommend using the websocket API for faster market and account data.
	*/

	lm := &limiterManager{
		Limiter: rate.NewLimiter(29, 1),
	}
	// 清空 token
	lm.Limiter.AllowN(time.Now(), 2)

	return lm
}

func (o *limiterManager) ThroughLimiters(request httptools.Request, meta extools.Meta) error {
	switch meta.ReqType() {
	case extools.ReqTypeWait:
		return o.WaitN(context.Background(), meta.Weight())
	case extools.ReqTypeAllow:
		if !o.AllowN(time.Now(), meta.Weight()) {
			return fmt.Errorf("request exceeds the frequency limit, discarded. ")
		} else {
			return nil
		}
	default:
		return fmt.Errorf("wrong reqType:%s. ", meta.ReqType().String())
	}
}
