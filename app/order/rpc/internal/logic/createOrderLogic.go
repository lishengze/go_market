package logic

import (
	"bcts/app/dataService/rpc/dataservice"
	"bcts/app/order/rpc/internal/ordercheck"
	"bcts/app/userCenter/rpc/usercenter"
	"bcts/common/globalKey"
	"bcts/common/utils"
	"bcts/common/xerror"
	"bcts/pkg/kafkaclient"
	"bcts/pkg/kafkaclient/mpupb"
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strconv"
	"time"

	"bcts/app/order/model"
	"bcts/app/order/rpc/internal/svc"
	"bcts/app/order/rpc/types/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateOrderLogic) CreateOrder(in *pb.CreateOrderReq) (rsp *pb.CreateOrderRsp, err error) {
	l.Logger.Info("CreateOrder in:%+v", in)

	//1. 报单参数校验
	if err = l.orderParamVerify(in); err != nil {
		return nil, err
	}

	//2. 用户信息校验
	userInfo, err := l.svcCtx.UserCenterRpc.GetUserInfo(l.ctx, &usercenter.UserInfoReq{UserID: in.UserID})
	if err != nil {
		return nil, err
	}

	//3. 用户手续费获取
	feeInfo, err := l.svcCtx.DataServiceRpc.GetUserFeeInfo(l.ctx, &dataservice.GetUserFeeInfoReq{
		UserID:  in.UserID,
		GroupID: userInfo.GroupID,
		Symbol:  in.Symbol,
	})
	if err != nil {
		return nil, err
	}

	//4. 精度检验
	volume, amount, amountExceptFee, err := ordercheck.CheckUniversal(in.OrderMode, in.Volume, in.Amount, feeInfo)
	if err != nil {
		return nil, err
	}
	logx.Infof("amountExceptFee:%s", amountExceptFee.String())

	//TODO 5. 检验限额, 开关，阈值

	var (
		price decimal.Decimal
		fee   decimal.Decimal //预估手续费(对于OTC订单，此预估手续费就是实际手续费)
	)
	//TODO: OTC订单不论是按volume/amount下单，fee预算值就是最终值
	//对于OTC订单而言，一定有price；volume和amount参数二选一
	takerFee, _ := decimal.NewFromString(feeInfo.TakerFee)
	if in.OrderType == globalKey.OrderTypeOTC { //订单类型 1:otc订单 2:聚合交易订单
		price, err = decimal.NewFromString(in.Price)
		if err != nil {
			err = errors.Wrapf(xerror.ErrorParamError, "params price:%s error, err:%+v", in.Price, err)
			return
		}
		if in.OrderMode == globalKey.OrderModeVolume { //Volume.GreaterThan(decimal.Zero) {
			if feeInfo.FeeKind == globalKey.FeeKindPercentage { //1表示百比分，2表示绝对值
				fee = volume.Mul(price).Mul(takerFee)
			} else {
				fee = takerFee
			}
			amount = volume.Mul(price).Add(fee) //预估amount
		} else if in.OrderMode == globalKey.OrderModeAmount {
			//如果是总额询价,要先扣除手续费
			if feeInfo.FeeKind == globalKey.FeeKindPercentage { //1表示百比分，2表示绝对值
				fee = amount.Mul(takerFee)
			} else {
				fee = takerFee
			}
			volume = amount.Sub(fee).Div(price) //预估volume

			//重新计算fee & amount，避免精度问题带来误差
			if feeInfo.FeeKind == globalKey.FeeKindPercentage { //1表示百比分，2表示绝对值
				fee = volume.Mul(price).Mul(takerFee)
			} else {
				fee = takerFee
			}
			amount = volume.Mul(price).Add(fee)
		} else {
			err = xerror.ErrorParamError
			return
		}
	} else {
		err = xerror.ErrorParamError
		return
	}
	//TODO: 应用精度
	orderLocalID := utils.GenOrderID()
	ts := time.Now()

	order := &model.Order{
		UserId:          in.UserID,
		OrderLocalId:    orderLocalID,
		Symbol:          in.Symbol,
		BaseCurrency:    feeInfo.BaseCurrency,
		TargetCurrency:  feeInfo.TargetCurrency,
		Direction:       int64(in.Direction),
		OrderType:       int64(in.OrderType),
		OrderMode:       int64(in.OrderMode),
		OrderPriceType:  int64(in.OrderPriceType),
		Price:           price,
		Volume:          volume,
		Amount:          amount,
		OrderMaker:      globalKey.MarketTaker, //目前统一写 taker // 1:taker 2:maker
		TradeAmount:     decimal.Zero,
		TradeVolume:     decimal.Zero,
		FeeKind:         int64(feeInfo.FeeKind),
		FeeRate:         takerFee,
		OrderStatus:     globalKey.OrderStatusSent,
		OrderCreateTime: ts,
		OrderModifyTime: ts,
	}
	logx.Infof("new order:%+v", order)

	//手续费，买入订单收取coin,卖出订单收取baseCoin
	var frozenCurrency string                //需要冻结的currency
	var frozenCurrencyAmount decimal.Decimal //需要冻结的currency的数量/金额
	var otherCurrency string
	//var otherCurrencyAmount decimal.Decimal
	if order.Direction == globalKey.DirectionBuy {
		frozenCurrency = feeInfo.BaseCurrency
		frozenCurrencyAmount = order.Amount
		otherCurrency = feeInfo.TargetCurrency
		//otherCurrencyAmount = order.Volume
	} else {
		frozenCurrency = feeInfo.TargetCurrency
		frozenCurrencyAmount = order.Volume
		otherCurrency = feeInfo.BaseCurrency
		//otherCurrencyAmount = order.Amount
	}
	logx.Infof("frozen currency:%s, other currency:%s", frozenCurrency, otherCurrency)
	whereBuilderForAccount := l.svcCtx.AccountModel.RowBuilder()
	whereBuilderForAccount = whereBuilderForAccount.Where(squirrel.Eq{"user_id": in.UserID, "currency": frozenCurrency})
	frozenAccount, err := l.svcCtx.AccountModel.FindOneByQuery(l.ctx, whereBuilderForAccount)
	if err != nil {
		logx.Error(err)
		if err == model.ErrNotFound {
			//用户没有该币种账户，余额不足
			return nil, errors.Wrapf(xerror.ErrorUserAccountNotEnough, "balance=0. user %s account not exists. params :%+v", frozenCurrency, in)
		}
		return nil, errors.Wrapf(xerror.ErrorDB, "get user %s account err: %+v , params :%+v", frozenCurrency, err, in)
	}
	whereBuilder := l.svcCtx.AccountModel.RowBuilder().Where(squirrel.Eq{"user_id": in.UserID, "currency": otherCurrency})
	otherAccount, err := l.svcCtx.AccountModel.FindOneByQuery(l.ctx, whereBuilder)
	if err != nil {
		if err == model.ErrNotFound {
			err = nil
			otherAccount = &model.Account{
				UserId:     in.UserID,
				Currency:   otherCurrency,
				CurrencyId: "", //TODO:
				Frozen:     decimal.Zero,
				Balance:    decimal.Zero,
			}
			_, e := l.svcCtx.AccountModel.Insert(l.ctx, nil, otherAccount)
			if e != nil {
				err = e
				logx.Error(err)
				return
			}
		} else {
			logx.Error(err)
			return nil, errors.Wrapf(xerror.ErrorDB, "get user %s account err: %+v , params :%+v", otherCurrency, err, in)
		}
	}
	if frozenAccount.Balance.LessThan(order.Amount) {
		err = errors.Wrapf(xerror.ErrorUserAccountNotEnough, "user %s account balance not enough. balance:%s, params :%+v", frozenCurrency, frozenAccount.Balance.String(), in)
		return nil, err
	}
	err = l.svcCtx.OrderModel.Trans(l.ctx, func(context context.Context, session sqlx.Session) error {
		// account 冻结
		//TODO: 冻结资金 需要加锁
		frozenAccount.Frozen = frozenAccount.Frozen.Add(frozenCurrencyAmount)
		//account.Available = account.Available.Sub(order.Amount)
		_, err = l.svcCtx.AccountModel.Update(l.ctx, nil, frozenAccount)
		if err != nil {
			logx.Error(err)
			return err
		}
		//冻结资金--交易流水
		tradeFlow := &model.TradeFlow{
			UserId:         order.UserId,
			AccountId:      0, //TODO: 暂定忽略此字段
			Currency:       frozenAccount.Currency,
			Source:         1, // 1:冻结 2:解冻 3:成交 4:线下交易 5:法币出入金  最小原子状态
			Type:           1, // 1:交易 2:交易手续费 3:交易修正 4:交易手续费修正  前台筛选用
			OrderLocalId:   order.OrderLocalId,
			TradeId:        "",
			OfflineTradeId: "",
			Amount:         frozenCurrencyAmount,
			FlowCreateTime: time.Now(),
		}
		_, err = l.svcCtx.TradeFlowModel.Insert(l.ctx, nil, tradeFlow)
		if err != nil {
			logx.Error(err)
			return err
		}

		// order
		insertResult, err := l.svcCtx.OrderModel.Insert(l.ctx, nil, order)
		if err != nil {
			logx.Error(err)
			return err
		}
		lastId, err := insertResult.LastInsertId()
		if err != nil {
			return errors.Wrapf(xerror.ErrorDB, "create order insertResult.LastInsertId err:%v,order:%+v", err, order)
		}
		//order.Id = lastId
		order, _ = l.svcCtx.OrderModel.FindOne(l.ctx, lastId) //refresh
		if err != nil {
			return errors.Wrapf(xerror.ErrorDB, "refresh order err:%v", err)
		}
		return nil
	})
	if err != nil {
		logx.Error(err)
		return
	}

	logx.Infof("created order_local_id:%s, order_id:%d, order:%+v", order.OrderLocalId, order.Id, order)
	//trades
	err = l.TradesHandle(order, fee, frozenAccount, otherAccount)
	if err != nil {
		logx.Error(err)
		return
	}

	return &pb.CreateOrderRsp{
		RequestID:      in.RequestID,
		OrderLocalID:   orderLocalID,
		UserID:         in.UserID,
		Symbol:         in.Symbol,
		OrderType:      in.OrderType,
		OrderMode:      in.OrderMode,
		OrderPriceType: in.OrderPriceType,
		Direction:      in.Direction,
		Volume:         in.Volume,
		Amount:         in.Amount,
		Price:          in.Price,
	}, nil

}

//报单参数校验 通用
func (l *CreateOrderLogic) orderParamVerify(in *pb.CreateOrderReq) error {

	volume, _ := decimal.NewFromString(in.Volume)
	amount, _ := decimal.NewFromString(in.Amount)

	//1. 下单数据校验
	if in.OrderType == globalKey.OrderTypeOTC {
		//价格校验
		userIdStr := strconv.FormatInt(in.UserID, 10)
		redisKey := userIdStr + globalKey.OTCRedisPriceSuffix
		quotePrice, err := l.svcCtx.RedisConn.Get(redisKey)
		if err != nil {
			return errors.Wrapf(xerror.ErrorOrderPriceExpired, "paramVerify redis price get err:%+v", err)
		}
		if quotePrice != in.Price {
			return errors.Wrapf(xerror.ErrorOTCOrderPrice, "paramVerify redis price:%s, orderPrice:%s", quotePrice, in.Price)
		}

		if in.OrderMode == globalKey.OrderModeVolume && volume.IsPositive() {
			return errors.Wrapf(xerror.ErrorOrderVolumeInvalid, "paramVerify error")
		}

		if in.OrderMode == globalKey.OrderModeAmount && amount.IsPositive() {
			return errors.Wrapf(xerror.ErrorOrderVolumeInvalid, "paramVerify error")
		}

		return nil
	}
	return nil
}

func ParamHandle(in *pb.CreateOrderReq) (err error) {
	// 用户维度限制
	// 品种限制
	// 币种限制

	return nil
}

func (l *CreateOrderLogic) TradesHandle(order *model.Order, fee decimal.Decimal, frozenAccount *model.Account, otherAccount *model.Account) (err error) {
	// 根据行情撮合
	// 生成成交记录(insert trade(s))
	// 解冻资金
	outerTrades, updatedOrder, err := trade(order)
	if err != nil {
		logx.Error(err)
		return
	}
	//logx.Info(updatedOrder)
	logx.Infof("insert trade count:%d", len(outerTrades))
	for _, trade := range outerTrades {
		_, err = l.svcCtx.TradeModel.Insert(l.ctx, nil, trade)
		if err != nil {
			logx.Error(err)
			return
		}
	}
	//成交（资金变动+流水记录）
	if order.Direction == 0 { //buy BTC_USDT: +BTC -USDT -USDT(fee)
		//+BTC
		otherAccount.Balance = otherAccount.Balance.Add(order.Volume)
		_, err = l.svcCtx.AccountModel.Update(l.ctx, nil, otherAccount)
		if err != nil {
			logx.Error(err)
			return
		}
		//流水
		if err = l.addTradeFlow(&model.TradeFlow{
			UserId:         order.UserId,
			AccountId:      0, //TODO: 暂定忽略此字段
			Currency:       otherAccount.Currency,
			Source:         3, // 1:冻结 2:解冻 3:成交 4:线下交易 5:法币出入金  最小原子状态
			Type:           1, // 1:交易 2:交易手续费 3:交易修正 4:交易手续费修正  前台筛选用
			OrderLocalId:   order.OrderLocalId,
			TradeId:        "",
			OfflineTradeId: "",
			Amount:         order.Volume,
			FlowCreateTime: time.Now(),
		}); err != nil {
			logx.Error(err)
			return
		}

		//-USDT -USDT(fee)  //卖出的USDT+手续费USDT一起扣除
		frozenAccount.Balance = frozenAccount.Balance.Sub(order.Amount)
		frozenAccount.Frozen = frozenAccount.Frozen.Sub(order.Amount)
		_, err = l.svcCtx.AccountModel.Update(l.ctx, nil, frozenAccount)
		if err != nil {
			logx.Error(err)
			return
		}
		//流水则要两条（不能合并记录）
		if err = l.addTradeFlow(&model.TradeFlow{
			UserId:         order.UserId,
			AccountId:      0, //TODO: 暂定忽略此字段
			Currency:       frozenAccount.Currency,
			Source:         3, // 1:冻结 2:解冻 3:成交 4:线下交易 5:法币出入金  最小原子状态
			Type:           1, // 1:交易 2:交易手续费 3:交易修正 4:交易手续费修正  前台筛选用
			OrderLocalId:   order.OrderLocalId,
			TradeId:        "",
			OfflineTradeId: "",
			Amount:         order.Volume.Mul(order.Price).Neg(),
			FlowCreateTime: time.Now(),
		}); err != nil {
			logx.Error(err)
			return
		}
		if err = l.addTradeFlow(&model.TradeFlow{
			UserId:         order.UserId,
			AccountId:      0, //TODO: 暂定忽略此字段
			Currency:       frozenAccount.Currency,
			Source:         3, // 1:冻结 2:解冻 3:成交 4:线下交易 5:法币出入金  最小原子状态
			Type:           2, // 1:交易 2:交易手续费 3:交易修正 4:交易手续费修正  前台筛选用
			OrderLocalId:   order.OrderLocalId,
			TradeId:        "",
			OfflineTradeId: "",
			Amount:         fee.Neg(),
			FlowCreateTime: time.Now(),
		}); err != nil {
			logx.Error(err)
			return
		}
	} else { //sell BTC_USDT: -BTC +USDT -USDT(fee)
		//-BTC(BTC已经冻结了）
		frozenAccount.Balance = frozenAccount.Balance.Sub(order.Volume)
		frozenAccount.Frozen = frozenAccount.Frozen.Sub(order.Volume)
		_, err = l.svcCtx.AccountModel.Update(l.ctx, nil, frozenAccount)
		if err != nil {
			logx.Error(err)
			return
		}
		//流水
		if err = l.addTradeFlow(&model.TradeFlow{
			UserId:         order.UserId,
			AccountId:      0, //TODO: 暂定忽略此字段
			Currency:       frozenAccount.Currency,
			Source:         3, // 1:冻结 2:解冻 3:成交 4:线下交易 5:法币出入金  最小原子状态
			Type:           1, // 1:交易 2:交易手续费 3:交易修正 4:交易手续费修正  前台筛选用
			OrderLocalId:   order.OrderLocalId,
			TradeId:        "",
			OfflineTradeId: "",
			Amount:         order.Volume.Neg(),
			FlowCreateTime: time.Now(),
		}); err != nil {
			logx.Error(err)
			return
		}

		//+USDT -USDT(fee) //扣除（获得的USDT-手续费USDT）
		diffValue := order.Volume.Mul(order.Price).Sub(fee)
		//if diffValue.LessThanOrEqual(decimal.Zero){
		//	//TODO:error 赚的的USDT不足扣手续费
		//}
		otherAccount.Balance = otherAccount.Balance.Add(diffValue)
		_, err = l.svcCtx.AccountModel.Update(l.ctx, nil, otherAccount)
		if err != nil {
			logx.Error(err)
			return
		}
		//流水则要两条（不能合并记录）
		if err = l.addTradeFlow(&model.TradeFlow{
			UserId:         order.UserId,
			AccountId:      0, //TODO: 暂定忽略此字段
			Currency:       otherAccount.Currency,
			Source:         3, // 1:冻结 2:解冻 3:成交 4:线下交易 5:法币出入金  最小原子状态
			Type:           1, // 1:交易 2:交易手续费 3:交易修正 4:交易手续费修正  前台筛选用
			OrderLocalId:   order.OrderLocalId,
			TradeId:        "",
			OfflineTradeId: "",
			Amount:         order.Volume.Mul(order.Price),
			FlowCreateTime: time.Now(),
		}); err != nil {
			logx.Error(err)
			return
		}
		if err = l.addTradeFlow(&model.TradeFlow{
			UserId:         order.UserId,
			AccountId:      0, //TODO: 暂定忽略此字段
			Currency:       otherAccount.Currency,
			Source:         3, // 1:冻结 2:解冻 3:成交 4:线下交易 5:法币出入金  最小原子状态
			Type:           2, // 1:交易 2:交易手续费 3:交易修正 4:交易手续费修正  前台筛选用
			OrderLocalId:   order.OrderLocalId,
			TradeId:        "",
			OfflineTradeId: "",
			Amount:         fee.Neg(),
			FlowCreateTime: time.Now(),
		}); err != nil {
			logx.Error(err)
			return
		}
	}

	//解冻资金（上文已经解冻，本次只用记录流水）
	//交易流水表(解冻资金)
	if err = l.addTradeFlow(&model.TradeFlow{
		UserId:         order.UserId,
		AccountId:      0, //TODO: 暂定忽略此字段
		Currency:       frozenAccount.Currency,
		Source:         2, // 1:冻结 2:解冻 3:成交 4:线下交易 5:法币出入金  最小原子状态
		Type:           1, // 1:交易 2:交易手续费 3:交易修正 4:交易手续费修正  前台筛选用
		OrderLocalId:   order.OrderLocalId,
		TradeId:        "",
		OfflineTradeId: "",
		Amount:         order.Amount,
		FlowCreateTime: time.Now(),
	}); err != nil {
		logx.Error(err)
		return
	}

	order.OrderStatus = 1 // 订单状态 0:已发送 1:全部成交 2:部分成交在队列 3:部分成交已撤单 4:撤单
	order.TradeVolume = updatedOrder.TradeVolume
	order.TradeAmount = updatedOrder.TradeAmount
	order.OrderModifyTime = time.Now()
	_, err = l.svcCtx.OrderModel.Update(l.ctx, nil, order)
	if err != nil {
		logx.Error(err)
		return
	}

	return
}

func trade(order *model.Order) (outerTrades []*model.Trade, updatedOrder *model.Order, e error) {
	if order.Volume.Cmp(decimal.Zero) <= 0 || order.Volume.Cmp(order.TradeVolume) <= 0 {
		logx.Error("Trade return: order.Volume <=0")
		e = errors.New("Trade return: order.Volume <=0")
		return
	}

	//if order.OrderType == 1 { //"OTC"
	priceVolumeList, err := getDepthData(order.Symbol, order.Direction, order.Volume, order.Amount)
	if err != nil {
		logx.Error(err)
		e = err
		return
	}
	return handleOTCTrade(*order, priceVolumeList)
	//}
	//return
}

func getDepthData(symbol string, direction int64, volume, amount decimal.Decimal) (priceVolumeList []*mpupb.PriceVolume, err error) {
	logx.Infof("depth len:%d", len(kafkaclient.DepthDataList))
	if len(kafkaclient.DepthDataList) > 0 {
		for k, v := range kafkaclient.DepthDataList {
			//logx.Info(k) //FTX:ETH_USDT
			//	logx.Info(v.Timestamp)
			//Exchange:FTX, Symbol:ETH_USDT
			logx.Infof("Key:%s, Exchange:%s, Symbol:%s, len(Asks):%d, len(Bids):%d", k, v.Exchange, v.Symbol, len(v.Asks), len(v.Bids))
		}
		if symbolDepthData, ok := kafkaclient.DepthDataList[symbol]; ok {
			if direction == 0 { //用户买 BUY
				priceVolumeList = symbolDepthData.Asks
			} else { //用户卖 SELL
				priceVolumeList = symbolDepthData.Bids
			}
			return
		}
	}
	err = errors.New("no DepthData") //TODO:...
	return
}

func handleOTCTrade(order model.Order, priceVolumeList []*mpupb.PriceVolume) (outerTrades []*model.Trade, updatedOrder *model.Order, e error) {
	if order.Price.Cmp(decimal.Zero) <= 0 {
		logx.Info("handleOTCTrade error return , order.Price <= 0")
		e = errors.New("handleOTCTrade error return , order.Price <= 0")
		return
	}
	logx.Info("handleOTCTrade order:%+v, priceVolumeList:%+v", order, priceVolumeList)
	updatedOrder = &order

	for _, v := range priceVolumeList {
		priceInDepth, _ := decimal.NewFromString(v.Price)
		volumeInDepth, _ := decimal.NewFromString(v.Volume)
		marketDepth := &MarketDepth{
			Price:        priceInDepth,
			Volume:       volumeInDepth,
			TradedVolume: decimal.Zero,
			//Fee:          decimal.Zero,
		}
		trade, err := processMatchDepth(updatedOrder, marketDepth)
		if err != nil {
			logx.Error(err)
			continue
		}
		logx.Infof("gen trade:%+v", trade)
		outerTrades = append(outerTrades, trade)
		//判断交易单是否完成
		if orderTradedIsCompleted(updatedOrder) {
			break
		}
	}

	return
}

func orderTradedIsCompleted(o *model.Order) bool {
	return o.Volume.LessThanOrEqual(o.TradeVolume) // (order.volume <= order.trade_volume)
}

/**
 * 处理委托订单和匹配的行情
 * @param order 当前订单
 * @param matchDepth 匹配行情
 */
func processMatchDepth(order *model.Order, matchDepth *MarketDepth) (trade *model.Trade, err error) {
	//需要交易的数量，成交量, 成交价，可用数量
	var (
		dealPrice, tradedVolume, turnover, fee decimal.Decimal
	)

	tradedVolume, dealPrice, turnover, fee, err = calculateTradedVolume(order, matchDepth)
	if err != nil {
		logx.Error("processMatchDepth calculateTradedVolume err(%+v)", err)
		return
	}
	logx.Info("processMatchDepth tradedVolume:%v dealPrice:%v turnover:%v fee:%v",
		tradedVolume.String(), dealPrice.String(), turnover.String(), fee.String())

	//matchDepth.AddTradedAmount(tradedVolume).AddTurnover(turnover).AddFee(fee)
	matchDepth.TradedVolume = matchDepth.TradedVolume.Add(tradedVolume)
	//matchDepth.Turnover = matchDepth.Turnover.Add(turnover)
	//matchDepth.Fee = matchDepth.Fee.Add(fee)

	//order.AddTradedAmount(tradedVolume).AddTurnover(turnover.Add(fee))
	order.TradeVolume = order.TradeVolume.Add(tradedVolume)
	order.TradeAmount = order.TradeAmount.Add(turnover).Add(fee)

	//创建成交记录
	ts := time.Now()
	trade = &model.Trade{
		UserId:              order.UserId,
		OrderLocalId:        order.OrderLocalId,
		TradeId:             utils.GenTradeID(),
		OrderType:           int64(order.OrderType),
		OrderMode:           int64(order.OrderMode),
		OrderPriceType:      int64(order.OrderPriceType),
		Symbol:              order.Symbol,
		BaseCurrency:        order.BaseCurrency,
		TargetCurrency:      order.TargetCurrency,
		Direction:           int64(order.Direction),
		Price:               order.Price,
		Volume:              order.Volume,
		Amount:              order.Amount,
		OrderMaker:          order.OrderMaker,
		TradeVolume:         tradedVolume,
		TradeAmount:         turnover.Add(fee),
		FeeKind:             order.FeeKind,
		FeeRate:             order.FeeRate,
		Source:              1, //成交来源 1:otc
		OfflineTradeId:      "",
		Fee:                 fee,
		TurnOver:            turnover,
		BaseCurrencyPrice:   decimal.Zero, //pms, 基础币种成交时的价格 TODO:...
		TargetCurrencyPrice: decimal.Zero, //pms, 目标币种成交时的价格 TODO:...
		TradeTime:           ts,
	}

	return
}

func calculateTradedVolume(order *model.Order, matchDepth *MarketDepth) (tradedVolume, dealPrice, turnover, fee decimal.Decimal, err error) {
	var needVolume decimal.Decimal

	if order.OrderType == 1 { //OTC
		dealPrice = order.Price
	}
	needVolume = order.Volume.Sub(order.TradeVolume)
	tradedVolume, err = getTradedVolume(matchDepth.GetRestVolume(), needVolume)
	turnover = tradedVolume.Mul(dealPrice)
	if order.FeeKind == 1 { //百分比
		fee = tradedVolume.Mul(dealPrice).Mul(order.FeeRate)
	} else { //绝对值
		fee = tradedVolume.Div(order.Volume).Mul(order.FeeRate)
	}

	return
}

func getTradedVolume(restDepthVolume, needVolume decimal.Decimal) (tradedVolume decimal.Decimal, err error) {
	//计算成交量
	if restDepthVolume.Cmp(needVolume) >= 0 {
		tradedVolume = needVolume
	} else {
		tradedVolume = restDepthVolume
	}
	// 如果成交额为0说明剩余额度无法成交，退出
	if tradedVolume.Cmp(decimal.Zero) <= 0 {
		err = errors.New("traded volume error")
	}
	return
}

func (l *CreateOrderLogic) addTradeFlow(tradeFlow *model.TradeFlow) (err error) {
	_, err = l.svcCtx.TradeFlowModel.Insert(l.ctx, nil, tradeFlow)
	if err != nil {
		logx.Error(err)
		return
	}
	return
}
