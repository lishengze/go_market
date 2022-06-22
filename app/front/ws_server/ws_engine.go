package ws_server

import (
	"encoding/json"
	"fmt"
	"market_server/app/front/config"
	"market_server/app/front/net"
	"market_server/app/front/worker"
	"market_server/common/datastruct"
	"market_server/common/util"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

var upgrader = websocket.Upgrader{}

type WSEngine struct {
	WSConSet      map[int64]*net.WSInfo
	WSConSetMutex sync.Mutex

	WsConfig *config.WSConfig

	next_worker worker.WorkerI

	IsTest bool
}

func NewWSEngine(ws_config *config.WSConfig) *WSEngine {
	return &WSEngine{
		WsConfig: ws_config,
		WSConSet: make(map[int64]*net.WSInfo),
	}
}

func (w *WSEngine) Start() {
	w.StartListen()
	w.StartHeartbeat()
}

func (w *WSEngine) SetTestFlag(value bool) {
	w.IsTest = value
}

func (w *WSEngine) SetNextWorker(next_worker worker.WorkerI) {
	w.next_worker = next_worker
}

func (w *WSEngine) StoreWS(ws *net.WSInfo) {
	w.WSConSetMutex.Lock()
	defer w.WSConSetMutex.Unlock()

	w.WSConSet[ws.ID] = ws
}

func (w *WSEngine) StartListen() {
	logx.Infof("Start Listen: %s:%s", w.WsConfig.Address, w.WsConfig.Url)
	http.ListenAndServe(w.WsConfig.Address, nil)
	http.HandleFunc(w.WsConfig.Url, w.ListenRequest)
}

func (w *WSEngine) StartHeartbeat() {
	logx.Info("---- StatisticTimeTask Start!")
	duration := time.Duration((time.Duration)(w.WsConfig.HeartbeatSendSecs) * time.Second)
	timer := time.Tick(duration)

	for {
		select {
		case <-timer:
			w.ChecktHeartbeat()
		}
	}
}

// {"time":"2022-06-20 08:23:54.20117945","type":"heartbeat"}

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

func (w *WSEngine) ChecktHeartbeat() {
	w.WSConSetMutex.Lock()
	defer w.WSConSetMutex.Unlock()

	var dead_ws = []*net.WSInfo{}

	for _, ws := range w.WSConSet {
		if !ws.CheckAlive(int64(w.WsConfig.HeartbeatLostSecs)) {
			dead_ws = append(dead_ws, ws)
		}
	}

	for _, ws := range dead_ws {
		delete(w.WSConSet, ws.ID)
	}

	for _, ws := range w.WSConSet {
		ws.Conn.WriteMessage(1, GetHeartbeatMsg())
	}
}

func (w *WSEngine) ListenRequest(h http.ResponseWriter, r *http.Request) {
	fmt.Printf("%+v, Echo: %+v \n", time.Now(), *r)

	c, err := upgrader.Upgrade(h, r, nil)

	ws := net.NewWSInfo(c)

	w.StoreWS(ws)

	if err != nil {
		logx.Errorf("upgrade err: %+v ", err)
		return
	}
	defer ws.Close()

	for {
		mt, message, err := ws.Conn.ReadMessage()
		if err != nil {
			logx.Errorf("read, type: %d, err: %+v\n", mt, err)
			break
		}
		w.ProcessMessage(message, ws)
	}
}

/*
   if (js["type"].get<string>() == "sub_symbol")
   {
       process_depth_req(ori_msg, socket_id, ws_safe);

       process_trade_req(ori_msg, socket_id, ws_safe);
   }

   if (js["type"].get<string>() == HEARTBEAT)
   {
       process_heartbeat(socket_id, ws_safe);
   }

   if (js["type"].get<string>() == KLINE_UPDATE)
   {
       process_kline_req(ori_msg, socket_id, ws_safe);
   }

   if (js["type"].get<string>() == TRADE)
   {
       process_trade_req(ori_msg, socket_id, ws_safe);
   }
*/

func (w *WSEngine) ProcessMessage(msg []byte, ws *net.WSInfo) {
	ws.SetLastReqTime(util.UTCNanoTime())

	var m map[string]interface{}
	if err := json.Unmarshal([]byte(msg), &m); err != nil {
		logx.Errorf("Error = %+v", err)
		return
	}

	if m["type"].(string) == net.SYMBOL_SUB {
		w.ProcessSubDepth(m, ws)
		w.ProcessSubTrade(m, ws)
	}

	if m["type"].(string) == net.KLINE_UPDATE {
		w.ProcessSubKline(m, ws)
	}

	if m["type"].(string) == net.TRADE {
		w.ProcessSubTrade(m, ws)
	}

	if m["type"].(string) == net.HEARTBEAT {
		w.ProcessHeartbeat(m, ws)
	}
}

/*
   sub_info = {
       "type":"sub_symbol",
       "symbol":[symbol]
   }
*/
func (w *WSEngine) ProcessSubDepth(m map[string]interface{}, ws *net.WSInfo) {
	if value, ok := m["symbol"]; ok {
		symbol_list := value.([]string)

		for _, symbol := range symbol_list {
			w.next_worker.SubDepth(symbol, ws)
		}

	} else {
		logx.Error("ProcessSubTrade: No Symbol Data %+v", m)
	}
}

/*
   sub_info = {
       "type":"trade",
       "symbol":[symbol]
   }
*/
func (w *WSEngine) ProcessSubTrade(m map[string]interface{}, ws *net.WSInfo) {
	if value, ok := m["symbol"]; ok {
		symbol_list := value.([]string)

		for _, symbol := range symbol_list {
			w.next_worker.SubTrade(symbol, ws)
		}
	} else {
		logx.Error("ProcessSubTrade: No Symbol Data %+v", m)
	}
}

/*
   sub_info = {
       "type":"kline_update",
       "symbol":symbol,
       "data_count":str(2),
       "frequency":str(frequency)
   }
*/
func (w *WSEngine) ProcessSubKline(m map[string]interface{}, ws *net.WSInfo) {
	var symbol string
	var resolution uint32
	count := uint32(0)
	start_time := uint64(0)
	end_time := uint64(0)

	if value, ok := m["symbol"]; ok {
		symbol = value.(string)
	} else {
		logx.Error("ProcessSubTrade: No Symbol Data %+v", m)
		return
	}

	if value, ok := m["frequency"]; ok {
		resolution = value.(uint32)
	} else {
		logx.Error("ProcessSubTrade: No frequency Data %+v", m)
		return
	}

	if value, ok := m["data_count"]; ok {
		count = value.(uint32)
	} else {
		logx.Error("ProcessSubTrade: No data_count Data %+v", m)
	}

	if value, ok := m["start_time"]; ok {
		start_time = value.(uint64)
	}
	if value, ok := m["end_time"]; ok {
		end_time = value.(uint64)
	}

	if uint64(count)+start_time+end_time == 0 {
		logx.Error("ProcessSubTrade: No data_count start_time end_time Data %+v", m)
		return
	}

	req_kline := &datastruct.ReqHistKline{
		Symbol:    symbol,
		Exchange:  datastruct.BCTS_EXCHANGE,
		StartTime: start_time,
		EndTime:   end_time,
		Count:     count,
		Frequency: resolution,
	}
	w.next_worker.SubKline(req_kline, ws)
}

/*

 */
func (w *WSEngine) ProcessHeartbeat(m map[string]interface{}, ws *net.WSInfo) {

}
