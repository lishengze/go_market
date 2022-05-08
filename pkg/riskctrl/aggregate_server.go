package riskctrl

import (
	"sync"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
)

type aggregator struct {
	depth_cache map[string]map[string]*DepthQuote
	kline_cache map[string]map[string]*Kline
	trade_cache map[string]map[string]*Trade

	depth_mutex sync.Mutex
	kline_mutex sync.Mutex
	trade_mutex sync.Mutex

	depth_aggregator_millsecs time.Duration
}

func (a *aggregator) Start() {
	a.start_aggregate_depth()
}

func (a *aggregator) start_aggregate_depth() {

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

func mix_depth(src *treemap.Map, other *treemap.Map) {

}

func (a *aggregator) aggregate_depth() {
	a.depth_mutex.Lock()

	for symbol, exchange_depth_map := range a.depth_cache {
		new_depth := NewDepth(nil)
		new_depth.Symbol = symbol
		new_depth.Exchange = BCTS_EXCHANGE
		new_depth.Time = uint64(time.Now().Unix())

		for _, cur_depth := range exchange_depth_map {
			mix_depth(new_depth.Asks, cur_depth.Asks)
			mix_depth(new_depth.Bids, cur_depth.Bids)
		}

	}

	defer a.depth_mutex.Unlock()

}

func (a *aggregator) cache_data(data_chan *DataChannel) {
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

func (a *aggregator) cache_depth(depth *DepthQuote) {

	new_depth := NewDepth(depth)

	a.depth_mutex.Lock()

	a.depth_cache[new_depth.Symbol][new_depth.Exchange] = new_depth

	defer a.depth_mutex.Unlock()
}

func (a *aggregator) cache_kline(kline *Kline) {

}

func (a *aggregator) cache_trade(trade *Trade) {
	new_trade := NewTrade(trade)
	new_trade.Exchange = BCTS_EXCHANGE
	a.publish_trade(new_trade)
}

func (a *aggregator) publish_depth(depth *DepthQuote) {

}

func (a *aggregator) publish_kline(depth *Kline) {

}

func (a *aggregator) publish_trade(depth *Trade) {

}
