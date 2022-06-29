/*
接受数据的接口;
*/

package data_engine

import (
	"fmt"
	"market_server/app/data_manager/rpc/marketservice"
	"market_server/common/datastruct"
	"sync"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/zeromicro/go-zero/core/logx"
)

type AtomData struct {
	price float64
	time  int64
}

type SortedList struct {
	list  []*AtomData
	IsAsc bool
}

func (s *SortedList) Size() int {
	return len(s.list)
}

func NewSortedList(is_asc bool) *SortedList {
	return &SortedList{
		IsAsc: is_asc,
	}
}

func (s *SortedList) Add(a *AtomData) {
	s.list = append(s.list, a)
	data_count := len(s.list)

	if s.IsAsc {
		for i := data_count - 1; i > 0; i-- {
			if s.list[i-1].price > s.list[i].price {
				tmp := s.list[i]
				s.list[i] = s.list[i-1]
				s.list[i-1] = tmp
			}
		}
	} else {
		for i := data_count - 1; i > 0; i-- {
			if s.list[i-1].price < s.list[i].price {
				tmp := s.list[i]
				s.list[i] = s.list[i-1]
				s.list[i-1] = tmp
			}
		}
	}
}

// seq = append(seq[:index], seq[index+1:]...)

func (s *SortedList) Del(a *AtomData) {
	index := -1
	for i := len(s.list) - 1; i > 0; i-- {
		if a.price == s.list[i].price && a.time == s.list[i].time {
			index = i
		}
	}

	if index != -1 {
		s.list = append(s.list[:index], s.list[index+1:]...)
	}
}

func (s *SortedList) Last() *AtomData {
	return s.list[len(s.list)-1]
}

type PeriodData struct {
	time_cache_data       *treemap.Map
	high_price_cache_data *SortedList
	low_price_cache_data  *SortedList

	Symbol  string
	Max     float64
	MaxTime int64

	Min     float64
	MinTime int64

	Start     float64
	StartTime int64

	Change float64

	TimeNanos int64
	Count     int

	CurTrade *datastruct.Trade

	mutex sync.Mutex
}

func (p *PeriodData) String() string {
	return fmt.Sprintf("%s, Max: %f, MaxTime: %+v, Min: %f, MinTime: %+v, Change: %f", p.Symbol,
		p.Max, time.Unix(int64(p.MaxTime/datastruct.NANO_PER_SECS), p.MaxTime%datastruct.NANO_PER_SECS),
		p.Min, time.Unix(int64(p.MinTime/datastruct.NANO_PER_SECS), p.MinTime%datastruct.NANO_PER_SECS),
		p.Change)
}

func NewPeriodData() *PeriodData {
	return nil
}

func (p *PeriodData) UpdateWithTrade(trade *datastruct.Trade) {

	p.UpdateMeta()
}

func (p *PeriodData) UpdateWithKline(kline *datastruct.Kline) {
	p.AddCacheData(kline)

	p.EraseOuttimeData()

	p.UpdateMeta()
}

func (p *PeriodData) AddTradeData(trade *datastruct.Trade) {
	p.mutex.Lock()

	p.CurTrade = trade

	defer p.mutex.Unlock()
}

func (p *PeriodData) AddCacheData(kline *datastruct.Kline) {
	p.mutex.Lock()

	logx.Statf("[Add] Kline: %s", kline.String())

	p.time_cache_data.Put(kline.Time, kline)

	p.high_price_cache_data.Add(&AtomData{
		time:  kline.Time,
		price: kline.High,
	})

	p.low_price_cache_data.Add(&AtomData{
		time:  kline.Time,
		price: kline.Low,
	})

	defer p.mutex.Unlock()
}

func (p *PeriodData) EraseOuttimeData() {

	p.mutex.Lock()

	defer p.mutex.Unlock()

	// type outtime_data struct {
	// 	time int64
	// 	high float64
	// 	low  float64
	// }

	var outtime_datalist []*datastruct.Kline

	begin_iter := p.time_cache_data.Iterator()
	if ok := begin_iter.First(); !ok {
		return
	}

	last_iter := p.time_cache_data.Iterator()
	if ok := last_iter.Last(); !ok {
		return
	}

	for begin_iter.Next() {
		if last_iter.Key().(int64)-begin_iter.Key().(int64) > p.TimeNanos {
			outtime_datalist = append(outtime_datalist, begin_iter.Value().(*datastruct.Kline))
		} else {
			break
		}
	}

	for _, outtime := range outtime_datalist {
		logx.Statf("[Erase] %s ", outtime.String())
		p.time_cache_data.Remove(outtime.Time)

		p.high_price_cache_data.Del(&AtomData{
			price: outtime.High,
			time:  outtime.Time})

		p.low_price_cache_data.Del(&AtomData{
			price: outtime.Low,
			time:  outtime.Time})
	}

}

func (p *PeriodData) InitCacheData(klines *marketservice.HistKlineData) {

	p.mutex.Lock()

	defer p.mutex.Unlock()

	for _, pb_kline := range klines.KlineData {
		kline := marketservice.NewKlineWithPbKline(pb_kline)
		if kline == nil {
			continue
		}

		p.time_cache_data.Put(kline.Time, kline)

		p.high_price_cache_data.Add(&AtomData{
			price: kline.High,
			time:  kline.Time})

		p.low_price_cache_data.Add(&AtomData{
			price: kline.Low,
			time:  kline.Time})
	}

	iter := p.time_cache_data.Iterator()
	if iter.First() {
		logx.Statf("[Init]First : %s ", iter.Value().(*datastruct.Kline).String())
	}

	if iter.Last() {
		logx.Statf("[Init]Last : %s ", iter.Value().(*datastruct.Kline).String())
	}

}

func (p *PeriodData) InitCacheDataWithTreeMap(klines *treemap.Map) {

	p.mutex.Lock()

	defer p.mutex.Unlock()

	iter := klines.Iterator()

	for iter.Begin(); iter.Next(); {
		kline := iter.Value().(*datastruct.Kline)

		p.time_cache_data.Put(kline.Time, kline)

		p.high_price_cache_data.Add(&AtomData{
			price: kline.High,
			time:  kline.Time})

		p.low_price_cache_data.Add(&AtomData{
			price: kline.Low,
			time:  kline.Time})
	}

	if iter.First() {
		logx.Statf("[Init] First: %s ", iter.Value().(*datastruct.Kline).String())
	}

	if iter.Last() {
		logx.Statf("[Init] Last: %s ", iter.Value().(*datastruct.Kline).String())
	}
}

func (p *PeriodData) UpdateMeta() {
	p.mutex.Lock()

	defer p.mutex.Unlock()

	if p.time_cache_data.Size() > 0 {
		first := p.time_cache_data.Iterator()
		first.First()
		p.Start = first.Value().(*datastruct.Kline).Open
		p.StartTime = first.Key().(int64)
	} else {
		return
	}

	if p.high_price_cache_data.Size() > 0 {
		p.Max = p.high_price_cache_data.Last().price
		p.MaxTime = p.high_price_cache_data.Last().time

		if p.CurTrade != nil && p.CurTrade.Time > p.StartTime && p.Max < p.CurTrade.Price {
			p.Max = p.CurTrade.Price
			p.MaxTime = p.CurTrade.Time
		}
	} else {
		return
	}

	if p.low_price_cache_data.Size() > 0 {
		p.Min = p.low_price_cache_data.Last().price
		p.MinTime = p.low_price_cache_data.Last().time

		if p.CurTrade != nil && p.CurTrade.Time > p.StartTime && p.Min > p.CurTrade.Price {
			p.Min = p.CurTrade.Price
			p.MinTime = p.CurTrade.Time
		}
	} else {
		return
	}

	p.Change = (p.Max - p.Start) / p.Start
	p.Count = p.time_cache_data.Size()

	// logx.Statf("[Meta] %s", p.String())
}

func (p *PeriodData) UpdateWithPbKlines(klines *marketservice.HistKlineData) {
	p.InitCacheData(klines)

	p.EraseOuttimeData()

	p.UpdateMeta()
}

func (p *PeriodData) UpdateWithKlines(klines *treemap.Map) {
	p.InitCacheDataWithTreeMap(klines)

	p.EraseOuttimeData()

	p.UpdateMeta()
}

func (p *PeriodData) GetChangeInfo() *datastruct.ChangeInfo {
	p.mutex.Lock()

	defer p.mutex.Unlock()

	return &datastruct.ChangeInfo{
		Symbol: p.Symbol,
		High:   p.Max,
		Low:    p.Min,
		Change: p.Change,
	}
}
