package db_config

import (
	"market_server/common/config"
	"market_server/common/datastruct"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

type DBConfig struct {
	zrpc.RpcServerConf

	Comm      config.CommConfig
	Nacos     config.NacosConfig
	LogConfig logx.LogConf
	Mysql     config.MysqlConfig

	CacheConfig datastruct.CacheConfig
}
