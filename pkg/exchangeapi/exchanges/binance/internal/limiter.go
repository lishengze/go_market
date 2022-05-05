package internal

import (
	"context"
	"exterior-interactor/pkg/exchangeapi/extools"
	"exterior-interactor/pkg/httptools"
	"fmt"
	"golang.org/x/time/rate"
	"strings"
	"time"
)

var DefaultLimiterManager = newDefaultLimiterManager()

type (
	limiterManager struct {
		apiLimiter  *rate.Limiter // api
		sapiLimiter *rate.Limiter // sapi
		fapiLimiter *rate.Limiter // fapi
		dapiLimiter *rate.Limiter // dapi
		vapiLimiter *rate.Limiter // vapi
	}
)

func newDefaultLimiterManager() extools.LimiterManager {
	var defaultBucketCap = 600

	lm := &limiterManager{
		apiLimiter:  rate.NewLimiter(18, defaultBucketCap),
		sapiLimiter: rate.NewLimiter(18, defaultBucketCap),
		fapiLimiter: rate.NewLimiter(18, defaultBucketCap),
		dapiLimiter: rate.NewLimiter(18, defaultBucketCap),
		vapiLimiter: rate.NewLimiter(18, defaultBucketCap),
	}
	// 清空 token
	lm.apiLimiter.AllowN(time.Now(), defaultBucketCap)
	lm.sapiLimiter.AllowN(time.Now(), defaultBucketCap)
	lm.fapiLimiter.AllowN(time.Now(), defaultBucketCap)
	lm.dapiLimiter.AllowN(time.Now(), defaultBucketCap)
	lm.vapiLimiter.AllowN(time.Now(), defaultBucketCap)

	return lm
}

func (o *limiterManager) ThroughLimiters(request httptools.Request, meta extools.Meta) error {
	if strings.HasPrefix(meta.Url(), "https://api.binance.com/api") {
		return o.through(o.apiLimiter, meta)
	} else if strings.HasPrefix(meta.Url(), "https://api.binance.com/sapi") {
		return o.through(o.sapiLimiter, meta)
	} else if strings.HasPrefix(meta.Url(), "https://fapi.binance.com/fapi") {
		return o.through(o.fapiLimiter, meta)
	} else if strings.HasPrefix(meta.Url(), "https://dapi.binance.com/dapi") {
		return o.through(o.dapiLimiter, meta)
	} else if strings.HasPrefix(meta.Url(), "https://vapi.binance.com/vapi") {
		return o.through(o.vapiLimiter, meta)
	} else {
		panic(fmt.Sprintf("can not parse ApiType, url:%s", meta.Url()))
	}
}

func (o *limiterManager) through(limiter *rate.Limiter, meta extools.Meta) error {
	switch meta.ReqType() {
	case extools.ReqTypeWait:
		return limiter.WaitN(context.Background(), meta.Weight())
	case extools.ReqTypeAllow:
		if !limiter.AllowN(time.Now(), meta.Weight()) {
			return fmt.Errorf("request exceeds the frequency limit, discarded. ")
		} else {
			return nil
		}
	default:
		return fmt.Errorf("wrong reqType:%s. ", meta.ReqType().String())
	}
}
