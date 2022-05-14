package comm

// func NewTradeWithProtoTrade(proto_trade *kafka.Trade) *Trade {
// 	price, err := strconv.ParseFloat(proto_trade.Price, 64)
// 	if err != nil {
// 		panic(err.Error)
// 	}

// 	volume, err := strconv.ParseFloat(proto_trade.Volume, 64)
// 	if err != nil {
// 		panic(err.Error)
// 	}

// 	return &Trade{
// 		Exchange: proto_trade.Exchange,
// 		Symbol:   proto_trade.Symbol,
// 		Price:    price,
// 		Volume:   volume,
// 		Time:     int64(proto_trade.Timestamp.Nanos),
// 	}
// }

// func NewDepthQuoteWithProtoDepth(proto_depth *kafka.Depth) *DepthQuote {
// 	asks := treemap.NewWith(utils.Float64Comparator)
// 	bids := treemap.NewWith(utils.Float64Comparator)

// 	// for item

// 	return &DepthQuote{
// 		Exchange: proto_depth.Exchange,
// 		Symbol:   proto_depth.Symbol,
// 		Time:     int64(proto_depth.Timestamp.Nanos),
// 		Asks:     asks,
// 		Bids:     bids,
// 	}
// }
