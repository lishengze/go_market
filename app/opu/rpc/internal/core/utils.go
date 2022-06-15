package core

import (
	"exterior-interactor/app/opu/model"
	"exterior-interactor/app/opu/rpc/opupb"
	"exterior-interactor/pkg/exchangeapi/exmodel"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func toPbOrder(order *model.Order) *opupb.Order {
	return &opupb.Order{
		Id:            order.Id,
		AccountId:     order.AccountId,
		ClientOrderId: order.ClientOrderId,
		AccountAlias:  order.AccountAlias,
		ExOrderId:     order.ExOrderId,
		Exchange:      order.Exchange,
		Volume:        order.Volume,
		Price:         order.Price,
		Symbol:        order.StdSymbol,
		Side:          order.Side,
		Status:        order.Status,
		Type:          order.Tp,
		CreateTime:    timestamppb.New(order.CreateTime),
		UpdateTime:    timestamppb.New(order.UpdateTime),
		FilledVolume:  order.FilledVolume,
	}
}

func toPbTrade(trade *model.Trade) *opupb.Trade {
	return &opupb.Trade{
		Id:          trade.Id,
		Volume:      trade.Volume,
		Price:       trade.Price,
		Liquidity:   trade.Liquidity,
		Fee:         trade.Fee,
		FeeCurrency: trade.FeeCurrency,
		TradeTime:   timestamppb.New(trade.TradeTime),
		ExTradeId:   trade.ExTradeId,
	}
}

func toPbTrades(trades []*model.Trade) []*opupb.Trade {
	var res []*opupb.Trade
	for _, trade := range trades {
		res = append(res, toPbTrade(trade))
	}
	return res
}

func toPbOrderUpdate(order *model.Order) *opupb.OrderTradesUpdate {
	return &opupb.OrderTradesUpdate{
		UpdateType:    exmodel.OrderUpdate.String(),
		OrderId:       order.Id,
		AccountId:     order.AccountId,
		ClientOrderId: order.ClientOrderId,
		ExOrderId:     order.ExOrderId,
		Exchange:      order.Exchange,
		Volume:        order.Volume,
		Price:         order.Price,
		Symbol:        order.StdSymbol,
		Side:          order.Side,
		OrderType:     order.Tp,
		CreateTime:    timestamppb.New(order.CreateTime),
		OrderUpdate: &opupb.OrderUpdate{
			Status:       order.Status,
			FilledVolume: order.FilledVolume,
			UpdateTime:   timestamppb.New(time.Now()),
		},
		TradeUpdate: nil,
	}
}

func toPbTradesUpdate(order *model.Order, trade *model.Trade) *opupb.OrderTradesUpdate {
	return &opupb.OrderTradesUpdate{
		UpdateType:    exmodel.TradesUpdate.String(),
		OrderId:       order.Id,
		AccountId:     order.AccountId,
		AccountAlias:  order.AccountAlias,
		ClientOrderId: order.ClientOrderId,
		ExOrderId:     order.ExOrderId,
		Exchange:      order.Exchange,
		Volume:        order.Volume,
		Price:         order.Price,
		Symbol:        order.StdSymbol,
		Side:          order.Side,
		OrderType:     order.Tp,
		CreateTime:    timestamppb.New(order.CreateTime),
		OrderUpdate:   nil,
		TradeUpdate: &opupb.TradesUpdate{
			Id:          trade.Id,
			ExTradeId:   trade.ExTradeId,
			Volume:      trade.Volume,
			Price:       trade.Price,
			Liquidity:   trade.Liquidity,
			TradeFee:    trade.Fee,
			FeeCurrency: trade.FeeCurrency,
			TradeTime:   timestamppb.New(trade.TradeTime),
		},
	}
}

func toPbBalances(balances []*exmodel.Balance) []*opupb.Balance {
	var res []*opupb.Balance

	for _, b := range balances {
		pbBalance := &opupb.Balance{
			WalletType:     b.WalletType.String(),
			BalanceDetails: nil,
		}

		for _, d := range b.Details {
			pbBalance.BalanceDetails = append(pbBalance.BalanceDetails, &opupb.BalanceDetail{
				Currency:  d.Currency.String(),
				Available: d.Available,
				Total:     d.Total,
			})
		}

		res = append(res, pbBalance)
	}

	return res
}

func toPbSymbols(symbols []*model.Symbol) []*opupb.Symbol {
	var res []*opupb.Symbol
	for _, s := range symbols {
		res = append(res, &opupb.Symbol{
			Exchange:      s.Exchange,
			ExFormat:      s.ExFormat,
			StdFormat:     s.StdSymbol,
			Type:          s.Tp,
			VolumeScale:   s.VolumeScale,
			PriceScale:    s.PriceScale,
			MinVolume:     s.MinVolume,
			ContractSize:  s.ContractSize,
			BaseCurrency:  s.BaseCurrency,
			QuoteCurrency: s.QuoteCurrency,
		})
	}
	return res
}
