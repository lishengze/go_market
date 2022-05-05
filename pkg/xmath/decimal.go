package xmath

import (
	"fmt"
	"github.com/shopspring/decimal"
)

func MustDecimal(num string) decimal.Decimal {
	d, err := decimal.NewFromString(num)
	if err != nil {
		panic(fmt.Sprintf("cant convert:%s to decimal,err:%v", num, err))
	}
	return d
}

func ZeroFillingString(d decimal.Decimal, place int) string {
	var fill string
	for i := 1; i <= place; i++ {
		fill += "0"
	}

	fracPart := d.Sub(decimal.NewFromInt(d.IntPart())) // 小数部分
	if fracPart.IsZero() && place > 0 {
		return fmt.Sprintf("%s.%s", d.String(), fill)
	}

	return d.String()
}
