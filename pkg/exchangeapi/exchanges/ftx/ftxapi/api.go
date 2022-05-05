package ftxapi

import (
	"encoding/json"
	"exterior-interactor/pkg/exchangeapi/exmodel"
	"exterior-interactor/pkg/exchangeapi/extools"
	"exterior-interactor/pkg/httptools"
	"net/http"
	"time"
)

type Api struct {
	ApiType exmodel.ApiType
	base    extools.ExBase
}

func NewApi(base extools.ExBase) *Api {
	return &Api{
		ApiType: exmodel.ApiTypeUnified,
		base:    base,
	}
}

// GetMarket 获取市场品种
func (o *Api) GetMarket() (*GetMarketRsp, error) {
	var (
		meta = extools.NewMetaWithOneWeight(http.MethodGet, urlMarket, extools.ReqTypeWait, false)
		res  = &GetMarketRsp{}
	)

	err := o.base.Request(httptools.EmptyReq, &res, meta)
	return res, err
}

// GetStreamMarketTrade 获取 逐笔成交
func (o *Api) GetStreamMarketTrade() (httptools.AutoWsSubscriber, error) {
	newSubscriberFn := func() (httptools.WsSubscriber, error) {
		return httptools.NewWsSubscriber(wsUrl,
			func(topics ...string) ([][]byte, error) {
				var msg [][]byte
				for _, topic := range topics {
					m, err := json.Marshal(WsSubscribeMsg{
						Op:      "subscribe",
						Channel: "trades",
						Market:  topic,
					})
					if err != nil {
						return nil, err
					}
					msg = append(msg, m)
				}
				return msg, nil
			}, func(topics ...string) ([][]byte, error) {
				var msg [][]byte
				for _, topic := range topics {
					m, err := json.Marshal(WsUnsubscribeMsg{
						Op:      "unsubscribe",
						Channel: "trades",
						Market:  topic,
					})
					if err != nil {
						return nil, err
					}
					msg = append(msg, m)
				}
				return msg, nil
			}, func(config *httptools.WsSubscriberConfig) {
				config.Tag = "ftx-wstrade"
				config.Proxy = o.base.Config().Proxy
				config.MaxTopicCount = 50
				config.KeepaliveInterval = time.Second * 15
			})
	}

	return httptools.NewAutoWsSubscriber(newSubscriberFn, func(msg []byte) (interface{}, error) {
		data := &StreamMarketTrade{}
		err := json.Unmarshal(msg, &data)
		return data, err
	})
}

// GetStreamDepth 获取 depth 推送
func (o *Api) GetStreamDepth() (httptools.AutoWsSubscriber, error) {
	newSubscriberFn := func() (httptools.WsSubscriber, error) {
		return httptools.NewWsSubscriber(wsUrl,
			func(topics ...string) ([][]byte, error) {
				var msg [][]byte
				for _, topic := range topics {
					m, err := json.Marshal(WsSubscribeMsg{
						Op:      "subscribe",
						Channel: "orderbook",
						Market:  topic,
					})
					if err != nil {
						return nil, err
					}
					msg = append(msg, m)
				}
				return msg, nil
			}, func(topics ...string) ([][]byte, error) {
				var msg [][]byte
				for _, topic := range topics {
					m, err := json.Marshal(WsUnsubscribeMsg{
						Op:      "unsubscribe",
						Channel: "orderbook",
						Market:  topic,
					})
					if err != nil {
						return nil, err
					}
					msg = append(msg, m)
				}
				return msg, nil
			}, func(config *httptools.WsSubscriberConfig) {
				config.Tag = "ftx-wsdepth"
				config.Proxy = o.base.Config().Proxy
				config.MaxTopicCount = 50
				config.KeepaliveInterval = time.Second * 10
			})
	}

	return httptools.NewAutoWsSubscriber(newSubscriberFn, func(msg []byte) (interface{}, error) {
		data := &StreamDepth{}
		err := json.Unmarshal(msg, &data)
		return data, err
	})
}
