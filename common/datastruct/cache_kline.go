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
	Klines map[string]map[int]*treemap.Map

	RealKlineByKline map[string]map[int]*Kline
	RealKlineByTrade map[string]map[int]*Kline

	Config *CacheConfig

	KlinesMutex           sync.Mutex
	RealKlineByKlineMutex sync.Mutex
	RealKlineByTradeMutex sync.Mutex
}

func NewKlineCache(config *CacheConfig) *KlineCache {
	return &KlineCache{
		Klines:           make(map[string]map[int]*treemap.Map),
		Config:           config,
		RealKlineByKline: make(map[string]map[int]*Kline),
		RealKlineByTrade: make(map[string]map[int]*Kline),
	}
}

//Undo -- Cache Must Init First!
func (k *KlineCache) InitWithKlines(klines []*Kline, symbol string, target_resolution int) {
	defer util.CatchExp("InitWithKlines")

	k.KlinesMutex.Lock()
	defer k.KlinesMutex.Unlock()

	if _, ok := k.Klines[symbol]; !ok {
		k.Klines[symbol] = make(map[int]*treemap.Map)
		logx.Infof("KlineCache Add New Symbol: %s", symbol)
	}

	if _, ok := k.Klines[symbol][target_resolution]; !ok {
		k.Klines[symbol][target_resolution] = NewTreeMapWithKlines(klines, target_resolution)

		if k.Klines[symbol][target_resolution] != nil {
			latest_kline := GetLastKline(k.Klines[symbol][target_resolution])
			k.SetKlineCacheKline(latest_kline)
			k.SetTradeCacheKline(latest_kline)
		} else {
			logx.Errorf("%s.%d kline Init Failed", symbol, target_resolution)
		}

	} else {
		logx.Errorf(" KlineCache %s.%d already cached", symbol, target_resolution)
	}
}

func (k *KlineCache) SetKlineCacheKline(kline *Kline) {
	defer util.CatchExp(fmt.Sprintf("SetKlineCacheKline %s", kline.String()))

	k.RealKlineByKlineMutex.Lock()
	defer k.RealKlineByKlineMutex.Unlock()

	if _, ok := k.RealKlineByKline[kline.Symbol]; !ok {
		k.RealKlineByKline[kline.Symbol] = make(map[int]*Kline)
	}
	k.RealKlineByKline[kline.Symbol][kline.Resolution] = kline
}

func (k *KlineCache) SetTradeCacheKline(kline *Kline) {
	defer util.CatchExp(fmt.Sprintf("SetTradeCacheKline %s", kline.String()))
	k.RealKlineByTradeMutex.Lock()
	defer k.RealKlineByTradeMutex.Unlock()

	if _, ok := k.RealKlineByTrade[kline.Symbol]; !ok {
		k.RealKlineByTrade[kline.Symbol] = make(map[int]*Kline)
	}
	k.RealKlineByTrade[kline.Symbol][kline.Resolution] = kline
}

func (k *KlineCache) UpdateWithKlines(ori_klines []*Kline, symbol string) error {
	defer util.CatchExp("UpdateWithKlines")

	// k.KlinesMutex.Lock()
	// defer k.KlinesMutex.Unlock()

	// if _, ok := k.Klines[symbol]; !ok {
	// 	return fmt.Errorf(" UpdateWithKlines Failed , symbol%s was not init", symbol)
	// }

	// for resolution, kline_tree := range k.Klines[symbol] {
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

func (k *KlineCache) UpdateWithKline(kline *Kline) (*Kline, error) {
	defer util.CatchExp("UpdateWithKline")

	k.RealKlineByKlineMutex.Lock()
	defer k.RealKlineByKlineMutex.Unlock()

	var pub_kline *Kline = nil
	var err error = nil

	if _, ok := k.RealKlineByKline[kline.Symbol]; !ok {
		return nil, fmt.Errorf(" UpdateWithKlines Failed , symbol%s was not init", kline.Symbol)
	}

	for resolution, cache_kline := range k.RealKlineByKline[kline.Symbol] {
		pub_kline, is_add := UpdateTreeWithKline(cache_kline, kline, resolution)

		if is_add {
			k.AddNewKline(pub_kline)
		}
	}

	k.EraseOutTimeKline()

	return pub_kline, err
}

func (k *KlineCache) UpdateWithTrade(trade *Trade) (*Kline, error) {
	defer util.CatchExp("UpdateWithKline")

	k.RealKlineByTradeMutex.Lock()
	defer k.RealKlineByTradeMutex.Unlock()

	var pub_kline *Kline = nil
	var err error = nil

	if _, ok := k.RealKlineByTrade[trade.Symbol]; !ok {
		return nil, fmt.Errorf(" UpdateWithKlines Failed , symbol%s was not init", trade.Symbol)
	}

	for resolution, cache_kline := range k.RealKlineByTrade[trade.Symbol] {
		pub_kline, is_add := UpdateTreeWithTrade(cache_kline, trade, resolution)

		if is_add {
			k.AddNewKline(pub_kline)
		}
	}

	k.EraseOutTimeKline()

	return pub_kline, err
}

// Undo
func (k *KlineCache) EraseOutTimeKline() {

}

//Undo
func (k *KlineCache) GetKlinesByCount(symbol string, resolution int, count int, get_most bool) []*Kline {
	defer util.CatchExp("GetKlinesByCount")

	k.KlinesMutex.Lock()
	defer k.KlinesMutex.Unlock()

	var rst []*Kline

	if _, ok := k.Klines[symbol]; !ok {
		return nil
	}

	if _, ok := k.Klines[symbol][resolution]; !ok {
		return nil
	}

	if k.Klines[symbol][resolution].Size() == 0 {
		return nil
	}

	// 缓存的数目不够，并且强制要求得到 count 数量的数据;
	if k.Klines[symbol][resolution].Size() < count && k.Klines[symbol][resolution].Size() >= k.Config.Count && get_most {
		return nil
	}

	// 将当前的缓存的数据发送出去;
	iter := k.Klines[symbol][resolution].Iterator()
	start_pos := k.Klines[symbol][resolution].Size() - count
	index := 0

	for iter.Begin(); iter.Next(); {
		if index > start_pos {
			rst = append(rst, iter.Value().(*Kline))
		}
		index++
	}

	return rst
}

//Undo
func (k *KlineCache) GetKlinesByTime(symbol string, resolution int, start_time int64, end_time int64, get_most bool) []*Kline {
	defer util.CatchExp("GetKlinesByTime")

	k.KlinesMutex.Lock()
	defer k.KlinesMutex.Unlock()

	var rst []*Kline

	if _, ok := k.Klines[symbol]; !ok {
		return nil
	}

	if _, ok := k.Klines[symbol][resolution]; !ok {
		return nil
	}

	if k.Klines[symbol][resolution].Size() == 0 {
		return nil
	}

	// 缓存的数目不够，并且强制要求得到 count 数量 的数据;
	iter := k.Klines[symbol][resolution].Iterator()
	iter.First()

	if iter.Key().(int64) > start_time && k.Klines[symbol][resolution].Size() >= k.Config.Count && get_most {
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
