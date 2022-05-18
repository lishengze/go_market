package config

import (
	"encoding/json"

	"github.com/zeromicro/go-zero/core/logx"
)

func ParseJsonHedgerConfig(hedgingContent string) ([]*HedgingConfig, error) {

	var hedgings []*HedgingConfig

	err := json.Unmarshal([]byte(hedgingContent), &hedgings)

	if err != nil {
		logx.Error(err)
		return nil, err
	}

	var retHedgings []*HedgingConfig
	for _, hedging := range hedgings {
		retHedgings = append(retHedgings, hedging)
	}
	return retHedgings, nil
}

func ParseJsonSymbolConfig(symbolContent string) ([]*SymbolConfig, error) {

	var symbols []*SymbolConfig

	err := json.Unmarshal([]byte(symbolContent), &symbols)
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	var retSymbols []*SymbolConfig
	for _, symbol := range symbols {
		retSymbols = append(retSymbols, symbol)
	}
	return retSymbols, nil
}

func ParseJsonMarketRiskConfig(currencyContent string) ([]*MarketRiskConfig, error) {
	var market_risk_configs []*MarketRiskConfig
	err := json.Unmarshal([]byte(currencyContent), &market_risk_configs)
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	var ret_market_risk_configs []*MarketRiskConfig
	for _, market_risk_config := range market_risk_configs {
		ret_market_risk_configs = append(ret_market_risk_configs, market_risk_config)
	}
	return ret_market_risk_configs, nil
}
