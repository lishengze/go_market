// Code generated by goctl. DO NOT EDIT!

package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/shopspring/decimal"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	exchangeTradesFieldNames          = builder.RawFieldNames(&ExchangeTrades{})
	exchangeTradesRows                = strings.Join(exchangeTradesFieldNames, ",")
	exchangeTradesRowsExpectAutoSet   = strings.Join(stringx.Remove(exchangeTradesFieldNames, "`create_time`", "`update_time`"), ",")
	exchangeTradesRowsWithPlaceHolder = strings.Join(stringx.Remove(exchangeTradesFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheBrokerExchangeTradesIdPrefix = "cache:broker:exchangeTrades:id:"
)

type (
	exchangeTradesModel interface {
		Insert(ctx context.Context, session sqlx.Session, data *ExchangeTrades) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*ExchangeTrades, error)
		Update(ctx context.Context, session sqlx.Session, data *ExchangeTrades) (sql.Result, error)
		Delete(ctx context.Context, session sqlx.Session, id int64) error
	}

	defaultExchangeTradesModel struct {
		sqlc.CachedConn
		table string
	}

	ExchangeTrades struct {
		Id int64 `db:"id"` // 主键ID

		Symbol sql.NullString `db:"symbol"` // 交易币种名称，格式：BTC_USDT

		CoinSymbol sql.NullString `db:"coin_symbol"` // 交易币单位

		BaseSymbol sql.NullString `db:"base_symbol"` // 结算单位

		TradeType sql.NullInt64 `db:"trade_type"` // 成交类型

		OfflineTradeId sql.NullString `db:"offline_trade_id"` // 线下交易ID

		Price decimal.Decimal `db:"price"` // 交易价格

		Amount decimal.Decimal `db:"amount"` // 成交量

		BuyTurnover decimal.Decimal `db:"buy_turnover"` // 买入成交额

		BuyFee decimal.Decimal `db:"buy_fee"` // 买入手续费

		SellTurnover decimal.Decimal `db:"sell_turnover"` // 卖出成交额

		SellFee decimal.Decimal `db:"sell_fee"` // 卖出手续费

		Direction sql.NullString `db:"direction"` // 交易方向

		BuyOrderId sql.NullString `db:"buy_order_id"` // 买入订单ID

		SellOrderId sql.NullString `db:"sell_order_id"` // 卖出订单ID

		Status sql.NullInt64 `db:"status"` // 交易状态

		Time sql.NullInt64 `db:"time"` // 成交时间

		UpdateTime time.Time `db:"update_time"` // 交易更新时间

		NeedHedgeAmount decimal.Decimal `db:"need_hedge_amount"` // 待对冲量

		MemberId sql.NullInt64 `db:"member_id"` // 用户ID

		CoinPrice decimal.Decimal `db:"coin_price"` // coin当前价格

		BasePrice decimal.Decimal `db:"base_price"` // base当前价格

		BuyTime sql.NullTime `db:"buy_time"` // 买订单时间

		SellTime sql.NullTime `db:"sell_time"` // 卖订单时间

	}
)

func newExchangeTradesModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultExchangeTradesModel {
	return &defaultExchangeTradesModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`exchange_trades`",
	}
}

func (m *defaultExchangeTradesModel) Insert(ctx context.Context, session sqlx.Session, data *ExchangeTrades) (sql.Result, error) {
	brokerExchangeTradesIdKey := fmt.Sprintf("%s%v", cacheBrokerExchangeTradesIdPrefix, data.Id)
	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, exchangeTradesRowsExpectAutoSet)
		if session != nil {
			return session.ExecCtx(ctx, query, data.Id, data.Symbol, data.CoinSymbol, data.BaseSymbol, data.TradeType, data.OfflineTradeId, data.Price, data.Amount, data.BuyTurnover, data.BuyFee, data.SellTurnover, data.SellFee, data.Direction, data.BuyOrderId, data.SellOrderId, data.Status, data.Time, data.NeedHedgeAmount, data.MemberId, data.CoinPrice, data.BasePrice, data.BuyTime, data.SellTime)
		}
		return conn.ExecCtx(ctx, query, data.Id, data.Symbol, data.CoinSymbol, data.BaseSymbol, data.TradeType, data.OfflineTradeId, data.Price, data.Amount, data.BuyTurnover, data.BuyFee, data.SellTurnover, data.SellFee, data.Direction, data.BuyOrderId, data.SellOrderId, data.Status, data.Time, data.NeedHedgeAmount, data.MemberId, data.CoinPrice, data.BasePrice, data.BuyTime, data.SellTime)
	}, brokerExchangeTradesIdKey)
}

func (m *defaultExchangeTradesModel) FindOne(ctx context.Context, id int64) (*ExchangeTrades, error) {
	brokerExchangeTradesIdKey := fmt.Sprintf("%s%v", cacheBrokerExchangeTradesIdPrefix, id)
	var resp ExchangeTrades
	err := m.QueryRowCtx(ctx, &resp, brokerExchangeTradesIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", exchangeTradesRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultExchangeTradesModel) Update(ctx context.Context, session sqlx.Session, data *ExchangeTrades) (sql.Result, error) {
	brokerExchangeTradesIdKey := fmt.Sprintf("%s%v", cacheBrokerExchangeTradesIdPrefix, data.Id)
	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, exchangeTradesRowsWithPlaceHolder)
		if session != nil {
			return session.ExecCtx(ctx, query, data.Symbol, data.CoinSymbol, data.BaseSymbol, data.TradeType, data.OfflineTradeId, data.Price, data.Amount, data.BuyTurnover, data.BuyFee, data.SellTurnover, data.SellFee, data.Direction, data.BuyOrderId, data.SellOrderId, data.Status, data.Time, data.NeedHedgeAmount, data.MemberId, data.CoinPrice, data.BasePrice, data.BuyTime, data.SellTime, data.Id)
		}
		return conn.ExecCtx(ctx, query, data.Symbol, data.CoinSymbol, data.BaseSymbol, data.TradeType, data.OfflineTradeId, data.Price, data.Amount, data.BuyTurnover, data.BuyFee, data.SellTurnover, data.SellFee, data.Direction, data.BuyOrderId, data.SellOrderId, data.Status, data.Time, data.NeedHedgeAmount, data.MemberId, data.CoinPrice, data.BasePrice, data.BuyTime, data.SellTime, data.Id)
	}, brokerExchangeTradesIdKey)
}

func (m *defaultExchangeTradesModel) Delete(ctx context.Context, session sqlx.Session, id int64) error {
	brokerExchangeTradesIdKey := fmt.Sprintf("%s%v", cacheBrokerExchangeTradesIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		if session != nil {
			return session.ExecCtx(ctx, query, id)
		}
		return conn.ExecCtx(ctx, query, id)
	}, brokerExchangeTradesIdKey)
	return err
}

func (m *defaultExchangeTradesModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheBrokerExchangeTradesIdPrefix, primary)
}
func (m *defaultExchangeTradesModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", exchangeTradesRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultExchangeTradesModel) tableName() string {
	return m.table
}
