package ftx

import (
	"context"
	"exterior-interactor/pkg/exchangeapi/exchanges/ftx/ftxapi"
	"exterior-interactor/pkg/exchangeapi/exmodel"
	"exterior-interactor/pkg/exchangeapi/extools"
	"exterior-interactor/pkg/httptools"
	"exterior-interactor/pkg/xmath"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
	"time"
)

type tradeManager struct {
	api                *NativeApi
	outputUpdateCh     chan *exmodel.OrderTradesUpdate
	orderWsTransceiver httptools.WsTransceiver
	tradeWsTransceiver httptools.WsTransceiver
	cancel             context.CancelFunc
	ctx                context.Context
}

func NewTradeManager(api *NativeApi) extools.TradeManager {
	orderWsTransceiver, err := api.GetOrderWsTransceiver()
	if err != nil {
		panic(err)
	}

	tradeWsTransceiver, err := api.GetTradeWsTransceiver()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	t := &tradeManager{
		api:                api,
		outputUpdateCh:     make(chan *exmodel.OrderTradesUpdate, 1024),
		orderWsTransceiver: orderWsTransceiver,
		tradeWsTransceiver: tradeWsTransceiver,
		cancel:             cancel,
		ctx:                ctx,
	}

	go t.run()

	return t
}

func (o *tradeManager) run() {
	orderCh := o.orderWsTransceiver.ReadCh()
	tradeCh := o.tradeWsTransceiver.ReadCh()
	for {
		select {
		case <-o.ctx.Done():
			return
		case data := <-orderCh:
			order := data.(*ftxapi.WsOrders)
			if order.Type == "subscribed" {
				continue
			}
			update := &exmodel.OrderTradesUpdate{
				Type:          exmodel.OrderUpdate,
				OrderId:       fmt.Sprint(order.Data.Id),
				ClientOrderId: order.Data.ClientId,
				OrderUpdateInfo: &exmodel.OrderUpdateInfo{
					OrderStatus:  o.parseOrderStatus(order.Data.Status, order.Data.Size == order.Data.FilledSize),
					FilledVolume: fmt.Sprint(order.Data.FilledSize),
					UpdateTime:   time.Now(),
				},
				TradesUpdateInfo: nil,
			}
			// 延迟 50ms 再推送订单更新
			time.AfterFunc(time.Millisecond*50, func() {
				o.outputUpdateCh <- update
			})
		case data := <-tradeCh:
			trade := data.(*ftxapi.WsFills)
			if trade.Type == "subscribed" {
				continue
			}
			update := &exmodel.OrderTradesUpdate{
				Type:            exmodel.TradesUpdate,
				OrderId:         fmt.Sprint(trade.Data.OrderId),
				ClientOrderId:   fmt.Sprint(trade.Data.ClientOrderId),
				OrderUpdateInfo: nil,
				TradesUpdateInfo: &exmodel.TradesUpdateInfo{
					TradeId:     fmt.Sprint(trade.Data.TradeId),
					Price:       fmt.Sprint(trade.Data.Price),
					Volume:      fmt.Sprint(trade.Data.Size),
					Fee:         fmt.Sprint(trade.Data.Fee),
					FeeCurrency: exmodel.NewCurrency(""), // todo,
					Liquidity:   o.parseLiquidity(trade.Data.Liquidity),
					TradeTime:   trade.Data.Time,
				},
			}
			o.outputUpdateCh <- update
		}
	}
}

func (o *tradeManager) PlaceOrder(req exmodel.PlaceOrderReq) (*exmodel.PlaceOrderRsp, error) {
	rsp, err := o.api.PlaceOrder(ftxapi.PlaceOrderReq{
		Market:            req.SymbolExFormat,
		Side:              req.Side.Lower(),
		Price:             req.Price,
		Type:              req.OrderType.Lower(),
		Size:              req.Volume,
		ClientId:          req.ClientOrderId,
		ReduceOnly:        false,
		Ioc:               false,
		PostOnly:          false,
		RejectOnPriceBand: false,
		RejectAfterTs:     "",
	})

	if err != nil {
		logx.Error(err)
		return nil, err
	}

	if !rsp.Success {
		logx.Error("ftx request not success")
		return nil, fmt.Errorf("ftx request not success")
	}

	isFilled := xmath.MustDecimal(req.Volume).Equal(decimal.NewFromFloat(rsp.Result.FilledSize))

	return &exmodel.PlaceOrderRsp{
		OrderId:       fmt.Sprint(rsp.Result.Id),
		ClientOrderId: req.ClientOrderId,
		FilledVolume:  fmt.Sprint(rsp.Result.FilledSize),
		Status:        o.parseOrderStatus(rsp.Result.Status, isFilled),
	}, nil
}

func (o *tradeManager) CancelOrder(req exmodel.CancelOrderReq) (*exmodel.CancelOrderRsp, error) {

	rsp, err := o.api.CancelOrderByClientOrderId(req.ClientOrderId)
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	if !rsp.Success {
		logx.Error("ftx request not success")
		return nil, fmt.Errorf("ftx request not success")
	}

	return &exmodel.CancelOrderRsp{}, nil

}

func (o *tradeManager) QueryOrder(req exmodel.QueryOrderReq) (*exmodel.Order, error) {
	rsp, err := o.api.QueryOrderByClientOrderId(req.ClientOrderId)
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	if !rsp.Success {
		logx.Error("ftx request not success")
		return nil, fmt.Errorf("ftx request not success")
	}

	order := &exmodel.Order{
		OrderId:        fmt.Sprint(rsp.Result.Id),
		ClientOrderId:  req.ClientOrderId,
		Volume:         fmt.Sprint(rsp.Result.Size),
		Price:          fmt.Sprint(rsp.Result.Price),
		FilledVolume:   fmt.Sprint(rsp.Result.FilledSize),
		Exchange:       Name,
		SymbolExFormat: rsp.Result.Market,
		ApiType:        exmodel.ApiTypeUnified,
		Side:           o.parseOrderSide(rsp.Result.Side),
		Type:           o.parseOrderType(rsp.Result.Type),
		Status:         o.parseOrderStatus(rsp.Result.Status, rsp.Result.Size == rsp.Result.FilledSize),
	}

	return order, nil

}

func (o *tradeManager) QueryTrades(req exmodel.QueryTradeReq) ([]*exmodel.Trade, error) {
	rsp, err := o.api.QueryTrades(ftxapi.QueryTradesReq{
		OrderId: req.OrderId,
	})

	if err != nil {
		logx.Error(err)
		return nil, err
	}

	if !rsp.Success {
		logx.Error("ftx request not success")
		return nil, fmt.Errorf("ftx request not success")
	}

	var trades []*exmodel.Trade

	for _, t := range rsp.Result {
		trades = append(trades, &exmodel.Trade{
			TradeId:        fmt.Sprint(t.TradeId),
			OrderId:        fmt.Sprint(t.OrderId),
			Exchange:       Name,
			SymbolExFormat: t.Market,
			ApiType:        exmodel.ApiTypeUnified,
			Liquidity:      o.parseLiquidity(t.Liquidity),
			Volume:         fmt.Sprint(t.Size),
			Price:          fmt.Sprint(t.Price),
			Fee:            fmt.Sprint(t.Fee),
			FeeCurrency:    exmodel.NewCurrency(t.FeeCurrency),
			TradeTime:      t.Time,
		})
	}

	return trades, nil
}

func (o *tradeManager) OutputUpdateCh() <-chan *exmodel.OrderTradesUpdate {
	return o.outputUpdateCh
}

func (o *tradeManager) Close() {
	o.cancel()
}

func (o *tradeManager) parseOrderStatus(in string, isFilled bool) exmodel.OrderStatus {
	switch strings.ToLower(in) {
	case "new": // accepted but not processed yet
		return exmodel.OrderStatusPending
	case "open":
		return exmodel.OrderStatusSent
	case "closed": // filled or cancelled
		if isFilled {
			return exmodel.OrderStatusFilled
		}
		return exmodel.OrderStatusCancelled
	default:
		logx.Errorf("ftx parse a unknown order status:%s", in)
		return exmodel.OrderStatusUnknown
	}
}

func (o *tradeManager) parseOrderSide(in string) exmodel.OrderSide {
	switch strings.ToLower(in) {
	case "buy":
		return exmodel.OrderSideBuy
	case "sell":
		return exmodel.OrderSideSell
	default:
		logx.Errorf("ftx parse a unknown order side:%s", in)
		return exmodel.OrderSideUnknown
	}
}

func (o *tradeManager) parseOrderType(in string) exmodel.OrderType {
	switch strings.ToLower(in) {
	case "limit":
		return exmodel.OrderTypeLimit
	case "market":
		return exmodel.OrderTypeMarket
	default:
		logx.Errorf("ftx parse a unknown order type:%s", in)
		return exmodel.OrderTypeUnknown
	}
}

func (o *tradeManager) parseLiquidity(in string) exmodel.Liquidity {
	switch strings.ToLower(in) {
	case "taker":
		return exmodel.LiquidityTaker
	case "maker":
		return exmodel.LiquidityMaker
	default:
		logx.Errorf("ftx parse a unknown Liquidity:%s", in)
		return exmodel.LiquidityUnknown
	}
}
