package config

import (
	"bcts/common/nacosAdapter"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	NacosConfig *nacosAdapter.Config
}
