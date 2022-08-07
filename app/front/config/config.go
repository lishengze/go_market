package config

import (
	comm_config "market_server/common/config"
	"market_server/common/datastruct"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

type CacheConfig struct {
	CacheDataCount int
}

type WSConfig struct {
	Address           string
	Url               string
	HeartbeatLostSecs int
	HeartbeatSendSecs int
}

type Config struct {
	RpcConfig zrpc.RpcClientConf
	Comm      comm_config.CommConfig
	LogConfig logx.LogConf
	Nacos     comm_config.NacosConfig

	WS          WSConfig
	CacheConfig datastruct.CacheConfig
}
