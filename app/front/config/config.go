package config

import (
	comm_config "market_server/common/config"

	"github.com/zeromicro/go-zero/core/logx"
)

type CacheConfig struct {
}

type Config struct {
	Comm      comm_config.CommConfig
	Cache     CacheConfig
	LogConfig logx.LogConf
}
