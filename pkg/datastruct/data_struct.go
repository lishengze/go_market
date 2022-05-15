package datastruct

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
)

const (
	BCTS_EXCHANGE = "_bcts_"
)

const (
	NANO_PER_SECS = 1000000000
)

type TestData struct {
	Name string
}

type TSymbol string
type TExchange string

// type RFloat float64
type TPrice float64
type TVolume float64

type Metadata struct {
	DepthMeta map[string](map[string]struct{})
	KlineMeta map[string](map[string]struct{})
	TradeMeta map[string](map[string]struct{})
}

type InnerDepth struct {
	Volume         float64
	ExchangeVolume map[string]float64
}

func (src *InnerDepth) Add(other *InnerDepth) {
	if src == other {
		return
	}

	src.Volume += other.Volume

	for exchange, volume := range other.ExchangeVolume {
		src.ExchangeVolume[exchange] += volume
	}
}

type DepthQuote struct {
	Exchange string
	Symbol   string
	Time     int64
	Asks     *treemap.Map
	Bids     *treemap.Map
}

type Kline struct {
	Exchange   string
	Symbol     string
	Time       int64
	Open       float64
	High       float64
	Low        float64
	Close      float64
	Volume     float64
	Resolution int
}

type Trade struct {
	Exchange string
	Symbol   string
	Time     int64
	Price    float64
	Volume   float64
}

func (t *Trade) String() string {
	res := fmt.Sprintf("%s.%s, %+v, p: %f v: %f \n", t.Exchange, t.Symbol,
		time.Unix(int64(t.Time/NANO_PER_SECS), t.Time%NANO_PER_SECS), t.Price, t.Volume)
	return res
}

func (k *Kline) String() string {
	res := fmt.Sprintf("%s.%s, %+v, o: %f, h: %f, l: %f, c: %f, v: %f\n",
		k.Exchange, k.Symbol, time.Unix(int64(k.Time/NANO_PER_SECS), k.Time%NANO_PER_SECS),
		k.Open, k.High, k.Low, k.Close, k.Volume)
	return res
}

func (d *DepthQuote) String(len int) string {

	res := fmt.Sprintf("%s.%s, %v\nAsks: %s\nBids: %s \n", d.Exchange, d.Symbol,
		time.Unix(int64(d.Time/NANO_PER_SECS), d.Time%NANO_PER_SECS),
		d.Asks.String(), d.Bids.String())

	return string(res)
}

func NewTrade(src *Trade) *Trade {
	if src != nil {
		rst := &Trade{
			Exchange: src.Exchange,
			Symbol:   src.Symbol,
			Time:     src.Time,
			Price:    src.Price,
			Volume:   src.Volume,
		}
		return rst
	} else {
		rst := &Trade{}
		return rst
	}
}

func NewKline(src *Kline) *Kline {
	if src != nil {
		rst := &Kline{
			Exchange:   src.Exchange,
			Symbol:     src.Symbol,
			Time:       src.Time,
			Open:       src.Open,
			High:       src.High,
			Low:        src.Low,
			Close:      src.Close,
			Volume:     src.Volume,
			Resolution: src.Resolution,
		}
		return rst
	} else {
		rst := &Kline{}
		return rst
	}
}

func InitKlineByTrade(src *Kline, trade *Trade) {
	src.Exchange = BCTS_EXCHANGE
	src.Symbol = trade.Symbol
	src.Time = trade.Time
	src.Resolution = 60
	src.Open = trade.Price
	src.High = trade.Price
	src.Low = trade.Price
	src.Close = trade.Price
	src.Volume = trade.Volume
}

func NewDepth(src *DepthQuote) *DepthQuote {
	if src != nil {
		rst := &DepthQuote{
			Exchange: src.Exchange,
			Symbol:   src.Symbol,
			Time:     src.Time,
			Asks:     treemap.NewWith(utils.Float64Comparator),
			Bids:     treemap.NewWith(utils.Float64Comparator),
		}

		ask_iter := src.Asks.Iterator()
		for ask_iter.Begin(); ask_iter.Next(); {
			rst.Asks.Put(ask_iter.Key(), ask_iter.Value())
		}

		bid_iter := src.Bids.Iterator()
		for bid_iter.Begin(); bid_iter.Next(); {
			rst.Bids.Put(bid_iter.Key(), bid_iter.Value())
		}

		return rst
	} else {
		rst := &DepthQuote{
			Asks: treemap.NewWith(utils.Float64Comparator),
			Bids: treemap.NewWith(utils.Float64Comparator),
		}
		return rst
	}
}

type DataChannel struct {
	TradeChannel chan *Trade
	KlineChannel chan *Kline
	DepthChannel chan *DepthQuote
}

func (d *DepthQuote) Init() {
	d.Asks = treemap.NewWith(utils.Float64Comparator)
	d.Bids = treemap.NewWith(utils.Float64Comparator)
}

func GetTestDepth() *DepthQuote {
	var rst DepthQuote
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(3)
	exchange_array := []string{"FTX", "HUOBI", "OKEX"}
	exchange_type := index % 3

	rst.Exchange = exchange_array[exchange_type]
	rst.Symbol = "BTC_USDT"
	rst.Time = time.Now().Unix()
	rst.Asks = treemap.NewWith(utils.Float64Comparator)
	rst.Bids = treemap.NewWith(utils.Float64Comparator)

	rst.Asks.Put(55000.0, &InnerDepth{5.5, map[string]float64{rst.Exchange: 5.5}})
	rst.Asks.Put(50000.0, &InnerDepth{5.0, map[string]float64{rst.Exchange: 5.0}})

	rst.Bids.Put(45000.0, &InnerDepth{4.5, map[string]float64{rst.Exchange: 4.5}})
	rst.Bids.Put(40000.0, &InnerDepth{4.0, map[string]float64{rst.Exchange: 4.0}})

	switch exchange_type {
	case 0:
		rst.Asks.Put(60000.0, &InnerDepth{6.0, map[string]float64{rst.Exchange: 6.0}})
		rst.Bids.Put(35000.0, &InnerDepth{3.5, map[string]float64{rst.Exchange: 3.5}})

	case 1:
		rst.Asks.Put(70000.0, &InnerDepth{7.0, map[string]float64{rst.Exchange: 7.0}})
		rst.Bids.Put(30000.0, &InnerDepth{3.0, map[string]float64{rst.Exchange: 3.0}})

	case 2:
		rst.Asks.Put(75000.0, &InnerDepth{7.5, map[string]float64{rst.Exchange: 7.5}})
		rst.Bids.Put(25000.0, &InnerDepth{2.5, map[string]float64{rst.Exchange: 2.5}})
	}

	return &rst
}

func GetTestTrade() *Trade {
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(3)

	exchange_array := []string{"FTX", "HUOBI", "OKEX"}
	cur_exchange := exchange_array[randomNum%3]
	symbol := "ETH_USDT"
	trade_price := float64(rand.Intn(1000))
	trade_volume := float64(rand.Intn(100))

	new_trade := NewTrade(nil)
	new_trade.Exchange = cur_exchange
	new_trade.Symbol = symbol
	new_trade.Price = trade_price
	new_trade.Volume = trade_volume
	new_trade.Time = time.Now().Unix()

	// fmt.Printf("Send Trade: %s\n", new_trade.String())

	return new_trade
}

func GetTestKline() *Kline {
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(3)

	exchange_array := []string{"FTX", "HUOBI", "OKEX"}
	cur_exchange := exchange_array[randomNum%3]
	symbol := "ETH_USDT"

	new_kline := Kline{
		Exchange:   cur_exchange,
		Symbol:     symbol,
		Time:       time.Now().Unix(),
		Resolution: 60,
		Volume:     1.1,
		Open:       3000,
		High:       4000,
		Low:        2800,
		Close:      3500,
	}

	return &new_kline
}
