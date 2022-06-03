package main

import (
	"bcts/pkg/kafkaclient"
	"flag"
	"fmt"
	"market_server/common/interceptor/rpcserver"

	"github.com/zeromicro/go-zero/core/logx"

	"market_server/app/order/rpc/internal/config"
	"market_server/app/order/rpc/internal/server"
	"market_server/app/order/rpc/internal/svc"
	"market_server/app/order/rpc/types/pb"

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
