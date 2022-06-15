package core

import (
	"context"
	"exterior-interactor/app/opu/model"
	"exterior-interactor/pkg/exchangeapi/exchanges/ftx"
	"exterior-interactor/pkg/exchangeapi/exmodel"
	"exterior-interactor/pkg/exchangeapi/extools"
	"exterior-interactor/pkg/xencrypt"
	"fmt"
)

type accountManager struct {
	*model.Account
	extools.TradeManager
	extools.WalletManager
	outputCh chan<- *exmodel.OrderTradesUpdate
	cancel   context.CancelFunc
	ctx      context.Context
}

func newAccountManager(account *model.Account, proxy string, outputCh chan<- *exmodel.OrderTradesUpdate) (*accountManager, error) {
	var (
		tm extools.TradeManager
		wm extools.WalletManager
	)

	switch account.Exchange {
	case exmodel.FTX.String():
		key, _ := xencrypt.DecryptByAes(account.Key)
		secret, _ := xencrypt.DecryptByAes(account.Secret)
		passphrase, _ := xencrypt.DecryptByAes(account.Passphrase)

		api := ftx.NewNativeApiWithProxy(exmodel.AccountConfig{
			Alias:          account.Alias,
			Key:            key,
			Secret:         secret,
			PassPhrase:     passphrase,
			SubAccountName: account.SubAccountName,
		}, proxy)

		tm = ftx.NewTradeManager(api)
		wm = ftx.NewWalletManager(api)

	default:
		return nil, fmt.Errorf("not support exchange:%s", account.Exchange)
	}

	ctx, cancel := context.WithCancel(context.Background())

	am := &accountManager{
		Account:       account,
		TradeManager:  tm,
		WalletManager: wm,
		outputCh:      outputCh,
		cancel:        cancel,
		ctx:           ctx,
	}

	go am.run()

	return am, nil
}

func (o *accountManager) run() {
	updateCh := o.TradeManager.OutputUpdateCh()
	for {
		select {
		case <-o.ctx.Done():
			return
		case update := <-updateCh:
			o.outputCh <- update
		}
	}
}

func (o *accountManager) close() {
	o.TradeManager.Close()
	o.cancel()
}
