package main

import (
	"exterior-interactor/app/mpu/rpc/internal/config"
	"exterior-interactor/app/mpu/rpc/internal/server"
	"exterior-interactor/app/mpu/rpc/internal/svc"
	"exterior-interactor/app/mpu/rpc/mpupb"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/mpu.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	logx.MustSetup(c.Log)
	ctx := svc.NewServiceContext(c)
	srv := server.NewMpuServer(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		mpupb.RegisterMpuServer(grpcServer, srv)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
