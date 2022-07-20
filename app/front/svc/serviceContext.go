package svc

import (
	"market_server/app/front/config"
	mkconfig "market_server/app/market_aggregate/config"
	"sync"
)

type ServiceContext struct {
	Config config.Config

	symbol_config_mutex sync.Mutex
	SymbolConfigs       map[string]*mkconfig.SymbolConfig
}

func (s *ServiceContext) SetConfig(symbol_configs map[string]*mkconfig.SymbolConfig) {
	s.symbol_config_mutex.Lock()

	defer s.symbol_config_mutex.Unlock()

	s.SymbolConfigs = symbol_configs
}

func (s *ServiceContext) UpdateSymbolConfigWithSlice(symbol_configs []*mkconfig.SymbolConfig) {
	s.symbol_config_mutex.Lock()

	defer s.symbol_config_mutex.Unlock()

	for _, configs := range symbol_configs {
		s.SymbolConfigs[configs.Symbol] = configs
	}

}

func (s *ServiceContext) GetSymbolConfig(symbol string) *mkconfig.SymbolConfig {
	s.symbol_config_mutex.Lock()

	defer s.symbol_config_mutex.Unlock()

	if config, ok := s.SymbolConfigs[symbol]; ok {
		return config
	}
	return nil
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		SymbolConfigs: make(map[string]*mkconfig.SymbolConfig),
	}
}
