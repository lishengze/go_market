package timeutils

import (
	"fmt"
	"github.com/shopspring/decimal"
	"time"
)

// TimeToSec 秒 10 位
func TimeToSec(t time.Time) int64 {
	return t.Unix()
}

// TimeToSecString 秒 10 位
func TimeToSecString(t time.Time) string {
	return fmt.Sprint(t.Unix())
}

// TimeToMillSec 毫秒 13 位
func TimeToMillSec(t time.Time) int64 {
	return decimal.NewFromInt(t.UnixNano() / 1e6).IntPart()
}

// TimeToMillSecString 毫秒 13 位
func TimeToMillSecString(t time.Time) string {
	return fmt.Sprint(t.UnixNano() / 1e6)
}

func MicroSecToTime(ts int64) time.Time {
	second := ts / 1000000
	microSecond := ts - second*1000000
	return time.Unix(second, microSecond*1e3)

}

func MillSecToTime(ts int64) time.Time {
	second := ts / 1000
	millSecond := ts - second*1000
	return time.Unix(second, millSecond*1e6)

}

func SecToTime(ts int64) time.Time {
	return time.Unix(ts, 0)
}

func MillSecStrToTime(tsStr string) (time.Time, error) {
	ts, err := decimal.NewFromString(tsStr)
	if err != nil {
		return time.Time{}, err
	}
	return MillSecToTime(ts.IntPart()), nil
}

func SecStrToTime(tsStr string) (time.Time, error) {
	ts, err := decimal.NewFromString(tsStr)
	if err != nil {
		return time.Time{}, err
	}
	return SecToTime(ts.IntPart()), nil
}

// TimeToExactMinute 转化成 整分钟  2022-04-22 17:48:16.802309 -> 2022-04-22 17:48:00
func TimeToExactMinute(t time.Time) time.Time {
	t = t.Add(-time.Nanosecond * time.Duration(t.Nanosecond()))
	t = t.Add(-time.Second * time.Duration(t.Second()))
	return t
}

//func TimeToRFC3339Mill(t time.Time) string {
//	s := t.Add(1).UTC().Format(time.RFC3339Nano) // 加一纳秒保证输出统一
//	if len(s) != 30 {
//		s = t.Add(2).UTC().Format(time.RFC3339Nano) // 加二纳秒
//		// 2019-05-17T15:07:08.941000001Z
//		if len(s) != 30 {
//			return "1970-01-01T00:00:00.000Z"
//		}
//	}
//	return s[:23] + s[29:]
//}
//
//func MillSecToRFC3339Mill(ts int64) string {
//	return TimeToRFC3339Mill(MillSecToTime(ts))
//}
