package binance

import (
	"exterior-interactor/pkg/exchangeapi/extools"
	"exterior-interactor/pkg/timeutils"
	"exterior-interactor/pkg/xmath"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
	"github.com/shopspring/decimal"
	"sync"

	"exterior-interactor/pkg/exchangeapi/exchanges/binance/api/binancespot"
	"exterior-interactor/pkg/exchangeapi/exmodel"
	"exterior-interactor/pkg/httptools"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
	"strings"
	"time"
)

const (
	depthLimit       = 200
	depthOutputLimit = 100
)

const (
	depthOutputDuration = time.Millisecond * 50
)

type (
	depthManager struct {
		extools.SymbolManager
		api      *NativeApi
		spot     httptools.AutoWsSubscriber
		store    map[exmodel.StdSymbol]*depthUnit
		outputCh chan *exmodel.StreamDepth
	}

	depthUnit struct {
		mutex          sync.Mutex
		spotInputCh    chan *binancespot.WsDiffDepth
		spotDepthCache []*binancespot.WsDiffDepth // 缓存
		api            *NativeApi
		exchange       exmodel.Exchange
		symbol         *exmodel.Symbol
		lastUpdateId   int64 // 从 rest 获取到的 lastUpdateId
		startUpdateId  int64 // 从 ws 获取到的 startUpdateId
		endUpdateId    int64 // 从 ws 获取到的 endUpdateId
		lastUpdateTime time.Time
		asks           *treemap.Map // 价格从高到低, key 必须为 float64
		bids           *treemap.Map // 价格从低到高, key 必须为 float64
		outputCh       chan *exmodel.StreamDepth
	}
)

func NewDepthManager(mgr extools.SymbolManager, api *NativeApi) extools.DepthManager {
	res, err := api.SpotApi.GetStreamDiffDepth()
	if err != nil {
		panic(err)
	}

	depthMgr := &depthManager{
		SymbolManager: mgr,
		api:           api,
		spot:          res,
		store:         map[exmodel.StdSymbol]*depthUnit{},
		outputCh:      make(chan *exmodel.StreamDepth, 5096),
	}

	go depthMgr.run()
	return depthMgr
}

func newDepthUnit(symbol *exmodel.Symbol, api *NativeApi, outputCh chan *exmodel.StreamDepth) *depthUnit {
	info := &depthUnit{
		mutex:          sync.Mutex{},
		api:            api,
		spotInputCh:    make(chan *binancespot.WsDiffDepth, 1024),
		spotDepthCache: make([]*binancespot.WsDiffDepth, 0),
		exchange:       Name,
		symbol:         symbol,
		lastUpdateTime: time.Now(),
		asks:           treemap.NewWith(utils.Float64Comparator), // 默认按 key 从小到大排序
		bids:           treemap.NewWith(utils.Float64Comparator), // 默认按 key 从小到大排序
		outputCh:       outputCh,
	}
	go info.run()
	info.spotReqRest()
	return info
}

func (o *depthUnit) run() {
	timeutils.Every(depthOutputDuration, func() {
		d, ok := o.generateDepth()
		if ok {
			o.outputCh <- d
		}
	})

	for {
		select {
		case depth := <-o.spotInputCh:
			dataLost := o.spotUpdate(depth)
			if dataLost {
				o.spotReqRest()
			}
		}
	}
}

func (o *depthUnit) generateDepth() (*exmodel.StreamDepth, bool) {
	defer o.mutex.Unlock()
	o.mutex.Lock()
	//fmt.Println((o.lastUpdateTime).Format(time.RFC3339Nano), "%%%%%%%%%%%%%%")
	//start:=time.Now()
	d := &exmodel.StreamDepth{
		Exchange:  o.exchange,
		Time:      o.lastUpdateTime,
		LocalTime: time.Now(),
		Symbol:    o.symbol,
		Asks:      make([][2]string, 0, depthOutputLimit),
		Bids:      make([][2]string, 0, depthOutputLimit),
	}

	asksIt := o.asks.Iterator()
	for asksIt.Begin(); asksIt.Next(); {
		if len(d.Asks) == depthOutputLimit {
			break
		}
		d.Asks = append(d.Asks, [2]string{fmt.Sprint(asksIt.Key()), asksIt.Value().(string)})
	}

	bidsIt := o.bids.Iterator()
	for bidsIt.End(); bidsIt.Prev(); {
		if len(d.Bids) == depthOutputLimit {
			break
		}

		d.Bids = append(d.Bids, [2]string{fmt.Sprint(bidsIt.Key()), bidsIt.Value().(string)})
	}

	if len(d.Asks) == 0 && len(d.Bids) == 0 {
		return nil, false
	}
	//fmt.Println("cost: ",time.Now().Sub(start))
	return d, true
}

func (o *depthUnit) spotUpdate(depth *binancespot.WsDiffDepth) (dataLost bool) {
	defer o.mutex.Unlock()
	o.mutex.Lock()

	if o.lastUpdateId == 0 { // 还未通过 rest 获取数据
		o.spotDepthCache = append(o.spotDepthCache, depth)
		return false
	}

	// 已通过 rest 获取数据，o.lastUpdateId 已有数据
	for _, diffDepth := range o.spotDepthCache {
		if diffDepth.Data.EndId <= o.lastUpdateId { // 丢弃
			continue
		}

		// 缓存中是否有 第一条 数据
		if diffDepth.Data.StartId <= o.lastUpdateId+1 && diffDepth.Data.EndId >= o.lastUpdateId+1 {
			o.updateAsksBids(diffDepth.Data.A, diffDepth.Data.B)
			o.lastUpdateTime = time.UnixMilli(diffDepth.Data.E1)
			o.startUpdateId = diffDepth.Data.StartId
			o.endUpdateId = diffDepth.Data.EndId
			continue
		}

		if diffDepth.Data.StartId != o.endUpdateId+1 {
			return true // 缓存中有丢包
		}

		o.updateAsksBids(diffDepth.Data.A, diffDepth.Data.B)
		o.lastUpdateTime = time.UnixMilli(diffDepth.Data.E1)
		o.startUpdateId = diffDepth.Data.StartId
		o.endUpdateId = diffDepth.Data.EndId
	}

	if len(o.spotDepthCache) != 0 {
		o.spotDepthCache = o.spotDepthCache[0:0] // 清空缓存
	}

	if depth.Data.EndId <= o.lastUpdateId { // 丢弃
		return false
	}

	if depth.Data.StartId <= o.lastUpdateId+1 && depth.Data.EndId >= o.lastUpdateId+1 {
		o.updateAsksBids(depth.Data.A, depth.Data.B)
		o.lastUpdateTime = time.UnixMilli(depth.Data.E1)
		o.startUpdateId = depth.Data.StartId
		o.endUpdateId = depth.Data.EndId
		return false
	}

	if depth.Data.StartId != o.endUpdateId+1 {
		return true // 有丢包
	}

	o.updateAsksBids(depth.Data.A, depth.Data.B)
	o.lastUpdateTime = time.UnixMilli(depth.Data.E1)
	o.startUpdateId = depth.Data.StartId
	o.endUpdateId = depth.Data.EndId
	return false
}

// spotReqRest 请求 rest 接口，获取全量数据
func (o *depthUnit) spotReqRest() {
	for {
		req := binancespot.GetDepthReq{
			Symbol: o.symbol.ExFormat,
			Limit:  depthLimit,
		}
		rsp, err := o.api.SpotApi.GetDepth(req)
		if err != nil {
			logx.Errorf("GetDepth err:%v, req:%+v", err, req)
			time.Sleep(time.Second * 3)
			continue
		}

		o.mutex.Lock()
		o.asks.Clear()
		o.bids.Clear()
		o.updateAsksBids(rsp.Asks, rsp.Bids)
		o.lastUpdateId = rsp.LastUpdateId
		o.lastUpdateTime = time.Now()
		o.mutex.Unlock()
		break
	}
}

func (o *depthUnit) updateAsksBids(asks, bids [][2]string) {
	for _, item := range asks {
		price, _ := xmath.MustDecimal(item[0]).Float64()
		volume := xmath.MustDecimal(item[1])
		if volume.Equal(decimal.Zero) {
			o.asks.Remove(price)
		} else {
			o.asks.Put(price, volume.String())
		}
	}

	for _, item := range bids {
		price, _ := xmath.MustDecimal(item[0]).Float64()
		volume := xmath.MustDecimal(item[1])
		if volume.Equal(decimal.Zero) {
			o.bids.Remove(price)
		} else {
			o.bids.Put(price, volume.String())
		}
	}

}

func (o *depthManager) run() {
	spotCh := o.spot.ReadCh()
	for {
		select {
		case msg := <-spotCh:
			depth := msg.(*binancespot.WsDiffDepth)
			symbol, err := o.SymbolManager.Convert(depth.Data.S, o.api.SpotApi.ApiType)
			if err != nil {
				logx.Errorf("binance spot WsDiffDepth can't parse symbol, depth: %+v", *depth)
				continue
			}
			info, ok := o.store[symbol.StdSymbol]
			if !ok { // 还未初始化完成，先丢弃
				continue
			}

			info.spotInputCh <- depth
		}
	}

}

func (o *depthManager) Sub(symbols ...exmodel.StdSymbol) {
	var valid []*exmodel.Symbol

	for _, s := range symbols {
		if _, ok := o.store[s]; ok {
			continue
		}

		r, err := o.SymbolManager.GetSymbol(s)
		if err != nil {
			logx.Infof("%s not support symbol:%s", Name, s)
		} else {
			valid = append(valid, r)
		}
	}

	//spotTopics, coinFutureTopics, usdtFutureTopics := o.parseSymbolToTopic(valid...)
	spotTopics, _, _ := o.parseSymbolToTopic(valid...)
	o.spot.Sub(spotTopics...)
	//o.coinFuture.Sub(coinFutureTopics...)
	//o.usdtFuture.Sub(usdtFutureTopics...)

	_, _ = mr.MapReduce(func(source chan<- interface{}) {
		for _, symbol := range valid {
			source <- symbol
		}
	}, func(item interface{}, writer mr.Writer, cancel func(error)) {
		symbol := item.(*exmodel.Symbol)
		info := newDepthUnit(symbol, o.api, o.outputCh)
		o.store[symbol.StdSymbol] = info
		writer.Write(struct{}{})
	}, func(pipe <-chan interface{}, writer mr.Writer, cancel func(error)) {
		writer.Write(struct{}{})
	})

}

func (o *depthManager) parseSymbolToTopic(symbols ...*exmodel.Symbol) (
	spotTopics, coinFutureTopics, usdtFutureTopics []string) {
	for _, symbol := range symbols {
		topic := fmt.Sprintf("%s@depth@100ms", strings.ToLower(symbol.ExFormat))
		switch symbol.Type {
		case exmodel.SymbolTypeSpot:
			spotTopics = append(spotTopics, topic)
		case exmodel.SymbolTypeUsdtPerp, exmodel.SymbolTypeUsdtDelivery:
			usdtFutureTopics = append(usdtFutureTopics, topic)
		case exmodel.SymbolTypeCoinPerp, exmodel.SymbolTypeCoinDelivery:
			coinFutureTopics = append(coinFutureTopics, topic)
		default:
			logx.Errorf("not support sub depth with symbol:%s", symbol.StdSymbol)
		}
	}
	return
}

func (o *depthManager) OutputCh() <-chan *exmodel.StreamDepth {
	return o.outputCh
}
