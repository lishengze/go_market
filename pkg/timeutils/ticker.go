package timeutils

import (
	"time"
)

func Every(d time.Duration, fn func()) {
	ticker := time.NewTicker(d)
	go func() {
		for range ticker.C {
			fn()
		}
	}()
}
