package main

import (
	"encoding/json"
	"market_server/app/front/net"

	"github.com/zeromicro/go-zero/core/logx"
)

type Client struct {
}

func GetTestTradeReqJson() []byte {
	symbol_list := []string{"BTC_USDT"}
	sub_info := map[string]interface{}{
		"type":   net.TRADE,
		"symbol": symbol_list,
	}
	rst, err := json.Marshal(sub_info)

	if err != nil {
		logx.Errorf("GetTestTradeReqJson: %+v \n", err)
		return nil
	} else {
		logx.Infof("SubJson: %s", string(rst))
		return rst
	}
}

func GetTestDepthReqJson() []byte {
	symbol_list := []string{"BTC_USDT"}
	sub_info := map[string]interface{}{
		"type":   net.SYMBOL_SUB,
		"symbol": symbol_list,
	}
	rst, err := json.Marshal(sub_info)

	if err != nil {
		logx.Errorf("GetTestDepthReqJson: %+v \n", err)
		return nil
	} else {
		logx.Infof("SubJson: %s", string(rst))
		return rst
	}
}

func GetTestKlineReqJson() []byte {
	sub_info := map[string]interface{}{
		"type":       net.KLINE_UPDATE,
		"symbol":     "BTC_USDT",
		"data_count": 2,
		"frequency":  500,
	}
	rst, err := json.Marshal(sub_info)
	if err != nil {
		logx.Errorf("GetTestKlineReqJson: %+v \n", err)
		return nil
	} else {
		logx.Infof("SubJson: %s", string(rst))
		return rst
	}
}
