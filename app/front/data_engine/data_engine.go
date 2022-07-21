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
		}
	} else {
		return &DataEngine{
			RecvDataChan:      recvDataChan,
			ctx:               ctx,
			cache_period_data: make(map[string]*PeriodData),
			cache_kline_data:  make(map[string]map[int]*treemap.Map),
			IsTest:            false,
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
	util.CatchExp("InitPeriodDara " + symbol)
	logx.Infof("Init PeriodData: %s", symbol)
	a.cache_period_data[symbol] = &PeriodData{
		Symbol:                symbol,
		TimeNanos:             datastruct.NANO_PER_DAY,
		Count:                 0,
		MaxTime:               0,
		MinTime:               0,
		time_cache_data:       treemap.NewWith(utils.Int64Comparator),
		high_price_cache_data: NewSortedList(true),
		low_price_cache_data:  NewSortedList(false),
	}

	if a.IsTest {
		req_hist_info := &datastruct.ReqHistKline{
			Symbol:    symbol,
			Exchange:  datastruct.BCTS_EXCHANGE,
			Count:     datastruct.MIN_PER_DAY,
			Frequency: datastruct.SECS_PER_MIN,
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
		Frequency: datastruct.SECS_PER_MIN,
	}

	hist_klines, err := a.msclient.RequestHistKlineData(context.Background(), req_hist_info)

	if err != nil || hist_klines == nil {
		fmt.Printf("err %+v \n", err)
		logx.Errorf("ReqHistKline Err: %+v\n", err)
		return err
	}

	logx.Infof("Init Period HistKline: %s", marketservice.HistKlineString(hist_klines))

	// fmt.Printf("Rst: %+v \n", hist_klines)

	a.cache_period_data[symbol].UpdateWithPbKlines(hist_klines)

	logx.Infof("Symbol Meta Info: %+v", a.cache_period_data[symbol].String())

	return nil

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

func (d *DataEngine) process_kline(kline *datastruct.Kline) error {
	defer catch_kline_exp(kline)
	// kline.Time = kline.Time / datastruct.NANO_PER_SECS

	// logx.Statf("Rcv kline: %s", kline.String())

	d.cache_period_data_mutex.Lock()
	if _, ok := d.cache_period_data[kline.Symbol]; !ok {
		err := d.InitPeriodDara(kline.Symbol)

		if err != nil {
			logx.Errorf("process_kline error: %+v", err)
			return err
		}

		symbol_list := d.get_symbol_list()
		d.PublishSymbol(symbol_list, nil)
	}
	d.cache_period_data_mutex.Unlock()

	d.cache_period_data[kline.Symbol].UpdateWithKline(kline)

	d.PublishKline(kline, nil)

	// d.PublishChangeinfo(d.cache_period_data[kline.Symbol].GetChangeInfo(), nil)

	return nil
}

func catch_trade_exp(trade *datastruct.Trade) {
	errMsg := recover()
	if errMsg != nil {
		// fmt.Println("This is catch_exp func")
		logx.Errorf("catch_exp trade:  %+v", trade.String())
		logx.Errorf("errMsg: %+v", errMsg)

		logx.Infof("catch_exp trade:  %+v\n", trade.String())
		logx.Infof("errMsg: %+v \n", errMsg)
		// fmt.Println(errMsg)
	}
}

func (d *DataEngine) process_trade(trade *datastruct.Trade) error {
	defer catch_trade_exp(trade)

	d.trade_cache_map.Store(trade.Symbol, trade)

	d.cache_period_data_mutex.Lock()
	if _, ok := d.cache_period_data[trade.Symbol]; !ok {
		err := d.InitPeriodDara(trade.Symbol)

		if err != nil {
			logx.Errorf("process_trade error: %+v", err)
			return err
		}

		symbol_list := d.get_symbol_list()
		d.PublishSymbol(symbol_list, nil)
	}
	d.cache_period_data_mutex.Unlock()

	d.cache_period_data[trade.Symbol].UpdateWithTrade(trade)

	usd_price := trade.Price * d.GetUsdPrice(trade.Symbol)

	symbol_config := d.ctx.GetSymbolConfig(trade.Symbol)
	precision := 4

	if symbol_config != nil {
		precision = symbol_config.PricePrecision
	}

	rsp_trade := datastruct.RspTrade{
		TradeData:     trade,
		ChangeData:    d.cache_period_data[trade.Symbol].GetChangeInfo(precision),
		UsdPrice:      usd_price,
		ReqArriveTime: util.UTCNanoTime(),
	}

	d.PublishTrade(&rsp_trade, nil)

	return nil
}

func (d *DataEngine) get_symbol_list() []string {
	// d.cache_period_data_mutex.Lock()
	// defer d.cache_period_data_mutex.Unlock()

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

	if trade, ok := d.trade_cache_map.Load(symbol); ok {
		d.cache_period_data_mutex.Lock()
		var init_period_err error
		if _, ok := d.cache_period_data[symbol]; !ok {
			init_period_err = d.InitPeriodDara(symbol)

			symbol_list := d.get_symbol_list()
			d.PublishSymbol(symbol_list, nil)
		}

		if d.GetUsdPrice(symbol) != 0.0 {
			trade_value := trade.(*datastruct.Trade)
			usd_price := trade_value.Price * d.GetUsdPrice(symbol)

			if init_period_err != nil {
				logx.Errorf("SubTrade error: %+v", init_period_err)
				rsp_trade := datastruct.RspTrade{
					TradeData:     trade_value,
					ChangeData:    nil,
					UsdPrice:      usd_price,
					ReqWSTime:     req_trade.ReqWSTime,
					ReqArriveTime: req_trade.ReqArriveTime,
				}
				// logx.Slowf("[DE Trade] %s, ReqWSTime %d ns", req_trade.Symbol, rsp_trade.ReqWSTime)
				go d.PublishTrade(&rsp_trade, ws)
			} else {
				symbol_config := d.ctx.GetSymbolConfig(req_trade.Symbol)
				precision := 4

				if symbol_config != nil {
					precision = symbol_config.PricePrecision
				}

				rsp_trade := datastruct.RspTrade{
					TradeData:     trade_value,
					ChangeData:    d.cache_period_data[symbol].GetChangeInfo(precision),
					UsdPrice:      usd_price,
					ReqWSTime:     req_trade.ReqWSTime,
					ReqArriveTime: req_trade.ReqArriveTime,
				}

				// logx.Slowf("[DE Trade] %s, ReqWSTime %d ns", req_trade.Symbol, rsp_trade.ReqWSTime)
				go d.PublishTrade(&rsp_trade, ws)
			}
		}

		d.cache_period_data_mutex.Unlock()
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

func (d *DataEngine) GetHistKlineData(req_kline_info *datastruct.ReqHistKline) *datastruct.RspHistKline {

	tmp := treemap.NewWith(utils.Int64Comparator)

	if d.IsTest {
		tmp = datastruct.GetTestHistKline(req_kline_info)

		// return &datastruct.RspHistKline{
		// 	Klines:  datastruct.GetTestHistKline(req_kline_info),
		// 	ReqInfo: req_kline_info,
		// }
	} else {
		rate := req_kline_info.Frequency / datastruct.SECS_PER_MIN

		req_hist_info := &marketservice.ReqHishKlineInfo{
			Symbol:    req_kline_info.Symbol,
			Exchange:  req_kline_info.Exchange,
			StartTime: req_kline_info.StartTime,
			EndTime:   req_kline_info.EndTime,
			Count:     req_kline_info.Count * rate,
			Frequency: datastruct.SECS_PER_MIN,
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

	logx.Slowf("\nOriKlineTime: %s", datastruct.HistKlineTimeList(tmp, 3))
	// fmt.Printf("OriKlineTime: %s \n", datastruct.HistKlineTimeList(tmp))

	trans_kline, is_last_complete := d.TrasOriKlineData(req_kline_info, tmp)

	logx.Slowf("\nTransKlineTime: %s", datastruct.HistKlineTimeList(trans_kline, 3))
	// fmt.Printf("TransKlineTime: %s\n", datastruct.HistKlineTimeList(trans_kline))

	return &datastruct.RspHistKline{
		ReqInfo:        req_kline_info,
		Klines:         trans_kline,
		IsLastComplete: is_last_complete,
	}
}

func catch_trasoriklinedaata_exp(req_kline_info *datastruct.ReqHistKline) {
	errMsg := recover()
	if errMsg != nil {
		fmt.Printf("catch_trasoriklinedaata_exp sub depth,  %+v\n", req_kline_info)
		fmt.Printf("errMsg: %+v \n", errMsg)

		logx.Errorf("catch_trasoriklinedaata_exp sub depth, %+v\n", req_kline_info)
		logx.Errorf("errMsg: %+v \n", errMsg)

		logx.Infof("catch_trasoriklinedaata_exp sub depth, %+v\n", req_kline_info)
		logx.Infof("errMsg: %+v \n", errMsg)
	}
}

func (d *DataEngine) TrasOriKlineData(req_kline_info *datastruct.ReqHistKline, ori_klines *treemap.Map) (*treemap.Map, bool) {
	defer catch_trasoriklinedaata_exp(req_kline_info)

	rst := treemap.NewWith(utils.Int64Comparator)
	resolution := int(req_kline_info.Frequency)

	if ori_klines.Size() == 0 {
		return rst, false
	}

	iter := ori_klines.Iterator()
	iter.First()

	cache_kline := iter.Value().(*datastruct.Kline)

	if !datastruct.IsNewKlineStart(cache_kline, int64(resolution)) {
		cache_kline.Time = cache_kline.Time - cache_kline.Time%(int64(resolution)*datastruct.NANO_PER_SECS)
	}

	iter.Begin()
	var pub_kline *datastruct.Kline
	for iter.Next() {
		cur_kline := iter.Value().(*datastruct.Kline)

		if datastruct.IsOldKlineEnd(cur_kline, int64(resolution)) {
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

func catch_sub_kline_exp(req_kline_info *datastruct.ReqHistKline, ws *net.WSInfo) {
	errMsg := recover()
	if errMsg != nil {
		fmt.Printf("catch_exp sub kline, %+v, %+v\n", req_kline_info, ws)
		fmt.Printf("errMsg: %+v \n", errMsg)

		logx.Errorf("catch_exp sub kline, %+v, %+v\n", req_kline_info, ws)
		logx.Errorf("errMsg: %+v \n", errMsg)

		logx.Infof("catch_exp sub kline, %+v, %+v\n", req_kline_info, ws)
		logx.Infof("errMsg: %+v \n", errMsg)
	}
}

func (d *DataEngine) SubKline(req_kline_info *datastruct.ReqHistKline, ws *net.WSInfo) (string, bool) {
	defer catch_sub_kline_exp(req_kline_info, ws)

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
