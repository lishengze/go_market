package svc

import (
	"market_server/app/admin/api/internal/config"
	"market_server/app/admin/api/internal/middleware"
	"market_server/app/admin/model"
	"market_server/common/nacosAdapter"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config                        config.Config
	JwtAuthMiddleware             rest.Middleware
	RbacMiddleware                rest.Middleware
	Mysql                         sqlx.SqlConn
	Parameters                    *nacosAdapter.Client
	PriceRedis                    *redis.Redis // 注意这里不要使用go-zero自带的redis,默认使用db是0的数据库
	AccountsModel                 model.AccountsModel
	AdminUserModel                model.AdminUserModel
	CoinPricesModel               model.CoinPricesModel
	ExchangeTradesModel           model.ExchangeTradesModel
	MembersModel                  model.MembersModel
	MemberTransactionsModel       model.MemberTransactionsModel
	MenuModel                     model.MenuModel
	OfflineTradeInputModel        model.OfflineTradeInputModel
	OfflineTradeInputHistoryModel model.OfflineTradeInputHistoryModel
	ProfitsModel                  model.ProfitsModel
	SecurityLogModel              model.SecurityLogModel
	UserRoleModel                 model.UserRoleModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.Addr)

	return &ServiceContext{
		Config: c,
		// 中间件
		JwtAuthMiddleware: middleware.NewJwtAuthMiddleware(c.JwtAuth.AccessSecret).Handle,
		RbacMiddleware:    middleware.NewRbacMiddleware(conn).Handle,
		Mysql:             conn,
		// nacos
		Parameters: nacosAdapter.NewClient(&nacosAdapter.Config{
			ServerConfig: &nacosAdapter.ServerConfig{
				IpAddr: c.Nacos.IpAddr,
				Port:   c.Nacos.Port,
			},
			ClientConfig: &nacosAdapter.ClientConfig{
				NamespaceId:         c.Nacos.NamespaceId,
				TimeoutMs:           c.Nacos.TimeoutMs,
				NotLoadCacheAtStart: c.Nacos.NotLoadCacheAtStart,
				LogDir:              c.Nacos.LogDir,
				CacheDir:            c.Nacos.CacheDir,
				LogLevel:            c.Nacos.LogLevel,
			},
		}),
		PriceRedis: redis.New(c.PriceRedis.Host, redis.WithPass(c.PriceRedis.Pass)),
		// Model(按字母顺序排)
		AccountsModel:                 model.NewAccountsModel(conn, c.Cache),
		AdminUserModel:                model.NewAdminUserModel(conn),
		CoinPricesModel:               model.NewCoinPricesModel(conn),
		ExchangeTradesModel:           model.NewExchangeTradesModel(conn, c.Cache),
		MembersModel:                  model.NewMembersModel(conn, c.Cache),
		MemberTransactionsModel:       model.NewMemberTransactionsModel(conn, c.Cache),
		MenuModel:                     model.NewMenuModel(conn),
		OfflineTradeInputModel:        model.NewOfflineTradeInputModel(conn, c.Cache),
		OfflineTradeInputHistoryModel: model.NewOfflineTradeInputHistoryModel(conn, c.Cache),
		ProfitsModel:                  model.NewProfitsModel(conn, c.Cache),
		SecurityLogModel:              model.NewSecurityLogModel(conn),
		UserRoleModel:                 model.NewUserRoleModel(conn),
	}
}
