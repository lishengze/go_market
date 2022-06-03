package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CoinPricesModel = (*customCoinPricesModel)(nil)

type (
	// CoinPricesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCoinPricesModel.
	CoinPricesModel interface {
		coinPricesModel
		CoinPriceByTradeTime(ctx context.Context, dateTime, coin string) (*CoinPrices, error)
	}

	customCoinPricesModel struct {
		*defaultCoinPricesModel
	}
)

func (c *customCoinPricesModel) CoinPriceByTradeTime(ctx context.Context, dateTime, coin string) (coinPrice *CoinPrices, err error) {
	var row CoinPrices
	err = c.conn.QueryRowCtx(ctx, &row, fmt.Sprintf("select * from %s where date_time = ? and coin_unit = ?", c.table), dateTime, coin)
	coinPrice = &row
	return
}

// NewCoinPricesModel returns a model for the database table.
func NewCoinPricesModel(conn sqlx.SqlConn) CoinPricesModel {
	return &customCoinPricesModel{
		defaultCoinPricesModel: newCoinPricesModel(conn),
	}
}
