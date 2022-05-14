package riskctrl

import (
	"fmt"
	"market_aggregate/pkg/datastruct"
	"math/rand"
	"sync"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
)

type Aggregator struct {
	depth_cache map[string]map[string]*datastruct.DepthQuote

	kline_cache map[string]map[string]*datastruct.Kline
	trade_cache map[string]map[string]*datastruct.Trade

	kline_aggregated map[string]*datastruct.Kline

	depth_mutex sync.Mutex
	kline_mutex sync.Mutex
	trade_mutex sync.Mutex

	depth_aggregator_millsecs time.Duration
}

func (a *Aggregator) Start(data_chan *datastruct.DataChannel) {
	a.start_receiver(data_chan)
	a.start_aggregate_depth()
	a.start_aggregate_kline()
}

func (a *Aggregator) start_aggregate_depth() {
	LOG_INFO("Aggregator start_aggregate_depth!")
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
	LOG_INFO("Aggregator start_aggregate_depth Over!")
}

func mix_depth(src *treemap.Map, other *treemap.Map, exchange string) {
	other_iter := other.Iterator()

	for other_iter.Begin(); other_iter.Next(); {
		if cur_iter, ok := src.Get(other_iter.Key()); ok {
			cur_innerdepth := cur_iter.(*datastruct.InnerDepth)
			cur_innerdepth.Volume += other_iter.Value().(*datastruct.InnerDepth).Volume
			cur_innerdepth.ExchangeVolume[exchange] = other_iter.Value().(*datastruct.InnerDepth).Volume
		} else {
			src_inner_depth := datastruct.InnerDepth{0, make(map[string]float64)}

			other_inner_depth := other_iter.Value().(*datastruct.InnerDepth)
			src_inner_depth.Volume = other_inner_depth.Volume
			src_inner_depth.ExchangeVolume[exchange] = other_inner_depth.Volume
			price := other_iter.Key().(float64)

			src.Put(price, &src_inner_depth)
		}
	}
}

func (a *Aggregator) start_aggregate_kline() {
	LOG_INFO("Aggregate datastruct.Kline Start!")
	go func() {
		for {
			WaitForNextMinute()

			a.aggregate_kline()
		}
	}()
	LOG_INFO("Aggregate datastruct.Kline Over!")
}

func (a *Aggregator) aggregate_depth() {
	a.depth_mutex.Lock()

	// LOG_INFO("----- Aggregate Depth Start ------ ")

	for symbol, exchange_depth_map := range a.depth_cache {
		new_depth := datastruct.NewDepth(nil)
		new_depth.Symbol = symbol
		new_depth.Exchange = datastruct.BCTS_EXCHANGE
		new_depth.Time = time.Now().Unix()

		for exchange, cur_depth := range exchange_depth_map {
			LOG_INFO("\n===== <<CurDepth>>: " + cur_depth.String(5))
			mix_depth(new_depth.Asks, cur_depth.Asks, exchange)
			mix_depth(new_depth.Bids, cur_depth.Bids, exchange)
		}

		LOG_INFO("\n^^^^^^^ <<aagregated_depth>>: " + new_depth.String(5))
		a.publish_depth(new_depth)

		// for _, cur_depth := range exchange_depth_map {
		// 	LOG_INFO("\n===== After <<CurDepth>>: " + cur_depth.String(5))
		// }
	}

	// LOG_INFO("----- Aggregate Depth Over!------ \n")

	defer a.depth_mutex.Unlock()

}

func (a *Aggregator) start_receiver(data_chan *datastruct.DataChannel) {
	LOG_INFO("Aggregator start_receiver")
	go func() {
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
	}()
	LOG_INFO("Aggregator start_receiver Over!")
}

func (a *Aggregator) aggregate_kline() {
	defer a.kline_mutex.Unlock()
	a.kline_mutex.Lock()

	for _, kline := range a.kline_aggregated {
		new_kline := datastruct.NewKline(kline)
		kline.Time = 0
		new_kline.Time = TimeMinute()
		a.publish_kline(new_kline)
	}
}

func (a *Aggregator) update_kline(trade *datastruct.Trade) {
	defer a.kline_mutex.Unlock()
	a.kline_mutex.Lock()

	if _, ok := a.kline_aggregated[trade.Symbol]; ok == false {
		a.kline_aggregated[trade.Symbol] = datastruct.NewKline(nil)
	}

	cur_kline := a.kline_aggregated[trade.Symbol]
	if cur_kline.Time == 0 {
		datastruct.InitKlineByTrade(cur_kline, trade)
		// fmt.Printf("\nUpdate datastruct.Kline: %s\n", cur_kline.String())
		return
	}

	if cur_kline.High < trade.Price {
		cur_kline.High = trade.Price
	}
	if cur_kline.Low > trade.Price {
		cur_kline.Low = trade.Price
	}

	cur_kline.Close = trade.Price
	cur_kline.Volume += trade.Volume

	// fmt.Printf("\nUpdate datastruct.Kline: %s\n", cur_kline.String())
}

func (a *Aggregator) cache_depth(depth *datastruct.DepthQuote) {

	a.depth_mutex.Lock()

	new_depth := datastruct.NewDepth(depth)

	LOG_INFO("\n******* <<Cache Depth>>: " + depth.String(5))

	if _, ok := a.depth_cache[new_depth.Symbol]; ok == false {
		a.depth_cache[new_depth.Symbol] = make(map[string]*datastruct.DepthQuote)
	}

	a.depth_cache[new_depth.Symbol][new_depth.Exchange] = new_depth

	defer a.depth_mutex.Unlock()
}

func (a *Aggregator) cache_kline(kline *datastruct.Kline) {

}

func (a *Aggregator) cache_trade(trade *datastruct.Trade) {
	new_trade := datastruct.NewTrade(trade)
	new_trade.Exchange = datastruct.BCTS_EXCHANGE

	fmt.Printf(" Recv datastruct.Trade: %s\n", trade.String())

	a.update_kline(trade)
	a.publish_trade(new_trade)
}

func (a *Aggregator) publish_depth(depth *datastruct.DepthQuote) {
	// LOG_INFO("publish_depth: " + depth.String(5))
}

func (a *Aggregator) publish_kline(kline *datastruct.Kline) {
	fmt.Printf("\nPub datastruct.Kline: \n%s\n", kline.String())
}

func (a *Aggregator) publish_trade(trade *datastruct.Trade) {
	fmt.Printf("Pub datastruct.Trade: %s\n", trade.String())
}

func GetTestDepthByType(index int) *datastruct.DepthQuote {
	var rst datastruct.DepthQuote

	exchange_array := []string{"FTX", "HUOBI", "OKEX"}
	exchange_type := index % 3

	rst.Exchange = exchange_array[exchange_type]
	rst.Symbol = "BTC_USDT"
	rst.Time = time.Now().Unix()
	rst.Asks = treemap.NewWith(utils.Float64Comparator)
	rst.Bids = treemap.NewWith(utils.Float64Comparator)

	rst.Asks.Put(55000.0, &datastruct.InnerDepth{5.5, map[string]float64{rst.Exchange: 5.5}})
	rst.Asks.Put(50000.0, &datastruct.InnerDepth{5.0, map[string]float64{rst.Exchange: 5.0}})

	rst.Bids.Put(45000.0, &datastruct.InnerDepth{4.5, map[string]float64{rst.Exchange: 4.5}})
	rst.Bids.Put(40000.0, &datastruct.InnerDepth{4.0, map[string]float64{rst.Exchange: 4.0}})

	switch exchange_type {
	case 0:
		rst.Asks.Put(60000.0, &datastruct.InnerDepth{6.0, map[string]float64{rst.Exchange: 6.0}})
		rst.Bids.Put(35000.0, &datastruct.InnerDepth{3.5, map[string]float64{rst.Exchange: 3.5}})

	case 1:
		rst.Asks.Put(70000.0, &datastruct.InnerDepth{7.0, map[string]float64{rst.Exchange: 7.0}})
		rst.Bids.Put(30000.0, &datastruct.InnerDepth{3.0, map[string]float64{rst.Exchange: 3.0}})

	case 2:
		rst.Asks.Put(75000.0, &datastruct.InnerDepth{7.5, map[string]float64{rst.Exchange: 7.5}})
		rst.Bids.Put(25000.0, &datastruct.InnerDepth{2.5, map[string]float64{rst.Exchange: 2.5}})
	}

	return &rst
}

func GetTestTrade() *datastruct.Trade {
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(3)

	exchange_array := []string{"FTX", "HUOBI", "OKEX"}
	cur_exchange := exchange_array[randomNum%3]
	symbol := "ETH_USDT"
	trade_price := float64(rand.Intn(1000))
	trade_volume := float64(rand.Intn(100))

	new_trade := datastruct.NewTrade(nil)
	new_trade.Exchange = cur_exchange
	new_trade.Symbol = symbol
	new_trade.Price = trade_price
	new_trade.Volume = trade_volume
	new_trade.Time = time.Now().Unix()

	fmt.Printf("Send datastruct.Trade: %s\n", new_trade.String())

	return new_trade
}

func PublishTest(data *datastruct.DataChannel) {
	timer := time.Tick(3 * time.Second)

	// index := 0
	for {
		select {
		case <-timer:
			// depth_quote := GetTestDepthByType(index)
			// index++
			// data.DepthChannel <- depth_quote

			data.TradeChannel <- GetTestTrade()
		}
	}
}

func TestAggregator() {
	aggregator := Aggregator{
		depth_cache:               make(map[string]map[string]*datastruct.DepthQuote),
		kline_cache:               make(map[string]map[string]*datastruct.Kline),
		trade_cache:               make(map[string]map[string]*datastruct.Trade),
		kline_aggregated:          make(map[string]*datastruct.Kline),
		depth_aggregator_millsecs: 5000,
	}

	data_chan := datastruct.DataChannel{
		DepthChannel: make(chan *datastruct.DepthQuote),
		KlineChannel: make(chan *datastruct.Kline),
		TradeChannel: make(chan *datastruct.Trade),
	}
	aggregator.Start(&data_chan)

	go PublishTest(&data_chan)

	time.Sleep(time.Hour)
}
