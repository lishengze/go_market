package aggregate

import (

	// config "market_aggregate/app/conf"
	"fmt"
	"market_server/app/market_aggregate/config"
	mkconfig "market_server/app/market_aggregate/config"
	"market_server/common/datastruct"
	"market_server/common/dingtalk"
	"market_server/common/monitorStruct"
	"market_server/common/util"
	"strings"
	"sync"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/zeromicro/go-zero/core/logx"
)

type Aggregator struct {
	depth_cache map[string]map[string]*datastruct.DepthQuote

	kline_cache map[string]map[string]*datastruct.Kline
	trade_cache map[string]map[string]*datastruct.Trade

	depth_seq_map       map[string]uint64
	depth_seq_map_mutex sync.Mutex

	tradekline_seq_map       map[string]uint64
	tradekline_seq_map_mutex sync.Mutex

	kline_aggregated map[string]*datastruct.Kline

	depth_mutex      sync.Mutex
	tradekline_mutex sync.Mutex

	AggConfig      mkconfig.AggregateConfig
	AggConfigMutex sync.RWMutex

	RecvDataChan *datastruct.DataChannel
	PubDataChan  *datastruct.DataChannel

	RiskWorker *RiskWorkerManager

	MonitorMarketDataWorker *monitorStruct.MonitorMarketData
	MonitorChan             *monitorStruct.MonitorChannel
	DingClient              *dingtalk.Client
	MetaInfo                string

	cfg *config.Config
}

func NewAggregator(RecvDataChan *datastruct.DataChannel, PubDataChan *datastruct.DataChannel,
	RiskWorker *RiskWorkerManager, cfg *config.Config) (a *Aggregator) {

	monitor_chan := monitorStruct.NewMonitorChannel()
	dingtalk := dingtalk.NewClient(cfg.DingConfigInfo.Token, cfg.DingConfigInfo.Secret)

	return &Aggregator{
		RiskWorker:       RiskWorker,
		RecvDataChan:     RecvDataChan,
		PubDataChan:      PubDataChan,
		cfg:              cfg,
		depth_cache:      make(map[string]map[string]*datastruct.DepthQuote),
		kline_cache:      make(map[string]map[string]*datastruct.Kline),
		trade_cache:      make(map[string]map[string]*datastruct.Trade),
		kline_aggregated: make(map[string]*datastruct.Kline),

		tradekline_seq_map: make(map[string]uint64),
		depth_seq_map:      make(map[string]uint64),

		MetaInfo:                "Aggragte",
		MonitorChan:             monitor_chan,
		MonitorMarketDataWorker: monitorStruct.NewMonitorMarketData("Aggragte", &cfg.MonitorConfigInfo, monitor_chan),
		DingClient:              dingtalk,

		AggConfig: mkconfig.AggregateConfig{
			DepthAggregatorConfigMap: make(map[string]mkconfig.AggregateConfigAtom),
		},
	}
}

// func (a *Aggregator) Init(RecvDataChan *datastruct.DataChannel, PubDataChan *datastruct.DataChannel, RiskWorker *RiskWorkerManager) {
// 	a.RecvDataChan = RecvDataChan
// 	a.PubDataChan = PubDataChan
// 	a.AggConfig = mkconfig.AggregateConfig{
// 		DepthAggregatorConfigMap: make(map[string]mkconfig.AggregateConfigAtom),
// 	}
// 	a.RiskWorker = RiskWorker

// 	// a.RiskWorker.Init()

// 	if a.depth_cache == nil {
// 		a.depth_cache = make(map[string]map[string]*datastruct.DepthQuote)
// 	}

// 	if a.kline_cache == nil {
// 		a.kline_cache = make(map[string]map[string]*datastruct.Kline)
// 	}

// 	if a.trade_cache == nil {
// 		a.trade_cache = make(map[string]map[string]*datastruct.Trade)
// 	}

// 	if a.kline_aggregated == nil {
// 		a.kline_aggregated = make(map[string]*datastruct.Kline)
// 	}
// }

func (a *Aggregator) UpdateConfig(config mkconfig.AggregateConfig) {
	defer a.AggConfigMutex.Unlock()

	a.AggConfigMutex.Lock()

	a.AggConfig = config
}

// UnTest
func (a *Aggregator) UpdateDepthSeq(symbol string) {
	defer util.CatchExp(fmt.Sprintf("UpdateDepthSeq %s", symbol))

	a.depth_seq_map_mutex.Lock()
	defer a.depth_seq_map_mutex.Unlock()

	if _, ok := a.depth_seq_map[symbol]; !ok {
		a.depth_seq_map[symbol] = 0
	}

	a.depth_seq_map[symbol] = a.depth_seq_map[symbol] + 1
}

// UnTest
func (a *Aggregator) UpdateTradeKlineSeq(symbol string) {
	defer util.CatchExp(fmt.Sprintf("UpdateTradeKlineSeq %s", symbol))

	a.tradekline_seq_map_mutex.Lock()
	defer a.tradekline_seq_map_mutex.Unlock()

	if _, ok := a.tradekline_seq_map[symbol]; !ok {
		a.tradekline_seq_map[symbol] = 0
	}

	a.tradekline_seq_map[symbol] = a.tradekline_seq_map[symbol] + 1

}

// UnTest
func (a *Aggregator) GetDepthSeq(symbol string) uint64 {
	defer util.CatchExp(fmt.Sprintf("GetDepthSeq %s", symbol))

	a.depth_seq_map_mutex.Lock()
	defer a.depth_seq_map_mutex.Unlock()

	if _, ok := a.depth_seq_map[symbol]; !ok {
		a.depth_seq_map[symbol] = 0
	}

	return a.depth_seq_map[symbol]
}

// UnTest
func (a *Aggregator) GetTradeKlineSeq(symbol string) uint64 {
	defer util.CatchExp(fmt.Sprintf("GetTradeKlineSeq %s", symbol))

	a.tradekline_seq_map_mutex.Lock()
	defer a.tradekline_seq_map_mutex.Unlock()

	if _, ok := a.tradekline_seq_map[symbol]; !ok {
		a.tradekline_seq_map[symbol] = 0
	}

	return a.tradekline_seq_map[symbol]
}

func (a *Aggregator) Start() {
	a.start_listen_recvdata()
	a.start_aggregate_depth()
	a.start_aggregate_kline()

	go a.MonitorMarketDataWorker.StartCheck()
	go a.StartListenInvalidData()
}

func (k *Aggregator) StartListenInvalidData() {
	logx.Info("[S] Aggregator StartListenInvalidData Start!")
	go func() {
		for {
			select {
			case invalid_depth := <-k.MonitorChan.DepthChan:
				go k.process_invalid_depth(invalid_depth)
			case invalid_trade := <-k.MonitorChan.TradeChan:
				go k.process_invalid_trade(invalid_trade)
			case invalid_kline := <-k.MonitorChan.KlineChan:
				go k.process_invalid_kline(invalid_kline)
			}
		}
	}()
	logx.Info("[S] Aggregator StartListenInvalidData Over!")
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

	a.AggConfigMutex.RLock()

	default_sleep_millsecs := 10 * 1000

	if a.AggConfig.DepthAggregatorConfigMap == nil || len(a.AggConfig.DepthAggregatorConfigMap) == 0 {
		sleep_millsecs = default_sleep_millsecs // 还未获取配置信息;
	} else {
		// logx.Info(a.AggConfig.String())
		min_sleep_secs := time.Duration(default_sleep_millsecs)

		for symbol, publish_config := range a.AggConfig.DepthAggregatorConfigMap {

			if publish_config.IsPublish == false {
				continue
			}

			cur_sleep_millsecs := a.get_wait_millsecs(curr_time, int(publish_config.AggregateFreq))

			if cur_sleep_millsecs == 0 {
				publish_list = append(publish_list, symbol)
			} else if sleep_millsecs == 0 || sleep_millsecs > cur_sleep_millsecs {
				sleep_millsecs = cur_sleep_millsecs
			}

			if min_sleep_secs > publish_config.AggregateFreq {
				min_sleep_secs = publish_config.AggregateFreq
			}
		}

		// 所有币对当前都需要处理，下次处理的时间就是发布等待时间最短的那个币对
		if sleep_millsecs == 0 {
			sleep_millsecs = int(min_sleep_secs)
		}
	}

	defer a.AggConfigMutex.RUnlock()

	return sleep_millsecs, publish_list
}

func (a *Aggregator) start_aggregate_depth() {
	logx.Info("Aggregator start_aggregate_depth!")
	go func() {
		logx.Info("start_aggregate_depth")
		// timer := time.Tick(time.Duration(a.depth_aggregator_millsecs * time.Millisecond))
		// timeout := time.After(time.Millisecond * a.depth_aggregator_millsecs)

		curr_time := 0
		for {
			sleep_millsecs, aggregate_symbol_list := a.get_sleep_millsecs(curr_time)

			// logx.Info(fmt.Sprintf("curr_time:%+v, sleep_millsecs: %+v, aggregate_symbol_list%+v",
			// 	curr_time, sleep_millsecs, aggregate_symbol_list))

			a.aggregate_depth(aggregate_symbol_list)

			curr_time += sleep_millsecs

			if curr_time > 24*60*60*1000 {
				curr_time = 0
			}

			time.Sleep(time.Millisecond * time.Duration(sleep_millsecs))

		}
	}()
	logx.Info("Aggregator start_aggregate_depth Over!")
}

func mix_depth(src *treemap.Map, other *treemap.Map, exchange string) {
	other_iter := other.Iterator()

	for other_iter.Begin(); other_iter.Next(); {
		if cur_iter, ok := src.Get(other_iter.Key()); ok {
			cur_innerdepth := cur_iter.(*datastruct.InnerDepth)
			cur_innerdepth.Volume += other_iter.Value().(*datastruct.InnerDepth).Volume
			cur_innerdepth.ExchangeVolume[exchange] = other_iter.Value().(*datastruct.InnerDepth).Volume
		} else {
			src_inner_depth := datastruct.InnerDepth{Volume: 0, ExchangeVolume: make(map[string]float64)}

			other_inner_depth := other_iter.Value().(*datastruct.InnerDepth)
			src_inner_depth.Volume = other_inner_depth.Volume
			src_inner_depth.ExchangeVolume[exchange] = other_inner_depth.Volume
			price := other_iter.Key().(float64)

			src.Put(price, &src_inner_depth)
		}
	}
}

func (a *Aggregator) start_aggregate_kline() {
	logx.Info("Aggregate datastruct.Kline Start!")
	go func() {
		for {
			util.WaitForNextMinute()

			a.aggregate_kline()
		}
	}()
	logx.Info("Aggregate datastruct.Kline Over!")
}

func (a *Aggregator) aggregate_depth(symbol_list []string) {

	for _, symbol := range symbol_list {

		a.depth_mutex.Lock()

		if exchange_depth_map, ok := a.depth_cache[symbol]; ok == true {
			new_depth := datastruct.NewDepth(nil)
			new_depth.Symbol = symbol
			new_depth.Exchange = datastruct.BCTS_EXCHANGE
			new_depth.Time = int64(util.UTCNanoTime())

			for exchange, cur_depth := range exchange_depth_map {
				// logx.Info("\n===== <<CurDepth>>: " + cur_depth.String(3))
				mix_depth(new_depth.Asks, cur_depth.Asks, exchange)
				mix_depth(new_depth.Bids, cur_depth.Bids, exchange)
			}

			// logx.Statf("\n[AD]: " + new_depth.String(3))

			a.publish_depth(new_depth)
		}

		a.depth_mutex.Unlock()
	}
}

func (a *Aggregator) start_listen_recvdata() {
	logx.Info("Aggregator start_listen_recvdata")
	go func() {
		for {
			select {
			case new_depth := <-a.RecvDataChan.DepthChannel:
				// if new_depth.Symbol == a.cfg.RiskTestConfig.TestSymbol {
				// 	logx.Info("\n[Rcv] Depth " + new_depth.String(3))
				// }
				a.cache_depth(new_depth)
			case new_kline := <-a.RecvDataChan.KlineChannel:
				// if new_kline.Symbol == a.cfg.RiskTestConfig.TestSymbol {
				// 	logx.Info("\n[Rcv] Kline " + new_kline.String())
				// }
				a.cache_kline(new_kline)
			case new_trade := <-a.RecvDataChan.TradeChannel:
				// if new_trade.Symbol == a.cfg.RiskTestConfig.TestSymbol {
				// 	logx.Info("\n[Rcv] Trade " + new_trade.String())
				// }
				a.cache_trade(new_trade)
			}
		}
	}()
	logx.Info("Aggregator start_receiver Over!")
}

func (a *Aggregator) aggregate_kline() {
	util.CatchExp(fmt.Sprintf("aggregate_kline"))

	a.tradekline_mutex.Lock()
	defer a.tradekline_mutex.Unlock()

	// logx.Info(fmt.Sprintf("------ Start aggregate_kline time: %v", time.Now().UTC()))

	for _, kline := range a.kline_aggregated {

		kline.SetHistoryFlag()

		if !kline.HasTrade() {
			sequence := a.GetTradeKlineSeq(kline.Symbol)

			kline.Sequence = sequence

			a.UpdateTradeKlineSeq(kline.Symbol)
		}

		kline.Time = util.LastUTCMinuteNano() // 当前时间超过了上一分钟;
		new_kline := datastruct.NewKline(kline)

		logx.Slowf("[PHK]: %s", new_kline.FullString())
		a.publish_kline(new_kline)

		kline.RestWithLastPrice()
	}

	// for k, _ := range a.kline_aggregated {
	// 	delete(a.kline_aggregated, k)
	// }
}

func (a *Aggregator) update_kline(trade *datastruct.Trade, sequence uint64) {
	defer util.CatchExp(fmt.Sprintf("update_kline %s", trade.String()))

	if _, ok := a.kline_aggregated[trade.Symbol]; !ok {
		a.kline_aggregated[trade.Symbol] = datastruct.NewKline(nil)
	}

	cur_kline := a.kline_aggregated[trade.Symbol]
	cur_kline.Sequence = sequence

	if !cur_kline.IsInited() {
		logx.Slowf("Kline Init By Trade: %s", trade.String())
		datastruct.InitKlineByTrade(cur_kline, trade)
	} else {
		if cur_kline.High < trade.Price {
			cur_kline.High = trade.Price
		}
		if cur_kline.Low > trade.Price {
			cur_kline.Low = trade.Price
		}

		cur_kline.Close = trade.Price
		cur_kline.Volume += trade.Volume
		cur_kline.LastVolume = trade.Volume
		cur_kline.Time = trade.Time
	}

	// logx.Slowf("[PK]: %s", cur_kline.FullString())

	a.publish_kline(cur_kline)
}

func (a *Aggregator) update_trade(trade *datastruct.Trade, sequence uint64) {
	defer util.CatchExp(fmt.Sprintf("update_trade %s", trade.String()))

	new_trade := datastruct.NewTrade(trade)
	new_trade.Exchange = datastruct.BCTS_EXCHANGE
	new_trade.Sequence = sequence

	// logx.Slowf("[PT]: %s", new_trade.String())

	a.publish_trade(new_trade)
}

func (a *Aggregator) cache_depth(depth *datastruct.DepthQuote) {

	a.MonitorMarketDataWorker.UpdateDepth(depth.Exchange + "_" + depth.Symbol)

	a.depth_mutex.Lock()

	new_depth := datastruct.NewDepth(depth)

	if _, ok := a.depth_cache[new_depth.Symbol]; !ok {
		a.depth_cache[new_depth.Symbol] = make(map[string]*datastruct.DepthQuote)
	}

	a.depth_cache[new_depth.Symbol][new_depth.Exchange] = new_depth

	defer a.depth_mutex.Unlock()
}

func (a *Aggregator) cache_trade(trade *datastruct.Trade) {
	defer util.CatchExp(fmt.Sprintf("cache_trade %s", trade.String()))

	a.tradekline_mutex.Lock()
	defer a.tradekline_mutex.Unlock()

	a.MonitorMarketDataWorker.UpdateTrade(trade.Exchange + "_" + trade.Symbol)

	sequence := a.GetTradeKlineSeq(trade.Symbol)
	a.UpdateTradeKlineSeq(trade.Symbol)

	a.update_kline(trade, sequence)
	a.update_trade(trade, sequence)

	logx.Info("")
}

func (a *Aggregator) cache_kline(kline *datastruct.Kline) {

}

func (a *Aggregator) publish_depth(depth *datastruct.DepthQuote) {

	aggregate_depth := datastruct.NewDepth(depth)
	aggregate_depth.Exchange = datastruct.BCTS_EXCHANGE_AGGREGATE
	a.PubDataChan.DepthChannel <- aggregate_depth

	if a.RiskWorker != nil {
		a.RiskWorker.Execute(depth)
		new_depth := datastruct.NewDepth(depth)

		a.PubDataChan.DepthChannel <- new_depth
	}

}

func (a *Aggregator) publish_kline(kline *datastruct.Kline) {
	// fmt.Printf("\nPub datastruct.Kline: \n%s\n", kline.String())
	a.PubDataChan.KlineChannel <- kline
}

func (a *Aggregator) publish_trade(trade *datastruct.Trade) {
	// fmt.Printf("Pub datastruct.Trade: %s\n", trade.String())

	a.PubDataChan.TradeChannel <- trade
}

func (a *Aggregator) process_invalid_depth(montior_atom *monitorStruct.MonitorAtom) error {
	logx.Info(montior_atom.InvalidInfo)
	a.DingClient.SendMessage("Depth: \n" + montior_atom.InvalidInfo)
	info_list := strings.Split(montior_atom.Symbol, "_")
	exchange := info_list[0]
	symbol := info_list[1]

	a.depth_mutex.Lock()

	delete(a.depth_cache[symbol], exchange)

	if len(a.depth_cache[symbol]) == 0 {
		delete(a.depth_cache, symbol)
	}

	defer a.depth_mutex.Unlock()

	return nil
}

func (k *Aggregator) process_invalid_trade(montior_atom *monitorStruct.MonitorAtom) error {
	logx.Info(montior_atom.InvalidInfo)
	k.DingClient.SendMessage("Trade: \n" + montior_atom.InvalidInfo)

	// info_list := strings.Split(montior_atom.Symbol, "_")
	// exchange := info_list[0]
	// symbol := info_list[1]

	return nil
}

func (k *Aggregator) process_invalid_kline(montior_atom *monitorStruct.MonitorAtom) error {
	logx.Info(montior_atom.InvalidInfo)
	k.DingClient.SendMessage("kline: \n" + montior_atom.InvalidInfo)

	// info_list := strings.Split(montior_atom.Symbol, "_")
	// exchange := info_list[0]
	// symbol := info_list[1]

	return nil
}
