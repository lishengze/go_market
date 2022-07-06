package monitor_market

import (
	"encoding/json"
	"log"
	"market_server/app/front/net"
	"market_server/common/util"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

type WSClient struct {
	Config     *WSConfig
	Client     *websocket.Conn
	SymbolList []string
}

func NewWSClient(config *WSConfig, symbol_list []string) *WSClient {
	return &WSClient{
		Config:     config,
		Client:     nil,
		SymbolList: symbol_list,
	}
}

func (w *WSClient) Start() {
	err := w.InitClient()

	if err != nil {
		return
	}

	go w.StartListenData()

	go w.StartSubData()
}

func (w *WSClient) InitClient() error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: w.Config.Address, Path: w.Config.Url}
	log.Printf("connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		logx.Errorf("websocket.DefaultDialer err: %+v \n", err)
		return err
	}

	w.Client = conn
	return nil
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

func (w *WSClient) StartListenData() {
	for {
		_, message, err := w.Client.ReadMessage()
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
			logx.Errorf("Msg Error, ori msg: %+v", string(message))
			return
		}

		if m["type"] == net.HEARTBEAT {
			w.ProcessHeartbeat()
		} else if m["type"] == net.DEPTH_UPDATE {

		} else if m["type"] == net.TRADE_UPATE {

		} else if m["type"] == net.KLINE_UPATE {

		}
	}
}

func (w *WSClient) ProcessHeartbeat() {
	err := w.Client.WriteMessage(websocket.TextMessage, GetHeartbeatMsg())
	if err != nil {
		logx.Errorf("WriteMessage:", err)
		return
	}
}

func GetTestTradeReqJson(symbol_list []string) []byte {
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

func GetTestDepthReqJson(symbol_list []string) []byte {
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

func GetTestKlineReqJson(symbol string) []byte {
	sub_info := map[string]interface{}{
		"type":      net.KLINE_SUB,
		"symbol":    symbol,
		"count":     "2",
		"frequency": "600",
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

func (w *WSClient) StartSubDepth() {
	send_msg := GetTestDepthReqJson(w.SymbolList)

	err := w.Client.WriteMessage(websocket.TextMessage, send_msg)
	if err != nil {
		logx.Errorf("StartSubDepth err: %+v, send_msg:%s", err, string(send_msg))
		return
	}
}

func (w *WSClient) StartSubTrade() {
	send_msg := GetTestTradeReqJson(w.SymbolList)

	err := w.Client.WriteMessage(websocket.TextMessage, send_msg)
	if err != nil {
		logx.Errorf("StartSubTrade err: %+v, send_msg:%s", err, string(send_msg))
		return
	}
}

func (w *WSClient) StartSubKline() {
	for _, symbol := range w.SymbolList {
		send_msg := GetTestKlineReqJson(symbol)

		err := w.Client.WriteMessage(websocket.TextMessage, send_msg)
		if err != nil {
			logx.Errorf("StartSubKline err: %+v, send_msg:%s", err, string(send_msg))
			return
		}
	}
}

func (w *WSClient) StartSubData() {
	w.StartSubDepth()
	w.StartSubTrade()
	w.StartSubKline()
}

func (w *WSClient) ProcessDepth(m map[string]interface{}) {

}

func (w *WSClient) ProcessTrade(m map[string]interface{}) {

}

func (w *WSClient) ProcessKline(m map[string]interface{}) {

}
