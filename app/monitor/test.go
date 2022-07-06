package main

import (
	"flag"
	"fmt"
	"market_server/app/monitor/monitor_market"
	"os"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

func main() {
	flag.Parse()
	env := os.Args[1]

	fmt.Printf("env: %+v \n", env)
	var configFile = flag.String("f", "etc/"+env+"/config.yaml", "the config file")

	fmt.Println(*configFile)

	var c monitor_market.Config
	conf.MustLoad(*configFile, &c)

	fmt.Printf("config: %+v \n", c)

	logx.MustSetup(c.LogConfig)

	logx.Infof("config: %+v \n", c)

	server := monitor_market.NewServerEngine(&c)
	server.Start()

	select {}
}
