package dbserver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"market_server/app/data_manager/rpc/types/pb"
	"market_server/common/config"
	"market_server/common/datastruct"
	"market_server/common/util"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/protobuf/types/known/timestamppb"

	_ "github.com/go-sql-driver/mysql"
)

type DBServer struct {
	kline_cache     *datastruct.KlineCache
	RecvDataChan    *datastruct.DataChannel
	conn            sqlx.SqlConn
	db              *sql.DB
	insert_stmt_map map[string]*sql.Stmt
	tables_         map[string]struct{}
}

func NewDBServer(recvDataChan *datastruct.DataChannel, mysql_config *config.MysqlConfig, cache_config *datastruct.CacheConfig) (*DBServer, error) {

	new_db, err := sql.Open("mysql", mysql_config.Addr)

	if err != nil {
		return nil, err
	}

	return &DBServer{
		RecvDataChan:    recvDataChan,
		db:              new_db,
		insert_stmt_map: make(map[string]*sql.Stmt),
		tables_:         make(map[string]struct{}),
		kline_cache:     datastruct.NewKlineCache(cache_config),
	}, nil
}

func (a *DBServer) StartListenRecvdata() {
	defer util.CatchExp(fmt.Sprintf("DBServer StartListenRecvdata"))

	logx.Info("[S] DBServer start_listen_recvdata")
	go func() {
		for {
			select {
			case new_depth := <-a.RecvDataChan.DepthChannel:
				a.process_depth(new_depth)
			case new_kline := <-a.RecvDataChan.KlineChannel:
				a.process_kline(new_kline)
			case new_trade := <-a.RecvDataChan.TradeChannel:
				a.process_trade(new_trade)
			}
		}
	}()
	logx.Info("[S] DBServer start_receiver Over!")
}

func (d *DBServer) get_insert_stmt(data_type string, symbol string, exchange string) (*sql.Stmt, error) {
	table_name := d.get_table_name(data_type, symbol, exchange)

	var stmt *sql.Stmt
	var ok bool
	var err error
	if stmt, ok = d.insert_stmt_map[table_name]; !ok {
		switch data_type {
		case datastruct.KLINE_TYPE:
			stmt, err = d.db.Prepare(fmt.Sprintf(`INSERT %s (exchange,symbol,time,open,high,low,close,volume,resolution) 
																			values (?,?,?,?,?,?,?,?,?)`, table_name))
		case datastruct.TRADE_TYPE:
			stmt, err = d.db.Prepare(fmt.Sprintf(`INSERT %s (exchange,symbol,time, price, volume) values (?,?,?,?,?)`, table_name))
		case datastruct.DEPTH_TYPE:
			stmt, err = d.db.Prepare(fmt.Sprintf(`INSERT %s (exchange,symbol,time, price, volume) values (?,?,?,?,?)`, table_name))
		}
	}

	if err != nil {
		return nil, err
	}

	d.insert_stmt_map[table_name] = stmt

	return stmt, err
}

func (d *DBServer) get_table_name(data_type string, symbol string, exchange string) string {
	switch data_type {
	case datastruct.KLINE_TYPE:
		return "kline" + "_" + symbol + "_" + exchange
	case datastruct.TRADE_TYPE:
		return "trade" + "_" + symbol + "_" + exchange
	case datastruct.DEPTH_TYPE:
		return "depth" + "_" + symbol + "_" + exchange
	}
	return ""
}

func (d *DBServer) update_table_list() {
	sql_str := "show tables;"
	rows, err := d.db.Query(sql_str)

	if err != nil {
		logx.Errorf("err: %+v", err)
	}

	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)

		if err != nil {
			logx.Error(err.Error())
			return
		}

		for _, col := range values {
			if col != nil {
				table := string(col.([]byte))

				if _, ok := d.tables_[table]; !ok {
					d.tables_[table] = struct{}{}
				}
			}
		}
	}
}

func (d *DBServer) check_table(data_type string, symbol string, exchange string) bool {
	table_name := d.get_table_name(data_type, symbol, exchange)

	if _, ok := d.tables_[table_name]; !ok || len(d.tables_) == 0 {
		d.update_table_list()
	} else {
		return true
	}

	if _, ok := d.tables_[table_name]; !ok {
		return false
	}

	return true
}

func (d *DBServer) create_table(data_type string, symbol string, exchange string) (bool, error) {

	var create_str string
	table_name := d.get_table_name(data_type, symbol, exchange)
	switch data_type {
	case datastruct.KLINE_TYPE:
		create_str = get_kline_create_str(table_name, symbol, exchange)
	case datastruct.TRADE_TYPE:
		create_str = get_trade_create_str(table_name, symbol, exchange)
	case datastruct.DEPTH_TYPE:
		create_str = get_depth_create_str(table_name, symbol, exchange)
	}

	fmt.Printf("create_str: %s \n", create_str)

	_, err := d.db.Exec(create_str)
	if err != nil {
		logx.Error(err.Error())
		return false, err
	}

	d.update_table_list()
	return d.check_table(data_type, symbol, exchange), nil
}

// UnTest
func (d *DBServer) process_kline(kline *datastruct.Kline) error {
	defer util.CatchExp(fmt.Sprintf("DBServer process_kline %s", kline.String()))

	d.kline_cache.UpdateAllKline(kline)

	if kline.IsHistory() {
		return d.store_kline(kline)
	} else {
		// logx.Slowf("[RK] %s", kline.FullString())

		trade := datastruct.NewTradeWithRealTimeKline(kline)
		return d.store_trade(trade)
	}
}

func (d *DBServer) store_kline(kline *datastruct.Kline) error {

	defer util.CatchExp(fmt.Sprintf("store_kline %s", kline.FullString()))

	// logx.Slowf("[HK] %s", kline.FullString())

	if ok := d.check_table(datastruct.KLINE_TYPE, kline.Symbol, kline.Exchange); !ok {
		if ok, err := d.create_table(datastruct.KLINE_TYPE, kline.Symbol, kline.Exchange); !ok {
			logx.Error(err.Error())
			return err
		}
	}

	stmt, err := d.get_insert_stmt(datastruct.KLINE_TYPE, kline.Symbol, kline.Exchange)

	if err != nil {
		logx.Error(err)
		return err
	}

	stmt.Exec(kline.Exchange, kline.Symbol, kline.Time, kline.Open,
		kline.High, kline.Low, kline.Close, kline.Volume, kline.Resolution)

	return err
}

func (d *DBServer) store_trade(trade *datastruct.Trade) error {
	defer util.CatchExp(fmt.Sprintf("store_trade %s", trade.String()))

	// logx.Slowf("[ST] %s", trade.String())

	if ok := d.check_table(datastruct.TRADE_TYPE, trade.Symbol, trade.Exchange); !ok {
		if ok, err := d.create_table(datastruct.TRADE_TYPE, trade.Symbol, trade.Exchange); !ok {
			logx.Error(err.Error())
			return err
		}
	}

	stmt, err := d.get_insert_stmt(datastruct.TRADE_TYPE, trade.Symbol, trade.Exchange)

	if err != nil {
		logx.Error(err)
		return err
	}

	stmt.Exec(trade.Exchange, trade.Symbol, trade.Time, trade.Price, trade.Volume)

	return err
}

func (d *DBServer) process_trade(trade *datastruct.Trade) error {

	// logx.Slowf("[RT] %s", trade.String())

	// if ok := d.check_table(datastruct.TRADE_TYPE, trade.Symbol, trade.Exchange); !ok {
	// 	if ok, err := d.create_table(datastruct.TRADE_TYPE, trade.Symbol, trade.Exchange); !ok {
	// 		logx.Error(err.Error())
	// 		return err
	// 	}
	// }

	// // d.kline_cache.UpdateWithTrade(trade)

	// stmt, err := d.get_insert_stmt(datastruct.TRADE_TYPE, trade.Symbol, trade.Exchange)

	// if err != nil {
	// 	logx.Error(err)
	// 	return err
	// }

	// stmt.Exec(trade.Exchange, trade.Symbol, trade.Time, trade.Price, trade.Volume)

	return nil
}

func (d *DBServer) store_depth(depth *datastruct.DepthQuote) error {
	return nil
}

func (d *DBServer) process_depth(depth *datastruct.DepthQuote) error {
	return nil
}

// UnTest
func (d *DBServer) GetTradesByCount(symbol string, count int) []*datastruct.Trade {
	defer util.CatchExp(fmt.Sprintf(" DBServer GetTradesByCount %s.%d Faled ", symbol, count))

	table_name := d.get_table_name(datastruct.TRADE_TYPE, symbol, datastruct.BCTS_EXCHANGE)
	sql_str := get_lastest_trades_by_count(table_name, count)

	rows, err := d.db.Query(sql_str)
	if err != nil {
		logx.Errorf("db.Query %s err: %+v", sql_str, err)
		return nil
	}

	trades := TransDBTrades(rows)

	return trades
}

// UnTest
func (d *DBServer) GetLastMinuteTrades(symbol string) []*datastruct.Trade {
	defer util.CatchExp(fmt.Sprintf("DBServer GetLastMinuteTrades %s", symbol))
	var rst []*datastruct.Trade = nil

	most_lastest_trades := d.GetTradesByCount(symbol, 1000)

	if most_lastest_trades == nil || len(most_lastest_trades) == 0 {
		return rst
	}

	last_trade := most_lastest_trades[len(most_lastest_trades)-1]
	last_min_nacos := datastruct.GetMiniuteNanos(last_trade.Time)

	for i := len(most_lastest_trades) - 1; i > 0; i-- {
		if datastruct.GetMiniuteNanos(most_lastest_trades[i].Time) == last_min_nacos {
			rst = append(rst, most_lastest_trades[i])
		} else {
			break
		}
	}

	return rst
}

// UnTest
func (d *DBServer) GetDBKlinesByCount(symbol string, resolution int, count int) []*datastruct.Kline {
	defer util.CatchExp(fmt.Sprintf("GetDBKlinesByCount %s, %d, %d", symbol, resolution, count))

	logx.Slowf("GetDBKlinesByCount %s, %d, %d", symbol, resolution, count)

	table_name := d.get_table_name(datastruct.KLINE_TYPE, symbol, datastruct.BCTS_EXCHANGE)
	sql_str := get_kline_sql_str_by_count(table_name, count)

	rows, err := d.db.Query(sql_str)

	if err != nil {
		logx.Errorf("err: %+v", err)
		return nil
	}

	rst := TransDBKlines(rows)

	return rst
}

func (d *DBServer) GetDBKlinesByTime(symbol string, resolution int, start_time int64, end_time int64) []*datastruct.Kline {
	defer util.CatchExp("DBServer GetKlinesByCount")

	table_name := d.get_table_name(datastruct.KLINE_TYPE, symbol, datastruct.BCTS_EXCHANGE)
	sql_str := get_kline_sql_str_by_time(table_name, uint64(start_time), uint64(end_time))

	rows, err := d.db.Query(sql_str)

	if err != nil {
		logx.Errorf("err: %+v", err)
		return nil
	}

	rst := TransDBKlines(rows)

	return rst
}

// UnTest
func (d *DBServer) GetKlinesByCount(symbol string, resolution int, count int) []*datastruct.Kline {
	defer util.CatchExp("DBServer GetKlinesByCount")

	logx.Slowf("GetKlinesByCount  %s, %d, %d", symbol, resolution, count)

	rst := d.kline_cache.GetKlinesByCount(symbol, resolution, count, true)

	if rst == nil {
		logx.Slowf("KlineCache does not have enough data")
		req_count := util.MaxInt(count, d.kline_cache.Config.Count)
		db_klines := d.GetDBKlinesByCount(symbol, resolution, req_count)
		d.kline_cache.InitWithHistKlines(db_klines, symbol, resolution)

		latest_kline := d.kline_cache.GetLatestRealTimeKline(symbol)
		d.kline_cache.UpdateWithKline(latest_kline, resolution)

		logx.Infof("\n--------------------------  KlineCache: ------------------------", d.kline_cache.String(symbol, resolution))
	}

	rst = d.kline_cache.GetKlinesByCount(symbol, resolution, count, false)

	return rst
}

// UnTest
func (d *DBServer) GetKlinesByTime(symbol string, resolution int, start_time int64, end_time int64) []*datastruct.Kline {
	defer util.CatchExp("DBServer GetKlinesByTime")

	rst := d.kline_cache.GetKlinesByTime(symbol, resolution, start_time, end_time, true)

	if rst != nil {
		db_klines := d.GetDBKlinesByTime(symbol, resolution, start_time, end_time)
		d.kline_cache.InitWithHistKlines(db_klines, symbol, resolution)
	}

	rst = d.kline_cache.GetKlinesByTime(symbol, resolution, start_time, end_time, false)

	return rst
}

// UnTest
func (d *DBServer) RequestHistKlineData(ctx context.Context, in *pb.ReqHishKlineInfo) (*pb.HistKlineData, error) {
	defer util.CatchExp("DBServer RequestHistKlineData")

	symbol := in.GetSymbol()
	exchange := in.GetExchange()

	if exchange != datastruct.BCTS_EXCHANGE {
		return nil, fmt.Errorf("only has %s data ", exchange)
	}

	count := in.GetCount()
	start_time := in.GetStartTime()
	end_time := in.GetEndTime()
	frequency := in.GetFrequency()

	if frequency%datastruct.SECS_PER_MIN != 0 {
		return nil, fmt.Errorf("frequency %d is error ", frequency)
	}

	logx.Slowf("ReqHishKlineInfo: %+v", in)

	var ori_klines []*datastruct.Kline = nil

	if count > 0 {
		ori_klines = d.GetKlinesByCount(symbol, int(frequency), int(count))
	} else if start_time > 0 && end_time > 0 && start_time <= end_time {
		ori_klines = d.GetKlinesByTime(symbol, int(frequency), int64(start_time), int64(end_time))
	} else {
		return nil, fmt.Errorf("invalid count %d, start_time: %d, end_time: %d ", count, start_time, end_time)
	}

	rst := &pb.HistKlineData{
		Symbol:    symbol,
		Exchange:  exchange,
		Frequency: frequency,
		Count:     count,
		StartTime: start_time,
		EndTime:   end_time,
	}

	rst.KlineData = TransKlineData(ori_klines)

	var err error = nil
	if rst.KlineData == nil {
		err = fmt.Errorf("empty Kline Data For %s, %d, %d, %d, %d",
			symbol, frequency, count, start_time, end_time)
	}

	return rst, err
}

func (d *DBServer) RequestHistKlineDataBak(ctx context.Context, in *pb.ReqHishKlineInfo) (*pb.HistKlineData, error) {

	rst := pb.HistKlineData{}

	table_name := d.get_table_name(datastruct.KLINE_TYPE, in.GetSymbol(), in.GetExchange())

	var sql_str string
	if in.GetCount() != 0 {
		sql_str = get_kline_sql_str_by_count(table_name, int(in.GetCount()))
	} else if in.GetStartTime() <= 0 || in.GetEndTime() <= 0 || in.GetEndTime() < in.GetStartTime() {
		return nil, errors.New("Time Invalid!")
	} else {
		sql_str = get_kline_sql_str_by_time(table_name, in.GetStartTime(), in.GetEndTime())
	}

	// frequency := in.GetFrequency()

	fmt.Printf("sql_str: %+v \n", sql_str)

	rows, err := d.db.Query(sql_str)

	if err != nil {
		logx.Errorf("err: %+v", err)
		return nil, err
	}

	columns, _ := rows.Columns()

	rst.Count = 0
	rst.Symbol = in.GetSymbol()
	rst.Exchange = in.GetExchange()
	rst.StartTime = in.GetStartTime()
	rst.EndTime = in.GetEndTime()
	rst.Frequency = in.GetFrequency()

	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		rst.Count = rst.Count + 1
		err = rows.Scan(scanArgs...)

		if err != nil {
			logx.Error(err.Error())
			return nil, err
		}

		tmp_kline := &pb.Kline{}
		tmp_kline.Exchange = string(values[0].([]byte))
		tmp_kline.Symbol = string(values[1].([]byte))

		time, _ := strconv.Atoi(string(values[2].([]byte)))
		tmp_kline.Timestamp = &timestamppb.Timestamp{Seconds: int64(time) / datastruct.NANO_PER_SECS, Nanos: int32(time % datastruct.NANO_PER_SECS)}

		tmp_kline.Open = string(values[3].([]byte))
		tmp_kline.High = string(values[4].([]byte))
		tmp_kline.Low = string(values[5].([]byte))
		tmp_kline.Close = string(values[6].([]byte))
		tmp_kline.Volume = string(values[7].([]byte))

		resolution, _ := strconv.Atoi(string(values[8].([]byte)))
		tmp_kline.Resolution = uint32(resolution)

		rst.KlineData = append(rst.KlineData, tmp_kline)

	}

	return &rst, nil
}

func (d *DBServer) GetAllTime(table_name string) []uint64 {
	sql_str := fmt.Sprintf("select time from %s", table_name)
	var rst []uint64

	rows, err := d.db.Query(sql_str)

	if err != nil {
		logx.Errorf("err: %+v", err)
		return rst
	}

	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)

		if err != nil {
			logx.Error(err.Error())
			return rst
		}

		for _, col := range values {
			if col != nil {
				time, _ := strconv.Atoi(string(col.([]byte)))

				rst = append(rst, uint64(time))

			}
		}
	}
	return rst
}

func (d *DBServer) RequestTradeData(ctx context.Context, in *pb.ReqTradeInfo) (*pb.Trade, error) {
	rst := pb.Trade{}
	table_name := d.get_table_name(datastruct.TRADE_TYPE, in.GetSymbol(), in.GetExchange())

	time_list := d.GetAllTime(table_name)

	if len(time_list) > 0 {
		nearest_time := time_list[0]
		minum_time_delta := time_list[0]
		for _, time := range time_list {
			// cur_delta :=  uint64(math.Abs(float64(float64(time)-float64(in.GetTime()))))
			cur_delta := uint64(math.Abs(float64(int64(time) - int64(in.GetTime()))))
			if cur_delta < minum_time_delta {
				nearest_time = time
				minum_time_delta = cur_delta
			}
		}

		logx.Infof("TimeList: %+v \n", time_list)
		logx.Infof("minum_time_delta: %+v, RequestTime: %+v, NearestTime: %+v \n",
			minum_time_delta, in.GetTime(), nearest_time)

		sql_str := get_trade_sql_str(table_name, nearest_time)
		rows, err := d.db.Query(sql_str)
		if err != nil {
			logx.Errorf("err: %+v", err)
			return &rst, err
		}

		columns, _ := rows.Columns()
		scanArgs := make([]interface{}, len(columns))
		values := make([]interface{}, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		for rows.Next() {
			err = rows.Scan(scanArgs...)

			if err != nil {
				logx.Error(err.Error())
				return &rst, err
			}

			rst.Exchange = string(values[0].([]byte))
			rst.Symbol = string(values[1].([]byte))

			time, _ := strconv.Atoi(string(values[2].([]byte)))
			rst.Timestamp = &timestamppb.Timestamp{Seconds: int64(time) / datastruct.NANO_PER_SECS, Nanos: int32(time % datastruct.NANO_PER_SECS)}
			rst.Price = string(values[3].([]byte))
			rst.Volume = string(values[4].([]byte))
		}

		return &rst, nil
	} else {
		return &rst, errors.New("TradeData Empty!")
	}

}

func test_basic(dbServer *DBServer) {
	dbServer.update_table_list()
	fmt.Printf("Original All Tables: %+v: \n\n", dbServer.tables_)

	ok, err := dbServer.create_table(datastruct.KLINE_TYPE, "ETH_USDT", "_bcts_")
	if !ok {
		fmt.Println(err)
		return
	}
	dbServer.update_table_list()
	fmt.Printf("All Tables: %+v: \n\n", dbServer.tables_)

	dbServer.create_table(datastruct.TRADE_TYPE, "ETH_USDT", "_bcts_")
	if !ok {
		fmt.Println(err)
		return
	}
	dbServer.update_table_list()
	fmt.Printf("All Tables: %+v: \n\n", dbServer.tables_)
}

func test_store(dbServer *DBServer) {
	// test_kline := datastruct.GetTestKline()
	// test_kline.Exchange = datastruct.BCTS_EXCHANGE
	// dbServer.process_kline(test_kline)

	test_trade := datastruct.GetTestTrade()
	test_trade.Exchange = datastruct.BCTS_EXCHANGE
	dbServer.store_trade(test_trade)
}

func store_data(dbServer *DBServer, data_count int) {
	for i := 0; i < data_count; i++ {

		test_kline := datastruct.GetTestKline()
		test_kline.Exchange = datastruct.BCTS_EXCHANGE
		dbServer.process_kline(test_kline)

		test_trade := datastruct.GetTestTrade()
		test_trade.Exchange = datastruct.BCTS_EXCHANGE
		dbServer.store_trade(test_trade)

		fmt.Printf("Store data %d \n", i)
		time.Sleep(time.Second * 1)
	}
}

func test_request_trade(dbServer *DBServer) {
	in := &pb.ReqTradeInfo{
		Symbol:   "BTC_USDT",
		Exchange: datastruct.BCTS_EXCHANGE,
		Time:     1654297013967604740,
	}

	ctx_bk := context.Background()
	rst, err := dbServer.RequestTradeData(ctx_bk, in)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Rst: %+v \n", rst)
}

func test_request_histkline(dbServer *DBServer) {
	in := &pb.ReqHishKlineInfo{
		Symbol:    "BTC_USDT",
		Exchange:  datastruct.BCTS_EXCHANGE,
		StartTime: 1654297007842658763,
		EndTime:   1654297013959596689,
		Count:     10,
		Frequency: 60,
	}

	ctx_bk := context.Background()
	rst, err := dbServer.RequestHistKlineData(ctx_bk, in)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Rst: %+v \n", rst)
}

func TestDB() {
	recv_data_chan := datastruct.NewDataChannel()

	mysql_config := config.MysqlConfig{
		Addr: "bcts:bcts@tcp(127.0.0.1:3306)/market",
	}

	fmt.Println(mysql_config)

	dbServer, err := NewDBServer(recv_data_chan, &mysql_config, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// test_basic(dbServer)

	// test_store(dbServer)

	// store_data(dbServer, 100)

	// test_request_trade(dbServer)

	test_request_histkline(dbServer)
}
