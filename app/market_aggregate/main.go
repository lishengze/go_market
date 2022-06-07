package main

import (
	"flag"
	"fmt"
	"market_server/app/market_aggregate/aggregate"
	"market_server/app/market_aggregate/config"
	"market_server/app/market_aggregate/svc"
	"os"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

func set_test_meta(s *aggregate.ServerEngine) {
	s.SetTestFlag(true)

	s.SetTestConfig()
}

func main() {
	fmt.Printf("Args: %+v \n", os.Args)
	env := os.Args[1]

	fmt.Printf("env: %+v \n", env)
	var configFile = flag.String("f", "etc/"+env+"/client.yaml", "the config file")

	var c config.Config
	conf.MustLoad(*configFile, &c)

	logx.MustSetup(c.LogConfig)

	fmt.Printf("Log: %s \n", c.String())

	ctx := svc.NewServiceContext(c)
	svr := aggregate.NewServerEngine(ctx)

	is_test := true

	if is_test {
		set_test_meta(svr)
	}

	svr.Start()

	if is_test {
		time.Sleep(time.Second * 6)
		svr.TestKafkaCancelListen()
	}

	select {}
}
