package main

import (
	"fmt"
	"market_aggregate/pkg/riskctrl"
	"reflect"
	"time"
)

// 测试用定时器向一个 channel 发送 depth_quote
func send_depth(channel_depth_quote chan *riskctrl.DepthQuote) {

	duration := time.Duration(3 * time.Second)
	timer := time.Tick(duration)
	for {
		select {
		case <-timer:
			// fmt.Println(time.Now())
			depth_quote := riskctrl.GetTestDepth()
			fmt.Printf("\nSendDepth: %s\n", depth_quote.String(5))
			channel_depth_quote <- &depth_quote
		}
	}
}

func recev_depth(channel_depth_quote chan *riskctrl.DepthQuote) {
	for {
		select {
		case depth_quote := <-channel_depth_quote:
			fmt.Printf("\nRecvDepth: %s\n", depth_quote.String(5))
			// fmt.Println(depth_quote.String(5))
		}
	}
}

func TestDepthChannel() {
	channel_depth_quote := make(chan *riskctrl.DepthQuote)

	go send_depth(channel_depth_quote)

	go recev_depth(channel_depth_quote)

	time.Sleep(time.Hour)

}

type EmptyInterface interface {
}

func process_depth(depth_quote *riskctrl.DepthQuote) {
	fmt.Println(depth_quote.String(5))
}

func TestDepthReflection(data EmptyInterface) {

	// data_value := reflect.ValueOf(data)

	// fmt.Println(data_value)

	// fmt.Println(reflect.TypeOf(data))

	// data_type

	// if reflect.TypeOf(data) == riskctrl.DepthQuote {

	// }

	fmt.Println(reflect.TypeOf(data).Name())

	if reflect.TypeOf(data).Name() == "DepthQuote" {
		// data_value := reflect.ValueOf(data)
		// process_depth(data)
	}

	// reflect.ValueOf(data).(riskctrl.DepthQuote)
}

func TestReflection() {
	depth_quote := riskctrl.GetTestDepth()

	TestDepthReflection(depth_quote)
}

func main() {
	fmt.Println("Test Risk Ctrl")

	// riskctrl.TestInnerDepth()

	// riskctrl.TestImport()

	// riskctrl.TestWorker()

	// TestDepthChannel()

	TestReflection()

	// TestTreeMap()
}
