/*
Author: lsz;
CreateTime: 2022.07.26
Des: Cache Multi Frequency Kline Data And Support Update Realtime;
*/

package datastruct

import (
	"fmt"
	"market_server/common/util"
	"sync"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/zeromicro/go-zero/core/logx"
)

type CacheConfig struct {
	Count int
}

type KlineCache struct {
	CompletedKlines map[string]map[int]*treemap.Map

	CacheKlines map[string]map[int]*Kline
	LastKlines  map[string]map[int]*Kline

	Config *CacheConfig

	KlinesMutex     sync.Mutex
	CacheKlineMutex sync.Mutex
	LastKlineMutex  sync.Mutex
}

func NewKlineCache(config *CacheConfig) *KlineCache {
	return &KlineCache{
		CompletedKlines: make(map[string]map[int]*treemap.Map),
		Config:          config,
		CacheKlines:     make(map[string]map[int]*Kline),
		LastKlines:      make(map[string]map[int]*Kline),
	}
}

//Cache Must Init First!
func (k *KlineCache) InitWithKlines(klines []*Kline, symbol string, target_resolution int) {
	defer util.CatchExp("InitWithKlines")

	k.KlinesMutex.Lock()
	defer k.KlinesMutex.Unlock()

	if _, ok := k.CompletedKlines[symbol]; !ok {
		k.CompletedKlines[symbol] = make(map[int]*treemap.Map)
		logx.Infof("KlineCache Add New Symbol: %s", symbol)
	}

	if _, ok := k.CompletedKlines[symbol][target_resolution]; !ok {
		k.CompletedKlines[symbol][target_resolution] = NewTreeMapWithKlines(klines, target_resolution)

		// if k.CompletedKlines[symbol][target_resolution] != nil {
		// 	latest_kline := GetLastKline(k.CompletedKlines[symbol][target_resolution])

		// } else {
		// 	logx.Errorf("%s.%d kline Init Failed", symbol, target_resolution)
		// }

	} else {
		logx.Errorf(" KlineCache %s.%d already cached", symbol, target_resolution)
	}
}

// Undo
func (k *KlineCache) UpdateWithKlines(ori_klines []*Kline, symbol string) error {
	defer util.CatchExp("UpdateWithKlines")

	// k.KlinesMutex.Lock()
	// defer k.KlinesMutex.Unlock()

	// if _, ok := k.CompletedKlines[symbol]; !ok {
	// 	return fmt.Errorf(" UpdateWithKlines Failed , symbol%s was not init", symbol)
	// }

	// for resolution, kline_tree := range k.CompletedKlines[symbol] {
	// 	UpdateTreeWithKlines(kline_tree, ori_klines, resolution)
	// }

	// k.EraseOutTimeKline()

	return nil
}

func (k *KlineCache) AddNewKline(kline *Kline) {
	util.CatchExp(fmt.Sprintf("AddNewKline: %s", kline.String()))
	k.KlinesMutex.Lock()
	defer k.KlinesMutex.Unlock()

}

func (k *KlineCache) GetCurCacheKline(symbol string, resolution int) *Kline {
	defer util.CatchExp("GetCurCacheKline")

	k.CacheKlineMutex.Lock()
	defer k.CacheKlineMutex.Unlock()

	return nil
}

func (k *KlineCache) AddNewCacheLine(new_kline *Kline, resolution int) *Kline {
	return nil
}

func (k *KlineCache) UpdateWithKline(new_kline *Kline, resolution int) (*Kline, error) {
	defer util.CatchExp("UpdateWithKline")

	var pub_kline *Kline = nil
	var err error = nil

	k.KlinesMutex.Lock()
	defer k.KlinesMutex.Unlock()

	cache_kline := k.GetCurCacheKline(new_kline.Symbol, resolution)

	if cache_kline == nil {

	} else {

	}

	k.EraseOutTimeKline()

	return pub_kline, err
}

// UnTest
func EraseTreeNode(tree *treemap.Map, target_count int) {
	defer util.CatchExp("EraseTreeNode")

	var erase_times []int64
	iter := tree.Iterator()
	iter.First()
	for i := 0; tree.Size()-i > target_count; i++ {
		erase_times = append(erase_times, iter.Key().(int64))
		iter.Next()
	}

	for _, outtimeKey := range erase_times {
		tree.Remove(outtimeKey)
	}
}

// UnTest
func (k *KlineCache) EraseOutTimeKline() {
	defer util.CatchExp("EraseOutTimeKline")
	k.KlinesMutex.Lock()
	defer k.KlinesMutex.Unlock()

	for _, s_klines := range k.CompletedKlines {
		for _, klineTree := range s_klines {
			EraseTreeNode(klineTree, k.Config.Count)
		}
	}
}

func (k *KlineCache) GetKlinesByCount(symbol string, resolution int, count int, get_most bool) []*Kline {
	defer util.CatchExp("GetKlinesByCount")

	k.KlinesMutex.Lock()
	defer k.KlinesMutex.Unlock()

	var rst []*Kline

	if _, ok := k.CompletedKlines[symbol]; !ok {
		return nil
	}

	if _, ok := k.CompletedKlines[symbol][resolution]; !ok {
		return nil
	}

	if k.CompletedKlines[symbol][resolution].Size() == 0 {
		return nil
	}

	// 缓存的数目不够，并且强制要求得到 count 数量的数据;
	if k.CompletedKlines[symbol][resolution].Size() < count && k.CompletedKlines[symbol][resolution].Size() >= k.Config.Count && get_most {
		return nil
	}

	// 将当前的缓存的数据发送出去;
	iter := k.CompletedKlines[symbol][resolution].Iterator()
	start_pos := k.CompletedKlines[symbol][resolution].Size() - count
	index := 0

	for iter.Begin(); iter.Next(); {
		if index > start_pos {
			rst = append(rst, iter.Value().(*Kline))
		}
		index++
	}

	return rst
}

func (k *KlineCache) GetKlinesByTime(symbol string, resolution int, start_time int64, end_time int64, get_most bool) []*Kline {
	defer util.CatchExp("GetKlinesByTime")

	k.KlinesMutex.Lock()
	defer k.KlinesMutex.Unlock()

	var rst []*Kline

	if _, ok := k.CompletedKlines[symbol]; !ok {
		return nil
	}

	if _, ok := k.CompletedKlines[symbol][resolution]; !ok {
		return nil
	}

	if k.CompletedKlines[symbol][resolution].Size() == 0 {
		return nil
	}

	// 缓存的数目不够，并且强制要求得到 count 数量 的数据;
	iter := k.CompletedKlines[symbol][resolution].Iterator()
	iter.First()

	if iter.Key().(int64) > start_time && k.CompletedKlines[symbol][resolution].Size() >= k.Config.Count && get_most {
		return nil
	}

	// 将当前的缓存的数据发送出去;
	for iter.Begin(); iter.Next(); {
		if iter.Key().(int64) >= start_time && iter.Key().(int64) < end_time {
			rst = append(rst, iter.Value().(*Kline))
		}
	}

	return rst
}
