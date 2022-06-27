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

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/protobuf/types/known/timestamppb"

	_ "github.com/go-sql-driver/mysql"
)

type DBServer struct {
	RecvDataChan    *datastruct.DataChannel
	conn            sqlx.SqlConn
	db              *sql.DB
	insert_stmt_map map[string]*sql.Stmt
	tables_         map[string]struct{}
}

func NewDBServer(recvDataChan *datastruct.DataChannel, mysql_config config.MysqlConfig) (*DBServer, error) {

	new_db, err := sql.Open("mysql", mysql_config.Addr)

	if err != nil {
		return nil, err
	}

	return &DBServer{
		RecvDataChan:    recvDataChan,
		db:              new_db,
		insert_stmt_map: make(map[string]*sql.Stmt),
		tables_:         make(map[string]struct{}),
	}, nil
}

func (a *DBServer) StartListenRecvdata() {
	logx.Info("[S] DBServer start_listen_recvdata")
	go func() {
		for {
			select {
			case new_depth := <-a.RecvDataChan.DepthChannel:
				a.store_depth(new_depth)
			case new_kline := <-a.RecvDataChan.KlineChannel:
				a.store_kline(new_kline)
			case new_trade := <-a.RecvDataChan.TradeChannel:
				a.store_trade(new_trade)
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

func (d *DBServer) store_kline(kline *datastruct.Kline) error {
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

func (d *DBServer) store_depth(depth *datastruct.DepthQuote) error {
	return nil
}

func (d *DBServer) RequestHistKlineData(ctx context.Context, in *pb.ReqHishKlineInfo) (*pb.HistKlineData, error) {

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

	fmt.Printf("sql_str: %+v \n", sql_str)

	rows, err := d.db.Query(sql_str)

	if err != nil {
		logx.Errorf("err: %+v", err)
		return nil, err
	}

	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	rst.Count = 0
	rst.Symbol = in.GetSymbol()
	rst.Exchange = in.GetExchange()
	rst.StartTime = in.GetStartTime()
	rst.EndTime = in.GetEndTime()
	rst.Frequency = in.GetFrequency()

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
	// dbServer.store_kline(test_kline)

	test_trade := datastruct.GetTestTrade()
	test_trade.Exchange = datastruct.BCTS_EXCHANGE
	dbServer.store_trade(test_trade)
}

func store_data(dbServer *DBServer, data_count int) {
	for i := 0; i < data_count; i++ {

		test_kline := datastruct.GetTestKline()
		test_kline.Exchange = datastruct.BCTS_EXCHANGE
		dbServer.store_kline(test_kline)

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

	dbServer, err := NewDBServer(recv_data_chan, mysql_config)
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
