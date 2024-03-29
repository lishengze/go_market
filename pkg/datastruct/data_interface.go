package datastruct

type SerializerI interface {
	EncodeDepth(*DepthQuote) ([]byte, error)
	EncodeKline(*Kline) ([]byte, error)
	EncodeTrade(*Trade) ([]byte, error)

	DecodeDepth([]byte) (*DepthQuote, error)
	DecodeKline([]byte) (*Kline, error)
	DecodeTrade([]byte) (*Trade, error)
}

type NetServerI interface {
	Init(SerializerI, *DataChannel, *DataChannel) error
	Start()
	UpdateMetaData(*Metadata)

	PublishDepth(*DepthQuote) error
	PublishKline(*Kline) error
	PublishTrade(*Trade) error

	SendRecvedDepth(*DepthQuote)
	SendRecvedKline(*Kline)
	SendRecvedTrade(*Trade)
}
