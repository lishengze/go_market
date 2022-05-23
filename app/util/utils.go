package util

import (
	"fmt"
	"time"
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

func TimeToExactMinute(t time.Time) time.Time {
	t = t.Add(-time.Nanosecond * time.Duration(t.Nanosecond()))
	t = t.Add(-time.Second * time.Duration(t.Second()))
	return t
}

func WaitForNextMinute() {
	utc_time_secs := time.Now().Unix()

	utc_time_min_secs := TimeToExactMinute(time.Unix(utc_time_secs, 0)).Unix()

	delta_secs := utc_time_secs - utc_time_min_secs

	// fmt.Printf("\nutc_time_secs: %d, utc_time_min_secs: %d, delta_secs: %d\n",
	// 	utc_time_secs, utc_time_min_secs, delta_secs)
	// fmt.Println(time.Unix(utc_time_secs, 0))
	// fmt.Println(time.Unix(utc_time_min_secs, 0))

	time.Sleep(time.Duration(60-delta_secs) * time.Second)
}

func TestTimeStr() {
	// fmt.Println(TimeString())
}

// func NanoTimeUInt64() uint64 {

// }
