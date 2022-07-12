package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"market_server/app/front/net"
	"market_server/common/util"
	"net/url"
	"os"
	"os/signal"
	"time"

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

func GetTestTradeReqJson() []byte {
	symbol_list := []string{"BTC_USDT"}
	sub_info := map[string]interface{}{
		"type":   net.TRADE_SUB,
		"symbol": symbol_list,
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

func GetTestDepthReqJson() []byte {
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

func GetTestKlineReqJson(frequency int) []byte {
	sub_info := map[string]interface{}{
		"type":      net.KLINE_SUB,
		"symbol":    "ETH_USDT",
		"count":     "1000",
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

func TestGetJsonData() {
	rst1 := GetTestTradeReqJson()
	rst2 := GetTestDepthReqJson()
	rst3 := GetTestKlineReqJson(600)

	fmt.Println(string(rst1))
	fmt.Println(string(rst2))
	fmt.Println(string(rst3))
}

// var addr = flag.String("addr", "127.0.0.1:8114", "http service address")

var addr = flag.String("addr", "18.162.42.238:8114", "http service address")

// var addr = flag.String("addr", "10.10.1.75:8114", "http service address")

func GetHeartbeat() []byte {
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

func read_func(c *websocket.Conn) {
	log.Println("Read_Func Start!")
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s", message)

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
			err := c.WriteMessage(websocket.TextMessage, GetHeartbeatMsg())
			if err != nil {
				log.Println("write:", err)
				return
			}

		}
	}
}

func GetHeartbeatMsg() []byte {

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

func write_func(c *websocket.Conn) {

	// send_msg := GetTestTradeReqJson()
	// send_msg := GetTestDepthReqJson()
	send_msg := GetTestKlineReqJson(60)

	err := c.WriteMessage(websocket.TextMessage, send_msg)
	if err != nil {
		log.Println("write:", err)
		return
	}

	time.Sleep(time.Second * 5)

	send_msg2 := GetTestKlineReqJson(300)

	err = c.WriteMessage(websocket.TextMessage, send_msg2)
	if err != nil {
		log.Println("write:", err)
		return
	}

	// ticker := time.NewTicker(time.Second)
	// defer ticker.Stop()

	// for {
	// 	select {
	// 	case t := <-ticker.C:
	// 		fmt.Printf("Write: %s \n", t.String())
	// 		err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
	// 		if err != nil {
	// 			log.Println("write:", err)
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

	// done := make(chan struct{})

	go read_func(c)

	go write_func(c)

	for {
		select {

		case <-interrupt:
			log.Println("interrupt")
			return
		}
	}
}

func main() {
	InitLogx()

	basic_func()

	// TestGetJsonData()
}
