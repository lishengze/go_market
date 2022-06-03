package offlinetrade

import (
	"context"
	"fmt"
	"market_server/app/admin/api/internal/repository"
	"market_server/app/admin/api/internal/svc"
	"market_server/app/admin/api/internal/types"
	"market_server/app/admin/model"
	"market_server/common/middleware"
	"market_server/common/nacosAdapter"
	"market_server/common/utils"
	"market_server/common/xerror"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOfflineTradeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateOfflineTradeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOfflineTradeLogic {
	return &CreateOfflineTradeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateOfflineTradeLogic) BeforeCreate(req *types.OfflineTradeReq) (err error) {
	req.Operator = l.ctx.Value(middleware.JWT_USER_NAME).(string)

	base := NewBaseLogic(l.ctx, l.svcCtx)
	validators := []validator{
		base.validateOperator,
		base.validateOfflineTradeID,
		base.validateDirection,
		base.validateUID,
		base.validateDecimal,
	}

	for _, validator := range validators {
		if err = validator(req); err != nil {
			return
		}
	}
	l.Logger.Info("BeforeCreate")
	return
}

func (l *CreateOfflineTradeLogic) CreateOfflineTrade(req *types.OfflineTradeReq) (resp *types.OfflineTradeInput, err error) {
	if err = l.BeforeCreate(req); err != nil {
		return
	}
	memberId, err := strconv.ParseInt(req.UID, 10, 64)
	if err != nil {
		return
	}
	member, err := l.svcCtx.MembersModel.FindOneByMemberId(l.ctx, utils.NewString(req.UID))
	if err != nil {
		return
	}

	symbolRepository := repository.NewSymbolRepository(l.ctx, l.svcCtx)
	symbolInfo, err := symbolRepository.GetSymbolInfoByName(req.Symbol)
	if err != nil {
		return
	}

	trade, err := l.svcCtx.OfflineTradeInputModel.FindOneByOfflineTradeId(l.ctx, req.OfflineTradeID)
	// 首次录入
	if err == model.ErrNotFound {
		err = l.NewFirstOfflineTradeInput(req, member, symbolInfo)
		if err != nil {
			//NewCodeError(1500, fmt.Sprintf("OfflineTradeID:%s create failed, %s", req.OfflineTradeID, err.Error()), "")
			err = xerror.ErrorOfflineTradeCreateFail
			return
		}

		return
	}

	//二次录入
	base := NewBaseLogic(l.ctx, l.svcCtx)
	if err = base.validateSecondTrade(req, trade); err != nil {
		return
	}

	err = l.NewSecondOfflineTradeInput(req, memberId, trade)
	if err != nil {
		return
	}

	return
}

func (l *CreateOfflineTradeLogic) NewFirstOfflineTradeInput(req *types.OfflineTradeReq, member *model.Members, symbolInfo *nacosAdapter.Symbol) (err error) {
	memberId := member.Id
	indexOfCurrency := strings.Index(req.Symbol, "_")
	leftCurrency := req.Symbol[:indexOfCurrency]
	rightCurrency := req.Symbol[indexOfCurrency+1:]

	//leftAccountEntity, e := logic.AccountRepository.AccountByMemeberIDAndCoinName(memberId, leftCurrency)
	_, e := l.svcCtx.AccountsModel.FindOneByMemberIdCoinUnit(l.ctx, utils.NewInt64(memberId), utils.NewString(leftCurrency))
	if e != nil {
		//xerror.NewCodeError(1301, fmt.Sprintf("Currency:%s is not found", leftCurrency), "")
		err = xerror.ErrorExchangeCoinNotExist
		return
	}

	//rightAccountEntity, e := logic.AccountRepository.AccountByMemeberIDAndCoinName(memberId, rightCurrency)
	_, e = l.svcCtx.AccountsModel.FindOneByMemberIdCoinUnit(l.ctx, utils.NewInt64(memberId), utils.NewString(rightCurrency))
	if e != nil {
		//NewCodeError(1301, fmt.Sprintf("Currency:%s is not found", rightCurrency), "")
		err = xerror.ErrorBaseCoinNotExist
		return
	}

	quantity, _ := decimal.NewFromString(req.Quantity)
	amount, _ := decimal.NewFromString(req.Amount)
	fee, _ := decimal.NewFromString(req.Fee)

	//买
	price := amount.Sub(fee).DivRound(quantity, int32(symbolInfo.PricePrecision))
	now := time.Now().In(time.UTC) //.Format("2006-01-02 15:04:05")
	tradeTime, _ := time.ParseInLocation("2006-01-02", req.TradeTime, time.UTC)
	baseLogic := NewBaseLogic(l.ctx, l.svcCtx)
	tradeInput := &model.OfflineTradeInput{
		TradeBusinessNo: baseLogic.GenerateTradeBusinessNo(),
		TradeTime:       utils.NewTime(tradeTime),
		OfflineTradeId:  req.OfflineTradeID,
		Uid:             utils.NewInt64(member.Id),
		UserName:        member.Username,
		Name:            member.RealName,
		Symbol:          utils.NewString(req.Symbol),
		Direction:       utils.NewInt64(int64(req.Direction)),
		Quantity:        quantity,
		Amount:          amount,
		Price:           price,
		Fee:             fee,
		TradeId:         utils.NewString(""),
		Remarks:         utils.NewString(req.Remarks),
		Status:          utils.NewInt64(int64(model.OfflineTradeInputStatusFirst.INT())),
		Operator1:       utils.NewString(req.Operator),
		OperateTime:     utils.NewTime(now),
	}

	_, err = l.svcCtx.OfflineTradeInputModel.Insert(l.ctx, nil, tradeInput)
	if err != nil {
		//NewCodeError(1500, fmt.Sprintf("OfflineTradeID:%s create failed, %s", req.OfflineTradeID, err.Error()), "")
		err = xerror.ErrorOfflineTradeCreateFail
		return
	}
	return
}

func (l *CreateOfflineTradeLogic) NewSecondOfflineTradeInput(req *types.OfflineTradeReq, memberId int64, trade *model.OfflineTradeInput) error {
	baseLogic := NewBaseLogic(l.ctx, l.svcCtx)
	leftAccountEntity, rightAccountEntity, err := baseLogic.GetAccountsByMemeberIDAndCoinName(req, memberId)
	if err != nil {
		return err
	}
	leftCurrency, rightCurrency := leftAccountEntity.CoinUnit, rightAccountEntity.CoinUnit

	quantity, _ := decimal.NewFromString(req.Quantity)
	amount, _ := decimal.NewFromString(req.Amount)
	fee, _ := decimal.NewFromString(req.Fee)

	if trade.Uid.Int64 != memberId {
		//NewCodeError(1305, fmt.Sprintf("UID:%s is not same", req.UID), "")
		return xerror.ErrorUidAlreadyUpdate
	}

	if trade.Symbol.String != req.Symbol {
		//NewCodeError(1306, fmt.Sprintf("Symbol:%s is not same", req.Symbol), "")
		return xerror.ErrorSymbolAlreadyUpdate
	}

	if trade.Direction.Int64 != int64(req.Direction) {
		//NewCodeError(1307, fmt.Sprintf("Direction:%d is not same", req.Direction), "")
		return xerror.ErrorDirectionInvalid
	}

	if !trade.Quantity.Equal(quantity) {
		//NewCodeError(1308, fmt.Sprintf("Quantity:%s is not same", req.Quantity), "")
		return xerror.ErrorQuantityAlreadyUpdate
	}

	if !trade.Amount.Equal(amount) {
		//return xerror.NewCodeError(1309, fmt.Sprintf("Amount:%s is not same", req.Amount), "")
		return xerror.ErrorAmountAlreadyUpdate
	}

	if !trade.Fee.Equal(fee) {
		//NewCodeError(1310, fmt.Sprintf("Fee:%s is not same", req.Fee), "")
		return xerror.ErrorFeeAlreadyUpdate
	}
	var coinMarketPrice decimal.Decimal
	var baseMarketPrice decimal.Decimal
	addDay := time.Time(trade.TradeTime.Time).AddDate(0, 0, 1).Format("2006-01-02") //取日终价格
	// 交易时间 < 当日时间,历史的币的市场价格存在在CoinPrice
	if utils.CompareDate(req.TradeTime) {
		coinMarketPrice, baseMarketPrice, err = baseLogic.GetCoinPricesByTradeTime(addDay, req.Symbol)
		if err != nil {
			return err
		}
	} else { // 交易未完成
		coinMarketPrice, baseMarketPrice, err = baseLogic.GetSymbolMarketPrices(memberId, req.Symbol)
		if err != nil {
			return err
		}
	}
	//校验通过
	e := l.svcCtx.Mysql.Transact(func(session sqlx.Session) error {
		id := baseLogic.GenerateDistributedId()
		idStr := strconv.FormatInt(id, 10)

		now := time.Now().In(time.UTC)

		tradeInputHistory := &model.OfflineTradeInputHistory{
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
			Status:          utils.NewInt64(int64(model.OfflineTradeInputStatusSecond.INT())),
			Operator1:       trade.Operator1,
			Operator2:       utils.NewString(req.Operator),
			ReviseOperator:  trade.ReviseOperator,
			OperateTime:     utils.NewTime(now),
			CreateTime:      utils.NewTime(now),
		}

		if _, e := l.svcCtx.OfflineTradeInputHistoryModel.Insert(l.ctx, session, tradeInputHistory); e != nil {
			return e
		}

		//1. 写offlinetradeinput
		tradeTime, _ := time.ParseInLocation("2006-01-02", req.TradeTime, time.UTC)

		dateA := time.Now().In(time.UTC).Format("2006-01-02")
		timeX, _ := time.ParseInLocation("2006-01-02", dateA, time.UTC)
		/* 待修改
		trade.SetOperator2(req.Operator)
		trade.SetStatus(entity.OfflineTradeInputStatusSecond.INT())
		trade.SetOperateTime(now)
		trade.SetTradeID(idStr)
		trade.SetRemarks(req.Remarks)
		*/
		trade.Operator2 = utils.NewString(req.Operator)
		trade.Status = utils.NewInt64(int64(model.OfflineTradeInputStatusSecond.INT()))
		trade.OperateTime = utils.NewTime(now)
		trade.TradeId = utils.NewString(idStr)
		trade.Remarks = utils.NewString(req.Remarks)

		if _, e := l.svcCtx.OfflineTradeInputModel.Update(l.ctx, session, trade); e != nil {
			return e
		}

		//2. trade表
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

		uid, err := strconv.ParseInt(req.UID, 10, 64)
		if err != nil {
			fmt.Println(err)
		}

		tradeEntity := &model.ExchangeTrades{
			Id:              id,
			Symbol:          utils.NewString(req.Symbol),
			CoinSymbol:      leftCurrency,
			BaseSymbol:      rightCurrency,
			TradeType:       utils.NewInt64(int64(model.TradeTypeOffline.INT())),
			OfflineTradeId:  utils.NewString(trade.OfflineTradeId),
			Price:           amount.Div(quantity),
			Amount:          amount,
			NeedHedgeAmount: decimal.Zero,
			BuyTurnover:     buyTurnover,
			BuyFee:          buyFee,
			SellTurnover:    sellTurnover,
			SellFee:         sellFee,
			Direction:       utils.NewString(fmt.Sprintf("%d", req.Direction)),
			BuyOrderId:      utils.NewString(""),
			SellOrderId:     utils.NewString(""),
			Status:          utils.NewInt64(int64(model.TradeStatusDone)),
			Time:            utils.NewInt64(tradeTime.Unix()),
			UpdateTime:      now,
			MemberId:        utils.NewInt64(uid),
			CoinPrice:       coinMarketPrice,
			BasePrice:       baseMarketPrice,
		}

		if _, e := l.svcCtx.ExchangeTradesModel.Insert(l.ctx, session, tradeEntity); e != nil {
			return e
		}

		var averagePrice decimal.Decimal
		var profit decimal.Decimal
		symbolRepository := repository.NewSymbolRepository(l.ctx, l.svcCtx)
		//均价计算 收益计算 todo
		if req.Direction == 0 { //买入，计算新的均价 还有 币种收益
			averagePrice, profit, err = symbolRepository.BuyAveragePrice(amount, quantity, leftAccountEntity, rightAccountEntity, baseMarketPrice)
			if err != nil {
				return err
			}
			profitEntity := &model.Profits{
				MemberId:        utils.NewInt64(memberId),
				DetailId:        utils.NewInt64(id),
				ProfitTime:      utils.NewTime(timeX),
				TransactionType: utils.NewInt64(int64(model.TransTypeExchange)),
				Symbol:          rightCurrency,
				Profit:          profit,
				Direction:       utils.NewInt64(int64(req.Direction)),
				Status:          1,
			}

			if _, e := l.svcCtx.ProfitsModel.Insert(l.ctx, session, profitEntity); e != nil {
				return e
			}
		} else {
			averagePrice, profit, err = symbolRepository.SellAveragePrice(amount.Sub(fee), quantity, leftAccountEntity.AveragePrice, rightAccountEntity, baseMarketPrice)
			if err != nil {
				return err
			}
			profitEntity := &model.Profits{
				MemberId:        utils.NewInt64(memberId),
				DetailId:        utils.NewInt64(id),
				ProfitTime:      utils.NewTime(timeX),
				TransactionType: utils.NewInt64(int64(model.TransTypeExchange)),
				Symbol:          leftCurrency,
				Profit:          profit,
				Direction:       utils.NewInt64(int64(req.Direction)),
				Status:          1,
			}
			if _, e := l.svcCtx.ProfitsModel.Insert(l.ctx, session, profitEntity); e != nil {
				return e
			}
		}
		fmt.Println("averagePrice", averagePrice)
		var leftAccountBuilder, rightAccountBuilder squirrel.UpdateBuilder
		//3 account
		if req.Direction == 0 { //买入 left 成本变化
			//leftAccountEntity.AddBalance(quantity)
			//leftAccountEntity.UpdateAveragePrice(averagePrice)
			//rightAccountEntity.SubBalance(amount)
			leftAccountBuilder = model.NewAccountsUpdater(leftAccountEntity).
				AddBalance(quantity).
				UpdateAveragePrice(averagePrice).
				Builder()
			rightAccountBuilder = model.NewAccountsUpdater(rightAccountEntity).
				SubBalance(amount).
				Builder()
		} else { //卖出 right 成本变化
			//leftAccountEntity.SubBalance(quantity)
			//rightAccountEntity.AddBalance(amount.Sub(fee))
			//rightAccountEntity.UpdateAveragePrice(averagePrice)
			leftAccountBuilder = model.NewAccountsUpdater(leftAccountEntity).
				SubBalance(quantity).
				Builder()
			rightAccountBuilder = model.NewAccountsUpdater(rightAccountEntity).
				AddBalance(amount.Sub(fee)).
				UpdateAveragePrice(averagePrice).
				Builder()
		}
		if _, e := l.svcCtx.AccountsModel.UpdateByBuilder(l.ctx, session, leftAccountEntity, leftAccountBuilder); e != nil {
			return e
		}

		if _, e := l.svcCtx.AccountsModel.UpdateByBuilder(l.ctx, session, rightAccountEntity, rightAccountBuilder); e != nil {
			return e
		}

		//4. 资金流水
		if req.Direction == 0 { //买
			leftTradeFeeTransaction := &model.MemberTransactions{
				Amount:          quantity,
				Symbol:          leftCurrency,
				MemberId:        utils.NewInt64(memberId),
				AccountId:       utils.NewInt64(leftAccountEntity.Id),
				DetailId:        utils.NewInt64(id),
				TransactionType: utils.NewInt64(int64(model.TransTypeExchange)),
				Fee:             quantity,
				DiscountFee:     decimal.Zero,
				RealFee:         quantity,
				PreBalance:      leftAccountEntity.Balance,
				Balance:         leftAccountEntity.Balance.Add(quantity),
				PreFrozenBal:    leftAccountEntity.FrozenBalance,
				FrozenBal:       leftAccountEntity.FrozenBalance,
				AveragePrice:    leftAccountEntity.AveragePrice,
				CoinPrice:       coinMarketPrice,
			}
			if _, err = l.svcCtx.MemberTransactionsModel.Insert(l.ctx, session, leftTradeFeeTransaction); err != nil {
				return err
			}

			rightTradeFeeTransaction := &model.MemberTransactions{
				Amount:          amount.Sub(fee).Neg(),
				Symbol:          rightCurrency,
				MemberId:        utils.NewInt64(memberId),
				AccountId:       utils.NewInt64(rightAccountEntity.Id),
				DetailId:        utils.NewInt64(id),
				TransactionType: utils.NewInt64(int64(model.TransTypeExchange)),
				Fee:             amount.Sub(fee).Neg(),
				DiscountFee:     decimal.Zero,
				RealFee:         amount.Sub(fee).Neg(),
				PreBalance:      rightAccountEntity.Balance,
				Balance:         rightAccountEntity.Balance.Sub(amount.Sub(fee)),
				PreFrozenBal:    rightAccountEntity.FrozenBalance,
				FrozenBal:       rightAccountEntity.FrozenBalance,
				AveragePrice:    rightAccountEntity.AveragePrice,
				CoinPrice:       baseMarketPrice,
			}
			if _, err = l.svcCtx.MemberTransactionsModel.Insert(l.ctx, session, rightTradeFeeTransaction); err != nil {
				return err
			}

			feeTransaction := &model.MemberTransactions{
				Amount:          fee.Neg(),
				Symbol:          rightCurrency,
				MemberId:        utils.NewInt64(memberId),
				AccountId:       utils.NewInt64(rightAccountEntity.Id),
				DetailId:        utils.NewInt64(id),
				TransactionType: utils.NewInt64(int64(model.TransTypeFee)),
				Fee:             fee.Neg(),
				DiscountFee:     decimal.Zero,
				RealFee:         fee.Neg(),
				PreBalance:      rightAccountEntity.Balance,
				Balance:         rightAccountEntity.Balance.Sub(fee),
				PreFrozenBal:    rightAccountEntity.FrozenBalance,
				FrozenBal:       rightAccountEntity.FrozenBalance,
				AveragePrice:    rightAccountEntity.AveragePrice,
				CoinPrice:       baseMarketPrice,
			}
			if _, err = l.svcCtx.MemberTransactionsModel.Insert(l.ctx, session, feeTransaction); err != nil {
				return err
			}
		} else { //卖
			leftTradeFeeTransaction := &model.MemberTransactions{
				Amount:          quantity.Neg(),
				Symbol:          leftCurrency,
				MemberId:        utils.NewInt64(memberId),
				AccountId:       utils.NewInt64(leftAccountEntity.Id),
				DetailId:        utils.NewInt64(id),
				TransactionType: utils.NewInt64(int64(model.TransTypeExchange)),
				Fee:             quantity.Neg(),
				DiscountFee:     decimal.Zero,
				RealFee:         quantity.Neg(),
				PreBalance:      leftAccountEntity.Balance,
				Balance:         leftAccountEntity.Balance.Sub(quantity),
				PreFrozenBal:    leftAccountEntity.FrozenBalance,
				FrozenBal:       leftAccountEntity.FrozenBalance,
				AveragePrice:    leftAccountEntity.AveragePrice,
				CoinPrice:       coinMarketPrice,
			}
			if _, err = l.svcCtx.MemberTransactionsModel.Insert(l.ctx, session, leftTradeFeeTransaction); err != nil {
				return err
			}

			rightTradeFeeTransaction := &model.MemberTransactions{
				Amount:          amount,
				Symbol:          rightCurrency,
				MemberId:        utils.NewInt64(memberId),
				AccountId:       utils.NewInt64(rightAccountEntity.Id),
				DetailId:        utils.NewInt64(id),
				TransactionType: utils.NewInt64(int64(model.TransTypeExchange)),
				Fee:             decimal.Zero,
				DiscountFee:     decimal.Zero,
				RealFee:         decimal.Zero,
				PreBalance:      rightAccountEntity.Balance,
				Balance:         rightAccountEntity.Balance.Add(amount.Sub(fee)),
				PreFrozenBal:    rightAccountEntity.FrozenBalance,
				FrozenBal:       rightAccountEntity.FrozenBalance,
				AveragePrice:    rightAccountEntity.AveragePrice,
				CoinPrice:       baseMarketPrice,
			}
			if _, err = l.svcCtx.MemberTransactionsModel.Insert(l.ctx, session, rightTradeFeeTransaction); err != nil {
				return err
			}

			feeTransaction := &model.MemberTransactions{
				Amount:          fee.Neg(),
				Symbol:          rightCurrency,
				MemberId:        utils.NewInt64(memberId),
				AccountId:       utils.NewInt64(rightAccountEntity.Id),
				DetailId:        utils.NewInt64(id),
				TransactionType: utils.NewInt64(int64(model.TransTypeFee)),
				Fee:             fee.Neg(),
				DiscountFee:     decimal.Zero,
				RealFee:         fee.Neg(),
				PreBalance:      rightAccountEntity.Balance,
				Balance:         rightAccountEntity.Balance.Sub(fee),
				PreFrozenBal:    rightAccountEntity.FrozenBalance,
				FrozenBal:       rightAccountEntity.FrozenBalance,
				AveragePrice:    rightAccountEntity.AveragePrice,
				CoinPrice:       baseMarketPrice,
			}
			if _, err = l.svcCtx.MemberTransactionsModel.Insert(l.ctx, session, feeTransaction); err != nil {
				return err
			}
		}
		return nil
	})
	if e != nil {
		//xerror.NewCodeError(1500, e.Error(), "")
		l.Logger.Errorf("NewSecondOfflineTradeInput err: %s", err)
		return xerror.ErrorOfflineTradeUpdateFail
	}
	return nil
}
