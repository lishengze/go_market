package model

import (
	"github.com/shopspring/decimal"
	"testing"
)

func TestBuilder(t *testing.T) {
	builder := NewAccountsUpdater(nil).
		AddFrozenBalance(decimal.NewFromInt(2)).
		SubBalance(decimal.NewFromInt(10)).
		UpdateAveragePrice(decimal.NewFromInt(5)).
		Builder()

	t.Log(builder.Table("accounts").ToSql())
}
