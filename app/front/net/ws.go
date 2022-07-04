package net

import (
	"encoding/json"
	"fmt"
	"market_server/common/datastruct"
	"market_server/common/util"
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

type WSInfo struct {
	ID          int64
	Conn        *websocket.Conn
	Alive       int32
	LastReqTime int64

	mutex sync.Mutex
}

func NewWSInfo(conn *websocket.Conn) *WSInfo {
	return &WSInfo{
		ID:          util.UTCNanoTime(),
		Conn:        conn,
		Alive:       1,
		LastReqTime: util.UTCNanoTime(),
	}
}

type ErrorMsg struct {
	TypeInfo string `json:"type"`
	Info     string `json:"info"`
}

func (w *WSInfo) SendErrorMsg(msg string) {
	json_data := ErrorMsg{
		TypeInfo: DEPTH_UPDATE,
		Info:     msg,
	}

	rst, err := json.Marshal(json_data)

	if err != nil {
		w.SendMsg(1, rst)
	}
}

func (w *WSInfo) SendMsg(messageType int, data []byte) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.Conn.WriteMessage(messageType, data)
}

func (w *WSInfo) Close() {
	logx.Infof("%d closed! \n", w.ID)
	atomic.StoreInt32(&w.Alive, 0)
}

func (w *WSInfo) String() string {
	return fmt.Sprintf("ID: %d ", w.ID)
}

func (w *WSInfo) SetAlive(value bool) {
	if value {
		atomic.StoreInt32(&w.Alive, 1)
	} else {
		atomic.StoreInt32(&w.Alive, 0)
	}

}

func (w *WSInfo) CheckAlive(HeartbeatLostSecs int64) bool {

	if util.UTCNanoTime()-atomic.LoadInt64(&w.LastReqTime) > HeartbeatLostSecs*datastruct.NANO_PER_SECS {
		atomic.StoreInt32(&w.Alive, 1)
		return false
	} else {
		return true
	}
}

func (w *WSInfo) IsAlive() bool {
	if atomic.LoadInt32(&w.Alive) == 1 {
		return true
	} else {
		return false
	}
}

func (w *WSInfo) SetLastReqTime(req_time int64) {
	atomic.StoreInt64(&w.LastReqTime, req_time)
	atomic.StoreInt32(&w.Alive, 1)
}
