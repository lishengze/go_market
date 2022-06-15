package exmodel

type Currency string

func NewCurrency(s string) Currency {
	return Currency(s)
}

func (o Currency) String() string {
	return string(o)
}
