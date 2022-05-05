package binancespot

import (
	"encoding/json"
	"exterior-interactor/pkg/exchangeapi/exmodel"
	"exterior-interactor/pkg/exchangeapi/extools"
	"exterior-interactor/pkg/httptools"
	"math/rand"
	"net/http"
	"time"
)

type Api struct {
	ApiType exmodel.ApiType
	base    extools.ExBase
}

func NewApi(base extools.ExBase) *Api {
	return &Api{
		ApiType: exmodel.ApiTypeSpot,
		base: base,
	}
}

func (o *Api) GetExchangeInfo() (*GetExchangeInfoRsp, error) {
	var (
		// GET /api/v3/exchangeInfo 权重: 10
		meta = extools.NewMeta(http.MethodGet, urlExchangeInfo, 10, extools.ReqTypeWait, false)
		res  = &GetExchangeInfoRsp{}
	)

	err := o.base.Request(httptools.EmptyReq, &res, meta)
	return res, err
}

func (o *Api) GetDepth(req GetDepthReq) (*GetDepthRsp, error) {
	var (
		// GET /api/v3/depth
		/*
			权重基于限制调整:
			限制	权重
			1-100	1
			101-500	5
			501-1000	10
			1001-5000	50
		*/
		meta = extools.NewMetaWithGetWeightFn(http.MethodGet, urlDepth, extools.ReqTypeWait, false, func() int {
			if req.Limit >= 0 && req.Limit <= 100 {
				return 1
			}
			if req.Limit >= 101 && req.Limit <= 500 {
				return 5
			}
			if req.Limit >= 500 && req.Limit <= 1000 {
				return 10
			}
			if req.Limit >= 1001 && req.Limit <= 5000 {
				return 50
			}
			return 1
		})
		res = &GetDepthRsp{}
	)

	err := o.base.Request(req, &res, meta)
	return res, err
}

// GetStreamDiffDepth 获取增量 depth 推送
func (o *Api) GetStreamDiffDepth() (httptools.AutoWsSubscriber, error) {
	newSubscriberFn := func() (httptools.WsSubscriber, error) {
		return httptools.NewWsSubscriber(wsUrlCombinedStream,
			func(topics ...string) ([][]byte, error) {
				m, err := json.Marshal(WsMarketStreamSendMsg{
					Method: "SUBSCRIBE",
					Params: topics,
					Id:     rand.Int31n(1000),
				})
				return [][]byte{m}, err
			}, func(topics ...string) ([][]byte, error) {
				m, err := json.Marshal(WsMarketStreamSendMsg{
					Method: "UNSUBSCRIBE",
					Params: topics,
					Id:     rand.Int31n(1000),
				})
				return [][]byte{m}, err
			}, func(config *httptools.WsSubscriberConfig) {
				config.Tag = "binance-spot-diff-wsdepth"
				config.Proxy = o.base.Config().Proxy
				config.MaxTopicCount = 50
				config.KeepaliveInterval = time.Second * 10
			})
	}

	return httptools.NewAutoWsSubscriber(newSubscriberFn, func(msg []byte) (interface{}, error) {
		data := &WsDiffDepth{}
		err := json.Unmarshal(msg, &data)
		return data, err
	})
}

// GetStreamKline 获取 kline stream
func (o *Api) GetStreamKline() (httptools.AutoWsSubscriber, error) {
	newSubscriberFn := func() (httptools.WsSubscriber, error) {
		return httptools.NewWsSubscriber(wsUrlCombinedStream,
			func(topics ...string) ([][]byte, error) {
				m, err := json.Marshal(WsMarketStreamSendMsg{
					Method: "SUBSCRIBE",
					Params: topics,
					Id:     rand.Int31n(1000),
				})
				return [][]byte{m}, err
			}, func(topics ...string) ([][]byte, error) {
				m, err := json.Marshal(WsMarketStreamSendMsg{
					Method: "UNSUBSCRIBE",
					Params: topics,
					Id:     rand.Int31n(1000),
				})
				return [][]byte{m}, err
			}, func(config *httptools.WsSubscriberConfig) {
				config.Tag = "binance-spot-kline"
				config.Proxy = o.base.Config().Proxy
				config.MaxTopicCount = 50
				config.KeepaliveInterval = time.Second * 10
			})
	}

	return httptools.NewAutoWsSubscriber(newSubscriberFn, func(msg []byte) (interface{}, error) {
		data := &WsKline{}
		err := json.Unmarshal(msg, &data)
		return data, err
	})
}

// GetStreamMarketTrade 获取 逐笔成交
func (o *Api) GetStreamMarketTrade() (httptools.AutoWsSubscriber, error) {
	newSubscriberFn := func() (httptools.WsSubscriber, error) {
		return httptools.NewWsSubscriber(wsUrlCombinedStream,
			func(topics ...string) ([][]byte, error) {
				m, err := json.Marshal(WsMarketStreamSendMsg{
					Method: "SUBSCRIBE",
					Params: topics,
					Id:     rand.Int31n(1000),
				})
				return [][]byte{m}, err
			}, func(topics ...string) ([][]byte, error) {
				m, err := json.Marshal(WsMarketStreamSendMsg{
					Method: "UNSUBSCRIBE",
					Params: topics,
					Id:     rand.Int31n(1000),
				})
				return [][]byte{m}, err
			}, func(config *httptools.WsSubscriberConfig) {
				config.Tag = "binance-spot-wstrade"
				config.Proxy = o.base.Config().Proxy
				config.MaxTopicCount = 50
				config.KeepaliveInterval = time.Second * 10
			})
	}

	return httptools.NewAutoWsSubscriber(newSubscriberFn, func(msg []byte) (interface{}, error) {
		data := &StreamMarketTrade{}
		err := json.Unmarshal(msg, &data)
		return data, err
	})
}
