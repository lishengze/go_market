package dbserver

import (
	"context"
	"market_server/app/dataManager/rpc/internal/config"
	"market_server/app/dataManager/rpc/types/pb"
	"market_server/common/datastruct"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type DBServer struct {
	RecvDataChan *datastruct.DataChannel
	conn         sqlx.SqlConn
}

func NewDBServer(recvDataChan *datastruct.DataChannel, mysql_config config.MysqlConfig) *DBServer {

	return &DBServer{
		RecvDataChan: recvDataChan,
		conn:         sqlx.NewMysql(mysql_config.Addr),
	}
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

func (d *DBServer) get_table_name(data_type string, symbol string, exchange string) string {
	switch data_type {
	case "kline":
		return "kline" + "_" + symbol + "_" + exchange
	case "trade":
		return "trade" + "_" + symbol + "_" + exchange
	case "depth":
		return "depth" + "_" + symbol + "_" + exchange
	}
	return ""
}

func (d *DBServer) check_table(data_type string, symbol string, exchange string) bool {
	return false
}

func (d *DBServer) create_table(data_type string, symbol string, exchange string) bool {
	return false
}

func (d *DBServer) get_kline_create_str(symbol string, exchange string) string {
	rst := ""
	return rst
}

func (d *DBServer) get_trade_create_str(symbol string, exchange string) string {
	rst := ""
	return rst
}

func (d *DBServer) get_depth_create_str(symbol string, exchange string) string {
	rst := ""
	return rst
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
