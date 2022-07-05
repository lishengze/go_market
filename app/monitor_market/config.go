package monitor_market

import "market_server/common/monitorStruct"

type DingConfig struct {
	secret string
	token  string
}

type Config struct {
	dingConfig    DingConfig
	monitorConfig monitorStruct.MonitorConfig
}
