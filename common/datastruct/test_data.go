package datastruct

import (
	"fmt"
	"market_server/common/util"
	"math/rand"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
)

func GetTestMetadata(symbols []string) *Metadata {
	symbol_set := make(map[string](map[string]struct{}))
	exchange_set := make(map[string]struct{})
	exchange_set["FTX"] = struct{}{}

	for _, symbol := range symbols {
		symbol_set[symbol] = exchange_set
	}

	// symbol_set["ETH_USDT"] = exchange_set

	MetaData := Metadata{}

	MetaData.DepthMeta = symbol_set
	MetaData.TradeMeta = symbol_set

	return &MetaData
}

func GetTestDepth() *DepthQuote {
	var rst DepthQuote
	rand.Seed(time.Now().UnixNano())
	exchange_type := rand.Intn(3)
	exchange_array := []string{"FTX", "HUOBI", "OKEX"}

	symbol_index := rand.Intn(3)
	symbol_array := []string{"BTC_USDT", "ETH_USDT", "DOT_USDT"}

	rst.Exchange = exchange_array[exchange_type]
	rst.Symbol = symbol_array[symbol_index]
	rst.Time = util.UTCNanoTime()
	rst.Asks = treemap.NewWith(utils.Float64Comparator)
	rst.Bids = treemap.NewWith(util.Float64ComparatorDsc)

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
	symbol := "BTC_USDT"
	trade_price := float64(rand.Intn(1000))
	trade_volume := float64(rand.Intn(100))

	new_trade := NewTrade(nil)
	new_trade.Exchange = cur_exchange
	new_trade.Symbol = symbol
	new_trade.Price = trade_price
	new_trade.Volume = trade_volume
	new_trade.Time = util.UTCNanoTime()

	return new_trade
}

func GetTestKline() *Kline {
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(3)

	exchange_array := []string{"FTX", "HUOBI", "OKEX"}
	cur_exchange := exchange_array[randomNum%3]
	symbol := "BTC_USDT"

	new_kline := Kline{
		Exchange:   cur_exchange,
		Symbol:     symbol,
		Time:       util.UTCNanoTime(),
		Resolution: 60,
		Volume:     1.1,
		Open:       3000,
		High:       4000,
		Low:        2800,
		Close:      3500,
	}

	return &new_kline
}

func GetTestDepthMultiSymbols(symbol_list []string, exchange string) *DepthQuote {
	var rst DepthQuote
	rand.Seed(time.Now().UnixNano())

	symbol_index := rand.Intn(len(symbol_list))

	rst.Exchange = exchange
	rst.Symbol = symbol_list[symbol_index]
	rst.Time = util.UTCNanoTime()
	rst.Asks = treemap.NewWith(utils.Float64Comparator)
	rst.Bids = treemap.NewWith(util.Float64ComparatorDsc)

	rst.Asks.Put(55000.0, &InnerDepth{5.5, map[string]float64{rst.Exchange: 5.5}})
	rst.Asks.Put(50000.0, &InnerDepth{5.0, map[string]float64{rst.Exchange: 5.0}})

	rst.Bids.Put(45000.0, &InnerDepth{4.5, map[string]float64{rst.Exchange: 4.5}})
	rst.Bids.Put(40000.0, &InnerDepth{4.0, map[string]float64{rst.Exchange: 4.0}})

	switch symbol_index {
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

func GetTestTradeMultiSymbols(symbol_list []string, exchange string) *Trade {
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(len(symbol_list))

	cur_exchange := exchange
	symbol := symbol_list[randomNum]
	trade_price := float64(rand.Intn(1000))
	trade_volume := float64(rand.Intn(100))

	new_trade := NewTrade(nil)
	new_trade.Exchange = cur_exchange
	new_trade.Symbol = symbol
	new_trade.Price = trade_price
	new_trade.Volume = trade_volume
	new_trade.Time = util.UTCNanoTime()

	return new_trade
}

func GetTestHistKline(req_kline_info *ReqHistKline) *treemap.Map {

	klines := treemap.NewWith(utils.Int64Comparator)

	real_count := int64(req_kline_info.Count * (req_kline_info.Frequency / SECS_PER_MIN))

	cur_time := util.TimeMinuteNanos() - real_count*NANO_PER_MIN

	if req_kline_info.Count != 0 {
		for i := 0; i < int(real_count); i++ {
			cur_time += NANO_PER_MIN

			klines.Put(cur_time, &Kline{
				Exchange:   req_kline_info.Exchange,
				Symbol:     req_kline_info.Symbol,
				Time:       cur_time,
				Open:       float64(1000 + rand.Intn(100)),
				High:       float64(1200 + rand.Intn(100)),
				Low:        float64(800 + rand.Intn(100)),
				Close:      float64(1000 + rand.Intn(100)),
				Volume:     float64(rand.Intn(100)),
				Resolution: SECS_PER_MIN,
			})
		}

		return klines
	}

	return nil
}

func GetTestKlineMultiSymbols(symbol_list []string, exchange string, last_time int64) *Kline {
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(len(symbol_list))

	cur_exchange := exchange
	symbol := symbol_list[randomNum]

	new_kline := Kline{
		Exchange:   cur_exchange,
		Symbol:     symbol,
		Time:       last_time + SECS_PER_MIN*NANO_PER_SECS,
		Resolution: SECS_PER_MIN,
		Open:       float64(1000 + rand.Intn(100)),
		High:       float64(1200 + rand.Intn(100)),
		Low:        float64(800 + rand.Intn(100)),
		Close:      float64(1000 + rand.Intn(100)),
		Volume:     float64(rand.Intn(100)),
	}

	return &new_kline
}

func TestKlineTime() {
	cur_time := util.TimeMinuteNanos()

	symbol_list := []string{"BTC_USDT"}

	kline := GetTestKlineMultiSymbols(symbol_list, BCTS_EXCHANGE, cur_time)

	fmt.Printf("Kline: %s \n", kline.String())

	fmt.Printf("IsOldKlineEnd: %+v , IsNewKline: %+v \n", IsOldKlineEnd(kline, 60), IsNewKlineStart(kline, 60))
}

func TestGetHistKlineData() {
	req_info := &ReqHistKline{
		Symbol:    "BTC_USDT",
		Exchange:  BCTS_EXCHANGE,
		Count:     MIN_PER_DAY,
		Frequency: SECS_PER_MIN,
	}

	klines := GetTestHistKline(req_info)

	iter := klines.Iterator()

	if iter.First() {
		fmt.Printf("First Kline: %s ", iter.Value().(*Kline).String())
	}

	if iter.Last() {
		fmt.Printf("Last Kline: %s ", iter.Value().(*Kline).String())
	}
}

func TestIsTargetTime() {
	// time_obj := time.Date(2022, 7, 18, 0, 0, 0, 0, time.Local)
	time_obj := time.Date(1970, 1, 5, 0, 0, 0, 0, time.UTC)

	nano_secs := time_obj.UTC().UnixNano()
	fmt.Printf("nano_secs: %d, CurTime: %s, IsTarget: %+v \n",
		nano_secs,
		util.TimeStrFromInt(nano_secs),
		IsTargetTime(nano_secs/NANO_PER_SECS, SECS_PER_DAY*7))

}

func TestGetLastStartTime() {
	time_obj := time.Date(1970, 1, 6, 10, 20, 0, 0, time.UTC)
	original_time := time_obj.UTC().UnixNano()

	// original_time := time.Now().UnixNano()

	last_start_time := GetLastStartTime(original_time, SECS_PER_DAY*7)

	fmt.Printf("original_time: %s, last_start_time: %s\n", util.TimeStrFromInt(original_time), util.TimeStrFromInt(last_start_time))

}
