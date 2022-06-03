package offlinetrade

import (
	repository2 "bcts/app/admin/api/internal/repository"
	"bcts/app/admin/api/internal/svc"
	"bcts/app/admin/api/internal/types"
	"bcts/app/admin/model"
	"bcts/common/googleAuth"
	"bcts/common/utils"
	"bcts/common/xerror"
	"context"
	"database/sql"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
	"strings"
)

type validator func(request *types.OfflineTradeReq) error

type BaseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BaseLogic {
	return &BaseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BaseLogic) validateOperator(req *types.OfflineTradeReq) (err error) {
	user, err := l.svcCtx.AdminUserModel.FindOneByName(l.ctx, req.Operator)
	if err != nil {
		err = xerror.NewCodeError(30000, fmt.Sprintf("Operator:%s dost not exsits", req.Operator), fmt.Sprintf("操作员:%s 不存在", req.Operator))
		return
	}

	googleCode := googleAuth.GenerateTwoFactorCode(user.GoogleAuthSecret)
	if l.svcCtx.Config.Environment.Env != "qa" && googleCode != req.TwoFactorCode {
		l.Logger.Errorf("validate TwoFactorCode fail, googleCode: %s, TwoFactorCode: %s", googleCode, req.TwoFactorCode)
		return xerror.ErrorGoogleAuthCode
	}
	return
}

func (l *BaseLogic) validateOfflineTradeID(req *types.OfflineTradeReq) (err error) {
	if req.OfflineTradeID == "" {
		//return 1312 errors.New(fmt.Sprintf("OfflineTradeID is empty"))
		return xerror.ErrorOfflineTradeIdIsEmpty
	}
	return
}

func (l *BaseLogic) validateDirection(req *types.OfflineTradeReq) (err error) {
	if !(req.Direction == 0 || req.Direction == 1) {
		//NewCodeError(300001, fmt.Sprintf("Direction:%d is invalid", req.Direction), "交易方向不正确")
		return xerror.ErrorDirectionInvalid
	}
	return
}

func (l *BaseLogic) validateUID(req *types.OfflineTradeReq) (err error) {
	memeberId, err := strconv.ParseInt(req.UID, 10, 64)
	if err != nil {
		// NewCodeError(1250, fmt.Sprintf("UID:%s can not found", req.UID), "该账号不存在")
		l.Logger.Errorf("validate uid fails. uid: %s", req.UID)
		err = xerror.ErrorMemberNotFound
		return
	}

	_, err = l.svcCtx.MembersModel.FindOne(l.ctx, memeberId)
	if err != nil {
		if err == model.ErrNotFound {
			err = xerror.ErrorMemberNotFound
			return
		}
		// NewCodeError(150, fmt.Sprintf("Unknown error, try again"), "未知错误,请重试")
		l.Logger.Errorf("validate uid fails. uid: %s, err: %s", memeberId, err)
		err = xerror.ErrorSystem
		return
	}
	return
}

func (l *BaseLogic) validateDecimal(req *types.OfflineTradeReq) (err error) {
	_, err = decimal.NewFromString(req.Quantity)
	if err != nil {
		//NewCodeError(1302, fmt.Sprintf("Quantity:%s is invalid", req.Quantity), "数量不合法")
		err = xerror.ErrorQuantityInvalid
		return
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		//NewCodeError(1303, fmt.Sprintf("Amount:%s is invalid", req.Amount), "金额不合法")
		err = xerror.ErrorAmountInvalid
		return
	}

	fee, err := decimal.NewFromString(req.Fee)
	if err != nil {
		//NewCodeError(1304, fmt.Sprintf("Fee:%s is invalid", req.Fee), "手续费不合法")
		err = xerror.ErrorFeeInvalid
		return
	}

	if fee.GreaterThan(amount) {
		//NewCodeError(1201, fmt.Sprintf("Fee:%s Amount:%s", req.Fee, req.Amount), "收费不合法")
		err = xerror.ErrorFeeGtAmount
		return
	}

	return
}

// delete
func (l *BaseLogic) validateFirstStatus(req *types.OfflineTradeReq) (err error) {
	if req.ID == 0 {
		err = xerror.ErrorOfflineTradeNotFound
		return
	}

	record, err := l.svcCtx.OfflineTradeInputModel.FindOne(l.ctx, req.ID)
	if err != nil {
		if err == model.ErrNotFound {
			err = xerror.ErrorOfflineTradeNotFound
		}
		return
	}

	if record.Status.Int64 != int64(model.OfflineTradeInputStatusFirst.INT()) {
		err = xerror.ErrorOfflineTradeNotDelete
		return
	}

	if record.Operator1.String != req.Operator {
		err = xerror.ErrorOfflineTradeNotAllowDelete
		return
	}

	return
}

func (l *BaseLogic) validateSecondTrade(req *types.OfflineTradeReq, trade *model.OfflineTradeInput) (err error) {
	if trade.Status.Int64 == int64(model.OfflineTradeInputStatusSecond.INT()) {
		//NewCodeError(1300, fmt.Sprintf("OfflineTradeID:%s is already end state", req.OfflineTradeID), "线下交易ID已存在")
		err = xerror.ErrorAlreadyEnd
		return
	}

	if trade.Operator1.String == req.Operator {
		//NewCodeError(1237, fmt.Sprintf("OfflineTradeID:%s operator can not the same", req.OfflineTradeID), "")
		err = xerror.ErrorOperatorSame
		return
	}

	return
}

func (l *BaseLogic) GetAccountsByMemeberIDAndCoinName(req *types.OfflineTradeReq, memberId int64) (exchangeCoin *model.Accounts, balanceCoin *model.Accounts, err error) {
	indexOfCurrency := strings.Index(req.Symbol, "_")
	leftCurrency := req.Symbol[:indexOfCurrency]
	rightCurrency := req.Symbol[indexOfCurrency+1:]

	exchangeCoin, err = l.svcCtx.AccountsModel.FindOneByMemberIdCoinUnit(l.ctx, sql.NullInt64{Int64: memberId, Valid: true}, sql.NullString{String: leftCurrency, Valid: true})
	if err != nil {
		//NewCodeError(1301, fmt.Sprintf("Currency:%s is not found", leftCurrency), "交易币不存在")
		err = xerror.ErrorExchangeCoinNotExist
		return
	}

	balanceCoin, err = l.svcCtx.AccountsModel.FindOneByMemberIdCoinUnit(l.ctx, sql.NullInt64{Int64: memberId, Valid: true}, sql.NullString{String: rightCurrency, Valid: true})
	if err != nil {
		//NewCodeError(1301, fmt.Sprintf("Currency:%s is not found", rightCurrency), "结算币不存在")
		err = xerror.ErrorBaseCoinNotExist
		return
	}

	return
}

func (l *BaseLogic) GetCoinPricesByTradeTime(addDay, symbol string) (coinMarketPrice decimal.Decimal, balanceCoin decimal.Decimal, err error) {
	indexOfCurrency := strings.Index(symbol, "_")
	leftCurrency := symbol[:indexOfCurrency]
	rightCurrency := symbol[indexOfCurrency+1:]

	coinPriceEntity, err := l.svcCtx.CoinPricesModel.CoinPriceByTradeTime(l.ctx, addDay, leftCurrency)
	if err != nil || coinPriceEntity == nil || !coinPriceEntity.Price.Valid {
		//NewCodeError(1324, fmt.Sprintf("Currency:%s is get marketPrice err", leftCurrency), "")
		err = xerror.ErrorExchangeCoinPriceNotExist
		return
	}
	coinMarketPrice = coinPriceEntity.Price.Decimal

	basePriceEntity, err := l.svcCtx.CoinPricesModel.CoinPriceByTradeTime(l.ctx, addDay, rightCurrency)
	if err != nil || basePriceEntity == nil || !basePriceEntity.Price.Valid {
		//NewCodeError(1324, fmt.Sprintf("Currency:%s is get marketPrice err", rightCurrency), "")
		err = xerror.ErrorBaseCoinPriceNotExist
		return
	}
	balanceCoin = basePriceEntity.Price.Decimal

	return
}

func (l *BaseLogic) GetSymbolMarketPrices(memberId int64, symbol string) (coinMarketPrice decimal.Decimal, baseMarketPrice decimal.Decimal, err error) {
	indexOfCurrency := strings.Index(symbol, "_")
	leftCurrency := symbol[:indexOfCurrency]
	rightCurrency := symbol[indexOfCurrency+1:]

	symbolRepository := repository2.NewSymbolRepository(l.ctx, l.svcCtx)
	coinMarketPrice, err = symbolRepository.GetSymbolMarketPrice(leftCurrency)
	alarmRepository := repository2.NewAlarmRepository(l.ctx, l.svcCtx)
	if err != nil {
		_, _ = alarmRepository.DingPmsAlertMsg(fmt.Sprintf("NewSecondOfflineTradeInput: member: %d, Currency:%s  get marketPrice error", memberId, leftCurrency))
		//xerror.NewCodeError(1324, fmt.Sprintf("Currency:%s is get marketPrice err", leftCurrency), "")
		err = xerror.ErrorExchangeCoinPriceNotExist
		return
	}

	baseMarketPrice, err = symbolRepository.GetSymbolMarketPrice(rightCurrency)
	if err != nil {
		_, _ = alarmRepository.DingPmsAlertMsg(fmt.Sprintf("NewSecondOfflineTradeInput: member: %d,  Currency:%s  get marketPrice error", memberId, rightCurrency))
		//NewCodeError(1324, fmt.Sprintf("Currency:%s is get marketPrice err", rightCurrency), "")
		err = xerror.ErrorBaseCoinPriceNotExist
		return
	}

	return
}

func (l *BaseLogic) GenerateDistributedId() int64 {
	did := utils.NewDid(&utils.Snowflake{
		Node:     l.svcCtx.Config.Snowflake.Node,
		Epoch:    l.svcCtx.Config.Snowflake.Epoch,
		NodeBits: l.svcCtx.Config.Snowflake.NodeBits,
	})
	return did.Generate()
}

func (l *BaseLogic) GenerateTradeBusinessNo() string {
	id := l.GenerateDistributedId()
	idStr := strconv.FormatInt(id, 10)
	return idStr[len(idStr)-8:]
}

func ConvertModel(model *model.OfflineTradeInput) *types.OfflineTradeInput {
	if model == nil {
		return nil
	}

	return &types.OfflineTradeInput{
		Id:              model.Id,
		TradeBusinessNo: model.TradeBusinessNo,
		TradeTime:       model.TradeTime.Time.Unix(),
		OfflineTradeId:  model.OfflineTradeId,
		Uid:             model.Uid.Int64,
		UserName:        model.UserName.String,
		Name:            model.Name.String,
		Symbol:          model.Symbol.String,
		Direction:       model.Direction.Int64,
		Quantity:        model.Quantity.String(),
		Amount:          model.Amount.String(),
		Price:           model.Price.String(),
		Fee:             model.Fee.String(),
		TradeId:         model.TradeId.String,
		Remarks:         model.Remarks.String,
		Status:          model.Status.Int64,
		Operator1:       model.Operator1.String,
		Operator2:       model.Operator2.String,
		ReviseOperator:  model.ReviseOperator.String,
		OperateTime:     model.OperateTime.Time.Unix(),
	}
}

func ConvertModels(models []*model.OfflineTradeInput) []*types.OfflineTradeInput {
	resp := make([]*types.OfflineTradeInput, 0, len(models))
	if models == nil || len(models) == 0 {
		return resp
	}

	for _, m := range models {
		if m != nil {
			resp = append(resp, ConvertModel(m))
		}
	}

	return resp
}
