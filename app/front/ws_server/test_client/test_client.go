package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"market_server/app/front/front_engine"
	"market_server/app/front/net"
	"market_server/common/datastruct"
	"market_server/common/util"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

func InitLogx() {

	LogConfig := logx.LogConf{
		Compress:            true,
		KeepDays:            0,
		Level:               "info",
		Mode:                "file",
		Path:                "./log",
		ServiceName:         "client",
		StackCooldownMillis: 100,
		TimeFormat:          "2006-01-02 15:04:05",
	}

	logx.MustSetup(LogConfig)
}

type TestMain struct {
	TradeUpdatedSymbolMap      map[string]struct{}
	TradeUpdatedSymbolMapMutex sync.Mutex
}

func NewTestMain() *TestMain {
	return &TestMain{
		TradeUpdatedSymbolMap: make(map[string]struct{}),
	}
}

func (t *TestMain) GetTestTradeReqJson() []byte {
	symbol_list := []string{"BTC_USDT", "ETH_USDT", "USDT_USD", "BTC_USD", "ETH_USD", "ETH_BTC"}
	// symbol_list := []string{"BTC_USDT", "ETH_USDT"}
	req_start_time := strconv.FormatInt(util.UTCNanoTime(), 10)
	sub_info := map[string]interface{}{
		"type":           net.TRADE_SUB,
		"symbol":         symbol_list,
		"req_start_time": req_start_time,
	}
	rst, err := json.Marshal(sub_info)

	if err != nil {
		logx.Errorf("GetTestTradeReqJson: %+v \n", err)
		return nil
	} else {
		logx.Infof("SubJson: %s", string(rst))
		return rst
	}
}

func (t *TestMain) GetTestDepthReqJson() []byte {
	symbol_list := []string{"BTC_USDT"}
	sub_info := map[string]interface{}{
		"type":   net.SYMBOL_SUB,
		"symbol": symbol_list,
	}
	rst, err := json.Marshal(sub_info)

	if err != nil {
		logx.Errorf("GetTestDepthReqJson: %+v \n", err)
		return nil
	} else {
		logx.Infof("SubJson: %s", string(rst))
		return rst
	}
}

func (t *TestMain) GetTestKlineReqJson(frequency int) []byte {
	sub_info := map[string]interface{}{
		"type":      net.KLINE_SUB,
		"symbol":    "BTC_USDT",
		"count":     "20",
		"frequency": fmt.Sprintf("%d", frequency),
	}
	rst, err := json.Marshal(sub_info)
	if err != nil {
		logx.Errorf("GetTestKlineReqJson: %+v \n", err)
		return nil
	} else {
		logx.Infof("SubJson: %s", string(rst))
		return rst
	}
}

func (t *TestMain) TestGetJsonData() {
	rst1 := t.GetTestTradeReqJson()
	rst2 := t.GetTestDepthReqJson()
	rst3 := t.GetTestKlineReqJson(600)

	fmt.Println(string(rst1))
	fmt.Println(string(rst2))
	fmt.Println(string(rst3))
}

// var addr = flag.String("addr", "127.0.0.1:8114", "http service address")

var addr = flag.String("addr", "18.162.42.238:8114", "http service address")

// var addr = flag.String("addr", "10.10.1.75:8114", "http service address")

func (t *TestMain) GetHeartbeat() []byte {
	info := map[string]interface{}{
		"type": "heartbeat",
	}

	rst, err := json.Marshal(info)

	if err != nil {
		logx.Errorf("GetTestDepthReqJson: %+v \n", err)
		return nil
	} else {
		return rst
	}
}

func (t *TestMain) read_func(c *websocket.Conn) {
	logx.Info("Read_Func Start!")
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			logx.Info("read:", err)
			return
		}
		// log.Printf("recv: %s", message)

		var m map[string]interface{}
		if err := json.Unmarshal([]byte(message), &m); err != nil {
			logx.Errorf("Error = %+v", err)
			return
		}

		if _, ok := m["type"]; !ok {
			logx.Error("Msg Error, ori msg: %+v", string(message))
			return
		}

		if m["type"] == "heartbeat" {
			err := c.WriteMessage(websocket.TextMessage, t.GetHeartbeatMsg())
			if err != nil {
				logx.Info("write:", err)
				return
			}
		}

		if m["type"] == net.KLINE_UPATE {
			t.process_kline(message)
		}

		if m["type"] == net.TRADE_UPATE {
			t.process_trade(message)
		}

	}
}

func (t *TestMain) process_kline(message []byte) {
	var kline_data front_engine.PubKlineJson
	if err := json.Unmarshal([]byte(message), &kline_data); err != nil {
		logx.Errorf("Error = %+v", err)
		return
	} else {
		delta_time := util.UTCNanoTime() - kline_data.ReqResponseTime
		logx.Infof("Kline: req_process_time: %d us, ws_time: %dus, kline_data: %s", kline_data.ReqProcessTime/datastruct.NANO_PER_MICR, delta_time/datastruct.NANO_PER_MICR, kline_data.TimeList())
	}
}

func (t *TestMain) process_trade(message []byte) {
	var trade_data front_engine.PubTradeJson
	if err := json.Unmarshal([]byte(message), &trade_data); err != nil {
		logx.Errorf("Error = %+v", err)
		return
	} else {
		t.TradeUpdatedSymbolMapMutex.Lock()
		delta_time := util.UTCNanoTime() - trade_data.ReqResponseTime
		if _, ok := t.TradeUpdatedSymbolMap[trade_data.Symbol]; !ok {
			logx.Infof("Trade %s, req_process_time: %d us, ws_time: %dus ", trade_data.Symbol, trade_data.ReqProcessTime/datastruct.NANO_PER_MICR, delta_time/datastruct.NANO_PER_MICR)
			fmt.Printf("Trade %s, req_ws:%d us, req_process_time: %d us, ws_time: %dus \n", trade_data.Symbol, trade_data.ReqWSTime/datastruct.NANO_PER_MICR, trade_data.ReqProcessTime/datastruct.NANO_PER_MICR, delta_time/datastruct.NANO_PER_MICR)
			t.TradeUpdatedSymbolMap[trade_data.Symbol] = struct{}{}
		}

		t.TradeUpdatedSymbolMapMutex.Unlock()
	}
}

func (t *TestMain) GetHeartbeatMsg() []byte {

	heartbeat_map := map[string]interface{}{
		"time": util.UTCNanoTime(),
		"type": net.HEARTBEAT,
	}
	rst, err := json.Marshal(heartbeat_map)

	if err != nil {
		logx.Errorf("GetHeartbeatMsg: %+v", err)
		return nil
	}
	return rst
}

func (t *TestMain) write_func(c *websocket.Conn) {

	send_msg := t.GetTestTradeReqJson()
	// send_msg := GetTestDepthReqJson()
	// send_msg := GetTestKlineReqJson(900)

	err := c.WriteMessage(websocket.TextMessage, send_msg)
	if err != nil {
		logx.Info("write:", err)
		return
	}

	// time.Sleep(time.Second * 5)

	// send_msg2 := GetTestKlineReqJson(300)

	// err = c.WriteMessage(websocket.TextMessage, send_msg2)
	// if err != nil {
	// 	logx.Info("write:", err)
	// 	return
	// }

	// ticker := time.NewTicker(time.Second)
	// defer ticker.Stop()

	// for {
	// 	select {
	// 	case t := <-ticker.C:
	// 		fmt.Printf("Write: %s \n", t.String())
	// 		err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
	// 		if err != nil {
	// 			logx.Info("write:", err)
	// 			return
	// 		}
	// 	}
	// }
}

func basic_func() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/trading/marketws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	t := NewTestMain()

	// done := make(chan struct{})

	go t.read_func(c)

	go t.write_func(c)

	for {
		select {

		case <-interrupt:
			logx.Info("interrupt")
			return
		}
	}
}

func test_split() {
	symbol := "BTC_USDT"
	symbol_list := strings.Split(symbol, "_")

	fmt.Println(symbol_list)
}

func main() {
	InitLogx()

	basic_func()

	// test_split()
}
