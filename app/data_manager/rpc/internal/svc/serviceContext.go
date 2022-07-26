package svc

import "market_server/common/config"

type ServiceContext struct {
	Config db_config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
