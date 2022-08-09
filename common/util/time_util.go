package util

import (
	"fmt"
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

func SetToStartTime(src_time int64, resolution uint64) int64 {

	if resolution < uint64(time.Second) {
		resolution = resolution * uint64(time.Second)
	}

	nano_per_day := uint64(24 * time.Hour)
	nano_per_week := uint64(7 * nano_per_day)

	new_src_time := src_time

	if resolution == nano_per_week {
		src_time = src_time + int64(3*nano_per_day)

		// left_nanos := src_time % int64(resolution)
		// left_secs := left_nanos / int64(time.Second)
		// left_days := left_secs / (24 * 3600)

		// fmt.Printf("FullTime: %s, left_secs: %d, left_day: %d\n", TimeStrFromInt(src_time), left_secs, left_days)
	}

	start_time := new_src_time - src_time%int64(resolution)

	return start_time
}

func IsNewResolutionStart(new_time int64, old_time int64, resolution uint64) bool {

	if resolution < uint64(time.Second) {
		resolution = resolution * uint64(time.Second)
	}

	new_start_time := SetToStartTime(new_time, resolution)
	old_start_time := SetToStartTime(old_time, resolution)

	return new_start_time != old_start_time
}

func TestNanoMinute() {
	UTCMinuteNano()
}

func TestSetToStartTime() {
	resolution := uint64(7 * 24 * time.Hour)
	time1 := time.Date(1970, 1, 7, 10, 20, 33, 0, time.UTC)
	int_time1 := time1.UTC().UnixNano()
	start_time1 := SetToStartTime(int_time1, resolution)

	time2 := time.Date(2022, 8, 3, 10, 20, 33, 0, time.UTC)
	int_time2 := time2.UTC().UnixNano()
	start_time2 := SetToStartTime(int_time2, resolution)

	fmt.Printf("time1: %s, start_time1: %s\ntime2: %s, start_time2: %s\n",
		TimeStrFromInt(int_time1), TimeStrFromInt(start_time1),
		TimeStrFromInt(int_time2), TimeStrFromInt(start_time2))

	resolution2 := uint64(24 * time.Hour)
	time3 := time.Date(1970, 1, 7, 10, 20, 33, 0, time.UTC)
	int_time3 := time3.UTC().UnixNano()
	start_time3 := SetToStartTime(int_time3, resolution2)

	time4 := time.Date(2022, 8, 3, 10, 20, 33, 0, time.UTC)
	int_time4 := time4.UTC().UnixNano()
	start_time4 := SetToStartTime(int_time4, resolution2)

	fmt.Printf("time3: %s, start_time1: %s\ntime2: %s, start_time2: %s\n",
		TimeStrFromInt(int_time3), TimeStrFromInt(start_time3),
		TimeStrFromInt(int_time4), TimeStrFromInt(start_time4))
}
