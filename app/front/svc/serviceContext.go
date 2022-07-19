package svc

import (
	"market_server/app/front/config"
	mkconfig "market_server/app/market_aggregate/config"
)

type ServiceContext struct {
	Config        config.Config
	SymbolConfigs []*mkconfig.SymbolConfig
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
