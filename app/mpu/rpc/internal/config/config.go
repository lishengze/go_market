package config

import "github.com/zeromicro/go-zero/zrpc"

type NacosConf struct {
	Server []struct {
		Host string
		Port uint64
	}
	Client struct {
		NamespaceId string
		TimeoutMs   uint64
		CacheDir    string
		LogDir      string
		LogLevel    string
	}

	DataId string
	Group  string
}

type KafkaConf struct {
	Address string
}

type Config struct {
	zrpc.RpcServerConf
	Exchange  string
	Proxy     string
	KafkaConf KafkaConf
	NacosConf NacosConf
}
