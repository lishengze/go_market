package svc

import "market_server/app/data_manager/rpc/internal/dmconfig"

type ServiceContext struct {
	Config *dmconfig.ServerConfig
}

func NewServiceContext(c *dmconfig.ServerConfig) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
