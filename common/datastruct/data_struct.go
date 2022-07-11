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
}

type Trade struct {
	Exchange string
	Symbol   string
	Time     int64
	Price    float64
	Volume   float64
}

type ReqHistKline struct {
	Symbol    string
	Exchange  string
	StartTime uint64
	EndTime   uint64
	Count     uint32
	Frequency uint32
}

func (r *ReqHistKline) String() string {
	result := fmt.Sprintf("%s, count: %d, resolution: %d", r.Symbol, r.Count, r.Frequency)
	return result
}

type RspHistKline struct {
	ReqInfo *ReqHistKline
	Klines  *treemap.Map
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
	res := fmt.Sprintf("%s.%s, %+v, p: %f v: %f \n", t.Exchange, t.Symbol,
		time.Unix(int64(t.Time/NANO_PER_SECS), t.Time%NANO_PER_SECS), t.Price, t.Volume)
	return res
}

func (k *Kline) String() string {
	res := fmt.Sprintf("%s.%s, %+v, o: %f, h: %f, l: %f, c: %f, v: %f\n",
		k.Exchange, k.Symbol, time.Unix(int64(k.Time/NANO_PER_SECS), k.Time%NANO_PER_SECS),
		k.Open, k.High, k.Low, k.Close, k.Volume)
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

	rst := fmt.Sprintf("HistKline, Size: %d, ", hist_line.Size())

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
	}
}

func IsNewKlineStart(kline *Kline, resolution int64) bool {
	tmp_time := kline.Time

	if resolution > NANO_PER_SECS {
		resolution = resolution / NANO_PER_SECS
	}

	if tmp_time > NANO_PER_SECS {
		tmp_time = tmp_time / NANO_PER_SECS
	}

	if tmp_time%resolution == 0 {
		return true
	} else {
		return false
	}
}

func IsOldKlineEnd(kline *Kline, resolution int64) bool {
	tmp_time := kline.Time
	tmp_resolution := kline.Resolution

	if resolution > NANO_PER_SECS {
		resolution = resolution / NANO_PER_SECS
	}

	if tmp_time > NANO_PER_SECS {
		tmp_time = tmp_time / NANO_PER_SECS
	}

	if tmp_resolution > NANO_PER_SECS {
		tmp_resolution = tmp_resolution / NANO_PER_SECS
	}

	tmp_time = tmp_time + int64(tmp_resolution)

	if tmp_time%resolution == 0 {
		// logx.Slowf("OldKlineEnd: kline: %s, resolution: %d, tmp_time: %d, ", kline.String(), resolution, tmp_time)
		return true
	} else {
		return false
	}
}

func GetDepthString(m *treemap.Map, numb int) string {
	str := "TreeMap\nmap["
	it := m.Iterator()
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
	src.Resolution = 60
	src.Open = trade.Price
	src.High = trade.Price
	src.Low = trade.Price
	src.Close = trade.Price
	src.Volume = trade.Volume
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
