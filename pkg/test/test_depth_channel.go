package main

import (
	"fmt"
	"market_aggregate/pkg/riskctrl"
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

func main() {
	fmt.Println("Test Risk Ctrl")

	// riskctrl.TestInnerDepth()

	// riskctrl.TestImport()

	// riskctrl.TestWorker()

	TestDepthChannel()

	// TestTreeMap()
}
