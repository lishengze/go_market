package riskctrl

import (
	"fmt"
	"math"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
	"github.com/shopspring/decimal"
)

// type TString string

func LOG_INFO(info string) {
	fmt.Println("INFO: " + info)
}

func LOG_WARN(info string) {
	fmt.Println("WARN: " + info)
}

func LOG_ERROR(info string) {
	fmt.Println("Error " + info)
}

func ExceptionFunc() {
	errMsg := recover()
	if errMsg != nil {
		fmt.Println(errMsg)
	}
}

// func catch_exp() {
// 	errMsg := recover()

// 	if errMsg != nil {
// 		fmt.Println(errMsg)

// 	}

// 	fmt.Println("This is catch_exp func")

// }

type HedgeConfig struct {
	FeeKind  int
	FeeValue float64
}

type RiskCtrlConfig struct {
	HedgeConfigMap map[string]HedgeConfig

	PricePrecison  uint32
	VolumePrecison uint32

	PriceBiasValue float64
	PriceBiasKind  int

	VolumeBiasValue float64
	VolumeBiasKind  int

	PriceMinumChange float64
}

type RiskCtrlConfigMap map[string]RiskCtrlConfig

type RiskWorkerInterface interface {
	Process(depth_quote *DepthQuote, configs *RiskCtrlConfigMap) bool
	Execute(depth_quote *DepthQuote, configs *RiskCtrlConfigMap) bool
	SetNext(nextWorker RiskWorkerInterface)
	GetWorkerName() string
}

type Worker struct {
	nextWorker RiskWorkerInterface

	WorkerName string
}

func (w *Worker) Process(depth_quote *DepthQuote, configs *RiskCtrlConfigMap) bool {
	fmt.Println("Original Worker Process")
	return false
}

func (w *Worker) Execute(depth_quote *DepthQuote, configs *RiskCtrlConfigMap) bool {
	defer ExceptionFunc()

	fmt.Println(w)

	w.Process(depth_quote, configs)

	if w.nextWorker != nil {
		w.nextWorker.Execute(depth_quote, configs)
	} else {
		fmt.Printf("%s Worker Has No Next Worker!\n", w.WorkerName)
	}

	return false
}

func (w *Worker) SetNext(next RiskWorkerInterface) {
	defer ExceptionFunc()

	w.nextWorker = next
}

func (w *Worker) GetWorkerName() string {
	return w.WorkerName
}

type FeeWorker struct {
	Worker
}

func get_bias_value(original_value float64, bias_kind int, bias_value float64) float64 {
	defer ExceptionFunc()

	rst := decimal.NewFromFloat(original_value)

	// biased_value := original_value

	if bias_kind == 1 {
		rst = rst.Mul(decimal.NewFromFloat(1 + bias_value))
	} else if bias_kind == 2 {
		rst = rst.Add(decimal.NewFromFloat(bias_value))
	} else {
		LOG_ERROR("Error Bias Kind")
	}

	rst_float, _ := rst.Float64()
	if rst_float < 0 {
		rst_float = 0
	}

	return rst_float
}

func get_new_fee_price(original_price float64, exchange string, hedge_config *RiskCtrlConfig, isAsk bool) float64 {
	defer ExceptionFunc()

	new_price := original_price
	if hedge_config, ok := hedge_config.HedgeConfigMap[exchange]; ok {

		if isAsk {
			new_price = get_bias_value(original_price, hedge_config.FeeKind, hedge_config.FeeValue)
		} else {
			new_price = get_bias_value(original_price, hedge_config.FeeKind, hedge_config.FeeValue*-1)
		}
	}
	return new_price
}

// 根据每个交易所的手续费率，计算手续费
func (w *FeeWorker) calc_depth_fee(depth *treemap.Map, config *RiskCtrlConfig, isAsk bool) *treemap.Map {
	defer ExceptionFunc()

	result := treemap.NewWith(utils.Float64Comparator)

	iter := depth.Iterator()

	for iter.Begin(); iter.Next(); {
		original_price := iter.Key().(float64)
		inner_depth := iter.Value().(InnerDepth)

		for exchange, _ := range inner_depth.ExchangeVolume {
			new_price := get_new_fee_price(original_price, exchange, config, isAsk)

			if new_inner_depth_iter, ok := result.Get(new_price); ok {

				new_inner_depth := new_inner_depth_iter.(InnerDepth)

				new_inner_depth.ExchangeVolume[exchange] += inner_depth.Volume
				new_inner_depth.Volume += inner_depth.Volume

			} else {
				new_inner_depth := InnerDepth{0, make(map[string]float64)}

				new_inner_depth.ExchangeVolume[exchange] = inner_depth.Volume
				new_inner_depth.Volume = inner_depth.Volume

				result.Put(new_price, new_inner_depth)
			}
		}
	}

	return result
}

func (w *FeeWorker) Execute(depth_quote *DepthQuote, configs *RiskCtrlConfigMap) bool {
	defer ExceptionFunc()

	fmt.Println(w)

	w.Process(depth_quote, configs)

	if w.nextWorker != nil {
		w.nextWorker.Execute(depth_quote, configs)
	} else {
		fmt.Printf("%s Worker Has No Next Worker!\n", w.WorkerName)
	}

	return false
}

func (w *FeeWorker) Process(depth_quote *DepthQuote, configs *RiskCtrlConfigMap) bool {
	defer ExceptionFunc()

	fmt.Println("-------- FeeWorker Process  ---------")

	if config, ok := (*configs)[depth_quote.Symbol]; ok {

		fmt.Printf("\nConfig:%v \n", configs)
		LOG_INFO("\nBefore FeeCtrl: \n" + depth_quote.String(5))
		// fmt.Println(depth_quote)

		new_asks := w.calc_depth_fee(depth_quote.Asks, &config, true)
		depth_quote.Asks = new_asks

		new_bids := w.calc_depth_fee(depth_quote.Bids, &config, false)
		depth_quote.Bids = new_bids

		LOG_INFO("\nAfter FeeCtrl: \n" + depth_quote.String(5))
		// fmt.Println(depth_quote)

	} else {

		LOG_ERROR("Symbol: " + string(depth_quote.Symbol) + " Has no config")
		return false
	}

	fmt.Println("FeeWorker Process")
	return false
}

type QuotebiasWorker struct {
	Worker
}

func calc_depth_bias(depth *treemap.Map, config *RiskCtrlConfig, isAsk bool) *treemap.Map {
	defer ExceptionFunc()

	result := treemap.NewWith(utils.Float64Comparator)

	depth_iter := depth.Iterator()

	PriceBiasValue := config.PriceBiasValue
	VolumeBiasValue := config.VolumeBiasValue

	if isAsk == false {
		PriceBiasValue *= -1
		VolumeBiasValue *= -1
	}

	for depth_iter.Begin(); depth_iter.Next(); {
		original_price := depth_iter.Key().(float64)
		inner_depth := depth_iter.Value().(InnerDepth)

		// var new_inner_depth InnerDepth

		new_inner_depth := InnerDepth{0, make(map[string]float64)}

		new_price := get_bias_value(original_price, config.PriceBiasKind, PriceBiasValue)
		new_volume := get_bias_value(inner_depth.Volume, config.VolumeBiasKind, VolumeBiasValue)

		new_inner_depth.Volume = new_volume

		for exchange, exchange_volume := range inner_depth.ExchangeVolume {
			new_exchange_volume := get_bias_value(exchange_volume, config.VolumeBiasKind, VolumeBiasValue)

			new_inner_depth.ExchangeVolume[exchange] = new_exchange_volume
		}

		result.Put(new_price, new_inner_depth)
	}

	return result
}

func (w *QuotebiasWorker) Process(depth_quote *DepthQuote, configs *RiskCtrlConfigMap) bool {
	defer ExceptionFunc()

	if config, ok := (*configs)[depth_quote.Symbol]; ok {

		fmt.Printf("config:%v \n", configs)
		LOG_INFO("\nBefore QuotebiasCtrl: \n" + depth_quote.String(5))
		// fmt.Println(depth_quote)

		new_asks := calc_depth_bias(depth_quote.Asks, &config, true)
		depth_quote.Asks = new_asks

		new_bids := calc_depth_bias(depth_quote.Bids, &config, false)
		depth_quote.Bids = new_bids

		LOG_INFO("\nAfter QuotebiasCtrl: \n" + depth_quote.String(5))
		// fmt.Println(depth_quote)

	} else {

		LOG_ERROR("Symbol: " + string(depth_quote.Symbol) + " Has no config")
		return false
	}

	fmt.Println("QuotebiasWorker Process")
	return false
}

func (w *QuotebiasWorker) Execute(depth_quote *DepthQuote, configs *RiskCtrlConfigMap) bool {
	defer ExceptionFunc()

	fmt.Println(w)

	w.Process(depth_quote, configs)

	if w.nextWorker != nil {
		w.nextWorker.Execute(depth_quote, configs)
	} else {
		fmt.Printf("%s Worker Has No Next Worker!\n", w.WorkerName)
	}

	return false
}

type WatermarkWorker struct {
	Worker
}

// 每个交易所的买一卖一档，然后取买一卖一中位数的均值
func calc_watermark(depth_quote *DepthQuote) float64 {
	var rst float64

	ask_iter := depth_quote.Asks.Iterator()
	bid_iter := depth_quote.Bids.Iterator()
	ask_iter.First()
	bid_iter.Last()

	ask_minum := ask_iter.Key().(float64)
	bid_maxum := bid_iter.Key().(float64)

	ask_crossed_price_list := []float64{}
	for ask_iter.Begin(); ask_iter.Next(); {
		if ask_iter.Key().(float64) <= bid_maxum {
			ask_crossed_price_list = append(ask_crossed_price_list, ask_iter.Key().(float64))
		} else {
			break
		}
	}

	fmt.Printf("\n------ ask_crossed_price_list:%v \n", ask_crossed_price_list)

	bid_crossed_price_list := []float64{}
	for bid_iter.End(); bid_iter.Prev(); {
		if bid_iter.Key().(float64) >= ask_minum {
			bid_crossed_price_list = append(bid_crossed_price_list, bid_iter.Key().(float64))
		} else {
			break
		}
	}
	fmt.Printf("\n------bid_crossed_price_list:%v \n", bid_crossed_price_list)

	rst = (ask_crossed_price_list[len(ask_crossed_price_list)/2] + bid_crossed_price_list[len(bid_crossed_price_list)/2]) / 2

	fmt.Printf("\n-------watermark: %v \n", rst)

	return rst
}

func filter_depth_by_watermark(depth *treemap.Map, watermark float64, price_minum_change float64, isAsk bool) {

	crossed_price := []float64{}
	new_inner_depth := InnerDepth{0, make(map[string]float64)}
	depth_iter := depth.Iterator()
	new_price := watermark + float64(price_minum_change)

	fmt.Printf("\nNewPrice: %v\n", new_price)

	if isAsk {

		for depth_iter.Begin(); depth_iter.Next(); {
			cur_price := depth_iter.Key().(float64)
			cur_innerdepth := depth_iter.Value().(InnerDepth)

			if cur_price <= new_price {
				crossed_price = append(crossed_price, cur_price)
				new_inner_depth.Add(&cur_innerdepth)
			} else {
				break
			}
		}

	} else {

		for depth_iter.End(); depth_iter.Prev(); {
			cur_price := depth_iter.Key().(float64)
			cur_innerdepth := depth_iter.Value().(InnerDepth)

			if cur_price >= new_price {
				crossed_price = append(crossed_price, cur_price)
				new_inner_depth.Add(&cur_innerdepth)
			} else {
				break
			}
		}
	}

	if len(crossed_price) > 0 {
		for _, price := range crossed_price {

			depth.Remove(price)
		}
	}

	if new_price != 0 {
		depth.Put(new_price, new_inner_depth)
	}
}

func get_sorted_key(depth *treemap.Map) []float64 {
	var keys []float64

	// for price, _ := range depth {

	// 	keys = append(keys, price)

	// 	for i := len(keys) - 1; i > 0; i-- {
	// 		if keys[i] < keys[i-1] {
	// 			tmp := keys[i]
	// 			keys[i] = keys[i-1]
	// 			keys[i-1] = tmp
	// 		} else {
	// 			break
	// 		}
	// 	}
	// }
	return keys
}

func check_cross(depth_quote *DepthQuote) bool {
	if depth_quote.Asks.Size() == 0 || depth_quote.Bids.Size() == 0 {
		return false
	}

	ask_iter := depth_quote.Asks.Iterator()
	bid_iter := depth_quote.Bids.Iterator()

	if ask_iter.First() && bid_iter.Last() && ask_iter.Key().(float64) <= bid_iter.Key().(float64) {
		return true
	}
	return false
}

func (w *WatermarkWorker) Execute(depth_quote *DepthQuote, configs *RiskCtrlConfigMap) bool {
	defer ExceptionFunc()

	fmt.Println(w)

	w.Process(depth_quote, configs)

	if w.nextWorker != nil {
		w.nextWorker.Execute(depth_quote, configs)
	} else {
		fmt.Printf("%s Worker Has No Next Worker!\n", w.WorkerName)
	}

	return false
}

func (w *WatermarkWorker) Process(depth_quote *DepthQuote, configs *RiskCtrlConfigMap) bool {
	defer ExceptionFunc()

	if check_cross(depth_quote) == false {
		return true
	}

	fmt.Println(configs)
	if config, ok := (*configs)[depth_quote.Symbol]; ok {

		LOG_INFO("\nBefore WatermarkWorker: \n" + depth_quote.String(5))
		// fmt.Println(depth_quote)

		watermark := calc_watermark(depth_quote)

		if watermark <= 0 {
			return true
		}

		filter_depth_by_watermark(depth_quote.Asks, watermark, config.PriceMinumChange, true)

		filter_depth_by_watermark(depth_quote.Bids, watermark, config.PriceMinumChange*-1, false)

		LOG_INFO("\nAfter WatermarkWorker: \n" + depth_quote.String(5))
		// fmt.Println(depth_quote)

	} else {

		LOG_ERROR("Symbol: " + string(depth_quote.Symbol) + " Has no config")
		return false
	}
	return true
}

func resize_float64(src float64, presion uint32) float64 {
	defer ExceptionFunc()

	x := math.Pow10(int(presion))
	return math.Trunc(src*x) / x
}

func resize_depth_precision(depth *treemap.Map, config *RiskCtrlConfig) *treemap.Map {
	defer ExceptionFunc()

	result := treemap.NewWith(utils.Float64Comparator)
	depth_iter := depth.Iterator()

	for depth_iter.Begin(); depth_iter.Next(); {
		original_price := depth_iter.Key().(float64)
		inner_depth := depth_iter.Value().(InnerDepth)

		new_inner_depth := InnerDepth{0, make(map[string]float64)}

		new_price := resize_float64(original_price, config.PricePrecison)
		new_volume := resize_float64(inner_depth.Volume, config.VolumePrecison)

		new_inner_depth.Volume = new_volume

		for exchange, exchange_volume := range inner_depth.ExchangeVolume {
			new_exchange_volume := float64(resize_float64(float64(exchange_volume), config.VolumePrecison))

			new_inner_depth.ExchangeVolume[exchange] = new_exchange_volume
		}

		result.Put(new_price, new_inner_depth)
	}

	return result
}

func (w *PrecisionWorker) Execute(depth_quote *DepthQuote, configs *RiskCtrlConfigMap) bool {
	defer ExceptionFunc()

	fmt.Println(w)

	w.Process(depth_quote, configs)

	if w.nextWorker != nil {
		w.nextWorker.Execute(depth_quote, configs)
	} else {
		fmt.Printf("%s Worker Has No Next Worker!\n", w.WorkerName)
	}

	return false
}

func (w *PrecisionWorker) Process(depth_quote *DepthQuote, configs *RiskCtrlConfigMap) bool {
	defer ExceptionFunc()

	if config, ok := (*configs)[depth_quote.Symbol]; ok {

		fmt.Printf("\nconfig:%v \n", configs)
		LOG_INFO("\nBefore PrecisionWorker: \n" + depth_quote.String(5))

		new_asks := resize_depth_precision(depth_quote.Asks, &config)
		depth_quote.Asks = new_asks

		new_bids := resize_depth_precision(depth_quote.Bids, &config)
		depth_quote.Bids = new_bids

		LOG_INFO("\nAfter PrecisionWorker: \n" + depth_quote.String(5))

	} else {

		LOG_ERROR("Symbol: " + string(depth_quote.Symbol) + " Has no config")
		return false
	}

	fmt.Println("PrecisionWorker Process")

	return false
}

type PrecisionWorker struct {
	Worker
}

type RiskWorkerManager struct {
	FeeWorker_       FeeWorker
	QuotebiasWorker_ QuotebiasWorker
	WatermarkWorker_ WatermarkWorker
	PrecisionWorker_ PrecisionWorker
}

func (r *RiskWorkerManager) Init() {
	r.FeeWorker_.WorkerName = "FeeWorker"
	r.QuotebiasWorker_.WorkerName = "QuotebiasWorker"
	r.WatermarkWorker_.WorkerName = "WatermarkWorker"
	r.PrecisionWorker_.WorkerName = "PrecisionWorker"

	// r.FeeWorker_.SetNext(&r.QuotebiasWorker_)
	// r.QuotebiasWorker_.SetNext(&r.WatermarkWorker_)
	// r.WatermarkWorker_.SetNext(&r.PrecisionWorker_)
}

func (r *RiskWorkerManager) Execute(depth_quote *DepthQuote, configs *RiskCtrlConfigMap) {
	fmt.Printf("\n------- RiskWorkerManager  Executing -------- \n\n")

	// r.FeeWorker_.Execute(depth_quote, configs)

	// r.QuotebiasWorker_.Execute(depth_quote, configs)

	r.WatermarkWorker_.Execute(depth_quote, configs)

	// r.PrecisionWorker_.Execute(depth_quote, configs)
}

func test_get_sorted_keys() {
	defer ExceptionFunc()

	test_map := treemap.NewWith(utils.Float64Comparator)

	depth1 := InnerDepth{0, make(map[string]float64)}
	depth1.Volume = 1
	depth1.ExchangeVolume["FTX"] = 1

	depth2 := InnerDepth{0, make(map[string]float64)}
	depth2.Volume = 2

	depth3 := InnerDepth{0, make(map[string]float64)}
	depth3.Volume = 3

	test_map.Put(1.1, depth1)
	test_map.Put(2.1, depth2)
	test_map.Put(3.1, depth3)

	// test_map_iter := test_map.Iterator()

	// for price, value := range test_map {
	// 	fmt.Println(price)
	// 	fmt.Println(value)
	// 	fmt.Printf("\n")
	// }

	// keys := get_sorted_key(test_map)
	// fmt.Println(keys)
}

func GetTestDepth() DepthQuote {
	var rst DepthQuote

	rst.Exchange = "FTX"
	rst.Symbol = "BTC_USDT"
	rst.Time = uint64(time.Now().Unix())
	rst.Asks = treemap.NewWith(utils.Float64Comparator)
	rst.Bids = treemap.NewWith(utils.Float64Comparator)

	rst.Asks.Put(41001.11111, InnerDepth{1.11111, map[string]float64{"FTX": 1.11111}})
	rst.Asks.Put(41002.22222, InnerDepth{2.22222, map[string]float64{"FTX": 2.22222}})

	rst.Asks.Put(41003.33333, InnerDepth{3.33333, map[string]float64{"FTX": 3.33333}})
	rst.Asks.Put(41004.44444, InnerDepth{4.44444, map[string]float64{"FTX": 4.44444}})
	rst.Asks.Put(41005.55555, InnerDepth{5.55555, map[string]float64{"FTX": 5.55555}})

	rst.Bids.Put(41004.44444, InnerDepth{4.44444, map[string]float64{"FTX": 4.44444}})
	rst.Bids.Put(41003.33333, InnerDepth{3.33333, map[string]float64{"FTX": 3.33333}})
	rst.Bids.Put(41002.22222, InnerDepth{2.22222, map[string]float64{"FTX": 2.22222}})
	rst.Bids.Put(41001.11111, InnerDepth{1.11111, map[string]float64{"FTX": 1.11111}})
	rst.Bids.Put(40009.99999, InnerDepth{9.99999, map[string]float64{"FTX": 9.99999}})

	// rst.Asks.Put(40002.222222].Volume = 2.222222
	// rst.Asks.Put(40003.222222].Volume = 2.222222
	// rst.Asks.Put(40004.222222].Volume = 2.222222

	return rst
}

func get_test_config() RiskCtrlConfigMap {
	rst := RiskCtrlConfigMap{
		"BTC_USDT": {
			HedgeConfigMap:   map[string]HedgeConfig{"FTX": {1, 0.1}},
			PricePrecison:    2,
			VolumePrecison:   3,
			PriceBiasValue:   0.1,
			PriceBiasKind:    1,
			VolumeBiasValue:  0.1,
			VolumeBiasKind:   1,
			PriceMinumChange: 1.0,
		},
	}

	return rst
}

func TestWorker() {
	risk_worker_manager := RiskWorkerManager{}
	risk_worker_manager.Init()

	depth_quote := GetTestDepth()
	config := get_test_config()

	// fmt.Printf("depth_quote: %v\n", depth_quote)
	// fmt.Printf("config: %v\n", config)

	risk_worker_manager.Execute(&depth_quote, &config)
}

func TestInnerDepth() {
	// a := InnerDepth{0, make(map[string]float64)}

	// var e1 = map[string]float64{
	// 	"FTX": 1.1,
	// }

	// e := map[string]float64{
	// 	"FTX": 1.1,
	// }

	a := InnerDepth{0, map[string]float64{"FTX": 1.1}}

	// fmt.Println(e)
	// fmt.Println(e1)
	fmt.Println(a)
}

func test_json() {
	depth_quote := GetTestDepth()
	fmt.Println(depth_quote.String(5))
}

func TestImport() {
	data := TestData{
		Name: "Tom",
	}

	fmt.Println(data)
}

// func main() {
// 	test_worker()

// 	// test_get_sorted_keys()

// 	// test_inner_depth()

// 	// test_json()
// }
