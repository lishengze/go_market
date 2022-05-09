package riskctrl

import (
	"fmt"
	"sync"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
)

type Aggregator struct {
	depth_cache map[string]map[string]*DepthQuote
	kline_cache map[string]map[string]*Kline
	trade_cache map[string]map[string]*Trade

	depth_mutex sync.Mutex
	kline_mutex sync.Mutex
	trade_mutex sync.Mutex

	depth_aggregator_millsecs time.Duration
}

func (a *Aggregator) Start(data_chan *DataChannel) {
	a.start_receiver(data_chan)
	a.start_aggregate_depth()

}

func (a *Aggregator) start_aggregate_depth() {

	go func() {
		LOG_INFO("start_aggregate_depth")
		timer := time.Tick(time.Duration(a.depth_aggregator_millsecs * time.Millisecond))

		for {
			select {
			case <-timer:
				a.aggregate_depth()
			}
		}
	}()
}

func mix_depth(src *treemap.Map, other *treemap.Map, exchange string) {
	other_iter := other.Iterator()

	for other_iter.Begin(); other_iter.Next(); {
		if cur_iter, ok := src.Get(other_iter.Key()); ok {
			cur_innerdepth := cur_iter.(*InnerDepth)
			cur_innerdepth.Volume += other_iter.Value().(*InnerDepth).Volume
			cur_innerdepth.ExchangeVolume[exchange] = other_iter.Value().(*InnerDepth).Volume
		} else {
			src.Put(other_iter.Key(), other_iter.Value())
		}
	}
}

func (a *Aggregator) aggregate_depth() {
	a.depth_mutex.Lock()

	LOG_INFO("aggregate depth")

	for symbol, exchange_depth_map := range a.depth_cache {
		new_depth := NewDepth(nil)
		new_depth.Symbol = symbol
		new_depth.Exchange = BCTS_EXCHANGE
		new_depth.Time = uint64(time.Now().Unix())

		for exchange, cur_depth := range exchange_depth_map {
			mix_depth(new_depth.Asks, cur_depth.Asks, exchange)
			mix_depth(new_depth.Bids, cur_depth.Bids, exchange)
		}

		a.publish_depth(new_depth)
	}

	defer a.depth_mutex.Unlock()

}

func (a *Aggregator) start_receiver(data_chan *DataChannel) {
	LOG_INFO("Aggregator start_receiver")
	for {
		select {
		case new_depth := <-data_chan.DepthChannel:
			a.cache_depth(new_depth)
		case new_kline := <-data_chan.KlineChannel:
			a.cache_kline(new_kline)
		case new_trade := <-data_chan.TradeChannel:
			a.cache_trade(new_trade)
		}
	}
}

func (a *Aggregator) cache_depth(depth *DepthQuote) {

	new_depth := NewDepth(depth)

	a.depth_mutex.Lock()

	a.depth_cache[new_depth.Symbol][new_depth.Exchange] = new_depth

	defer a.depth_mutex.Unlock()
}

func (a *Aggregator) cache_kline(kline *Kline) {

}

func (a *Aggregator) cache_trade(trade *Trade) {
	new_trade := NewTrade(trade)
	new_trade.Exchange = BCTS_EXCHANGE
	a.publish_trade(new_trade)
}

func (a *Aggregator) publish_depth(depth *DepthQuote) {
	LOG_INFO("publish_depth: \n" + depth.String(5))
}

func (a *Aggregator) publish_kline(kline *Kline) {
	fmt.Printf("Publish Kline: \n%+v\n", kline)
}

func (a *Aggregator) publish_trade(trade *Trade) {
	fmt.Printf("Publish Trade: \n%+v\n", trade)
}

func PublishTest(data *DataChannel) {
	timer := time.Tick(3 * time.Second)

	index := 0
	for {
		select {
		case <-timer:
			GetTestDepth()
		}
	}
}

func TestAggregator() {
	aggregator := Aggregator{
		depth_aggregator_millsecs: 5000,
	}

	data_chan := DataChannel{
		DepthChannel: make(chan *DepthQuote),
		KlineChannel: make(chan *Kline),
		TradeChannel: make(chan *Trade),
	}
	aggregator.Start(&data_chan)

	go PublishTest(&data_chan)
}
