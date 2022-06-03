package nacosAdapter

const BCTS_GROUP = "BCTS"

const (
	Test_PARAMS         = "test_by_huangtao"
	CURRENCY_PARAMS     = "CurrencyParams"    //币种参数
	CURRENCY_PARAMSV2   = "CurrencyParamsV2"  //币种参数
	HEDGE_PARAMS        = "HedgeParams"       //对冲平台参数
	HEDGE_PARAMS_V2     = "HedgeParamsV2"     //对冲平台参数（新版，float64=>string）
	SYMBOL_PARAMS       = "SymbolParams"      //品种参数
	SYMBOL_PARAMS_V2    = "SymbolParamsV2"    //品种参数（新版，float64=>string）
	TRADE_SWITCH_PARAMS = "tradeSwitchParams" //交易总开关
	//MEMBER_GROUP_TRADE_PARAMS_PARAMS = "MemberGroupTradeParams"
	MEMBER_GROUP_TRADE_PARAMS_PARAMS_V2 = "MemberGroupTradeParamsV2"
	MEMBER_GROUP_TRADE_AUTH_PARAMS      = "MemberGroupTradeAuth"
	BROKER_CONF                         = "broker-conf.yaml"
)
