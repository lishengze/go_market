package ftx

import (
	"exterior-interactor/pkg/exchangeapi/exchanges/ftx/ftxapi"
	"exterior-interactor/pkg/exchangeapi/exmodel"
	"exterior-interactor/pkg/exchangeapi/extools"
	"exterior-interactor/pkg/httptools"
	"exterior-interactor/pkg/xmath"
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logx"
	"hash/crc32"
	"strings"
	"sync"
	"time"
)

const (
	depthOutputLimit = 100

	checksumInterval = 20 // 20 次 update 校验一次 checksum
)

const (
	depthUnitStatusNormal    = 1
	depthUnitStatusResetting = 2
)

const (
	depthOutputDuration = time.Millisecond * 50

	depthUpdateCheckDuration = time.Second * 30 // depth 超过此时间未更新，将触发重连
)

type (
	depthManager struct {
		mutex sync.Mutex
		extools.SymbolManager
		api        *NativeApi
		subscriber httptools.AutoWsSubscriber
		store      map[exmodel.StdSymbol]*depthUnit
		outputCh   chan *exmodel.StreamDepth
	}

	depthUnit struct {
		msgSeq         int
		mutex          sync.Mutex
		inputCh        chan *ftxapi.StreamDepth
		api            *NativeApi
		subscriber     httptools.AutoWsSubscriber
		exchange       exmodel.Exchange
		symbol         *exmodel.Symbol
		lastUpdateTime time.Time
		asks           *treemap.Map //  key 必须为 float64
		bids           *treemap.Map //  key 必须为 float64
		outputCh       chan *exmodel.StreamDepth
		priceScale     int32
		status         int32
	}
)

func NewDepthManager(mgr extools.SymbolManager, api *NativeApi) extools.DepthManager {
	res, err := api.Api.GetStreamDepth()
	if err != nil {
		panic(err)
	}

	depthMgr := &depthManager{
		mutex:         sync.Mutex{},
		SymbolManager: mgr,
		api:           api,
		subscriber:    res,
		store:         map[exmodel.StdSymbol]*depthUnit{},
		outputCh:      make(chan *exmodel.StreamDepth, 5096),
	}

	go depthMgr.run()
	return depthMgr
}

func newDepthUnit(symbol *exmodel.Symbol, api *NativeApi, outputCh chan *exmodel.StreamDepth, subscriber httptools.AutoWsSubscriber) *depthUnit {
	info := &depthUnit{
		mutex:          sync.Mutex{},
		api:            api,
		inputCh:        make(chan *ftxapi.StreamDepth, 1024),
		exchange:       Name,
		symbol:         symbol,
		lastUpdateTime: time.Now(),
		subscriber:     subscriber,
		asks:           treemap.NewWith(extools.DecimalStringComparator), // 默认按 key 从小到大排序
		bids:           treemap.NewWith(extools.DecimalStringComparator), // 默认按 key 从小到大排序
		outputCh:       outputCh,
		priceScale:     0,
		status:         depthUnitStatusNormal,
	}

	s, err := decimal.NewFromString(symbol.PriceScale)
	if err != nil {
		panic(fmt.Sprintf("cat parse PriceScale, symbol:%+v \n", *symbol))
	}
	info.priceScale = -s.Exponent()
	if info.priceScale == 0 {
		info.priceScale = 1
	}

	go info.run()
	return info
}

func (o *depthUnit) run() {
	ticker := time.NewTicker(depthOutputDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			d, ok := o.generateDepth()
			if ok {
				o.outputCh <- d
			}
		case depth := <-o.inputCh:
			if depth.Type == "subscribed" {
				continue
			}
			//fmt.Println(depth.Data.Checksum,">>>>>>>>>>>>>>")
			_ = o.update(depth)
		}
	}
}

func (o *depthUnit) generateDepth() (*exmodel.StreamDepth, bool) {
	defer o.mutex.Unlock()
	o.mutex.Lock()

	if o.status != depthUnitStatusNormal {
		return nil, false
	}

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
		d.Asks = append(d.Asks, [2]string{asksIt.Key().(string), asksIt.Value().(string)})
	}

	bidsIt := o.bids.Iterator()
	for bidsIt.End(); bidsIt.Prev(); {
		if len(d.Bids) == depthOutputLimit {
			break
		}
		d.Bids = append(d.Bids, [2]string{bidsIt.Key().(string), bidsIt.Value().(string)})
	}

	if len(d.Asks) == 0 && len(d.Bids) == 0 {
		return nil, false
	}

	return d, true
}

func (o *depthUnit) checkSum() uint32 {
	/*
		The checksum operates on a string that represents the first 100 orders on the orderbook on either side. The format of the string is:

		<best_bid_price>:<best_bid_size>:<best_ask_price>:<best_ask_size>:<second_best_bid_price>:<second_best_ask_price>:...
		For example, if the orderbook was comprised of the following two bids and asks:

		bids: [[5000.5, 10], [4995.0, 5]]
		asks: [[5001.0, 6], [5002.0, 7]]
		The string would be '5005.5:10:5001.0:6:4995.0:5:5002.0:7'
	*/
	asksIt := o.asks.Iterator()
	bidsIt := o.bids.Iterator()
	bidsIt.End()
	count := 0
	var list []string
	for {
		if count == 100 { // 只取 ask bid 最多一百条
			break
		}

		isAsksItNext := asksIt.Next()
		isBidsItPrev := bidsIt.Prev()

		if isBidsItPrev {
			list = append(list, bidsIt.Key().(string), bidsIt.Value().(string))
		}

		if isAsksItNext {
			list = append(list, asksIt.Key().(string), asksIt.Value().(string))
		}

		if !isAsksItNext && !isBidsItPrev {
			break
		}

		count++
	}
	str := strings.Join(list, ":")
	//fmt.Println(str)
	//fmt.Println(len(list))
	//fmt.Println(len(o.asks.Keys()))
	//fmt.Println(len(o.bids.Keys()))
	return crc32.ChecksumIEEE([]byte(str))
}

func (o *depthUnit) update(data *ftxapi.StreamDepth) (ok bool) {
	defer o.mutex.Unlock()
	o.mutex.Lock()

	updateFn := func() {
		tsDecimal := decimal.NewFromFloat(data.Data.Time)
		microSec := tsDecimal.Mul(decimal.NewFromInt(1_000_000)).IntPart()
		o.lastUpdateTime = time.UnixMicro(microSec)
		for _, item := range data.Data.Asks {
			price := item[0].String()
			volume := item[1].String()
			if xmath.MustDecimal(volume).Equal(decimal.Zero) {
				o.asks.Remove(price)
			} else {
				o.asks.Put(price, volume)
			}
		}
		for _, item := range data.Data.Bids {
			price := item[0].String()
			volume := item[1].String()
			if xmath.MustDecimal(volume).Equal(decimal.Zero) {
				o.bids.Remove(price)
			} else {
				o.bids.Put(price, volume)
			}
		}
	}

	switch o.status {
	case depthUnitStatusNormal:
		updateFn() // 此状态下数据都要更新
		o.msgSeq++
		// 只有 update 时 验证数据
		if data.Type != "update" {
			return true
		}

		if o.msgSeq%checksumInterval != 0 {
			return true
		}

		// 每20 条 update 算一次 checksum
		sum := o.checkSum()
		if data.Data.Checksum != sum {
			logx.Errorf("Inconsistent depth,symbol:%s ,type:%s ,local checksum:%v, current checksum:%v",
				o.symbol.ExFormat, data.Type, sum, data.Data.Checksum)
			o.status = depthUnitStatusResetting
			o.resub()
			return false
		}
		//fmt.Printf("***right***,symbol:%s checksum:%v, local_checksum:%v\n", o.symbol.ExFormat, data.Data.Checksum, sum)
		return true

	case depthUnitStatusResetting:
		if data.Type == "partial" {
			// 只接收 首次推送的数据
			o.reset() // 置为 normal 并清空数据
			updateFn()
			return true
		}
		return false
	default:
		panic(fmt.Sprintf("wrong depthUnitStatus:%v", o.status))
	}
}

// reset 置为 normal 并清空数据
func (o *depthUnit) reset() {
	o.asks.Clear()
	o.bids.Clear()
	o.status = depthUnitStatusNormal
}

// resub 重新订阅
func (o *depthUnit) resub() {
	o.subscriber.Unsub(o.symbol.ExFormat) // 取消订阅
	o.subscriber.Sub(o.symbol.ExFormat)   // 再次订阅
}

func (o *depthManager) run() {
	ch := o.subscriber.ReadCh()
	ticker := time.NewTicker(depthUpdateCheckDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fn := func() {
				defer o.mutex.Unlock()
				o.mutex.Lock()

				var resetConnTopics []string

				for _, unit := range o.store {
					if time.Now().Sub(unit.lastUpdateTime) > depthUpdateCheckDuration {
						resetConnTopics = append(resetConnTopics, unit.symbol.ExFormat)
					}
				}

				o.subscriber.ResetConn(resetConnTopics...)
			}

			fn()
		case msg := <-ch:
			depth := msg.(*ftxapi.StreamDepth)
			symbol, err := o.SymbolManager.Convert(depth.Market, o.api.Api.ApiType)
			if err != nil {
				logx.Errorf("ftx stream depth can't parse symbol, depth: %+v", *depth)
				return
			}

			fn := func() {
				defer o.mutex.Unlock()
				o.mutex.Lock()

				info, ok := o.store[symbol.StdSymbol]
				if !ok { // 还未初始化完成，先丢弃
					info = newDepthUnit(symbol, o.api, o.outputCh, o.subscriber)
					o.store[symbol.StdSymbol] = info
				}

				info.inputCh <- depth
			}

			fn()
		}
	}

}

func (o *depthManager) Sub(symbols ...exmodel.StdSymbol) {
	defer o.mutex.Unlock()
	o.mutex.Lock()

	var topics []string
	for _, s := range symbols {
		if _, ok := o.store[s]; ok {
			continue
		}

		r, err := o.SymbolManager.GetSymbol(s)
		if err != nil {
			logx.Infof("%s not support symbol:%s", Name, s)
		} else {
			topics = append(topics, r.ExFormat)
		}
	}

	o.subscriber.Sub(topics...)
}

func (o *depthManager) OutputCh() <-chan *exmodel.StreamDepth {
	return o.outputCh
}
