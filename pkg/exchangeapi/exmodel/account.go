package exmodel

var EmptyAccountConfig = AccountConfig{}

type AccountConfig struct {
	Proxy          string
	Alias          string
	Key            string
	Secret         string
	PassPhrase     string
	SubAccountName string // FTX sub_account option
}
