package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"market_server/app/admin/api/internal/svc"
	"market_server/app/admin/model"
	"market_server/common/nacosAdapter"
	"market_server/common/xerror"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	"github.com/zeromicro/go-zero/core/logx"
)

type SymbolRepository struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSymbolRepository(ctx context.Context, svcCtx *svc.ServiceContext) *SymbolRepository {
	return &SymbolRepository{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type MarketDetail struct {
	MarketPrice float64
	Timestamp   uint64
}

var TryCnt = 2

func (l *SymbolRepository) GetSymbolMarketPrice(coin string) (decimal.Decimal, error) {
	var cnt int
	for {
		cnt = cnt + 1
		price, err := l.GetSymbolMarketPrice1(coin)
		if err != nil {
			if cnt <= TryCnt {
				continue
			} else {
				return price, err
			}
		} else {
			return price, err
		}
	}
}

//获取市价
func (l *SymbolRepository) GetSymbolMarketPrice1(coin string) (decimal.Decimal, error) {
	var redisConn = l.svcCtx.PriceRedis
	var price decimal.Decimal
	if coin == "USD" {
		return decimal.NewFromInt(1), nil
	}
	//HUSD_USD, 没有直接的数据，你可以通过 BTC_HUSD和 BTC_USD 这两个币对的数据算出来
	if coin == "HUSD" {
		btcHusd, err := l.GetRedisPrice("BTC_HUSD")
		if err != nil {
			return decimal.Decimal{}, err
		}
		btcUsd, err := l.GetRedisPrice("BTC_USD")
		if err != nil {
			return decimal.Decimal{}, err
		}
		if btcHusd.Cmp(decimal.Zero) != 0 {
			price = btcUsd.DivRound(btcHusd, 16)
			return price, nil
		}
		return decimal.Decimal{}, errors.New("get HUSD market price err")
	}
	pre := "marketPrice:"
	coinUSD := pre + coin + "_USD"
	coinUSDT := pre + coin + "_USDT"
	symbol := pre + "USDT_USD"
	var basicInfo, basicInfoV2, basicInfoV3 MarketDetail
	coinV1, _ := redisConn.Get(coinUSD)
	coinV2, _ := redisConn.Get(coinUSDT)
	coinV3, _ := redisConn.Get(symbol)
	if coinV1 != "" {
		if err := json.Unmarshal([]byte(coinV1), &basicInfo); err != nil {
			return decimal.Zero, errors.New(fmt.Sprintf("%s get marketPrice from cache error", coinUSD))
		}
	}
	if coinV2 != "" {
		if err := json.Unmarshal([]byte(coinV2), &basicInfoV2); err != nil {
			return decimal.Zero, errors.New(fmt.Sprintf("%s get marketPrice from cache error", coinUSDT))
		}
	}
	if coinV3 != "" {
		if err := json.Unmarshal([]byte(coinV3), &basicInfoV3); err != nil {
			return decimal.Zero, errors.New(fmt.Sprintf("%s get marketPrice from cache error", symbol))
		}
	}
	if coinV1 == "" {
		if coinV2 != "" && coinV3 != "" {
			_, e := l.Judge(basicInfoV2.Timestamp, coin)
			if e != nil {
				return decimal.Decimal{}, e
			}
			price1 := decimal.NewFromFloat(basicInfoV2.MarketPrice)
			price2 := decimal.NewFromFloat(basicInfoV3.MarketPrice)
			price = price1.Mul(price2)
			return price, nil
		} else {
			return decimal.Decimal{}, errors.New(fmt.Sprintf("redis get Symbol %s  MarketPrice error", coin))
		}
	}
	//判断取USD还是USDT
	timeS := int64(basicInfo.Timestamp) - int64(basicInfoV2.Timestamp)
	//取USDT
	if basicInfoV2.Timestamp > 0 && timeS > 0 && timeS-int64(l.svcCtx.Config.TimeoutMarketPrice.Ctime) > 0 {
		_, e := l.Judge(basicInfoV2.Timestamp, coin)
		if e != nil {
			return decimal.Decimal{}, e
		}
		price1 := decimal.NewFromFloat(basicInfoV2.MarketPrice)
		price2 := decimal.NewFromFloat(basicInfoV3.MarketPrice)
		price = price1.Mul(price2)
		return price, nil
	}
	//取USD
	_, e := l.Judge(basicInfo.Timestamp, coin)
	if e != nil {
		return decimal.Decimal{}, e
	}
	return decimal.NewFromFloat(basicInfo.MarketPrice), nil
}

func (l *SymbolRepository) GetRedisPrice(symbol string) (decimal.Decimal, error) {
	redisSymbol := "marketPrice:" + symbol
	var redisConn = l.svcCtx.PriceRedis
	coinV1, _ := redisConn.Get(redisSymbol)
	if coinV1 != "" {
		basicInfo := &MarketDetail{}
		if err := json.Unmarshal([]byte(coinV1), &basicInfo); err != nil {
			return decimal.Zero, errors.New(fmt.Sprintf("%s get marketPrice from cache error", redisSymbol))
		}
		//infra.Logger.Info("member redis get",
		//	zap.Any("coinV1", coinV1), zap.Any("basicInfo", basicInfo))
		price := decimal.NewFromFloat(basicInfo.MarketPrice)
		return price, nil
	}
	return decimal.Decimal{}, errors.New("get" + symbol + "market price error")
}

func (l *SymbolRepository) Judge(marketTime uint64, coin string) (bool, error) {
	nowTime := time.Now().UnixNano()
	timeAdd := uint64(nowTime) - marketTime
	if len(l.svcCtx.Config.TimeoutMarketPrice.Coins) > 0 {
		for _, v := range l.svcCtx.Config.TimeoutMarketPrice.Coins {
			if coin == v {
				return true, nil
			}
		}
	}
	if timeAdd > l.svcCtx.Config.TimeoutMarketPrice.Mtime {
		return false, errors.New(fmt.Sprintf("redis  Symbol: %s market price not update for long time", coin))
	}
	return true, nil
}

func (l *SymbolRepository) GetSymbolInfoByName(symbol string) (symbolInfo *nacosAdapter.Symbol, err error) {
	symbols, err := l.svcCtx.Parameters.GetSymbol(symbol)
	if err != nil {
		err = xerror.NewCodeError(1201, fmt.Sprintf("Get Symbol info error"), "")
		l.Logger.Errorf("GetSymbolInfoByName fail. symbol: %s, err: %s", symbol, err)
		return
	}
	symbolFound := false
	for _, v := range symbols {
		if strings.ToUpper(symbol) == v.SymbolID {
			symbolFound = true
			symbolInfo = v
			break
		}
	}
	if !symbolFound {
		//xerror.NewCodeError(1201, fmt.Sprintf("Symbol:%s is not configure", symbol), "")
		err = xerror.ErrorSymbolNotFound
		l.Logger.Errorf("GetSymbolInfoByName fail. symbol: %s not found", symbol)
		return
	}
	return
}

/**
eg:ETH_BTC
amount : 花费的钱 BTC (包含手续费)
quantity：买入的量 ETH
account: ETH account
baseAccount: BTC account
baseMarketPrice: BTC市价
*/
func (l *SymbolRepository) BuyAveragePrice(amount decimal.Decimal, quantity decimal.Decimal, account *model.Accounts, baseAccount *model.Accounts, baseMarketPrice decimal.Decimal) (averagePrice decimal.Decimal, profit decimal.Decimal, err error) {
	amountAdd := amount.Mul(baseMarketPrice)
	balanceBefore := account.Balance.Add(account.FrozenBalance)
	if baseAccount.AveragePrice.Cmp(decimal.Zero) == 0 {
		aPrice, err := l.GetSymbolMarketPrice(baseAccount.CoinUnit.String)
		if err != nil {
			aPrice = decimal.Zero //todo 应该报错
			return aPrice, decimal.Zero, err
		}
		baseAccount.AveragePrice = aPrice
	}
	amountBefore := balanceBefore.Mul(account.AveragePrice)

	balanceTotal := balanceBefore.Add(quantity)

	averagePrice = amountAdd.Add(amountBefore).DivRound(balanceTotal, 16)

	profit = baseMarketPrice.Sub(baseAccount.AveragePrice).Mul(amount)
	return
}

func (l *SymbolRepository) AveragePriceAndProfitFromBuy(amount decimal.Decimal, quantity decimal.Decimal, account *model.Accounts, baseAccount *model.Accounts, baseMarketPrice decimal.Decimal) (averagePrice decimal.Decimal, profit decimal.Decimal, err error) {
	baseAveragePrice := baseAccount.AveragePrice
	if baseAveragePrice.Cmp(decimal.Zero) == 0 {
		baseAveragePriceLatest, err := l.GetSymbolMarketPrice(baseAccount.CoinUnit.String)
		if err != nil {
			baseAveragePriceLatest = decimal.Zero //todo 应该报错
			return baseAveragePriceLatest, decimal.Zero, err
		}
		baseAccount.AveragePrice = baseAveragePriceLatest
		baseAveragePrice = baseAveragePriceLatest
	}

	profit = l.ProfitFromBuy(amount, baseMarketPrice, baseAveragePrice)
	averagePrice, err = AveragePriceFromBuy(amount, baseMarketPrice, quantity, account.Balance.Add(account.FrozenBalance), account.AveragePrice)

	return
}

// (baseMarketPrice - baseAveragePrice) * amount
func (l *SymbolRepository) ProfitFromBuy(amount, baseMarketPrice, baseAveragePrice decimal.Decimal) decimal.Decimal {
	return baseMarketPrice.Sub(baseAveragePrice).Mul(amount)
}

// ((amount * baseMarketPrice) + quantity * averagePrice) / (quantity + balance)
func AveragePriceFromBuy(amount, baseMarketPrice decimal.Decimal, quantity, balance, averagePrice decimal.Decimal) (averagePriceLatest decimal.Decimal, err error) {
	amountNow := amount.Mul(baseMarketPrice)
	amountBefore := balance.Mul(averagePrice)
	balanceNow := quantity.Add(balance)
	if balanceNow.Cmp(decimal.Zero) == 0 {
		err = errors.New("Both quantity and balance are not zero")
		return
	}
	fmt.Println("averagePriceLatest", amountNow, amountBefore, balanceNow)
	averagePriceLatest = amountNow.Add(amountBefore).DivRound(balanceNow, 16)
	return
}

/**
eg:ETH_BTC
amount : 增量 (增加的BTC 减去 fee)
quantity: ETH卖出量
coinAveragePrice: ETH成本价
account: BTC account
baseMarketPrice: BTC市价
*/
func (l *SymbolRepository) SellAveragePrice(amount decimal.Decimal, quantity decimal.Decimal, coinAveragePrice decimal.Decimal, account *model.Accounts, baseMarketPrice decimal.Decimal) (averagePrice decimal.Decimal, profit decimal.Decimal, err error) {
	amountAdd := amount.Mul(baseMarketPrice)
	balanceBefore := account.Balance.Add(account.FrozenBalance)
	if account.AveragePrice.Cmp(decimal.Zero) == 0 {
		aPrice, err := l.GetSymbolMarketPrice(account.CoinUnit.String)
		if err != nil {
			aPrice = decimal.Zero //todo 应该报错
			return aPrice, decimal.Zero, err
		}
		account.AveragePrice = aPrice
	}
	amountBefore := balanceBefore.Mul(account.AveragePrice)
	balanceTotal := balanceBefore.Add(amount)
	averagePrice = amountAdd.Add(amountBefore).DivRound(balanceTotal, 16)

	coinPriceOld := quantity.Mul(coinAveragePrice)
	profit = amountAdd.Sub(coinPriceOld)
	return
}

func (l *SymbolRepository) AveragePriceAndProfitFromSell(amount decimal.Decimal, quantity decimal.Decimal, account *model.Accounts, baseAccount *model.Accounts, baseMarketPrice decimal.Decimal) (averagePrice decimal.Decimal, profit decimal.Decimal, err error) {
	if baseAccount.AveragePrice.Cmp(decimal.Zero) == 0 {
		averagePriceLatest, err := l.GetSymbolMarketPrice(baseAccount.CoinUnit.String)
		if err != nil {
			averagePriceLatest = decimal.Zero //todo 应该报错
			return averagePriceLatest, decimal.Zero, err
		}
		baseAccount.AveragePrice = averagePrice
	}

	averagePrice, err = AveragePriceFromSell(amount, baseAccount.Balance.Add(baseAccount.FrozenBalance), baseMarketPrice, baseAccount.AveragePrice)
	if err != nil {
		return
	}

	profit = ProfitFromSell(amount, baseMarketPrice, quantity, account.AveragePrice)
	return
}

// ((amount * baseMarketPrice) + balance * baseAveragePrice) / (balance + amount)
func AveragePriceFromSell(amount, balance, baseMarketPrice, baseAveragePrice decimal.Decimal) (averagePriceLatest decimal.Decimal, err error) {
	in := amount.Mul(baseMarketPrice)
	priceOld := balance.Mul(baseAveragePrice)
	priceNew := in.Add(priceOld)

	balanceNew := balance.Add(amount)
	averagePriceLatest = priceNew.DivRound(balanceNew, 16)
	return
}

// baseMarketPrice * amount - averagePrice * quantity
func ProfitFromSell(amount, baseMarketPrice decimal.Decimal, quantity, averagePrice decimal.Decimal) decimal.Decimal {
	in := amount.Mul(baseMarketPrice)
	out := quantity.Mul(averagePrice)

	return in.Sub(out)
}

/**
因为线下交易修改，导致均价变化
eg ETH_BTC
amount : 花费的钱 BTC (包含手续费)
quantity：买入的量 ETH
ETH account
baseMarketPrice: BTC市价
*/
func (l *SymbolRepository) BuyAveragePriceForUpdate(amount decimal.Decimal, quantity decimal.Decimal, account *model.Accounts, baseMarketPrice decimal.Decimal) (averagePrice decimal.Decimal) {
	amountAdd := amount.Mul(baseMarketPrice)
	fmt.Println("amountAdd", amountAdd, amount, baseMarketPrice)
	balanceTotal := account.Balance.Add(account.FrozenBalance)
	amountTotal := balanceTotal.Mul(account.AveragePrice)
	amountSub := amountTotal.Sub(amountAdd)
	fmt.Println("amountSub", amountSub)
	balanceSub := balanceTotal.Sub(quantity)
	fmt.Println("balanceSub", balanceSub)
	if balanceSub.Cmp(decimal.Zero) == 0 {
		averagePrice = account.AveragePrice
		return
	}
	averagePrice = amountSub.DivRound(balanceSub, 16)
	return
}

/**
因为线下交易修改，导致均价变化 base 币种增加
eg:ETH_BTC 卖出ETH 买入BTC
amount：增加的BTC减去fee
account: BTC account
baseMarketPrice: BTC市价
*/
func (l *SymbolRepository) SellAveragePriceForUpdate(amount decimal.Decimal, account *model.Accounts, baseMarketPrice decimal.Decimal) (averagePrice decimal.Decimal) {
	amountAdd := amount.Mul(baseMarketPrice)
	balanceTotal := account.Balance.Add(account.FrozenBalance)
	amountTotal := balanceTotal.Mul(account.AveragePrice)
	amountSub := amountTotal.Sub(amountAdd)
	balanceSub := balanceTotal.Sub(amount)
	if balanceSub.Cmp(decimal.Zero) == 0 {
		averagePrice = account.AveragePrice
		return
	}
	averagePrice = amountSub.DivRound(balanceSub, 16)
	return
}
