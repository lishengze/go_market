Insert(data *{{.upperStartCamelObject}}) (sql.Result,error)
TxInsert(data *{{.upperStartCamelObject}}) func() (interface{}, error)
BulkInsert(list []*{{.upperStartCamelObject}}) error
TxBulkInsert(list []*{{.upperStartCamelObject}}) func() (interface{}, error)