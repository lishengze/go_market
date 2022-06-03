package svc

import (
	"market_server/app/client/api/internal/config"
	"market_server/app/order/model"
	"market_server/app/order/rpc/order"
	"market_server/app/userCenter/rpc/usercenter"
	"market_server/common/crypto"
	"market_server/common/nacosAdapter"
	"os"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	//nacos
	NacosClient *nacosAdapter.Client

	PemFileBase *PemFileBase64

	//rpc
	UserCenterRpc usercenter.UserCenter
	OrderRpc      order.Order
	OrderModel    model.OrderModel
}

//Pem的string格式
type PemFileBase64 struct {
	PriPemStr []byte
	PubPemStr []byte
	AESKey    string
}

func NewServiceContext(c config.Config) *ServiceContext {
	pubGatewayPem, err := os.ReadFile(c.PemFileConfig.PubPemFilePath)
	if err != nil {
		panic(err)
	}

	priBCTSPem, err := os.ReadFile(c.PemFileConfig.PriPemFilePath)
	if err != nil {
		panic(err)
	}

	svcCtx := &ServiceContext{
		Config:      c,
		NacosClient: nacosAdapter.NewClient(c.NacosConfig),
		PemFileBase: &PemFileBase64{
			PriPemStr: priBCTSPem,
			PubPemStr: pubGatewayPem,
		},

		UserCenterRpc: usercenter.NewUserCenter(zrpc.MustNewClient(c.UserCenterRpc)),
		OrderRpc:      order.NewOrder(zrpc.MustNewClient(c.OrderRpc)),
		OrderModel:    model.NewOrderModel(sqlx.NewMysql(c.Mysql.DataSource), c.Cache),
	}

	aesKey, err := crypto.GenShareKey(priBCTSPem, pubGatewayPem)
	if err != nil {
		panic(err)
	}

	svcCtx.PemFileBase.AESKey = aesKey

	return svcCtx
}
