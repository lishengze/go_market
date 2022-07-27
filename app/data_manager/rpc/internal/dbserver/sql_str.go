package dbserver

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func get_kline_create_str(table_name string, symbol string, exchange string) string {
	result := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (exchange VARCHAR(32), 
							symbol VARCHAR(64), time BIGINT PRIMARY KEY, 
						   open DECIMAL(32, 8), high DECIMAL(32, 8), low DECIMAL(32, 8),
						   close DECIMAL(32, 8), volume DECIMAL(32, 8), resolution BIGINT) 
						DEFAULT CHARSET utf8`,
		table_name)

	return result
}

func get_trade_create_str(table_name string, symbol string, exchange string) string {
	result := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (exchange VARCHAR(32), 
						   symbol VARCHAR(64), time BIGINT PRIMARY KEY, 
						   price DECIMAL(32, 8), volume DECIMAL(32, 8))  DEFAULT CHARSET utf8`,
		table_name)

	return result
}

func get_depth_create_str(table_name string, symbol string, exchange string) string {
	result := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (exchange VARCHAR(32), 
						symbol VARCHAR(64), time BIGINT PRIMARY KEY, 
						price DECIMAL(32, 8), volume DECIMAL(32, 8))  DEFAULT CHARSET utf8`,
		table_name)

	return result
}

func get_kline_sql_str_by_time(table_name string, start_time uint64, end_time uint64) string {
	return fmt.Sprintf(`select * from %s where time<=%d and time>=%d;`, table_name, end_time, start_time)
}

func get_kline_sql_str_by_count(table_name string, data_count int) string {
	return fmt.Sprintf(`select * from %s order by time desc limit %d;`, table_name, data_count)
}

func get_trade_sql_str(table_name string, time uint64) string {
	return fmt.Sprintf(`select * from %s where time=%d;`, table_name, time)
}

func get_lastest_trades_by_count(table_name string, count int) string {
	return fmt.Sprintf("select * from %s order by time desc limit %d", table_name, count)
}
