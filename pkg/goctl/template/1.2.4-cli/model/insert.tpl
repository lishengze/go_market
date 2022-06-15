
func (m *default{{.upperStartCamelObject}}Model) Insert(data *{{.upperStartCamelObject}}) (sql.Result,error) {
	{{- if .withCache}}{{.keys}}
    ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values ({{.expression}})", m.table, {{.lowerStartCamelObject}}RowsExpectAutoSet)
		return conn.Exec(query, {{.expressionValues}})
	}, {{.keyValues}})
	{{- else}}query := fmt.Sprintf("insert into %s (%s) values ({{.expression}})", m.table, {{.lowerStartCamelObject}}RowsExpectAutoSet)
    ret,err:=m.conn.Exec(query, {{.expressionValues}}){{end}}
	return ret,err
}

func (m *default{{.upperStartCamelObject}}Model) TxInsert(data *{{.upperStartCamelObject}}) func() (interface{}, error) {
	keys := make([]string, 0)
	args := []interface{}{ {{.expressionValues}} }
    insertSql := fmt.Sprintf("insert into %s (%s) values ({{.expression}})", m.table, {{.lowerStartCamelObject}}RowsExpectAutoSet)
	{{- if .withCache}}
	{{.keys}}
    keys = append(keys, {{.keyValues}} )
    {{- end}}
	return func() (interface{}, error) {
		return &txPrepare_{
			sql:  insertSql,
			args: args,
			keys: keys,
		}, nil
	}
}

func (m *default{{.upperStartCamelObject}}Model) BulkInsert({{.lowerStartCamelObject}}s []*{{.upperStartCamelObject}}) error {
	if len({{.lowerStartCamelObject}}s) == 0 {
		return nil
	}

	var (
		insertSql      = fmt.Sprintf("insert into %s (%s) values ", m.table, {{.lowerStartCamelObject}}RowsExpectAutoSet)
		args           = make([]interface{}, 0)
		keys           = make([]string, 0)
		placeHolders   = make([]string, 0)
	)
	for _, data := range {{.lowerStartCamelObject}}s {
        {{if .withCache}}{{.keys}}
		keys = append(keys, {{.keyValues}}){{end}}
		args = append(args, {{.expressionValues}})
		placeHolders = append(placeHolders, "({{.expression}})")
	}

	insertSql += strings.Join(placeHolders, ",")
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.Exec(insertSql, args...)
	}, keys...)

	return  err
}


func (m *default{{.upperStartCamelObject}}Model) TxBulkInsert({{.lowerStartCamelObject}}s []*{{.upperStartCamelObject}}) func() (interface{}, error) {
	if len({{.lowerStartCamelObject}}s) == 0 {
		return func() (interface{}, error) {
			return newEmptyTxPrepare_(), nil
		}
	}
	var (
		insertSql      = fmt.Sprintf("insert into %s (%s) values ", m.table, {{.lowerStartCamelObject}}RowsExpectAutoSet)
		args           = make([]interface{}, 0)
		keys           = make([]string, 0)
		placeHolders   = make([]string, 0)
	)
	for _, data := range {{.lowerStartCamelObject}}s {
        {{if .withCache}}{{.keys}}
		keys = append(keys, {{.keyValues}}){{end}}
		args = append(args, {{.expressionValues}})
		placeHolders = append(placeHolders, "({{.expression}})")
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