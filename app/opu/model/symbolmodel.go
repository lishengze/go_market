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
	symbolFieldNames          = builder.RawFieldNames(&Symbol{})
	symbolRows                = strings.Join(symbolFieldNames, ",")
	symbolRowsExpectAutoSet   = strings.Join(stringx.Remove(symbolFieldNames, "`create_time`", "`update_time`"), ",")
	symbolRowsWithPlaceHolder = strings.Join(stringx.Remove(symbolFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheSymbolIdPrefix                      = "cache:symbol:id:"
	cacheSymbolExFormatApiTypeExchangePrefix = "cache:symbol:exFormat:apiType:exchange:"
	cacheSymbolStdSymbolExchangePrefix       = "cache:symbol:stdSymbol:exchange:"
)

type (
	SymbolModel interface {
		Insert(data *Symbol) (sql.Result, error)
		TxInsert(data *Symbol) func() (interface{}, error)
		BulkInsert(list []*Symbol) error
		TxBulkInsert(list []*Symbol) func() (interface{}, error)
		FindOne(id string) (*Symbol, error)
		FindMany(q Query_) ([]*Symbol, error)
		FindPks(q Query_) ([]string, error)
		FindManyByPks(pks []string) ([]*Symbol, error)

		FindOneByExFormatApiTypeExchange(exFormat string, apiType string, exchange string) (*Symbol, error)
		FindOneByStdSymbolExchange(stdSymbol string, exchange string) (*Symbol, error)
		Update(data *Symbol, update func()) error
		TxUpdate(data *Symbol, update func()) func() (interface{}, error)
		Delete(id string) error
		TxDelete(id string) func() (interface{}, error)
	}

	defaultSymbolModel struct {
		sqlc.CachedConn
		table string
	}

	Symbol struct {
		Id            string    `db:"id"`             // id
		Tp            string    `db:"type"`           // type
		ApiType       string    `db:"api_type"`       // api_type
		StdSymbol     string    `db:"std_symbol"`     // std_symbol
		ExFormat      string    `db:"ex_format"`      // std_symbol
		BaseCurrency  string    `db:"base_currency"`  // base_currency
		QuoteCurrency string    `db:"quote_currency"` // quote_currency
		Exchange      string    `db:"exchange"`       // exchange
		VolumeScale   string    `db:"volume_scale"`   // volume_scale
		PriceScale    string    `db:"price_scale"`    // price_scale
		MinVolume     string    `db:"min_volume"`     // min_volume
		ContractSize  string    `db:"contract_size"`  // contract_size
		CreateTime    time.Time `db:"create_time"`    // 创建时间
		UpdateTime    time.Time `db:"update_time"`    // 更新时间
	}
)

func NewSymbolModel(conn sqlx.SqlConn, c cache.CacheConf) SymbolModel {
	return &defaultSymbolModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`symbol`",
	}
}

func (m *defaultSymbolModel) Insert(data *Symbol) (sql.Result, error) {
	symbolIdKey := fmt.Sprintf("%s%v", cacheSymbolIdPrefix, data.Id)
	symbolExFormatApiTypeExchangeKey := fmt.Sprintf("%s%v:%v:%v", cacheSymbolExFormatApiTypeExchangePrefix, data.ExFormat, data.ApiType, data.Exchange)
	symbolStdSymbolExchangeKey := fmt.Sprintf("%s%v:%v", cacheSymbolStdSymbolExchangePrefix, data.StdSymbol, data.Exchange)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, symbolRowsExpectAutoSet)
		return conn.Exec(query, data.Id, data.Tp, data.ApiType, data.StdSymbol, data.ExFormat, data.BaseCurrency, data.QuoteCurrency, data.Exchange, data.VolumeScale, data.PriceScale, data.MinVolume, data.ContractSize)
	}, symbolStdSymbolExchangeKey, symbolIdKey, symbolExFormatApiTypeExchangeKey)
	return ret, err
}

func (m *defaultSymbolModel) TxInsert(data *Symbol) func() (interface{}, error) {
	keys := make([]string, 0)
	args := []interface{}{data.Id, data.Tp, data.ApiType, data.StdSymbol, data.ExFormat, data.BaseCurrency, data.QuoteCurrency, data.Exchange, data.VolumeScale, data.PriceScale, data.MinVolume, data.ContractSize}
	insertSql := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, symbolRowsExpectAutoSet)
	symbolIdKey := fmt.Sprintf("%s%v", cacheSymbolIdPrefix, data.Id)
	symbolExFormatApiTypeExchangeKey := fmt.Sprintf("%s%v:%v:%v", cacheSymbolExFormatApiTypeExchangePrefix, data.ExFormat, data.ApiType, data.Exchange)
	symbolStdSymbolExchangeKey := fmt.Sprintf("%s%v:%v", cacheSymbolStdSymbolExchangePrefix, data.StdSymbol, data.Exchange)
	keys = append(keys, symbolStdSymbolExchangeKey, symbolIdKey, symbolExFormatApiTypeExchangeKey)
	return func() (interface{}, error) {
		return &txPrepare_{
			sql:  insertSql,
			args: args,
			keys: keys,
		}, nil
	}
}

func (m *defaultSymbolModel) BulkInsert(symbols []*Symbol) error {
	if len(symbols) == 0 {
		return nil
	}

	var (
		insertSql    = fmt.Sprintf("insert into %s (%s) values ", m.table, symbolRowsExpectAutoSet)
		args         = make([]interface{}, 0)
		keys         = make([]string, 0)
		placeHolders = make([]string, 0)
	)
	for _, data := range symbols {
		symbolIdKey := fmt.Sprintf("%s%v", cacheSymbolIdPrefix, data.Id)
		symbolExFormatApiTypeExchangeKey := fmt.Sprintf("%s%v:%v:%v", cacheSymbolExFormatApiTypeExchangePrefix, data.ExFormat, data.ApiType, data.Exchange)
		symbolStdSymbolExchangeKey := fmt.Sprintf("%s%v:%v", cacheSymbolStdSymbolExchangePrefix, data.StdSymbol, data.Exchange)
		keys = append(keys, symbolStdSymbolExchangeKey, symbolIdKey, symbolExFormatApiTypeExchangeKey)
		args = append(args, data.Id, data.Tp, data.ApiType, data.StdSymbol, data.ExFormat, data.BaseCurrency, data.QuoteCurrency, data.Exchange, data.VolumeScale, data.PriceScale, data.MinVolume, data.ContractSize)
		placeHolders = append(placeHolders, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	}

	insertSql += strings.Join(placeHolders, ",")
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.Exec(insertSql, args...)
	}, keys...)

	return err
}

func (m *defaultSymbolModel) TxBulkInsert(symbols []*Symbol) func() (interface{}, error) {
	if len(symbols) == 0 {
		return func() (interface{}, error) {
			return newEmptyTxPrepare_(), nil
		}
	}
	var (
		insertSql    = fmt.Sprintf("insert into %s (%s) values ", m.table, symbolRowsExpectAutoSet)
		args         = make([]interface{}, 0)
		keys         = make([]string, 0)
		placeHolders = make([]string, 0)
	)
	for _, data := range symbols {
		symbolIdKey := fmt.Sprintf("%s%v", cacheSymbolIdPrefix, data.Id)
		symbolExFormatApiTypeExchangeKey := fmt.Sprintf("%s%v:%v:%v", cacheSymbolExFormatApiTypeExchangePrefix, data.ExFormat, data.ApiType, data.Exchange)
		symbolStdSymbolExchangeKey := fmt.Sprintf("%s%v:%v", cacheSymbolStdSymbolExchangePrefix, data.StdSymbol, data.Exchange)
		keys = append(keys, symbolStdSymbolExchangeKey, symbolIdKey, symbolExFormatApiTypeExchangeKey)
		args = append(args, data.Id, data.Tp, data.ApiType, data.StdSymbol, data.ExFormat, data.BaseCurrency, data.QuoteCurrency, data.Exchange, data.VolumeScale, data.PriceScale, data.MinVolume, data.ContractSize)
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

func (m *defaultSymbolModel) FindOne(id string) (*Symbol, error) {
	symbolIdKey := fmt.Sprintf("%s%v", cacheSymbolIdPrefix, id)
	var resp Symbol
	err := m.QueryRow(&resp, symbolIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", symbolRows, m.table)
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

func (m *defaultSymbolModel) FindMany(q Query_) ([]*Symbol, error) {
	pks, err := m.FindPks(q)
	if err != nil {
		return nil, err
	}
	return m.FindManyByPks(pks)
}

func (m *defaultSymbolModel) FindPks(q Query_) ([]string, error) {
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

func (m *defaultSymbolModel) FindManyByPks(pks []string) ([]*Symbol, error) {
	// 使用MapReduce 处理并发
	r, err := mr.MapReduce(func(source chan<- interface{}) {
		for index, pk := range pks {
			source <- []interface{}{index, pk}
		}
	}, func(item interface{}, writer mr.Writer, cancel func(error)) {
		indexAndPk := item.([]interface{})
		Symbol, err := m.FindOne(indexAndPk[1].(string))
		if err != nil {
			cancel(err)
			return
		}
		writer.Write([]interface{}{indexAndPk[0], Symbol})
	}, func(pipe <-chan interface{}, writer mr.Writer, cancel func(error)) {
		var symbols = make([]*Symbol, len(pks))
		for content := range pipe {
			symbols[content.([]interface{})[0].(int)] = content.([]interface{})[1].(*Symbol)
		}
		writer.Write(symbols)
	}, mr.WithWorkers(len(pks)))
	if err != nil {
		return nil, err
	}
	return r.([]*Symbol), nil
}

func (m *defaultSymbolModel) FindOneByExFormatApiTypeExchange(exFormat string, apiType string, exchange string) (*Symbol, error) {
	symbolExFormatApiTypeExchangeKey := fmt.Sprintf("%s%v:%v:%v", cacheSymbolExFormatApiTypeExchangePrefix, exFormat, apiType, exchange)
	var resp Symbol
	err := m.QueryRowIndex(&resp, symbolExFormatApiTypeExchangeKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `ex_format` = ? and `api_type` = ? and `exchange` = ? limit 1", symbolRows, m.table)
		if err := conn.QueryRow(&resp, query, exFormat, apiType, exchange); err != nil {
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

func (m *defaultSymbolModel) FindOneByStdSymbolExchange(stdSymbol string, exchange string) (*Symbol, error) {
	symbolStdSymbolExchangeKey := fmt.Sprintf("%s%v:%v", cacheSymbolStdSymbolExchangePrefix, stdSymbol, exchange)
	var resp Symbol
	err := m.QueryRowIndex(&resp, symbolStdSymbolExchangeKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `std_symbol` = ? and `exchange` = ? limit 1", symbolRows, m.table)
		if err := conn.QueryRow(&resp, query, stdSymbol, exchange); err != nil {
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

func (m *defaultSymbolModel) Update(data *Symbol, update func()) error {
	symbolIdKey := fmt.Sprintf("%s%v", cacheSymbolIdPrefix, data.Id)
	symbolExFormatApiTypeExchangeKey := fmt.Sprintf("%s%v:%v:%v", cacheSymbolExFormatApiTypeExchangePrefix, data.ExFormat, data.ApiType, data.Exchange)
	symbolStdSymbolExchangeKey := fmt.Sprintf("%s%v:%v", cacheSymbolStdSymbolExchangePrefix, data.StdSymbol, data.Exchange)
	update()
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, symbolRowsWithPlaceHolder)
		return conn.Exec(query, data.Tp, data.ApiType, data.StdSymbol, data.ExFormat, data.BaseCurrency, data.QuoteCurrency, data.Exchange, data.VolumeScale, data.PriceScale, data.MinVolume, data.ContractSize, data.Id)
	}, symbolIdKey, symbolExFormatApiTypeExchangeKey, symbolStdSymbolExchangeKey)
	return err
}

func (m *defaultSymbolModel) TxUpdate(data *Symbol, update func()) func() (interface{}, error) {
	keys := make([]string, 0)
	insertSql := fmt.Sprintf("update %s set %s where `id` = ?", m.table, symbolRowsWithPlaceHolder)
	symbolIdKey := fmt.Sprintf("%s%v", cacheSymbolIdPrefix, data.Id)
	symbolExFormatApiTypeExchangeKey := fmt.Sprintf("%s%v:%v:%v", cacheSymbolExFormatApiTypeExchangePrefix, data.ExFormat, data.ApiType, data.Exchange)
	symbolStdSymbolExchangeKey := fmt.Sprintf("%s%v:%v", cacheSymbolStdSymbolExchangePrefix, data.StdSymbol, data.Exchange)
	keys = append(keys, symbolIdKey, symbolExFormatApiTypeExchangeKey, symbolStdSymbolExchangeKey)
	update()
	args := []interface{}{data.Tp, data.ApiType, data.StdSymbol, data.ExFormat, data.BaseCurrency, data.QuoteCurrency, data.Exchange, data.VolumeScale, data.PriceScale, data.MinVolume, data.ContractSize, data.Id}
	return func() (interface{}, error) {
		return &txPrepare_{
			sql:  insertSql,
			args: args,
			keys: keys,
		}, nil
	}
}

func (m *defaultSymbolModel) Delete(id string) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	symbolStdSymbolExchangeKey := fmt.Sprintf("%s%v:%v", cacheSymbolStdSymbolExchangePrefix, data.StdSymbol, data.Exchange)
	symbolIdKey := fmt.Sprintf("%s%v", cacheSymbolIdPrefix, id)
	symbolExFormatApiTypeExchangeKey := fmt.Sprintf("%s%v:%v:%v", cacheSymbolExFormatApiTypeExchangePrefix, data.ExFormat, data.ApiType, data.Exchange)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, symbolIdKey, symbolExFormatApiTypeExchangeKey, symbolStdSymbolExchangeKey)
	return err
}

func (m *defaultSymbolModel) TxDelete(id string) func() (interface{}, error) {
	data, err := m.FindOne(id)
	if err != nil {
		return func() (interface{}, error) {
			return nil, err
		}
	}
	symbolStdSymbolExchangeKey := fmt.Sprintf("%s%v:%v", cacheSymbolStdSymbolExchangePrefix, data.StdSymbol, data.Exchange)
	symbolIdKey := fmt.Sprintf("%s%v", cacheSymbolIdPrefix, id)
	symbolExFormatApiTypeExchangeKey := fmt.Sprintf("%s%v:%v:%v", cacheSymbolExFormatApiTypeExchangePrefix, data.ExFormat, data.ApiType, data.Exchange)
	deleteSql := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	args := []interface{}{id}
	keys := make([]string, 0)
	keys = append(keys, symbolIdKey, symbolExFormatApiTypeExchangeKey, symbolStdSymbolExchangeKey)

	return func() (interface{}, error) {
		return &txPrepare_{
			sql:  deleteSql,
			args: args,
			keys: keys,
		}, nil
	}
}

func (m *defaultSymbolModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheSymbolIdPrefix, primary)
}

func (m *defaultSymbolModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", symbolRows, m.table)
	return conn.QueryRow(v, query, primary)
}
