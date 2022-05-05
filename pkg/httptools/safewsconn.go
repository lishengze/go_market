package httptools


import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
	"time"
)

type (
	SafeWsConn interface {
		Done() <-chan struct{}
		WriteCh() chan<- []byte
		ReadCh() <-chan []byte
		Close()
	}

	safeWsConn struct {
		mu                sync.RWMutex
		conn              *websocket.Conn
		keepaliveInterval time.Duration
		tag               string
		readCh            chan []byte
		writeCh           chan []byte
		writeWait         time.Duration
		cancel            context.CancelFunc
		ctx               context.Context
		lastPongTime      time.Time
	}
)

func NewSafeWsConn(conn *websocket.Conn, tag string, keepaliveInterval time.Duration) SafeWsConn {
	ctx, cancel := context.WithCancel(context.Background())

	swc := &safeWsConn{
		conn:              conn,
		keepaliveInterval: keepaliveInterval,
		tag:               tag,
		readCh:            make(chan []byte, 256),
		writeCh:           make(chan []byte, 256),
		writeWait:         time.Second * 10,
		cancel:            cancel,
		ctx:               ctx,
		lastPongTime:      time.Now(),
	}
	go swc.read()
	if keepaliveInterval > 0 {
		go swc.writeWithKeepalive()
	} else {
		go swc.write()
	}
	return swc
}

func (o *safeWsConn) Done() <-chan struct{} {
	return o.ctx.Done()
}

func (o *safeWsConn) WriteCh() chan<- []byte {
	return o.writeCh
}

func (o *safeWsConn) ReadCh() <-chan []byte {
	return o.readCh
}

func (o *safeWsConn) Close() {
	o.cancel()
	_ = o.conn.Close()
}

// read 是一个接收消息的 go routine
func (o *safeWsConn) read() {
	logx.Infof("[wsconn] tag:%s,start read go routine", o.tag)

	defer func() {
		o.cancel()
		_ = o.conn.Close() // 关闭conn
		close(o.readCh)
	}()

	for {
		t, message, err := o.conn.ReadMessage()
		if err != nil {
			logx.Errorf("[wsconn] tag:%s,ReadMessage t:%d msg:%s err:%v", o.tag, t, string(message), err)
			break // 退出
		}
		//println(string(message))
		o.readCh <- message
	}

	logx.Infof("[wsconn] tag:%s,quit read go routine", o.tag)
}

// writeWithKeepalive 负责写的go routine, 并且会 定时主动发送 ping 帧
func (o *safeWsConn) writeWithKeepalive() {
	logx.Infof("[wsconn] tag:%s, start writeWithKeepalive go routine", o.tag)
	o.conn.SetPongHandler(o.pongHandler)
	ticker := time.NewTicker(o.keepaliveInterval)
	defer func() {
		o.cancel()
		ticker.Stop()
		_ = o.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-o.writeCh:
			_ = o.conn.SetWriteDeadline(time.Now().Add(o.writeWait))
			if !ok {
				_ = o.conn.WriteMessage(websocket.CloseMessage, []byte{})
				goto exit
			}
			logx.Infof("[wsconn] tag:%s, write msg:%s ", o.tag, string(msg))
			err := o.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				logx.Errorf("[wsconn] tag:%s, write msg:%s ,err:%v", o.tag, string(msg), err)
				goto exit
			}
		case <-ticker.C:
			if time.Since(o.getLastPongTime()) > 3*o.keepaliveInterval {
				logx.Errorf("[wsconn] tag:%s, %v not rcv pong ", o.tag, time.Since(o.lastPongTime))
				goto exit
			}

			_ = o.conn.SetWriteDeadline(time.Now().Add(o.writeWait))
			err := o.conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second*10))
			if err != nil {
				logx.Errorf("[wsconn] tag:%s, send keepalive ping err:%v", o.tag, err)
				goto exit
			}
		}
	}
exit:
	logx.Infof("[wsconn] tag:%s,quit writeWithKeepalive go routine", o.tag)
}

// write 负责写的go routine
func (o *safeWsConn) write() {
	logx.Infof("[wsconn] tag:%s, start write go routine", o.tag)
	defer func() {
		o.cancel()
		_ = o.conn.SetReadDeadline(time.Now())
		//_ = o.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-o.writeCh:
			_ = o.conn.SetWriteDeadline(time.Now().Add(o.writeWait))
			if !ok {
				_ = o.conn.WriteMessage(websocket.CloseMessage, []byte{})
				goto exit
			}
			logx.Infof("[wsconn] tag:%s, write msg:%s ", o.tag, string(msg))
			err := o.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				logx.Errorf("[wsconn] tag:%s, write msg:%s ,err:%v", o.tag, string(msg), err)
				goto exit
			}
		}
	}
exit:
	logx.Infof("[wsconn] tag:%s,quit write go routine", o.tag)
}

func (o *safeWsConn) pongHandler(msg string) error {
	o.setLastPongTime(time.Now())
	return nil
}

func (o *safeWsConn) getLastPongTime() time.Time {
	defer o.mu.RUnlock()
	o.mu.RLock()
	return o.lastPongTime
}

func (o *safeWsConn) setLastPongTime(t time.Time) {
	defer o.mu.Unlock()
	o.mu.Lock()
	o.lastPongTime = t
}

