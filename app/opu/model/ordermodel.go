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
	orderFieldNames          = builder.RawFieldNames(&Order{})
	orderRows                = strings.Join(orderFieldNames, ",")
	orderRowsExpectAutoSet   = strings.Join(stringx.Remove(orderFieldNames, "`create_time`", "`update_time`"), ",")
	orderRowsWithPlaceHolder = strings.Join(stringx.Remove(orderFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheOrderIdPrefix                     = "cache:order:id:"
	cacheOrderAccountIdClientOrderIdPrefix = "cache:order:accountId:clientOrderId:"
)

type (
	OrderModel interface {
		Insert(data *Order) (sql.Result, error)
		TxInsert(data *Order) func() (interface{}, error)
		BulkInsert(list []*Order) error
		TxBulkInsert(list []*Order) func() (interface{}, error)
		FindOne(id string) (*Order, error)
		FindMany(q Query_) ([]*Order, error)
		FindPks(q Query_) ([]string, error)
		FindManyByPks(pks []string) ([]*Order, error)

		FindOneByAccountIdClientOrderId(accountId string, clientOrderId string) (*Order, error)
		Update(data *Order, update func()) error
		TxUpdate(data *Order, update func()) func() (interface{}, error)
		Delete(id string) error
		TxDelete(id string) func() (interface{}, error)
	}

	defaultOrderModel struct {
		sqlc.CachedConn
		table string
	}

	Order struct {
		Id            string    `db:"id"`              // id, 同时是报给交易所的 order id
		AccountId     string    `db:"account_id"`      // account_id
		AccountAlias  string    `db:"account_alias"`   // account_alias
		ClientOrderId string    `db:"client_order_id"` // client_order_id
		ExOrderId     string    `db:"ex_order_id"`     // ex_order_id
		ApiType       string    `db:"api_type"`        // api_type
		Side          string    `db:"side"`            // side
		Status        string    `db:"status"`          // status
		Volume        string    `db:"volume"`          // volume
		FilledVolume  string    `db:"filled_volume"`   // filled_volume
		Price         string    `db:"price"`           // price
		Tp            string    `db:"tp"`              // type
		StdSymbol     string    `db:"std_symbol"`      // std_symbol
		ExSymbol      string    `db:"ex_symbol"`       // ex_symbol
		Exchange      string    `db:"exchange"`        // exchange
		RejectReason  string    `db:"reject_reason"`   // reject_reason
		CancelFlag    string    `db:"cancel_flag"`     // cancel_flag,表示客户是否下达撤单指令 UNSET|SET
		CreateTime    time.Time `db:"create_time"`     // 创建时间
		UpdateTime    time.Time `db:"update_time"`     // 更新时间
	}
)

func NewOrderModel(conn sqlx.SqlConn, c cache.CacheConf) OrderModel {
	return &defaultOrderModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`order`",
	}
}

func (m *defaultOrderModel) Insert(data *Order) (sql.Result, error) {
	orderIdKey := fmt.Sprintf("%s%v", cacheOrderIdPrefix, data.Id)
	orderAccountIdClientOrderIdKey := fmt.Sprintf("%s%v:%v", cacheOrderAccountIdClientOrderIdPrefix, data.AccountId, data.ClientOrderId)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, orderRowsExpectAutoSet)
		return conn.Exec(query, data.Id, data.AccountId, data.AccountAlias, data.ClientOrderId, data.ExOrderId, data.ApiType, data.Side, data.Status, data.Volume, data.FilledVolume, data.Price, data.Tp, data.StdSymbol, data.ExSymbol, data.Exchange, data.RejectReason, data.CancelFlag)
	}, orderIdKey, orderAccountIdClientOrderIdKey)
	return ret, err
}

func (m *defaultOrderModel) TxInsert(data *Order) func() (interface{}, error) {
	keys := make([]string, 0)
	args := []interface{}{data.Id, data.AccountId, data.AccountAlias, data.ClientOrderId, data.ExOrderId, data.ApiType, data.Side, data.Status, data.Volume, data.FilledVolume, data.Price, data.Tp, data.StdSymbol, data.ExSymbol, data.Exchange, data.RejectReason, data.CancelFlag}
	insertSql := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, orderRowsExpectAutoSet)
	orderIdKey := fmt.Sprintf("%s%v", cacheOrderIdPrefix, data.Id)
	orderAccountIdClientOrderIdKey := fmt.Sprintf("%s%v:%v", cacheOrderAccountIdClientOrderIdPrefix, data.AccountId, data.ClientOrderId)
	keys = append(keys, orderIdKey, orderAccountIdClientOrderIdKey)
	return func() (interface{}, error) {
		return &txPrepare_{
			sql:  insertSql,
			args: args,
			keys: keys,
		}, nil
	}
}

func (m *defaultOrderModel) BulkInsert(orders []*Order) error {
	if len(orders) == 0 {
		return nil
	}

	var (
		insertSql    = fmt.Sprintf("insert into %s (%s) values ", m.table, orderRowsExpectAutoSet)
		args         = make([]interface{}, 0)
		keys         = make([]string, 0)
		placeHolders = make([]string, 0)
	)
	for _, data := range orders {
		orderIdKey := fmt.Sprintf("%s%v", cacheOrderIdPrefix, data.Id)
		orderAccountIdClientOrderIdKey := fmt.Sprintf("%s%v:%v", cacheOrderAccountIdClientOrderIdPrefix, data.AccountId, data.ClientOrderId)
		keys = append(keys, orderIdKey, orderAccountIdClientOrderIdKey)
		args = append(args, data.Id, data.AccountId, data.AccountAlias, data.ClientOrderId, data.ExOrderId, data.ApiType, data.Side, data.Status, data.Volume, data.FilledVolume, data.Price, data.Tp, data.StdSymbol, data.ExSymbol, data.Exchange, data.RejectReason, data.CancelFlag)
		placeHolders = append(placeHolders, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	}

	insertSql += strings.Join(placeHolders, ",")
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.Exec(insertSql, args...)
	}, keys...)

	return err
}

func (m *defaultOrderModel) TxBulkInsert(orders []*Order) func() (interface{}, error) {
	if len(orders) == 0 {
		return func() (interface{}, error) {
			return newEmptyTxPrepare_(), nil
		}
	}
	var (
		insertSql    = fmt.Sprintf("insert into %s (%s) values ", m.table, orderRowsExpectAutoSet)
		args         = make([]interface{}, 0)
		keys         = make([]string, 0)
		placeHolders = make([]string, 0)
	)
	for _, data := range orders {
		orderIdKey := fmt.Sprintf("%s%v", cacheOrderIdPrefix, data.Id)
		orderAccountIdClientOrderIdKey := fmt.Sprintf("%s%v:%v", cacheOrderAccountIdClientOrderIdPrefix, data.AccountId, data.ClientOrderId)
		keys = append(keys, orderIdKey, orderAccountIdClientOrderIdKey)
		args = append(args, data.Id, data.AccountId, data.AccountAlias, data.ClientOrderId, data.ExOrderId, data.ApiType, data.Side, data.Status, data.Volume, data.FilledVolume, data.Price, data.Tp, data.StdSymbol, data.ExSymbol, data.Exchange, data.RejectReason, data.CancelFlag)
		placeHolders = append(placeHolders, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
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

func (m *defaultOrderModel) FindOne(id string) (*Order, error) {
	orderIdKey := fmt.Sprintf("%s%v", cacheOrderIdPrefix, id)
	var resp Order
	err := m.QueryRow(&resp, orderIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", orderRows, m.table)
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

func (m *defaultOrderModel) FindMany(q Query_) ([]*Order, error) {
	pks, err := m.FindPks(q)
	if err != nil {
		return nil, err
	}
	return m.FindManyByPks(pks)
}

func (m *defaultOrderModel) FindPks(q Query_) ([]string, error) {
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

func (m *defaultOrderModel) FindManyByPks(pks []string) ([]*Order, error) {
	// 使用MapReduce 处理并发
	r, err := mr.MapReduce(func(source chan<- interface{}) {
		for index, pk := range pks {
			source <- []interface{}{index, pk}
		}
	}, func(item interface{}, writer mr.Writer, cancel func(error)) {
		indexAndPk := item.([]interface{})
		Order, err := m.FindOne(indexAndPk[1].(string))
		if err != nil {
			cancel(err)
			return
		}
		writer.Write([]interface{}{indexAndPk[0], Order})
	}, func(pipe <-chan interface{}, writer mr.Writer, cancel func(error)) {
		var orders = make([]*Order, len(pks))
		for content := range pipe {
			orders[content.([]interface{})[0].(int)] = content.([]interface{})[1].(*Order)
		}
		writer.Write(orders)
	}, mr.WithWorkers(len(pks)))
	if err != nil {
		return nil, err
	}
	return r.([]*Order), nil
}

func (m *defaultOrderModel) FindOneByAccountIdClientOrderId(accountId string, clientOrderId string) (*Order, error) {
	orderAccountIdClientOrderIdKey := fmt.Sprintf("%s%v:%v", cacheOrderAccountIdClientOrderIdPrefix, accountId, clientOrderId)
	var resp Order
	err := m.QueryRowIndex(&resp, orderAccountIdClientOrderIdKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `account_id` = ? and `client_order_id` = ? limit 1", orderRows, m.table)
		if err := conn.QueryRow(&resp, query, accountId, clientOrderId); err != nil {
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

func (m *defaultOrderModel) Update(data *Order, update func()) error {
	orderIdKey := fmt.Sprintf("%s%v", cacheOrderIdPrefix, data.Id)
	orderAccountIdClientOrderIdKey := fmt.Sprintf("%s%v:%v", cacheOrderAccountIdClientOrderIdPrefix, data.AccountId, data.ClientOrderId)
	update()
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, orderRowsWithPlaceHolder)
		return conn.Exec(query, data.AccountId, data.AccountAlias, data.ClientOrderId, data.ExOrderId, data.ApiType, data.Side, data.Status, data.Volume, data.FilledVolume, data.Price, data.Tp, data.StdSymbol, data.ExSymbol, data.Exchange, data.RejectReason, data.CancelFlag, data.Id)
	}, orderIdKey, orderAccountIdClientOrderIdKey)
	return err
}

func (m *defaultOrderModel) TxUpdate(data *Order, update func()) func() (interface{}, error) {
	keys := make([]string, 0)
	insertSql := fmt.Sprintf("update %s set %s where `id` = ?", m.table, orderRowsWithPlaceHolder)
	orderIdKey := fmt.Sprintf("%s%v", cacheOrderIdPrefix, data.Id)
	orderAccountIdClientOrderIdKey := fmt.Sprintf("%s%v:%v", cacheOrderAccountIdClientOrderIdPrefix, data.AccountId, data.ClientOrderId)
	keys = append(keys, orderIdKey, orderAccountIdClientOrderIdKey)
	update()
	args := []interface{}{data.AccountId, data.AccountAlias, data.ClientOrderId, data.ExOrderId, data.ApiType, data.Side, data.Status, data.Volume, data.FilledVolume, data.Price, data.Tp, data.StdSymbol, data.ExSymbol, data.Exchange, data.RejectReason, data.CancelFlag, data.Id}
	return func() (interface{}, error) {
		return &txPrepare_{
			sql:  insertSql,
			args: args,
			keys: keys,
		}, nil
	}
}

func (m *defaultOrderModel) Delete(id string) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	orderIdKey := fmt.Sprintf("%s%v", cacheOrderIdPrefix, id)
	orderAccountIdClientOrderIdKey := fmt.Sprintf("%s%v:%v", cacheOrderAccountIdClientOrderIdPrefix, data.AccountId, data.ClientOrderId)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, orderIdKey, orderAccountIdClientOrderIdKey)
	return err
}

func (m *defaultOrderModel) TxDelete(id string) func() (interface{}, error) {
	data, err := m.FindOne(id)
	if err != nil {
		return func() (interface{}, error) {
			return nil, err
		}
	}
	orderIdKey := fmt.Sprintf("%s%v", cacheOrderIdPrefix, id)
	orderAccountIdClientOrderIdKey := fmt.Sprintf("%s%v:%v", cacheOrderAccountIdClientOrderIdPrefix, data.AccountId, data.ClientOrderId)
	deleteSql := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	args := []interface{}{id}
	keys := make([]string, 0)
	keys = append(keys, orderIdKey, orderAccountIdClientOrderIdKey)

	return func() (interface{}, error) {
		return &txPrepare_{
			sql:  deleteSql,
			args: args,
			keys: keys,
		}, nil
	}
}

func (m *defaultOrderModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheOrderIdPrefix, primary)
}

func (m *defaultOrderModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", orderRows, m.table)
	return conn.QueryRow(v, query, primary)
}
