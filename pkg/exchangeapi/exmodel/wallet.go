package exmodel

type WalletType string

const (
	WalletTypeSpot  = "SPOT"
	WalletTypeUnify = "UNIFY"
)

func (o WalletType) String() string {
	return string(o)
}

type Balance struct {
	WalletType
	Details []*BalanceDetail
}

type BalanceDetail struct {
	Currency
	Total     string
	Available string
}

type QueryBalanceReq struct {
	WalletType
}

type QueryBalanceRsp struct {
	Balances []*Balance
}
