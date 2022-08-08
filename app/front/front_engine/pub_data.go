package front_engine

import (
	"encoding/json"
	"fmt"
	"market_server/app/front/net"
	"market_server/common/datastruct"
	"market_server/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type PubSymbolistJson struct {
	TypeInfo   string   `json:"type"`
	SymbolList []string `json:"symbol"`
}

func NewSymbolListMsg(symbol_list []string) []byte {
	json_data := PubSymbolistJson{
		TypeInfo:   net.SYMBOL_UPDATE,
		SymbolList: symbol_list,
	}

	rst, err := json.Marshal(json_data)

	if err != nil {
		logx.Errorf("NewSymbolListMsg: json_data: %+v, symbol_list: %+v, error: %s",
			json_data, symbol_list, err.Error())
	} else {
		// logx.Slowf("SymbolList json_data: %+v ", json_data)
	}

	return rst
}

/*
   string result;
   nlohmann::json json_data;
   json_data["symbol"] = string(depth_data_.symbol);
   json_data["exchange"] = string(depth_data_.exchange);
   json_data["tick"] = depth_data_.tick;
   json_data["seqno"] = depth_data_.seqno;
   json_data["ask_length"] = depth_data_.ask_length;
   json_data["bid_length"] = depth_data_.bid_length;

   nlohmann::json asks_json;
   for (int i = 0; i < depth_data_.ask_length && i < DEPCH_LEVEL_COUNT; ++i)
   {
       nlohmann::json depth_level_atom;
       depth_level_atom[0] = append_zero(depth_data_.asks[i].price.get_str_value(), depth_data_.precise);
       depth_level_atom[1] = append_zero(depth_data_.asks[i].volume.get_str_value(), depth_data_.vprecise);
       depth_level_atom[2] = append_zero(ask_accumulated_volume_[i].get_str_value(), depth_data_.vprecise);
       asks_json[i] = depth_level_atom;
   }
   json_data["asks"] = asks_json;

   nlohmann::json bids_json;
   for (int i = 0; i < depth_data_.bid_length && i < DEPCH_LEVEL_COUNT; ++i)
   {
       nlohmann::json depth_level_atom;
       depth_level_atom[0] = append_zero(depth_data_.bids[i].price.get_str_value(), depth_data_.precise);
       depth_level_atom[1] = append_zero(depth_data_.bids[i].volume.get_str_value(), depth_data_.vprecise);
       depth_level_atom[2] = append_zero(bid_accumulated_volume_[i].get_str_value(), depth_data_.vprecise);
       bids_json[i] = depth_level_atom;
   }
   json_data["bids"] = bids_json;
   json_data["type"] = MARKET_DATA_UPDATE;

   result = json_data.dump();
*/
type PubDepthJson struct {
	TypeInfo  string      `json:"type"`
	Symbol    string      `json:"symbol"`
	Exchange  string      `json:"exchange"`
	AskLength int         `json:"ask_length"`
	BidLength int         `json:"bid_length"`
	Asks      [][]float64 `json:"asks"`
	Bids      [][]float64 `json:"bids"`
	Time      int64       `json:"tick"`
}

func NewDepthJsonMsg(depth *datastruct.DepthQuote) []byte {
	var rst []byte
	var asks [][]float64
	var bids [][]float64

	// asks := make([][]float64, depth.Asks.Size())
	// bids := make([][]float64, depth.Bids.Size())

	volume_sum := float64(-1)
	ask_iter := depth.Asks.Iterator()
	for ask_iter.Begin(); ask_iter.Next(); {

		price := ask_iter.Key().(float64)
		volume := ask_iter.Value().(*datastruct.InnerDepth).Volume
		if volume_sum < 0 {
			volume_sum = volume
		} else {
			volume_sum += volume
		}

		tmp_data := []float64{price, volume, volume_sum}
		asks = append(asks, tmp_data)
	}

	volume_sum = float64(-1)
	bid_iter := depth.Bids.Iterator()
	for bid_iter.Begin(); bid_iter.Next(); {
		price := bid_iter.Key().(float64)
		volume := bid_iter.Value().(*datastruct.InnerDepth).Volume
		if volume_sum < 0 {
			volume_sum = volume
		} else {
			volume_sum += volume
		}

		tmp_data := []float64{price, volume, volume_sum}
		bids = append(bids, tmp_data)
	}

	json_data := PubDepthJson{
		TypeInfo:  net.DEPTH_UPDATE,
		Symbol:    depth.Symbol,
		Exchange:  depth.Exchange,
		AskLength: depth.Asks.Size(),
		BidLength: depth.Bids.Size(),
		Asks:      asks,
		Bids:      bids,
		Time:      depth.Time,
	}

	rst, err := json.Marshal(json_data)

	if err != nil {
		logx.Errorf("NewDepthJsonMsg: depth: %s, json_data: %+v, error: %s",
			depth.String(3), json_data, err.Error())
	} else {
		// logx.Slowf("Depth json_data: %+v ", json_data)
	}

	return rst
}

/*
json_data["type"] = TRADE;
json_data["symbol"] = string(symbol_);
json_data["price"] = price_.get_str_value();
json_data["volume"] = volume_.get_str_value();
json_data["change"] = std::to_string(change_);
json_data["change_rate"] = std::to_string(change_rate_);
json_data["high"] = high_.get_str_value();
json_data["low"] = low_.get_str_value();
*/
type PubTradeJson struct {
	TypeInfo        string  `json:"type"`
	Symbol          string  `json:"symbol"`
	Price           float64 `json:"price"`
	Volume          float64 `json:"volume"`
	Change          float64 `json:"change"`
	ChangeRate      float64 `json:"change_rate"`
	High            float64 `json:"high"`
	Low             float64 `json:"low"`
	Time            int64   `json:"date"`
	USDPrice        float64 `json:"usdPrice"`
	ReqWSTime       int64   `json:"req_ws_time"`
	ReqProcessTime  int64   `json:"req_process_time"`
	ReqResponseTime int64   `json:"req_response_time"`
	Sequence        uint64  `json:"sequence"`
}

func (t *PubTradeJson) String() string {
	return fmt.Sprintf("%s, %s, %d,v : %f,p: %f,up: %f,c: %f,cr: %f,h: %f,l:%f",
		t.Symbol, util.TimeStrFromInt(t.Time*datastruct.NANO_PER_SECS), t.Sequence,
		t.Volume, t.Price, t.USDPrice, t.Change, t.ChangeRate, t.High, t.Low)
}

func NewTradeJsonMsg(trade *datastruct.RspTrade) []byte {

	json_data := PubTradeJson{
		TypeInfo:        net.TRADE_UPATE,
		Symbol:          trade.TradeData.Symbol,
		Price:           trade.TradeData.Price,
		Volume:          trade.TradeData.Volume,
		Time:            trade.TradeData.Time / datastruct.NANO_PER_SECS,
		USDPrice:        trade.UsdPrice,
		ReqWSTime:       trade.ReqWSTime,
		ReqProcessTime:  util.UTCNanoTime() - trade.ReqArriveTime,
		ReqResponseTime: util.UTCNanoTime(),
		Sequence:        trade.TradeData.Sequence,
	}

	if trade.ChangeData != nil {
		json_data.Change = trade.ChangeData.Change
		json_data.ChangeRate = trade.ChangeData.ChangeRate
		json_data.High = trade.ChangeData.High
		json_data.Low = trade.ChangeData.Low
	}

	var rst []byte

	rst, err := json.Marshal(json_data)

	if err != nil {
		logx.Errorf("NewTradeJsonMsg: trade: %s, change_info: %s, error: %s",
			trade.TradeData.String(), trade.ChangeData.String(), err.Error())
	} else {
		// logx.Slowf("Trade %s, ReqProcessTime: %d, ReqResponseTime: %d",
		// 	json_data.Symbol, json_data.ReqProcessTime, json_data.ReqResponseTime)
	}

	return rst
}

/*
   string result;
   nlohmann::json json_data;
   if (is_update_)
   {
       json_data["type"] = KLINE_SUB;
   }
   else
   {
       json_data["type"] = KLINE_RSP;
   }

   json_data["symbol"] = string(symbol_);
   json_data["start_time"] = start_time_;
   json_data["end_time"] = end_time_;
   json_data["frequency"] = frequency_;
   json_data["data_count"] = data_count_;

   int i = 0;
   nlohmann::json detail_data;
   for (KlineDataPtr atom_data:kline_data_vec_)
   {
       nlohmann::json tmp_json;
       tmp_json["open"] = atom_data->px_open.get_value();
       tmp_json["high"] = atom_data->px_high.get_value();
       tmp_json["low"] = atom_data->px_low.get_value();
       tmp_json["close"] = atom_data->px_close.get_value();
       tmp_json["volume"] = atom_data->volume.get_value();
       tmp_json["tick"] = atom_data->index;
       detail_data[i++] = tmp_json;
   }
   json_data["data"] = detail_data;
   return json_data.dump();
*/

type PubKlineDetail struct {
	Open       float64 `json:"open"`
	High       float64 `json:"high"`
	Low        float64 `json:"low"`
	Close      float64 `json:"close"`
	Volume     float64 `json:"volume"`
	Tick       int64   `json:"tick"`
	Sequence   uint64  `json:"sequence"`
	LastVolume float64 `json:"last_volume"`
}

type PubKlineJson struct {
	TypeInfo        string           `json:"type"`
	DataType        string           `json:"data_type"`
	Symbol          string           `json:"symbol"`
	StartTime       int64            `json:"start_time"`
	EndTime         int64            `json:"end_time"`
	Resolution      uint64           `json:"frequency"`
	DataCount       int              `json:"data_count"`
	Data            []PubKlineDetail `json:"data"`
	ReqProcessTime  int64            `json:"req_process_time"`
	ReqResponseTime int64            `json:"req_response_time"`
}

func (p *PubKlineJson) TimeList() string {
	rst := ""
	for _, kline_detail := range p.Data {
		rst = rst + fmt.Sprintf("%s, %d,lv: %f,c: %f,o: %f,h: %f,l: %f,v: %f;\n",
			util.TimeStrFromInt(kline_detail.Tick*datastruct.NANO_PER_SECS),
			kline_detail.Sequence, kline_detail.LastVolume, kline_detail.Close,
			kline_detail.Open, kline_detail.High,
			kline_detail.Low, kline_detail.Volume)
	}
	rst = rst[0 : len(rst)-1]
	return rst
}

func (p *PubKlineJson) UTCTimeList() string {
	rst := ""
	for _, kline_detail := range p.Data {
		rst = rst + fmt.Sprintf("%s, %d,lv: %f,c: %f,o: %f,h: %f,l: %f,v: %f;\n",
			util.TimeStrFromInt(kline_detail.Tick*datastruct.NANO_PER_SECS),
			kline_detail.Sequence, kline_detail.LastVolume, kline_detail.Close,
			kline_detail.Open, kline_detail.High,
			kline_detail.Low, kline_detail.Volume)
	}

	rst = rst[0 : len(rst)-1]
	return rst
}

func NewHistKlineJsonMsg(hist_kline *datastruct.RspHistKline) []byte {

	var kline_data []PubKlineDetail

	for _, tmp_kline := range hist_kline.Klines {
		tmp_detail := PubKlineDetail{
			Open:       tmp_kline.Open,
			High:       tmp_kline.High,
			Low:        tmp_kline.Low,
			Close:      tmp_kline.Close,
			Volume:     tmp_kline.Volume,
			Tick:       tmp_kline.Time / datastruct.NANO_PER_SECS,
			Sequence:   tmp_kline.Sequence,
			LastVolume: tmp_kline.LastVolume,
		}

		kline_data = append(kline_data, tmp_detail)
	}

	if hist_kline.ReqInfo.Frequency > datastruct.NANO_PER_SECS {
		hist_kline.ReqInfo.Frequency = hist_kline.ReqInfo.Frequency / datastruct.NANO_PER_SECS
	}

	json_data := PubKlineJson{
		TypeInfo:        net.KLINE_UPATE,
		DataType:        net.KLINE_HIST,
		Symbol:          hist_kline.ReqInfo.Symbol,
		StartTime:       int64(hist_kline.ReqInfo.StartTime),
		EndTime:         int64(hist_kline.ReqInfo.EndTime),
		Resolution:      hist_kline.ReqInfo.Frequency,
		DataCount:       len(hist_kline.Klines),
		Data:            kline_data,
		ReqProcessTime:  util.UTCNanoTime() - hist_kline.ReqInfo.ReqArriveTime,
		ReqResponseTime: util.UTCNanoTime(),
	}

	var rst []byte

	rst, err := json.Marshal(json_data)

	if err != nil {
		logx.Errorf("NewHistKlineJsonMsg: hist_kline: %+v, json_data: %+v, error: %s",
			hist_kline, json_data, err.Error())
	} else {
		// logx.Slowf("klinersp json_data: %+v ", json_data)
	}

	return rst
}

func NewKlineUpdateJsonMsg(kline *datastruct.Kline) []byte {
	var kline_data []PubKlineDetail
	tmp_detail := PubKlineDetail{
		Open:       kline.Open,
		High:       kline.High,
		Low:        kline.Low,
		Close:      kline.Close,
		Volume:     kline.Volume,
		Tick:       kline.Time / datastruct.NANO_PER_SECS,
		Sequence:   kline.Sequence,
		LastVolume: kline.LastVolume,
	}
	kline_data = append(kline_data, tmp_detail)

	json_data := PubKlineJson{
		TypeInfo:        net.KLINE_UPATE,
		DataType:        net.KLINE_REAL,
		Symbol:          kline.Symbol,
		StartTime:       kline.Time,
		EndTime:         kline.Time,
		Resolution:      kline.Resolution,
		DataCount:       1,
		Data:            kline_data,
		ReqResponseTime: util.UTCNanoTime(),
	}

	var rst []byte

	rst, err := json.Marshal(json_data)

	if err != nil {
		logx.Errorf("NewKlineUpdateJsonMsg: kline: %+v, json_data: %+v, error: %s",
			kline, json_data, err.Error())
	} else {
		// logx.Slowf("klineupdate json_data: %+v ", json_data)
	}

	return rst
}
