package ftx

import (
	"exterior-interactor/pkg/exchangeapi/exchanges/ftx/ftxapi"
	"exterior-interactor/pkg/exchangeapi/exmodel"
	"exterior-interactor/pkg/exchangeapi/extools"
	"exterior-interactor/pkg/httptools"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type marketTradeManager struct {
	extools.SymbolManager
	api        *NativeApi
	subscriber httptools.AutoWsSubscriber
	outputCh   chan *exmodel.StreamMarketTrade
}

func NewMarketTradeManager(mgr extools.SymbolManager, api *NativeApi) extools.MarketTradeManager {
	res, err := api.Api.GetStreamMarketTrade()
	if err != nil {
		panic(err)
	}

	tradeMgr := &marketTradeManager{
		SymbolManager: mgr,
		api:           api,
		subscriber:    res,
		outputCh:      make(chan *exmodel.StreamMarketTrade, 1024),
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
	for {
		select {
		case msg := <-ch:
			trade := msg.(*ftxapi.StreamMarketTrade)
			for _, t := range trade.Data {
				symbol, err := o.SymbolManager.Convert(trade.Market, o.api.Api.ApiType)
				if err != nil {
					logx.Errorf("ftx  trade can't parse symbol, data:%+v", t)
					continue
				}

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

		}
	}
}

func (o *marketTradeManager) Sub(symbols ...exmodel.StdSymbol) {
	var topics []string
	for _, s := range symbols {
		r, err := o.SymbolManager.GetSymbol(s)
		if err != nil {
			logx.Infof("%s not support symbol:%s", Name, s)
		} else {
			topics = append(topics, r.ExFormat)
		}
	}

	o.subscriber.Sub(topics...)
}
