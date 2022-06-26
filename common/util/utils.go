package util

import (
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

func LOG_INFO(info string) {
	fmt.Println("INFO: " + info)
}

func LOG_WARN(info string) {
	fmt.Println("WARN: " + info)
}

func LOG_ERROR(info string) {
	fmt.Println("Error " + info)
}

func ExceptionFunc() {
	errMsg := recover()
	if errMsg != nil {
		fmt.Println(errMsg)
	}
}

func CurrTimeString() string {
	return TimeToSecString(time.Now())
}

func TimeFormat() string {
	return "2006-01-02 15:04:05"
}

// TimeToSecString 秒 10 位
func TimeToSecString(t time.Time) string {
	// fmt.Println(t.GoString())
	// t.GoString()
	return t.Format("2006-01-02 15:04:05")
}

func UTCNanoTime() int64 {

	return time.Now().UTC().UnixNano()
}

func TimeMinute() int64 {
	utc_time_secs := time.Now().Unix()

	utc_time_min_secs := TimeToExactMinute(time.Unix(utc_time_secs, 0)).Unix()

	return utc_time_min_secs
}

func TimeMinuteNanos() int64 {
	utc_time_secs := time.Now().Unix()

	utc_time_min_nanos := TimeToExactMinute(time.Unix(utc_time_secs, 0)).UTC().UnixNano()

	return utc_time_min_nanos
}

func TimeToExactMinute(t time.Time) time.Time {
	t = t.Add(-time.Nanosecond * time.Duration(t.Nanosecond()))
	t = t.Add(-time.Second * time.Duration(t.Second()))
	return t
}

func WaitForNextMinute() {
	utc_time_nano_secs := time.Now().UnixNano()

	utc_time_min_nano_secs := TimeToExactMinute(time.Unix(utc_time_nano_secs/1e9, utc_time_nano_secs%1e9)).UnixNano()

	delta_nano_secs := utc_time_nano_secs - utc_time_min_nano_secs

	// fmt.Printf("\nutc_time_secs: %d, utc_time_min_secs: %d, delta_secs: %d\n",
	// 	utc_time_secs, utc_time_min_secs, delta_secs)
	// fmt.Println(time.Unix(utc_time_secs, 0))
	// fmt.Println(time.Unix(utc_time_min_secs, 0))

	time.Sleep(time.Duration(1e9*60-delta_nano_secs) * time.Nanosecond)
}

func TestTimeStr() {
	// fmt.Println(TimeString())
}

func MinFloat64(a float64, b float64) float64 {
	if a < b {
		return a
	} else {
		return b
	}
}
func MaxFloat64(a float64, b float64) float64 {
	if a > b {
		return a
	} else {
		return b
	}
}

// func NanoTimeUInt64() uint64 {

// }

func TestUTCMinuteNano() {
	t := TimeMinuteNanos()
	t2 := time.Unix(int64(t/int64(time.Second)), t%int64(time.Second))

	fmt.Printf("IntTime: %d, Time: %+v\n", t, t2)

	WaitForNextMinute()

	tt := UTCNanoTime()
	tt2 := time.Unix(int64(tt/int64(time.Second)), tt%int64(time.Second))

	fmt.Printf("IntTime: %d, Time: %+v\n", tt, tt2)
}

func Float64ComparatorDsc(a, b interface{}) int {
	aAsserted := a.(float64)
	bAsserted := b.(float64)
	switch {
	case aAsserted < bAsserted:
		return 1
	case aAsserted > bAsserted:
		return -1
	default:
		return 0
	}
}

func InitTestLogx() {

	LogConfig := logx.LogConf{
		Compress:            true,
		KeepDays:            0,
		Level:               "info",
		Mode:                "file",
		Path:                "./log",
		ServiceName:         "client",
		StackCooldownMillis: 100,
		TimeFormat:          "2006-01-02 15:04:05",
	}

	logx.MustSetup(LogConfig)
}
