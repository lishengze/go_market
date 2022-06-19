package net

import (
	"fmt"
	"market_server/common/util"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

var upgrader = websocket.Upgrader{}

type WSInfo struct {
	ID          int64
	Conn        *websocket.Conn
	IsAlive     bool
	LastReqTime int64
}

func NewWSInfo(conn *websocket.Conn) *WSInfo {
	return &WSInfo{
		ID:   util.UTCNanoTime(),
		Conn: conn,
	}
}

func (w *WSInfo) Close() {
	fmt.Printf("%d closed! \n", w.ID)
}

func (s *WSInfo) String() string {
	return fmt.Sprintf("ID: %d ", s.ID)
}

func (s *WSInfo) CheckAlive(HeartbeatLostSecs int64) bool {

	if util.UTCNanoTime()-s.LastReqTime > HeartbeatLostSecs {
		s.IsAlive = false
		return false
	} else {
		return true
	}
}

type WSEngine struct {
	WSConSet          map[int64]*WSInfo
	WSConSetMutex     sync.Mutex
	HeartbeatSendSecs int64
	HeartbeatLostSecs int64
}

func (w *WSEngine) Start() {
	w.StartListen()
	w.StartHeartbeat()
}

func (w *WSEngine) StoreWS(ws *WSInfo) {
	w.WSConSetMutex.Lock()
	defer w.WSConSetMutex.Unlock()

	w.WSConSet[ws.ID] = ws
}

func (w *WSEngine) StartListen() {
	http.HandleFunc("/trading/marketws", w.ListenRequest)
}

func (w *WSEngine) StartHeartbeat() {
	logx.Info("---- StatisticTimeTask Start!")
	duration := time.Duration((time.Duration)(w.HeartbeatSendSecs) * time.Second)
	timer := time.Tick(duration)

	for {
		select {
		case <-timer:
			w.ChecktHeartbeat()
		}
	}
}

func GetHeartbeatMsg() []byte {
	rst := "data"
	return []byte(rst)
}

func (w *WSEngine) ChecktHeartbeat() {
	w.WSConSetMutex.Lock()
	defer w.WSConSetMutex.Unlock()

	var dead_ws = []*WSInfo{}

	for _, ws := range w.WSConSet {
		if !ws.CheckAlive(w.HeartbeatLostSecs) {
			dead_ws = append(dead_ws, ws)
		}
	}

	for _, ws := range dead_ws {
		delete(w.WSConSet, ws.ID)
		//
	}

	for _, ws := range w.WSConSet {
		ws.Conn.WriteMessage(1, GetHeartbeatMsg())
	}
}

func (w *WSEngine) ListenRequest(h http.ResponseWriter, r *http.Request) {
	fmt.Printf("%+v, Echo: %+v \n", time.Now(), *r)

	c, err := upgrader.Upgrade(h, r, nil)

	ws := NewWSInfo(c)

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

func (w *WSEngine) ProcessMessage(msg []byte, ws *WSInfo) {

}

func (w *WSEngine) ProcessSubDepth(msg []byte, ws *WSInfo) {

}

func (w *WSEngine) ProcessSubTrade(msg []byte, ws *WSInfo) {

}

func (w *WSEngine) ProcessSubKline(msg []byte, ws *WSInfo) {

}

func (w *WSEngine) ProcessHeartbeat(msg []byte, ws *WSInfo) {

}
