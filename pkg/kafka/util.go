package kafka

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
