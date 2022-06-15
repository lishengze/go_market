package httptools

import (
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type (
	// WsTransceiver websocket 收发器, 是并发安全的
	WsTransceiver interface {
		WriteMsg(msg []byte) error
		ReadCh() <-chan interface{}
		ResetConn() // 重置网络链接
		Done() <-chan struct{}
		Close()
	}

	WsTransceiverConfig struct {
		Tag               string
		Proxy             string
		KeepaliveInterval time.Duration
	}

	wsTransceiver struct {
		WsSubscriber
		readCh       chan interface{}
		msgConvertFn func(msg []byte) (interface{}, error)
	}

	MsgConvertFn func(msg []byte) (interface{}, error)
)

func NewWsTransceiver(url_ string, msgConvertFn MsgConvertFn, configFn func(config *WsTransceiverConfig)) (WsTransceiver, error) {

	c := &WsTransceiverConfig{}
	if configFn != nil {
		configFn(c)
	}

	wsSubscriber, err := NewWsSubscriber(url_,
		func(topics ...string) ([][]byte, error) {
			return [][]byte{}, nil
		}, func(topics ...string) ([][]byte, error) {
			return [][]byte{}, nil
		}, func(config *WsSubscriberConfig) {
			config.Tag = c.Tag
			config.Proxy = c.Proxy
			config.KeepaliveInterval = c.KeepaliveInterval
		})

	if err != nil {
		return nil, err
	}
	wt := &wsTransceiver{
		WsSubscriber: wsSubscriber,
		readCh:       make(chan interface{}, 256),
		msgConvertFn: msgConvertFn,
	}
	go wt.read()
	return wt, nil
}

func (o *wsTransceiver) WriteMsg(msg []byte) error {
	return o.WsSubscriber.writeMsg(msg)
}

func (o *wsTransceiver) ReadCh() <-chan interface{} {
	return o.readCh
}

func (o *wsTransceiver) read() {
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

func (o *wsTransceiver) Close() {
	o.WsSubscriber.Close()
	close(o.readCh)
}
