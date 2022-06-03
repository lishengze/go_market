package offlinetrade

import (
	repository2 "bcts/app/admin/api/internal/repository"
	"bcts/app/admin/api/internal/svc"
	"bcts/app/admin/api/internal/types"
	"bcts/app/admin/model"
	"bcts/common/middleware"
	"bcts/common/nacosAdapter"
	"bcts/common/utils"
	"bcts/common/xerror"
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strconv"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOfflineTradeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateOfflineTradeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOfflineTradeLogic {
	return &UpdateOfflineTradeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateOfflineTradeLogic) BeforeUpdate(req *types.OfflineTradeReq) (err error) {
	req.Operator = l.ctx.Value(middleware.JWT_USER_NAME).(string)
	base := NewBaseLogic(l.ctx, l.svcCtx)
	validators := []validator{
		base.validateOperator,
		base.validateOfflineTradeID,
		base.validateUID,
		base.validateDecimal,
	}

	for _, validator := range validators {
		if err = validator(req); err != nil {
			return
		}
	}
	l.Logger.Info("BeforeUpdate")
	return
}

func (l *UpdateOfflineTradeLogic) UpdateOfflineTrade(req *types.OfflineTradeReq) (resp *types.OfflineTradeInput, err error) {
	if err = l.BeforeUpdate(req); err != nil {
		return
	}

	memeberId, err := strconv.ParseInt(req.UID, 10, 64)
	if err != nil {
		return
	}

	trade, err := l.svcCtx.OfflineTradeInputModel.FindOneByOfflineTradeId(l.ctx, req.OfflineTradeID)

	if err == model.ErrNotFound {
		//NewCodeError(1213, fmt.Sprintf("TradeBusinessNo:%s record not found", req.TradeBusinessNo), "")
		err = xerror.ErrorTradeBusinessNoNotFound
		return
	}

	if int(trade.Status.Int64) == model.OfflineTradeInputStatusFirst.INT() {
		//NewCodeError(1313, fmt.Sprintf("TradeBusinessNo:%s can not update", req.TradeBusinessNo), "")
		err = xerror.ErrorOfflineTradeNotUpdate
		return
	}

	if memeberId != trade.Uid.Int64 {
		//NewCodeError(1305, "memberID can not change", "")
		err = xerror.ErrorUidAlreadyUpdate
		return
	}

	//6. offlineTradeID 已经存在了
	if trade.OfflineTradeId != req.OfflineTradeID {
		if _, err = l.svcCtx.OfflineTradeInputModel.FindOneByOfflineTradeId(l.ctx, req.OfflineTradeID); err == nil {
			//NewCodeError(1315, "offline trade ID duplicate", "")
			err = xerror.ErrorOfflineTradeAlreadyExist
			return
		}
	}
	symbolRepository := repository2.NewSymbolRepository(l.ctx, l.svcCtx)
	symbolInfo, err := symbolRepository.GetSymbolInfoByName(req.Symbol)
	if err != nil {
		return
	}

	err = l.UpdateOfflineTradeInput(req, trade, memeberId, symbolInfo)
	if err != nil {
		return
	}

	return
}

func (l *UpdateOfflineTradeLogic) UpdateOfflineTradeInput(req *types.OfflineTradeReq, trade *model.OfflineTradeInput, memberId int64, symbol *nacosAdapter.Symbol) error {
	base := NewBaseLogic(l.ctx, l.svcCtx)

	quantity, _ := decimal.NewFromString(req.Quantity)
	amount, _ := decimal.NewFromString(req.Amount)
	fee, _ := decimal.NewFromString(req.Fee)

	if fee.GreaterThan(amount) {
		return xerror.ErrorFeeGtAmount
	}

	dateA := time.Now().In(time.UTC).Format("2006-01-02")
	timeX, _ := time.ParseInLocation("2006-01-02", dateA, time.UTC)

	//默认都是吧修改的
	tradeFix := l.TradeHasFixed(req, trade)
	remarksFix := l.RemarksHasFixed(req, trade)
	tradeTimeFix := l.TradeTimeHasFixed(req, trade)
	offlineTradeIDFix := l.OfflineTradeIDHasFixed(req, trade)
	//symbolFix := true

	if !tradeTimeFix && !offlineTradeIDFix && !tradeFix && !remarksFix {
		//NewCodeError(1316, "record is same, nothing changes", "")
		return xerror.ErrorNotingChange
	}

	//只需要改动offline_input
	//只有remarks变动
	if remarksFix && !tradeFix && !tradeTimeFix && !offlineTradeIDFix {
		e := l.svcCtx.Mysql.Transact(func(session sqlx.Session) error {

			//1. 写offlinetradeinput
			//now := time.Now().In(time.UTC)
			tradeHistory := generateOfflineTradeInputHistory(req, trade, remarksFix, tradeTimeFix, offlineTradeIDFix, tradeFix)
			_, e := l.svcCtx.OfflineTradeInputHistoryModel.Insert(l.ctx, session, tradeHistory)
			if e != nil {
				l.Logger.Errorf("remarksFix OfflineTradeInputHistoryModel Insert err: %s", e.Error())
				return e
			}
			//待修改
			trade.Remarks = utils.NewString(req.Remarks)
			trade.ReviseOperator = utils.NewString(req.Operator)
			trade.OperateTime = tradeHistory.CreateTime
			if _, e := l.svcCtx.OfflineTradeInputModel.Update(l.ctx, session, trade); e != nil {
				l.Logger.Errorf("remarksFix OfflineTradeInputHistoryModel Update err: %s", e.Error())
				return e
			}
			return nil
		})
		if e != nil {
			return xerror.ErrorOfflineTradeUpdateFail
		}
		return nil
	}

	//只需要改动offline + trade
	//只有tradeTime 或者 OfflineTradeID 变动
	if !tradeFix && (tradeTimeFix || offlineTradeIDFix || remarksFix) {
		e := l.svcCtx.Mysql.Transact(func(session sqlx.Session) error {

			tradeTime, _ := time.ParseInLocation("2006-01-02", req.TradeTime, time.UTC)
			//1. 写offlinetradeinput
			//now := time.Now().In(time.UTC)
			// 覆盖tradeTime
			//trade.SetTradeTime(tradeTime)
			tradeHistory := generateOfflineTradeInputHistory(req, trade, remarksFix, tradeTimeFix, offlineTradeIDFix, tradeFix)
			_, e := l.svcCtx.OfflineTradeInputHistoryModel.Insert(l.ctx, session, tradeHistory)
			if e != nil {
				l.Logger.Errorf("tradeTimeFix OfflineTradeInputHistoryModel Insert err: %s", e.Error())
				return e
			}
			//待修改
			trade.OfflineTradeId = req.OfflineTradeID
			trade.TradeTime = utils.NewTime(tradeTime)
			trade.ReviseOperator = utils.NewString(req.Operator)
			trade.OperateTime = tradeHistory.CreateTime
			trade.Remarks = utils.NewString(req.Remarks)
			if _, e := l.svcCtx.OfflineTradeInputModel.Update(l.ctx, session, trade); e != nil {
				l.Logger.Errorf("tradeTimeFix OfflineTradeInputHistoryModel Update err: %s", e.Error())
				return e
			}
			//2. 修改trade
			tradeID, _ := strconv.ParseInt(trade.TradeId.String, 10, 64)
			tradeEntity, err := l.svcCtx.ExchangeTradesModel.FindOne(l.ctx, tradeID)
			if err != nil && err == model.ErrNotFound {
				// NewCodeError(500, fmt.Sprintf("tradeID:%s is not found", trade.TradeId.String), "")
				return xerror.ErrorOfflineTradeNotFound
			}

			symbolRepository := repository2.NewSymbolRepository(l.ctx, l.svcCtx)
			alarmRepository := repository2.NewAlarmRepository(l.ctx, l.svcCtx)

			if tradeTimeFix { //todo 修改了时间，收益均价都需要变化
				indexOfCurrency := strings.Index(trade.Symbol.String, "_")
				//leftCurrency := trade.Symbol[:indexOfCurrency]
				rightCurrency := trade.Symbol.String[indexOfCurrency+1:]
				//老的市价
				rightBaseMarketPriceOld := tradeEntity.CoinPrice
				var averagePrice decimal.Decimal
				var profit decimal.Decimal
				var rightBaseMarketPriceNew decimal.Decimal
				var leftCoinMarketPriceNew decimal.Decimal
				//新的市价
				if utils.CompareDate(req.TradeTime) {
					addDay := tradeTime.AddDate(0, 0, 1).Format("2006-01-02") //取日终价格
					leftCoinMarketPriceNew, rightBaseMarketPriceNew, err = base.GetCoinPricesByTradeTime(addDay, req.Symbol)
					if err != nil {
						l.Logger.Errorf("tradeTimeFix GetCoinPricesByTradeTime err: %s", err.Error())
						return err
					}
				} else {
					rightBaseMarketPriceNew, err = symbolRepository.GetSymbolMarketPrice(rightCurrency)
					if err != nil {
						l.Logger.Errorf("tradeTimeFix GetSymbolMarketPrice err: %s", err.Error())
						_, _ = alarmRepository.DingPmsAlertMsg(fmt.Sprintf("UpdateOfflineTradeInput: memberId %d, Currency:%s  get marketPrice error", memberId, rightCurrency))
						return xerror.ErrorExchangeCoinPriceNotExist
					}
				}

				leftAccountEntity, rightAccountEntity, err := base.GetAccountsByMemeberIDAndCoinName(req, memberId)
				if err != nil {
					l.Logger.Errorf("tradeTimeFix GetAccountsByMemeberIDAndCoinName err: %s", err.Error())
					return err
				}

				if req.Direction == 0 { //买
					averagePrice = symbolRepository.BuyAveragePriceForUpdate(trade.Amount, trade.Quantity, leftAccountEntity, rightBaseMarketPriceOld)
					amountAdd := trade.Amount.Mul(rightBaseMarketPriceNew)
					balanceBefore := leftAccountEntity.Balance.Add(leftAccountEntity.FrozenBalance).Sub(trade.Quantity)
					amountBefore := balanceBefore.Mul(averagePrice)
					balanceTotal := balanceBefore.Add(quantity)
					if balanceTotal.Cmp(decimal.Zero) == 0 {
						// NewCodeError(500, fmt.Sprintf("offline trade %d change tradeTime err", tradeEntity.Id), "")
						return xerror.ErrorTotalBalanceIsEmpty
					}
					averagePrice = amountAdd.Add(amountBefore).DivRound(balanceTotal, 16)
					profit = rightBaseMarketPriceNew.Sub(rightAccountEntity.AveragePrice).Mul(trade.Amount)
					// 待修改
					//leftAccountEntity.UpdateAveragePrice(averagePrice)
					builder := model.NewAccountsUpdater(leftAccountEntity).UpdateAveragePrice(averagePrice).Builder()
					if _, e = l.svcCtx.AccountsModel.UpdateByBuilder(l.ctx, session, leftAccountEntity, builder); e != nil {
						l.Logger.Errorf("tradeTimeFix UpdateByBuilder err: %s", e.Error())
						return e
					}
				} else {
					amount1 := trade.Amount.Sub(trade.Fee)
					averagePrice = symbolRepository.SellAveragePriceForUpdate(amount1, rightAccountEntity, rightBaseMarketPriceOld)
					amountAdd := amount1.Mul(rightBaseMarketPriceNew)
					balanceBefore := rightAccountEntity.Balance.Add(rightAccountEntity.FrozenBalance).Sub(amount1)
					amountBefore := balanceBefore.Mul(averagePrice)
					balanceTotal := balanceBefore.Add(amount1)
					if balanceTotal.Cmp(decimal.Zero) == 0 {
						// NewCodeError(500, fmt.Sprintf("offline trade %d change tradeTime err", tradeEntity.Id), "")
						return xerror.ErrorTotalBalanceIsEmpty
					}
					averagePrice = amountAdd.Add(amountBefore).DivRound(balanceTotal, 16)
					coinPriceOld := quantity.Mul(leftAccountEntity.AveragePrice)
					profit = amountAdd.Sub(coinPriceOld)
					//待修改
					//rightAccountEntity.UpdateAveragePrice(averagePrice)
					builder := model.NewAccountsUpdater(rightAccountEntity).UpdateAveragePrice(averagePrice).Builder()
					if _, e = l.svcCtx.AccountsModel.UpdateByBuilder(l.ctx, session, rightAccountEntity, builder); e != nil {
						l.Logger.Errorf("tradeTimeFix UpdateByBuilder err: %s", e.Error())
						return e
					}
				}

				profitEntity, err := l.svcCtx.ProfitsModel.FindOneByTradeId(l.ctx, tradeID)
				if err != nil {
					l.Logger.Errorf("tradeTimeFix FindOneByTradeId err: %s", err.Error())
					return xerror.ErrorProfitNotFound
				}
				// 待修改
				profitEntity.Profit = profit
				profitEntity.ProfitTime = utils.NewTime(timeX)
				if _, e = l.svcCtx.ProfitsModel.Update(l.ctx, session, profitEntity); e != nil {
					l.Logger.Errorf("tradeTimeFix ProfitsModel Update err: %s", e.Error())
					return e
				}
				// 待修改
				tradeEntity.BasePrice = rightBaseMarketPriceNew
				tradeEntity.CoinPrice = leftCoinMarketPriceNew
			}
			// 待修改
			tradeEntity.OfflineTradeId = utils.NewString(req.OfflineTradeID)
			tradeEntity.Time = utils.NewInt64(tradeTime.Unix())
			if _, e = l.svcCtx.ExchangeTradesModel.Update(l.ctx, session, tradeEntity); e != nil {
				l.Logger.Errorf("tradeTimeFix ExchangeTradesModel Update err: %s", e.Error())
				return e
			}
			return nil
		})
		if e != nil {
			return e
		}
		return e
	}

	//其他情况所有表都要改动
	//offline/trade/流水/
	e := l.svcCtx.Mysql.Transact(func(session sqlx.Session) error {
		indexOfCurrency := strings.Index(trade.Symbol.String, "_")
		leftCurrency := trade.Symbol.String[:indexOfCurrency]
		rightCurrency := trade.Symbol.String[indexOfCurrency+1:]

		leftAccountEntity, rightAccountEntity, err := base.GetAccountsByMemeberIDAndCoinName(req, memberId)
		if err != nil {
			return err
		}

		tradeTime, _ := time.ParseInLocation("2006-01-02", req.TradeTime, time.UTC)
		price := decimal.Zero
		//var averagePrice decimal.Decimal
		if req.Direction == 0 { //买
			price = amount.Sub(fee).DivRound(quantity, int32(symbol.PricePrecision))
			//left计算均价
			//amounts := trade.Amount.Sub(trade.Fee)
			//averagePrice = CalcAveragePriceForUpdate(leftAccountEntity, amounts, trade.Quantity)
		} else {
			price = amount.Sub(fee).DivRound(quantity, int32(symbol.PricePrecision))
		}
		//1. offlinetradehistory
		tradeHistory := generateOfflineTradeInputHistory(req, trade, remarksFix, tradeTimeFix, offlineTradeIDFix, tradeFix)
		_, e := l.svcCtx.OfflineTradeInputHistoryModel.Insert(l.ctx, session, tradeHistory)
		if e != nil {
			l.Logger.Errorf("tradeFix OfflineTradeInputHistoryModel Insert err: %s", e.Error())
			return e
		}

		//2. 修改trade
		buyTurnover := decimal.Zero
		buyFee := decimal.Zero
		sellTurnover := decimal.Zero
		sellFee := decimal.Zero

		if req.Direction == 0 {
			buyTurnover = quantity
			buyFee = fee
		} else {
			sellTurnover = amount.Sub(fee)
			sellFee = fee
		}

		tradeID, _ := strconv.ParseInt(trade.TradeId.String, 10, 64)
		tradeEntity, err := l.svcCtx.ExchangeTradesModel.FindOne(l.ctx, tradeID)
		if err != nil && err == model.ErrNotFound {
			//NewCodeError(500, fmt.Sprintf("tradeID:%s is not found", trade.TradeId.String), "")
			return xerror.ErrorOfflineTradeNotFound
		}
		symbolRepository := repository2.NewSymbolRepository(l.ctx, l.svcCtx)
		alarmRepository := repository2.NewAlarmRepository(l.ctx, l.svcCtx)
		//todo 回滚之前的 获取之前输入的 市价;重新计算均价，作废之前的profit记录 获取历史币种的买入量
		amountHistory := trade.Amount.Sub(trade.Fee)
		//coinMarketPriceOld := trade.CoinPrice //之前录入时的市价
		baseMarketPriceOld := tradeEntity.BasePrice
		var averagePrice decimal.Decimal
		if trade.Direction.Int64 == int64(0) {
			averagePrice = symbolRepository.BuyAveragePriceForUpdate(trade.Amount, trade.Quantity, leftAccountEntity, baseMarketPriceOld)
		} else {
			averagePrice = symbolRepository.SellAveragePriceForUpdate(amountHistory, rightAccountEntity, baseMarketPriceOld)
		}
		indexOfCurrency = strings.Index(req.Symbol, "_")
		leftCurrencyNew := req.Symbol[:indexOfCurrency]
		rightCurrencyNew := req.Symbol[indexOfCurrency+1:]
		//获取新的币种市价
		var coinNewMarketPrice decimal.Decimal
		var baseNewMarketPrice decimal.Decimal
		addDay := tradeTime.AddDate(0, 0, 1).Format("2006-01-02") //取日终价格，所以日期加一
		if utils.CompareDate(req.TradeTime) {
			coinPriceEntity, err := l.svcCtx.CoinPricesModel.CoinPriceByTradeTime(l.ctx, addDay, leftCurrencyNew)
			if err == nil && coinPriceEntity != nil {
				coinNewMarketPrice = coinPriceEntity.Price.Decimal
			} else {
				//NewCodeError(500, fmt.Sprintf("Currency:%s is get marketPrice err", leftCurrencyNew), "")
				return xerror.ErrorExchangeCoinPriceNotExist
			}
			basePriceEntity, err := l.svcCtx.CoinPricesModel.CoinPriceByTradeTime(l.ctx, addDay, rightCurrencyNew)
			if err == nil && basePriceEntity != nil {
				baseNewMarketPrice = basePriceEntity.Price.Decimal
			} else {
				//NewCodeError(500, fmt.Sprintf("Currency:%s is get marketPrice err", rightCurrencyNew), "")
				return xerror.ErrorBaseCoinPriceNotExist
			}
		} else {
			coinNewMarketPrice, err = symbolRepository.GetSymbolMarketPrice(leftCurrencyNew)
			if err != nil {
				_, _ = alarmRepository.DingPmsAlertMsg(fmt.Sprintf("UpdateOfflineTradeInput: memberId: %d, Currency:%s  get marketPrice error", memberId, leftCurrencyNew))
				//NewCodeError(500, fmt.Sprintf("Currency:%s is get marketPrice err", leftCurrencyNew), "")
				return xerror.ErrorExchangeCoinPriceNotExist
			}
			baseNewMarketPrice, err = symbolRepository.GetSymbolMarketPrice(rightCurrencyNew)
			if err != nil {
				_, _ = alarmRepository.DingPmsAlertMsg(fmt.Sprintf("UpdateOfflineTradeInput: memberId: %d,  Currency:%s  get marketPrice error", memberId, rightCurrencyNew))
				//NewCodeError(500, fmt.Sprintf("Currency:%s is get marketPrice err", rightCurrencyNew), "")
				return xerror.ErrorBaseCoinPriceNotExist
			}
		}
		//待修改
		tradeEntity.OfflineTradeId = utils.NewString(req.OfflineTradeID)
		tradeEntity.Symbol = utils.NewString(req.Symbol)
		tradeEntity.CoinSymbol = utils.NewString(rightCurrencyNew)
		tradeEntity.BaseSymbol = utils.NewString(leftCurrencyNew)
		tradeEntity.Price = price
		tradeEntity.Amount = amount
		tradeEntity.BuyTurnover = buyTurnover
		tradeEntity.BuyFee = buyFee
		tradeEntity.SellTurnover = sellTurnover
		tradeEntity.SellFee = sellFee
		tradeEntity.Direction = utils.NewString(fmt.Sprintf("%d", req.Direction))
		tradeEntity.Time = utils.NewInt64(tradeTime.Unix())
		tradeEntity.CoinPrice = coinNewMarketPrice
		tradeEntity.BasePrice = baseNewMarketPrice
		if _, e = l.svcCtx.ExchangeTradesModel.Update(l.ctx, session, tradeEntity); e != nil {
			l.Logger.Errorf("tradeFix ExchangeTradesModel Update err: %s", e.Error())
			return e
		}

		//3 资金流水 恢复
		if trade.Direction.Int64 == 0 { //买
			leftTradeFeeTransaction := &model.MemberTransactions{
				Amount:          trade.Quantity.Neg(),
				Symbol:          utils.NewString(leftCurrency),
				MemberId:        utils.NewInt64(memberId),
				AccountId:       utils.NewInt64(leftAccountEntity.Id),
				DetailId:        utils.NewInt64(tradeID), //成交ID
				TransactionType: utils.NewInt64(int64(model.TransTypeTradeFix)),
				Fee:             trade.Quantity.Neg(),
				DiscountFee:     decimal.Zero,
				RealFee:         trade.Quantity.Neg(),
				PreBalance:      leftAccountEntity.Balance,
				Balance:         leftAccountEntity.Balance.Sub(trade.Quantity),
				PreFrozenBal:    leftAccountEntity.FrozenBalance,
				FrozenBal:       leftAccountEntity.FrozenBalance,
			}
			if _, err = l.svcCtx.MemberTransactionsModel.Insert(l.ctx, session, leftTradeFeeTransaction); err != nil {
				l.Logger.Errorf("tradeFix MemberTransactionsModel Insert err: %s", err.Error())
				return err
			}

			rightTradeFeeTransaction := &model.MemberTransactions{
				Amount:          trade.Amount.Sub(trade.Fee),
				Symbol:          utils.NewString(rightCurrency),
				MemberId:        utils.NewInt64(memberId),
				AccountId:       utils.NewInt64(rightAccountEntity.Id),
				DetailId:        utils.NewInt64(tradeID),
				TransactionType: utils.NewInt64(int64(model.TransTypeTradeFix)),
				Fee:             trade.Amount.Sub(trade.Fee),
				DiscountFee:     decimal.Zero,
				RealFee:         trade.Amount.Sub(trade.Fee),
				PreBalance:      rightAccountEntity.Balance,
				Balance:         rightAccountEntity.Balance.Add(trade.Amount.Sub(trade.Fee)),
				PreFrozenBal:    rightAccountEntity.FrozenBalance,
				FrozenBal:       rightAccountEntity.FrozenBalance,
			}
			if _, err = l.svcCtx.MemberTransactionsModel.Insert(l.ctx, session, rightTradeFeeTransaction); err != nil {
				l.Logger.Errorf("tradeFix MemberTransactionsModel Insert err: %s", err.Error())
				return err
			}

			feeTransaction := &model.MemberTransactions{
				Amount:          trade.Fee,
				Symbol:          utils.NewString(rightCurrency),
				MemberId:        utils.NewInt64(memberId),
				AccountId:       utils.NewInt64(rightAccountEntity.Id),
				DetailId:        utils.NewInt64(tradeID),
				TransactionType: utils.NewInt64(int64(model.TransTypeFeeFix)),
				Fee:             trade.Fee,
				DiscountFee:     decimal.Zero,
				RealFee:         trade.Fee,
				PreBalance:      rightAccountEntity.Balance,
				Balance:         rightAccountEntity.Balance.Add(trade.Fee),
				PreFrozenBal:    rightAccountEntity.FrozenBalance,
				FrozenBal:       rightAccountEntity.FrozenBalance,
			}
			if _, err = l.svcCtx.MemberTransactionsModel.Insert(l.ctx, session, feeTransaction); err != nil {
				l.Logger.Errorf("tradeFix MemberTransactionsModel Insert err: %s", err.Error())
				return err
			}
		} else { //卖
			leftTradeFeeTransaction := &model.MemberTransactions{
				Amount:          trade.Quantity,
				Symbol:          utils.NewString(leftCurrency),
				MemberId:        utils.NewInt64(memberId),
				AccountId:       utils.NewInt64(leftAccountEntity.Id),
				DetailId:        utils.NewInt64(tradeID),
				TransactionType: utils.NewInt64(int64(model.TransTypeTradeFix)),
				Fee:             trade.Quantity,
				DiscountFee:     decimal.Zero,
				RealFee:         trade.Quantity,
				PreBalance:      leftAccountEntity.Balance,
				Balance:         leftAccountEntity.Balance.Add(trade.Quantity),
				PreFrozenBal:    leftAccountEntity.FrozenBalance,
				FrozenBal:       leftAccountEntity.FrozenBalance,
			}
			if _, err := l.svcCtx.MemberTransactionsModel.Insert(l.ctx, session, leftTradeFeeTransaction); err != nil {
				l.Logger.Errorf("tradeFix MemberTransactionsModel Insert err: %s", err.Error())
				return err
			}

			rightTradeFeeTransaction := &model.MemberTransactions{
				Amount:          trade.Amount.Neg(),
				Symbol:          utils.NewString(rightCurrency),
				MemberId:        utils.NewInt64(memberId),
				AccountId:       utils.NewInt64(rightAccountEntity.Id),
				DetailId:        utils.NewInt64(tradeID),
				TransactionType: utils.NewInt64(int64(model.TransTypeTradeFix)),
				Fee:             decimal.Zero,
				DiscountFee:     decimal.Zero,
				RealFee:         decimal.Zero,
				PreBalance:      rightAccountEntity.Balance,
				Balance:         rightAccountEntity.Balance.Sub(trade.Amount.Sub(trade.Fee)),
				PreFrozenBal:    rightAccountEntity.FrozenBalance,
				FrozenBal:       rightAccountEntity.FrozenBalance,
			}
			if _, err := l.svcCtx.MemberTransactionsModel.Insert(l.ctx, session, rightTradeFeeTransaction); err != nil {
				l.Logger.Errorf("tradeFix MemberTransactionsModel Insert err: %s", err.Error())
				return err
			}

			feeTransaction := &model.MemberTransactions{
				Amount:          trade.Fee,
				Symbol:          utils.NewString(rightCurrency),
				MemberId:        utils.NewInt64(memberId),
				AccountId:       utils.NewInt64(rightAccountEntity.Id),
				DetailId:        utils.NewInt64(tradeID),
				TransactionType: utils.NewInt64(int64(model.TransTypeFeeFix)),
				Fee:             trade.Fee,
				DiscountFee:     decimal.Zero,
				RealFee:         trade.Fee,
				PreBalance:      rightAccountEntity.Balance,
				Balance:         rightAccountEntity.Balance.Add(trade.Fee),
				PreFrozenBal:    rightAccountEntity.FrozenBalance,
				FrozenBal:       rightAccountEntity.FrozenBalance,
			}
			if _, err := l.svcCtx.MemberTransactionsModel.Insert(l.ctx, session, feeTransaction); err != nil {
				l.Logger.Errorf("tradeFix MemberTransactionsModel Insert err: %s", err.Error())
				return err
			}
		}
		var leftAccountBuilder, rightAccountBuilder squirrel.UpdateBuilder
		//3 account表 恢复原始的资金
		if trade.Direction.Int64 == 0 { //买入
			//leftAccountEntity.SubBalance(trade.Quantity)
			//leftAccountEntity.UpdateAveragePrice(averagePrice)
			//rightAccountEntity.AddBalance(trade.Amount)
			leftAccountBuilder = model.NewAccountsUpdater(leftAccountEntity).
				SubBalance(trade.Quantity).
				UpdateAveragePrice(averagePrice).
				Builder()
			rightAccountBuilder = model.NewAccountsUpdater(rightAccountEntity).
				AddBalance(amount).
				Builder()
		} else {
			//leftAccountEntity.AddBalance(trade.Quantity)
			//rightAccountEntity.UpdateAveragePrice(averagePrice)
			//rightAccountEntity.SubBalance(trade.Amount.Sub(trade.Fee))
			leftAccountBuilder = model.NewAccountsUpdater(leftAccountEntity).
				AddBalance(trade.Quantity).
				Builder()
			rightAccountBuilder = model.NewAccountsUpdater(rightAccountEntity).
				UpdateAveragePrice(averagePrice).
				SubBalance(trade.Amount.Sub(trade.Fee)).
				Builder()
		}
		if _, e = l.svcCtx.AccountsModel.UpdateByBuilder(l.ctx, session, leftAccountEntity, leftAccountBuilder); e != nil {
			l.Logger.Errorf("tradeFix UpdateByBuilder err: %s", e.Error())
			return e
		}

		if _, e = l.svcCtx.AccountsModel.UpdateByBuilder(l.ctx, session, rightAccountEntity, rightAccountBuilder); e != nil {
			l.Logger.Errorf("tradeFix UpdateByBuilder err: %s", e.Error())
			return e
		}

		//通过tradeID查找profit 将之前线下录入的订单收益制为无效
		profitEntity, err := l.svcCtx.ProfitsModel.FindOneByTradeId(l.ctx, tradeID)
		if err != nil {
			//NewCodeError(500, "FindProfitByTradeID err", "")
			return xerror.ErrorProfitNotFound
		}
		// 待修改
		profitEntity.Status = 2
		if _, e = l.svcCtx.ProfitsModel.Update(l.ctx, session, profitEntity); e != nil {
			return e
		}

		// symbol可能会换掉 这里要取出新的account
		leftAccountEntityNew, err := l.svcCtx.AccountsModel.FindOneByMemberIdCoinUnit(l.ctx, utils.NewInt64(memberId), utils.NewString(leftCurrencyNew))
		if err != nil {
			//NewCodeError(500, fmt.Sprintf("Currency:%s is not found", leftCurrencyNew), "")
			return xerror.ErrorExchangeCoinNotExist
		}

		rightAccountEntityNew, err := l.svcCtx.AccountsModel.FindOneByMemberIdCoinUnit(l.ctx, utils.NewInt64(memberId), utils.NewString(rightCurrencyNew))
		if err != nil {
			//NewCodeError(500, fmt.Sprintf("Currency:%s is not found", rightCurrencyNew), "")
			return xerror.ErrorBaseCoinNotExist
		}

		//均价计算 收益计算 todo
		var averagePriceNew decimal.Decimal
		var profitNew decimal.Decimal
		if req.Direction == 0 { //买入，计算新的均价 还有 币种收益
			averagePriceNew, profitNew, err = symbolRepository.BuyAveragePrice(amount, quantity, leftAccountEntityNew, rightAccountEntityNew, baseNewMarketPrice)
			if err != nil {
				//return err
			}
			profitEntity := &model.Profits{
				MemberId:        utils.NewInt64(memberId),
				DetailId:        utils.NewInt64(tradeID),
				ProfitTime:      utils.NewTime(timeX),
				TransactionType: utils.NewInt64(int64(model.TransTypeExchange)),
				Symbol:          utils.NewString(rightCurrencyNew),
				Profit:          profitNew,
				Direction:       utils.NewInt64(int64(req.Direction)),
				Status:          1,
			}
			if _, err = l.svcCtx.ProfitsModel.Insert(l.ctx, session, profitEntity); err != nil {
				return err
			}
		} else {
			averagePriceNew, profitNew, err = symbolRepository.SellAveragePrice(amount.Sub(fee), quantity, leftAccountEntityNew.AveragePrice, rightAccountEntityNew, baseNewMarketPrice)
			if err != nil {
				//return err
			}
			profitEntity := &model.Profits{
				MemberId:        utils.NewInt64(memberId),
				DetailId:        utils.NewInt64(tradeID),
				ProfitTime:      utils.NewTime(timeX),
				TransactionType: utils.NewInt64(int64(model.TransTypeExchange)),
				Symbol:          utils.NewString(leftCurrencyNew),
				Profit:          profitNew,
				Direction:       utils.NewInt64(int64(req.Direction)),
				Status:          1,
			}
			if _, err = l.svcCtx.ProfitsModel.Insert(l.ctx, session, profitEntity); err != nil {
				return err
			}
		}
		fmt.Println("averagePriceNew", averagePriceNew)
		//4 资金流水新增
		if req.Direction == 0 { //买
			leftTradeFeeTransaction := &model.MemberTransactions{
				Amount:          quantity,
				Symbol:          utils.NewString(leftCurrencyNew),
				MemberId:        utils.NewInt64(memberId),
				AccountId:       utils.NewInt64(leftAccountEntityNew.Id),
				DetailId:        utils.NewInt64(tradeID),
				TransactionType: utils.NewInt64(int64(model.TransTypeTradeFix)),
				Fee:             quantity,
				DiscountFee:     decimal.Zero,
				RealFee:         quantity,
				PreBalance:      leftAccountEntityNew.Balance,
				Balance:         leftAccountEntityNew.Balance.Add(quantity),
				PreFrozenBal:    leftAccountEntityNew.FrozenBalance,
				FrozenBal:       leftAccountEntityNew.FrozenBalance,
				AveragePrice:    leftAccountEntityNew.AveragePrice,
				CoinPrice:       coinNewMarketPrice,
			}
			if _, err := l.svcCtx.MemberTransactionsModel.Insert(l.ctx, session, leftTradeFeeTransaction); err != nil {
				return err
			}

			rightTradeFeeTransaction := &model.MemberTransactions{
				Amount:          amount.Sub(fee).Neg(),
				Symbol:          utils.NewString(rightCurrencyNew),
				MemberId:        utils.NewInt64(memberId),
				AccountId:       utils.NewInt64(rightAccountEntityNew.Id),
				DetailId:        utils.NewInt64(tradeID),
				TransactionType: utils.NewInt64(int64(model.TransTypeTradeFix)),
				Fee:             amount.Sub(fee).Neg(),
				DiscountFee:     decimal.Zero,
				RealFee:         amount.Sub(fee).Neg(),
				PreBalance:      rightAccountEntityNew.Balance,
				Balance:         rightAccountEntityNew.Balance.Sub(amount.Sub(fee)),
				PreFrozenBal:    rightAccountEntityNew.FrozenBalance,
				FrozenBal:       rightAccountEntityNew.FrozenBalance,
				AveragePrice:    rightAccountEntityNew.AveragePrice,
				CoinPrice:       baseNewMarketPrice,
			}
			if _, err = l.svcCtx.MemberTransactionsModel.Insert(l.ctx, session, rightTradeFeeTransaction); err != nil {
				return err
			}

			feeTransaction := &model.MemberTransactions{
				Amount:          fee.Neg(),
				Symbol:          utils.NewString(rightCurrencyNew),
				MemberId:        utils.NewInt64(memberId),
				AccountId:       utils.NewInt64(rightAccountEntityNew.Id),
				DetailId:        utils.NewInt64(tradeID),
				TransactionType: utils.NewInt64(int64(model.TransTypeFeeFix)),
				Fee:             fee.Neg(),
				DiscountFee:     decimal.Zero,
				RealFee:         fee.Neg(),
				PreBalance:      rightAccountEntityNew.Balance,
				Balance:         rightAccountEntityNew.Balance.Sub(fee),
				PreFrozenBal:    rightAccountEntityNew.FrozenBalance,
				FrozenBal:       rightAccountEntityNew.FrozenBalance,
				AveragePrice:    rightAccountEntityNew.AveragePrice,
				CoinPrice:       baseNewMarketPrice,
			}
			if _, err := l.svcCtx.MemberTransactionsModel.Insert(l.ctx, session, feeTransaction); err != nil {
				return err
			}
		} else { //卖
			leftTradeFeeTransaction := &model.MemberTransactions{
				Amount:          quantity.Neg(),
				Symbol:          utils.NewString(leftCurrencyNew),
				MemberId:        utils.NewInt64(memberId),
				AccountId:       utils.NewInt64(leftAccountEntityNew.Id),
				DetailId:        utils.NewInt64(tradeID),
				TransactionType: utils.NewInt64(int64(model.TransTypeTradeFix)),
				Fee:             quantity.Neg(),
				DiscountFee:     decimal.Zero,
				RealFee:         quantity.Neg(),
				PreBalance:      leftAccountEntityNew.Balance,
				Balance:         leftAccountEntityNew.Balance.Sub(quantity),
				PreFrozenBal:    leftAccountEntityNew.FrozenBalance,
				FrozenBal:       leftAccountEntityNew.FrozenBalance,
				AveragePrice:    leftAccountEntityNew.AveragePrice,
				CoinPrice:       coinNewMarketPrice,
			}
			if _, err = l.svcCtx.MemberTransactionsModel.Insert(l.ctx, session, leftTradeFeeTransaction); err != nil {
				return err
			}

			rightTradeFeeTransaction := &model.MemberTransactions{
				Amount:          amount,
				Symbol:          utils.NewString(rightCurrencyNew),
				MemberId:        utils.NewInt64(memberId),
				AccountId:       utils.NewInt64(rightAccountEntityNew.Id),
				DetailId:        utils.NewInt64(tradeID),
				TransactionType: utils.NewInt64(int64(model.TransTypeTradeFix)),
				Fee:             decimal.Zero,
				DiscountFee:     decimal.Zero,
				RealFee:         decimal.Zero,
				PreBalance:      rightAccountEntityNew.Balance,
				Balance:         rightAccountEntityNew.Balance.Add(amount.Sub(fee)),
				PreFrozenBal:    rightAccountEntityNew.FrozenBalance,
				FrozenBal:       rightAccountEntityNew.FrozenBalance,
				AveragePrice:    rightAccountEntityNew.AveragePrice,
				CoinPrice:       baseNewMarketPrice,
			}
			if _, err = l.svcCtx.MemberTransactionsModel.Insert(l.ctx, session, rightTradeFeeTransaction); err != nil {
				return err
			}

			feeTransaction := &model.MemberTransactions{
				Amount:          fee.Neg(),
				Symbol:          utils.NewString(rightCurrencyNew),
				MemberId:        utils.NewInt64(memberId),
				AccountId:       utils.NewInt64(rightAccountEntity.Id),
				DetailId:        utils.NewInt64(tradeID),
				TransactionType: utils.NewInt64(int64(model.TransTypeFeeFix)),
				Fee:             fee.Neg(),
				DiscountFee:     decimal.Zero,
				RealFee:         fee.Neg(),
				PreBalance:      rightAccountEntityNew.Balance,
				Balance:         rightAccountEntityNew.Balance.Sub(fee),
				PreFrozenBal:    rightAccountEntityNew.FrozenBalance,
				FrozenBal:       rightAccountEntityNew.FrozenBalance,
				AveragePrice:    rightAccountEntityNew.AveragePrice,
				CoinPrice:       baseNewMarketPrice,
			}
			if _, err = l.svcCtx.MemberTransactionsModel.Insert(l.ctx, session, feeTransaction); err != nil {
				return err
			}
		}

		//账户操作
		//3 account表 恢复原始的资金
		if req.Direction == 0 { //买入
			//leftAccountEntityNew.AddBalance(quantity)
			//leftAccountEntityNew.UpdateAveragePrice(averagePriceNew)
			//rightAccountEntityNew.SubBalance(amount)
			leftAccountBuilder = model.NewAccountsUpdater(leftAccountEntity).
				AddBalance(quantity).
				UpdateAveragePrice(averagePrice).
				Builder()
			rightAccountBuilder = model.NewAccountsUpdater(rightAccountEntity).
				SubBalance(amount).
				Builder()
		} else {
			//leftAccountEntityNew.SubBalance(quantity)
			//rightAccountEntityNew.AddBalance(amount.Sub(fee))
			//rightAccountEntityNew.UpdateAveragePrice(averagePriceNew)
			leftAccountBuilder = model.NewAccountsUpdater(leftAccountEntity).
				SubBalance(quantity).
				Builder()
			rightAccountBuilder = model.NewAccountsUpdater(rightAccountEntity).
				AddBalance(amount.Sub(fee)).
				UpdateAveragePrice(averagePrice).
				Builder()
		}

		if _, e = l.svcCtx.AccountsModel.UpdateByBuilder(l.ctx, session, leftAccountEntityNew, leftAccountBuilder); e != nil {
			l.Logger.Errorf("tradeFix UpdateByBuilder err: %s", e.Error())
			return e
		}

		if _, e = l.svcCtx.AccountsModel.UpdateByBuilder(l.ctx, session, rightAccountEntityNew, rightAccountBuilder); e != nil {
			l.Logger.Errorf("tradeFix UpdateByBuilder err: %s", e.Error())
			return e
		}
		//offlinetrade表
		trade.OfflineTradeId = req.OfflineTradeID
		trade.TradeTime = utils.NewTime(tradeTime)
		trade.Direction = utils.NewInt64(int64(req.Direction))
		trade.Symbol = utils.NewString(req.Symbol)
		trade.Quantity = quantity
		trade.Amount = amount
		trade.Price = price
		trade.Fee = fee
		trade.Remarks = utils.NewString(req.Remarks)
		trade.ReviseOperator = utils.NewString(req.Operator)
		trade.OperateTime = tradeHistory.CreateTime

		if _, e = l.svcCtx.OfflineTradeInputModel.Update(l.ctx, session, trade); e != nil {
			return e
		}

		return nil
	})
	if e != nil {
		l.Logger.Errorf("UpdateOfflineTradeInput: %d, tradeFix err:", trade.Id, e.Error())
		return xerror.ErrorOfflineTradeUpdateFail
	}
	return nil
}

func (l *UpdateOfflineTradeLogic) RemarksHasFixed(req *types.OfflineTradeReq, trade *model.OfflineTradeInput) bool {
	return req.Remarks != trade.Remarks.String
}

func (l *UpdateOfflineTradeLogic) OfflineTradeIDHasFixed(req *types.OfflineTradeReq, trade *model.OfflineTradeInput) bool {
	return req.OfflineTradeID != trade.OfflineTradeId
}

func (l *UpdateOfflineTradeLogic) TradeTimeHasFixed(req *types.OfflineTradeReq, trade *model.OfflineTradeInput) bool {
	// req.TradeTime,前端传YYYY-MM-DD,后台是YYYY-MM-DD HH:mm:ss比较,肯定不相等
	// 老版本默认是true,不会有这个问题
	// 重构之后的默认是false,此处按YYYY-MM-DD比较
	// .Format("2006-01-02 15:04:05")
	return req.TradeTime != trade.TradeTime.Time.Format("2006-01-02")
}

func (l *UpdateOfflineTradeLogic) TradeHasFixed(req *types.OfflineTradeReq, trade *model.OfflineTradeInput) bool {
	quantity, _ := decimal.NewFromString(req.Quantity)
	amount, _ := decimal.NewFromString(req.Amount)
	fee, _ := decimal.NewFromString(req.Fee)

	if int64(req.Direction) != trade.Direction.Int64 {
		return true
	}

	if req.Symbol != trade.Symbol.String {
		return true
	}

	if !quantity.Equal(trade.Quantity) {
		return true
	}

	if !amount.Equal(trade.Amount) {
		return true
	}

	if !fee.Equal(trade.Fee) {
		return true
	}

	return false
}

func generateOfflineTradeInputHistory(req *types.OfflineTradeReq, trade *model.OfflineTradeInput, remarksFix, tradeTimeFix, offlineTradeIDFix, tradeFix bool) *model.OfflineTradeInputHistory {
	now := time.Now().In(time.UTC)

	tradeHistory := model.OfflineTradeInputHistory{
		TradeBusinessNo: trade.TradeBusinessNo,
		TradeTime:       trade.TradeTime,
		OfflineTradeId:  trade.OfflineTradeId,
		Uid:             trade.Uid,
		UserName:        trade.UserName,
		Name:            trade.Name,
		Symbol:          trade.Symbol,
		Direction:       trade.Direction,
		Quantity:        trade.Quantity,
		Amount:          trade.Amount,
		Price:           trade.Price,
		Fee:             trade.Fee,
		TradeId:         trade.TradeId,
		Remarks:         utils.NewString(req.Remarks),
		Status:          trade.Status,
		Operator1:       trade.Operator1,
		Operator2:       trade.Operator2,
		ReviseOperator:  utils.NewString(req.Operator),
		OperateTime:     utils.NewTime(now),
		CreateTime:      utils.NewTime(now),
	}

	if remarksFix {
		tradeHistory.Remarks = utils.NewString(req.Remarks)
	}

	if tradeTimeFix {
		tradeTime, _ := time.ParseInLocation("2006-01-02", req.TradeTime, time.UTC)
		tradeHistory.TradeTime = utils.NewTime(tradeTime)
	}

	if offlineTradeIDFix {
		tradeHistory.OfflineTradeId = req.OfflineTradeID
	}

	if tradeFix {
		quantity, _ := decimal.NewFromString(req.Quantity)
		amount, _ := decimal.NewFromString(req.Amount)
		fee, _ := decimal.NewFromString(req.Fee)

		tradeHistory.Symbol = utils.NewString(req.Symbol)
		tradeHistory.Direction = utils.NewInt64(int64(req.Direction))
		tradeHistory.Quantity = quantity
		tradeHistory.Amount = amount
		tradeHistory.Fee = fee
	}

	return &tradeHistory
}
