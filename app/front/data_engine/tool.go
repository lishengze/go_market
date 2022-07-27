package data_engine

import (
	"market_server/app/data_manager/rpc/marketservice"
	"market_server/app/data_manager/rpc/types/pb"
	"market_server/common/datastruct"
)

func TrasPbKlines(pb_klines []*pb.Kline) []*datastruct.Kline {
	var rst []*datastruct.Kline = nil

	for _, pb_kline := range pb_klines {
		kline := marketservice.NewKlineWithPbKline(pb_kline)
		if kline == nil {
			continue
		}

		rst = append(rst, kline)
	}
	return rst
}
