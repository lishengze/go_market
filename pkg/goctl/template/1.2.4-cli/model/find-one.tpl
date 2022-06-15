
func (m *default{{.upperStartCamelObject}}Model) FindOne({{.lowerStartCamelPrimaryKey}} {{.dataType}}) (*{{.upperStartCamelObject}}, error) {
	{{if .withCache}}{{.cacheKey}}
	var resp {{.upperStartCamelObject}}
	err := m.QueryRow(&resp, {{.cacheKeyVariable}}, func(conn sqlx.SqlConn, v interface{}) error {
		query :=  fmt.Sprintf("select %s from %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}} limit 1", {{.lowerStartCamelObject}}Rows, m.table)
		return conn.QueryRow(v, query, {{.lowerStartCamelPrimaryKey}})
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}{{else}}query := fmt.Sprintf("select %s from %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}} limit 1", {{.lowerStartCamelObject}}Rows, m.table)
	var resp {{.upperStartCamelObject}}
	err := m.conn.QueryRow(&resp, query, {{.lowerStartCamelPrimaryKey}})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}{{end}}
}

func (m *default{{.upperStartCamelObject}}Model) FindMany(q Query_) ([]*{{.upperStartCamelObject}}, error) {
	pks,err:= m.FindPks(q)
	if err!=nil{
	    return nil, err
	}
	return m.FindManyByPks(pks)
}

func (m *default{{.upperStartCamelObject}}Model) FindPks(q Query_) ([]{{.dataType}}, error) {
	var (
		pks   = make([]{{.dataType}}, 0)
	)
	findPksSql, countSql, args := q.FindPksSql("`{{.lowerStartCamelPrimaryKey}}`", m.table)
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

func (m *default{{.upperStartCamelObject}}Model) FindManyByPks(pks []{{.dataType}}) ([]*{{.upperStartCamelObject}}, error) {
	// 使用MapReduce 处理并发
	r, err := mr.MapReduce(func(source chan<- interface{}) {
		for index, pk := range pks {
			source <- []interface{}{index, pk}
		}
	}, func(item interface{}, writer mr.Writer, cancel func(error)) {
		indexAndPk := item.([]interface{})
		{{.upperStartCamelObject}}, err := m.FindOne(indexAndPk[1].({{.dataType}}))
		if err != nil {
			cancel(err)
			return
		}
		writer.Write([]interface{}{indexAndPk[0], {{.upperStartCamelObject}}})
	}, func(pipe <-chan interface{}, writer mr.Writer, cancel func(error)) {
		var {{.lowerStartCamelObject}}s = make([]*{{.upperStartCamelObject}}, len(pks))
		for content := range pipe {
			{{.lowerStartCamelObject}}s[content.([]interface{})[0].(int)] = content.([]interface{})[1].(*{{.upperStartCamelObject}})
		}
		writer.Write({{.lowerStartCamelObject}}s)
	}, mr.WithWorkers(len(pks)))
	if err != nil {
		return nil, err
	}
	return r.([]*{{.upperStartCamelObject}}), nil
}
