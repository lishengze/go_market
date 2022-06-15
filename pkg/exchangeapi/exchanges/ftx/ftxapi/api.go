package ftxapi

import (
	"encoding/json"
	"exterior-interactor/pkg/exchangeapi/exmodel"
	"exterior-interactor/pkg/exchangeapi/extools"
	"exterior-interactor/pkg/httptools"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
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

// GetBalance 获取余额
func (o *Api) GetBalance() (*GetBalanceRsp, error) {
	var (
		meta = extools.NewMetaWithOneWeight(http.MethodGet, urlBalances, extools.ReqTypeWait, true)
		res  = &GetBalanceRsp{}
	)

	err := o.base.Request(httptools.EmptyReq, &res, meta)
	return res, err
}

// PlaceOrder 下单
func (o *Api) PlaceOrder(req PlaceOrderReq) (*PlaceOrderRsp, error) {
	var (
		meta = extools.NewMetaWithOneWeight(http.MethodPost, urlOrders, extools.ReqTypeAllow, true)
		res  = &PlaceOrderRsp{}
	)

	err := o.base.Request(req, &res, meta)
	return res, err
}

// QueryOrder 查单
func (o *Api) QueryOrder(orderId string) (*QueryOrderRsp, error) {
	var (
		meta = extools.NewMetaWithOneWeight(http.MethodGet,
			fmt.Sprintf("%s/%s", urlOrders, orderId), extools.ReqTypeWait, true)
		res = &QueryOrderRsp{}
	)

	err := o.base.Request(httptools.EmptyReq, &res, meta)
	return res, err
}

// CancelOrder 撤单
func (o *Api) CancelOrder(orderId string) (*CancelOrderRsp, error) {
	var (
		meta = extools.NewMetaWithOneWeight(http.MethodDelete,
			fmt.Sprintf("%s/%s", urlOrders, orderId), extools.ReqTypeWait, true)
		res = &CancelOrderRsp{}
	)

	err := o.base.Request(httptools.EmptyReq, &res, meta)
	return res, err
}

// CancelOrderByClientOrderId 根据客户订单id撤单
func (o *Api) CancelOrderByClientOrderId(clientOrderId string) (*CancelOrderRsp, error) {
	var (
		meta = extools.NewMetaWithOneWeight(http.MethodDelete,
			fmt.Sprintf("%s/by_client_id/%s", urlOrders, clientOrderId), extools.ReqTypeWait, true)
		res = &CancelOrderRsp{}
	)

	err := o.base.Request(httptools.EmptyReq, &res, meta)
	return res, err
}

// QueryTrades 查询成交
func (o *Api) QueryTrades(req QueryTradesReq) (*QueryTradesRsp, error) {
	var (
		meta = extools.NewMetaWithOneWeight(http.MethodGet, urlFills, extools.ReqTypeWait, true)
		res  = &QueryTradesRsp{}
	)

	err := o.base.Request(req, &res, meta)
	return res, err
}

// GetTradeWsTransceiver 获取 成交推送器
func (o *Api) GetTradeWsTransceiver() (httptools.WsTransceiver, error) {
	wt, err := httptools.NewWsTransceiver(wsUrl, func(msg []byte) (interface{}, error) {
		logx.Infof("[ORIGINAL TRADE UPDATE]: %s", string(msg))
		var data = &WsFills{}
		err := json.Unmarshal(msg, &data)
		return data, err
	}, func(config *httptools.WsTransceiverConfig) {
		config.Tag = "ftx-wsfills"
		config.Proxy = o.base.Config().Proxy
		config.KeepaliveInterval = time.Second * 15
	})

	if err != nil {
		return nil, err
	}

	now := time.Now().UnixMilli()
	signStr := fmt.Sprintf("%dwebsocket_login", now)
	signature, err := extools.GetParamHmacSHA256Sign(o.base.Config().AccountConfig.Secret, signStr)
	if err != nil {
		wt.Close()
		return nil, err
	}

	loginReq := WsLoginReq{
		WsLoginArgs: WsLoginArgs{
			Key:        o.base.Config().AccountConfig.Key,
			Sign:       signature,
			Time:       now,
			SubAccount: o.base.Config().AccountConfig.SubAccountName,
		},
		Op: "login",
	}

	msg, _ := json.Marshal(loginReq)
	err = wt.WriteMsg(msg)
	if err != nil {
		wt.Close()
		return nil, err
	}

	subReq := WsSubscribeMsg{
		Op:      "subscribe",
		Channel: "fills",
	}
	subMsg, _ := json.Marshal(subReq)
	err = wt.WriteMsg(subMsg)
	if err != nil {
		wt.Close()
		return nil, err
	}

	return wt, err
}

// GetOrderWsTransceiver 获取 订单推送器
func (o *Api) GetOrderWsTransceiver() (httptools.WsTransceiver, error) {
	wt, err := httptools.NewWsTransceiver(wsUrl, func(msg []byte) (interface{}, error) {
		logx.Infof("[ORIGINAL ORDER UPDATE]: %s", string(msg))
		var data = &WsOrders{}
		err := json.Unmarshal(msg, &data)
		return data, err
	}, func(config *httptools.WsTransceiverConfig) {
		config.Tag = "ftx-wsorders"
		config.Proxy = o.base.Config().Proxy
		config.KeepaliveInterval = time.Second * 15
	})

	if err != nil {
		return nil, err
	}

	now := time.Now().UnixMilli()
	signStr := fmt.Sprintf("%dwebsocket_login", now)
	signature, err := extools.GetParamHmacSHA256Sign(o.base.Config().AccountConfig.Secret, signStr)
	if err != nil {
		wt.Close()
		return nil, err
	}

	loginReq := WsLoginReq{
		WsLoginArgs: WsLoginArgs{
			Key:        o.base.Config().AccountConfig.Key,
			Sign:       signature,
			Time:       now,
			SubAccount: o.base.Config().AccountConfig.SubAccountName,
		},
		Op: "login",
	}

	msg, _ := json.Marshal(loginReq)
	err = wt.WriteMsg(msg)
	if err != nil {
		wt.Close()
		return nil, err
	}

	subReq := WsSubscribeMsg{
		Op:      "subscribe",
		Channel: "orders",
	}
	subMsg, _ := json.Marshal(subReq)
	err = wt.WriteMsg(subMsg)
	if err != nil {
		wt.Close()
		return nil, err
	}

	return wt, err
}
