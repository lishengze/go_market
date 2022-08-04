package util

import (
	"time"
)

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

func GetNextResolutionNanoTime(resolution int) int64 {
	return time.Now().UTC().UnixNano()
}

func UTCNanoTime() int64 {

	return time.Now().UTC().UnixNano()
}

func UTCMinuteNano() int64 {
	cur_time := UTCNanoTime()
	// fmt.Printf("ori_time: %s, %d \n", TimeStrFromInt(cur_time), cur_time)
	cur_time = cur_time - cur_time%int64(time.Minute)
	// fmt.Printf("trans_time: %s, %d \n", TimeStrFromInt(cur_time), cur_time)
	return cur_time
}

func LastUTCMinuteNano() int64 {
	cur_time := UTCMinuteNano()
	// fmt.Printf("ori_time: %s, %d \n", TimeStrFromInt(cur_time), cur_time)
	cur_time = cur_time - int64(time.Minute)
	// fmt.Printf("trans_time: %s, %d \n", TimeStrFromInt(cur_time), cur_time)
	return cur_time
}

func GetTimeFromtInt(int_time int64) time.Time {
	return time.Unix(int64(int_time/int64(time.Second)), int_time%int64(time.Second))
}

func TimeStrFromInt(int_time int64) string {
	dst_time := GetTimeFromtInt(int_time)
	return TimeToSecString(dst_time)
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

	time.Sleep(time.Duration(int64(time.Minute)-delta_nano_secs) * time.Nanosecond)
}

func IsNewMinuteStart(new_time int64, old_time int64) bool {

	new_time_minute := new_time - new_time%int64(time.Minute)
	old_time_minute := old_time - old_time%int64(time.Minute)

	return new_time_minute != old_time_minute
}

func IsNewResolutionStart(new_time int64, old_time int64, resolution int) bool {

	if resolution < int(time.Second) {
		resolution = resolution * int(time.Second)
	}

	new_time_minute := new_time - new_time%int64(resolution)
	old_time_minute := old_time - old_time%int64(resolution)

	return new_time_minute != old_time_minute
}

func TestNanoMinute() {
	UTCMinuteNano()
}
