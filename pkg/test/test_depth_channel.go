package main

import (
	"fmt"
	"market_aggregate/pkg/datastruct"
	"market_aggregate/pkg/util"
	"reflect"
	"time"
)

// 测试用定时器向一个 channel 发送 depth_quote
func send_depth(channel_depth_quote chan *datastruct.DepthQuote) {

	duration := time.Duration(3 * time.Second)
	timer := time.Tick(duration)
	for {
		select {
		case <-timer:
			// fmt.Println(time.Now())
			depth_quote := datastruct.GetTestDepth()
			fmt.Printf("\nSendDepth: %s\n", depth_quote.String(5))
			channel_depth_quote <- depth_quote
		}
	}
}

func recev_depth(channel_depth_quote chan *datastruct.DepthQuote) {
	for {
		select {
		case depth_quote := <-channel_depth_quote:
			fmt.Printf("\nRecvDepth: %s\n", depth_quote.String(5))
			// fmt.Println(depth_quote.String(5))
		}
	}
}

func TestDepthChannel() {
	channel_depth_quote := make(chan *datastruct.DepthQuote)

	go send_depth(channel_depth_quote)

	go recev_depth(channel_depth_quote)

	time.Sleep(time.Hour)

}

type EmptyInterface interface {
}

func process_depth(depth_quote *datastruct.DepthQuote) {
	fmt.Println(depth_quote.String(5))
}

func TestDepthReflection(data EmptyInterface) {

	// data_value := reflect.ValueOf(data)

	// fmt.Println(data_value)

	// fmt.Println(reflect.TypeOf(data))

	// data_type

	// if reflect.TypeOf(data) == datastruct.DepthQuote {

	// }

	fmt.Println(reflect.TypeOf(data).Name())

	if reflect.TypeOf(data).Name() == "DepthQuote" {
		// data_value := reflect.ValueOf(data)
		// process_depth(data)
	}

	// reflect.ValueOf(data).(datastruct.DepthQuote)
}

func TestReflection() {
	depth_quote := datastruct.GetTestDepth()

	TestDepthReflection(depth_quote)
}

func TestInnerDepth() {
	new_depth := datastruct.GetTestDepth()

	if result, ok := new_depth.Asks.Get(41001.11111); ok {
		// fmt.Println()
		trans := result.(*datastruct.InnerDepth)

		// trans.
		fmt.Printf("Original Depth:%+v \n", trans)

		trans.Volume += 100

		fmt.Printf("After Add Depth:%+v \n", trans)

	}

	fmt.Println(new_depth)
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

	fmt.Printf("\nutc_time_secs: %d, utc_time_min_secs: %d, delta_secs: %d\n",
		utc_time_secs, utc_time_min_secs, delta_secs)
	fmt.Println(time.Unix(utc_time_secs, 0))
	fmt.Println(time.Unix(utc_time_min_secs, 0))

	time.Sleep(time.Duration(60-delta_secs) * time.Second)
}

func TestTime() {
	// utc_time_secs := time.Now().Unix()

	// utc_time_min_secs := TimeToExactMinute(time.Unix(utc_time_secs, 0)).Unix()

	// delta_secs := utc_time_secs - utc_time_min_secs

	// fmt.Println(utc_time_secs)

	// fmt.Println(utc_time_secs)

	// fmt.Printf("utc_time_secs: %d, utc_time_min_secs: %d, delta_secs: %d\n",
	// 	utc_time_secs, utc_time_min_secs, delta_secs)

	// fmt.Println(time.Unix(utc_time_secs, 0))
	// fmt.Println(time.Unix(utc_time_min_secs, 0))

	for {
		util.WaitForNextMinute()

		fmt.Println(time.Now())
	}

	// time := time.Unix(int_time, 0)

	// fmt.Println(int_time)

	// fmt.Println(time)
}

// func main() {
// 	// fmt.Println("Test Risk Ctrl")

// 	// aggregate.TestInnerDepth()

// 	// aggregate.TestImport()

// 	// aggregate.TestWorker()

// 	// TestDepthChannel()

// 	// TestReflection()

// 	// TestInnerDepth()

// 	// TestTreeMap()

// 	// aggregate.TestAggregator()

// 	// kafkaClient.TestConsumer()

// 	// TestTime()

// 	// comm.TestSeDepth()

// 	// comm.TestSeTrade()

// }
