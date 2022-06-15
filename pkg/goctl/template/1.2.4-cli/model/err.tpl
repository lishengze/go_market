package {{.pkg}}

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var ErrNotFound = sqlx.ErrNotFound

type (
	counter_ struct {
		count int64 `db:"count(*)"`
	}

	Query_ interface {
		FindPksSql(pk, table string) (string, string, []interface{})
		HasPaginator() bool
		Size() int64
		Page() int64
		SetCount(count int64)
	}

	txPrepare_ struct {
		sql  string
		args []interface{}
		keys []string
	}
)

func newEmptyTxPrepare_() *txPrepare_ {
	return &txPrepare_{}
}

func (o *txPrepare_) Sql() string {
	return o.sql
}

func (o *txPrepare_) Args() []interface{} {
	return o.args
}

func (o *txPrepare_) DelKeys() []string {
	return o.keys
}

func (o *txPrepare_) IsEmpty() bool {
	return o.sql == ""
}
