package monitor_market

import (
	comm_config "market_server/common/config"
	"market_server/common/monitorStruct"

	"github.com/zeromicro/go-zero/core/logx"
)

type DingConfig struct {
	Secret string
	Token  string
}

type WSConfig struct {
	Address string
	Url     string
}

type MonitorMeta struct {
	Symbols []string
}

type Config struct {
	WS   WSConfig
	Comm comm_config.CommConfig

	LogConfig logx.LogConf

	MonitorConfigInfo monitorStruct.MonitorConfig
	MonitorMetaInfo   MonitorMeta

	DingConfigInfo DingConfig
}