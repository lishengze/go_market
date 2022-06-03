package ordercheck

import (
	"bcts/app/dataService/rpc/dataservice"
	"bcts/common/globalKey"
	"bcts/common/xerror"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

//amountExceptFee decimal.Decimal //金额询价，扣除预估手续费后的可用于实际交易的金额
func CheckUniversal(quoteType int32, volumeStr, amountStr string, feeInfo *dataservice.GetUserFeeInfoRsp) (volume, amount, amountExceptFee decimal.Decimal, err error) {
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
