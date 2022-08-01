package dmconfig

import (
	"market_server/common/config"
	"market_server/common/datastruct"
	"market_server/common/monitorStruct"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServerConfig struct {
	zrpc.RpcServerConf

	Comm      config.CommConfig
	Nacos     config.NacosConfig
	LogConfig logx.LogConf
	Mysql     config.MysqlConfig

	CacheConfig datastruct.CacheConfig

	MonitorConfigInfo monitorStruct.MonitorConfig
	DingConfigInfo    config.DingConfig
}
