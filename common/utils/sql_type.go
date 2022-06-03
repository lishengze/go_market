package utils

import (
	"database/sql"
	"time"
)

func NewInt64(val int64) sql.NullInt64 {
	return sql.NullInt64{Int64: val, Valid: true}
}

func NewString(val string) sql.NullString {
	return sql.NullString{String: val, Valid: true}
}

func NewTime(val time.Time) sql.NullTime {
	return sql.NullTime{Time: val, Valid: true}
}
