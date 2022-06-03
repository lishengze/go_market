package logic

import (
	"bcts/app/dataService/rpc/dataservice"
	"bcts/app/order/model"
	"bcts/app/order/rpc/internal/ordercheck"
	"bcts/app/order/rpc/internal/svc"
	"bcts/app/order/rpc/types/pb"
	"bcts/app/userCenter/rpc/usercenter"
	"bcts/common/globalKey"
	"bcts/common/nacosAdapter"
	"bcts/common/xerror"
	"bcts/pkg/kafkaclient"
	"bcts/pkg/kafkaclient/mpupb"
	"context"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
)

type OtcQuoteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOtcQuoteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtcQuoteLogic {
	return &OtcQuoteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// OTC询价
func (l *OtcQuoteLogic) OtcQuote(in *pb.QuoteReq) (*pb.QuoteRsp, error) {
	l.Logger.Infof("order rpc OtcQuote, params:%+v", in)
	//TODO: 1.kafka断开超过一定时间不让询价/下单
	//TODO: 2.需要区分kafka里数据最后更新时间&行情数据创建时间（需要区分是行情系统有问题还是原始交易所的数据没有变化）

	symbol := in.Symbol
	if len(symbol) == 0 {
		return nil, xerror.ErrorSymbolIsEmpty
	}

	//用户信息校验
	userInfo, err := l.svcCtx.UserCenterRpc.GetUserInfo(l.ctx, &usercenter.UserInfoReq{UserID: in.UserID})
	if err != nil {
		return nil, err
	}

	//手续费
	feeInfo, err := l.svcCtx.DataServiceRpc.GetUserFeeInfo(l.ctx, &dataservice.GetUserFeeInfoReq{
		UserID:  in.UserID,
		GroupID: userInfo.GroupID,
		Symbol:  in.Symbol,
	})
	if err != nil {
		return nil, err
	}

	//feeInfo, err := QuerySymbolFee(in.UserID, symbol, l.svcCtx.NacosClient)
	//if err != nil {
	//	return nil, err
	//}

	l.Logger.Infof("OtcQuote, in:%+v, feeInfo:%+v", in, feeInfo)

	volumePrecision := feeInfo.SymbolInfo.VolumePrecision
	pricePrecision := feeInfo.SymbolInfo.PricePrecision
	amountPrecision := feeInfo.SymbolInfo.AmountPrecision
	//minUnit, _ = decimal.NewFromString(feeInfo.MinUnit)//最小交易单位

	//volume, amount, amountExceptFee, err := checkParams(in.QuoteType, in.Volume, in.Amount, feeInfo)
	volume, amount, amountExceptFee, err := ordercheck.CheckUniversal(in.QuoteType, in.Volume, in.Amount, feeInfo)
	if err != nil {
		return nil, err
	}

	takerFee, _ := decimal.NewFromString(feeInfo.TakerFee)
	priceRsp, volumeRsp, err := quote(symbol, in.Direction, in.QuoteType, volume, amountExceptFee, volumePrecision, pricePrecision, amountPrecision)
	if err != nil {
		return nil, err
	}

	var (
		turnOverOutput decimal.Decimal
		volumeOutput   decimal.Decimal
	)

	userIdStr := strconv.FormatInt(in.UserID, 10)
	redisKey := userIdStr + globalKey.OTCRedisPriceSuffix
	if err = l.svcCtx.RedisConn.Setex(redisKey, priceRsp.String(), 10); err != nil { //询价有效期 10s
		l.Logger.Errorf("quote set redis error, err:%+v", err)
		return nil, xerror.ErrorTryAgain
	}

	if in.QuoteType == globalKey.QuoteTypeVolume { //根据数量
		volumeOutput = volume
		var estimateFee decimal.Decimal
		if feeInfo.FeeKind == globalKey.FeeKindPercentage { //1表示百比分，2表示绝对值
			estimateFee = priceRsp.Mul(volume).Mul(takerFee)
		} else {
			estimateFee = takerFee
		}
		if in.Direction == globalKey.DirectionBuy {
			turnOverOutput = priceRsp.Mul(volume).Add(estimateFee).Truncate(int32(amountPrecision))
		} else {
			turnOverOutput = priceRsp.Mul(volume).Sub(estimateFee).Truncate(int32(amountPrecision))
		}
	} else { //根据总额
		turnOverOutput = amount
		volumeOutput = volumeRsp.Truncate(int32(volumePrecision))
	}

	rsp := &pb.QuoteRsp{
		QuoteID:   in.QuoteID,
		UserID:    in.UserID,
		Symbol:    symbol,
		Direction: in.Direction,
		QuoteType: in.QuoteType,
		Price:     priceRsp.String(),
		Volume:    volumeOutput.String(),
		Amount:    turnOverOutput.String(),
	}

	return rsp, nil
}

//amountExceptFee decimal.Decimal //金额询价，扣除预估手续费后的可用于实际交易的金额
func checkParams(quoteType int32, volumeStr, amountStr string, feeInfo *dataservice.GetUserFeeInfoRsp) (volume, amount, amountExceptFee decimal.Decimal, err error) {
	volumePrecision := feeInfo.SymbolInfo.VolumePrecision
	//pricePrecision := feeInfo.PricePrecision
	amountPrecision := feeInfo.SymbolInfo.AmountPrecision
	minUnit, _ := decimal.NewFromString(feeInfo.SymbolInfo.MinUnit)
	takerFee, _ := decimal.NewFromString(feeInfo.TakerFee)

	//数量询价
	if quoteType == globalKey.QuoteTypeVolume {
		volume, err = decimal.NewFromString(volumeStr)
		if err != nil {
			err = errors.Wrapf(xerror.ErrorOrderVolumeInvalid, "volume:%s, err:%+v", volumeStr, err)
			return
		}
		//1.
		if volume.LessThanOrEqual(decimal.Zero) {
			err = errors.Wrapf(xerror.ErrorOrderVolumeInvalid, "volume:%s, err:%+v", volumeStr, err)
			return
		}

		//2.
		volumeP := volume.Truncate(int32(volumePrecision))
		if !volume.Equal(volumeP) {
			//l.Logger.Errorf("volumePrecision:%d, volume:%s, feeInfo:%+v", volumePrecision, volume.String(), feeInfo)
			err = xerror.ErrorOrderVolumeInvalid
			return
		}

		//3. 最小交易单位
		if volume.Mod(minUnit).GreaterThan(decimal.Zero) {
			//l.Logger.Errorf("minUnit:%s, volume:%s, feeInfo:%+v", minUnit.String(), volume.String(), feeInfo)
			err = errors.Wrapf(xerror.ErrorOrderVolumeMinUnitInvalid, "min trade unit:%s", minUnit.String())
			return
		}

	} else if quoteType == globalKey.QuoteTypeAmount { //金额询价
		amount, err = decimal.NewFromString(amountStr)
		if err != nil {
			err = errors.Wrapf(xerror.ErrorOrderAmountInvalid, "amount:%s, err:%+v", amountStr, err)
			return
		}
		//1.
		if amount.LessThanOrEqual(decimal.Zero) {
			err = errors.Wrapf(xerror.ErrorOrderVolumeInvalid, "amount:%s, err:%+v", amountStr, err)
			return
		}

		amountP := amount.Truncate(int32(amountPrecision))
		if !amount.Equal(amountP) {
			//l.Logger.Errorf("amountPrecision:%d, amount:%s, feeInfo:%+v", amountPrecision, amount.String(), feeInfo)
			err = xerror.ErrorOrderAmountPrecisionInvalid
			return
		}
		//如果是总额询价,要先扣除手续费
		if feeInfo.FeeKind == globalKey.FeeKindPercentage { //1表示百比分，2表示绝对值
			amountExceptFee = amount.Mul(decimal.NewFromInt(1).Sub(takerFee))
		} else {
			amountExceptFee = amount.Sub(takerFee)
		}

	} else {
		//l.Logger.Errorf("OtcQuote quote type error, in:%+v", in)
		err = xerror.ErrorParamError
		return
	}
	return
}

func QuerySymbolFee(userId int64, symbol string, nacosClient *nacosAdapter.Client) (feeInfo *FeeInfo, err error) {
	feeInfos, err := QueryFee(userId, symbol, nacosClient)
	if err != nil {
		//这种error 前端并不能解析出来
		return nil, errors.Wrapf(xerror.ErrorTryAgain, "QueryFee err:%+v", err)
	}
	if len(feeInfos) > 0 {
		for _, v := range feeInfos {
			if v.Symbol == symbol {
				feeInfo = v
				break
			}
		}
		if feeInfo == nil {
			return nil, errors.Wrapf(xerror.ErrorSymbolNotFound, "QueryFee symbol not found symbol:%s", symbol)
		}
	}
	return
}
func QueryFee(userId int64, symbol string, nacosClient *nacosAdapter.Client) (feeInfos RspQueryFee, err error) {
	//symbolParam, err := l.svcCtx.NacosClient.GetSymbol(symbol)
	symbolParam, err := nacosClient.GetSymbol(symbol)
	if err != nil {
		//l.Logger.Error(err)
		return
	}
	hundred := decimal.NewFromInt(100)
	for i, symParam := range symbolParam {
		if symParam.FeeKind == 1 { //百分比 0.001
			takerFee, _ := decimal.NewFromString(symbolParam[i].TakerFee)
			makerFee, _ := decimal.NewFromString(symbolParam[i].MakerFee)
			symbolParam[i].TakerFee = takerFee.Mul(hundred).String()
			symbolParam[i].MakerFee = makerFee.Mul(hundred).String()
		}
	}
	for _, param := range symbolParam {
		feeInfo := &FeeInfo{
			SymbolKind:      param.SymbolKind,
			Underlying:      param.Underlying,
			PrimaryCurrency: param.PrimaryCurrency,
			BidCurrency:     param.BidCurrency,
			SettleCurrency:  param.SettleCurrency,

			BuyCommissionSwitch:  false,
			SellCommissionSwitch: false,
			Switch:               param.Switch,
			VolumePrecision:      param.VolumePrecision, //<= param.AmountPrecision
			PricePrecision:       param.PricePrecision,
			AmountPrecision:      param.AmountPrecision, //<= param.AmountPrecision
			MinUnit:              param.MinUnit,
			MinChangePrice:       param.MinChangePrice,
			Spread:               param.Spread,
			SingleMinOrderVolume: param.SingleMinOrderVolume,
			SingleMaxOrderVolume: param.SingleMaxOrderVolume,
			SingleMinOrderAmount: param.SingleMinOrderAmount,
			SingleMaxOrderAmount: param.SingleMaxOrderAmount,
			BuyPriceLimit:        param.BuyPriceLimit,
			SellPriceLimit:       param.SellPriceLimit,
			MaxMatchGear:         param.MaxMatchGear,
			Symbol:               param.SymbolID,
			FeeKind:              param.FeeKind,
			TakerFee:             param.TakerFee,
			MakerFee:             param.MakerFee,

			OtcMinOrderVolume: param.OtcMinOrderVolume, //param.OtcMinOrder
			OtcMaxOrderVolume: param.OtcMaxOrderVolume, //param.OtcMaxOrder,
			OtcMinOrderAmount: param.OtcMinOrderAmount, //param.OtcMinPrice,
			OtcMaxOrderAmount: param.OtcMaxOrderAmount, //param.OtcMaxPrice,
		}
		feeInfos = append(feeInfos, feeInfo)
	}

	//先不考虑用户组
	/*
		if userId != 0 {
			groupEntity, e := l.getMemberGroup(userId)
			if e != nil {
				if e != model.ErrNotFound {
					err = e
					l.Logger.Error(err)
					return
				}
			}
			if e == nil {
				l.Logger.Infof("group code:%s", groupEntity.GroupCode)
				tradeParams, e := l.getGroupTradeParamBySymbol(groupEntity.GroupCode, symbol)
				if e != nil {
					err = e
					l.Logger.Error(err)
					return
				}
				for _, param := range tradeParams {
					for _, feeInfo := range feeInfos {
						if feeInfo.Symbol == param.SymbolId {
							feeInfo.FeeKind = param.FeeKind
							feeInfo.TakerFee = param.TakerFee
							feeInfo.MakerFee = param.MakerFee

							feeInfo.DailyTradingLimitSwitch = true
							dailyBuyingLimitPerClientWeekday, _ := decimal.NewFromString(param.DailyBuyingLimitPerClientWeekday)
							dailyBuyingLimitPerClientWeekend, _ := decimal.NewFromString(param.DailyBuyingLimitPerClientWeekend)
							dailySellingLimitPerClientWeekday, _ := decimal.NewFromString(param.DailySellingLimitPerClientWeekday)
							dailySellingLimitPerClientWeekend, _ := decimal.NewFromString(param.DailySellingLimitPerClientWeekend)
							feeInfo.DailyBuyingLimitPerClientWeekday = dailyBuyingLimitPerClientWeekday
							feeInfo.DailyBuyingLimitPerClientWeekend = dailyBuyingLimitPerClientWeekend
							feeInfo.DailySellingLimitPerClientWeekday = dailySellingLimitPerClientWeekday
							feeInfo.DailySellingLimitPerClientWeekend = dailySellingLimitPerClientWeekend
							break
						}
					}
				}

				tradeAuth, e := l.getGroupTradeAuthBySymbol(groupEntity.GroupCode, symbol)
				if e != nil {
					err = e
					l.Logger.Error(err)
					return
				}
				for _, param := range tradeAuth {
					for _, feeInfo := range feeInfos {
						if feeInfo.Symbol == param.SymbolId {
							feeInfo.SellCommissionSwitch = param.SellCommissionSwitch
							feeInfo.BuyCommissionSwitch = param.BuyCommissionSwitch
							break
						}
					}
				}
			}
		}

	*/

	return
}

func (l *OtcQuoteLogic) getGroupTradeParamBySymbol(groupCode string, symbol string) (rspTradeParams []*nacosAdapter.GroupTradeParam, err error) {
	tradeParams, err := l.svcCtx.NacosClient.GetGroupTradeParam(groupCode)
	if err != nil {
		l.Logger.Error(err)
		return
	}
	hundred := decimal.NewFromInt(100)
	for i, param := range tradeParams {
		if param.FeeKind == 1 {
			takerFee, _ := decimal.NewFromString(tradeParams[i].TakerFee)
			makerFee, _ := decimal.NewFromString(tradeParams[i].MakerFee)
			tradeParams[i].TakerFee = takerFee.Mul(hundred).String()
			tradeParams[i].MakerFee = makerFee.Mul(hundred).String()
		}
	}
	if symbol == "" {
		rspTradeParams = tradeParams
		return
	}
	for _, param := range tradeParams {
		if param.SymbolId == symbol {
			rspTradeParams = append(rspTradeParams, param)
			break
		}
	}
	return
}

func (l *OtcQuoteLogic) getGroupTradeAuthBySymbol(groupCode string, symbol string) (rspTradeAuth []*nacosAdapter.GroupTradeAuth, err error) {
	tradeAuth, err := l.svcCtx.NacosClient.GetGroupTradeAuth(groupCode)
	if err != nil {
		l.Logger.Error(err)
		return
	}
	if symbol == "" {
		rspTradeAuth = tradeAuth
		return
	}
	for _, param := range tradeAuth {
		if param.SymbolId == symbol {
			rspTradeAuth = append(rspTradeAuth, param)
			break
		}
	}
	return
}

func (l *OtcQuoteLogic) getMemberGroup(userId int64) (*model.UserGroup, error) {

	//TODO: members表要重构
	groupId := int64(7)
	group, err := l.svcCtx.UserGroupModel.FindOne(l.ctx, groupId)
	return group, err
}

func quote(symbol string, direction int32, quoteType int32, volume, amount decimal.Decimal, volumePrecision, pricePrecision, amountPrecision int64) (priceRsp, volumeRsp decimal.Decimal, err error) {
	logx.Info(len(kafkaclient.DepthDataList))
	if len(kafkaclient.DepthDataList) > 0 {
		for k, v := range kafkaclient.DepthDataList {
			//l.Logger.Info(k) //FTX:ETH_USDT
			//	l.Logger.Info(v.Timestamp)
			//Exchange:FTX, Symbol:ETH_USDT
			logx.Infof("Key:%s, Exchange:%s, Symbol:%s, len(Asks):%d, len(Bids):%d", k, v.Exchange, v.Symbol, len(v.Asks), len(v.Bids))
		}
		if symbolDepthData, ok := kafkaclient.DepthDataList[symbol]; ok {
			var priceVolumeList []*mpupb.PriceVolume
			//fmt.Printf("asks:%+v\n", symbolDepthData.Asks)
			//fmt.Printf("bids:%+v\n", symbolDepthData.Bids)
			if direction == 0 { //用户买 BUY
				priceVolumeList = symbolDepthData.Asks
			} else { //用户卖 SELL
				priceVolumeList = symbolDepthData.Bids
			}
			var price decimal.Decimal
			var completed bool
			if quoteType == globalKey.QuoteTypeVolume {
				price, completed = matchPriceWithDepthByVolume(volume, priceVolumeList)
			} else {
				price, completed = matchPriceWithDepthByAmount(amount, priceVolumeList)
			}

			logx.Infof("final price:%s, completed:%+v", price.String(), completed)
			if completed {
				priceRsp = price.Truncate(int32(pricePrecision))
				if quoteType == globalKey.QuoteTypeVolume {
					volumeRsp = volume
				} else {
					volumeRsp = amount.Div(priceRsp)
				}
			} else {
				err = errors.New("行情档不够/no enough balance/......") //TODO:...
			}
		}
	} else {
		return decimal.Zero, decimal.Zero, errors.Wrapf(xerror.ErrorQuoteError, "kafka have no data")
	}

	return
}

func matchPriceWithDepthByVolume(volume decimal.Decimal, priceVolumeList []*mpupb.PriceVolume) (priceRsp decimal.Decimal, completed bool) {
	matchedVolume := decimal.Zero
	matchedAmount := decimal.Zero

	for _, v := range priceVolumeList {
		logx.Infof("matchPriceWithDepthByVolume market:%+v", v)
		priceInDepth, _ := decimal.NewFromString(v.Price)
		volumeInDepth, _ := decimal.NewFromString(v.Volume)
		if matchedVolume.Add(volumeInDepth).GreaterThanOrEqual(volume) {
			diffVolume := volume.Sub(matchedVolume)
			matchedVolume = matchedVolume.Add(diffVolume)

			matchedAmount = matchedAmount.Add(diffVolume.Mul(priceInDepth))
			priceRsp = matchedAmount.Div(volume) //价格为平均值
			completed = true
			logx.Infof("matchPriceWithDepthByVolume finally priceInDepth:%s, volumeInDepth:%s, volume:%s(matched:%s, matchedAmount:%s)",
				v.Price, v.Volume, volume, matchedVolume.String(), matchedAmount.String())
			break
		} else {
			matchedVolume = matchedVolume.Add(volumeInDepth)
			matchedAmount = matchedAmount.Add(volumeInDepth.Mul(priceInDepth))
			logx.Infof("matchPriceWithDepthByVolume priceInDepth:%s, volumeInDepth:%s, volume:%s(matched:%s, matchedAmount:%s)",
				v.Price, v.Volume, volume, matchedVolume.String(), matchedAmount.String())
		}
	}
	return
}

func matchPriceWithDepthByAmount(amount decimal.Decimal, priceVolumeList []*mpupb.PriceVolume) (priceRsp decimal.Decimal, completed bool) {
	matchedVolume := decimal.Zero
	matchedAmount := decimal.Zero
	for _, v := range priceVolumeList {
		logx.Infof("matchPriceWithDepthByAmount market:%+v", v)
		priceInDepth, _ := decimal.NewFromString(v.Price)
		volumeInDepth, _ := decimal.NewFromString(v.Volume)
		if matchedAmount.Add(priceInDepth.Mul(volumeInDepth)).GreaterThanOrEqual(amount) {
			diffAmount := amount.Sub(matchedAmount)
			matchedAmount = matchedAmount.Add(diffAmount)
			matchedVolume = matchedVolume.Add(diffAmount.Div(priceInDepth))

			priceRsp = matchedAmount.Div(matchedVolume) //价格为平均值
			completed = true
			logx.Infof("matchPriceWithDepthByAmount finally priceInDepth:%s, volumeInDepth:%s, amount:%s(matched:%s, matchedVolume:%s)",
				v.Price, v.Volume, amount, matchedAmount.String(), matchedVolume.String())
			break
		} else {
			matchedVolume = matchedVolume.Add(volumeInDepth)
			matchedAmount = matchedAmount.Add(volumeInDepth.Mul(priceInDepth))
			logx.Infof("matchPriceWithDepthByAmount priceInDepth:%s, volumeInDepth:%s, amount:%s(matched:%s, matchedVolume:%s)",
				v.Price, v.Volume, amount, matchedAmount.String(), matchedVolume.String())
		}
	}

	return
}
