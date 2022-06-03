package utils

import "time"

/**
传入的日期早于当天，返回true
否则返回false
*/
func CompareDate(day string) bool {
	//获取当前日期
	today := time.Now().Format("2006-01-02")
	t1, err1 := time.Parse("2006-01-02", day)
	t2, err2 := time.Parse("2006-01-02", today)

	if err1 == nil && err2 == nil && t1.Before(t2) {
		return true
	} else {
		return false
	}
}
