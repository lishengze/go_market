package main

import (
	"fmt"
	"market_server/common/util"
)

func test_nano() {
	int_time := util.UTCNanoTime()

	fmt.Printf("int_time: %d, trans_time: %+v \n", int_time, util.GetTimeFromtInt(int_time))
}
func main() {
	// util.TestUTCMinuteNano()

	util.TestNanoMinute()

	// test_nano()
}
