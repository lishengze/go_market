package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	IdSrvRpcConf zrpc.RpcClientConf

	Mysql struct {
		DataSource string
	}
	Cache cache.CacheConf

	KafkaConf KafkaConf
	Exchange  string
	Proxy     string
}

type KafkaConf struct {
	Address string
}
