package httptools

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

type (
	// WsSubscriber websocket 订阅器, 封装了断线重连和重新订阅的功能，且是并发安全的
	WsSubscriber interface {
		writeMsg(msg []byte) error
		LeftTopics() int // 剩余可订阅数量
		Sub(topics ...string) error
		UnSub(topics ...string) error
		HaveTopic(topic string) bool
		ResetConn() // 重置网络链接
		ReadCh() <-chan []byte
		Done() <-chan struct{}
		Close()
	}

	WsSubscriberConfig struct {
		MaxTopicCount     int32 // 每个conn 最多订阅的topic数量 默认 50
		Tag               string
		Proxy             string
		KeepaliveInterval time.Duration
	}

	TopicsToMsgFn func(topics ...string) ([][]byte, error) // 将订阅的topic转化成发给对端的msg

	wsSubscriber struct {
		*WsSubscriberConfig
		safeWsConn    SafeWsConn
		dialer        *websocket.Dialer
		rwMutex       sync.RWMutex
		m             sync.Map // 存储所有订阅的topic
		url           string
		curTopicCount int32 // 当前订阅的topic数量
		writeCh       chan []byte
		readCh        chan []byte
		closed        bool
		subFn         TopicsToMsgFn
		unsubFn       TopicsToMsgFn
		cancel        context.CancelFunc
		ctx           context.Context
	}
)

func NewWsSubscriber(url_ string, subFn, unsubFn TopicsToMsgFn,
	configFn func(config *WsSubscriberConfig)) (WsSubscriber, error) {
	var (
		dialer = &websocket.Dialer{}
		config = &WsSubscriberConfig{
			MaxTopicCount: 50,
		}
	)

	if configFn != nil {
		configFn(config)
	}

	if config.Proxy != "" {
		uProxy, err := url.Parse(config.Proxy)
		if err != nil {
			return nil, err
		}

		dialer.Proxy = http.ProxyURL(uProxy)
	}

	conn, _, err := dialer.Dial(url_, nil)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	ws := &wsSubscriber{
		WsSubscriberConfig: config,
		safeWsConn:         NewSafeWsConn(conn, config.Tag, config.KeepaliveInterval),
		dialer:             dialer,
		rwMutex:            sync.RWMutex{},
		m:                  sync.Map{},
		url:                url_,
		curTopicCount:      0,
		writeCh:            make(chan []byte, 256),
		readCh:             make(chan []byte, 256),
		closed:             false,
		subFn:              subFn,
		unsubFn:            unsubFn,
		cancel:             cancel,
		ctx:                ctx,
	}

	go ws.readAndWrite()
	return ws, nil
}

func (o *wsSubscriber) writeMsg(msg []byte) error {
	if !o.isClosed() {
		o.writeCh <- msg
		return nil
	} else {
		return fmt.Errorf("[wsSubscriber] tag:%s is closed ", o.Tag)
	}
}

func (o *wsSubscriber) LeftTopics() int {
	return int(atomic.LoadInt32(&o.MaxTopicCount) - atomic.LoadInt32(&o.curTopicCount))
}

func (o *wsSubscriber) HaveTopic(topic string) bool {
	_, ok := o.m.Load(topic)
	return ok
}

func (o *wsSubscriber) Sub(topics ...string) error {
	var unStoredTopics []string // 有效的topic

	for _, topic := range topics {
		if _, ok := o.m.Load(topic); !ok {
			unStoredTopics = append(unStoredTopics, topic)
			o.m.Store(topic, struct{}{})
		}
	}

	if len(unStoredTopics) == 0 {
		return nil
	}

	if o.LeftTopics() < len(unStoredTopics) {
		return fmt.Errorf("[wsSubscriber] tag:%s maximum topic limit exceeded. ", o.Tag)
	}

	atomic.AddInt32(&o.curTopicCount, int32(len(unStoredTopics)))

	msgs, err := o.subFn(unStoredTopics...)
	if err != nil {
		return err
	}

	if !o.isClosed() {
		//logx.Infof("[wsSubscriber] tag:%s send msg:%s", o.Tag, string(msg))
		for _, msg := range msgs {
			o.writeCh <- msg
		}
		return nil
	} else {
		return fmt.Errorf("[wsSubscriber] tag:%s is closed ", o.Tag)
	}
}

func (o *wsSubscriber) UnSub(topics ...string) error {
	var storedTopic []string // 有效的topic

	for _, topic := range topics {
		if _, ok := o.m.Load(topic); ok {
			storedTopic = append(storedTopic, topic)
		}
	}

	if len(storedTopic) == 0 {
		return nil
	}

	atomic.AddInt32(&o.curTopicCount, -int32(len(storedTopic)))

	msgs, err := o.unsubFn(storedTopic...)
	if err != nil {
		return err
	}

	for _, topic := range storedTopic {
		o.m.Delete(topic)
	}

	if !o.isClosed() {
		//logx.Infof("[wsSubscriber] tag:%s send msg:%s", o.Tag, string(msg))
		for _, msg := range msgs {
			o.writeCh <- msg
		}
		return nil
	} else {
		return fmt.Errorf("[wsSubscriber] tag:%s is closed ", o.Tag)
	}
}

func (o *wsSubscriber) ReadCh() <-chan []byte {
	return o.readCh
}

func (o *wsSubscriber) Close() {
	defer o.rwMutex.Unlock()
	o.rwMutex.Lock()
	o.cancel()
	o.closed = true
	o.safeWsConn.Close()
}

func (o *wsSubscriber) Done() <-chan struct{} {
	return o.ctx.Done()
}

func (o *wsSubscriber) ResetConn() {
	defer o.rwMutex.Unlock()
	o.rwMutex.Lock()
	logx.Infof("[wsSubscriber] tag:%s, start reset conn, this conn will be closed", o.Tag)
	o.safeWsConn.Close()
}

func (o *wsSubscriber) isClosed() bool {
	defer o.rwMutex.RUnlock()
	o.rwMutex.RLock()
	return o.closed
}

func (o *wsSubscriber) reconnect() {
	logx.Infof("[wsSubscriber] tag:%s, start reconnect ", o.Tag)
	for !o.isClosed() {
		conn, _, err := o.dialer.Dial(o.url, nil)
		if err != nil {
			logx.Errorf("[wsSubscriber] tag:%s, CreateConn err:%v", o.Tag, err)
			time.Sleep(time.Second * 3)
			continue
		}
		if !o.isClosed() {
			o.rwMutex.Lock()
			close(o.safeWsConn.WriteCh())
			o.safeWsConn = NewSafeWsConn(conn, o.Tag, time.Second*10)
			go o.readAndWrite()
			o.rwMutex.Unlock()
			o.reSub() //  reSub
		}
		break
	}

	logx.Infof("[wsSubscriber] tag:%s, reconnect completed ", o.Tag)
}

// reSub 重连时，重新订阅记录的 topic
func (o *wsSubscriber) reSub() {
	var topics []string
	o.m.Range(func(key, value interface{}) bool {
		topics = append(topics, key.(string))
		return true
	})

	if len(topics) == 0 {
		return
	}

	msgs, err := o.subFn(topics...)
	if err != nil {
		logx.Errorf("[wsSubscriber] tag:%s err:%v", o.Tag, err)
	}

	if !o.isClosed() {
		//logx.Infof("[wsSubscriber] tag:%s write msg:%s", o.Tag, string(msg))
		for _, msg := range msgs {
			o.writeCh <- msg
		}
	}
}

func (o *wsSubscriber) readAndWrite() {
	logx.Infof("[wsSubscriber] tag:%s start readAndWrite", o.Tag)
	readCh := o.safeWsConn.ReadCh()
	for {
		select {
		case <-o.ctx.Done():
			goto exit
		case <-o.safeWsConn.Done():
			if !o.isClosed() {
				go o.reconnect()
			}
			goto exit
		case msg, ok := <-o.writeCh:
			if !ok { // 被关闭了
				goto exit
			}
			o.safeWsConn.WriteCh() <- msg
		case msg, ok := <-readCh:
			if !ok {
				if !o.isClosed() {
					go o.reconnect()
				}
				goto exit
			}
			o.readCh <- msg
		}
	}
exit:
	logx.Infof("[wsSubscriber] tag:%s quit readAndWrite", o.Tag)
}
