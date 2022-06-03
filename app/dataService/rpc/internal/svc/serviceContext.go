package svc

import (
	"encoding/json"
	"market_server/app/dataService/rpc/internal/config"
	"market_server/common/nacosAdapter"
	"sync"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	Config      config.Config
	NacosClient *nacosAdapter.Client

	SymbolInfoMap   sync.Map
	CurrencyInfoMap sync.Map
}

func NewServiceContext(c config.Config) *ServiceContext {
	sc := &ServiceContext{
		Config:          c,
		NacosClient:     nacosAdapter.NewClient(c.NacosConfig),
		SymbolInfoMap:   sync.Map{},
		CurrencyInfoMap: sync.Map{},
	}

	//symbol
	sc.getSymbol()
	//currency
	sc.getCurrency()

	return sc
}

func (o *ServiceContext) getSymbol() {
	res, err := o.NacosClient.GetSymbol("")
	if err != nil {
		logx.Errorf("DataService get Symbol err:%+v", err)
		panic(err)
	}

	if len(res) > 0 {
		for _, v := range res {
			var symbol nacosAdapter.Symbol
			copier.Copy(&symbol, v)
			o.SymbolInfoMap.Store(v.SymbolID, &symbol)
		}
	}
	o.NacosClient.ListenSymbol(o.symbolChange)
}

func (o *ServiceContext) symbolChange(namespace, group, dataId, data string) {
	logx.Infof("DataService symbolChange namespace:%s, group:%s, dataId:%s", namespace, group, dataId)

	var SymbolParam []*nacosAdapter.Symbol
	err := json.Unmarshal([]byte(data), &SymbolParam)
	if err != nil {
		logx.Errorf("DataService symbolChange namespace:%s, group:%s, dataId:%s, data:%s  Unmarshal err:%+v", err)
	}

	for _, v := range SymbolParam {
		o.SymbolInfoMap.Store(v.SymbolID, v)
	}
}

func (o *ServiceContext) getCurrency() {
	res, err := o.NacosClient.GetCurrency("")
	if err != nil {
		logx.Errorf("DataService get currency err:%+v", err)
		panic(err)
	}

	if len(res) > 0 {
		for _, v := range res {
			var currency nacosAdapter.Currency
			copier.Copy(&currency, v)
			o.CurrencyInfoMap.Store(v.CurrencyID, &currency)
		}
	}
	o.NacosClient.ListenCurrency(o.currencyChange)
}

func (o *ServiceContext) currencyChange(namespace, group, dataId, data string) {
	logx.Infof("DataService currencyChange namespace:%s, group:%s, dataId:%s", namespace, group, dataId)

	var CurrencyParam []*nacosAdapter.Currency
	err := json.Unmarshal([]byte(data), &CurrencyParam)
	if err != nil {
		logx.Errorf("DataService currencyChange namespace:%s, group:%s, dataId:%s, data:%s  Unmarshal err:%+v", err)
	}

	for _, v := range CurrencyParam {
		o.SymbolInfoMap.Store(v.CurrencyID, v)
	}
}
