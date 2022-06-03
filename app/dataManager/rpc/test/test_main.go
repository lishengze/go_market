package main

import (
	"fmt"
	"market_server/app/dataManager/rpc/internal/dbserver"
)

func main() {
	fmt.Println("------- Test Main --------")

	dbserver.TestDB()
}
