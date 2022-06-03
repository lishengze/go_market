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
	result, err := d.db.Query(sql_str)

	if err != nil {
		logx.Errorf("err: %+v", err)
	}

	// while(result->next())
	// {
	// 	string table = result->getString(1);

	// 	// LOG_INFO("table: " + table);

	// 	table_set_.emplace(table);
	// }
}

func (d *DBServer) check_table(data_type string, symbol string, exchange string) bool {

	return false
}

func (d *DBServer) create_table(data_type string, symbol string, exchange string) bool {
	return false
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
	result := ""
	return result
}

func (d *DBServer) store_kline(kline *datastruct.Kline) error {
	return nil
}

func (d *DBServer) store_trade(trade *datastruct.Trade) error {
	return nil
}

func (d *DBServer) store_depth(depth *datastruct.DepthQuote) error {
	return nil
}

func (s *DBServer) RequestHistKlineData(ctx context.Context, in *pb.ReqHishKlineInfo) (*pb.HistKlineData, error) {

	rst := pb.HistKlineData{}
	return &rst, nil
}

func (s *DBServer) RequestTradeData(ctx context.Context, in *pb.ReqTradeInfo) (*pb.Trade, error) {

	rst := pb.Trade{}

	return &rst, nil
}
