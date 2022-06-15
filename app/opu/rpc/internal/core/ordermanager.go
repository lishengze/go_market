package core

import (
	"context"
	"exterior-interactor/app/opu/model"
	"exterior-interactor/app/opu/rpc/opupb"
	"exterior-interactor/pkg/exchangeapi/exmodel"
	"exterior-interactor/pkg/xmath"
	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
	"time"
)

type orderManager struct {
	svcCtx              *svcCtx
	getAccountManagerFn func(order *model.Order) (*accountManager, error)
	symbol              *model.Symbol
	order               *model.Order
	trades              []*model.Trade
	inputCh             chan *exmodel.OrderTradesUpdate
	outputCh            chan<- *opupb.OrderTradesUpdate
	cancel              context.CancelFunc
	ctx                 context.Context
	mutex               sync.Mutex
}

func newOrderManager(getAccountManagerFn func(order *model.Order) (*accountManager, error), svcCtx *svcCtx,
	outputCh chan<- *opupb.OrderTradesUpdate, symbol *model.Symbol, order *model.Order, trades []*model.Trade) *orderManager {

	ctx, cancel := context.WithCancel(context.Background())

	om := &orderManager{
		svcCtx:              svcCtx,
		getAccountManagerFn: getAccountManagerFn,
		order:               order,
		trades:              trades,
		inputCh:             make(chan *exmodel.OrderTradesUpdate, 256),
		outputCh:            outputCh,
		cancel:              cancel,
		ctx:                 ctx,
		symbol:              symbol,
		mutex:               sync.Mutex{},
	}

	om.sendOrder() // 发送订单

	go om.run()
	go om.startStatusHeartBeat()

	return om

}

// newOrderManagerWithoutSend 不向交易所发送订单，用于加载历史未完成的订单
func newOrderManagerWithoutSend(getAccountManagerFn func(order *model.Order) (*accountManager, error), svcCtx *svcCtx,
	outputCh chan<- *opupb.OrderTradesUpdate, symbol *model.Symbol, order *model.Order, trades []*model.Trade) *orderManager {

	ctx, cancel := context.WithCancel(context.Background())

	om := &orderManager{
		svcCtx:              svcCtx,
		getAccountManagerFn: getAccountManagerFn,
		order:               order,
		trades:              trades,
		inputCh:             make(chan *exmodel.OrderTradesUpdate, 256),
		outputCh:            outputCh,
		cancel:              cancel,
		ctx:                 ctx,
		symbol:              symbol,
		mutex:               sync.Mutex{},
	}

	go om.run()
	go om.startStatusHeartBeat()

	return om

}

// sendOrder 发送订单
func (o *orderManager) sendOrder() {
	account, err := o.getAccountManagerFn(o.order)
	if err != nil {
		logx.Errorf("getAccountManagerFn err:%s, order:%+v", err, o.order)
		return
	}

	req := exmodel.PlaceOrderReq{
		SymbolExFormat: o.order.StdSymbol,
		ApiType:        exmodel.ApiType(o.order.ApiType),
		ClientOrderId:  o.order.Id,
		OrderType:      exmodel.OrderType(o.order.Tp),
		Side:           exmodel.OrderSide(o.order.Side),
		Volume:         o.order.Volume,
		Price:          o.order.Price,
	}

	_, err = account.TradeManager.PlaceOrder(req)
	if err != nil {
		logx.Errorf("place order err:%s, req:%+v, order:%+v", err, req, o.order)
		// todo 向交易所查询订单是否下单真的失败了
		o.updateOrder("", "0", err.Error(), exmodel.OrderStatusRejected)
		return
	}
}

func (o *orderManager) cancelOrder() {
	o.mutex.Lock()
	cancelFlag := o.order.CancelFlag
	o.mutex.Unlock()

	if cancelFlag == "SET" { // 已经下下达撤单指令
		return
	}

	err := o.svcCtx.OrderModel.Update(o.order, func() {
		o.order.CancelFlag = "SET"
	})

	if err != nil {
		logx.Errorf("order set CancelFlag err:%v, order:%+v", err, *o.order)
	}

	// 循环撤单
	for {
		o.mutex.Lock()
		orderStatus := o.order.Status
		o.mutex.Unlock()

		if exmodel.OrderStatus(orderStatus).IsClosed() { // 订单关闭就不需要重复撤单了
			return
		}

		account, err := o.getAccountManagerFn(o.order)
		if err != nil {
			logx.Errorf("getAccountManagerFn err:%s, order:%+v", err, o.order)
			return
		}

		req := exmodel.CancelOrderReq{
			OrderId:        o.order.ExOrderId,
			ClientOrderId:  o.order.Id,
			SymbolExFormat: o.order.ExSymbol,
			ApiType:        exmodel.ApiType(o.order.ApiType),
		}

		_, err = account.TradeManager.CancelOrder(req)
		if err != nil {
			logx.Errorf("cancel order err:%s, req:%+v, order:%+v", err, req, o.order)
		}

		time.Sleep(time.Second)
	}
}

func (o *orderManager) run() {
	for {
		select {
		case <-o.ctx.Done():
			goto exit
		case update := <-o.inputCh:
			switch update.Type {
			case exmodel.OrderUpdate:
				o.updateOrder(update.OrderId, update.FilledVolume, "", update.OrderStatus)
			case exmodel.TradesUpdate:
				feeCurrency := o.getFeeCurrency(update)
				o.mutex.Lock() // 此处加锁
				o.updateTrade(update.TradeId, update.Fee, feeCurrency, update.Volume, update.Price, update.Liquidity, update.TradeTime)
				o.mutex.Unlock()
			}
		}
	}
exit:
	logx.Infof("[OrderManager Quit]orderId:%s, orderStats:%s ", o.order.Id, o.order.Status)
}

// startStatusHeartBeat  定时去交易所同步状态
// 心跳 同步间隔逐步从 3 秒 提升到 3min
func (o *orderManager) startStatusHeartBeat() {
	logx.Infof("[OrderStatusHeartBeat Start] order_id:%s", o.order.Id)
	d := time.Second * 3
	maxD := time.Minute * 3
	ticker := time.NewTicker(d)
	defer ticker.Stop()
	for {
		select {
		case <-o.ctx.Done():
			goto exit
		case <-ticker.C:
			logx.Infof("[OrderStatusHeartBeat SyncOrder], orderId:%s", o.order.Id)
			o.syncOrder()
			if d < maxD {
				d += d + 1
				ticker.Reset(d)
			}
		}
	}
exit:
	logx.Infof("[OrderStatusHeartBeat Quit]orderId:%s ", o.order.Id)
}

// syncOrder 同步订单信息
func (o *orderManager) syncOrder() {
	account, err := o.getAccountManagerFn(o.order)
	if err != nil {
		logx.Errorf("getAccountManagerFn err:%s, order:%+v", err, o.order)
		return
	}

	req := exmodel.QueryOrderReq{
		OrderId:        o.order.ExOrderId,
		ClientOrderId:  o.order.Id,
		SymbolExFormat: o.order.ExSymbol,
		ApiType:        exmodel.ApiType(o.order.ApiType),
	}

	order, err := account.TradeManager.QueryOrder(req)
	if err != nil {
		logx.Errorf("queryOrder err:%s, req:%+v, order:%+v", err, req, o.order)
		return
	}

	o.updateOrder(order.OrderId, order.FilledVolume, "", order.Status)
}

// inputUpdate 把订单的更新传给 orderManager
func (o *orderManager) inputUpdate(update *exmodel.OrderTradesUpdate) {
	o.inputCh <- update
}

// updateOrder 更新订单, 需要获取锁
func (o *orderManager) updateOrder(exOrderId, filledVolume, rejectReason string, status exmodel.OrderStatus) {
	defer o.mutex.Unlock()
	o.mutex.Lock()

	//  保存订单信息 并 推送订单
	saveAndOutputOrderUpdate := func() {
		err := o.svcCtx.OrderModel.Update(o.order, func() {
			o.order.ExOrderId = exOrderId
			o.order.FilledVolume = filledVolume
			o.order.Status = status.String()
			o.order.RejectReason = rejectReason
		})
		if err != nil {
			logx.Errorf("update order to DB failed, err:%v, order:%+v", err, *o.order)
		}

		o.outputOrderUpdate() // 推送订单更新
	}

	switch exmodel.OrderStatus(o.order.Status) {
	case exmodel.OrderStatusPending: // 当前为 pending
		switch status {
		case exmodel.OrderStatusPending: // 不管
		case exmodel.OrderStatusRejected, exmodel.OrderStatusSent,
			exmodel.OrderStatusCancelling, exmodel.OrderStatusPartial:
			saveAndOutputOrderUpdate()
		case exmodel.OrderStatusCancelled, exmodel.OrderStatusFilled:
			if !o.tradesFilledVolume().Equal(xmath.MustDecimal(filledVolume)) {
				logx.Errorf("[tradesFilledVolume:%s != filledVolume:%s ] order:%+v", o.tradesFilledVolume(), filledVolume, o.order)
				trades := o.mustQueryTrades()
				o.processQueryTradeRsp(trades) // 这里会先推送遗漏的trade
			}
			saveAndOutputOrderUpdate()

		default:
			logx.Errorf("order:%+v, rcv wrong status :%s", o.order, status)
		}
	case exmodel.OrderStatusSent: // 当前为 sent
		switch status {
		case exmodel.OrderStatusSent: //  再次收到 sent, 不管
		case exmodel.OrderStatusCancelling, exmodel.OrderStatusPartial, exmodel.OrderStatusRejected:
			saveAndOutputOrderUpdate()
		case exmodel.OrderStatusCancelled, exmodel.OrderStatusFilled:
			if !o.tradesFilledVolume().Equal(xmath.MustDecimal(filledVolume)) {
				logx.Errorf("[tradesFilledVolume:%s is not equal filledVolume:%s ] order:%+v", o.tradesFilledVolume(), filledVolume, o.order)
				trades := o.mustQueryTrades()
				o.processQueryTradeRsp(trades) // 这里会先推送遗漏的trade
			}
			saveAndOutputOrderUpdate()
		default:
			logx.Errorf("order:%+v, rcv wrong status :%s", o.order, status)
		}
	case exmodel.OrderStatusPartial: // 当前为 partial
		switch status {
		case exmodel.OrderStatusPartial:
			// 比较一下filledVolume
			if xmath.MustDecimal(filledVolume).GreaterThan(xmath.MustDecimal(o.order.FilledVolume)) {
				saveAndOutputOrderUpdate()
			} else {
				logx.Errorf("rcv a delay partial status msg, order:%+v, [filledVolume:%s]", o.order, filledVolume)
			}
		case exmodel.OrderStatusCancelling, exmodel.OrderStatusRejected:
			saveAndOutputOrderUpdate()
		case exmodel.OrderStatusCancelled, exmodel.OrderStatusFilled:
			if !o.tradesFilledVolume().Equal(xmath.MustDecimal(filledVolume)) {
				logx.Errorf("[tradesFilledVolume:%s != filledVolume:%s ] order:%+v", o.tradesFilledVolume(), filledVolume, o.order)
				trades := o.mustQueryTrades()
				o.processQueryTradeRsp(trades) // 这里会先推送遗漏的trade
			}
			saveAndOutputOrderUpdate()
		default:
			logx.Errorf("order:%+v, rcv wrong status :%s", o.order, status)
		}
	case exmodel.OrderStatusCancelling: // 当前为 cancelling
		switch status {
		case exmodel.OrderStatusCancelled, exmodel.OrderStatusFilled:
			if !o.tradesFilledVolume().Equal(xmath.MustDecimal(filledVolume)) {
				logx.Errorf("[tradesFilledVolume:%s != filledVolume:%s ] order:%+v", o.tradesFilledVolume(), filledVolume, o.order)
				trades := o.mustQueryTrades()
				o.processQueryTradeRsp(trades) // 这里会先推送遗漏的trade
			}
			saveAndOutputOrderUpdate()
		default:
			logx.Errorf("order:%+v, rcv wrong status :%s", o.order, status)
		}
	default:
		// 这个状态的 order 不再接收 order update 了
		logx.Errorf("order:%+v, rcv wrong status :%s", o.order, status)
	}

	if exmodel.OrderStatus(o.order.Status).IsClosed() {
		o.close()
	}
}

// updateTrade 更新成交，此方法中未加锁
func (o *orderManager) updateTrade(tradeId, fee, feeCurrency, volume, price string,
	liquidity exmodel.Liquidity, tradeTime time.Time) {

	var m = map[string]struct{}{}
	for _, t := range o.trades {
		m[t.ExTradeId] = struct{}{}
	}

	if _, ok := m[tradeId]; ok {
		// 已存在
		logx.Infof("trade exists, tradeId:%s", tradeId)
		return
	}
	trade := &model.Trade{
		Id:          o.svcCtx.IdSrv.MustGetId(),
		OrderId:     o.order.Id,
		ExTradeId:   tradeId,
		Exchange:    o.order.Exchange,
		StdSymbol:   o.order.StdSymbol,
		Liquidity:   liquidity.String(),
		Side:        o.order.Side,
		Volume:      volume,
		Price:       price,
		Fee:         fee,
		FeeCurrency: feeCurrency,
		TradeTime:   tradeTime,
		CreateTime:  time.Time{},
		UpdateTime:  time.Time{},
	}

	o.trades = append(o.trades, trade)

	o.outputTradesUpdate() // 推送成交
	_, err := o.svcCtx.TradeModel.Insert(trade)
	if err != nil {
		logx.Errorf("insert trade err:%v, trade:%+v", err, *trade)
	}
}

func (o *orderManager) tradesFilledVolume() decimal.Decimal {
	var total decimal.Decimal
	for _, t := range o.trades {
		total = total.Add(xmath.MustDecimal(t.Volume))
	}
	return total
}

func (o *orderManager) mustQueryTrades() []*exmodel.Trade {
	for {
		account, err := o.getAccountManagerFn(o.order)
		if err != nil {
			logx.Errorf("getAccountManagerFn err:%s, order:%+v", err, o.order)
			time.Sleep(time.Second)
			continue
		}

		req := exmodel.QueryTradeReq{
			OrderId:        o.order.ExOrderId,
			SymbolExFormat: o.order.ExSymbol,
			ApiType:        exmodel.ApiType(o.order.ApiType),
			StartTime:      o.order.CreateTime.Add(-time.Minute),
		}

		trades, err := account.TradeManager.QueryTrades(req)

		if err != nil {
			logx.Errorf("queryTrades err:%s, req:%+v, order:%+v", err, req, o.order)
			time.Sleep(time.Second * 5)
			continue
		}

		return trades
	}
}

func (o *orderManager) processQueryTradeRsp(trades []*exmodel.Trade) {
	for _, trade := range trades {
		o.updateTrade(trade.TradeId, trade.Fee, trade.FeeCurrency.String(),
			trade.Volume, trade.Price, trade.Liquidity, trade.TradeTime)
	}
}

// outputOrderUpdate 向外输出订单的更新
func (o *orderManager) outputOrderUpdate() {
	update := toPbOrderUpdate(o.order)
	logx.Infof("[ORDER UPDATE]: %s", update.String())
	o.outputCh <- update
}

// outputTradesUpdate 向外输出成交的更新
func (o *orderManager) outputTradesUpdate() {
	update := toPbTradesUpdate(o.order, o.trades[len(o.trades)-1]) // 最后一条trade
	logx.Infof("[TRADE UPDATE]: %s", update.String())
	o.outputCh <- update
}

func (o *orderManager) close() {
	o.cancel()
}

func (o *orderManager) orderIsClosed() bool {
	defer o.mutex.Unlock()
	o.mutex.Lock()

	return exmodel.OrderStatus(o.order.Status).IsClosed()
}

func (o *orderManager) getFeeCurrency(update *exmodel.OrderTradesUpdate) string {

	// ftx 不推送 fee currency, 根据ftx规则推断一下
	if o.order.Exchange == exmodel.FTX.String() {
		/*
			期货: USD
			现货：
				maker :
					sell: quote currency
					buy:  base currency
				taker :
					quote currency
		*/

		switch o.symbol.Tp {
		case exmodel.SymbolTypeSpot.String():
			switch update.TradesUpdateInfo.Liquidity {
			case exmodel.LiquidityMaker:
				switch o.order.Side {
				case exmodel.OrderSideBuy.String():
					return o.symbol.BaseCurrency
				case exmodel.OrderSideSell.String():
					return o.symbol.QuoteCurrency
				}
			case exmodel.LiquidityTaker:
				return o.symbol.QuoteCurrency
			}
		default:
			return "USD"
		}

		logx.Errorf("can't parse ftx fee currency, update:%+v, symbol:%+v", *update, *o.symbol)
		return ""
	}

	return update.TradesUpdateInfo.FeeCurrency.String()
}
