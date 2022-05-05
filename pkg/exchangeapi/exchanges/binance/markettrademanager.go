package binance

import (
	"exterior-interactor/pkg/exchangeapi/exchanges/binance/api/binancespot"
	"exterior-interactor/pkg/exchangeapi/exmodel"
	"exterior-interactor/pkg/exchangeapi/extools"
	"exterior-interactor/pkg/httptools"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
	"time"
)

type marketTradeManager struct {
	extools.SymbolManager
	api      *NativeApi
	spot     httptools.AutoWsSubscriber
	outputCh chan *exmodel.StreamMarketTrade
}

func NewMarketTradeManager(mgr extools.SymbolManager, api *NativeApi) extools.MarketTradeManager {
	res, err := api.SpotApi.GetStreamMarketTrade()
	if err != nil {
		panic(err)
	}

	tradeMgr := &marketTradeManager{
		SymbolManager: mgr,
		api:           api,
		spot:          res,
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
	spotCh := o.spot.ReadCh()
	for {
		select {
		case msg := <-spotCh:
			t := msg.(*binancespot.StreamMarketTrade)
			symbol, err := o.SymbolManager.Convert(t.Data.S, o.api.SpotApi.ApiType)
			if err != nil {
				logx.Errorf("binance spot ws market trade can't parse symbol, data:%+v", *t)
				continue
			}

			o.outputCh <- &exmodel.StreamMarketTrade{
				TradeId:  fmt.Sprint(t.Data.T),
				Exchange: Name,
				Time:     time.UnixMilli(t.Data.T1),
				Symbol:   symbol,
				Price:    t.Data.P,
				Volume:   t.Data.Q,
			}
		}
	}
}

func (o *marketTradeManager) Sub(symbols ...exmodel.StdSymbol) {
	var valid []*exmodel.Symbol

	for _, s := range symbols {
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
}

func (o *marketTradeManager) parseSymbolToTopic(symbols ...*exmodel.Symbol) (spotTopics, coinFutureTopics, usdtFutureTopics []string) {
	for _, symbol := range symbols {
		topic := fmt.Sprintf("%s@trade", strings.ToLower(symbol.ExFormat))
		switch symbol.Type {
		case exmodel.SymbolTypeSpot:
			spotTopics = append(spotTopics, topic)
		case exmodel.SymbolTypeUsdtPerp, exmodel.SymbolTypeUsdtDelivery:
			usdtFutureTopics = append(usdtFutureTopics, topic)
		case exmodel.SymbolTypeCoinPerp, exmodel.SymbolTypeCoinDelivery:
			coinFutureTopics = append(coinFutureTopics, topic)
		default:
			logx.Errorf("not support sub trade with symbol:%s", symbol.StdSymbol)
		}
	}
	return
}
