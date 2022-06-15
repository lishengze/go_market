package core

import (
	"errors"
	"exterior-interactor/app/idsrv/rpc/idsrv"
	"exterior-interactor/app/opu/model"
	"exterior-interactor/app/opu/rpc/internal/config"
	"exterior-interactor/app/opu/rpc/opupb"
	"exterior-interactor/pkg/exchangeapi/exmodel"
	"exterior-interactor/pkg/modelext"
	"exterior-interactor/pkg/xencrypt"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/gogo/protobuf/proto"
	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sync"
	"time"
)

type (
	OPU interface {
		RegisterAccount(req *opupb.RegisterAccountReq) (*opupb.RegisterAccountRsp, error)
		UpdateAccount(in *opupb.UpdateAccountReq) (*opupb.EmptyRsp, error)
		GetAccount(req *opupb.GetAccountReq) (*opupb.GetAccountRsp, error)
		PlaceOrder(req *opupb.PlaceOrderReq) (*opupb.EmptyRsp, error)
		QueryOrder(req *opupb.QueryOrderReq) (*opupb.QueryOrderRsp, error)
		CancelOrder(req *opupb.CancelOrderReq) (*opupb.EmptyRsp, error)
		QueryBalance(req *opupb.QueryBalanceReq) (*opupb.QueryBalanceRsp, error)
		QueryTrade(req *opupb.QueryTradeReq) (*opupb.QueryTradeRsp, error)
		GetSymbol(req *opupb.GetSymbolReq) (*opupb.GetSymbolRsp, error)
	}

	opu struct {
		exchange          string
		proxy             string
		unClosedOrdersMap sync.Map // 存储未完成的订单 key: Order.Id & Order.ClientOrderId value: *orderManager
		accountMap        sync.Map // 存储账户 key:Account.Id value: *accountManager
		svcCtx            *svcCtx
		kafkaSyncProducer sarama.SyncProducer
		symbolManager     *symbolManager
		orderUpdateCh     chan *exmodel.OrderTradesUpdate // 所有 accountManager 的推送统一推送到此 chan
		outputCh          chan *opupb.OrderTradesUpdate   // 所有 orderManager 的推送统一推送到此 chan
	}

	svcCtx struct {
		model.AccountModel
		model.OrderModel
		model.TradeModel
		model.SymbolModel
		idsrv.IdSrv
	}
)

func NewOpu(c config.Config) OPU {
	svcCtx := newSvcCtx(c)
	o := &opu{
		exchange:          c.Exchange,
		proxy:             c.Proxy,
		unClosedOrdersMap: sync.Map{},
		accountMap:        sync.Map{},
		svcCtx:            svcCtx,
		kafkaSyncProducer: newKafkaSyncProducer(c),
		symbolManager:     newSymbolManager(svcCtx, c.Exchange, c.Proxy),
		orderUpdateCh:     make(chan *exmodel.OrderTradesUpdate, 10000),
		outputCh:          make(chan *opupb.OrderTradesUpdate, 10000),
	}

	err := o.loadAccounts()
	if err != nil {
		panic(err)
	}

	err = o.loadUnClosedOrders()
	if err != nil {
		panic(err)
	}

	go o.run()
	go o.removeClosedOrderManager()

	return o
}

func newKafkaSyncProducer(c config.Config) sarama.SyncProducer {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	kafkaConfig.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	kafkaConfig.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回

	client, err := sarama.NewSyncProducer([]string{c.KafkaConf.Address}, kafkaConfig)
	if err != nil {
		panic(fmt.Sprintf("create kafka client err:%v", err))
	}

	return client
}

func newSvcCtx(c config.Config) *svcCtx {
	return &svcCtx{
		SymbolModel:  model.NewSymbolModel(sqlx.NewMysql(c.Mysql.DataSource), c.Cache),
		AccountModel: model.NewAccountModel(sqlx.NewMysql(c.Mysql.DataSource), c.Cache),
		OrderModel:   model.NewOrderModel(sqlx.NewMysql(c.Mysql.DataSource), c.Cache),
		TradeModel:   model.NewTradeModel(sqlx.NewMysql(c.Mysql.DataSource), c.Cache),
		IdSrv:        idsrv.NewIdSrv(zrpc.MustNewClient(c.IdSrvRpcConf)),
	}
}

// loadAccounts 将本交易所所有账户加载到内存中
func (o *opu) loadAccounts() error {
	accounts, err := o.svcCtx.AccountModel.FindMany(modelext.NewQuery().Equal("exchange", o.exchange))
	if err != nil {
		logx.Error(err)
		return err
	}

	for _, account := range accounts {
		if _, ok := o.accountMap.Load(account.Id); !ok {
			o.loadAccountToMemory(account)
		}
	}

	return nil
}

func (o *opu) RegisterAccount(req *opupb.RegisterAccountReq) (*opupb.RegisterAccountRsp, error) {
	key, err := xencrypt.EncryptByAes(req.Key)
	if err != nil {
		return nil, err
	}
	secret, err := xencrypt.EncryptByAes(req.Secret)
	if err != nil {
		return nil, err
	}
	passphrase, err := xencrypt.EncryptByAes(req.Passphrase)
	if err != nil {
		return nil, err
	}

	account, err := o.svcCtx.AccountModel.FindOneByAlias(req.Alias)
	switch err {
	case model.ErrNotFound:
		account = &model.Account{
			Id:             o.svcCtx.IdSrv.MustGetId(),
			Alias:          req.Alias,
			Key:            key,
			Secret:         secret,
			Passphrase:     passphrase,
			Exchange:       req.Exchange,
			SubAccountName: req.SubAccountName,
		}

		_, err = o.svcCtx.AccountModel.Insert(account)

		if err != nil {
			return nil, err
		}

		o.loadAccountToMemory(account)

		return &opupb.RegisterAccountRsp{
			AccountId: account.Id,
		}, nil
	case nil:
		if key == account.Key && secret == req.Secret && passphrase == account.Passphrase && req.SubAccountName == account.SubAccountName {
			return &opupb.RegisterAccountRsp{
				AccountId: account.Id,
			}, nil
		}
		return nil, errors.New(fmt.Sprintf("alias:%s has been used. ", req.Alias))

	default:
		return nil, err
	}
}

func (o *opu) GetAccount(req *opupb.GetAccountReq) (*opupb.GetAccountRsp, error) {
	accounts, err := o.svcCtx.AccountModel.FindMany(modelext.NewQuery().
		WhereFunc(func(wheres *modelext.Wheres) {
			if req.AccountId != "" {
				wheres.Equal("id", req.AccountId)
			}

			if req.AccountAlias != "" {
				wheres.Equal("account_alias", req.AccountAlias)
			}
		}).OrderBy("create_time", modelext.DESC))
	if err != nil {
		return nil, err
	}

	var res = &opupb.GetAccountRsp{}

	for _, account := range accounts {
		res.Accounts = append(res.Accounts, &opupb.Account{
			Id:         account.Id,
			Exchange:   account.Exchange,
			Alias:      account.Alias,
			CreateTime: timestamppb.New(account.CreateTime),
			UpdateTime: timestamppb.New(account.UpdateTime),
		})
	}

	return res, nil
}

func (o *opu) UpdateAccount(in *opupb.UpdateAccountReq) (*opupb.EmptyRsp, error) {
	// todo
	return nil, fmt.Errorf("not implement")
}

func (o *opu) GetSymbol(req *opupb.GetSymbolReq) (*opupb.GetSymbolRsp, error) {
	symbols, err := o.svcCtx.SymbolModel.FindMany(modelext.NewQuery().WhereFunc(
		func(wheres *modelext.Wheres) {
			if req.Symbol != "" {
				wheres.Equal("std_symbol", req.Symbol)
			}
			if req.Exchange != "" {
				wheres.Equal("exchange", req.Exchange)

			}
		}).OrderBy("create_time", modelext.DESC))

	if err != nil {
		logx.Error(err)
		return nil, err
	}

	return &opupb.GetSymbolRsp{Symbols: toPbSymbols(symbols)}, nil
}

func (o *opu) QueryBalance(req *opupb.QueryBalanceReq) (*opupb.QueryBalanceRsp, error) {
	account, err := o.getAccountManager(req.AccountId, req.AccountAlias)
	if err != nil {
		return nil, err
	}

	// todo 校验 wallet type
	rsp, err := account.WalletManager.QueryBalance(exmodel.QueryBalanceReq{
		WalletType: exmodel.WalletType(req.WalletType),
	})

	if err != nil {
		logx.Errorf("QueryBalance err:%v ,req:%v", err, req)
		return nil, err
	}

	return &opupb.QueryBalanceRsp{Balances: toPbBalances(rsp.Balances)}, nil
}

func (o *opu) QueryTrade(req *opupb.QueryTradeReq) (*opupb.QueryTradeRsp, error) {
	//account, err := o.getAccountManager(req.AccountId, req.AccountAlias)
	//if err != nil {
	//	return nil, err
	//}

	//account.QueryTrades()

	return nil, fmt.Errorf("not implement")
}

func (o *opu) QueryOrder(req *opupb.QueryOrderReq) (*opupb.QueryOrderRsp, error) {
	account, err := o.getAccountManager(req.AccountId, req.AccountAlias)
	if err != nil {
		return nil, err
	}

	order, err := o.svcCtx.OrderModel.FindOneByAccountIdClientOrderId(account.Account.Id, req.ClientOrderId)
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	om, ok := o.unClosedOrdersMap.Load(order.Id)
	if ok {
		om.(*orderManager).syncOrder() // 主动向交易所同步一次订单
	}

	order, err = o.svcCtx.OrderModel.FindOneByAccountIdClientOrderId(account.Account.Id, req.ClientOrderId)
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	trades, err := o.svcCtx.TradeModel.FindMany(modelext.NewQuery().Equal("order_id", order.Id))
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	return &opupb.QueryOrderRsp{
		Order:  toPbOrder(order),
		Trades: toPbTrades(trades),
	}, nil
}

func (o *opu) PlaceOrder(req *opupb.PlaceOrderReq) (*opupb.EmptyRsp, error) {
	account, err := o.getAccountManager(req.AccountId, req.AccountAlias)
	if err != nil {
		return nil, err
	}

	err = o.verifyPlaceOrderParams(req, account.Account.Id)
	if err != nil {
		return nil, err
	}

	symbol, err := o.svcCtx.SymbolModel.FindOneByStdSymbolExchange(req.StdSymbol, req.Exchange)
	switch err {
	case model.ErrNotFound:
		return nil, fmt.Errorf("wrong symbol or exchange")
	case nil:
	default:
		logx.Error(err)
		return nil, err
	}

	order := &model.Order{
		Id:            o.svcCtx.IdSrv.MustGetId(),
		AccountId:     account.Account.Id,
		ClientOrderId: req.ClientOrderId,
		AccountAlias:  account.Alias,
		ExOrderId:     "",
		ApiType:       symbol.ApiType,
		Side:          req.Side,
		Status:        exmodel.OrderStatusPending.String(),
		Volume:        req.Volume,
		FilledVolume:  "0",
		Price:         req.Price,
		Tp:            req.Type,
		CancelFlag:    "UNSET",
		SendFlag:      "UNSENT",
		StdSymbol:     symbol.StdSymbol,
		ExSymbol:      symbol.ExFormat,
		Exchange:      req.Exchange,
	}

	_, err = o.svcCtx.OrderModel.Insert(order)
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	order, _ = o.svcCtx.OrderModel.FindOne(order.Id)

	om := newOrderManager(func(order *model.Order) (*accountManager, error) {
		return o.getAccountManager(order.AccountId, "")
	}, o.svcCtx, o.outputCh, symbol, order, []*model.Trade{})

	o.unClosedOrdersMap.Store(order.Id, om)

	return &opupb.EmptyRsp{}, nil
}

func (o *opu) CancelOrder(req *opupb.CancelOrderReq) (*opupb.EmptyRsp, error) {

	account, err := o.getAccountManager(req.AccountId, req.AccountAlias)
	if err != nil {
		return nil, err
	}

	order, err := o.svcCtx.OrderModel.FindOneByAccountIdClientOrderId(account.Account.Id, req.ClientOrderId)
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	om, ok := o.unClosedOrdersMap.Load(order.Id)
	if ok {
		om.(*orderManager).cancelOrder() // 撤单
		return &opupb.EmptyRsp{}, nil
	}

	return &opupb.EmptyRsp{}, fmt.Errorf("order is closed")
}

func (o *opu) verifyPlaceOrderParams(req *opupb.PlaceOrderReq, accountId string) error {
	if req.Side != exmodel.OrderSideBuy.String() && req.Side != exmodel.OrderSideSell.String() {
		return fmt.Errorf("wrong side:%s", req.Side)
	}

	if req.Type != exmodel.OrderTypeLimit.String() && req.Type != exmodel.OrderTypeMarket.String() {
		return fmt.Errorf("wrong type:%s", req.Type)
	}

	_, err := o.svcCtx.OrderModel.FindOneByAccountIdClientOrderId(accountId, req.ClientOrderId)
	switch err {
	case model.ErrNotFound:
	case nil:
		return fmt.Errorf("Duplicated ClientOrderId:%s. ", req.ClientOrderId)
	default:
		logx.Error(err)
		return err
	}

	_, err = decimal.NewFromString(req.Price)
	if err != nil {
		return fmt.Errorf("wrong price:%s", req.Price)
	}

	_, err = decimal.NewFromString(req.Volume)
	if err != nil {
		return fmt.Errorf("wrong volume:%s", req.Volume)
	}

	return nil
}

func (o *opu) getAccountManager(accountId, accountAlias string) (*accountManager, error) {
	var (
		manager       *accountManager
		realAccountId string
	)
	if accountId != "" {
		realAccountId = accountId
	} else {
		account, err := o.svcCtx.AccountModel.FindOneByAlias(accountAlias)
		if err != nil {
			logx.Error(err)
			return nil, err
		}
		realAccountId = account.Id
	}

	res, ok := o.accountMap.Load(realAccountId)
	if !ok {
		return nil, fmt.Errorf("wrong accountId:%s", accountId)
	}
	manager = res.(*accountManager)

	return manager, nil
}

// loadAccountToMemory 加载账户到内存中，并且监听订单推送
func (o *opu) loadAccountToMemory(account *model.Account) {
	am, err := newAccountManager(account, o.proxy, o.orderUpdateCh)
	if err != nil {
		logx.Error(err)
		return
	}

	o.accountMap.Store(account.Id, am)
}

// removeClosedOrderManager 定期清除已关闭的order
func (o *opu) removeClosedOrderManager() {
	for {
		o.unClosedOrdersMap.Range(func(key, value interface{}) bool {
			om := value.(*orderManager)
			if om.orderIsClosed() {
				o.unClosedOrdersMap.Delete(key)
			}
			return true
		})

		time.Sleep(time.Minute)
	}
}

// loadUnClosedOrders 加载未完成的订单
func (o *opu) loadUnClosedOrders() error {
	orders, err := o.svcCtx.OrderModel.FindMany(modelext.NewQuery().
		NotEqual("status", exmodel.OrderStatusFilled).
		NotEqual("status", exmodel.OrderStatusRejected).
		NotEqual("status", exmodel.OrderStatusCancelled))

	if err != nil {
		logx.Error(err)
		return err
	}

	for _, order := range orders {
		symbol, err := o.svcCtx.SymbolModel.FindOneByStdSymbolExchange(order.StdSymbol, order.Exchange)
		if err != nil {
			logx.Error(err)
			return err
		}

		om := newOrderManagerWithoutSend(func(order *model.Order) (*accountManager, error) {
			return o.getAccountManager(order.AccountId, "")
		}, o.svcCtx, o.outputCh, symbol, order, []*model.Trade{})

		o.unClosedOrdersMap.Store(order.Id, om)
	}

	return nil
}

// run
func (o *opu) run() {
	for {
		select {
		case update := <-o.orderUpdateCh: // 分发 各个 accountManager 的 update 到 相应的 orderManager
			om, ok := o.unClosedOrdersMap.Load(update.ClientOrderId)
			if !ok {
				logx.Errorf("can't find orderManager, update:%+v", *update)
				continue
			}
			om.(*orderManager).inputUpdate(update)

		case update := <-o.outputCh: // 将信息推送到 kafka
			bytes, err := proto.Marshal(update)
			if err != nil {
				logx.Errorf("convert update err:%v, update:%v ", err, update)
				continue
			}

			var (
				kafkaTopic = fmt.Sprintf("ORDER.%s", update.AccountAlias)
				msg        = &sarama.ProducerMessage{
					Topic: kafkaTopic,
					Value: sarama.ByteEncoder(bytes),
				}
			)

			_, _, err = o.kafkaSyncProducer.SendMessage(msg)

			if err != nil {
				logx.Errorf("kafkaConn write update err:%s, update:%v", err, update)
			}
		}
	}
}
