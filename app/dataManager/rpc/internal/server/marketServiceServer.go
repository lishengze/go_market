// Code generated by goctl. DO NOT EDIT!
// Source: marketData.proto

package server

import (
	"context"
	"fmt"
	"os"
	"flag"

	"market_server/app/dataManager/rpc/internal/dbserver"
	"market_server/app/dataManager/rpc/internal/svc"
	"market_server/app/dataManager/rpc/config"
	"market_server/app/dataManager/rpc/types/pb"
	"market_server/common/comm"

	"market_server/common/datastruct"
	"github.com/zeromicro/go-zero/core/logx"


	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"	
)



type MarketServiceServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedMarketServiceServer

	dbServer *dbserver.DBServer
	commer *comm.Comm
	NacosClientWorker *config.NacosClient

	recvDataChan *datastruct.DataChannel
	pubDataChan  *datastruct.DataChannel
}

func NewMarketServiceServer(svcCtx *svc.ServiceContext) (*MarketServiceServer) {
	recv_data_chan := datastruct.NewDataChannel()
	pub_data_chan := datastruct.NewDataChannel()

	dbServer, err := dbserver.NewDBServer(recv_data_chan, svcCtx.Config.Mysql)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	rst := &MarketServiceServer{
		svcCtx: svcCtx,
		recvDataChan: recv_data_chan,
		pubDataChan: pub_data_chan,
		commer : comm.NewComm(recv_data_chan, pub_data_chan, svcCtx.Config.Comm),
		dbServer: dbServer,
	}

	return rst
}

func (m *MarketServiceServer) Start() {
	m.dbServer.StartListenRecvdata()

	m.commer.Start()

	m.SetInitMeta()

	// go m.StartNacosClient()
}

func (s *MarketServiceServer) StartNacosClient() {
	logx.Info("****************** StartNacosClient *****************")
	s.NacosClientWorker = config.NewNacosClient(&config.NATIVE_CONFIG().Nacos)

	logx.Info("Connect Nacos Successfully!")

	SymbolConfigStr, err := s.NacosClientWorker.GetConfigContent("SymbolParams", datastruct.BCTS_GROUP)
	if err != nil {
		logx.Error(err.Error())
	}
	logx.Info("Requested SymbolConfigStr: " + SymbolConfigStr)
	s.ProcessSymbolConfigStr(SymbolConfigStr)

	s.NacosClientWorker.ListenConfig("SymbolParams", datastruct.BCTS_GROUP, s.SymbolParamsChanged)
}

func (s *MarketServiceServer) SymbolParamsChanged(namespace, group, dataId, data string) {
	logx.Infof("SymbolContent: %s\n", data)
	s.ProcessSymbolConfigStr(data)
}

func (s *MarketServiceServer) ProcessSymbolConfigStr(data string) {
	symbol_configs, err := config.ParseJsonSymbolConfig(data)

	if err != nil {
		logx.Error(err.Error())
		return
	}

	symbol_exchange_set := make(map[string](map[string]struct{}))
	new_meta := datastruct.Metadata{}

	for _, symbol_config := range symbol_configs {
		if _,ok := symbol_exchange_set[symbol_config.Symbol];!ok {
			symbol_exchange_set[symbol_config.Symbol] = make(map[string]struct{})
		}
		if _,ok := symbol_exchange_set[symbol_config.Symbol][datastruct.BCTS_EXCHANGE]; !ok {
			symbol_exchange_set[symbol_config.Symbol][datastruct.BCTS_EXCHANGE] = struct{}{}
		}

		
	}

	new_meta.TradeMeta = symbol_exchange_set
	new_meta.KlineMeta = symbol_exchange_set	

	logx.Infof("NewMeta: %v \n", new_meta)

	s.commer.UpdateMetaData(&new_meta)
}

func (s *MarketServiceServer) SetInitMeta() {
	init_symbol_list := []string{"BTC_USDT", "ETH_USDT", "USD_USDT","BTC_USD", "ETH_USD"}

	symbol_exchange_set := make(map[string](map[string]struct{}))
	new_meta := datastruct.Metadata{}
	for _,symbol := range init_symbol_list {
		if _,ok := symbol_exchange_set[symbol];!ok {
			symbol_exchange_set[symbol] = make(map[string]struct{})
		}
		if _,ok := symbol_exchange_set[symbol][datastruct.BCTS_EXCHANGE]; !ok {
			symbol_exchange_set[symbol][datastruct.BCTS_EXCHANGE] = struct{}{}
		}
	}

	new_meta.TradeMeta = symbol_exchange_set
	new_meta.KlineMeta = symbol_exchange_set	

	logx.Infof("[I] InitMeta: %v \n", new_meta)

	s.commer.UpdateMetaData(&new_meta)
}

//  服务名字还是用 MarketService, 整个行情系统都要用；
func (m *MarketServiceServer) RequestHistKlineData(ctx context.Context, in *pb.ReqHishKlineInfo) (*pb.HistKlineData, error) {
	
	return m.dbServer.RequestHistKlineData(ctx, in)
}

func (m *MarketServiceServer) RequestTradeData(ctx context.Context, in *pb.ReqTradeInfo) (*pb.Trade, error) {
	return m.dbServer.RequestTradeData(ctx, in)
}


func TestMain() {
	flag.Parse()

	env := "local"

	for _, v := range os.Args {
		env = v
	}

	fmt.Printf("env: %+v \n", env)
	var configFile = flag.String("f", "etc/"+env+"/marketData.yaml", "the config file")

	fmt.Println(*configFile)

	var c config.Config
	conf.MustLoad(*configFile, &c)

	fmt.Printf("config: %+v \n",c)

	logx.MustSetup(c.LogConfig)

	ctx := svc.NewServiceContext(c)
	svr := NewMarketServiceServer(ctx)

	// return

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

	select{}
}
