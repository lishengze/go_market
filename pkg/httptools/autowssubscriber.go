package httptools

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
)

type (
	// AutoWsSubscriber 是一个可以自动扩容的 WsSubscriber, 是并发安全的
	AutoWsSubscriber interface {
		Sub(topics ...string)
		Unsub(topics ...string)
		ReadCh() <-chan interface{}
		ResetConn(topics ...string) // 重置订阅这些 topic 的网络链接
		Done() <-chan struct{}
		Close()
	}

	autoWsSubscriber struct {
		m            sync.Map // 存储所有订阅的topic
		mutex        sync.Mutex
		subscribers  []WsSubscriber
		new          func() (WsSubscriber, error)
		msgConvertFn func(msg []byte) (interface{}, error)
		readCh       chan interface{}
		cancel       context.CancelFunc
		ctx          context.Context
	}
)

/*
NewAutoWsSubscriber
msgConvertFn: 将接收到的 bytes 转换成特定格式的函数
*/
func NewAutoWsSubscriber(newFn func() (WsSubscriber, error), msgConvertFn MsgConvertFn) (AutoWsSubscriber, error) {

	ctx, cancel := context.WithCancel(context.Background())

	return &autoWsSubscriber{
		mutex:        sync.Mutex{},
		subscribers:  []WsSubscriber{},
		new:          newFn,
		msgConvertFn: msgConvertFn,
		readCh:       make(chan interface{}, 256),
		cancel:       cancel,
		ctx:          ctx,
	}, nil
}

func (o *autoWsSubscriber) Sub(topics ...string) {
	defer o.mutex.Unlock()
	o.mutex.Lock()

	var unStoredTopics []string // 有效的topic

	for _, topic := range topics {
		if _, ok := o.m.Load(topic); !ok {
			unStoredTopics = append(unStoredTopics, topic)
			o.m.Store(topic, struct{}{})
		}
	}

	if len(unStoredTopics) == 0 {
		return
	}

	for _, subscriber := range o.subscribers {
		if subscriber.LeftTopics() >= len(unStoredTopics) {
			// subscriber 可以完全订阅 所有的 topics
			err := subscriber.Sub(unStoredTopics...)
			if err != nil {
				logx.Errorf("err:%v", err)
			}
			return
		} else {
			// subscriber 可以订阅 一部分 topics
			var left = subscriber.LeftTopics()
			err := subscriber.Sub(unStoredTopics[0:left]...)
			if err != nil {
				logx.Errorf("err:%v", err)
			}

			unStoredTopics = append(unStoredTopics[:0], unStoredTopics[left:]...)
			// 继续找下个 subscriber
		}
	}

	// 还有未订阅的topic 循环创建 subscriber
	for len(unStoredTopics) > 0 {
		subscriber, err := o.new()
		if err != nil {
			logx.Errorf("create new subscriber err:%v", err)
			return
		}
		go o.read(subscriber)

		if subscriber.LeftTopics() >= len(unStoredTopics) {
			// subscriber 可以完全订阅 所有的 topics
			err := subscriber.Sub(unStoredTopics...)
			if err != nil {
				logx.Errorf("err:%v", err)
			}

			o.subscribers = append(o.subscribers, subscriber)
			return
		} else {
			// subscriber 可以订阅 一部分 topics
			var left = subscriber.LeftTopics()
			err := subscriber.Sub(unStoredTopics[0:left]...)
			if err != nil {
				logx.Errorf("err:%v", err)
			}

			unStoredTopics = append(unStoredTopics[:0], unStoredTopics[left:]...)
		}
		o.subscribers = append(o.subscribers, subscriber)
	}

}

func (o *autoWsSubscriber) Unsub(topics ...string) {
	defer o.mutex.Unlock()
	o.mutex.Lock()
	for _, subscriber := range o.subscribers {
		for _, topic := range topics {
			if subscriber.HaveTopic(topic) {
				err := subscriber.UnSub(topics...)
				if err != nil {
					logx.Errorf("err:%v", err)
				}
			}
		}
	}

	for _, topic := range topics {
		o.m.Delete(topic)
	}
}

func (o *autoWsSubscriber) ResetConn(topics ...string) {
	defer o.mutex.Unlock()
	o.mutex.Lock()

	var subscribers = make(map[WsSubscriber]struct{}, 0) // 记录需要重置的 WsSubscriber
	for _, subscriber := range o.subscribers {
		for _, topic := range topics {
			if subscriber.HaveTopic(topic) {
				subscribers[subscriber] = struct{}{}
			}
		}
	}

	for subscriber := range subscribers {
		go func(s WsSubscriber) {
			s.ResetConn()
		}(subscriber)
	}
}

func (o *autoWsSubscriber) Close() {
	defer o.mutex.Unlock()
	o.mutex.Lock()

	for _, subscriber := range o.subscribers {
		go func(s WsSubscriber) {
			s.Close()
		}(subscriber)
	}
	close(o.readCh)
	o.cancel()
}

func (o *autoWsSubscriber) Done() <-chan struct{} {
	return o.ctx.Done()
}

func (o *autoWsSubscriber) ReadCh() <-chan interface{} {
	return o.readCh
}

func (o *autoWsSubscriber) read(subscriber WsSubscriber) {
	ch := subscriber.ReadCh()
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				return
			}
			res, err := o.msgConvertFn(msg)
			if err != nil {
				logx.Errorf("msgConvertFn err:%v, msg:%s", err, string(msg))
			} else {
				o.readCh <- res
			}
		case <-subscriber.Done():
			return
		case <-o.ctx.Done():
			return
		}
	}
}
