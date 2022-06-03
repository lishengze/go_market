package main

import (
	"bcts/common/interceptor/rpcserver"
	"bcts/pkg/kafkaclient"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"

	"bcts/app/order/rpc/internal/config"
	"bcts/app/order/rpc/internal/server"
	"bcts/app/order/rpc/internal/svc"
	"bcts/app/order/rpc/types/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/order.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	svr := server.NewOrderServer(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterOrderServer(grpcServer, svr)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}

	})

	//rpc log
	s.AddUnaryInterceptors(rpcserver.LoggerInterceptor)
	defer s.Stop()
	logx.DisableStat() //关闭输出的统计日志(stat)

	TopicSymbolList := []string{"BTC_USDT", "ETH_USDT"}
	kclient := kafkaclient.NewKafkaClient(TopicSymbolList, "152.32.254.76", "9117")
	go kclient.FetchDepthWorkWather()
	//depthInfo := kclient.GetDepthInfo()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
