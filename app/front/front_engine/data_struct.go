package front_engine

import (
	"encoding/json"
	"fmt"
	"market_server/app/front/net"
	"market_server/common/datastruct"
	"sync"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

type DepthPubInfo struct {
	ws_info *net.WSInfo
	data    []byte
}

func (d *DepthPubInfo) String() string {
	// return fmt.Sprintf("ws_info: %s, depth: %s", d.ws_info.String(), d.data.String(3))

	return fmt.Sprintf("ws_info: %s, depth: %s", d.ws_info.String(), string(d.data))
}

type SymbolPubInfo struct {
	ws_info *net.WSInfo
	data    []byte
}

func (d *SymbolPubInfo) String() string {
	return fmt.Sprintf("ws_info: %s, symbol_list: %+v", d.ws_info.String(), string(d.data))
}

type TradePubInfo struct {
	ws_info *net.WSInfo
	data    []byte
}

func (d *TradePubInfo) String() string {
	// return fmt.Sprintf("ws_info: %s, trade: %s", d.ws_info.String(), d.data.String())

	return fmt.Sprintf("[ws]:%s,[T]:%s", d.ws_info.String(), string(d.data))
}

type KlinePubInfo struct {
	ws_info    *net.WSInfo
	Symbol     string
	Resolution int
	data       []byte
}

func (d *KlinePubInfo) String() string {
	var kline_data PubKlineJson
	if err := json.Unmarshal([]byte(d.data), &kline_data); err != nil {
		logx.Errorf("Error = %+v", err)
		return fmt.Sprintf("[ws]:%s,[k]: nil", d.ws_info.String())
	} else {
		// logx.Infof("Rsp Update New Kline Time: %s", kline_data.TimeList())
		return fmt.Sprintf("[ws]:%s,[k]:%s %s", d.ws_info.String(), kline_data.Symbol, kline_data.TimeList())
	}
}

type SymbolSubInfo struct {
	mutex   sync.Mutex
	ws_info *treemap.Map
}

func (d *SymbolSubInfo) String() string {
	rst := "SymbolSubInfo: "

	iter := d.ws_info.Iterator()
	for iter.Begin(); iter.Next(); {
		rst = rst + fmt.Sprintf("ws: %d", iter.Key().(int64)) + ","
	}

	return rst
}

func NewSymbolSubInfo() *SymbolSubInfo {
	return &SymbolSubInfo{
		ws_info: treemap.NewWith(utils.Int64Comparator),
	}
}

type DepthSubInfo struct {
	mutex sync.Mutex
	Info  map[string]*treemap.Map
}

func (d *DepthSubInfo) String() string {
	rst := "DepthSubInfo: "
	for symbol, ws_map := range d.Info {
		rst = rst + symbol + ":"

		iter := ws_map.Iterator()
		for iter.Begin(); iter.Next(); {
			rst = rst + fmt.Sprintf("ws: %d", iter.Key().(int64)) + ","
		}
	}
	return rst
}

func NewDepthSubInfo() *DepthSubInfo {
	return &DepthSubInfo{
		Info: make(map[string]*treemap.Map),
	}
}

type KlineSubItem struct {
	ws_info    *treemap.Map
	cache_data *datastruct.Kline
}

func NewKlineWithKline() *KlineSubItem {
	return &KlineSubItem{
		ws_info: treemap.NewWith(utils.Int64Comparator),
	}
}

type KlineSubInfo struct {
	mutex sync.Mutex
	Info  map[string](map[int]*KlineSubItem)
}

func (d *KlineSubInfo) String() string {
	rst := "KlineSubInfo: "
	for symbol, resolution_map := range d.Info {
		rst = rst + symbol + ", "

		for resolution, ws_map := range resolution_map {
			rst = rst + fmt.Sprintf(" resolution: %d, ", resolution)
			iter := ws_map.ws_info.Iterator()
			for iter.Begin(); iter.Next(); {
				rst = rst + fmt.Sprintf("ws: %d, ", iter.Key().(int64))
			}
		}

	}
	return rst
}

func NewKlineSubInfo() *KlineSubInfo {
	return &KlineSubInfo{
		Info: make(map[string](map[int]*KlineSubItem)),
	}
}

type TradeSubInfo struct {
	mutex sync.Mutex
	Info  map[string]*treemap.Map
}

func (d *TradeSubInfo) String() string {
	rst := "TradeSubInfo: "
	for symbol, ws_map := range d.Info {
		rst = rst + symbol + ":"

		iter := ws_map.Iterator()
		for iter.Begin(); iter.Next(); {
			rst = rst + fmt.Sprintf("ws: %d, ", iter.Key().(int64))
		}
	}
	return rst
}

func NewTradeSubInfo() *TradeSubInfo {
	return &TradeSubInfo{
		Info: make(map[string]*treemap.Map),
	}
}
