package main

import (
	"flag"
	"fmt"
	"market_server/app/market_aggregate/aggregate"
	"market_server/app/market_aggregate/config"
	"market_server/app/market_aggregate/svc"
	"os"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

func main() {
	fmt.Printf("Args: %+v \n", os.Args)
	env := os.Args[1]

	fmt.Printf("env: %+v \n", env)
	var configFile = flag.String("f", "etc/"+env+"/client.yaml", "the config file")

	var c config.Config
	conf.MustLoad(*configFile, &c)

	logx.MustSetup(c.LogConfig)

	fmt.Printf("Log: %+v \n", c)

	ctx := svc.NewServiceContext(c)
	svr := aggregate.NewServerEngine(ctx)

	svr.Start()
}
