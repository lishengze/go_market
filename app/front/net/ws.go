package net

import (
	"fmt"
	"market_server/common/datastruct"
	"market_server/common/util"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

type WSInfo struct {
	ID          int64
	Conn        *websocket.Conn
	Alive       int32
	LastReqTime int64
}

func NewWSInfo(conn *websocket.Conn) *WSInfo {
	return &WSInfo{
		ID:    util.UTCNanoTime(),
		Conn:  conn,
		Alive: 1,
	}
}

func (w *WSInfo) Close() {
	logx.Infof("%d closed! \n", w.ID)
	atomic.StoreInt32(&w.Alive, 0)
}

func (w *WSInfo) String() string {
	return fmt.Sprintf("ID: %d ", w.ID)
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
