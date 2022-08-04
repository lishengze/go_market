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
	"github.com/emirpasic/gods/utils"
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
func (k *KlineCache) InitWithHistKlines(klines []*Kline, symbol string, target_resolution int) {
	defer util.CatchExp("InitWithKlines")

	// k.KlinesMutex.Lock()
	// defer k.KlinesMutex.Unlock()

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

// UnTest
func (k *KlineCache) GetCurCacheKline(symbol string, resolution int) *Kline {
	defer util.CatchExp(fmt.Sprintf("GetCurCacheKline %s, %d", symbol, resolution))

	k.CacheKlineMutex.Lock()
	defer k.CacheKlineMutex.Unlock()

	if _, ok := k.CacheKlines[symbol]; !ok {
		return nil
	}

	if _, ok := k.CacheKlines[symbol][resolution]; !ok {
		return nil
	}

	return k.CacheKlines[symbol][resolution]
}

// UnTest
func (k *KlineCache) SetCacheKline(new_kline *Kline, resolution int) *Kline {
	defer util.CatchExp(fmt.Sprintf("SetCacheKline %s, %d", new_kline.String(), resolution))

	k.CacheKlineMutex.Lock()
	defer k.CacheKlineMutex.Unlock()

	if _, ok := k.CacheKlines[new_kline.Symbol]; !ok {
		k.CacheKlines[new_kline.Symbol] = make(map[int]*Kline)
	}

	if _, ok := k.CacheKlines[new_kline.Symbol][resolution]; !ok {
		k.CacheKlines[new_kline.Symbol][resolution] = NewKlineWithKline(new_kline)
	}

	k.CacheKlines[new_kline.Symbol][resolution] = NewKlineWithKline(new_kline)

	return new_kline
}

// UnTest
func (k *KlineCache) GetCurLastKline(symbol string, resolution int) *Kline {
	defer util.CatchExp(fmt.Sprintf("GetCurLastKline %s, %d", symbol, resolution))

	k.LastKlineMutex.Lock()
	defer k.LastKlineMutex.Unlock()

	if _, ok := k.LastKlines[symbol]; !ok {
		return nil
	}

	if _, ok := k.LastKlines[symbol][resolution]; !ok {
		return nil
	}

	return k.LastKlines[symbol][resolution]
}

// UnTest
func (k *KlineCache) SetLastKline(new_kline *Kline, resolution int) *Kline {
	defer util.CatchExp(fmt.Sprintf("SetLastKline %s, %d", new_kline.String(), resolution))

	k.LastKlineMutex.Lock()
	defer k.LastKlineMutex.Unlock()

	if _, ok := k.LastKlines[new_kline.Symbol]; !ok {
		k.LastKlines[new_kline.Symbol] = make(map[int]*Kline)
	}

	if _, ok := k.LastKlines[new_kline.Symbol][resolution]; !ok {
		k.LastKlines[new_kline.Symbol][resolution] = NewKlineWithKline(new_kline)
	}

	k.LastKlines[new_kline.Symbol][resolution] = NewKlineWithKline(new_kline)

	return new_kline
}

// UnTest
func (k *KlineCache) AddCompletedKline(new_kline *Kline, resolution int) {
	defer util.CatchExp("")

	k.KlinesMutex.Lock()
	defer k.KlinesMutex.Unlock()

	if _, ok := k.CompletedKlines[new_kline.Symbol]; !ok {
		k.CompletedKlines[new_kline.Symbol] = make(map[int]*treemap.Map)
	}

	if _, ok := k.CompletedKlines[new_kline.Symbol][resolution]; !ok {
		k.CompletedKlines[new_kline.Symbol][resolution] = treemap.NewWith(utils.Int64Comparator)
	}
	logx.Slowf("AddComplteKline: %s", new_kline.FullString())

	k.CompletedKlines[new_kline.Symbol][resolution].Put(new_kline.Time, new_kline)
}

func (k *KlineCache) CheckStoredKline(kline *Kline) bool {
	defer util.CatchExp(fmt.Sprintf("CheckStoredKline %s", kline.FullString()))

	k.KlinesMutex.Lock()
	defer k.KlinesMutex.Unlock()

	if _, ok := k.CompletedKlines[kline.Symbol]; !ok {
		return false
	}

	if _, ok := k.CompletedKlines[kline.Symbol][kline.Resolution]; !ok {
		return false
	}

	_, ok := k.CompletedKlines[kline.Symbol][kline.Resolution].Get(kline.Time)

	return ok
}

// UnTest
func (k *KlineCache) ProcessOldKline(new_kline *Kline, cache_kline *Kline, last_kline *Kline, resolution int) *Kline {
	defer util.CatchExp(fmt.Sprintf("ProcessOldKline \n%s\n%s\n%d", new_kline.FullString(), cache_kline.FullString(), resolution))

	if new_kline.Sequence == 0 {
		logx.Infof("KlineSource Restarted")
		k.ProcessLaterKline(new_kline, cache_kline, last_kline, resolution)

	} else {
		logx.Errorf("NewKline is Old: %s\nCacheKLine is: %s\n", new_kline.FullString(), cache_kline.FullString())
		logx.Slowf("NewKline is Old: NewKline.Seq: %d, CacheKline.Seq: %d", new_kline.Sequence, cache_kline.Sequence)

	}
	return nil
}

func (k *KlineCache) ProcessEqualKline(new_kline *Kline, cache_kline *Kline, last_kline *Kline, resolution int) *Kline {
	defer util.CatchExp(fmt.Sprintf("ProcessEqualKline \n%s\n%s\n%d", new_kline.FullString(), cache_kline.FullString(), resolution))

	var pub_kline *Kline = nil

	if !new_kline.IsHistory() {
		logx.Errorf("NewKLine Is Real: %s\nCacheKline: %s\nIs Same!", new_kline.FullString(), cache_kline.FullString())
		logx.Slowf("NewKLine Is Real: Is Same with CacheKline!", new_kline.FullString(), cache_kline.FullString())
		return nil
	}

	//Tested
	cache_kline.Volume = cache_kline.Volume + new_kline.Volume
	cache_kline.Sequence = new_kline.Sequence

	k.SetLastKline(new_kline, resolution)
	k.SetCacheKline(cache_kline, resolution)

	logx.Slowf("UpdateHistEqual CacheKline:%s", cache_kline.FullString())

	pub_kline = NewKlineWithKline(cache_kline)

	if IsOldKlineEnd(new_kline, int64(resolution)) {
		logx.Slowf("Old Kline End!")
		k.AddCompletedKline(cache_kline, resolution)
	}

	return pub_kline
}

//UnTested
func (k *KlineCache) ProcessOldMinuteWork(cache_kline *Kline, last_kline *Kline) {
	defer util.CatchExp(fmt.Sprintf("ProcessOldMinuteWork \ncache: %s\n%s", cache_kline.FullString(), last_kline.FullString()))

	if !last_kline.IsHistory() {
		cache_kline.Volume = cache_kline.Volume + last_kline.Volume

		logx.Slowf("OldMinuteNotFinished LastKline is Real:\nLastKLine:%s\nNewCache: %s", last_kline.FullString(), cache_kline.FullString())
	}

}

func (k *KlineCache) ProcessNewMinuteWork(new_kline *Kline, cache_kline *Kline, last_kline *Kline, resolution int) *Kline {
	defer util.CatchExp(fmt.Sprintf("ProcessNewMinuteWork \ncache: %s\n%s", cache_kline.FullString(), last_kline.FullString()))

	var pub_kline *Kline = nil

	if util.IsNewResolutionStart(new_kline.Time, last_kline.Time, resolution) { // Tested;
		logx.Slowf("NewResolutionStart!")

		if !k.CheckStoredKline(cache_kline) {
			k.AddCompletedKline(cache_kline, resolution)
		}

		cache_kline.ResetWithNewKline(new_kline)
		cache_kline.SetPerfectTime(int64(resolution))
		cache_kline.Volume = 0
		cache_kline.LastVolume = 0
		cache_kline.Resolution = resolution

		logx.Slowf("SetNewCache:%s", cache_kline.FullString())

		pub_kline = NewKlineWithKline(cache_kline)
		pub_kline.Volume = new_kline.Volume
	} else { // Tested
		cache_kline.UpdateInfoByRealKline(new_kline)
		logx.Slowf("UpdateCache: %s", cache_kline.FullString())

		pub_kline = NewKlineWithKline(cache_kline)
		pub_kline.Volume = cache_kline.Volume + new_kline.Volume
	}

	return pub_kline
}

// Tested
func (k *KlineCache) ProcessLaterRealKline(new_kline *Kline, cache_kline *Kline, last_kline *Kline, resolution int) *Kline {
	defer util.CatchExp(fmt.Sprintf("ProcessLaterRealKline \n%s\n%s\n%d", new_kline.FullString(), cache_kline.FullString(), resolution))

	var pub_kline *Kline = nil
	logx.Slowf("ProcessLaterRealKline")

	if util.IsNewMinuteStart(new_kline.Time, last_kline.Time) { // Tested
		logx.Slowf("NewMinuteStart: NewTime: %s", util.TimeStrFromInt(new_kline.Time))
		k.ProcessOldMinuteWork(cache_kline, last_kline)

		pub_kline = k.ProcessNewMinuteWork(new_kline, cache_kline, last_kline, resolution)
	} else { // Tested
		cache_kline.UpdateInfoByRealKline(new_kline)
		logx.Slowf("MiddleSecs UpdateCache:%s", cache_kline.FullString())

		pub_kline = NewKlineWithKline(cache_kline)
		pub_kline.Volume = cache_kline.Volume + new_kline.Volume
	}

	k.SetCacheKline(cache_kline, resolution)
	k.SetLastKline(new_kline, resolution)

	return pub_kline
}

// UnTest
func (k *KlineCache) ProcessLaterHistKline(new_kline *Kline, cache_kline *Kline, last_kline *Kline, resolution int) *Kline {
	defer util.CatchExp(fmt.Sprintf("ProcessLaterHistKline \n%s\n%s\n%d", new_kline.FullString(), cache_kline.FullString(), resolution))

	var pub_kline *Kline = nil

	logx.Slowf("ProcessLaterHistKline")

	if IsOldKlineEnd(new_kline, int64(resolution)) {

		logx.Slowf("OldKlineEnd,Resolution:%d, NewTime: %s", resolution, util.TimeStrFromInt(new_kline.Time))

		if new_kline.Resolution != resolution {
			cache_kline.UpdateInfoByHistKline(new_kline)
			logx.Slowf("UpdateLastCache: %s", cache_kline.FullString())
		} else {
			cache_kline.ResetWithNewKline(new_kline)
		}

		k.AddCompletedKline(cache_kline, resolution)

	} else if util.IsNewResolutionStart(new_kline.Time, last_kline.Time, resolution) {

		cache_kline = NewKlineWithKline(new_kline)
		cache_kline.Resolution = resolution
		cache_kline.LastVolume = 0
		cache_kline.SetPerfectTime(int64(resolution))

		logx.Slowf("NewKlineStart, SetCache: %s", cache_kline.FullString())

	} else { //Tested

		cache_kline.UpdateInfoByHistKline(new_kline)
		logx.Slowf("UpdateMiddleCache: %s", cache_kline.FullString())
	}

	k.SetCacheKline(cache_kline, resolution)
	k.SetLastKline(new_kline, resolution)
	pub_kline = NewKlineWithKline(cache_kline)

	return pub_kline
}

// UnTest
func (k *KlineCache) ProcessLaterKline(new_kline *Kline, cache_kline *Kline, last_kline *Kline, resolution int) *Kline {
	defer util.CatchExp(fmt.Sprintf("ProcessLaterKline \n%s\n%s\n%d", new_kline.FullString(), cache_kline.FullString(), resolution))
	var pub_kline *Kline = nil

	if new_kline.IsHistory() {
		pub_kline = k.ProcessLaterHistKline(new_kline, cache_kline, last_kline, resolution)
	} else { // Tested
		pub_kline = k.ProcessLaterRealKline(new_kline, cache_kline, last_kline, resolution)
	}
	return pub_kline
}

// Tested
func (k *KlineCache) InitCacheKline(new_kline *Kline, resolution int) *Kline {
	defer util.CatchExp(fmt.Sprintf("InitCacheKline: %d, %s", resolution, new_kline.FullString()))
	var pub_kline *Kline = nil

	kline := NewKlineWithKline(new_kline)
	kline.SetPerfectTime(int64(resolution))
	kline.Resolution = resolution
	kline.LastVolume = 0

	pub_kline = NewKlineWithKline(kline)

	if new_kline.IsHistory() {
		if IsOldKlineEnd(new_kline, int64(resolution)) {
			k.AddCompletedKline(kline, resolution)
		}
	} else {
		kline.Volume = 0
	}

	k.SetCacheKline(kline, resolution)
	k.SetLastKline(new_kline, resolution)

	logx.Slowf("InitCache: %s", kline.FullString())

	return pub_kline
}

// UnTest - 70%
func (k *KlineCache) UpdateWithKline(new_kline *Kline, resolution int) (*Kline, error) {
	defer util.CatchExp("UpdateWithKline")

	var pub_kline *Kline = nil
	var err error = nil

	// k.KlinesMutex.Lock()
	// defer k.KlinesMutex.Unlock()

	logx.Slowf("NewKline: %s", new_kline.FullString())

	cache_kline := k.GetCurCacheKline(new_kline.Symbol, resolution)
	last_kline := k.GetCurLastKline(new_kline.Symbol, resolution)

	if cache_kline == nil {
		pub_kline = k.InitCacheKline(new_kline, resolution)
	} else {
		logx.Slowf("cache_kline: %s", cache_kline.FullString())
		logx.Slowf("last_kline : %s", last_kline.FullString())

		if new_kline.Sequence < cache_kline.Sequence {
			pub_kline = k.ProcessOldKline(new_kline, cache_kline, resolution)
		} else if new_kline.Sequence == cache_kline.Sequence {
			pub_kline = k.ProcessEqualKline(new_kline, cache_kline, last_kline, resolution)
		} else {
			pub_kline = k.ProcessLaterKline(new_kline, cache_kline, last_kline, resolution)
		}
	}

	logx.Slowf("PubKline: %s\n", pub_kline.FullString())

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
