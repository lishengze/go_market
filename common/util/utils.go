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

func CatchExp(func_name string) {
	errMsg := recover()
	if errMsg != nil {
		logx.Errorf("%s errMsg: %+v \n", func_name, errMsg)
		logx.Infof("%s errMsg: %+v \n", func_name, errMsg)
	}
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

func MaxInt(a int, b int) int {
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

// func GetFloorTime(ori_time int64, resolution_secs int64) int64 {

// }

func InitTestLogx() {

	LogConfig := logx.LogConf{
		Compress:            true,
		Encoding:            "plain",
		KeepDays:            0,
		Level:               "info",
		Mode:                "file",
		Path:                "./log",
		ServiceName:         "client",
		StackCooldownMillis: 100,
		TimeFormat:          "2006-01-02 15:04:05",
	}

	logx.MustSetup(LogConfig)
	logx.Infof("----- InitTestLogx -----\n")
}
