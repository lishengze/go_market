package main

import (
	"flag"
	"fmt"
	"market_server/app/data_manager/rpc/internal/dmconfig"
	"market_server/app/data_manager/rpc/internal/server"
	"market_server/app/data_manager/rpc/internal/svc"
	"market_server/app/data_manager/rpc/types/pb"
	"market_server/common/config"
	"os"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	flag.Parse()

	fmt.Printf("Args: %+v \n", os.Args)
	env := os.Args[1]

	is_test := false
	if len(os.Args) > 2 {
		is_test = true
	}

	fmt.Printf("env: %+v \n", env)
	var configFile = flag.String("f", "etc/"+env+"/marketData.yaml", "the config file")
	fmt.Printf("configFile: %s", *configFile)

	var c dmconfig.ServerConfig
	conf.MustLoad(*configFile, &c)

	logx.MustSetup(c.LogConfig)
	logx.Infof("Log: %+v \n", c)

	fmt.Printf("Log: %+v \n", c)

	ctx := svc.NewServiceContext(&c)
	svr := server.NewMarketServiceServer(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterMarketServiceServer(grpcServer, svr)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)

	svr.SetTestValue(is_test)

	svr.Start()
	s.Start()

	// server.TestMain()
	// test_config()
}
