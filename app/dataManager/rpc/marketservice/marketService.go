// Code generated by goctl. DO NOT EDIT!
// Source: marketData.proto

package marketservice

import (
	"context"

	"bcts/app/dataManager/rpc/types/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	Depth            = pb.Depth
	EmptyReq         = pb.EmptyReq
	EmptyRsp         = pb.EmptyRsp
	HistKlineData    = pb.HistKlineData
	Kline            = pb.Kline
	PriceVolume      = pb.PriceVolume
	ReqHishKlineInfo = pb.ReqHishKlineInfo
	ReqTradeInfo     = pb.ReqTradeInfo
	Trade            = pb.Trade

	MarketService interface {
		//  服务名字还是用 MarketService, 整个行情系统都要用；
		RequestHistKlineData(ctx context.Context, in *ReqHishKlineInfo, opts ...grpc.CallOption) (*HistKlineData, error)
		RequestTradeData(ctx context.Context, in *ReqTradeInfo, opts ...grpc.CallOption) (*Trade, error)
	}

	defaultMarketService struct {
		cli zrpc.Client
	}
)

func NewMarketService(cli zrpc.Client) MarketService {
	return &defaultMarketService{
		cli: cli,
	}
}

//  服务名字还是用 MarketService, 整个行情系统都要用；
func (m *defaultMarketService) RequestHistKlineData(ctx context.Context, in *ReqHishKlineInfo, opts ...grpc.CallOption) (*HistKlineData, error) {
	client := pb.NewMarketServiceClient(m.cli.Conn())
	return client.RequestHistKlineData(ctx, in, opts...)
}

func (m *defaultMarketService) RequestTradeData(ctx context.Context, in *ReqTradeInfo, opts ...grpc.CallOption) (*Trade, error) {
	client := pb.NewMarketServiceClient(m.cli.Conn())
	return client.RequestTradeData(ctx, in, opts...)
}
