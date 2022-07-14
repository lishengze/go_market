package monitor_market

import (
	"encoding/json"
	"market_server/app/front/net"
	"market_server/common/datastruct"
	"market_server/common/monitorStruct"
	"market_server/common/util"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

type WSClient struct {
	Config      *WSConfig
	Client      *websocket.Conn
	SymbolList  []string
	MonitorChan chan *monitorStruct.MonitorData

	statistic_secs     int
	rcv_statistic_info sync.Map
	statistic_start    time.Time
}

func NewWSClient(config *WSConfig, symbol_list []string, monitor_channel chan *monitorStruct.MonitorData) *WSClient {
	return &WSClient{
		Config:         config,
		Client:         nil,
		SymbolList:     symbol_list,
		MonitorChan:    monitor_channel,
		statistic_secs: 22,
	}
}

func (w *WSClient) Start() {
	logx.Info("------- WSClient Start -------")
	err := w.InitClient()

	if err != nil {
		return
	}

	go w.StartListenData()

	go w.StartSubData()

	go w.StatisticTimeTaskMain()
}

func (k *WSClient) StatisticTimeTaskMain() {
	logx.Info("---- WSClient StatisticTimeTask Start!")
	duration := time.Duration((time.Duration)(k.statistic_secs) * time.Second)
	timer := time.Tick(duration)

	k.statistic_start = time.Now()
	for {
		select {
		case <-timer:
			k.UpdateStatisticInfo()
		}
	}
}

func (k *WSClient) OutputRcvInfo(key, value interface{}) bool {
	if value.(int) != 0 {
		logx.Statf("[rcv] %s : %d ", key, value)
		k.rcv_statistic_info.Store(key, 0)
	}

	return true
}

func (k *WSClient) UpdateStatisticInfo() {

	logx.Statf("Websocket Statistic Start: %+v \n", util.TimeToSecString(k.statistic_start))

	k.rcv_statistic_info.Range(k.OutputRcvInfo)

	k.statistic_start = time.Now()

	logx.Statf("Websocket Statistic End: %+v \n", util.TimeToSecString(k.statistic_start))
}

func (w *WSClient) InitClient() error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: w.Config.Address, Path: w.Config.Url}
	logx.Infof("connecting to %s", u.String())

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
			logx.Infof("WS Read Err: %+v", err)
			logx.Errorf("WS Read Err: %+v", err)
			return
		}
		logx.Infof("WSClient Msg: %s", message)

		go w.ProcessMsg(message)
	}
}

func (w *WSClient) ProcessMsg(message []byte) {
	defer util.CatchExp("WSClient::ProcessMsg")

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
		w.ProcessDepth(m)
	} else if m["type"] == net.TRADE_UPATE {
		w.ProcessTrade(m)
	} else if m["type"] == net.KLINE_UPATE {
		w.ProcessKline(m)
	}
}

func (w *WSClient) UpdateRecvInfo(msg string) {
	if value, ok := w.rcv_statistic_info.Load(msg); ok {
		w.rcv_statistic_info.Store(msg, value.(int)+1)
	} else {
		w.rcv_statistic_info.Store(msg, 1)
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
		// logx.Infof("SubJson: %s", string(rst))
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
		// logx.Infof("SubJson: %s", string(rst))
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
		// logx.Infof("SubJson: %s", string(rst))
		return rst
	}
}

func (w *WSClient) StartSubDepth() {
	send_msg := GetTestDepthReqJson(w.SymbolList)
	logx.Infof("WS SubInfo %s ", string(send_msg))

	err := w.Client.WriteMessage(websocket.TextMessage, send_msg)
	if err != nil {
		logx.Errorf("StartSubDepth err: %+v, send_msg:%s", err, string(send_msg))
		return
	}
}

func (w *WSClient) StartSubTrade() {
	send_msg := GetTestTradeReqJson(w.SymbolList)
	logx.Infof("WS SubInfo %s ", string(send_msg))

	err := w.Client.WriteMessage(websocket.TextMessage, send_msg)
	if err != nil {
		logx.Errorf("StartSubTrade err: %+v, send_msg:%s", err, string(send_msg))
		return
	}
}

func (w *WSClient) StartSubKline() {
	for _, symbol := range w.SymbolList {
		send_msg := GetTestKlineReqJson(symbol)
		logx.Infof("WS SubInfo %s ", string(send_msg))

		err := w.Client.WriteMessage(websocket.TextMessage, send_msg)
		if err != nil {
			logx.Errorf("StartSubKline err: %+v, send_msg:%s", err, string(send_msg))
			return
		}
	}
}

func (w *WSClient) StartSubData() {
	// w.StartSubDepth()
	w.StartSubTrade()
	w.StartSubKline()
}

func (w *WSClient) ProcessDepth(m map[string]interface{}) {
	defer util.CatchExp("WSClient::ProcessDepth")

	if value, ok := m["symbol"]; ok {
		symbol := value.(string)

		msg := "depth." + datastruct.BCTS_EXCHANGE + "_" + symbol
		w.UpdateRecvInfo(msg)

		// logx.Slowf("WS depth: %s", msg)

		w.MonitorChan <- &monitorStruct.MonitorData{
			Symbol:   msg,
			DataType: datastruct.DEPTH_TYPE,
		}
	}
}

func (w *WSClient) ProcessTrade(m map[string]interface{}) {
	defer util.CatchExp("WSClient::ProcessTrade")

	if value, ok := m["symbol"]; ok {
		symbol := value.(string)

		msg := "trade." + datastruct.BCTS_EXCHANGE + "_" + symbol
		w.UpdateRecvInfo(msg)

		// logx.Slowf("WS trade: %s", msg)

		w.MonitorChan <- &monitorStruct.MonitorData{
			Symbol:   msg,
			DataType: datastruct.TRADE_TYPE,
		}
	}
}

func (w *WSClient) ProcessKline(m map[string]interface{}) {
	defer util.CatchExp("WSClient::ProcessKline")

	if value, ok := m["symbol"]; ok {
		symbol := value.(string)

		msg := "kline." + datastruct.BCTS_EXCHANGE + "_" + symbol
		w.UpdateRecvInfo(msg)

		// logx.Slowf("WS kline: %s", msg)

		w.MonitorChan <- &monitorStruct.MonitorData{
			Symbol:   msg,
			DataType: datastruct.KLINE_TYPE,
		}
	}
}
