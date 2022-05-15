package kafka

import (
	"context"
	"market_aggregate/pkg/datastruct"
)

// DEPTH_TYPE + TYPE_SEPARATOR + symbol + SYMBOL_EXCHANGE_SEPARATOR  + exchange
const (
	DEPTH_TYPE                = "DEPTH"
	TRADE_TYPE                = "TRADE"
	KLINE_TYPE                = "KLINE"
	TYPE_SEPARATOR            = "."
	SYMBOL_EXCHANGE_SEPARATOR = "."
)

func GetDepthTopic(symbol string, exchange string) string {
	return DEPTH_TYPE + TYPE_SEPARATOR + symbol + SYMBOL_EXCHANGE_SEPARATOR + exchange
}

func GetKlineTopic(symbol string, exchange string) string {
	return KLINE_TYPE + TYPE_SEPARATOR + symbol + SYMBOL_EXCHANGE_SEPARATOR + exchange
}

func GetTradeTopic(symbol string, exchange string) string {
	return TRADE_TYPE + TYPE_SEPARATOR + symbol + SYMBOL_EXCHANGE_SEPARATOR + exchange
}

func GetTopics(MetaData datastruct.Metadata) map[string]struct{} {
	rst := make(map[string]struct{})

	for symbol, exchange_map := range MetaData.DepthMeta {
		for exchange, _ := range exchange_map {
			topic := GetDepthTopic(symbol, exchange)
			if _, ok := rst[topic]; ok == false {
				rst[topic] = struct{}{}
			}
		}
	}

	for symbol, exchange_map := range MetaData.KlineMeta {
		for exchange, _ := range exchange_map {
			topic := GetKlineTopic(symbol, exchange)
			if _, ok := rst[topic]; ok == false {
				rst[topic] = struct{}{}
			}
		}
	}

	for symbol, exchange_map := range MetaData.TradeMeta {
		for exchange, _ := range exchange_map {
			topic := GetTradeTopic(symbol, exchange)
			if _, ok := rst[topic]; ok == false {
				rst[topic] = struct{}{}
			}
		}
	}
	return rst
}

func GetConsumeSet(MetaData datastruct.Metadata) map[string](*ConsumeItem) {
	rst := make(map[string](*ConsumeItem))

	for symbol, exchange_map := range MetaData.DepthMeta {
		for exchange, _ := range exchange_map {
			topic := GetDepthTopic(symbol, exchange)
			if _, ok := rst[topic]; ok == false {

				base_ctx := context.Background()
				base_child_ctx, base_child_cancel_func := context.WithCancel(base_ctx)

				consume_item := ConsumeItem{
					Topic:      topic,
					Ctx:        base_child_ctx,
					CancelFunc: base_child_cancel_func,
				}
				rst[topic] = &consume_item
			}
		}
	}

	for symbol, exchange_map := range MetaData.KlineMeta {
		for exchange, _ := range exchange_map {
			topic := GetKlineTopic(symbol, exchange)
			if _, ok := rst[topic]; ok == false {
				base_ctx := context.Background()
				base_child_ctx, base_child_cancel_func := context.WithCancel(base_ctx)

				consume_item := ConsumeItem{
					Topic:      topic,
					Ctx:        base_child_ctx,
					CancelFunc: base_child_cancel_func,
				}
				rst[topic] = &consume_item
			}
		}
	}

	for symbol, exchange_map := range MetaData.TradeMeta {
		for exchange, _ := range exchange_map {
			topic := GetTradeTopic(symbol, exchange)
			if _, ok := rst[topic]; ok == false {
				base_ctx := context.Background()
				base_child_ctx, base_child_cancel_func := context.WithCancel(base_ctx)

				consume_item := ConsumeItem{
					Topic:      topic,
					Ctx:        base_child_ctx,
					CancelFunc: base_child_cancel_func,
				}
				rst[topic] = &consume_item
			}
		}
	}
	return rst
}
