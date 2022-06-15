package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/mr"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	tradeFieldNames          = builder.RawFieldNames(&Trade{})
	tradeRows                = strings.Join(tradeFieldNames, ",")
	tradeRowsExpectAutoSet   = strings.Join(stringx.Remove(tradeFieldNames, "`create_time`", "`update_time`"), ",")
	tradeRowsWithPlaceHolder = strings.Join(stringx.Remove(tradeFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheTradeIdPrefix               = "cache:trade:id:"
	cacheTradeOrderIdExTradeIdPrefix = "cache:trade:orderId:exTradeId:"
)

type (
	TradeModel interface {
		Insert(data *Trade) (sql.Result, error)
		TxInsert(data *Trade) func() (interface{}, error)
		BulkInsert(list []*Trade) error
		TxBulkInsert(list []*Trade) func() (interface{}, error)
		FindOne(id string) (*Trade, error)
		FindMany(q Query_) ([]*Trade, error)
		FindPks(q Query_) ([]string, error)
		FindManyByPks(pks []string) ([]*Trade, error)

		FindOneByOrderIdExTradeId(orderId string, exTradeId string) (*Trade, error)
		Update(data *Trade, update func()) error
		TxUpdate(data *Trade, update func()) func() (interface{}, error)
		Delete(id string) error
		TxDelete(id string) func() (interface{}, error)
	}

	defaultTradeModel struct {
		sqlc.CachedConn
		table string
	}

	Trade struct {
		Id          string    `db:"id"`           // id
		OrderId     string    `db:"order_id"`     // order_id
		ExTradeId   string    `db:"ex_trade_id"`  // ex_trade_id
		Exchange    string    `db:"exchange"`     // exchange
		StdSymbol   string    `db:"std_symbol"`   // std_symbol
		Liquidity   string    `db:"liquidity"`    // liquidity
		Side        string    `db:"side"`         // side
		Volume      string    `db:"volume"`       // volume
		Price       string    `db:"price"`        // price
		Fee         string    `db:"fee"`          // fee
		FeeCurrency string    `db:"fee_currency"` // fee_currency
		TradeTime   time.Time `db:"trade_time"`   // trade_time
		CreateTime  time.Time `db:"create_time"`  // 创建时间
		UpdateTime  time.Time `db:"update_time"`  // 更新时间
	}
)

func NewTradeModel(conn sqlx.SqlConn, c cache.CacheConf) TradeModel {
	return &defaultTradeModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`trade`",
	}
}

func (m *defaultTradeModel) Insert(data *Trade) (sql.Result, error) {
	tradeIdKey := fmt.Sprintf("%s%v", cacheTradeIdPrefix, data.Id)
	tradeOrderIdExTradeIdKey := fmt.Sprintf("%s%v:%v", cacheTradeOrderIdExTradeIdPrefix, data.OrderId, data.ExTradeId)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, tradeRowsExpectAutoSet)
		return conn.Exec(query, data.Id, data.OrderId, data.ExTradeId, data.Exchange, data.StdSymbol, data.Liquidity, data.Side, data.Volume, data.Price, data.Fee, data.FeeCurrency, data.TradeTime)
	}, tradeIdKey, tradeOrderIdExTradeIdKey)
	return ret, err
}

func (m *defaultTradeModel) TxInsert(data *Trade) func() (interface{}, error) {
	keys := make([]string, 0)
	args := []interface{}{data.Id, data.OrderId, data.ExTradeId, data.Exchange, data.StdSymbol, data.Liquidity, data.Side, data.Volume, data.Price, data.Fee, data.FeeCurrency, data.TradeTime}
	insertSql := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, tradeRowsExpectAutoSet)
	tradeIdKey := fmt.Sprintf("%s%v", cacheTradeIdPrefix, data.Id)
	tradeOrderIdExTradeIdKey := fmt.Sprintf("%s%v:%v", cacheTradeOrderIdExTradeIdPrefix, data.OrderId, data.ExTradeId)
	keys = append(keys, tradeIdKey, tradeOrderIdExTradeIdKey)
	return func() (interface{}, error) {
		return &txPrepare_{
			sql:  insertSql,
			args: args,
			keys: keys,
		}, nil
	}
}

func (m *defaultTradeModel) BulkInsert(trades []*Trade) error {
	if len(trades) == 0 {
		return nil
	}

	var (
		insertSql    = fmt.Sprintf("insert into %s (%s) values ", m.table, tradeRowsExpectAutoSet)
		args         = make([]interface{}, 0)
		keys         = make([]string, 0)
		placeHolders = make([]string, 0)
	)
	for _, data := range trades {
		tradeIdKey := fmt.Sprintf("%s%v", cacheTradeIdPrefix, data.Id)
		tradeOrderIdExTradeIdKey := fmt.Sprintf("%s%v:%v", cacheTradeOrderIdExTradeIdPrefix, data.OrderId, data.ExTradeId)
		keys = append(keys, tradeIdKey, tradeOrderIdExTradeIdKey)
		args = append(args, data.Id, data.OrderId, data.ExTradeId, data.Exchange, data.StdSymbol, data.Liquidity, data.Side, data.Volume, data.Price, data.Fee, data.FeeCurrency, data.TradeTime)
		placeHolders = append(placeHolders, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	}

	insertSql += strings.Join(placeHolders, ",")
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.Exec(insertSql, args...)
	}, keys...)

	return err
}

func (m *defaultTradeModel) TxBulkInsert(trades []*Trade) func() (interface{}, error) {
	if len(trades) == 0 {
		return func() (interface{}, error) {
			return newEmptyTxPrepare_(), nil
		}
	}
	var (
		insertSql    = fmt.Sprintf("insert into %s (%s) values ", m.table, tradeRowsExpectAutoSet)
		args         = make([]interface{}, 0)
		keys         = make([]string, 0)
		placeHolders = make([]string, 0)
	)
	for _, data := range trades {
		tradeIdKey := fmt.Sprintf("%s%v", cacheTradeIdPrefix, data.Id)
		tradeOrderIdExTradeIdKey := fmt.Sprintf("%s%v:%v", cacheTradeOrderIdExTradeIdPrefix, data.OrderId, data.ExTradeId)
		keys = append(keys, tradeIdKey, tradeOrderIdExTradeIdKey)
		args = append(args, data.Id, data.OrderId, data.ExTradeId, data.Exchange, data.StdSymbol, data.Liquidity, data.Side, data.Volume, data.Price, data.Fee, data.FeeCurrency, data.TradeTime)
		placeHolders = append(placeHolders, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	}

	insertSql += strings.Join(placeHolders, ",")

	return func() (interface{}, error) {
		return &txPrepare_{
			sql:  insertSql,
			args: args,
			keys: keys,
		}, nil
	}
}

func (m *defaultTradeModel) FindOne(id string) (*Trade, error) {
	tradeIdKey := fmt.Sprintf("%s%v", cacheTradeIdPrefix, id)
	var resp Trade
	err := m.QueryRow(&resp, tradeIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", tradeRows, m.table)
		return conn.QueryRow(v, query, id)
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

func (m *defaultTradeModel) FindMany(q Query_) ([]*Trade, error) {
	pks, err := m.FindPks(q)
	if err != nil {
		return nil, err
	}
	return m.FindManyByPks(pks)
}

func (m *defaultTradeModel) FindPks(q Query_) ([]string, error) {
	var (
		pks = make([]string, 0)
	)
	findPksSql, countSql, args := q.FindPksSql("`id`", m.table)
	err := mr.Finish(func() error {
		return m.QueryRowsNoCache(&pks, findPksSql, args...)
	}, func() error {
		if !q.HasPaginator() {
			return nil
		}
		if q.Size() <= 0 || q.Page() <= 0 {
			return fmt.Errorf("size or page must > 0. ")
		}
		var c = counter_{}
		err := m.QueryRowNoCache(&c.count, countSql, args...)
		q.SetCount(c.count)
		return err
	})
	return pks, err
}

func (m *defaultTradeModel) FindManyByPks(pks []string) ([]*Trade, error) {
	// 使用MapReduce 处理并发
	r, err := mr.MapReduce(func(source chan<- interface{}) {
		for index, pk := range pks {
			source <- []interface{}{index, pk}
		}
	}, func(item interface{}, writer mr.Writer, cancel func(error)) {
		indexAndPk := item.([]interface{})
		Trade, err := m.FindOne(indexAndPk[1].(string))
		if err != nil {
			cancel(err)
			return
		}
		writer.Write([]interface{}{indexAndPk[0], Trade})
	}, func(pipe <-chan interface{}, writer mr.Writer, cancel func(error)) {
		var trades = make([]*Trade, len(pks))
		for content := range pipe {
			trades[content.([]interface{})[0].(int)] = content.([]interface{})[1].(*Trade)
		}
		writer.Write(trades)
	}, mr.WithWorkers(len(pks)))
	if err != nil {
		return nil, err
	}
	return r.([]*Trade), nil
}

func (m *defaultTradeModel) FindOneByOrderIdExTradeId(orderId string, exTradeId string) (*Trade, error) {
	tradeOrderIdExTradeIdKey := fmt.Sprintf("%s%v:%v", cacheTradeOrderIdExTradeIdPrefix, orderId, exTradeId)
	var resp Trade
	err := m.QueryRowIndex(&resp, tradeOrderIdExTradeIdKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `order_id` = ? and `ex_trade_id` = ? limit 1", tradeRows, m.table)
		if err := conn.QueryRow(&resp, query, orderId, exTradeId); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultTradeModel) Update(data *Trade, update func()) error {
	tradeIdKey := fmt.Sprintf("%s%v", cacheTradeIdPrefix, data.Id)
	tradeOrderIdExTradeIdKey := fmt.Sprintf("%s%v:%v", cacheTradeOrderIdExTradeIdPrefix, data.OrderId, data.ExTradeId)
	update()
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, tradeRowsWithPlaceHolder)
		return conn.Exec(query, data.OrderId, data.ExTradeId, data.Exchange, data.StdSymbol, data.Liquidity, data.Side, data.Volume, data.Price, data.Fee, data.FeeCurrency, data.TradeTime, data.Id)
	}, tradeIdKey, tradeOrderIdExTradeIdKey)
	return err
}

func (m *defaultTradeModel) TxUpdate(data *Trade, update func()) func() (interface{}, error) {
	keys := make([]string, 0)
	insertSql := fmt.Sprintf("update %s set %s where `id` = ?", m.table, tradeRowsWithPlaceHolder)
	tradeIdKey := fmt.Sprintf("%s%v", cacheTradeIdPrefix, data.Id)
	tradeOrderIdExTradeIdKey := fmt.Sprintf("%s%v:%v", cacheTradeOrderIdExTradeIdPrefix, data.OrderId, data.ExTradeId)
	keys = append(keys, tradeIdKey, tradeOrderIdExTradeIdKey)
	update()
	args := []interface{}{data.OrderId, data.ExTradeId, data.Exchange, data.StdSymbol, data.Liquidity, data.Side, data.Volume, data.Price, data.Fee, data.FeeCurrency, data.TradeTime, data.Id}
	return func() (interface{}, error) {
		return &txPrepare_{
			sql:  insertSql,
			args: args,
			keys: keys,
		}, nil
	}
}

func (m *defaultTradeModel) Delete(id string) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	tradeOrderIdExTradeIdKey := fmt.Sprintf("%s%v:%v", cacheTradeOrderIdExTradeIdPrefix, data.OrderId, data.ExTradeId)
	tradeIdKey := fmt.Sprintf("%s%v", cacheTradeIdPrefix, id)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, tradeIdKey, tradeOrderIdExTradeIdKey)
	return err
}

func (m *defaultTradeModel) TxDelete(id string) func() (interface{}, error) {
	data, err := m.FindOne(id)
	if err != nil {
		return func() (interface{}, error) {
			return nil, err
		}
	}
	tradeOrderIdExTradeIdKey := fmt.Sprintf("%s%v:%v", cacheTradeOrderIdExTradeIdPrefix, data.OrderId, data.ExTradeId)
	tradeIdKey := fmt.Sprintf("%s%v", cacheTradeIdPrefix, id)
	deleteSql := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	args := []interface{}{id}
	keys := make([]string, 0)
	keys = append(keys, tradeIdKey, tradeOrderIdExTradeIdKey)

	return func() (interface{}, error) {
		return &txPrepare_{
			sql:  deleteSql,
			args: args,
			keys: keys,
		}, nil
	}
}

func (m *defaultTradeModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheTradeIdPrefix, primary)
}

func (m *defaultTradeModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", tradeRows, m.table)
	return conn.QueryRow(v, query, primary)
}
