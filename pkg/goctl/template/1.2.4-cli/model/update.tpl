
func (m *default{{.upperStartCamelObject}}Model) Update(data *{{.upperStartCamelObject}}, update func()) error {
	{{if .withCache}}{{.keys}}
	update()
    _, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}}", m.table, {{.lowerStartCamelObject}}RowsWithPlaceHolder)
		return conn.Exec(query, {{.expressionValues}})
	}, {{.keyValues}}){{else}}query := fmt.Sprintf("update %s set %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}}", m.table, {{.lowerStartCamelObject}}RowsWithPlaceHolder)
    _,err:=m.conn.Exec(query, {{.expressionValues}}){{end}}
	return err
}

func (m *default{{.upperStartCamelObject}}Model) TxUpdate(data *{{.upperStartCamelObject}}, update func()) func() (interface{}, error) {
    keys := make([]string, 0)
   	insertSql := fmt.Sprintf("update %s set %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}}", m.table, {{.lowerStartCamelObject}}RowsWithPlaceHolder)
    {{- if .withCache}}
    {{.keys}}
    keys = append(keys, {{.keyValues}})
    update()
    args := []interface{}{ {{.expressionValues}} }
	{{- else}}
	update(data)
    args := []interface{}{ {{.expressionValues}} }
	{{- end}}
    return func() (interface{}, error) {
		return &txPrepare_{
			sql:  insertSql,
			args: args,
			keys: keys,
		}, nil
	}
}
