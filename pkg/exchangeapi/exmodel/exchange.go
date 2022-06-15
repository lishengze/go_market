package exmodel

type Exchange string

const (
	ALAMEDA       Exchange = "ALAMEDA"
	AMBER         Exchange = "AMBER"
	B2C2          Exchange = "B2C2"
	BINANCE       Exchange = "BINANCE"
	BITMEX        Exchange = "BITMEX"
	BLOCKFI       Exchange = "BLOCKFI"
	BYBIT         Exchange = "BYBIT"
	CUMBERLAND    Exchange = "CUMBERLAND"
	DERIBIT       Exchange = "DERIBIT"
	EXTERNAL      Exchange = "EXTERNAL"
	FALCONX       Exchange = "FALCONX"
	FTX           Exchange = "FTX"
	GALAXYDIGITA  Exchange = "GALAXYDIGITA"
	HKCUSTODY     Exchange = "HKCUSTODY"
	HUOBI         Exchange = "HUOBI"
	JUMP          Exchange = "JUMP"
	MATRIXPORT    Exchange = "MATRIXPORT"
	OKCOIN        Exchange = "OKCOIN"
	OKEX          Exchange = "OKEX"
	OKLINK        Exchange = "OKLINK"
	OSL           Exchange = "OSL"
	OWS           Exchange = "OWS"
	PAXOS         Exchange = "PAXOS"
	POLONIEX      Exchange = "POLONIEX"
	PRIMETRUST    Exchange = "PRIMETRUST"
	SIGNET        Exchange = "SIGNET"
	SILVERGATE    Exchange = "SILVERGATE"
	XDAEX         Exchange = "XDAEX"
	CoinMarketCap Exchange = "CMC"
)


func (o Exchange) String() string {
	return string(o)
}
