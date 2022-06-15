
func (m *default{{.upperStartCamelObject}}Model) Delete({{.lowerStartCamelPrimaryKey}} {{.dataType}}) error {
	{{if .withCache}}{{if .containsIndexCache}}data, err:=m.FindOne({{.lowerStartCamelPrimaryKey}})
	if err!=nil{
		return err
	}{{end}}

	{{.keys}}
    _, err {{if .containsIndexCache}}={{else}}:={{end}} m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}}", m.table)
		return conn.Exec(query, {{.lowerStartCamelPrimaryKey}})
	}, {{.keyValues}}){{else}}query := fmt.Sprintf("delete from %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}}", m.table)
		_,err:=m.conn.Exec(query, {{.lowerStartCamelPrimaryKey}}){{end}}
	return err
}


func (m *default{{.upperStartCamelObject}}Model) TxDelete({{.lowerStartCamelPrimaryKey}} {{.dataType}}) func() (interface{}, error) {
	{{- if .withCache}}
	{{- if .containsIndexCache}}
    data, err:=m.FindOne({{.lowerStartCamelPrimaryKey}})
	if err!=nil{
		return func() (interface{}, error) {
			return nil, err
		}
	}
    {{- end}}
	{{.keys}}
	deleteSql:=  fmt.Sprintf("delete from %s where {{.originalPrimaryKey}} = ?", m.table)
	args:= []interface{}{ {{.lowerStartCamelPrimaryKey}} }
	keys:= make([]string,0)
	keys= append(keys,{{.keyValues}})

	return func() (interface{}, error) {
		return &txPrepare_{
			sql:  deleteSql,
			args: args,
			keys: keys,
		}, nil
	}
   {{- end}}
}