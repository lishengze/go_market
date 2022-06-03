package main

import (
	"flag"
	"fmt"

	"market_server/app/pms/rpc/internal/config"
	"market_server/app/pms/rpc/internal/server"
	"market_server/app/pms/rpc/internal/svc"
	"market_server/app/pms/rpc/types/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/pms.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	//kafka服务初始化
	//k := kafka.NewKafka(ctx)
	//k.AddRouter(kafka.Router{
	//	Handler: logic.NewTradeLogic,
	//})
	//k.RunKafka()
	//defer f()

	svr := server.NewPmsServer(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterPmsServer(grpcServer, svr)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
