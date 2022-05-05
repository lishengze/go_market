package svc

import (
	"exterior-interactor/app/mpu/rpc/internal/config"
	"exterior-interactor/app/mpu/rpc/internal/core"
)

type ServiceContext struct {
	Config config.Config
	core.Mpu
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Mpu:    core.NewMpu(c),
	}
}
