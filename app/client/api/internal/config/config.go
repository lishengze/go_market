package config

import (
	"bcts/common/nacosAdapter"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf

	Mysql struct {
		DataSource string
	}
	Cache cache.CacheConf

	Log logx.LogConf

	ResponseEncrypted bool //是否将返回消息进行加密

	NacosConfig *nacosAdapter.Config

	PemFileConfig struct {
		PriPemFilePath string
		PubPemFilePath string
	}

	UserCenterRpc zrpc.RpcClientConf //用户中心rpc
	OrderRpc      zrpc.RpcClientConf
}
