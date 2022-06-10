package data_engine

import (
	"market_server/app/front/config"
	"market_server/common/datastruct"
)

/*
接受数据的接口;
*/

type DataEngine struct {
	recvDataChan *datastruct.DataChannel
	cache_config *config.CacheConfig

	kline_cache map[string][]*datastruct.Kline
}

func NewDataEngine(recvDataChan *datastruct.DataChannel, config *config.CacheConfig) *DataEngine {
	rst := &DataEngine{
		recvDataChan: recvDataChan,
		cache_config: config,
	}

	return rst
}
