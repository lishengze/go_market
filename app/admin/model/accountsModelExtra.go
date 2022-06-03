package model

import (
	"github.com/Masterminds/squirrel"
	"github.com/shopspring/decimal"
)

type AccountsUpdater struct {
	builder squirrel.UpdateBuilder
	account *Accounts
}

func NewAccountsUpdater(account *Accounts) *AccountsUpdater {
	builder := squirrel.UpdateBuilder{}.
		PlaceholderFormat(squirrel.Question)
	return &AccountsUpdater{
		builder: builder,
		account: account,
	}
}

func (b *AccountsUpdater) AddBalance(balance decimal.Decimal) *AccountsUpdater {
	b.builder = b.builder.Set("balance", squirrel.Expr("balance + cast(? as decimal(26,16))", balance))
	return b
}

func (b *AccountsUpdater) SubBalance(balance decimal.Decimal) *AccountsUpdater {
	b.builder = b.builder.Set("balance", squirrel.Expr("balance - cast(? as decimal(26,16))", balance))
	return b
}

func (b *AccountsUpdater) AddFrozenBalance(balance decimal.Decimal) *AccountsUpdater {
	b.builder = b.builder.Set("frozen_balance", squirrel.Expr("frozen_balance + cast(? as decimal(26,16))", balance))
	return b
}

func (b *AccountsUpdater) SubFrozenBalance(balance decimal.Decimal) *AccountsUpdater {
	b.builder = b.builder.Set("frozen_balance", squirrel.Expr("frozen_balance - cast(? as decimal(26,16))", balance))
	return b
}

func (b *AccountsUpdater) UpdateAveragePrice(price decimal.Decimal) *AccountsUpdater {
	b.builder = b.builder.Set("average_price", squirrel.Expr("cast(? as decimal(26,16))", price))
	return b
}

func (b *AccountsUpdater) Builder() squirrel.UpdateBuilder {
	return b.builder
}
