package aggregate

import (
	"fmt"
	"market_aggregate/pkg/conf"
	"market_aggregate/pkg/datastruct"
	"market_aggregate/pkg/util"
	"sync"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
)

type Aggregator struct {
	depth_cache map[string]map[string]*datastruct.DepthQuote

	kline_cache map[string]map[string]*datastruct.Kline
	trade_cache map[string]map[string]*datastruct.Trade

	kline_aggregated map[string]*datastruct.Kline

	depth_mutex sync.Mutex
	kline_mutex sync.Mutex
	trade_mutex sync.Mutex

	AggConfig      conf.AggregateConfig
	AggConfigMutex sync.RWMutex

	RecvDataChan *datastruct.DataChannel
	SendDataChan *datastruct.DataChannel

	RiskWorker *RiskWorkerManager
}

func (a *Aggregator) Init(RecvDataChan *datastruct.DataChannel, SendDataChan *datastruct.DataChannel) {
	a.RecvDataChan = RecvDataChan
	a.SendDataChan = SendDataChan
	a.AggConfig = conf.AggregateConfig{
		DepthAggregatorMillsecs: 5,
	}
	a.RiskWorker = &RiskWorkerManager{}

	a.RiskWorker.Init()

	if a.depth_cache == nil {
		a.depth_cache = make(map[string]map[string]*datastruct.DepthQuote)
	}

	if a.kline_cache == nil {
		a.kline_cache = make(map[string]map[string]*datastruct.Kline)
	}

	if a.trade_cache == nil {
		a.trade_cache = make(map[string]map[string]*datastruct.Trade)
	}

	if a.kline_aggregated == nil {
		a.kline_aggregated = make(map[string]*datastruct.Kline)
	}
}

func (a *Aggregator) UpdateConfig(config conf.AggregateConfig) {
	defer a.AggConfigMutex.Unlock()

	a.AggConfigMutex.Lock()

	a.AggConfig = config
}

func (a *Aggregator) Start() {
	a.start_listen_recvdata()
	a.start_aggregate_depth()
	a.start_aggregate_kline()
}

func (a *Aggregator) get_wait_millsecs(curr_time int, fre int) int {
	if curr_time == 0 {
		return fre
	} else if curr_time%fre == 0 {
		return 0
	} else {
		return fre - curr_time%fre
	}
}

func (a *Aggregator) get_sleep_millsecs(curr_time int) (int, []string) {
	publish_list := make([]string, 10)
	sleep_millsecs := 0

	if a.AggConfig.DepthAggregatorMillsecsMap == nil || len(a.AggConfig.DepthAggregatorMillsecsMap) == 0 {
		sleep_millsecs = 3 * 1000
	} else {
		min_sleep_secs := time.Hour

		for symbol, cur_publish_millsecs := range a.AggConfig.DepthAggregatorMillsecsMap {
			cur_sleep_millsecs := a.get_wait_millsecs(curr_time, int(cur_publish_millsecs))

			if cur_sleep_millsecs == 0 {
				publish_list = append(publish_list, symbol)
			} else if sleep_millsecs == 0 || sleep_millsecs > cur_sleep_millsecs {
				sleep_millsecs = cur_sleep_millsecs
			}

			if min_sleep_secs > cur_publish_millsecs {
				min_sleep_secs = cur_publish_millsecs
			}
		}

		// 所有币对当前都需要处理，下次处理的时间就是发布等待时间最短的那个币对
		if sleep_millsecs == 0 {
			sleep_millsecs = int(min_sleep_secs)
		}
	}

	return sleep_millsecs, publish_list
}

func (a *Aggregator) start_aggregate_depth() {
	util.LOG_INFO("Aggregator start_aggregate_depth!")
	go func() {
		util.LOG_INFO("start_aggregate_depth")
		// timer := time.Tick(time.Duration(a.depth_aggregator_millsecs * time.Millisecond))
		// timeout := time.After(time.Millisecond * a.depth_aggregator_millsecs)

		curr_time := 0
		for {
			sleep_millsecs, aggregate_symbol_list := a.get_sleep_millsecs(curr_time)

			a.aggregate_depth(aggregate_symbol_list)

			curr_time += sleep_millsecs

			if curr_time > 24*60*60*1000 {
				curr_time = 0
			}

			time.Sleep(time.Millisecond * time.Duration(sleep_millsecs))

		}
	}()
	util.LOG_INFO("Aggregator start_aggregate_depth Over!")
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
	util.LOG_INFO("Aggregate datastruct.Kline Start!")
	go func() {
		for {
			util.WaitForNextMinute()

			a.aggregate_kline()
		}
	}()
	util.LOG_INFO("Aggregate datastruct.Kline Over!")
}

func (a *Aggregator) aggregate_depth(symbol_list []string) {

	// util.LOG_INFO("----- Aggregate Depth Start ------ ")

	for _, symbol := range symbol_list {

		a.depth_mutex.Lock()

		if exchange_depth_map, ok := a.depth_cache[symbol]; ok == true {
			new_depth := datastruct.NewDepth(nil)
			new_depth.Symbol = symbol
			new_depth.Exchange = datastruct.BCTS_EXCHANGE
			new_depth.Time = int64(util.UTCNanoTime())

			for exchange, cur_depth := range exchange_depth_map {
				util.LOG_INFO("\n===== <<CurDepth>>: " + cur_depth.String(5))
				mix_depth(new_depth.Asks, cur_depth.Asks, exchange)
				mix_depth(new_depth.Bids, cur_depth.Bids, exchange)
			}

			util.LOG_INFO("\n^^^^^^^ <<aagregated_depth>>: " + new_depth.String(5))

			a.publish_depth(new_depth)
		}

		a.depth_mutex.Unlock()
	}

	// for symbol, exchange_depth_map := range a.depth_cache {

	// 	// for _, cur_depth := range exchange_depth_map {
	// 	// 	util.LOG_INFO("\n===== After <<CurDepth>>: " + cur_depth.String(5))
	// 	// }
	// }

	// util.LOG_INFO("----- Aggregate Depth Over!------ \n")

	defer a.depth_mutex.Unlock()
}

func (a *Aggregator) start_listen_recvdata() {
	util.LOG_INFO("Aggregator start_listen_recvdata")
	go func() {
		for {
			select {
			case new_depth := <-a.RecvDataChan.DepthChannel:
				a.cache_depth(new_depth)
			case new_kline := <-a.RecvDataChan.KlineChannel:
				a.cache_kline(new_kline)
			case new_trade := <-a.RecvDataChan.TradeChannel:
				a.cache_trade(new_trade)
			}
		}
	}()
	util.LOG_INFO("Aggregator start_receiver Over!")
}

func (a *Aggregator) aggregate_kline() {
	defer a.kline_mutex.Unlock()
	a.kline_mutex.Lock()

	for _, kline := range a.kline_aggregated {
		new_kline := datastruct.NewKline(kline)
		kline.Time = 0
		new_kline.Time = int64(util.UTCNanoTime())
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

	util.LOG_INFO("\n******* <<Cache Depth>>: " + depth.String(5))

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
	// util.LOG_INFO("publish_depth: " + depth.String(5))
}

func (a *Aggregator) publish_kline(kline *datastruct.Kline) {
	fmt.Printf("\nPub datastruct.Kline: \n%s\n", kline.String())

	a.SendDataChan.KlineChannel <- kline
}

func (a *Aggregator) publish_trade(trade *datastruct.Trade) {
	fmt.Printf("Pub datastruct.Trade: %s\n", trade.String())

	a.SendDataChan.TradeChannel <- trade
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

			data.TradeChannel <- datastruct.GetTestTrade()
		}
	}
}

func TestAggregator() {
	aggregator := Aggregator{}

	RecvDataChan := &datastruct.DataChannel{
		DepthChannel: make(chan *datastruct.DepthQuote),
		KlineChannel: make(chan *datastruct.Kline),
		TradeChannel: make(chan *datastruct.Trade),
	}

	PubDataChan := &datastruct.DataChannel{
		DepthChannel: make(chan *datastruct.DepthQuote),
		KlineChannel: make(chan *datastruct.Kline),
		TradeChannel: make(chan *datastruct.Trade),
	}

	AggConfig := conf.AggregateConfig{
		DepthAggregatorMillsecs: 5,
	}

	aggregator.Init(RecvDataChan, PubDataChan)

	aggregator.UpdateConfig(AggConfig)

	aggregator.Start()

	go PublishTest(RecvDataChan)

	time.Sleep(time.Hour)
}
