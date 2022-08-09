package datastruct

import (
	"fmt"
	"market_server/common/util"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

//Undo
//UnTest
func NewTradeWithRealTimeKline(kline *Kline) *Trade {
	return &Trade{
		Exchange: kline.Exchange,
		Symbol:   kline.Symbol,
		Time:     kline.Time,
		Price:    kline.Close,
		Volume:   kline.LastVolume,
		Sequence: kline.Sequence,
	}
}

func TransSliceKlines(ori_klines []*Kline) *treemap.Map {
	var rst *treemap.Map

	for _, kline := range ori_klines {
		rst.Put(kline.Time, kline)
	}

	return rst
}

func HistKlineString(hist_line *treemap.Map) string {

	rst := fmt.Sprintf("HistKline, Size: %d, ", hist_line.Size())

	iter := hist_line.Iterator()

	if iter.First() {
		rst = rst + fmt.Sprintf("First : %s ", iter.Value().(*Kline).String())
	}

	if iter.Last() {
		rst = rst + fmt.Sprintf("Last : %s ", iter.Value().(*Kline).String())
	}

	return rst
}

func HistKlineSimpleTime(klines []*Kline) string {

	if len(klines) == 0 {
		return "klines is empty!"
	}

	return fmt.Sprintf("Size: %d, First : %s, Last: %s ", len(klines), klines[0], klines[len(klines)-1])
}

func HistKlineTimeList(klines []*Kline, size int) string {

	rst := fmt.Sprintf("Size: %d; \n", len(klines))

	if size == 0 || size*2 > len(klines) {
		for _, kline := range klines {
			rst = rst + fmt.Sprintf("%s, \n", util.TimeStrFromInt(kline.Time))
		}
	} else {
		first_data := klines[0:size]
		rst = rst + fmt.Sprintf("First %d data: \n", size)
		for _, kline := range first_data {
			rst = rst + fmt.Sprintf("%s, \n", util.TimeStrFromInt(kline.Time))
		}

		end_data := klines[len(klines)-size:]

		rst = rst + fmt.Sprintf("Last %d data: \n", size)
		for _, kline := range end_data {
			rst = rst + fmt.Sprintf("%s, \n", util.TimeStrFromInt(kline.Time))
		}

	}
	return rst
}

func HistKlineList(hist_line *treemap.Map, size int) string {

	rst := fmt.Sprintf("Size: %d; \n", hist_line.Size())

	if size == 0 || size*2 > hist_line.Size() {
		iter := hist_line.Iterator()

		for iter.Begin(); iter.Next(); {
			rst = rst + fmt.Sprintf("%s, \n", iter.Value().(*Kline).FullString())
		}
	} else {
		first_count := 0
		iter := hist_line.Iterator()

		rst = rst + fmt.Sprintf("First %d data: \n", size)
		for iter.Begin(); iter.Next() && first_count < size; {
			rst = rst + fmt.Sprintf("%s, \n", iter.Value().(*Kline).FullString())
			first_count += 1
		}

		first_count = 0
		rst = rst + fmt.Sprintf("Last %d data: \n", size)
		for iter.End(); iter.Prev() && first_count < size; {
			rst = rst + fmt.Sprintf("%s, \n", iter.Value().(*Kline).FullString())
			first_count += 1
		}
	}
	return rst
}

func HistKlineListWithSlice(klines []*Kline, size int) string {

	rst := fmt.Sprintf("Size: %d; \n", len(klines))

	if size == 0 || size*2 > len(klines) {
		for _, kline := range klines {
			rst = rst + fmt.Sprintf("%s, \n", kline.FullString())
		}
	} else {
		first_data := klines[0:size]
		rst = rst + fmt.Sprintf("First %d data: \n", size)
		for _, kline := range first_data {
			rst = rst + fmt.Sprintf("%s, \n", kline.FullString())
		}

		end_data := klines[len(klines)-size:]

		rst = rst + fmt.Sprintf("Last %d data: \n", size)
		for _, kline := range end_data {
			rst = rst + fmt.Sprintf("%s, \n", kline.FullString())
		}

	}
	return rst
}

func NewKlineWithKline(kline *Kline) *Kline {
	return &Kline{
		Exchange:   kline.Exchange,
		Symbol:     kline.Symbol,
		Time:       kline.Time,
		Open:       kline.Open,
		High:       kline.High,
		Low:        kline.Low,
		Close:      kline.Close,
		Volume:     kline.Volume,
		Resolution: kline.Resolution,
		LastVolume: kline.LastVolume,
		Sequence:   kline.Sequence,
	}
}

func IsNewKlineStart(kline *Kline, resolution uint64) bool {
	return IsNewKlineStartTime(kline.Time, resolution)
}

func IsTargetTime(time_secs uint64, resolution_secs uint64) bool {

	if resolution_secs == SECS_PER_DAY*7 {
		// fmt.Printf("Ori Days: %d\n", time_secs/SECS_PER_DAY)

		time_secs = time_secs + SECS_PER_DAY*3

		// fmt.Printf("Tras Days: %d\n", time_secs/SECS_PER_DAY)
	}

	// fmt.Printf("time_secs: %d \n", time_secs)

	if time_secs%resolution_secs == 0 {
		return true
	} else {
		return false
	}
}

func IsNewKlineStartTime(tmp_time int64, resolution uint64) bool {

	if tmp_time > NANO_PER_HOUR {
		tmp_time = tmp_time - tmp_time%NANO_PER_MIN
		tmp_time = tmp_time / NANO_PER_SECS
	}

	if resolution > NANO_PER_SECS {
		resolution = resolution / NANO_PER_SECS
	}

	tmp_time = tmp_time - tmp_time%SECS_PER_MIN

	return IsTargetTime(uint64(tmp_time), resolution)

}

func IsOldKlineEnd(kline *Kline, resolution uint64) bool {

	return IsOldKlineEndTime(kline.Time, kline.Resolution, resolution)
}

func IsOldKlineEndTime(tmp_time int64, src_resolution uint64, dst_resolution uint64) bool {

	if tmp_time > NANO_PER_HOUR {
		tmp_time = tmp_time - tmp_time%NANO_PER_MIN
		tmp_time = tmp_time / NANO_PER_SECS
	}

	if dst_resolution > NANO_PER_SECS {
		dst_resolution = dst_resolution / NANO_PER_SECS
	}

	if src_resolution > NANO_PER_SECS {
		src_resolution = src_resolution / NANO_PER_SECS
	}

	tmp_time = tmp_time - tmp_time%SECS_PER_MIN

	tmp_time = tmp_time + int64(src_resolution)

	return IsTargetTime(uint64(tmp_time), dst_resolution)
}

func GetLastStartTime(tmp_time int64, resolution int64) int64 {
	if tmp_time > NANO_PER_DAY {
		tmp_time = tmp_time / NANO_PER_SECS
	}

	if resolution > NANO_PER_SECS {
		resolution = resolution / NANO_PER_SECS
	}

	// fmt.Printf("tmp_time days : %d \n", tmp_time/SECS_PER_DAY)

	original_time := tmp_time

	if resolution == SECS_PER_DAY*7 {

		tmp_time = tmp_time - tmp_time%SECS_PER_DAY
		original_time = tmp_time
		tmp_time = tmp_time + SECS_PER_DAY*3
	}

	// fmt.Printf("trans_time days : %d , res: %d\n", tmp_time/SECS_PER_DAY, resolution/SECS_PER_DAY)
	// fmt.Printf("trans_time: %s, left_days: %d \n", util.TimeStrFromInt(tmp_time*NANO_PER_SECS), tmp_time%(resolution)/SECS_PER_DAY)

	tmp_time = original_time - tmp_time%(resolution)

	tmp_time = tmp_time * NANO_PER_SECS

	return tmp_time
}

func GetMiniuteNanos(ori_nanos int64) int64 {
	return ori_nanos - ori_nanos%NANO_PER_MIN
}

func ResetFirstKline(latest_kline *Kline, target_resolution uint64) {
	defer util.CatchExp("ResetFirstKline")
	latest_kline.Time = GetLastStartTime(latest_kline.Time, int64(target_resolution))
	latest_kline.Resolution = target_resolution
}

func ProcessOldEndKline(cur_kline *Kline, latest_kline *Kline, target_resolution uint64) *Kline {
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

func ProcessNewStartKline(cur_kline *Kline, latest_kline *Kline, target_resolution uint64) {
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

func NewTreeMapWithKlines(ori_klines []*Kline, target_resolution uint64) *treemap.Map {
	defer util.CatchExp("NewTreeMapWithKlines")

	rst := treemap.NewWith(utils.Int64Comparator)

	var latest_kline *Kline = nil
	var pub_kline *Kline = nil

	for _, cur_kline := range ori_klines {
		if latest_kline == nil {
			latest_kline = cur_kline
			ResetFirstKline(latest_kline, target_resolution)
		}

		if IsOldKlineEndTime(cur_kline.Time, cur_kline.Resolution, target_resolution) {
			pub_kline = ProcessOldEndKline(cur_kline, latest_kline, target_resolution)
			rst.Put(pub_kline.Time, pub_kline)
		} else if IsNewKlineStartTime(cur_kline.Time, target_resolution) {
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

func OutputDetailHistKlines(klines []*Kline) {
	defer util.CatchExp("OutputDetailHistKlines")

	if len(klines) > 3 {
		pub_kline := klines[len(klines)-3]
		cache_kline := klines[len(klines)-2]
		last_kline := klines[len(klines)-1]

		logx.Slowf("\npub_kline: %s\ncache_kline:%s\nlast_kline:%s\n",
			pub_kline.FullString(), cache_kline.FullString(), last_kline.FullString())
	}

	for _, kline := range klines {
		logx.Slowf("%s", kline.FullString())
	}
}

type TestKlineTool struct {
}

func (t *TestKlineTool) TestNewTradeWithRealTimeKline() {
	ori_klines := GetTestKline()
	trans_trade := NewTradeWithRealTimeKline(ori_klines)
	fmt.Printf("OriKline: %s\nTransTrade: %s\n", ori_klines.FullString(), trans_trade.String())
}

func TestKlineToolMain() {
	test_obj := &TestKlineTool{}
	test_obj.TestNewTradeWithRealTimeKline()
}
