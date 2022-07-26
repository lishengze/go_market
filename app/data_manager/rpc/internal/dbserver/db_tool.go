package dbserver

import (
	"database/sql"
	"market_server/app/data_manager/rpc/types/pb"
	"market_server/common/datastruct"
	"market_server/common/util"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewPbKlineWithPbKline(ori_kline *pb.Kline) *pb.Kline {
	defer util.CatchExp("NewPbKlineWithPbKline")
	return nil
}

func NewPbKlineWithKline(ori_kline *datastruct.Kline) *pb.Kline {
	defer util.CatchExp("NewPbKlineWithPbKline")

	if ori_kline == nil {
		return nil
	}

	return &pb.Kline{
		Exchange:   ori_kline.Exchange,
		Symbol:     ori_kline.Symbol,
		Timestamp:  &timestamppb.Timestamp{Seconds: int64(ori_kline.Time) / datastruct.NANO_PER_SECS, Nanos: int32(ori_kline.Time % datastruct.NANO_PER_SECS)},
		Resolution: uint32(ori_kline.Resolution),
		Open:       strconv.FormatFloat(ori_kline.Open, 'g', 8, 64),
		High:       strconv.FormatFloat(ori_kline.High, 'g', 8, 64),
		Low:        strconv.FormatFloat(ori_kline.Low, 'g', 8, 64),
		Close:      strconv.FormatFloat(ori_kline.Close, 'g', 8, 64),
		Volume:     strconv.FormatFloat(ori_kline.Volume, 'g', 8, 64),
	}
}

func GetTimestamp(time int64) *timestamppb.Timestamp {
	defer util.CatchExp("GetTimestamp")
	return &timestamppb.Timestamp{Seconds: int64(time) / datastruct.NANO_PER_SECS, Nanos: int32(time % datastruct.NANO_PER_SECS)}
}

func GetKline(values []interface{}) *datastruct.Kline {
	defer util.CatchExp("GetPbKline")
	tmp_kline := &datastruct.Kline{}
	tmp_kline.Exchange = string(values[0].([]byte))
	tmp_kline.Symbol = string(values[1].([]byte))

	time, _ := strconv.Atoi(string(values[2].([]byte)))
	tmp_kline.Time = int64(time)

	open, err := strconv.ParseFloat(string(values[3].([]byte)), 64)
	if err != nil {
		logx.Error("strconv.ParseFloat open %s, Failed; Err: %s", string(values[3].([]byte)), err.Error())
		return nil
	}
	tmp_kline.Open = open

	high, err := strconv.ParseFloat(string(values[4].([]byte)), 64)
	if err != nil {
		logx.Error("strconv.ParseFloat high %s, Failed; Err: %s", string(values[4].([]byte)), err.Error())
		return nil
	}
	tmp_kline.High = high

	low, err := strconv.ParseFloat(string(values[5].([]byte)), 64)
	if err != nil {
		logx.Error("strconv.ParseFloat low %s, Failed; Err: %s", string(values[5].([]byte)), err.Error())
		return nil
	}
	tmp_kline.Low = low

	close, err := strconv.ParseFloat(string(values[6].([]byte)), 64)
	if err != nil {
		logx.Error("strconv.ParseFloat close %s, Failed; Err: %s", string(values[6].([]byte)), err.Error())
		return nil
	}
	tmp_kline.Close = close

	volume, err := strconv.ParseFloat(string(values[7].([]byte)), 64)
	if err != nil {
		logx.Error("strconv.ParseFloat volume %s, Failed; Err: %s", string(values[7].([]byte)), err.Error())
		return nil
	}
	tmp_kline.Volume = volume

	resolution, _ := strconv.Atoi(string(values[8].([]byte)))
	tmp_kline.Resolution = resolution

	return tmp_kline
}

func GetOriPbKline(kline_db_row *sql.Rows) []*datastruct.Kline {
	var rst []*datastruct.Kline

	columns, _ := kline_db_row.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for kline_db_row.Next() {
		err := kline_db_row.Scan(scanArgs...)

		if err != nil {
			logx.Error(err.Error())
			return rst
		}

		tmp_kline := GetKline(values)
		rst = append(rst, tmp_kline)
	}

	return rst
}

func TransKlineData(klines []*datastruct.Kline) []*pb.Kline {
	var rst []*pb.Kline = nil

	for _, kline := range klines {
		pb_kline := NewPbKlineWithKline(kline)
		rst = append(rst, pb_kline)
	}

	return rst
}
