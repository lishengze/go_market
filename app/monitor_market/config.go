package monitor_market

import "market_server/common/monitorStruct"

type DingConfig struct {
	secret string
	token  string
}

type WSConfig struct {
	Address string
	Url     string
}

type Config struct {
	WsConfig      WSConfig
	dingConfig    DingConfig
	monitorConfig monitorStruct.MonitorConfig
}
