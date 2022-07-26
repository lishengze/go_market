/*
Author: lsz;
CreateTime: 2022.07.26
Des: Cache Multi Frequency Kline Data And Support Update Realtime;
*/

package datastruct

import (
	"market_server/common/util"

	"github.com/emirpasic/gods/maps/treemap"
)

type CacheConfig struct {
	Count int
}

type KlineCache struct {
	Klines map[string]map[int]*treemap.Map
	Config *CacheConfig
}

func NewKlineCache(config *CacheConfig) *KlineCache {
	return &KlineCache{
		Klines: make(map[string]map[int]*treemap.Map),
		Config: config,
	}
}

func ResetFirstKline(cache_kline *Kline, target_resolution int) {
	defer util.CatchExp("ResetFirstKline")
	cache_kline.Time = GetLastStartTime(cache_kline.Time, int64(target_resolution))
}

func ProcessOldEndKline(cur_kline *Kline, cache_kline *Kline, target_resolution int) *Kline {
	defer util.CatchExp("ProcessOldEndKline")
	var pub_kline *Kline
	if cur_kline.Resolution != target_resolution {
		cache_kline.Close = cur_kline.Close
		cache_kline.Low = util.MinFloat64(cache_kline.Low, cur_kline.Low)
		cache_kline.High = util.MaxFloat64(cache_kline.High, cur_kline.High)
		cache_kline.Volume += cur_kline.Volume

		pub_kline = NewKlineWithKline(cache_kline)
	} else {
		pub_kline = NewKlineWithKline(cur_kline)
		cache_kline = NewKlineWithKline(pub_kline)
	}
	return pub_kline
}

func ProcessNewStartKline(cur_kline *Kline, cache_kline *Kline, target_resolution int) {
	defer util.CatchExp("ProcessNewStartKline")

	cache_kline = NewKlineWithKline(cur_kline)
	cache_kline.Resolution = target_resolution
}

func ProcessCachingKline(cur_kline *Kline, cache_kline *Kline) {
	defer util.CatchExp("ProcessCachingKline")

	cache_kline.Close = cur_kline.Close
	cache_kline.Low = util.MinFloat64(cache_kline.Low, cur_kline.Low)
	cache_kline.High = util.MaxFloat64(cache_kline.High, cur_kline.High)
	cache_kline.Volume = cache_kline.Volume + cur_kline.Volume
}

func NewTreeMapWithKlines(ori_klines []*Kline, target_resolution int) *treemap.Map {
	defer util.CatchExp("NewTreeMapWithKlines")

	var rst *treemap.Map

	var cache_kline *Kline = nil
	var pub_kline *Kline = nil

	for _, cur_kline := range ori_klines {
		if cache_kline == nil {
			cache_kline = cur_kline
			ResetFirstKline(cache_kline, target_resolution)
		}

		if IsOldKlineEndTime(cur_kline.Time, int(cur_kline.Resolution), int64(target_resolution)) {
			pub_kline = ProcessOldEndKline(cur_kline, cache_kline, int(target_resolution))
			rst.Put(pub_kline.Time, pub_kline)
		} else if IsNewKlineStartTime(cur_kline.Time, int64(target_resolution)) {
			ProcessNewStartKline(cur_kline, cache_kline, target_resolution)
		} else {
			ProcessCachingKline(cur_kline, cache_kline)
		}
	}

	if cache_kline.Time != pub_kline.Time {
		rst.Put(cache_kline.Time, pub_kline)
	}

	return rst
}

func UpdateTreeWithKlines(kline_tree *treemap.Map, klines []*Kline, target_resolution int) *Kline {
	defer util.CatchExp("UpdateTreeWithKlines")

	var rst *Kline

	return rst
}

func UpdateTreeWithKline(kline_tree *treemap.Map, kline *Kline, target_resolution int) *Kline {
	defer util.CatchExp("UpdateTreeWithKline")

	var rst *Kline

	return rst
}

func UpdateTreeWithTrade(kline_tree *treemap.Map, trade *Trade, target_resolution int) *Kline {
	defer util.CatchExp("UpdateTreeWithKline")

	var rst *Kline

	return rst
}

func TreeGetKlinesByCount(kline_tree *treemap.Map, count int, resolution int) []*Kline {
	defer util.CatchExp("TreeGetKlinesByCount")

	var rst []*Kline

	return rst
}

func (k *KlineCache) InitWithKlines(klines []*Kline, target_resolution int) {
	defer util.CatchExp("InitWithKlines")

}

func (k *KlineCache) UpdateWithKlines(klines []*Kline) {
	defer util.CatchExp("UpdateWithKlines")

}

func (k *KlineCache) UpdateWithKline(kline *Kline) {
	defer util.CatchExp("UpdateWithKline")

}

func (k *KlineCache) GetKlinesByCount(symbol string, resolution int, count int) []*Kline {
	defer util.CatchExp("GetKlinesByCount")
	var rst []*Kline

	return rst
}

func (k *KlineCache) GetKlinesByTime(symbol string, resolution int, start_time int64, end_time int64) []*Kline {
	defer util.CatchExp("GetKlinesByTime")

	var rst []*Kline

	return rst
}
