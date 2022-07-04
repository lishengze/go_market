package ftx

import (
	"exterior-interactor/pkg/exchangeapi/exchanges/ftx/ftxapi"
	"exterior-interactor/pkg/exchangeapi/exmodel"
	"exterior-interactor/pkg/exchangeapi/extools"
	"exterior-interactor/pkg/httptools"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
	"time"
)

const (
	marketTradeUpdateCheckDuration = time.Second * 5 //  超过此时间未更新，将触发重连
)

type marketTradeManager struct {
	mutex sync.Mutex
	extools.SymbolManager
	api        *NativeApi
	subscriber httptools.AutoWsSubscriber
	outputCh   chan *exmodel.StreamMarketTrade
	store      map[string]time.Time // 记录每个 topic 的更新时间
}

func NewMarketTradeManager(mgr extools.SymbolManager, api *NativeApi) extools.MarketTradeManager {
	res, err := api.Api.GetStreamMarketTrade()
	if err != nil {
		panic(err)
	}

	tradeMgr := &marketTradeManager{
		mutex:         sync.Mutex{},
		SymbolManager: mgr,
		api:           api,
		subscriber:    res,
		outputCh:      make(chan *exmodel.StreamMarketTrade, 1024),
		store:         map[string]time.Time{},
	}

	go tradeMgr.run()

	return tradeMgr
}

func (o *marketTradeManager) OutputCh() <-chan *exmodel.StreamMarketTrade {
	return o.outputCh
}

// run 开启 goroutine 转发数据
func (o *marketTradeManager) run() {
	ch := o.subscriber.ReadCh()
	ticker := time.NewTicker(marketTradeUpdateCheckDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fn := func() {
				defer o.mutex.Unlock()
				o.mutex.Lock()

				var resetConnTopics []string

				for topic, t := range o.store {
					if time.Now().Sub(t) > marketTradeUpdateCheckDuration {
						resetConnTopics = append(resetConnTopics, topic)
					}
				}

				o.subscriber.ResetConn(resetConnTopics...)
			}

			fn()
		case msg := <-ch:
			trade := msg.(*ftxapi.StreamMarketTrade)
			symbol, err := o.SymbolManager.Convert(trade.Market, o.api.Api.ApiType)
			if err != nil {
				logx.Errorf("ftx  trade can't parse symbol, data:%+v", *trade)
				continue
			}

			for _, t := range trade.Data {
				//time_, err := time.Parse(time.RFC3339Nano, t.Time)
				//if err != nil {
				//	logx.Errorf("ftx  trade can't parse time, data:%+v", t)
				//	continue
				//}
				o.outputCh <- &exmodel.StreamMarketTrade{
					TradeId:   fmt.Sprint(t.Id),
					Exchange:  Name,
					Time:      t.Time,
					LocalTime: time.Now(),
					Symbol:    symbol,
					Price:     fmt.Sprint(t.Price),
					Volume:    fmt.Sprint(t.Size),
				}
			}

			fn := func() {
				{
					o.store[symbol.ExFormat] = time.Now()
					defer o.mutex.Unlock()
					o.mutex.Lock()
				}
			}

			fn()

		}
	}
}

func (o *marketTradeManager) Sub(symbols ...exmodel.StdSymbol) {
	defer o.mutex.Unlock()
	o.mutex.Lock()

	var topics []string
	for _, s := range symbols {
		r, err := o.SymbolManager.GetSymbol(s)
		if err != nil {
			logx.Infof("%s not support symbol:%s", Name, s)
		} else {
			topics = append(topics, r.ExFormat)
			o.store[r.ExFormat] = time.Now()
		}
	}

	o.subscriber.Sub(topics...)
}
