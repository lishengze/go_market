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
	memberTransactionsFieldNames          = builder.RawFieldNames(&MemberTransactions{})
	memberTransactionsRows                = strings.Join(memberTransactionsFieldNames, ",")
	memberTransactionsRowsExpectAutoSet   = strings.Join(stringx.Remove(memberTransactionsFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	memberTransactionsRowsWithPlaceHolder = strings.Join(stringx.Remove(memberTransactionsFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheBrokerMemberTransactionsIdPrefix = "cache:broker:memberTransactions:id:"
)

type (
	memberTransactionsModel interface {
		Insert(ctx context.Context, session sqlx.Session, data *MemberTransactions) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*MemberTransactions, error)
		Update(ctx context.Context, session sqlx.Session, data *MemberTransactions) (sql.Result, error)
		Delete(ctx context.Context, session sqlx.Session, id int64) error
	}

	defaultMemberTransactionsModel struct {
		sqlc.CachedConn
		table string
	}

	MemberTransactions struct {
		Id int64 `db:"id"` // 主键ID

		MemberId sql.NullInt64 `db:"member_id"` // 会员ID

		AccountId sql.NullInt64 `db:"account_id"` // 账户ID

		DetailId sql.NullInt64 `db:"detail_id"` // 订单明细ID

		Amount decimal.Decimal `db:"amount"` // 交易金额

		TransactionType sql.NullInt64 `db:"transaction_type"` // 交易类型

		Symbol sql.NullString `db:"symbol"` // 币种名称

		TxId sql.NullString `db:"tx_id"` // 提现交易编号

		Address sql.NullString `db:"address"` // 充值或提现地址、或转账地址

		Fee decimal.Decimal `db:"fee"` // 交易手续费

		Flag sql.NullInt64 `db:"flag"` // 标识位

		RealFee decimal.Decimal `db:"real_fee"` // 实收手续费

		DiscountFee decimal.Decimal `db:"discount_fee"` // 折扣手续费

		PreBalance decimal.Decimal `db:"pre_balance"` // 处理前余额

		Balance decimal.Decimal `db:"balance"` // 处理前余额

		PreFrozenBal decimal.Decimal `db:"pre_frozen_bal"` // 处理前冻结余额

		FrozenBal decimal.Decimal `db:"frozen_bal"` // 处理后冻结余额

		CreateTime time.Time `db:"create_time"` // 创建时间

		CoinPrice decimal.Decimal `db:"coin_price"` // base当前价格

		AveragePrice decimal.Decimal `db:"average_price"` // 当前均价

		Profit decimal.Decimal `db:"profit"` // 收益

	}
)

func newMemberTransactionsModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultMemberTransactionsModel {
	return &defaultMemberTransactionsModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`member_transactions`",
	}
}

func (m *defaultMemberTransactionsModel) Insert(ctx context.Context, session sqlx.Session, data *MemberTransactions) (sql.Result, error) {
	brokerMemberTransactionsIdKey := fmt.Sprintf("%s%v", cacheBrokerMemberTransactionsIdPrefix, data.Id)
	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, memberTransactionsRowsExpectAutoSet)
		if session != nil {
			return session.ExecCtx(ctx, query, data.MemberId, data.AccountId, data.DetailId, data.Amount, data.TransactionType, data.Symbol, data.TxId, data.Address, data.Fee, data.Flag, data.RealFee, data.DiscountFee, data.PreBalance, data.Balance, data.PreFrozenBal, data.FrozenBal, data.CoinPrice, data.AveragePrice, data.Profit)
		}
		return conn.ExecCtx(ctx, query, data.MemberId, data.AccountId, data.DetailId, data.Amount, data.TransactionType, data.Symbol, data.TxId, data.Address, data.Fee, data.Flag, data.RealFee, data.DiscountFee, data.PreBalance, data.Balance, data.PreFrozenBal, data.FrozenBal, data.CoinPrice, data.AveragePrice, data.Profit)
	}, brokerMemberTransactionsIdKey)
}

func (m *defaultMemberTransactionsModel) FindOne(ctx context.Context, id int64) (*MemberTransactions, error) {
	brokerMemberTransactionsIdKey := fmt.Sprintf("%s%v", cacheBrokerMemberTransactionsIdPrefix, id)
	var resp MemberTransactions
	err := m.QueryRowCtx(ctx, &resp, brokerMemberTransactionsIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", memberTransactionsRows, m.table)
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

func (m *defaultMemberTransactionsModel) Update(ctx context.Context, session sqlx.Session, data *MemberTransactions) (sql.Result, error) {
	brokerMemberTransactionsIdKey := fmt.Sprintf("%s%v", cacheBrokerMemberTransactionsIdPrefix, data.Id)
	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, memberTransactionsRowsWithPlaceHolder)
		if session != nil {
			return session.ExecCtx(ctx, query, data.MemberId, data.AccountId, data.DetailId, data.Amount, data.TransactionType, data.Symbol, data.TxId, data.Address, data.Fee, data.Flag, data.RealFee, data.DiscountFee, data.PreBalance, data.Balance, data.PreFrozenBal, data.FrozenBal, data.CoinPrice, data.AveragePrice, data.Profit, data.Id)
		}
		return conn.ExecCtx(ctx, query, data.MemberId, data.AccountId, data.DetailId, data.Amount, data.TransactionType, data.Symbol, data.TxId, data.Address, data.Fee, data.Flag, data.RealFee, data.DiscountFee, data.PreBalance, data.Balance, data.PreFrozenBal, data.FrozenBal, data.CoinPrice, data.AveragePrice, data.Profit, data.Id)
	}, brokerMemberTransactionsIdKey)
}

func (m *defaultMemberTransactionsModel) Delete(ctx context.Context, session sqlx.Session, id int64) error {
	brokerMemberTransactionsIdKey := fmt.Sprintf("%s%v", cacheBrokerMemberTransactionsIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		if session != nil {
			return session.ExecCtx(ctx, query, id)
		}
		return conn.ExecCtx(ctx, query, id)
	}, brokerMemberTransactionsIdKey)
	return err
}

func (m *defaultMemberTransactionsModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheBrokerMemberTransactionsIdPrefix, primary)
}
func (m *defaultMemberTransactionsModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", memberTransactionsRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultMemberTransactionsModel) tableName() string {
	return m.table
}
