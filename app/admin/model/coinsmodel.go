package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CoinsModel = (*customCoinsModel)(nil)

type (
	// CoinsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCoinsModel.
	CoinsModel interface {
		coinsModel
	}

	customCoinsModel struct {
		*defaultCoinsModel
	}
)

// NewCoinsModel returns a model for the database table.
func NewCoinsModel(conn sqlx.SqlConn) CoinsModel {
	return &customCoinsModel{
		defaultCoinsModel: newCoinsModel(conn),
	}
}
