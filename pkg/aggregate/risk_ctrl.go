package aggregate

import (
	"fmt"
	config "market_aggregate/pkg/conf"
	"market_aggregate/pkg/datastruct"
	"market_aggregate/pkg/util"
	"math"
	"sync"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
	"github.com/shopspring/decimal"
)

// type TString string

// func catch_exp() {
// 	errMsg := recover()

// 	if errMsg != nil {
// 		fmt.Println(errMsg)

// 	}

// 	fmt.Println("This is catch_exp func")

// }

type RiskCtrlConfigMap map[string]*config.RiskCtrlConfig

func GetRiskCtrlConfigMapString(r *RiskCtrlConfigMap) string {
	result := ""
	for symbol, risk_config := range *r {
		result += symbol + ":\n" + risk_config.String()
	}
	return result
}

type RiskWorkerInterface interface {
	Process(depth_quote *datastruct.DepthQuote, configs *RiskCtrlConfigMap) bool
	Execute(depth_quote *datastruct.DepthQuote, configs *RiskCtrlConfigMap) bool
	SetNext(nextWorker RiskWorkerInterface)
	GetWorkerName() string
	GetNextWoker() RiskWorkerInterface
}

type Worker struct {
	nextWorker RiskWorkerInterface

	WorkerName string
}

func (w *Worker) Process(depth_quote *datastruct.DepthQuote, configs *RiskCtrlConfigMap) bool {
	fmt.Println("Original Worker Process")
	return false
}

func (w *Worker) Execute(depth_quote *datastruct.DepthQuote, configs *RiskCtrlConfigMap) bool {
	defer util.ExceptionFunc()

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
	defer util.ExceptionFunc()

	w.nextWorker = next
}

func (w *Worker) GetWorkerName() string {
	return w.WorkerName
}

func (w *Worker) GetNextWoker() RiskWorkerInterface {
	return w.nextWorker
}

type FeeWorker struct {
	Worker
}

func get_bias_value(original_value float64, bias_kind int, bias_value float64) float64 {
	defer util.ExceptionFunc()

	rst := decimal.NewFromFloat(original_value)

	// biased_value := original_value

	if bias_kind == 1 {
		rst = rst.Mul(decimal.NewFromFloat(1 + bias_value))
	} else if bias_kind == 2 {
		rst = rst.Add(decimal.NewFromFloat(bias_value))
	} else {
		util.LOG_ERROR("Error Bias Kind")
	}

	rst_float, _ := rst.Float64()
	if rst_float < 0 {
		rst_float = 0
	}

	return rst_float
}

func get_new_fee_price(original_price float64, exchange string, hedge_config *config.RiskCtrlConfig, isAsk bool) float64 {
	defer util.ExceptionFunc()

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
func (w *FeeWorker) calc_depth_fee(depth *treemap.Map, config *config.RiskCtrlConfig, isAsk bool) *treemap.Map {
	defer util.ExceptionFunc()

	result := treemap.NewWith(utils.Float64Comparator)

	iter := depth.Iterator()

	for iter.Begin(); iter.Next(); {
		original_price := iter.Key().(float64)
		inner_depth := iter.Value().(*datastruct.InnerDepth)

		for exchange, _ := range inner_depth.ExchangeVolume {
			new_price := get_new_fee_price(original_price, exchange, config, isAsk)

			if new_inner_depth_iter, ok := result.Get(new_price); ok {

				new_inner_depth := new_inner_depth_iter.(*datastruct.InnerDepth)

				new_inner_depth.ExchangeVolume[exchange] += inner_depth.Volume
				new_inner_depth.Volume += inner_depth.Volume

			} else {
				new_inner_depth := datastruct.InnerDepth{Volume: 0, ExchangeVolume: make(map[string]float64)}

				new_inner_depth.ExchangeVolume[exchange] = inner_depth.Volume
				new_inner_depth.Volume = inner_depth.Volume

				result.Put(new_price, &new_inner_depth)
			}
		}
	}

	return result
}

func (w *FeeWorker) Execute(depth_quote *datastruct.DepthQuote, configs *RiskCtrlConfigMap) bool {
	defer util.ExceptionFunc()

	fmt.Println(w)

	w.Process(depth_quote, configs)

	if w.nextWorker != nil {
		return w.nextWorker.Execute(depth_quote, configs)
	} else {
		fmt.Printf("%s Worker Has No Next Worker!\n", w.WorkerName)
		return true
	}
}

func (w *FeeWorker) Process(depth_quote *datastruct.DepthQuote, configs *RiskCtrlConfigMap) bool {
	defer util.ExceptionFunc()

	fmt.Println("-------- FeeWorker Process  ---------")
	util.LOG_INFO(fmt.Sprintf("\n------- configs: %+v\n\n", *configs))

	if config, ok := (*configs)[depth_quote.Symbol]; ok {
		fmt.Printf("Symbol:%s, Config:%+v \n", depth_quote.Symbol, config)

		util.LOG_INFO("\nBefore FeeCtrl: \n" + depth_quote.String(5))
		// fmt.Println(depth_quote)

		new_asks := w.calc_depth_fee(depth_quote.Asks, config, true)
		depth_quote.Asks = new_asks

		new_bids := w.calc_depth_fee(depth_quote.Bids, config, false)
		depth_quote.Bids = new_bids

		util.LOG_INFO("\nAfter FeeCtrl: \n" + depth_quote.String(5))
		// fmt.Println(depth_quote)

	} else {

		util.LOG_ERROR("Symbol: " + string(depth_quote.Symbol) + " Has no config")
		return false
	}

	fmt.Println("FeeWorker Process")
	return true
}

type QuotebiasWorker struct {
	Worker
}

func calc_depth_bias(depth *treemap.Map, config *config.RiskCtrlConfig, isAsk bool) *treemap.Map {
	defer util.ExceptionFunc()

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
		inner_depth := depth_iter.Value().(*datastruct.InnerDepth)

		// var new_inner_depth datastruct.InnerDepth

		new_inner_depth := datastruct.InnerDepth{Volume: 0, ExchangeVolume: make(map[string]float64)}

		new_price := get_bias_value(original_price, config.PriceBiasKind, PriceBiasValue)
		new_volume := get_bias_value(inner_depth.Volume, config.VolumeBiasKind, VolumeBiasValue)

		new_inner_depth.Volume = new_volume

		for exchange, exchange_volume := range inner_depth.ExchangeVolume {
			new_exchange_volume := get_bias_value(exchange_volume, config.VolumeBiasKind, VolumeBiasValue)

			new_inner_depth.ExchangeVolume[exchange] = new_exchange_volume
		}

		result.Put(new_price, &new_inner_depth)
	}

	return result
}

func (w *QuotebiasWorker) Process(depth_quote *datastruct.DepthQuote, configs *RiskCtrlConfigMap) bool {
	defer util.ExceptionFunc()

	if config, ok := (*configs)[depth_quote.Symbol]; ok {

		fmt.Printf("config:%v \n", configs)
		util.LOG_INFO("\nBefore QuotebiasCtrl: \n" + depth_quote.String(5))
		// fmt.Println(depth_quote)

		new_asks := calc_depth_bias(depth_quote.Asks, config, true)
		depth_quote.Asks = new_asks

		new_bids := calc_depth_bias(depth_quote.Bids, config, false)
		depth_quote.Bids = new_bids

		util.LOG_INFO("\nAfter QuotebiasCtrl: \n" + depth_quote.String(5))
		// fmt.Println(depth_quote)

	} else {

		util.LOG_ERROR("Symbol: " + string(depth_quote.Symbol) + " Has no config")
		return false
	}

	fmt.Println("QuotebiasWorker Process")
	return true
}

func (w *QuotebiasWorker) Execute(depth_quote *datastruct.DepthQuote, configs *RiskCtrlConfigMap) bool {
	defer util.ExceptionFunc()

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
func calc_watermark(depth_quote *datastruct.DepthQuote) float64 {
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
	new_inner_depth := datastruct.InnerDepth{Volume: 0, ExchangeVolume: make(map[string]float64)}
	depth_iter := depth.Iterator()
	new_price := watermark + float64(price_minum_change)

	fmt.Printf("\nNewPrice: %v\n", new_price)

	if isAsk {

		for depth_iter.Begin(); depth_iter.Next(); {
			cur_price := depth_iter.Key().(float64)
			cur_innerdepth := depth_iter.Value().(*datastruct.InnerDepth)

			if cur_price <= new_price {
				crossed_price = append(crossed_price, cur_price)
				new_inner_depth.Add(cur_innerdepth)
			} else {
				break
			}
		}

	} else {

		for depth_iter.End(); depth_iter.Prev(); {
			cur_price := depth_iter.Key().(float64)
			cur_innerdepth := depth_iter.Value().(*datastruct.InnerDepth)

			if cur_price >= new_price {
				crossed_price = append(crossed_price, cur_price)
				new_inner_depth.Add(cur_innerdepth)
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

func check_cross(depth_quote *datastruct.DepthQuote) bool {
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

func (w *WatermarkWorker) Execute(depth_quote *datastruct.DepthQuote, configs *RiskCtrlConfigMap) bool {
	defer util.ExceptionFunc()

	fmt.Println(w)

	w.Process(depth_quote, configs)

	if w.nextWorker != nil {
		w.nextWorker.Execute(depth_quote, configs)
	} else {
		fmt.Printf("%s Worker Has No Next Worker!\n", w.WorkerName)
	}

	return false
}

func (w *WatermarkWorker) Process(depth_quote *datastruct.DepthQuote, configs *RiskCtrlConfigMap) bool {
	defer util.ExceptionFunc()

	if check_cross(depth_quote) == false {
		util.LOG_INFO("++++++ WatermarkWorker:Process Has No Cross Depth! +++++\n\n")
		return true
	}

	fmt.Println(configs)
	if config, ok := (*configs)[depth_quote.Symbol]; ok {

		util.LOG_INFO("\nBefore WatermarkWorker: \n" + depth_quote.String(5))
		// fmt.Println(depth_quote)

		watermark := calc_watermark(depth_quote)

		if watermark <= 0 {
			return true
		}

		filter_depth_by_watermark(depth_quote.Asks, watermark, config.PriceMinumChange, true)

		filter_depth_by_watermark(depth_quote.Bids, watermark, config.PriceMinumChange*-1, false)

		util.LOG_INFO("\nAfter WatermarkWorker: \n" + depth_quote.String(5))
		// fmt.Println(depth_quote)

	} else {

		util.LOG_ERROR("Symbol: " + string(depth_quote.Symbol) + " Has no config")
		return false
	}
	return true
}

func resize_float64(src float64, presion uint32) float64 {
	defer util.ExceptionFunc()

	x := math.Pow10(int(presion))
	return math.Trunc(src*x) / x
}

func resize_depth_precision(depth *treemap.Map, config *config.RiskCtrlConfig) *treemap.Map {
	defer util.ExceptionFunc()

	result := treemap.NewWith(utils.Float64Comparator)
	depth_iter := depth.Iterator()

	for depth_iter.Begin(); depth_iter.Next(); {
		original_price := depth_iter.Key().(float64)
		inner_depth := depth_iter.Value().(*datastruct.InnerDepth)

		new_inner_depth := datastruct.InnerDepth{Volume: 0, ExchangeVolume: make(map[string]float64)}

		new_price := resize_float64(original_price, config.PricePrecison)
		new_volume := resize_float64(inner_depth.Volume, config.VolumePrecison)

		new_inner_depth.Volume = new_volume

		for exchange, exchange_volume := range inner_depth.ExchangeVolume {
			new_exchange_volume := float64(resize_float64(float64(exchange_volume), config.VolumePrecison))

			new_inner_depth.ExchangeVolume[exchange] = new_exchange_volume
		}

		result.Put(new_price, &new_inner_depth)
	}

	return result
}

func (w *PrecisionWorker) Execute(depth_quote *datastruct.DepthQuote, configs *RiskCtrlConfigMap) bool {
	defer util.ExceptionFunc()

	fmt.Println(w)

	w.Process(depth_quote, configs)

	if w.nextWorker != nil {
		w.nextWorker.Execute(depth_quote, configs)
	} else {
		fmt.Printf("%s Worker Has No Next Worker!\n", w.WorkerName)
	}

	return false
}

func (w *PrecisionWorker) Process(depth_quote *datastruct.DepthQuote, configs *RiskCtrlConfigMap) bool {
	defer util.ExceptionFunc()

	if config, ok := (*configs)[depth_quote.Symbol]; ok {

		fmt.Printf("\nconfig:%v \n", configs)
		util.LOG_INFO("\nBefore PrecisionWorker: \n" + depth_quote.String(5))

		new_asks := resize_depth_precision(depth_quote.Asks, config)
		depth_quote.Asks = new_asks

		new_bids := resize_depth_precision(depth_quote.Bids, config)
		depth_quote.Bids = new_bids

		util.LOG_INFO("\nAfter PrecisionWorker: \n" + depth_quote.String(5))

	} else {

		util.LOG_ERROR("Symbol: " + string(depth_quote.Symbol) + " Has no config")
		return false
	}

	fmt.Println("PrecisionWorker Process")

	return true
}

type PrecisionWorker struct {
	Worker
}

type RiskWorkerManager struct {
	FeeWorker_       FeeWorker
	QuotebiasWorker_ QuotebiasWorker
	WatermarkWorker_ WatermarkWorker
	PrecisionWorker_ PrecisionWorker

	Worker      RiskWorkerInterface
	RiskConfig  RiskCtrlConfigMap
	ConfigMutex *sync.RWMutex
}

func (r *RiskWorkerManager) Init() {
	// r.FeeWorker_.WorkerName = "FeeWorker"
	// r.QuotebiasWorker_.WorkerName = "QuotebiasWorker"
	// r.WatermarkWorker_.WorkerName = "WatermarkWorker"
	// r.PrecisionWorker_.WorkerName = "PrecisionWorker"

	fee_worker := &FeeWorker{
		Worker{WorkerName: "FeeWorker"},
	}

	quotebias_worker := &QuotebiasWorker{
		Worker{WorkerName: "QuotebiasWorker"},
	}

	watermark_worker := &WatermarkWorker{
		Worker{WorkerName: "WatermarkWorker"},
	}

	precision_worker := &PrecisionWorker{
		Worker{WorkerName: "PrecisionWorker"},
	}

	r.ConfigMutex = new(sync.RWMutex)
	r.RiskConfig = make(map[string]*config.RiskCtrlConfig)
	r.Worker = nil

	if config.TESTCONFIG().FeeRiskctrlOpen {
		r.AddWorker(fee_worker)
	}

	if config.TESTCONFIG().BiasRiskctrlOpen {
		r.AddWorker(quotebias_worker)
	}

	if config.TESTCONFIG().WatermarkRiskctrlOpen {
		r.AddWorker(watermark_worker)
	}

	if config.TESTCONFIG().PricesionRiskctrlOpen {
		r.AddWorker(precision_worker)
	}

	// r.FeeWorker_.SetNext(&r.QuotebiasWorker_)
	// r.QuotebiasWorker_.SetNext(&r.WatermarkWorker_)
	// r.WatermarkWorker_.SetNext(&r.PrecisionWorker_)
}

func (r *RiskWorkerManager) UpdateConfig(RiskConfig *RiskCtrlConfigMap) {
	defer r.ConfigMutex.Unlock()
	r.ConfigMutex.Lock()

	for symbol, value := range *RiskConfig {
		// r.RiskConfig[symbol] = value

		r.RiskConfig[symbol] = &config.RiskCtrlConfig{
			HedgeConfigMap: value.HedgeConfigMap,

			PricePrecison:  value.PricePrecison,
			VolumePrecison: value.VolumePrecison,

			PriceBiasValue: value.PriceBiasValue,
			PriceBiasKind:  value.PriceBiasKind,

			VolumeBiasValue: value.VolumeBiasValue,
			VolumeBiasKind:  value.VolumeBiasKind,

			PriceMinumChange: value.PriceMinumChange,
		}
	}

	util.LOG_INFO(fmt.Sprintf("\n------- r.RiskConfig: %+v\n\n", r.RiskConfig))
}

func (r *RiskWorkerManager) AddWorker(NewWorker RiskWorkerInterface) {
	util.LOG_INFO("Try Add Worker " + NewWorker.GetWorkerName())
	if r.Worker == nil {
		r.Worker = NewWorker
		util.LOG_INFO("Init First Worker " + NewWorker.GetWorkerName() + "\n")
		return
	}

	var tmp RiskWorkerInterface
	for tmp = r.Worker; tmp.GetNextWoker() != nil; tmp = tmp.GetNextWoker() {

		// time.Sleep(time.Second * 3)
		util.LOG_INFO("Stored Worker " + tmp.GetWorkerName())

		if tmp.GetWorkerName() == NewWorker.GetWorkerName() {
			util.LOG_WARN("Repeated Worker : " + tmp.GetWorkerName())
			return
		}
	}
	tmp.SetNext(NewWorker)
	util.LOG_INFO("Add Worker " + NewWorker.GetWorkerName() + "\n")
}

func (r *RiskWorkerManager) Execute(depth_quote *datastruct.DepthQuote) {
	defer r.ConfigMutex.RUnlock()

	fmt.Printf("\n------- RiskWorkerManager  Executing -------- \n\n")

	r.ConfigMutex.RLock()

	if r.Worker == nil {
		util.LOG_ERROR("No Worker Available")
		return
	}

	if len(r.RiskConfig) == 0 {
		util.LOG_ERROR("RiskConfig Not Available")
		return
	}

	r.Worker.Execute(depth_quote, &r.RiskConfig)
}

func GetTestRiskConfig() RiskCtrlConfigMap {
	rst := RiskCtrlConfigMap{
		"BTC_USDT": {
			HedgeConfigMap: map[string]*config.HedgeConfig{"FTX": &config.HedgeConfig{FeeKind: 1, FeeValue: 0.1},
				"OKEX":  {FeeKind: 1, FeeValue: 0.2},
				"HUOBI": {FeeKind: 1, FeeValue: 0.3}},
			PricePrecison:    2,
			VolumePrecison:   3,
			PriceBiasValue:   0.1,
			PriceBiasKind:    1,
			VolumeBiasValue:  0.1,
			VolumeBiasKind:   1,
			PriceMinumChange: 1.0,
		},
		"ETH_USDT": {
			HedgeConfigMap:   map[string]*config.HedgeConfig{"FTX": {FeeKind: 1, FeeValue: 0.1}, "OKEX": {FeeKind: 1, FeeValue: 0.2}, "HUOBI": {FeeKind: 1, FeeValue: 0.3}},
			PricePrecison:    2,
			VolumePrecison:   3,
			PriceBiasValue:   0.1,
			PriceBiasKind:    1,
			VolumeBiasValue:  0.1,
			VolumeBiasKind:   1,
			PriceMinumChange: 1.0,
		},
		"DOT_USDT": {
			HedgeConfigMap:   map[string]*config.HedgeConfig{"FTX": {FeeKind: 1, FeeValue: 0.1}, "OKEX": {FeeKind: 1, FeeValue: 0.2}, "HUOBI": {FeeKind: 1, FeeValue: 0.3}},
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
	depth_quote := datastruct.GetTestDepth()
	config := GetTestRiskConfig()

	risk_worker_manager := RiskWorkerManager{}
	risk_worker_manager.Init()

	risk_worker_manager.UpdateConfig(&config)

	// fmt.Printf("depth_quote: %v\n", depth_quote)
	// fmt.Printf("config: %v\n", config)

	risk_worker_manager.Execute(depth_quote)
}

func TestInnerDepth() {
	// a := datastruct.InnerDepth{0, make(map[string]float64)}

	// var e1 = map[string]float64{
	// 	"FTX": 1.1,
	// }

	// e := map[string]float64{
	// 	"FTX": 1.1,
	// }

	a := datastruct.InnerDepth{0, map[string]float64{"FTX": 1.1}}

	// fmt.Println(e)
	// fmt.Println(e1)
	fmt.Println(a)
}

func test_json() {
	depth_quote := datastruct.GetTestDepth()
	fmt.Println(depth_quote.String(5))
}

func TestImport() {
	data := datastruct.TestData{
		Name: "Tom",
	}

	fmt.Println(data)
}

func TestAddWorker() {
	risk_work := RiskWorkerManager{}
	risk_work.Init()
}

// func main() {
// 	test_worker()

// 	// test_get_sorted_keys()

// 	// test_inner_depth()

// 	// test_json()
// }
