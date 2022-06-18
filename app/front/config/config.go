package config

import (
	comm_config "market_server/common/config"

	"github.com/zeromicro/go-zero/core/logx"
)

type CacheConfig struct {
	CacheDataCount int
}

type Config struct {
	// RpcConfig zrpc.RpcClientConf
	Comm      comm_config.CommConfig
	LogConfig logx.LogConf
	Nacos     comm_config.NacosConfig

	Cache CacheConfig
}
