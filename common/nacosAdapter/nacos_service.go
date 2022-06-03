package nacosAdapter

import (
	"encoding/json"
	"market_server/common/xerror"

	"github.com/zeromicro/go-zero/core/logx"
)

//var ErrConfigNotFound = xerror.NewE("config not found")

// GetCurrency 币种参数配置
func (c *Client) GetCurrency(currencyId string) ([]*Currency, error) {
	var currencies []*Currency
	currencyContent, err := c.getConfigContent(CURRENCY_PARAMSV2, BCTS_GROUP)
	if err != nil {
		logx.Errorf("nacos GetCurrency err:%+v", err)
		return nil, err
	}
	err = json.Unmarshal([]byte(currencyContent), &currencies)
	if err != nil {
		logx.Errorf("nacos Unmarshal err:%+v", err)
		return nil, err
	}
	if currencyId == "" {
		return currencies, nil
	}
	var retCurrencies []*Currency
	for _, currency := range currencies {
		if currency.CurrencyID == currencyId {
			retCurrencies = append(retCurrencies, currency)
			break
		}
	}
	if len(retCurrencies) == 0 {
		return retCurrencies, xerror.ErrorRecordNotFound
	}
	return retCurrencies, nil
}

func (c *Client) SetCurrency(currencyParams []*Currency) error {
	currencyBytes, _ := json.Marshal(currencyParams)
	currencyStr := string(currencyBytes)
	err := c.publishConfig(CURRENCY_PARAMSV2, BCTS_GROUP, currencyStr)
	return err
}

func (c *Client) ListenCurrency(currencyChanged func(namespace, group, dataId, data string)) {
	go c.listenConfig(CURRENCY_PARAMSV2, BCTS_GROUP, currencyChanged)
}

// GetHedging 对冲平台参数
func (c *Client) GetHedging(hedgingId string) ([]*Hedging, error) {
	var hedgings []*Hedging
	hedgingContent, err := c.getConfigContent(HEDGE_PARAMS_V2, BCTS_GROUP)
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	err = json.Unmarshal([]byte(hedgingContent), &hedgings)
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	if hedgingId == "" {
		return hedgings, nil
	}
	var retHedgings []*Hedging
	for _, hedging := range hedgings {
		if hedging.PlatformID == hedgingId {
			retHedgings = append(retHedgings, hedging)
		}
	}
	if len(retHedgings) == 0 {
		return retHedgings, xerror.ErrorRecordNotFound
	}
	return retHedgings, nil
}

func (c *Client) GetSymbol(symbolId string) ([]*Symbol, error) {
	var symbols []*Symbol
	symbolContent, err := c.getConfigContent(SYMBOL_PARAMS_V2, BCTS_GROUP)
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	err = json.Unmarshal([]byte(symbolContent), &symbols)
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	if symbolId == "" {
		return symbols, nil
	}
	var retSymbols []*Symbol
	for _, symbol := range symbols {
		if symbol.SymbolID == symbolId {
			retSymbols = append(retSymbols, symbol)
			break
		}
	}
	if len(retSymbols) == 0 {
		return retSymbols, xerror.ErrorRecordNotFound
	}
	return retSymbols, nil
}

func (c *Client) SetSymbol(symbolParams []*Symbol) error {
	symbolBytes, _ := json.Marshal(symbolParams)
	symbolStr := string(symbolBytes)
	err := c.publishConfig(SYMBOL_PARAMS, BCTS_GROUP, symbolStr)
	return err
}

func (c *Client) ListenSymbol(symbolChanged func(namespace, group, dataId, data string)) {
	go c.listenConfig(SYMBOL_PARAMS_V2, BCTS_GROUP, symbolChanged)
}

//交易总开关
func (c *Client) GetTradeSwitch() (tradeSwitch *TradeSwitch, e error) {
	var trade TradeSwitch
	switchContent, e := c.getConfigContent(TRADE_SWITCH_PARAMS, BCTS_GROUP)
	if e != nil {
		logx.Error(e)
		return nil, e
	}
	logx.Info(switchContent)
	e = json.Unmarshal([]byte(switchContent), &trade)
	if e != nil {
		logx.Error(e)
		return nil, e
	}
	return &trade, nil
}

func (c *Client) SetTradeSwitch(tradeSwitch *TradeSwitch) error {
	tradeSwitchBytes, _ := json.Marshal(tradeSwitch)
	switchStr := string(tradeSwitchBytes)
	err := c.publishConfig(TRADE_SWITCH_PARAMS, BCTS_GROUP, switchStr)
	return err
}

func (c *Client) ListenTradeSwitch(symbolChanged func(namespace, group, dataId, data string)) {
	go c.listenConfig(TRADE_SWITCH_PARAMS, BCTS_GROUP, symbolChanged)
}

func (c *Client) GetGroupTradeParam(group string) ([]*GroupTradeParam, error) {
	var groupTradeParams []*GroupTradeParam
	groupTradeParamContent, err := c.getConfigContent(MEMBER_GROUP_TRADE_PARAMS_PARAMS_V2, BCTS_GROUP)
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	err = json.Unmarshal([]byte(groupTradeParamContent), &groupTradeParams)
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	if group == "" {
		return groupTradeParams, nil
	}
	var retGroupTradeParams []*GroupTradeParam
	for _, groupTradeParam := range groupTradeParams {
		if groupTradeParam.GroupCode == group {
			retGroupTradeParams = append(retGroupTradeParams, groupTradeParam)
		}
	}
	if len(retGroupTradeParams) == 0 {
		return retGroupTradeParams, xerror.ErrorRecordNotFound
	}
	return retGroupTradeParams, nil
}

func (c *Client) GetGroupTradeAuth(group string) ([]*GroupTradeAuth, error) {
	var groupTradeAuths []*GroupTradeAuth
	groupTradeAuthContent, err := c.getConfigContent(MEMBER_GROUP_TRADE_AUTH_PARAMS, BCTS_GROUP)
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	err = json.Unmarshal([]byte(groupTradeAuthContent), &groupTradeAuths)
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	if group == "" {
		return groupTradeAuths, nil
	}
	var retGroupTradeAuths []*GroupTradeAuth
	for _, groupTradeAuth := range groupTradeAuths {
		if groupTradeAuth.GroupCode == group {
			retGroupTradeAuths = append(retGroupTradeAuths, groupTradeAuth)
		}
	}
	if len(retGroupTradeAuths) == 0 {
		return retGroupTradeAuths, xerror.ErrorRecordNotFound
	}
	return retGroupTradeAuths, nil
}

// GetBrokerConf 钉钉告警
func (c *Client) GetBrokeConf() (*BrokerConf, error) {
	var confs *BrokerConf
	confContent, err := c.getConfigContent(BROKER_CONF, BCTS_GROUP)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(confContent), &confs)
	if err != nil {
		return nil, err
	}
	return confs, nil
}
