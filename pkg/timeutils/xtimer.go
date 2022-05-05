package timeutils

import (
	"context"
	"time"
)

type (
	// XTimer 定时执行
	XTimer interface {
		OnTimer(fn func())
		Close()
	}

	xTimer struct {
		cancel   context.CancelFunc
		ctx      context.Context
		fn       func()
		ticker   *time.Ticker
		interval time.Duration
	}
)

func NewXTimer(interval time.Duration, fn func()) XTimer {
	if fn == nil {
		fn = func() {}
	}
	ctx, cancel := context.WithCancel(context.Background())
	ifs := &xTimer{
		cancel: cancel,
		ctx:    ctx,
		fn:     fn,
		ticker: time.NewTicker(interval),
	}
	go ifs.run()
	return ifs
}

func (o *xTimer) run() {
	for {
		select {
		case <-o.ticker.C:
			o.fn()
		case <-o.ctx.Done():
			return
		}
	}
}

func (o *xTimer) OnTimer(fn func()) {
	if fn != nil {
		o.fn = fn
	}
}

func (o *xTimer) Close() {
	defer o.ticker.Stop()
	o.cancel()
}
