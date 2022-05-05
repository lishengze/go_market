package httptools

import (
	"github.com/zeromicro/go-zero/core/logx"
)

type (
	// WsReceiver websocket 接收器, 继承于 WsSubscriber
	WsReceiver interface {
		ReadCh() <-chan interface{}
		ResetConn() // 重置网络链接
		Done() <-chan struct{}
		Close()
	}

	wsReceiver struct {
		WsSubscriber
		readCh       chan interface{}
		msgConvertFn func(msg []byte) (interface{}, error)
	}
)

//func NewWsReceiver(url_, tag string, keepaliveInterval time.Duration, proxy string,
//	msgConvertFn func(msg []byte) (interface{}, error)) (WsReceiver, error) {
//	wsSubscriber, err := NewWsSubscriber(url_, tag, proxy, 0, keepaliveInterval,
//		func(topics ...string) ([]byte, error) {
//			return []byte{}, nil
//		}, func(topics ...string) ([]byte, error) {
//			return []byte{}, nil
//		})
//	if err != nil {
//		return nil, err
//	}
//	wr := &wsReceiver{
//		WsSubscriber: wsSubscriber,
//		readCh:       make(chan interface{}, 256),
//		msgConvertFn: msgConvertFn,
//	}
//	go wr.read()
//	return wr, nil
//}

func (o *wsReceiver) ReadCh() <-chan interface{} {
	return o.readCh
}

func (o *wsReceiver) read() {
	for {
		select {
		case msg, ok := <-o.WsSubscriber.ReadCh():
			if !ok {
				return
			}
			res, err := o.msgConvertFn(msg)
			if err != nil {
				logx.Error(err)
			} else {
				o.readCh <- res
			}
		case <-o.WsSubscriber.Done():
			return
		}
	}
}

func (o *wsReceiver) Close() {
	o.WsSubscriber.Close()
	close(o.readCh)
}
