package datastruct

import (
	"fmt"
	"market_server/common/util"
	"strings"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
)

const (
	BCTS_EXCHANGE           = "_bcts_"
	BCTS_EXCHANGE_AGGREGATE = "_bcts_aggregate_"
	BINANCE                 = "BINANCE"
	FTX                     = "FTX"

	BCTS_GROUP = "BCTS"

	KLINE_TYPE = "kline"
	DEPTH_TYPE = "depth"
	TRADE_TYPE = "trade"
)

const (
	NANO_PER_MICR = 1000
	MICR_PER_MILL = 1000
	MILL_PER_SECS = 1000
	NANO_PER_MILL = NANO_PER_MICR * MICR_PER_MILL
	NANO_PER_SECS = NANO_PER_MICR * MICR_PER_MILL * MILL_PER_SECS
	SECS_PER_MIN  = 60
	MIN_PER_HOUR  = 60
	HOUR_PER_DAY  = 24

	NANO_PER_MIN  = NANO_PER_SECS * SECS_PER_MIN
	NANO_PER_HOUR = NANO_PER_MIN * MIN_PER_HOUR
	NANO_PER_DAY  = NANO_PER_HOUR * HOUR_PER_DAY

	SECS_PER_HOUR = SECS_PER_MIN * MIN_PER_HOUR
	SECS_PER_DAY  = SECS_PER_HOUR * HOUR_PER_DAY
	MILL_PER_DAY  = MILL_PER_SECS * SECS_PER_DAY

	MIN_PER_DAY = MIN_PER_HOUR * HOUR_PER_DAY
)

type TestData struct {
	Name string
}

type TSymbol string
type TExchange string

// type RFloat float64
type TPrice float64
type TVolume float64

type Metadata struct {
	DepthMeta map[string](map[string]struct{})
	KlineMeta map[string](map[string]struct{})
	TradeMeta map[string](map[string]struct{})
}

func get_meta_string(meta *map[string](map[string]struct{})) string {
	result := ""

	for symbol, exchange_set := range *meta {
		result += symbol + " \n"
		for exchange := range exchange_set {
			result += exchange + " \n"
		}
	}

	return result
}

func (m *Metadata) String() string {
	result := fmt.Sprintf("DepthMetaInfo: %s\nKlineMeta: %s\nTradeMeta: %s\n",
		get_meta_string(&m.DepthMeta), get_meta_string(&m.KlineMeta), get_meta_string(&m.TradeMeta))
	return result
}

type InnerDepth struct {
	Volume         float64
	ExchangeVolume map[string]float64
}

func (src *InnerDepth) Add(other *InnerDepth) {
	if src == other {
		return
	}

	src.Volume += other.Volume

	for exchange, volume := range other.ExchangeVolume {
		src.ExchangeVolume[exchange] += volume
	}
}

type DepthQuote struct {
	Exchange string
	Symbol   string
	Time     int64
	Asks     *treemap.Map
	Bids     *treemap.Map
	Sequence uint64
}

type Kline struct {
	Exchange   string
	Symbol     string
	Time       int64
	Open       float64
	High       float64
	Low        float64
	Close      float64
	Volume     float64
	Resolution int
	Sequence   uint64
	LastVolume float64
}

// UnTest
func (k *Kline) IsHistory() bool {

	if k.LastVolume <= 0 {
		return true
	} else {
		return false
	}
}

//UnTest
func (k *Kline) HasTrade() bool {

	if k.Volume > 0 {
		return true
	} else {
		return false
	}
}

func (k *Kline) IsInited() bool {
	if k.Time > 0 {
		return true
	} else {
		return false
	}
}

// UnTest
func (k *Kline) UpdateInfoByHistKline(new_kline *Kline) {
	defer util.CatchExp(fmt.Sprintf("UpdateInfoByHistKline %s", new_kline.FullString()))
	if k == new_kline {
		return
	}

	k.Close = new_kline.Close
	k.Low = util.MinFloat64(k.Low, new_kline.Low)
	k.High = util.MaxFloat64(k.High, new_kline.High)
	k.Volume = k.Volume + new_kline.Volume
	k.Sequence = new_kline.Sequence
}

// UnTest
func (k *Kline) UpdateInfoByRealKline(new_kline *Kline) {
	defer util.CatchExp(fmt.Sprintf("UpdateInfoByRealKline %s", new_kline.FullString()))
	if k == new_kline {
		return
	}

	k.Close = new_kline.Close
	k.Low = util.MinFloat64(k.Low, new_kline.Low)
	k.High = util.MaxFloat64(k.High, new_kline.High)
	k.Sequence = new_kline.Sequence
}

func (k *Kline) ResetWithNewKline(new_kline *Kline) {
	defer util.CatchExp(fmt.Sprintf("ResetWithNewKline %s", new_kline.FullString()))

	if k == new_kline {
		return
	}

	k.Symbol = new_kline.Symbol
	k.Exchange = new_kline.Exchange
	k.Time = new_kline.Time
	k.Open = new_kline.Open
	k.Close = new_kline.Close
	k.Low = new_kline.Low
	k.High = new_kline.High
	k.Volume = new_kline.Volume
	k.Resolution = new_kline.Resolution
	k.Sequence = new_kline.Sequence
	k.LastVolume = new_kline.LastVolume
}

func (k *Kline) SetHistoryFlag() {
	k.LastVolume = -1
}

func (k *Kline) RestWithLastPrice() {
	k.Open = k.Close
	k.High = k.Close
	k.Low = k.Close
	k.Volume = 0
	k.LastVolume = -1
	k.Time = 0
}

func (k *Kline) SetPerfectTime(resolution int64) {
	k.Time = k.Time - k.Time%resolution
}

type Trade struct {
	Exchange string
	Symbol   string
	Time     int64
	Price    float64
	Volume   float64
	Sequence uint64
}

type ReqTrade struct {
	Symbol          string
	ReqWSTime       int64
	ReqArriveTime   int64
	ReqResponseTime int64
}

type RspTrade struct {
	TradeData       *Trade
	ChangeData      *ChangeInfo
	UsdPrice        float64
	ReqWSTime       int64
	ReqArriveTime   int64
	ReqResponseTime int64
}

type ReqDepth struct {
	Symbol          string
	ReqArriveTime   int64
	ReqResponseTime int64
}

type ReqHistKline struct {
	Symbol    string
	Exchange  string
	StartTime uint64
	EndTime   uint64
	Count     uint32
	Frequency uint32

	ReqArriveTime   int64
	ReqResponseTime int64
}

func (r *ReqHistKline) String() string {
	var result string
	if r.StartTime == 0 {
		result = fmt.Sprintf("%s, start: %d, end: %d;resolution: %d", r.Symbol, r.StartTime, r.EndTime, r.Frequency)
	} else {
		result = fmt.Sprintf("%s, count: %d, resolution: %d", r.Symbol, r.Count, r.Frequency)
	}
	return result
}

type RspHistKline struct {
	ReqInfo        *ReqHistKline
	Klines         *treemap.Map
	IsLastComplete bool
}

type ChangeInfo struct {
	Symbol     string
	High       float64
	Low        float64
	Change     float64
	ChangeRate float64
}

func (r *RspHistKline) TimeList() string {
	rst := HistKlineTimeList(r.Klines, 0)
	return rst
}

func (r *RspHistKline) SimpleTimeList() string {
	rst := HistKlineSimpleTime(r.Klines)
	return rst
}

func (c *ChangeInfo) String() string {
	return fmt.Sprintf("Symbol: %s, High: %f, Low: %f, Change: %f, ChangeRate: %f \n",
		c.Symbol, c.High, c.Low, c.Change, c.ChangeRate)
}

func (t *Trade) String() string {
	res := fmt.Sprintf("%s.%s, %s, %d, v : %f, p: %f ", t.Exchange, t.Symbol,
		util.TimeStrFromInt(t.Time), t.Sequence, t.Volume, t.Price)
	return res
}

func (k *Kline) FullString() string {
	res := fmt.Sprintf("%s.%s, %s, %d, lv: %f, v: %f, o: %f, h: %f, l: %f, c: %f, %d",
		k.Exchange, k.Symbol, util.TimeStrFromInt(k.Time), k.Sequence, k.LastVolume, k.Volume,
		k.Open, k.High, k.Low, k.Close, k.Resolution)
	return res
}

func (k *Kline) String() string {
	res := fmt.Sprintf("%s.%s, %s, v: %f",
		k.Exchange, k.Symbol, util.TimeStrFromInt(k.Time), k.Volume)
	return res
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

func HistKlineSimpleTime(hist_line *treemap.Map) string {

	rst := fmt.Sprintf("Size: %d, ", hist_line.Size())

	iter := hist_line.Iterator()

	if iter.First() {
		rst = rst + fmt.Sprintf("First : %s ", util.TimeStrFromInt(iter.Value().(*Kline).Time))
	}

	if iter.Last() {
		rst = rst + fmt.Sprintf("Last : %s ", util.TimeStrFromInt(iter.Value().(*Kline).Time))
	}

	return rst
}

func HistKlineTimeList(hist_line *treemap.Map, size int) string {

	rst := fmt.Sprintf("Size: %d; \n", hist_line.Size())

	if size == 0 || size*2 > hist_line.Size() {
		iter := hist_line.Iterator()

		for iter.Begin(); iter.Next(); {
			rst = rst + fmt.Sprintf("%s, \n", util.TimeStrFromInt(iter.Value().(*Kline).Time))
		}
	} else {
		first_count := 0
		iter := hist_line.Iterator()

		rst = rst + fmt.Sprintf("First %d data: \n", size)
		for iter.Begin(); iter.Next() && first_count < size; {
			rst = rst + fmt.Sprintf("%s, \n", util.TimeStrFromInt(iter.Value().(*Kline).Time))
			first_count += 1
		}

		first_count = 0
		rst = rst + fmt.Sprintf("Last %d data: \n", size)
		for iter.End(); iter.Prev() && first_count < size; {
			rst = rst + fmt.Sprintf("%s, \n", util.TimeStrFromInt(iter.Value().(*Kline).Time))
			first_count += 1
		}
	}

	// if iter.First() {
	// 	rst = rst + fmt.Sprintf("First : %s ", iter.Value().(*Kline).String())
	// }

	// if iter.Last() {
	// 	rst = rst + fmt.Sprintf("Last : %s ", iter.Value().(*Kline).String())
	// }

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

func IsNewKlineStart(kline *Kline, resolution int64) bool {
	return IsNewKlineStartTime(kline.Time, resolution)

	// tmp_time := kline.Time

	// if resolution > NANO_PER_SECS {
	// 	resolution = resolution / NANO_PER_SECS
	// }

	// if tmp_time > NANO_PER_SECS {
	// 	tmp_time = tmp_time / NANO_PER_SECS
	// }

	// if tmp_time%resolution == 0 {
	// 	return true
	// } else {
	// 	return false
	// }
}

func IsTargetTime(time_secs int64, resolution_secs int64) bool {

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

func IsNewKlineStartTime(tmp_time int64, resolution int64) bool {

	if tmp_time > NANO_PER_HOUR {
		tmp_time = tmp_time - tmp_time%NANO_PER_MIN
		tmp_time = tmp_time / NANO_PER_SECS
	}

	if resolution > NANO_PER_SECS {
		resolution = resolution / NANO_PER_SECS
	}

	tmp_time = tmp_time - tmp_time%SECS_PER_MIN

	return IsTargetTime(tmp_time, resolution)

	// if tmp_time%resolution == 0 {
	// 	return true
	// } else {
	// 	return false
	// }
}

func IsOldKlineEnd(kline *Kline, resolution int64) bool {

	return IsOldKlineEndTime(kline.Time, kline.Resolution, resolution)

	// tmp_time := kline.Time
	// src_resolution := kline.Resolution

	// if resolution > NANO_PER_SECS {
	// 	resolution = resolution / NANO_PER_SECS
	// }

	// if tmp_time > NANO_PER_SECS {
	// 	tmp_time = tmp_time / NANO_PER_SECS
	// }

	// if src_resolution > NANO_PER_SECS {
	// 	src_resolution = src_resolution / NANO_PER_SECS
	// }

	// tmp_time = tmp_time + int64(src_resolution)

	// if tmp_time%resolution == 0 {
	// 	// logx.Slowf("OldKlineEnd: kline: %s, resolution: %d, tmp_time: %d, ", kline.String(), resolution, tmp_time)
	// 	return true
	// } else {
	// 	return false
	// }
}

func IsOldKlineEndTime(tmp_time int64, src_resolution int, dst_resolution int64) bool {

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

	return IsTargetTime(tmp_time, dst_resolution)

	// if tmp_time%dst_resolution == 0 {
	// 	return true
	// } else {
	// 	return false
	// }
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

func GetDepthString(m *treemap.Map, numb int) string {
	str := "TreeMap\nmap["
	it := m.Iterator()
	it.Begin()
	count := 0
	for it.Next() {
		str += fmt.Sprintf("%v:%v ", it.Key(), it.Value())
		count += 1
		if count > numb {
			break
		}
	}
	return strings.TrimRight(str, " ") + "]"

}

func (d *DepthQuote) String(len int) string {

	res := fmt.Sprintf("%s.%s, %v\nAsks.%d: %s\nBids.%d: %s \n", d.Exchange, d.Symbol,
		time.Unix(int64(d.Time/NANO_PER_SECS), d.Time%NANO_PER_SECS),
		d.Asks.Size(), GetDepthString(d.Asks, len), d.Bids.Size(), GetDepthString(d.Bids, len))

	return string(res)
}

func DepthListString(m *treemap.Map, numb int) string {
	str := ""
	it := m.Iterator()
	it.Begin()
	count := 0
	for it.Next() {
		str += fmt.Sprintf("%v:%v \n", it.Key(), it.Value())
		count += 1
		if count > numb {
			break
		}
	}
	return str
}

func NewTrade(src *Trade) *Trade {
	if src != nil {
		rst := &Trade{
			Exchange: src.Exchange,
			Symbol:   src.Symbol,
			Time:     src.Time,
			Price:    src.Price,
			Volume:   src.Volume,
		}
		return rst
	} else {
		rst := &Trade{}
		return rst
	}
}

func NewKline(src *Kline) *Kline {
	if src != nil {
		rst := &Kline{
			Exchange:   src.Exchange,
			Symbol:     src.Symbol,
			Time:       src.Time,
			Open:       src.Open,
			High:       src.High,
			Low:        src.Low,
			Close:      src.Close,
			Volume:     src.Volume,
			Resolution: src.Resolution,
			LastVolume: src.LastVolume,
			Sequence:   src.Sequence,
		}
		return rst
	} else {
		rst := &Kline{}
		return rst
	}
}

func InitKlineByTrade(src *Kline, trade *Trade) {
	src.Exchange = BCTS_EXCHANGE
	src.Symbol = trade.Symbol
	src.Time = trade.Time
	src.Resolution = SECS_PER_MIN
	src.Open = trade.Price
	src.High = trade.Price
	src.Low = trade.Price
	src.Close = trade.Price
	src.Volume = trade.Volume
	src.LastVolume = trade.Volume
}

func NewDepth(src *DepthQuote) *DepthQuote {
	if src != nil {
		rst := &DepthQuote{
			Exchange: src.Exchange,
			Symbol:   src.Symbol,
			Time:     src.Time,
			Asks:     treemap.NewWith(utils.Float64Comparator),
			Bids:     treemap.NewWith(util.Float64ComparatorDsc),
		}

		ask_iter := src.Asks.Iterator()
		for ask_iter.Begin(); ask_iter.Next(); {
			rst.Asks.Put(ask_iter.Key(), ask_iter.Value())
		}

		bid_iter := src.Bids.Iterator()
		for bid_iter.Begin(); bid_iter.Next(); {
			rst.Bids.Put(bid_iter.Key(), bid_iter.Value())
		}

		return rst
	} else {
		rst := &DepthQuote{
			Asks: treemap.NewWith(utils.Float64Comparator),
			Bids: treemap.NewWith(util.Float64ComparatorDsc),
		}
		return rst
	}
}

type DataChannel struct {
	TradeChannel chan *Trade
	KlineChannel chan *Kline
	DepthChannel chan *DepthQuote
}

func NewDataChannel() *DataChannel {
	return &DataChannel{
		DepthChannel: make(chan *DepthQuote),
		KlineChannel: make(chan *Kline),
		TradeChannel: make(chan *Trade),
	}
}

func (d *DepthQuote) Init() {
	d.Asks = treemap.NewWith(utils.Float64Comparator)
	d.Bids = treemap.NewWith(utils.Float64Comparator)
}
