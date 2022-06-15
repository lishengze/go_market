package ftx

import (
	"exterior-interactor/pkg/exchangeapi/exmodel"
	"exterior-interactor/pkg/exchangeapi/extools"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
)

type walletManager struct {
	api *NativeApi
}

func NewWalletManager(api *NativeApi) extools.WalletManager {
	return &walletManager{
		api: api,
	}
}

func (o *walletManager) QueryBalance(req exmodel.QueryBalanceReq) (*exmodel.QueryBalanceRsp, error) {
	rsp, err := o.api.GetBalance()
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	if !rsp.Success {
		logx.Error("ftx request not success")
		return nil, fmt.Errorf("ftx request not success")
	}

	balance := &exmodel.Balance{
		WalletType: exmodel.ApiTypeUnified,
	}

	for _, b := range rsp.Result {
		balance.Details = append(balance.Details, &exmodel.BalanceDetail{
			Currency:  exmodel.NewCurrency(b.Coin),
			Total:     fmt.Sprint(b.Total),
			Available: fmt.Sprint(b.Free),
		})
	}

	res := &exmodel.QueryBalanceRsp{
		Balances: []*exmodel.Balance{
			balance,
		}}

	return res, nil
}
