package main

import (
	"flag"
	"fmt"
	"market_server/app/dataManager/rpc/internal/config"
	"market_server/app/dataManager/rpc/internal/server"
	"market_server/app/dataManager/rpc/internal/svc"
	"market_server/app/dataManager/rpc/types/pb"
	"os"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func test_config() {
	type TestConfig struct {
		zrpc.RpcServerConf

		IP            string
		NetServerType string
		SerialType    string

		Nacos     config.NacosConfig
		LogConfig logx.LogConf
		Mysql     config.MysqlConfig
	}

	var c TestConfig
	conf.MustLoad("marketData.yaml", &c)

	fmt.Printf("config: %+v \n", c)
}

func main() {
	flag.Parse()

	env := "local"

	for _, v := range os.Args {
		env = v
	}

	fmt.Printf("env: %+v \n", env)
	var configFile = flag.String("f", "etc/"+env+"/marketData.yaml", "the config file")

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	svr := server.NewMarketServiceServer(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterMarketServiceServer(grpcServer, svr)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)

	svr.Start()
	s.Start()

	// server.TestMain()

	// test_config()
}
