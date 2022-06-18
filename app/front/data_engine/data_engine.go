package data_engine

import (
	"context"
	"fmt"
	"market_server/app/dataManager/rpc/marketservice"
	"market_server/app/front/config"
	"market_server/app/front/net"
	"market_server/app/front/worker"
	"market_server/common/datastruct"
	"market_server/common/util"
	"sync"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

type DataEngine struct {
	RecvDataChan *datastruct.DataChannel
	config       *config.Config

	msclient marketservice.MarketService

	next_worker worker.WorkerI

	depth_cache_map sync.Map
	trade_cache_map sync.Map

	cache_period_data       map[string]*PeriodData // 缓存24小时1分频率的 k 线数据，用来计算24小时的涨跌幅;
	cache_period_data_mutex sync.Mutex

	cache_kline_data map[string](map[int]*treemap.Map)

	IsTest bool
}

func NewDataEngine(recvDataChan *datastruct.DataChannel, config *config.Config) *DataEngine {

	rst := &DataEngine{
		RecvDataChan:      recvDataChan,
		config:            config,
		cache_period_data: make(map[string]*PeriodData),
		cache_kline_data:  make(map[string]map[int]*treemap.Map),
		msclient:          marketservice.NewMarketService(zrpc.MustNewClient(config.RpcConfig)),
	}

	return rst
}

func (d *DataEngine) SetNextWorker(next_worker worker.WorkerI) {
	d.next_worker = next_worker
}

func (d *DataEngine) UpdateMeta(symbols []string) {
	for _, symbol := range symbols {
		if _, ok := d.cache_period_data[symbol]; !ok {
			d.InitPeriodDara(symbol)
		}
	}
}

func (a *DataEngine) InitPeriodDara(symbol string) {
	a.cache_period_data[symbol] = &PeriodData{
		TimeNanos:       24 * 60 * 60 * 1000000000,
		Count:           0,
		MaxTime:         0,
		MinTime:         0,
		time_cache_data: treemap.NewWith(utils.Int64Comparator),
	}

	end_time_nanos := uint64(util.TimeMinuteNanos())
	start_time_nanos := uint64(end_time_nanos - 24*60*1000000000)

	req_hist_info := &marketservice.ReqHishKlineInfo{
		Symbol:    symbol,
		Exchange:  datastruct.BCTS_EXCHANGE,
		StartTime: start_time_nanos,
		EndTime:   end_time_nanos,
		Count:     24 * 60,
		Frequency: 60,
	}

	hist_klines, err := a.msclient.RequestHistKlineData(context.Background(), req_hist_info)

	if err != nil {
		fmt.Printf("err %+v \n", err)
		logx.Errorf("ReqHistKline Err: %+v\n", err)
	}

	fmt.Printf("Rst: %+v \n", hist_klines)

	a.cache_period_data[symbol].UpdateWithPbKlines(hist_klines)
}

func (a *DataEngine) StartListenRecvdata() {
	logx.Info("[S] DBServer start_listen_recvdata")
	go func() {
		for {
			select {
			case new_depth := <-a.RecvDataChan.DepthChannel:
				go a.process_depth(new_depth)
			case new_kline := <-a.RecvDataChan.KlineChannel:
				go a.process_kline(new_kline)
			case new_trade := <-a.RecvDataChan.TradeChannel:
				go a.process_trade(new_trade)
			}
		}
	}()
	logx.Info("[S] DBServer start_receiver Over!")
}

func (d *DataEngine) process_depth(depth *datastruct.DepthQuote) error {
	depth.Time = depth.Time / datastruct.NANO_PER_SECS

	d.depth_cache_map.Store(depth.Symbol, depth)

	d.PublishDepth(depth, nil)
	return nil
}

func (d *DataEngine) process_kline(kline *datastruct.Kline) error {
	kline.Time = kline.Time / datastruct.NANO_PER_SECS

	if _, ok := d.cache_period_data[kline.Symbol]; !ok {
		d.InitPeriodDara(kline.Symbol)

		symbol_list := d.get_symbol_list()
		d.PublishSymbol(symbol_list, nil)
	}

	d.cache_period_data[kline.Symbol].UpdateWithKline(kline)

	d.PublishKline(kline, nil)

	d.PublishChangeinfo(d.cache_period_data[kline.Symbol].GetChangeInfo(), nil)

	return nil
}

func (d *DataEngine) process_trade(trade *datastruct.Trade) error {
	trade.Time = trade.Time / datastruct.NANO_PER_SECS

	d.trade_cache_map.Store(trade.Symbol, trade)

	d.cache_period_data[trade.Symbol].UpdateWithTrade(trade)

	d.PublishChangeinfo(d.cache_period_data[trade.Symbol].GetChangeInfo(), nil)

	d.PublishTrade(trade, nil)
	return nil
}

func (d *DataEngine) get_symbol_list() []string {
	d.cache_period_data_mutex.Lock()
	defer d.cache_period_data_mutex.Unlock()

	var rst []string
	for key := range d.cache_period_data {
		rst = append(rst, key)
	}
	return rst
}

func (d *DataEngine) SubSymbol(ws *net.WSInfo) {
	symbol_list := d.get_symbol_list()

	d.next_worker.PublishSymbol(symbol_list, ws)
}

func (d *DataEngine) PublishSymbol(symbol_list []string, ws *net.WSInfo) {
	d.next_worker.PublishSymbol(symbol_list, ws)
}

func (d *DataEngine) PublishDepth(depth *datastruct.DepthQuote, ws *net.WSInfo) {
	d.next_worker.PublishDepth(depth, ws)
}

func (d *DataEngine) PublishTrade(trade *datastruct.Trade, ws *net.WSInfo) {
	d.next_worker.PublishTrade(trade, ws)
}

func (d *DataEngine) PublishKline(kline *datastruct.Kline, ws *net.WSInfo) {
	d.next_worker.PublishKline(kline, ws)
}

func (d *DataEngine) PublishHistKline(kline *datastruct.RspHistKline, ws *net.WSInfo) {
	d.next_worker.PublishHistKline(kline, ws)
}

func (d *DataEngine) PublishChangeinfo(change_info *datastruct.ChangeInfo, ws *net.WSInfo) {
	d.next_worker.PublishChangeinfo(change_info, ws)
}

func (d *DataEngine) SubTrade(symbol string, ws *net.WSInfo) {

	if trade, ok := d.trade_cache_map.Load(symbol); ok {
		d.PublishTrade(trade.(*datastruct.Trade), ws)
	}
}

func (d *DataEngine) SubDepth(symbol string, ws *net.WSInfo) {

	if depth, ok := d.depth_cache_map.Load(symbol); ok {
		d.PublishDepth(depth.(*datastruct.DepthQuote), ws)
	}
}

func (d *DataEngine) GetHistKlineData(req_kline_info *datastruct.ReqHistKline) *datastruct.RspHistKline {

	if d.IsTest {
		return datastruct.GetTestHistKline(req_kline_info)
	}

	req_hist_info := &marketservice.ReqHishKlineInfo{
		Symbol:    req_kline_info.Symbol,
		Exchange:  req_kline_info.Exchange,
		StartTime: req_kline_info.StartTime,
		EndTime:   req_kline_info.EndTime,
		Count:     req_kline_info.Count,
		Frequency: req_kline_info.Frequency,
	}

	hist_klines, err := d.msclient.RequestHistKlineData(context.Background(), req_hist_info)

	if err != nil {
		return nil
	}

	tmp := treemap.NewWith(utils.Int64Comparator)

	for _, pb_kline := range hist_klines.KlineData {
		kline := marketservice.NewKlineWithPbKline(pb_kline)
		if kline == nil {
			continue
		}

		tmp.Put(kline.Time, kline)
	}

	d.UpdateCacheKlinesWithHist(tmp)

	trans_kline := d.TrasOriKlineData(req_kline_info, tmp)

	return &datastruct.RspHistKline{
		ReqInfo: req_kline_info,
		Klines:  trans_kline,
	}
}

func (d *DataEngine) TrasOriKlineData(req_kline_info *datastruct.ReqHistKline, ori_klines *treemap.Map) *treemap.Map {
	rst := treemap.NewWith(utils.Int64Comparator)
	resolution := int(req_kline_info.Frequency)

	if ori_klines.Size() == 0 {
		return rst
	}

	iter := ori_klines.Iterator()
	iter.Begin()

	cache_kline := iter.Value().(*datastruct.Kline)

	if !datastruct.IsNewKlineStart(cache_kline, int64(resolution)) {
		cache_kline.Time = cache_kline.Time - cache_kline.Time%int64(resolution)
	}

	for iter.Next() {
		cur_kline := iter.Value().(*datastruct.Kline)

		if datastruct.IsOldKlineEnd(cur_kline, int64(resolution)) {
			var pub_kline *datastruct.Kline

			if cur_kline.Resolution != resolution {
				cache_kline.Close = cur_kline.Close
				cache_kline.Low = util.MinFloat64(cache_kline.Low, cur_kline.Low)
				cache_kline.High = util.MaxFloat64(cache_kline.High, cur_kline.High)
				cache_kline.Volume += cur_kline.Volume

				pub_kline = datastruct.NewKlineWithKline(cache_kline)
			} else {
				pub_kline = datastruct.NewKlineWithKline(cur_kline)
			}

			cache_kline = datastruct.NewKlineWithKline(pub_kline)
			rst.Put(pub_kline.Time, pub_kline)
		} else if datastruct.IsNewKlineStart(cur_kline, int64(resolution)) {
			cache_kline = cur_kline
			cache_kline.Resolution = resolution
		} else {
			cache_kline.Close = cur_kline.Close
			cache_kline.Low = util.MinFloat64(cache_kline.Low, cur_kline.Low)
			cache_kline.High = util.MaxFloat64(cache_kline.High, cur_kline.High)
		}
	}

	return rst
}

func (d *DataEngine) UpdateCacheKlinesWithHist(klines *treemap.Map) {

}

func (d *DataEngine) SubKline(req_kline_info *datastruct.ReqHistKline, ws *net.WSInfo) {

	rst := d.GetHistKlineData(req_kline_info)

	d.next_worker.PublishHistKline(rst, ws)
}
