package main

import (
	"fmt"
	"market_server/app/front/data_engine"
	"market_server/common/datastruct"

	"github.com/zeromicro/go-zero/core/logx"
)

func InitLogx() {

	LogConfig := logx.LogConf{
		Compress:            true,
		KeepDays:            0,
		Level:               "info",
		Mode:                "file",
		Path:                "./log",
		ServiceName:         "client",
		StackCooldownMillis: 100,
		TimeFormat:          "2006-01-02 15:04:05",
	}

	logx.MustSetup(LogConfig)
}

func test_main() {
	InitLogx()

	data_worker := data_engine.NewDataEngine(nil, nil)
	data_worker.SetTestFlag(true)

	req_kline_info := &datastruct.ReqHistKline{
		Symbol:    "BTC_USDT",
		Exchange:  "FTX",
		StartTime: 0,
		EndTime:   0,
		Count:     2,
		Frequency: 300,
	}

	data_worker.GetHistKlineData(req_kline_info)

}

func main() {

	fmt.Println("Test DataEngine!")

	test_main()
}
