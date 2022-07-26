package datastruct

import (
	"market_server/common/util"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

func TransSliceKlines(ori_klines []*Kline) *treemap.Map {
	var rst *treemap.Map

	for _, kline := range ori_klines {
		rst.Put(kline.Time, kline)
	}

	return rst
}

func ResetFirstKline(latest_kline *Kline, target_resolution int) {
	defer util.CatchExp("ResetFirstKline")
	latest_kline.Time = GetLastStartTime(latest_kline.Time, int64(target_resolution))
	latest_kline.Resolution = target_resolution
}

func ProcessOldEndKline(cur_kline *Kline, latest_kline *Kline, target_resolution int) *Kline {
	defer util.CatchExp("ProcessOldEndKline")
	var pub_kline *Kline
	if cur_kline.Resolution != target_resolution {
		latest_kline.Close = cur_kline.Close
		latest_kline.Low = util.MinFloat64(latest_kline.Low, cur_kline.Low)
		latest_kline.High = util.MaxFloat64(latest_kline.High, cur_kline.High)
		latest_kline.Volume += cur_kline.Volume

		pub_kline = NewKlineWithKline(latest_kline)
	} else {
		pub_kline = NewKlineWithKline(cur_kline)
		latest_kline = NewKlineWithKline(pub_kline)
	}
	return pub_kline
}

func ProcessNewStartKline(cur_kline *Kline, latest_kline *Kline, target_resolution int) {
	defer util.CatchExp("ProcessNewStartKline")

	latest_kline = NewKlineWithKline(cur_kline)
	latest_kline.Resolution = target_resolution
}

func ProcessCachingKline(cur_kline *Kline, latest_kline *Kline) {
	defer util.CatchExp("ProcessCachingKline")

	latest_kline.Close = cur_kline.Close
	latest_kline.Low = util.MinFloat64(latest_kline.Low, cur_kline.Low)
	latest_kline.High = util.MaxFloat64(latest_kline.High, cur_kline.High)
	latest_kline.Volume = latest_kline.Volume + cur_kline.Volume
}

func NewTreeMapWithKlines(ori_klines []*Kline, target_resolution int) *treemap.Map {
	defer util.CatchExp("NewTreeMapWithKlines")

	rst := treemap.NewWith(utils.Int64Comparator)

	var latest_kline *Kline = nil
	var pub_kline *Kline = nil

	for _, cur_kline := range ori_klines {
		if latest_kline == nil {
			latest_kline = cur_kline
			ResetFirstKline(latest_kline, target_resolution)
		}

		if IsOldKlineEndTime(cur_kline.Time, int(cur_kline.Resolution), int64(target_resolution)) {
			pub_kline = ProcessOldEndKline(cur_kline, latest_kline, int(target_resolution))
			rst.Put(pub_kline.Time, pub_kline)
		} else if IsNewKlineStartTime(cur_kline.Time, int64(target_resolution)) {
			ProcessNewStartKline(cur_kline, latest_kline, target_resolution)
		} else {
			ProcessCachingKline(cur_kline, latest_kline)
		}
	}

	// 不完整的 k 线数据;ex: 5min, 最后的时间只更新到 20:23;
	if latest_kline.Time != pub_kline.Time {
		rst.Put(latest_kline.Time, pub_kline)
	}

	return rst
}

// Undo
func ProcessOutdatedKline(kline_tree *treemap.Map, last_kline *Kline, new_kline *Kline, resolution int) {
	defer util.CatchExp("ProcessUpdatedKline")
	logx.Errorf("NewKLineTime: %d, CachedKlineTime: %d, kline already updated", new_kline.Time, last_kline.Time)
}

func GetLastKline(kline_tree *treemap.Map) *Kline {
	iter := kline_tree.Iterator()
	ok := iter.Last()

	if !ok {
		logx.Errorf("UpdateTreeWithKline Tree is Empty!")
		return nil
	}

	latest_kline := iter.Value().(*Kline)

	return latest_kline
}

/*
功能: 更新某个频率下的 K线树;
流程: 1. 获取已经存储最新的一根Kline - latest_kline;
     2. 若 new_kline 是当前频率的最后一分的数据，当前这个频率时段已经结束; ex: 5min, new_kline: 20:24 即为 20:00 ~ 20:24 的最后一根kline
		再次判断 new_kline 与 latest_kline 原始频率是否相等; ex: latest_kline: 5min, new_kline: 5min
		2.1 相等: 无需从 latest_kline 获取信息，直接发布;
		2.2 不相等: latest_kline: 5min, new_kline: 1min, [new_kline 频率一定小于 latest_kline - 用低频率聚合高频率]
				 需要进行行情的比较与吸收 - close, high, low, volume，再发布;
	3. 是否是下一个频率时段的起始时间；ex: 5min, new_kline: 20:25
	   将其存储进 Kline Tree; 可以直接发布;
	4. 是当前频率时段的中间某个时间的 Kline:
	   进行行情的比较与吸收 - close, high, low, volume，再发布;
*/
func UpdateTreeWithKline(latest_kline *Kline, kline *Kline, resolution int) (*Kline, bool) {
	defer util.CatchExp("UpdateTreeWithKline")

	var pub_kline *Kline = nil
	is_add := false

	if kline.Time <= latest_kline.Time {
		return nil, false
	}

	if IsOldKlineEnd(kline, int64(resolution)) {
		logx.Slowf("Old Kline End: rsl:%d, %s", resolution, kline.String())

		if kline.Resolution != resolution {
			latest_kline.Close = kline.Close
			latest_kline.Low = util.MinFloat64(latest_kline.Low, kline.Low)
			latest_kline.High = util.MaxFloat64(latest_kline.High, kline.High)
			latest_kline.Volume += kline.Volume

			pub_kline = NewKlineWithKline(latest_kline)
		} else {
			pub_kline = NewKlineWithKline(kline)
		}
	} else if IsNewKlineStart(kline, int64(resolution)) {
		logx.Slowf("New Kline Start: rsl:%d, %s", resolution, kline.String())
		new_add_kline := NewKlineWithKline(kline)
		new_add_kline.Resolution = resolution
		pub_kline = new_add_kline
		is_add = true
	} else {
		latest_kline.Close = kline.Close
		latest_kline.Low = util.MinFloat64(latest_kline.Low, kline.Low)
		latest_kline.High = util.MaxFloat64(latest_kline.High, kline.High)
		latest_kline.Volume = latest_kline.Volume + kline.Volume

		logx.Slowf("Cached Kline:%d, %s", resolution, kline.String())

		pub_kline = NewKlineWithKline(latest_kline)
	}

	return pub_kline, is_add
}

func UpdateTreeWithTrade(latest_kline *Kline, trade *Trade, resolution int) (*Kline, bool) {
	defer util.CatchExp("UpdateTreeWithKline")

	var pub_kline *Kline = nil
	is_add := false
	NextKlineTime := latest_kline.Time + int64(resolution)*NANO_PER_SECS

	if trade.Time/NANO_PER_MIN == NextKlineTime/NANO_PER_MIN {

		tmp_kline := &Kline{
			Exchange:   trade.Exchange,
			Symbol:     trade.Symbol,
			Time:       trade.Time - trade.Time%NANO_PER_MIN,
			Open:       trade.Price,
			High:       trade.Price,
			Low:        trade.Price,
			Close:      trade.Price,
			Volume:     trade.Volume,
			Resolution: resolution,
		}

		logx.Slowf("New Kline With: \nTrade %s\nkline: %s \n", trade.String(), tmp_kline.FullString())
		pub_kline = NewKlineWithKline(tmp_kline)
		is_add = true
	} else {

		if resolution == 60 {
			NextKlineTime = NextKlineTime + int64(resolution)*NANO_PER_SECS
		}

		if trade.Time <= latest_kline.Time {
			logx.Errorf("Trade.Time %s, earlier than CachedKlineTime: %s", util.TimeStrFromInt(trade.Time), util.TimeStrFromInt(latest_kline.Time))
			return nil, false
		}

		if trade.Time > NextKlineTime {
			logx.Errorf("Trade.Time %s, later than NextKlineTime: %s", util.TimeStrFromInt(trade.Time), util.TimeStrFromInt(NextKlineTime))
			return nil, false
		}

		latest_kline.Close = trade.Price
		if latest_kline.Low > trade.Price {
			latest_kline.Low = trade.Price

		}
		if latest_kline.High < trade.Price {
			latest_kline.High = trade.Price
		}

		latest_kline.Volume = latest_kline.Volume + trade.Volume

		pub_kline = NewKlineWithKline(latest_kline)
	}

	return pub_kline, is_add
}

// func UpdateTreeWithKlines(kline_tree *treemap.Map, ori_klines []*Kline, target_resolution int) *Kline {
// 	defer util.CatchExp("UpdateTreeWithKlines")

// 	var rst *Kline

// 	for _, kline := range ori_klines {
// 		rst = UpdateTreeWithKline(kline_tree, kline, target_resolution)
// 	}

// 	return rst
// }

//Undo
func TreeGetKlinesByCount(kline_tree *treemap.Map, resolution int, count int) []*Kline {
	defer util.CatchExp("TreeGetKlinesByCount")

	var rst []*Kline

	return rst
}

//Undo
func TreeGetKlinesByTime(kline_tree *treemap.Map, resolution int, start_time int64, end_time int64) []*Kline {
	defer util.CatchExp("TreeGetKlinesByCount")

	var rst []*Kline

	return rst
}
