package exmodel

import "strings"

type Currency string

func NewCurrency(s string) Currency {
	return Currency(strings.ToUpper(s))
}

func (o Currency) String() string {
	return string(o)
}
