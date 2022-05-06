package extools

import (
	"exterior-interactor/pkg/exchangeapi/exmodel"
	"exterior-interactor/pkg/timeutils"
	"exterior-interactor/pkg/xmath"
	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type (
	KlineGenerator interface {
		InputMarketTrade(t *exmodel.StreamMarketTrade)
		OutputCh() <-chan *exmodel.Kline
	}

	klineGenerator struct {
		store    map[exmodel.StdSymbol]*klineInfo
		inputCh  chan *exmodel.StreamMarketTrade
		outputCh chan *exmodel.Kline
	}

	klineInfo struct {
		minuteSeq int       // 标记第几分钟
		exactTime time.Time // 整分钟时间
		high      decimal.Decimal
		low       decimal.Decimal
		volume    decimal.Decimal
		value     decimal.Decimal
		trades    []*exmodel.StreamMarketTrade
	}
)

func NewKlineGenerator() KlineGenerator {
	g := &klineGenerator{
		store:    make(map[exmodel.StdSymbol]*klineInfo),
		inputCh:  make(chan *exmodel.StreamMarketTrade, 100000),
		outputCh: make(chan *exmodel.Kline, 100000),
	}
	go g.run()
	return g
}

func NewKlineInfo(time time.Time) *klineInfo {
	return &klineInfo{
		minuteSeq: 1,
		exactTime: timeutils.TimeToExactMinute(time),
		trades:    make([]*exmodel.StreamMarketTrade, 0),
	}
}

func (o *klineInfo) generateKline() *exmodel.Kline {
	if len(o.trades) == 0 {
		return &exmodel.Kline{}
	}
	open := xmath.MustDecimal(o.trades[0].Price).String()
	close_ := xmath.MustDecimal(o.trades[len(o.trades)-1].Price).String()
	//fmt.Println("_____start:", o.trades[0].Time)
	//fmt.Println("_______end:", o.trades[len(o.trades)-1].Time)
	//fmt.Println("_____total:", len(o.trades))
	return &exmodel.Kline{
		Exchange:   o.trades[0].Exchange,
		Time:       timeutils.TimeToExactMinute(o.trades[0].Time),
		LocalTime:  time.Now(),
		Symbol:     o.trades[0].Symbol,
		Resolution: 60,
		Open:       open,
		High:       o.high.String(),
		Low:        o.low.String(),
		Close:      close_,
		Volume:     o.volume.String(),
		Value:      o.value.String(),
	}
}

func (o *klineInfo) reset(currentTradeTime time.Time) {
	o.exactTime = timeutils.TimeToExactMinute(currentTradeTime)
	o.minuteSeq += 1
	o.low = decimal.Zero
	o.high = decimal.Zero
	o.value = decimal.Zero
	o.volume = decimal.Zero
	o.trades = o.trades[0:0]
}

func (o *klineInfo) addTrade(trade *exmodel.StreamMarketTrade) {
	price, err := decimal.NewFromString(trade.Price)
	if err != nil {
		logx.Errorf("parse price err:%v, StreamMarketTrade:%+v", err, *trade)
		return
	}

	volume, err := decimal.NewFromString(trade.Volume)
	if err != nil {
		logx.Errorf("parse volume err:%v, StreamMarketTrade:%+v", err, *trade)
		return
	}

	value := price.Mul(volume)

	if price.LessThan(o.low) || o.low.Equal(decimal.Zero) {
		o.low = price
	}

	if price.GreaterThan(o.high) || o.high.Equal(decimal.Zero) {
		o.high = price
	}

	o.volume = o.volume.Add(volume)
	o.value = o.value.Add(value)

	if len(o.trades) > 0 {
		lastTrade := o.trades[len(o.trades)-1]
		if trade.Time.Before(lastTrade.Time) {
			logx.Errorf("recv an unordered trade, current:%+v, last:%+v", *trade, *lastTrade)
		}
	}

	o.trades = append(o.trades, trade)
}

func (o *klineGenerator) run() {
	for {
		select {
		case t := <-o.inputCh:
			info, ok := o.store[t.Symbol.StdSymbol]
			if !ok {
				info = NewKlineInfo(t.Time)
				o.store[t.Symbol.StdSymbol] = info
			}
			if t.Time.Sub(info.exactTime) >= time.Minute {
				// 已经到了下一分钟，推送 kline 数据，并清空上一分钟数据
				kline := info.generateKline()
				if info.minuteSeq != 1 { // 启动的第一分钟的数据基本是不全的，所以不推送
					o.outputCh <- kline
				}
				info.reset(t.Time)

			}

			info.addTrade(t)
		}
	}
}

// InputMarketTrade 输入 逐笔成交
func (o *klineGenerator) InputMarketTrade(t *exmodel.StreamMarketTrade) {
	o.inputCh <- t
}

func (o *klineGenerator) OutputCh() <-chan *exmodel.Kline {
	return o.outputCh
}
