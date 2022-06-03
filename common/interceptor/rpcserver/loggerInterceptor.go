package rpcserver

import (
	"bcts/common/xerror"
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//func LoggerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
//	resp, err = handler(ctx, req)
//	if err != nil {
//		//errMsg := err.Error()
//		causeErr := errors.Cause(err)                  // err类型
//		if e, ok := causeErr.(*xerror.CodeError); ok { //自定义错误类型
//			logx.WithContext(ctx).Errorf("[RPC-SRV-ERR] %+v", err)
//
//			//转成grpc err
//			//err = status.Error(codes.Code(e.ErrCode()), e.ErrorMsg())
//			err = status.Error(codes.Code(e.ErrCode()), err.Error())
//
//		} else {
//			logx.WithContext(ctx).Errorf("[RPC-SRV-ERR] %+v", err)
//		}
//	}
//	return resp, err
//}

func LoggerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	resp, err = handler(ctx, req)
	if err != nil {
		causeErr := errors.Cause(err)                  // err类型
		if e, ok := causeErr.(*xerror.CodeError); ok { //自定义错误类型
			logx.WithContext(ctx).Errorf("[RPC-SRV-ERR] %+v", err)

			//转成grpc err
			err = status.Error(codes.Code(e.ErrCode()), e.ErrorMsg())
		} else {
			logx.WithContext(ctx).Errorf("[RPC-SRV-ERR] %+v", err)
		}
	}
	return resp, err
}
