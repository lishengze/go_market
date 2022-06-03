package svc

import (
	"market_server/app/pms/model"
	"market_server/app/pms/rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config                        config.Config
	RedisConn                     *redis.Redis
	PortfolioInvestmentItemsModel model.PortfolioInvestmentItemsModel
	PortfolioInvestmentsModel     model.PortfolioInvestmentsModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:                    c,
		RedisConn:                 c.Redis.NewRedis(),
		PortfolioInvestmentsModel: model.NewPortfolioInvestmentsModel(conn, c.Cache),
	}
}
