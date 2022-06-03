package config

import (
	"bcts/common/nacosAdapter"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Cache       cache.CacheConf
	Environment *Environment
	JwtAuth     struct {
		AccessSecret string
		AccessExpire int64
	}
	Mysql *struct {
		Addr string
	}
	Nacos              nacosAdapter.Config
	PriceRedis         redis.RedisConf
	Snowflake          *Snowflake
	TimeoutMarketPrice TimeoutMarketPrice
}

type Snowflake struct {
	Node     int64
	Epoch    int64
	NodeBits uint8
}

type TimeoutMarketPrice struct {
	Mtime uint64 //行情为获取时间
	Ctime uint64 //判断取USD还是USDT
	Coins []string
}

type Environment struct {
	Env string
}
