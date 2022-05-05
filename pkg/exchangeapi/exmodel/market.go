package exmodel

import (
	"time"
)

const (
	MTypeDepth MType = "M-DEPTH"
	MTypeTrade MType = "M-TRADE"
)

type (
	MType string

	// StreamDepth 流式 深度行情
	StreamDepth struct {
		Exchange Exchange
		Time     time.Time
		Symbol   *Symbol
		Asks     [][2]string
		Bids     [][2]string
	}

	// StreamMarketTrade 流式 行情 trade
	StreamMarketTrade struct {
		TradeId  string
		Exchange Exchange
		Time     time.Time
		Symbol   *Symbol
		Price    string
		Volume   string
	}

	Kline struct {
		Exchange   Exchange
		Time       time.Time
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
