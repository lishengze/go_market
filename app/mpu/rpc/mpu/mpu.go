// Code generated by goctl. DO NOT EDIT!
// Source: mpu.proto

package mpu

import (
	"context"

	"exterior-interactor/app/mpu/rpc/mpupb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	Depth            = mpupb.Depth
	EmptyReq         = mpupb.EmptyReq
	EmptyRsp         = mpupb.EmptyRsp
	GetAllSymbolsRsp = mpupb.GetAllSymbolsRsp
	GetSymbolReq     = mpupb.GetSymbolReq
	GetSymbolRsp     = mpupb.GetSymbolRsp
	Kline            = mpupb.Kline
	PriceVolume      = mpupb.PriceVolume
	Symbol           = mpupb.Symbol
	Trade            = mpupb.Trade

	Mpu interface {
		GetAllSymbols(ctx context.Context, in *EmptyReq, opts ...grpc.CallOption) (*GetAllSymbolsRsp, error)
		GetSymbol(ctx context.Context, in *GetSymbolReq, opts ...grpc.CallOption) (*GetSymbolRsp, error)
	}

	defaultMpu struct {
		cli zrpc.Client
	}
)

func NewMpu(cli zrpc.Client) Mpu {
	return &defaultMpu{
		cli: cli,
	}
}

func (m *defaultMpu) GetAllSymbols(ctx context.Context, in *EmptyReq, opts ...grpc.CallOption) (*GetAllSymbolsRsp, error) {
	client := mpupb.NewMpuClient(m.cli.Conn())
	return client.GetAllSymbols(ctx, in, opts...)
}

func (m *defaultMpu) GetSymbol(ctx context.Context, in *GetSymbolReq, opts ...grpc.CallOption) (*GetSymbolRsp, error) {
	client := mpupb.NewMpuClient(m.cli.Conn())
	return client.GetSymbol(ctx, in, opts...)
}