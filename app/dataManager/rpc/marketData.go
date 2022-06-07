package main

import (
	"flag"
	"fmt"
	"market_server/app/dataManager/rpc/internal/server"
	"market_server/common/config"
	"os"

	"github.com/zeromicro/go-zero/core/conf"
)

func test_config() {
	var c config.Config
	flag.Parse()

	fmt.Printf("Args: %+v \n", os.Args)
	env := os.Args[1]

	fmt.Printf("env: %+v \n", env)
	var configFile = flag.String("f", "etc/"+env+"/marketData.yaml", "the config file")

	conf.MustLoad(*configFile, &c)
	// fmt.Println(c.Nacos.IpAddr, ": ", c.Nacos.Port)
	fmt.Printf("%+v\n", c)
}

func main() {
	// flag.Parse()

	// fmt.Printf("Args: %+v \n", os.Args)
	// env := os.Args[1]

	// fmt.Printf("env: %+v \n", env)
	// var configFile = flag.String("f", "etc/"+env+"/marketData.yaml", "the config file")

	// var c config.Config
	// conf.MustLoad(*configFile, &c)

	// logx.MustSetup(c.LogConfig)

	// fmt.Printf("Log: %+v \n", c)

	// ctx := svc.NewServiceContext(c)
	// svr := server.NewMarketServiceServer(ctx)

	// s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
	// 	pb.RegisterMarketServiceServer(grpcServer, svr)

	// 	if c.Mode == service.DevMode || c.Mode == service.TestMode {
	// 		reflection.Register(grpcServer)
	// 	}
	// })
	// defer s.Stop()

	// fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)

	// svr.Start()
	// s.Start()

	server.TestMain()

	// test_config()
}