package datastruct

import "market_aggregate/pkg/conf"

type SerializerI interface {
	EncodeDepth(*DepthQuote) ([]byte, error)
	EncodeKline(*Kline) ([]byte, error)
	EncodeTrade(*Trade) ([]byte, error)

	DecodeDepth([]byte) (*DepthQuote, error)
	DecodeKline([]byte) (*Kline, error)
	DecodeTrade([]byte) (*Trade, error)
}

type NetServerI interface {
	Init(*conf.Config, SerializerI, *DataChannel)
	// Start()
	// SetMetaData()

	PublishDepth(*DepthQuote)
	PublishKline(*Kline)
	PublishTrade(*Trade)

	SendRecvedDepth(*DepthQuote)
	SendRecvedKline(*Kline)
	SendRedvedTrade(*Trade)
}
