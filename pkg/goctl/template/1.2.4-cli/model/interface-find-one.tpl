FindOne({{.lowerStartCamelPrimaryKey}} {{.dataType}}) (*{{.upperStartCamelObject}}, error)
FindMany(q Query_) ([]*{{.upperStartCamelObject}}, error)
FindPks(q Query_) ([]{{.dataType}}, error)
FindManyByPks(pks []{{.dataType}}) ([]*{{.upperStartCamelObject}}, error)
