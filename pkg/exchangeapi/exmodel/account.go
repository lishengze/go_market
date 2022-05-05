package exmodel

var EmptyAccountConfig = AccountConfig{}

type AccountConfig struct {
	Alias          string
	Key            string
	Secret         string
	PassPhrase     string
	SubAccountName string // FTX sub_account option
}
