package svc

import (
	"exterior-interactor/app/opu/rpc/internal/config"
	"exterior-interactor/app/opu/rpc/internal/core"
)

type ServiceContext struct {
	Config config.Config
	core.OPU
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		OPU:    core.NewOpu(c),
	}
}
