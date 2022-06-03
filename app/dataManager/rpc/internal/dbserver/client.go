package dbserver

import (
	"context"
	"database/sql"
	"fmt"

	"market_server/app/dataManager/rpc/internal/config"
	"market_server/app/dataManager/rpc/types/pb"
	"market_server/common/datastruct"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

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

func (a *DBServer) start_listen_recvdata() {
	logx.Info("Aggregator start_listen_recvdata")
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
	logx.Info("Aggregator start_receiver Over!")
}

func (d *DBServer) get_insert_stmt(data_type string, symbol string, exchange string) (*sql.Stmt, error) {
	table_name := d.get_table_name(data_type, symbol, exchange)

	var stmt *sql.Stmt
	var ok bool
	var err error
	if stmt, ok = d.insert_stmt_map[table_name]; ok == false {
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
		d.insert_stmt_map[table_name] = stmt
	}

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

				if _, ok := d.tables_[table]; ok == false {
					d.tables_[table] = struct{}{}
				}
			}
		}
	}
}

func (d *DBServer) check_table(data_type string, symbol string, exchange string) bool {
	table_name := d.get_table_name(data_type, symbol, exchange)

	if _, ok := d.tables_[table_name]; ok == false || len(d.tables_) == 0 {
		d.update_table_list()
	} else {
		return true
	}

	if _, ok := d.tables_[table_name]; ok == false {
		return false
	}

	return true
}

func (d *DBServer) create_table(data_type string, symbol string, exchange string) (bool, error) {

	var create_str string
	switch data_type {
	case datastruct.KLINE_TYPE:
		create_str = d.get_kline_create_str(symbol, exchange)
	case datastruct.TRADE_TYPE:
		create_str = d.get_trade_create_str(symbol, exchange)
	case datastruct.DEPTH_TYPE:
		create_str = d.get_depth_create_str(symbol, exchange)
	}

	_, err := d.db.Exec(create_str)
	if err != nil {
		return false, err
	}

	d.update_table_list()
	return d.check_table(data_type, symbol, exchange), nil
}

func (d *DBServer) get_kline_create_str(symbol string, exchange string) string {
	result := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s exchange VARCHAR(32), 
							symbol VARCHAR(64), time BIGINT PRIMARY KEY, 
						   open DECIMAL(32, 8), high DECIMAL(32, 8), low DECIMAL(32, 8),
						   close DECIMAL(32, 8), volume DECIMAL(32, 8)), 
						   resolution BIGINT, DEFAULT CHARSET utf8`,
		d.get_table_name(datastruct.KLINE_TYPE, symbol, exchange))

	return result
}

func (d *DBServer) get_trade_create_str(symbol string, exchange string) string {
	result := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s exchange VARCHAR(32), 
						   symbol VARCHAR(64), time BIGINT PRIMARY KEY, 
						   price DECIMAL(32, 8), volume DECIMAL(32, 8)),  DEFAULT CHARSET utf8`,
		d.get_table_name(datastruct.TRADE_TYPE, symbol, exchange))

	return result
}

func (d *DBServer) get_depth_create_str(symbol string, exchange string) string {
	result := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s exchange VARCHAR(32), 
						symbol VARCHAR(64), time BIGINT PRIMARY KEY, 
						price DECIMAL(32, 8), volume DECIMAL(32, 8)),  DEFAULT CHARSET utf8`,
		d.get_table_name(datastruct.DEPTH_TYPE, symbol, exchange))

	return result
}

func (d *DBServer) store_kline(kline *datastruct.Kline) error {
	if ok := d.check_table(datastruct.KLINE_TYPE, kline.Symbol, kline.Exchange); ok == false {
		if ok, err := d.create_table(datastruct.KLINE_TYPE, kline.Symbol, kline.Exchange); ok == false {
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
	if ok := d.check_table(datastruct.KLINE_TYPE, trade.Symbol, trade.Exchange); ok == false {
		if ok, err := d.create_table(datastruct.TRADE_TYPE, trade.Symbol, trade.Exchange); ok == false {
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

func (s *DBServer) RequestHistKlineData(ctx context.Context, in *pb.ReqHishKlineInfo) (*pb.HistKlineData, error) {

	rst := pb.HistKlineData{}

	// string sql_str = get_kline_sql_str(req_kline_info.exchange, req_kline_info.symbol, req_kline_info.start_time, req_kline_info.end_time);

	return &rst, nil
}

func (s *DBServer) RequestTradeData(ctx context.Context, in *pb.ReqTradeInfo) (*pb.Trade, error) {

	rst := pb.Trade{}

	return &rst, nil
}
