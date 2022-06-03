package svc

import (
	"bcts/app/dataService/rpc/dataservice"
	"bcts/app/order/model"
	"bcts/app/order/rpc/internal/config"
	"bcts/app/userCenter/rpc/usercenter"
	"bcts/common/nacosAdapter"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config      config.Config
	NacosClient *nacosAdapter.Client
	RedisConn   *redis.Redis

	//rpc
	UserCenterRpc  usercenter.UserCenter
	DataServiceRpc dataservice.DataService

	//model
	UserGroupModel model.UserGroupModel
	OrderModel     model.OrderModel
	TradeModel     model.TradeModel
	TradeFlowModel model.TradeFlowModel
	AccountModel   model.AccountModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{

		Config:         c,
		NacosClient:    nacosAdapter.NewClient(c.Nacos),
		UserGroupModel: model.NewUserGroupModel(conn, c.Cache),
		RedisConn: redis.New(c.Redis.Host, func(r *redis.Redis) {
			r.Pass = c.Redis.Pass
			r.Type = c.Redis.Type
		}),

		//rpc
		UserCenterRpc:  usercenter.NewUserCenter(zrpc.MustNewClient(c.UserCenterRpc)),
		DataServiceRpc: dataservice.NewDataService(zrpc.MustNewClient(c.DataServiceRpc)),

		//model
		OrderModel:     model.NewOrderModel(sqlx.NewMysql(c.Mysql.DataSource), c.Cache),
		TradeModel:     model.NewTradeModel(sqlx.NewMysql(c.Mysql.DataSource), c.Cache),
		TradeFlowModel: model.NewTradeFlowModel(sqlx.NewMysql(c.Mysql.DataSource), c.Cache),
		AccountModel:   model.NewAccountModel(sqlx.NewMysql(c.Mysql.DataSource), c.Cache),
	}
}
