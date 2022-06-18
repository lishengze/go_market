package config

import (
	comm_config "market_server/common/config"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

type CacheConfig struct {
}

type Config struct {
	RpcConfig zrpc.RpcClientConf
	Comm      comm_config.CommConfig
	Cache     CacheConfig
	LogConfig logx.LogConf
	Nacos     comm_config.NacosConfig
}
