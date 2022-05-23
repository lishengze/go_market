package main

import (
	"fmt"
	"market_aggregate/app/aggregate"
)

func main() {
	fmt.Println("Test Risk Ctrl")

	// aggregate.TestInnerDepth()

	// aggregate.TestImport()

	// aggregate.TestWorker()

	// TestDepthChannel()

	// TestTreeMap()

	// comm.TestSeKline()

	aggregate.TestServerEngine()
}
