package modelext

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"sync"
)

type (
	TxFactory struct {
		conn sqlc.CachedConn
	}

	Tx struct {
		mutex    sync.Mutex
		conn     sqlc.CachedConn
		prepares []func() (interface{}, error)
	}

	Prepare interface {
		Sql() string
		Args() []interface{}
		DelKeys() []string
		IsEmpty() bool
	}
)

func (t *TxFactory) NewTx() *Tx {
	return &Tx{
		conn:     t.conn,
		prepares: make([]func() (interface{}, error), 0),
	}
}

func (t *Tx) Prepare(fns ...func() (interface{}, error)) *Tx {
	defer t.mutex.Unlock()
	t.mutex.Lock()
	t.prepares = append(t.prepares, fns...)
	return t
}

func (t *Tx) Execute() error {
	if len(t.prepares) == 0 {
		return nil
	}
	var (
		keys     = make([]string, 0)
		prepares []Prepare
	)

	for _, fn := range t.prepares {
		res, err := fn()
		if err != nil {
			return err
		}

		prepare, ok := res.(Prepare)
		if !ok {
			return fmt.Errorf("res is not implement Prepare")
		}

		if prepare.IsEmpty() {
			continue
		}

		keys = append(keys, prepare.DelKeys()...)
		prepares = append(prepares, prepare)
	}

	err := t.conn.DelCache(keys...)
	if err != nil {
		return err
	}

	return t.conn.Transact(func(session sqlx.Session) error {
		for _, prepare := range prepares {
			stmt, err := session.Prepare(prepare.Sql())
			if err != nil {
				return err
			}
			// 返回任何错误都会回滚事务
			if _, err := stmt.Exec(prepare.Args()...); err != nil {
				return err
			}
			_ = stmt.Close()
		}
		return nil
	})
}

func NewTxFactory(conn sqlx.SqlConn, c cache.CacheConf) TxFactory {
	return TxFactory{
		conn: sqlc.NewConn(conn, c),
	}
}
