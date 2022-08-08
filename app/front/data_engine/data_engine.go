package data_engine

import (
	"context"
	"fmt"
	"market_server/app/data_manager/rpc/marketservice"
	"market_server/app/front/net"
	"market_server/app/front/svc"
	"market_server/app/front/worker"
	"market_server/common/datastruct"
	"market_server/common/util"
	"strings"
	"sync"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

type DataEngine struct {
	RecvDataChan *datastruct.DataChannel
	// config       *config.Config
	ctx *svc.ServiceContext

	msclient marketservice.MarketService

	next_worker worker.WorkerI

	depth_cache_map sync.Map
	trade_cache_map sync.Map

	cache_period_data       map[string]*PeriodData // 缓存24小时1分频率的 k 线数据，用来计算24小时的涨跌幅;
	cache_period_data_mutex sync.Mutex

	cache_kline_data map[string](map[int]*treemap.Map)

	kline_cache *datastruct.KlineCache

	IsTest bool
}

func NewDataEngine(recvDataChan *datastruct.DataChannel, ctx *svc.ServiceContext) *DataEngine {

	if ctx != nil {
		return &DataEngine{
			RecvDataChan:      recvDataChan,
			ctx:               ctx,
			cache_period_data: make(map[string]*PeriodData),
			cache_kline_data:  make(map[string]map[int]*treemap.Map),
			IsTest:            false,
			msclient:          marketservice.NewMarketService(zrpc.MustNewClient(ctx.Config.RpcConfig)),
			kline_cache:       datastruct.NewKlineCache(&ctx.Config.CacheConfig),
		}
	} else {
		return &DataEngine{
			RecvDataChan:      recvDataChan,
			ctx:               ctx,
			cache_period_data: make(map[string]*PeriodData),
			cache_kline_data:  make(map[string]map[int]*treemap.Map),
			IsTest:            false,
			kline_cache:       datastruct.NewKlineCache(&ctx.Config.CacheConfig),
		}
	}
}

func (d *DataEngine) Start() {
	logx.Infof("DataEngine Start!")

	d.StartListenRecvdata()
}

func (d *DataEngine) SetNextWorker(next_worker worker.WorkerI) {
	d.next_worker = next_worker
}

func (d *DataEngine) SetTestFlag(value bool) {
	d.IsTest = value
}

func (d *DataEngine) UpdateMeta(symbols []string) {
	for _, symbol := range symbols {
		if _, ok := d.cache_period_data[symbol]; !ok {
			err := d.InitPeriodDara(symbol)
			if err != nil {
				logx.Errorf("UpdateMeta error: %+v", err)
			}
		}
	}
}

func (a *DataEngine) InitPeriodDara(symbol string) error {
	defer util.CatchExp("InitPeriodDara " + symbol)
	logx.Slowf("Init PeriodData: %s", symbol)

	a.cache_period_data[symbol] = NewPeriodData(symbol)

	if a.IsTest {
		req_hist_info := &datastruct.ReqHistKline{
			Symbol:    symbol,
			Exchange:  datastruct.BCTS_EXCHANGE,
			Count:     datastruct.MIN_PER_DAY,
			Frequency: datastruct.NANO_PER_MIN,
		}

		Klines := datastruct.GetTestHistKline(req_hist_info)

		a.cache_period_data[symbol].UpdateWithKlines(Klines)

		return nil
	}

	end_time_nanos := uint64(util.TimeMinuteNanos())
	start_time_nanos := uint64(end_time_nanos - datastruct.NANO_PER_DAY)

	req_hist_info := &marketservice.ReqHishKlineInfo{
		Symbol:    symbol,
		Exchange:  datastruct.BCTS_EXCHANGE,
		StartTime: start_time_nanos,
		EndTime:   end_time_nanos,
		Count:     datastruct.MIN_PER_DAY,
		Frequency: datastruct.NANO_PER_MIN,
	}

	hist_klines, err := a.msclient.RequestHistKlineData(context.Background(), req_hist_info)

	if err != nil || hist_klines == nil {
		fmt.Printf("err %+v \n", err)
		logx.Errorf("ReqHistKline Err: %+v\n", err)
		return err
	}

	logx.Slowf("Init Period HistKline: %s", marketservice.HistKlineString(hist_klines))

	a.cache_period_data[symbol].UpdateWithPbKlines(hist_klines)

	logx.Slowf("Symbol Meta Info: %+v", a.cache_period_data[symbol].String())

	return nil

}

func (a *DataEngine) StartListenRecvdata() {
	logx.Info("[S] DataEngine start_listen_recvdata")
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
	logx.Info("[S] DataEngine start_receiver Over!")
}

func catch_depth_exp(depth *datastruct.DepthQuote) {
	errMsg := recover()
	if errMsg != nil {
		// fmt.Println("This is catch_exp func")
		logx.Errorf("catch_exp depth:  %+v\n", depth.String(3))
		logx.Errorf("errMsg: %+v \n", errMsg)

		logx.Infof("catch_exp depth:  %+v\n", depth.String(3))
		logx.Infof("errMsg: %+v \n", errMsg)
		// fmt.Println(errMsg)
	}
}

func (d *DataEngine) process_depth(depth *datastruct.DepthQuote) error {

	defer catch_depth_exp(depth)

	// depth.Time = depth.Time / datastruct.NANO_PER_SECS

	d.depth_cache_map.Store(depth.Symbol, depth)

	// logx.Statf("Rcv Depth: %s", depth.String(3))

	d.PublishDepth(depth, nil)
	return nil
}

func catch_kline_exp(kline *datastruct.Kline) {
	errMsg := recover()
	if errMsg != nil {
		// fmt.Println("This is catch_exp func")
		logx.Errorf("catch_exp kline:  %+v\n", kline.String())
		logx.Errorf("errMsg: %+v \n", errMsg)

		logx.Infof("catch_exp kline:  %+v\n", kline.String())
		logx.Infof("errMsg: %+v \n", errMsg)
		// fmt.Println(errMsg)
	}
}

//UnTest
func (d *DataEngine) process_kline(kline *datastruct.Kline) error {
	defer catch_kline_exp(kline)

	logx.Slowf("[RK]: %s", kline.FullString())

	d.InitPeriodDaraMain(kline.Symbol)

	d.cache_period_data[kline.Symbol].UpdateWithKline(kline)

	pubklines, err := d.kline_cache.UpdateAllKline(kline)

	if !kline.IsHistory() {
		trade := datastruct.NewTradeWithRealTimeKline(kline)
		d.process_trade(trade)
	}

	if err != nil {
		logx.Errorf("UpdateAllKline Failed: %+v", err)
		logx.Slowf("UpdateAllKline Failed: %+v", err)
		logx.Infof("UpdateAllKline Failed: %+v", err)
		return err
	}

	for _, tmp_kline := range pubklines {
		// logx.Slowf("[PK] %s", tmp_kline.FullString())
		d.PublishKline(tmp_kline, nil)
	}

	return nil
}

func catch_trade_exp(msg string, trade *datastruct.Trade) {
	errMsg := recover()
	if errMsg != nil {
		// fmt.Println("This is catch_exp func")
		logx.Errorf("%s catch_exp trade:  %+v", msg, trade.String())
		logx.Errorf("%s errMsg: %+v", msg, errMsg)

		logx.Infof("%s catch_exp trade:  %+v\n", msg, trade.String())
		logx.Infof("%s errMsg: %+v \n", msg, errMsg)
		// fmt.Println(errMsg)
	}
}

func (d *DataEngine) InitPeriodDaraMain(symbol string) {
	defer util.CatchExp("InitPeriodDaraMain " + symbol)

	d.cache_period_data_mutex.Lock()
	defer d.cache_period_data_mutex.Unlock()

	if _, ok := d.cache_period_data[symbol]; !ok {
		err := d.InitPeriodDara(symbol)

		if err != nil {
			logx.Errorf("process_trade error: %+v", err)
		}
		symbol_list := d.get_symbol_list()
		d.PublishSymbol(symbol_list, nil)
	}

}

//UnTest
func (d *DataEngine) process_trade(trade *datastruct.Trade) error {
	defer catch_trade_exp("process_trade", trade)

	d.trade_cache_map.Store(trade.Symbol, trade)
	logx.Slowf("[RT]: %s", trade.String())

	usd_price := trade.Price * d.GetUsdPrice(trade.Symbol)

	symbol_config := d.ctx.GetSymbolConfig(trade.Symbol)
	precision := 4
	if symbol_config != nil {
		precision = symbol_config.PricePrecision
	} else {
		logx.Slowf("%s, has no config ", trade.Symbol)
	}

	var change_data *datastruct.ChangeInfo
	if _, ok := d.cache_period_data[trade.Symbol]; !ok {
		change_data = nil
	} else {
		change_data = d.cache_period_data[trade.Symbol].GetChangeInfo(precision)
	}
	logx.Slowf("change_data: %s", change_data.String())

	rsp_trade := datastruct.RspTrade{
		TradeData:     trade,
		ChangeData:    change_data,
		UsdPrice:      usd_price,
		ReqArriveTime: util.UTCNanoTime(),
	}

	d.PublishTrade(&rsp_trade, nil)

	return nil
}

func (d *DataEngine) get_symbol_list() []string {
	var rst []string
	for key := range d.cache_period_data {
		rst = append(rst, key)
	}
	return rst
}

func catch_sub_symbol_exp(ws *net.WSInfo) {
	errMsg := recover()
	if errMsg != nil {
		fmt.Printf("catch_exp sub symbol:  %+v\n", ws)
		fmt.Printf("errMsg: %+v \n", errMsg)

		logx.Errorf("catch_exp sub symbol:  %+v\n", ws)
		logx.Errorf("errMsg: %+v \n", errMsg)

		logx.Infof("catch_exp sub symbol:  %+v\n", ws)
		logx.Infof("errMsg: %+v \n", errMsg)
	}
}

func (d *DataEngine) SubSymbol(ws *net.WSInfo) {
	defer catch_sub_symbol_exp(ws)
	symbol_list := d.get_symbol_list()

	d.next_worker.PublishSymbol(symbol_list, ws)
}

func (d *DataEngine) PublishSymbol(symbol_list []string, ws *net.WSInfo) {
	d.next_worker.PublishSymbol(symbol_list, ws)
}

func (d *DataEngine) PublishDepth(depth *datastruct.DepthQuote, ws *net.WSInfo) {
	d.next_worker.PublishDepth(depth, ws)
}

func (d *DataEngine) PublishTrade(trade *datastruct.RspTrade, ws *net.WSInfo) {

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

func catch_sub_trade_exp(symbol string, ws *net.WSInfo) {
	errMsg := recover()
	if errMsg != nil {
		fmt.Printf("catch_exp sub trade, %s, %+v\n", symbol, ws)
		fmt.Printf("errMsg: %+v \n", errMsg)

		logx.Errorf("catch_exp sub trade, %s, %+v\n", symbol, ws)
		logx.Errorf("errMsg: %+v \n", errMsg)

		logx.Infof("catch_exp sub trade, %s, %+v\n", symbol, ws)
		logx.Infof("errMsg: %+v \n", errMsg)
	}
}

func (d *DataEngine) GetUsdPrice(symbol string) float64 {
	usd_price := 1.0

	symbol_list := strings.Split(symbol, "_")

	if len(symbol_list) != 2 {
		return usd_price
	}

	if symbol_list[1] != "USD" {
		trans_symbol := symbol_list[1] + "_USD"

		if trade, ok := d.trade_cache_map.Load(trans_symbol); ok {
			tmp_trade := trade.(*datastruct.Trade)
			usd_price = tmp_trade.Price
		}

		// logx.Slowf("trans_symbol: %s, usd_price: %f", trans_symbol, usd_price)
	}

	return usd_price
}

func (d *DataEngine) SubTrade(req_trade *datastruct.ReqTrade, ws *net.WSInfo) (string, bool) {
	symbol := req_trade.Symbol
	defer catch_sub_trade_exp(symbol, ws)

	if trade_iter, ok := d.trade_cache_map.Load(symbol); ok {

		d.InitPeriodDaraMain(symbol)

		trade := trade_iter.(*datastruct.Trade)
		usd_price := trade.Price * d.GetUsdPrice(symbol)

		symbol_config := d.ctx.GetSymbolConfig(req_trade.Symbol)
		precision := 4
		if symbol_config != nil {
			precision = symbol_config.PricePrecision
		}

		var change_data *datastruct.ChangeInfo
		if _, ok := d.cache_period_data[symbol]; !ok {
			change_data = nil
		} else {
			change_data = d.cache_period_data[symbol].GetChangeInfo(precision)
		}

		rsp_trade := datastruct.RspTrade{
			TradeData:     trade,
			ChangeData:    change_data,
			UsdPrice:      usd_price,
			ReqWSTime:     req_trade.ReqWSTime,
			ReqArriveTime: req_trade.ReqArriveTime,
		}

		go d.PublishTrade(&rsp_trade, ws)

		return "", true
	} else {
		logx.Errorf("trade %s not cached", symbol)
		return fmt.Sprintf("trade %s not cached", symbol), false
	}
}

func catch_sub_depth_exp(symbol string, ws *net.WSInfo) {
	errMsg := recover()
	if errMsg != nil {
		fmt.Printf("catch_exp sub depth, %s, %+v\n", symbol, ws)
		fmt.Printf("errMsg: %+v \n", errMsg)

		logx.Errorf("catch_exp sub depth, %s, %+v\n", symbol, ws)
		logx.Errorf("errMsg: %+v \n", errMsg)

		logx.Infof("catch_exp sub depth, %s, %+v\n", symbol, ws)
		logx.Infof("errMsg: %+v \n", errMsg)
	}
}

func (d *DataEngine) SubDepth(symbol string, ws *net.WSInfo) (string, bool) {
	defer catch_sub_depth_exp(symbol, ws)

	if depth, ok := d.depth_cache_map.Load(symbol); ok {
		d.PublishDepth(depth.(*datastruct.DepthQuote), ws)
		return "", true
	} else {
		logx.Errorf("depth %s not cached", symbol)
		return fmt.Sprintf("depth %s not cached", symbol), false
	}
}

// Undo
// Untest
func (d *DataEngine) GetDBKlinesByCount(symbol string, resolution uint64, count int) []*datastruct.Kline {
	util.CatchExp(fmt.Sprintf("DataEngine GetDBKlinesByCount %s,%d,%d", symbol, resolution, count))
	var rst []*datastruct.Kline

	req_hist_info := &marketservice.ReqHishKlineInfo{
		Symbol:    symbol,
		Exchange:  datastruct.BCTS_EXCHANGE,
		Count:     uint32(count),
		Frequency: resolution,
	}

	logx.Infof("req_hist_info: %+v", req_hist_info)

	hist_klines, err := d.msclient.RequestHistKlineData(context.Background(), req_hist_info)
	if err != nil {
		logx.Errorf("GetHistData Failed: %+v, %+v\n", req_hist_info, err)
		logx.Slowf("GetHistData Failed: %+v, %+v\n", req_hist_info, err)
		logx.Infof("GetHistData Failed: %+v, %+v\n", req_hist_info, err)
		return nil
	}

	rst = TrasPbKlines(hist_klines.KlineData)

	return rst
}

// Undo
// Untest
func (d *DataEngine) GetDBKlinesByTime(symbol string, resolution uint64, start_time int64, end_time int64) []*datastruct.Kline {
	util.CatchExp(fmt.Sprintf("DataEngine GetDBKlinesByTime %s,%d,%d~%d", symbol, resolution, start_time, end_time))
	var rst []*datastruct.Kline

	req_hist_info := &marketservice.ReqHishKlineInfo{
		Symbol:    symbol,
		Exchange:  datastruct.BCTS_EXCHANGE,
		StartTime: uint64(start_time),
		EndTime:   uint64(end_time),
		Frequency: resolution,
	}

	logx.Infof("req_hist_info: %+v", req_hist_info)

	hist_klines, err := d.msclient.RequestHistKlineData(context.Background(), req_hist_info)
	if err != nil {
		logx.Errorf("GetHistData Failed: %+v, %+v\n", req_hist_info, err)
		logx.Slowf("GetHistData Failed: %+v, %+v\n", req_hist_info, err)
		logx.Infof("GetHistData Failed: %+v, %+v\n", req_hist_info, err)
		return nil
	}

	rst = TrasPbKlines(hist_klines.KlineData)
	return rst
}

// Untest
func (d *DataEngine) GetKlinesByCount(symbol string, resolution uint64, count int) []*datastruct.Kline {
	defer util.CatchExp("DataEngine GetKlinesByCount")

	rst := d.kline_cache.GetKlinesByCount(symbol, resolution, count, false, true)

	if rst == nil {

		req_count := util.MaxInt(count, d.kline_cache.Config.Count)

		logx.Slowf("ReqDB: %s, %d, %d \n", symbol, resolution, req_count)

		db_klines := d.GetDBKlinesByCount(symbol, resolution, req_count)
		datastruct.OutputDetailHistKlines(db_klines)

		d.kline_cache.ReleaseInputKlines(db_klines, symbol, resolution)
	}

	rst = d.kline_cache.GetKlinesByCount(symbol, resolution, count, false, false)

	return rst
}

// Untest
func (d *DataEngine) GetKlinesByTime(symbol string, resolution uint64, start_time int64, end_time int64) []*datastruct.Kline {
	defer util.CatchExp("DataEngine GetKlinesByTime")

	rst := d.kline_cache.GetKlinesByTime(symbol, resolution, start_time, end_time, true)

	if rst == nil {
		db_klines := d.GetDBKlinesByTime(symbol, resolution, start_time, end_time)
		d.kline_cache.InitWithHistKlines(db_klines, symbol, resolution)
	}

	rst = d.kline_cache.GetKlinesByTime(symbol, resolution, start_time, end_time, false)

	return rst
}

func (d *DataEngine) GetHistKlineData(req_kline_info *datastruct.ReqHistKline) *datastruct.RspHistKline {
	var klines []*datastruct.Kline

	if d.IsTest {
		tmp := datastruct.GetTestHistKline(req_kline_info)
		fmt.Println(tmp)
	} else {
		klines = d.GetKlinesByCount(req_kline_info.Symbol, req_kline_info.Frequency, int(req_kline_info.Count))

		logx.Infof("req_hist_info: %+v", req_kline_info)
	}

	return &datastruct.RspHistKline{
		ReqInfo: req_kline_info,
		Klines:  klines,
	}
}

func (d *DataEngine) SubKline(req_kline_info *datastruct.ReqHistKline, ws *net.WSInfo) (string, bool) {
	defer util.CatchExp(fmt.Sprintf("SubKline %s, %s", req_kline_info.String(), ws.String()))

	logx.Slowf("[DE] SubK %s,", req_kline_info.String())

	rst := d.GetHistKlineData(req_kline_info)

	if rst != nil {
		logx.Slowf("[DE] HistK, rsl:%d, %s", req_kline_info.Frequency, datastruct.HistKlineSimpleTime(rst.Klines))

		d.next_worker.PublishHistKline(rst, ws)

		return "", true
	} else {
		logx.Errorf("kline %s get hist data failed! ", req_kline_info.String())
		return fmt.Sprintf("depth %s get hist data failed!", req_kline_info.String()), false
	}

}

func (f *DataEngine) UnSubTrade(symbol string, ws *net.WSInfo) {

}

func (f *DataEngine) UnSubDepth(symbol string, ws *net.WSInfo) {
}

func (f *DataEngine) UnSubKline(req_kline_info *datastruct.ReqHistKline, ws *net.WSInfo) {
}

// Untest
func (d *DataEngine) GetHistKlineDataNew(req_kline_info *datastruct.ReqHistKline) *datastruct.RspHistKline {
	// defer util.CatchExp(fmt.Sprintf("GetHistKlineDataNew %s", req_kline_info.String()))
	// var ori_klines []*datastruct.Kline

	// if req_kline_info.Count > 0 {
	// 	ori_klines = d.GetKlinesByCount(req_kline_info.Symbol, req_kline_info.Frequency, int(req_kline_info.Count))
	// } else if req_kline_info.StartTime > 0 && req_kline_info.EndTime > req_kline_info.StartTime {
	// 	ori_klines = d.GetKlinesByTime(req_kline_info.Symbol, req_kline_info.Frequency,
	// 		int64(req_kline_info.StartTime), int64(req_kline_info.EndTime))
	// } else {
	// 	logx.Errorf("error req_hist_kline %s", req_kline_info.String())
	// }

	// trans_kline := datastruct.TransSliceKlines(ori_klines)

	// return &datastruct.RspHistKline{
	// 	ReqInfo: req_kline_info,
	// 	Klines:  trans_kline,
	// }

	return nil
}

func (d *DataEngine) TrasOriKlineData(req_kline_info *datastruct.ReqHistKline, ori_klines *treemap.Map) (*treemap.Map, bool) {
	rst := treemap.NewWith(utils.Int64Comparator)
	resolution := req_kline_info.Frequency

	if ori_klines.Size() == 0 {
		return rst, false
	}

	iter := ori_klines.Iterator()
	iter.First()

	cache_kline := iter.Value().(*datastruct.Kline)

	// Error
	if !datastruct.IsNewKlineStart(cache_kline, resolution) {
		cache_kline.Time = datastruct.GetLastStartTime(cache_kline.Time, int64(resolution))
	}

	iter.Begin()
	var pub_kline *datastruct.Kline
	for iter.Next() {
		cur_kline := iter.Value().(*datastruct.Kline)

		if datastruct.IsOldKlineEnd(cur_kline, resolution) {
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
		} else if datastruct.IsNewKlineStart(cur_kline, resolution) {
			cache_kline = cur_kline
			cache_kline.Resolution = resolution
		} else {
			cache_kline.Close = cur_kline.Close
			cache_kline.Low = util.MinFloat64(cache_kline.Low, cur_kline.Low)
			cache_kline.High = util.MaxFloat64(cache_kline.High, cur_kline.High)
			cache_kline.Volume = cache_kline.Volume + cur_kline.Volume
		}
	}

	is_last_complete := true
	if cache_kline.Time != pub_kline.Time {
		pub_kline = datastruct.NewKlineWithKline(cache_kline)
		rst.Put(pub_kline.Time, pub_kline)
		is_last_complete = false
	}

	return rst, is_last_complete
}

func (d *DataEngine) UpdateCacheKlinesWithHist(klines *treemap.Map) {

}

func (d *DataEngine) GetHistKlineDataBak(req_kline_info *datastruct.ReqHistKline) *datastruct.RspHistKline {

	tmp := treemap.NewWith(utils.Int64Comparator)

	if d.IsTest {
		// tmp := datastruct.GetTestHistKline(req_kline_info)
	} else {

		// rate := req_kline_info.Frequency / datastruct.SECS_PER_MIN

		req_hist_info := &marketservice.ReqHishKlineInfo{
			Symbol:    req_kline_info.Symbol,
			Exchange:  req_kline_info.Exchange,
			StartTime: req_kline_info.StartTime,
			EndTime:   req_kline_info.EndTime,
			Count:     req_kline_info.Count,
			Frequency: req_kline_info.Frequency,
		}

		if req_kline_info.Frequency%datastruct.NANO_PER_MIN != 0 {
			return nil
		}

		logx.Infof("req_hist_info: %+v", req_kline_info)

		hist_klines, err := d.msclient.RequestHistKlineData(context.Background(), req_hist_info)

		if err != nil {
			logx.Errorf("GetHistData Failed: %+v, %+v\n", req_hist_info, err)
			return nil
		} else {
			// logx.Infof("Original hist_klines : %+v", hist_klines.KlineData)
		}

		for _, pb_kline := range hist_klines.KlineData {
			kline := marketservice.NewKlineWithPbKline(pb_kline)
			if kline == nil {
				continue
			}

			tmp.Put(kline.Time, kline)
		}
	}

	d.UpdateCacheKlinesWithHist(tmp)

	// logx.Slowf("\nOriKlineTime: %s", datastruct.HistKlineTimeList(tmp, 3))
	// // fmt.Printf("OriKlineTime: %s \n", datastruct.HistKlineTimeList(tmp))

	// trans_kline, is_last_complete := d.TrasOriKlineData(req_kline_info, tmp)

	// logx.Slowf("\nTransKlineTime: %s", datastruct.HistKlineTimeList(trans_kline, 3))
	// // fmt.Printf("TransKlineTime: %s\n", datastruct.HistKlineTimeList(trans_kline))

	// return &datastruct.RspHistKline{
	// 	ReqInfo:        req_kline_info,
	// 	Klines:         trans_kline,
	// 	IsLastComplete: is_last_complete,
	// }

	return nil
}
