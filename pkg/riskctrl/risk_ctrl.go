package riskctrl

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
)

// type TString string

type TSymbol string
type TExchange string

// type RFloat float64
type TPrice float64
type TVolume float64

type InnerDepth struct {
	Volume         TVolume
	ExchangeVolume map[TExchange]TVolume
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
	Exchange TExchange             `json:"Exchange"`
	Symbol   TSymbol               `json:"Symbol"`
	Time     uint64                `json:"Time"`
	Asks     map[TPrice]InnerDepth `json:"Asks"`
	Bids     map[TPrice]InnerDepth `json:"Bids"`
}

type TestTreeMap struct {
	data treemap.Map
}

func (depth_quote *DepthQuote) String(len int) string {

	res, _ := json.Marshal(*depth_quote)

	return string(res)
}

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
	HedgeConfigMap map[TExchange]HedgeConfig

	PricePrecison  uint32
	VolumePrecison uint32

	PriceBiasValue float64
	PriceBiasKind  int

	VolumeBiasValue float64
	VolumeBiasKind  int

	PriceMinumChange float64
}

type RiskCtrlConfigMap map[TSymbol]RiskCtrlConfig

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

	biased_value := original_value

	if bias_kind == 1 {
		biased_value *= (1 + bias_value)
	} else if bias_kind == 2 {
		biased_value += bias_value
	} else {
		LOG_ERROR("Error Bias Kind")
	}

	if biased_value < 0 {
		biased_value = 0
	}

	return biased_value
}

func get_new_fee_price(original_price TPrice, exchange TExchange, hedge_config *RiskCtrlConfig, isAsk bool) TPrice {
	defer ExceptionFunc()

	new_price := original_price
	if hedge_config, ok := hedge_config.HedgeConfigMap[exchange]; ok {

		if isAsk {
			new_price = TPrice(get_bias_value(float64(original_price), hedge_config.FeeKind, hedge_config.FeeValue))
		} else {
			new_price = TPrice(get_bias_value(float64(original_price), hedge_config.FeeKind, hedge_config.FeeValue*-1))
		}
	}
	return new_price
}

// 根据每个交易所的手续费率，计算手续费
func (w *FeeWorker) calc_depth_fee(depth map[TPrice]InnerDepth, config *RiskCtrlConfig, isAsk bool) map[TPrice]InnerDepth {
	defer ExceptionFunc()

	result := make(map[TPrice]InnerDepth)

	for original_price, inner_depth := range depth {
		for exchange, _ := range inner_depth.ExchangeVolume {
			new_price := get_new_fee_price(original_price, exchange, config, isAsk)

			if new_inner_depth, ok := result[new_price]; ok {
				new_inner_depth.ExchangeVolume[exchange] += inner_depth.Volume
				new_inner_depth.Volume += inner_depth.Volume
			} else {
				new_inner_depth := InnerDepth{0, make(map[TExchange]TVolume)}

				new_inner_depth.ExchangeVolume[exchange] = inner_depth.Volume
				new_inner_depth.Volume = inner_depth.Volume

				result[new_price] = new_inner_depth
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

		fmt.Printf("config:%v \n", configs)
		LOG_INFO("\nBefore FeeCtrl: \n" + depth_quote.String(5))
		fmt.Println(depth_quote)

		new_asks := w.calc_depth_fee(depth_quote.Asks, &config, true)
		depth_quote.Asks = new_asks

		new_bids := w.calc_depth_fee(depth_quote.Bids, &config, false)
		depth_quote.Bids = new_bids

		LOG_INFO("\nAfter FeeCtrl: \n" + depth_quote.String(5))
		fmt.Println(depth_quote)

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

func calc_depth_bias(depth map[TPrice]InnerDepth, config *RiskCtrlConfig, isAsk bool) map[TPrice]InnerDepth {
	defer ExceptionFunc()

	result := make(map[TPrice]InnerDepth)

	for original_price, inner_depth := range depth {

		// var new_inner_depth InnerDepth

		new_inner_depth := InnerDepth{0, make(map[TExchange]TVolume)}

		new_price := TPrice(get_bias_value(float64(original_price), config.PriceBiasKind, config.PriceBiasValue))
		new_volume := TVolume(get_bias_value(float64(inner_depth.Volume), config.VolumeBiasKind, config.VolumeBiasValue))

		new_inner_depth.Volume = new_volume

		for exchange, exchange_volume := range inner_depth.ExchangeVolume {
			new_exchange_volume := TVolume(get_bias_value(float64(exchange_volume), config.VolumeBiasKind, config.VolumeBiasValue))

			new_inner_depth.ExchangeVolume[exchange] = new_exchange_volume
		}

		result[new_price] = new_inner_depth
	}

	return result
}

func (w *QuotebiasWorker) Process(depth_quote *DepthQuote, configs *RiskCtrlConfigMap) bool {
	defer ExceptionFunc()

	if config, ok := (*configs)[depth_quote.Symbol]; ok {

		fmt.Printf("config:%v \n", configs)
		LOG_INFO("\nBefore QuotebiasCtrl: \n" + depth_quote.String(5))
		fmt.Println(depth_quote)

		new_asks := calc_depth_bias(depth_quote.Asks, &config, true)
		depth_quote.Asks = new_asks

		new_bids := calc_depth_bias(depth_quote.Bids, &config, false)
		depth_quote.Bids = new_bids

		LOG_INFO("\nAfter QuotebiasCtrl: \n" + depth_quote.String(5))
		fmt.Println(depth_quote)

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
func calc_watermark(depth_quote *DepthQuote) TPrice {
	var rst TPrice

	return rst
}

func filter_depth_by_watermark(depth map[TPrice]InnerDepth, watermark TPrice, price_minum_change float64, isAsk bool) {

	sorted_keys := get_sorted_key(depth)
	crossed_price := []TPrice{}
	new_inner_depth := InnerDepth{0, make(map[TExchange]TVolume)}
	var new_price TPrice = 0

	if isAsk {
		new_price = watermark + TPrice(price_minum_change)

		for i := 0; i < len(sorted_keys); i++ {
			cur_price := sorted_keys[i]
			cur_innerdepth := depth[cur_price]
			if cur_price <= new_price {
				crossed_price = append(crossed_price, cur_price)

				new_inner_depth.Add(&cur_innerdepth)

			} else {
				break
			}
		}

	} else {
		new_price = watermark - TPrice(price_minum_change)

		for i := len(sorted_keys); i > -1; i-- {
			cur_price := sorted_keys[i]
			cur_innerdepth := depth[cur_price]
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
			delete(depth, price)
		}
	}

	if new_price != 0 {
		depth[new_price] = new_inner_depth
	}
}

func get_sorted_key(depth map[TPrice]InnerDepth) []TPrice {
	var keys []TPrice

	for price, _ := range depth {

		keys = append(keys, price)

		for i := len(keys) - 1; i > 0; i-- {
			if keys[i] < keys[i-1] {
				tmp := keys[i]
				keys[i] = keys[i-1]
				keys[i-1] = tmp
			} else {
				break
			}
		}
	}
	return keys
}

func check_cross(depth_quote *DepthQuote) bool {
	if len(depth_quote.Asks) == 0 || len(depth_quote.Bids) == 0 {
		return false
	}

	ask_keys := get_sorted_key(depth_quote.Asks)

	bid_keys := get_sorted_key(depth_quote.Bids)

	if ask_keys[0] <= bid_keys[len(bid_keys)-1] {
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

	if config, ok := (*configs)[depth_quote.Symbol]; ok {

		LOG_INFO("\nBefore WatermarkWorker: \n" + depth_quote.String(5))
		fmt.Println(depth_quote)

		watermark := calc_watermark(depth_quote)

		if watermark <= 0 {
			return true
		}

		filter_depth_by_watermark(depth_quote.Asks, watermark, config.PriceMinumChange, true)

		filter_depth_by_watermark(depth_quote.Bids, watermark, config.PriceMinumChange*-1, false)

		LOG_INFO("\nAfter WatermarkWorker: \n" + depth_quote.String(5))
		fmt.Println(depth_quote)

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

func resize_depth_precision(depth map[TPrice]InnerDepth, config *RiskCtrlConfig) map[TPrice]InnerDepth {
	defer ExceptionFunc()

	result := make(map[TPrice]InnerDepth)

	for original_price, inner_depth := range depth {

		new_inner_depth := InnerDepth{0, make(map[TExchange]TVolume)}

		new_price := TPrice(resize_float64(float64(original_price), config.PricePrecison))
		new_volume := TVolume(resize_float64(float64(inner_depth.Volume), config.VolumePrecison))

		new_inner_depth.Volume = new_volume

		for exchange, exchange_volume := range inner_depth.ExchangeVolume {
			new_exchange_volume := TVolume(resize_float64(float64(exchange_volume), config.VolumePrecison))

			new_inner_depth.ExchangeVolume[exchange] = new_exchange_volume
		}

		result[new_price] = new_inner_depth
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

		fmt.Printf("config:%v \n", configs)
		LOG_INFO("\nBefore PrecisionWorker: \n" + depth_quote.String(5))
		fmt.Println(depth_quote)

		new_asks := resize_depth_precision(depth_quote.Asks, &config)
		depth_quote.Asks = new_asks

		new_bids := resize_depth_precision(depth_quote.Bids, &config)
		depth_quote.Bids = new_bids

		LOG_INFO("\nAfter PrecisionWorker: \n" + depth_quote.String(5))
		fmt.Println(depth_quote)

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

	r.QuotebiasWorker_.Execute(depth_quote, configs)

	// r.WatermarkWorker_.Execute(depth_quote, configs)

	// r.PrecisionWorker_.Execute(depth_quote, configs)
}

func test_get_sorted_keys() {
	defer ExceptionFunc()

	test_map := make(map[TPrice]InnerDepth)

	depth1 := InnerDepth{0, make(map[TExchange]TVolume)}
	depth1.Volume = 1
	depth1.ExchangeVolume["FTX"] = 1

	depth2 := InnerDepth{0, make(map[TExchange]TVolume)}
	depth2.Volume = 2

	depth3 := InnerDepth{0, make(map[TExchange]TVolume)}
	depth3.Volume = 3

	test_map[1.1] = depth1

	test_map[2.1] = depth2

	test_map[3.1] = depth3

	for price, value := range test_map {
		fmt.Println(price)
		fmt.Println(value)
		fmt.Printf("\n")
	}

	keys := get_sorted_key(test_map)
	fmt.Println(keys)
}

func get_test_depth() DepthQuote {
	var rst DepthQuote

	rst.Exchange = "FTX"
	rst.Symbol = "BTC_USDT"
	rst.Time = uint64(time.Now().Unix())
	rst.Asks = make(map[TPrice]InnerDepth)
	rst.Bids = make(map[TPrice]InnerDepth)

	rst.Asks[41001.11111] = InnerDepth{1.11111, map[TExchange]TVolume{"FTX": 1.11111}}
	rst.Asks[41002.22222] = InnerDepth{2.22222, map[TExchange]TVolume{"FTX": 2.22222}}
	rst.Asks[41003.33333] = InnerDepth{3.33333, map[TExchange]TVolume{"FTX": 3.33333}}
	rst.Asks[41004.44444] = InnerDepth{4.44444, map[TExchange]TVolume{"FTX": 4.44444}}
	rst.Asks[41005.55555] = InnerDepth{5.55555, map[TExchange]TVolume{"FTX": 5.55555}}

	rst.Bids[41004.44444] = InnerDepth{4.44444, map[TExchange]TVolume{"FTX": 4.44444}}
	rst.Bids[41003.33333] = InnerDepth{3.33333, map[TExchange]TVolume{"FTX": 3.33333}}
	rst.Bids[41002.22222] = InnerDepth{2.22222, map[TExchange]TVolume{"FTX": 2.22222}}
	rst.Bids[41001.11111] = InnerDepth{1.11111, map[TExchange]TVolume{"FTX": 1.11111}}
	rst.Bids[40009.99999] = InnerDepth{9.99999, map[TExchange]TVolume{"FTX": 9.99999}}

	// rst.Asks[40002.222222].Volume = 2.222222
	// rst.Asks[40003.222222].Volume = 2.222222
	// rst.Asks[40004.222222].Volume = 2.222222

	return rst
}

func get_test_config() RiskCtrlConfigMap {
	rst := RiskCtrlConfigMap{
		"BTC_USDT": {
			HedgeConfigMap:   map[TExchange]HedgeConfig{"FTX": {1, 0.1}},
			PricePrecison:    2,
			VolumePrecison:   3,
			PriceBiasValue:   0.1,
			PriceBiasKind:    1,
			VolumeBiasValue:  0.1,
			VolumeBiasKind:   1,
			PriceMinumChange: 100.0,
		},
	}

	return rst
}

func TestWorker() {
	risk_worker_manager := RiskWorkerManager{}
	risk_worker_manager.Init()

	depth_quote := get_test_depth()
	config := get_test_config()

	fmt.Printf("depth_quote: %v\n", depth_quote)
	fmt.Printf("config: %v\n", config)

	risk_worker_manager.Execute(&depth_quote, &config)
}

func TestInnerDepth() {
	// a := InnerDepth{0, make(map[TExchange]TVolume)}

	// var e1 = map[TExchange]TVolume{
	// 	"FTX": 1.1,
	// }

	// e := map[string]float64{
	// 	"FTX": 1.1,
	// }

	a := InnerDepth{0, map[TExchange]TVolume{"FTX": 1.1}}

	// fmt.Println(e)
	// fmt.Println(e1)
	fmt.Println(a)
}

func test_json() {
	depth_quote := get_test_depth()
	fmt.Println(depth_quote.String(5))
}

// func main() {
// 	test_worker()

// 	// test_get_sorted_keys()

// 	// test_inner_depth()

// 	// test_json()
// }
