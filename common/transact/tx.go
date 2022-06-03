package transact

import (
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
		mutex sync.Mutex
		conn  sqlc.CachedConn
		fns   []func() (*Prepare, error)
	}

	Prepare struct {
		Sql  string
		Args []interface{}
		Keys []string
	}
)

func (t *TxFactory) NewTx() *Tx {
	return &Tx{
		conn: t.conn,
		fns:  make([]func() (*Prepare, error), 0),
	}
}

func (t *Tx) Prepare(fns ...func() (*Prepare, error)) *Tx {
	defer t.mutex.Unlock()
	t.mutex.Lock()
	t.fns = append(t.fns, fns...)
	return t
}

func (t *Tx) Execute() error {
	if len(t.fns) == 0 {
		return nil
	}
	var (
		keys     = make([]string, 0)
		prepares []*Prepare
	)

	for _, fn := range t.fns {
		prepare, err := fn()
		if err != nil {
			return err
		}
		if prepare == nil {
			continue
		}
		keys = append(keys, prepare.Keys...)
		prepares = append(prepares, prepare)
	}

	err := t.conn.DelCache(keys...)
	if err != nil {
		return err
	}

	return t.conn.Transact(func(session sqlx.Session) error {
		for _, prepare := range prepares {
			stmt, err := session.Prepare(prepare.Sql)
			if err != nil {
				return err
			}
			// 返回任何错误都会回滚事务
			if _, err := stmt.Exec(prepare.Args...); err != nil {
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
