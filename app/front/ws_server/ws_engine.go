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
	"strconv"
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
	go w.StartListen()
	go w.StartHeartbeat()
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

func (w *WSEngine) CloseWS(ws *net.WSInfo) {
	ws.Close()
	w.RemoveWS(ws)
}

func (w *WSEngine) RemoveWS(ws *net.WSInfo) {
	w.WSConSetMutex.Lock()
	defer w.WSConSetMutex.Unlock()

	if _, ok := w.WSConSet[ws.ID]; ok {
		delete(w.WSConSet, ws.ID)
	}
}

func (w *WSEngine) StartListen() {
	logx.Infof("Start Listen: %s:%s", w.WsConfig.Address, w.WsConfig.Url)

	http.HandleFunc(w.WsConfig.Url, w.ListenRequest)
	http.ListenAndServe(w.WsConfig.Address, nil)
}

func (w *WSEngine) StartHeartbeat() {
	logx.Infof("---- StartHeartbeat Start, HeartbeatSendSecs: %d", w.WsConfig.HeartbeatSendSecs)
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

	logx.Statf("ChecktHeartbeat, WSConSet.Size: %d", len(w.WSConSet))

	var dead_ws = []*net.WSInfo{}

	for _, ws := range w.WSConSet {
		if !ws.IsAlive() || !ws.CheckAlive(int64(w.WsConfig.HeartbeatLostSecs)) {
			logx.Errorf("ws: %s is dead! last_time: %+v, cur_time: %+v,  HeartbeatLostSecs: %d\n",
				ws.String(), util.GetTimeFromtInt(ws.LastReqTime), time.Now(), w.WsConfig.HeartbeatLostSecs)
			dead_ws = append(dead_ws, ws)
		}
	}

	for _, ws := range dead_ws {
		delete(w.WSConSet, ws.ID)
	}

	for _, ws := range w.WSConSet {
		logx.Statf("Pub Heartbeat To %s", ws.String())
		ws.Conn.WriteMessage(1, GetHeartbeatMsg())
	}
}

func (w *WSEngine) Close(ws *net.WSInfo) {
	ws.SetAlive(false)
	ws.Close()
}

func (w *WSEngine) ListenRequest(h http.ResponseWriter, r *http.Request) {
	logx.Infof("RequestInfo: %+v, Echo: %+v \n", time.Now(), *r)

	c, err := upgrader.Upgrade(h, r, nil)

	ws := net.NewWSInfo(c)

	w.StoreWS(ws)

	if err != nil {
		logx.Errorf("upgrade err: %+v ", err)
		logx.Infof("upgrade err: %+v ", err)
		w.CloseWS(ws)
		return
	}
	defer w.CloseWS(ws)

	w.ProcessSubSymbol(ws)

	for {
		mt, message, err := ws.Conn.ReadMessage()
		if err != nil {
			w.CloseWS(ws)
			logx.Errorf("read, type: %d, err: %+v\n", mt, err)
			logx.Infof("read, type: %d, err: %+v\n", mt, err)
		} else {
			w.ProcessMessage(message, ws)
		}
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

   if (js["type"].get<string>() == KLINE_SUB)
   {
       process_kline_req(ori_msg, socket_id, ws_safe);
   }

   if (js["type"].get<string>() == TRADE)
   {
       process_trade_req(ori_msg, socket_id, ws_safe);
   }
*/

func catch_exp(msg []byte, ws *net.WSInfo) {
	errMsg := recover()
	if errMsg != nil {
		fmt.Println("This is catch_exp func")
		logx.Errorf("catch_exp OriginalMsg: %s, WSInfo: %+v\n", msg, *ws)
		logx.Errorf("errMsg: %+v \n", errMsg)
		fmt.Println(errMsg)
	}

}

func (w *WSEngine) ProcessMessage(msg []byte, ws *net.WSInfo) {
	defer catch_exp(msg, ws)

	logx.Infof("Original Msg: %s \n", string(msg))
	ws.SetLastReqTime(util.UTCNanoTime())

	if len(string(msg)) < 8 {
		logx.Errorf("Unknown Msg %s ", string(msg))
		return
	}

	var m map[string]interface{}
	if err := json.Unmarshal([]byte(msg), &m); err != nil {
		logx.Errorf("Error = %+v", err)
		return
	}

	logx.Infof("msg: %s, mapping: %+v\n", string(msg), m)

	if _, ok := m["type"]; !ok {
		logx.Error("Msg Error, ori msg: %+v", string(msg))
		return
	}

	if m["type"].(string) == net.SYMBOL_SUB {
		w.ProcessSubDepth(m, ws)
		w.ProcessSubTrade(m, ws)
	}

	if m["type"].(string) == net.DEPTH_SUB {
		w.ProcessSubDepth(m, ws)
	}

	if m["type"].(string) == net.DEPTH_UNSUB {
		w.ProcessUnSubDepth(m, ws)
	}

	if m["type"].(string) == net.TRADE_SUB {
		w.ProcessSubTrade(m, ws)
	}

	if m["type"].(string) == net.TRADE_UNSUB {
		w.ProcessUnSubTrade(m, ws)
	}

	if m["type"].(string) == net.KLINE_SUB {
		w.ProcessSubKline(m, ws)
	}

	if m["type"].(string) == net.KLINE_UNSUMB {
		w.ProcessUnSubKline(m, ws)
	}

	if m["type"].(string) == net.HEARTBEAT {
		w.ProcessHeartbeat(m, ws)
	}
}

func (w *WSEngine) ProcessSubSymbol(ws *net.WSInfo) {
	logx.Infof("WS %+v, SubStart!", ws)
	w.next_worker.SubSymbol(ws)
}

/*
   sub_info = {
       "type":"sub_symbol",
       "symbol":[symbol]
   }
*/
func (w *WSEngine) ProcessSubDepth(m map[string]interface{}, ws *net.WSInfo) {
	if value, ok := m["symbol"]; ok {
		symbol_list := value.([]interface{})

		for _, symbol := range symbol_list {
			w.next_worker.SubDepth(symbol.(string), ws)
		}

	} else {
		logx.Error("ProcessSubTrade: No Symbol Data %+v", m)
	}
}

/*
   unsub_info = {
       "type":"sub_symbol",
       "symbol":[symbol]
   }
*/
func (w *WSEngine) ProcessUnSubDepth(m map[string]interface{}, ws *net.WSInfo) {
	if value, ok := m["symbol"]; ok {
		symbol_list := value.([]interface{})

		for _, symbol := range symbol_list {
			w.next_worker.UnSubDepth(symbol.(string), ws)
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
	logx.Infof("SubTradeInfo: %+v", m)
	if value, ok := m["symbol"]; ok {
		// logx.Infof("value: %+v", value)
		symbol_list := value.([]interface{})

		for _, symbol := range symbol_list {
			w.next_worker.SubTrade(symbol.(string), ws)
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
func (w *WSEngine) ProcessUnSubTrade(m map[string]interface{}, ws *net.WSInfo) {
	logx.Infof("UnSubTradeInfo: %+v", m)
	if value, ok := m["symbol"]; ok {
		// logx.Infof("value: %+v", value)
		symbol_list := value.([]interface{})

		for _, symbol := range symbol_list {
			w.next_worker.UnSubTrade(symbol.(string), ws)
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
		// value_type := reflect.TypeOf(value)
		// true_value := reflect.ValueOf(value)
		tmp, err := strconv.Atoi(value.(string))
		if err != nil {
			ws.SendErrorMsg(fmt.Sprintf("subkline frequency is not string: %+v ", m))
			logx.Errorf("subkline frequency is not string: %+v ", m)
			return
		}

		resolution = uint32(tmp)

		if resolution%datastruct.SECS_PER_MIN != 0 {
			ws.SendErrorMsg(fmt.Sprintf("resolution : %d mod 60 != 0  ", resolution))
			logx.Errorf("resolution : %d mod 60 != 0  ", resolution)
			return
		}
	} else {
		ws.SendErrorMsg(fmt.Sprintf("sub_kline: No frequency Data %+v", m))
		logx.Error("ProcessSubTrade: No frequency Data %+v", m)
		return
	}

	if value, ok := m["count"]; ok {
		tmp, err := strconv.Atoi(value.(string))
		if err != nil {
			ws.SendErrorMsg(fmt.Sprintf("subkline count is not string: %+v ", m))
			logx.Errorf("subkline count is not string: %+v ", m)
			return
		}
		count = uint32(tmp)
	} else {
		ws.SendErrorMsg(fmt.Sprintf("subkline: No count Data %+v", m))
		logx.Errorf("subkline: No count Data %+v", m)
	}

	if value, ok := m["start_time"]; ok {
		tmp, err := strconv.Atoi(value.(string))
		if err != nil {
			logx.Errorf("subkline start_time is not string: %+v ", m)
			return
		}
		start_time = uint64(tmp)
	}
	if value, ok := m["end_time"]; ok {
		tmp, err := strconv.Atoi(value.(string))
		if err != nil {
			logx.Errorf("subkline end_time is not string: %+v ", m)
			return
		}

		end_time = uint64(tmp)
	}

	if uint64(count)+start_time+end_time == 0 {
		logx.Errorf("subkline: No data_count start_time end_time Data %+v", m)
		ws.SendErrorMsg(fmt.Sprintf("subkline: No data_count start_time end_time Data %+v", m))
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

func (w *WSEngine) ProcessUnSubKline(m map[string]interface{}, ws *net.WSInfo) {
	var symbol string
	var resolution uint32

	if value, ok := m["symbol"]; ok {
		symbol = value.(string)
	} else {
		logx.Error("ProcessSubTrade: No Symbol Data %+v", m)
		return
	}

	if value, ok := m["frequency"]; ok {
		tmp, err := strconv.Atoi(value.(string))
		if err != nil {
			logx.Errorf("subkline frequency is not string: %+v ", m)
			return
		}
		resolution = uint32(tmp)
	} else {
		logx.Error("ProcessSubTrade: No frequency Data %+v", m)
		return
	}

	req_kline := &datastruct.ReqHistKline{
		Symbol:    symbol,
		Exchange:  datastruct.BCTS_EXCHANGE,
		Frequency: resolution,
	}
	w.next_worker.UnSubKline(req_kline, ws)
}

/*

 */
func (w *WSEngine) ProcessHeartbeat(m map[string]interface{}, ws *net.WSInfo) {

}
