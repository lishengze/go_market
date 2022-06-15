package exmodel

import (
	"time"
)

type (
	// StreamDepth 流式 深度行情
	StreamDepth struct {
		Exchange  Exchange
		Time      time.Time // 交易所时间
		LocalTime time.Time // 服务器本地时间
		Symbol    *Symbol
		Asks      [][2]string // 价格从低到高
		Bids      [][2]string // 价格从高到低
	}

	// StreamMarketTrade 流式 行情 trade
	StreamMarketTrade struct {
		TradeId   string
		Exchange  Exchange
		Time      time.Time // 交易所时间
		LocalTime time.Time // 服务器本地时间
		Symbol    *Symbol
		Price     string
		Volume    string
	}

	Kline struct {
		Exchange   Exchange
		Time       time.Time
		LocalTime  time.Time // 服务器本地时间
		Symbol     *Symbol
		Resolution uint32 // 分辨率 分钟k线，小时k线等
		Open       string
		High       string
		Low        string
		Close      string
		Volume     string // 量
		Value      string // 额
	}
)
