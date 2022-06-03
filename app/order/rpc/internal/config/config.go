package config

import (
	"market_server/common/nacosAdapter"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	Mysql struct {
		DataSource string
	}
	Cache cache.CacheConf

	Redis redis.RedisConf

	Log logx.LogConf

	Nacos *nacosAdapter.Config

	UserCenterRpc zrpc.RpcClientConf //userCenter rpc

	DataServiceRpc zrpc.RpcClientConf //dataService rpc
}
