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

func ResetFirstKline(cache_kline *datastruct.Kline, target_resolution uint32) {
	defer util.CatchExp("ResetFirstKline")
	cache_kline.Time = datastruct.GetLastStartTime(cache_kline.Time, int64(target_resolution))
}

func ProcessOldEndKline(cur_kline *datastruct.Kline, cache_kline *datastruct.Kline, target_resolution int) *datastruct.Kline {
	defer util.CatchExp("ProcessOldEndKline")
	var pub_kline *datastruct.Kline
	if cur_kline.Resolution != target_resolution {
		cache_kline.Close = cur_kline.Close
		cache_kline.Low = util.MinFloat64(cache_kline.Low, cur_kline.Low)
		cache_kline.High = util.MaxFloat64(cache_kline.High, cur_kline.High)
		cache_kline.Volume += cur_kline.Volume

		pub_kline = datastruct.NewKlineWithKline(cache_kline)
	} else {
		pub_kline = datastruct.NewKlineWithKline(cur_kline)
	}

	cache_kline = datastruct.NewKlineWithKline(pub_kline)
	return pub_kline
}

func ProcessNewStartKline(tmp_kline *datastruct.Kline, cache_kline *datastruct.Kline) {

}

func ProcessCachingKline(tmp_kline *datastruct.Kline, cache_kline *datastruct.Kline) {

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

func TrasOriKlineData(ori_klines []*datastruct.Kline, target_resolution uint32) []*pb.Kline {
	var rst []*pb.Kline

	var cache_kline *datastruct.Kline = nil
	var pub_kline *datastruct.Kline = nil

	for _, tmp_kline := range ori_klines {
		if cache_kline == nil {
			cache_kline = tmp_kline
			ResetFirstKline(cache_kline, target_resolution)
		}

		if datastruct.IsOldKlineEndTime(tmp_kline.Time, int(tmp_kline.Resolution), int64(target_resolution)) {
			pub_kline = ProcessOldEndKline(tmp_kline, cache_kline, int(target_resolution))
			// rst = append(rst, pub_kline)
		} else if datastruct.IsNewKlineStartTime(tmp_kline.Time, int64(target_resolution)) {
			ProcessNewStartKline(tmp_kline, cache_kline)
		} else {
			ProcessCachingKline(tmp_kline, cache_kline)
		}
	}

	if cache_kline.Time != pub_kline.Time {
		// rst = append(rst, cache_kline)
	}

	return rst
}
