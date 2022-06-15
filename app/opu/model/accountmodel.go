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
	accountFieldNames          = builder.RawFieldNames(&Account{})
	accountRows                = strings.Join(accountFieldNames, ",")
	accountRowsExpectAutoSet   = strings.Join(stringx.Remove(accountFieldNames, "`create_time`", "`update_time`"), ",")
	accountRowsWithPlaceHolder = strings.Join(stringx.Remove(accountFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheAccountIdPrefix       = "cache:account:id:"
	cacheAccountAliasKeyPrefix = "cache:account:alias:key:"
	cacheAccountAliasPrefix    = "cache:account:alias:"
)

type (
	AccountModel interface {
		Insert(data *Account) (sql.Result, error)
		TxInsert(data *Account) func() (interface{}, error)
		BulkInsert(list []*Account) error
		TxBulkInsert(list []*Account) func() (interface{}, error)
		FindOne(id string) (*Account, error)
		FindMany(q Query_) ([]*Account, error)
		FindPks(q Query_) ([]string, error)
		FindManyByPks(pks []string) ([]*Account, error)

		FindOneByAliasKey(alias string, key string) (*Account, error)
		FindOneByAlias(alias string) (*Account, error)
		Update(data *Account, update func()) error
		TxUpdate(data *Account, update func()) func() (interface{}, error)
		Delete(id string) error
		TxDelete(id string) func() (interface{}, error)
	}

	defaultAccountModel struct {
		sqlc.CachedConn
		table string
	}

	Account struct {
		Id             string    `db:"id"`               // id
		Alias          string    `db:"alias"`            // alias
		Key            string    `db:"key"`              // key
		Secret         string    `db:"secret"`           // secret
		Passphrase     string    `db:"passphrase"`       // passphrase
		SubAccountName string    `db:"sub_account_name"` // sub_account_name
		Exchange       string    `db:"exchange"`         // exchange
		CreateTime     time.Time `db:"create_time"`      // 创建时间
		UpdateTime     time.Time `db:"update_time"`      // 更新时间
	}
)

func NewAccountModel(conn sqlx.SqlConn, c cache.CacheConf) AccountModel {
	return &defaultAccountModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`account`",
	}
}

func (m *defaultAccountModel) Insert(data *Account) (sql.Result, error) {
	accountIdKey := fmt.Sprintf("%s%v", cacheAccountIdPrefix, data.Id)
	accountAliasKeyKey := fmt.Sprintf("%s%v:%v", cacheAccountAliasKeyPrefix, data.Alias, data.Key)
	accountAliasKey := fmt.Sprintf("%s%v", cacheAccountAliasPrefix, data.Alias)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?)", m.table, accountRowsExpectAutoSet)
		return conn.Exec(query, data.Id, data.Alias, data.Key, data.Secret, data.Passphrase, data.SubAccountName, data.Exchange)
	}, accountIdKey, accountAliasKeyKey, accountAliasKey)
	return ret, err
}

func (m *defaultAccountModel) TxInsert(data *Account) func() (interface{}, error) {
	keys := make([]string, 0)
	args := []interface{}{data.Id, data.Alias, data.Key, data.Secret, data.Passphrase, data.SubAccountName, data.Exchange}
	insertSql := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?)", m.table, accountRowsExpectAutoSet)
	accountIdKey := fmt.Sprintf("%s%v", cacheAccountIdPrefix, data.Id)
	accountAliasKeyKey := fmt.Sprintf("%s%v:%v", cacheAccountAliasKeyPrefix, data.Alias, data.Key)
	accountAliasKey := fmt.Sprintf("%s%v", cacheAccountAliasPrefix, data.Alias)
	keys = append(keys, accountIdKey, accountAliasKeyKey, accountAliasKey)
	return func() (interface{}, error) {
		return &txPrepare_{
			sql:  insertSql,
			args: args,
			keys: keys,
		}, nil
	}
}

func (m *defaultAccountModel) BulkInsert(accounts []*Account) error {
	if len(accounts) == 0 {
		return nil
	}

	var (
		insertSql    = fmt.Sprintf("insert into %s (%s) values ", m.table, accountRowsExpectAutoSet)
		args         = make([]interface{}, 0)
		keys         = make([]string, 0)
		placeHolders = make([]string, 0)
	)
	for _, data := range accounts {
		accountIdKey := fmt.Sprintf("%s%v", cacheAccountIdPrefix, data.Id)
		accountAliasKeyKey := fmt.Sprintf("%s%v:%v", cacheAccountAliasKeyPrefix, data.Alias, data.Key)
		accountAliasKey := fmt.Sprintf("%s%v", cacheAccountAliasPrefix, data.Alias)
		keys = append(keys, accountIdKey, accountAliasKeyKey, accountAliasKey)
		args = append(args, data.Id, data.Alias, data.Key, data.Secret, data.Passphrase, data.SubAccountName, data.Exchange)
		placeHolders = append(placeHolders, "(?, ?, ?, ?, ?, ?, ?)")
	}

	insertSql += strings.Join(placeHolders, ",")
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.Exec(insertSql, args...)
	}, keys...)

	return err
}

func (m *defaultAccountModel) TxBulkInsert(accounts []*Account) func() (interface{}, error) {
	if len(accounts) == 0 {
		return func() (interface{}, error) {
			return newEmptyTxPrepare_(), nil
		}
	}
	var (
		insertSql    = fmt.Sprintf("insert into %s (%s) values ", m.table, accountRowsExpectAutoSet)
		args         = make([]interface{}, 0)
		keys         = make([]string, 0)
		placeHolders = make([]string, 0)
	)
	for _, data := range accounts {
		accountIdKey := fmt.Sprintf("%s%v", cacheAccountIdPrefix, data.Id)
		accountAliasKeyKey := fmt.Sprintf("%s%v:%v", cacheAccountAliasKeyPrefix, data.Alias, data.Key)
		accountAliasKey := fmt.Sprintf("%s%v", cacheAccountAliasPrefix, data.Alias)
		keys = append(keys, accountIdKey, accountAliasKeyKey, accountAliasKey)
		args = append(args, data.Id, data.Alias, data.Key, data.Secret, data.Passphrase, data.SubAccountName, data.Exchange)
		placeHolders = append(placeHolders, "(?, ?, ?, ?, ?, ?, ?)")
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

func (m *defaultAccountModel) FindOne(id string) (*Account, error) {
	accountIdKey := fmt.Sprintf("%s%v", cacheAccountIdPrefix, id)
	var resp Account
	err := m.QueryRow(&resp, accountIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", accountRows, m.table)
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

func (m *defaultAccountModel) FindMany(q Query_) ([]*Account, error) {
	pks, err := m.FindPks(q)
	if err != nil {
		return nil, err
	}
	return m.FindManyByPks(pks)
}

func (m *defaultAccountModel) FindPks(q Query_) ([]string, error) {
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

func (m *defaultAccountModel) FindManyByPks(pks []string) ([]*Account, error) {
	// 使用MapReduce 处理并发
	r, err := mr.MapReduce(func(source chan<- interface{}) {
		for index, pk := range pks {
			source <- []interface{}{index, pk}
		}
	}, func(item interface{}, writer mr.Writer, cancel func(error)) {
		indexAndPk := item.([]interface{})
		Account, err := m.FindOne(indexAndPk[1].(string))
		if err != nil {
			cancel(err)
			return
		}
		writer.Write([]interface{}{indexAndPk[0], Account})
	}, func(pipe <-chan interface{}, writer mr.Writer, cancel func(error)) {
		var accounts = make([]*Account, len(pks))
		for content := range pipe {
			accounts[content.([]interface{})[0].(int)] = content.([]interface{})[1].(*Account)
		}
		writer.Write(accounts)
	}, mr.WithWorkers(len(pks)))
	if err != nil {
		return nil, err
	}
	return r.([]*Account), nil
}

func (m *defaultAccountModel) FindOneByAliasKey(alias string, key string) (*Account, error) {
	accountAliasKeyKey := fmt.Sprintf("%s%v:%v", cacheAccountAliasKeyPrefix, alias, key)
	var resp Account
	err := m.QueryRowIndex(&resp, accountAliasKeyKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `alias` = ? and `key` = ? limit 1", accountRows, m.table)
		if err := conn.QueryRow(&resp, query, alias, key); err != nil {
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

func (m *defaultAccountModel) FindOneByAlias(alias string) (*Account, error) {
	accountAliasKey := fmt.Sprintf("%s%v", cacheAccountAliasPrefix, alias)
	var resp Account
	err := m.QueryRowIndex(&resp, accountAliasKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `alias` = ? limit 1", accountRows, m.table)
		if err := conn.QueryRow(&resp, query, alias); err != nil {
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

func (m *defaultAccountModel) Update(data *Account, update func()) error {
	accountIdKey := fmt.Sprintf("%s%v", cacheAccountIdPrefix, data.Id)
	accountAliasKeyKey := fmt.Sprintf("%s%v:%v", cacheAccountAliasKeyPrefix, data.Alias, data.Key)
	accountAliasKey := fmt.Sprintf("%s%v", cacheAccountAliasPrefix, data.Alias)
	update()
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, accountRowsWithPlaceHolder)
		return conn.Exec(query, data.Alias, data.Key, data.Secret, data.Passphrase, data.SubAccountName, data.Exchange, data.Id)
	}, accountIdKey, accountAliasKeyKey, accountAliasKey)
	return err
}

func (m *defaultAccountModel) TxUpdate(data *Account, update func()) func() (interface{}, error) {
	keys := make([]string, 0)
	insertSql := fmt.Sprintf("update %s set %s where `id` = ?", m.table, accountRowsWithPlaceHolder)
	accountIdKey := fmt.Sprintf("%s%v", cacheAccountIdPrefix, data.Id)
	accountAliasKeyKey := fmt.Sprintf("%s%v:%v", cacheAccountAliasKeyPrefix, data.Alias, data.Key)
	accountAliasKey := fmt.Sprintf("%s%v", cacheAccountAliasPrefix, data.Alias)
	keys = append(keys, accountIdKey, accountAliasKeyKey, accountAliasKey)
	update()
	args := []interface{}{data.Alias, data.Key, data.Secret, data.Passphrase, data.SubAccountName, data.Exchange, data.Id}
	return func() (interface{}, error) {
		return &txPrepare_{
			sql:  insertSql,
			args: args,
			keys: keys,
		}, nil
	}
}

func (m *defaultAccountModel) Delete(id string) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	accountIdKey := fmt.Sprintf("%s%v", cacheAccountIdPrefix, id)
	accountAliasKeyKey := fmt.Sprintf("%s%v:%v", cacheAccountAliasKeyPrefix, data.Alias, data.Key)
	accountAliasKey := fmt.Sprintf("%s%v", cacheAccountAliasPrefix, data.Alias)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, accountIdKey, accountAliasKeyKey, accountAliasKey)
	return err
}

func (m *defaultAccountModel) TxDelete(id string) func() (interface{}, error) {
	data, err := m.FindOne(id)
	if err != nil {
		return func() (interface{}, error) {
			return nil, err
		}
	}
	accountIdKey := fmt.Sprintf("%s%v", cacheAccountIdPrefix, id)
	accountAliasKeyKey := fmt.Sprintf("%s%v:%v", cacheAccountAliasKeyPrefix, data.Alias, data.Key)
	accountAliasKey := fmt.Sprintf("%s%v", cacheAccountAliasPrefix, data.Alias)
	deleteSql := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	args := []interface{}{id}
	keys := make([]string, 0)
	keys = append(keys, accountIdKey, accountAliasKeyKey, accountAliasKey)

	return func() (interface{}, error) {
		return &txPrepare_{
			sql:  deleteSql,
			args: args,
			keys: keys,
		}, nil
	}
}

func (m *defaultAccountModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheAccountIdPrefix, primary)
}

func (m *defaultAccountModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", accountRows, m.table)
	return conn.QueryRow(v, query, primary)
}
