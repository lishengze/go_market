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

func ParseJsonCurrencyConfig(currencyContent string) ([]*CurrencyConfig, error) {
	var currencies []*CurrencyConfig
	err := json.Unmarshal([]byte(currencyContent), &currencies)
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	var retCurrencies []*CurrencyConfig
	for _, currency := range currencies {
		retCurrencies = append(retCurrencies, currency)
	}
	return retCurrencies, nil
}
