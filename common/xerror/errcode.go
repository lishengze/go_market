package xerror

var (
	//base
	ErrorTryAgain       = NewCodeError(100001, "server is busy, please try again", "服务器开小差了, 请稍后重试")
	ErrorParamError     = NewCodeError(100002, "param error", "参数错误")
	ErrorOverLimit      = NewCodeError(100003, "request over limit", "超出频率限制")
	ErrorDB             = NewCodeError(100004, "DB error", "数据库错误")
	ErrorRecordNotFound = NewCodeError(100005, "record not found", "记录不存在")

	ErrorSymbolIsEmpty   = NewCodeError(100100, "symbol is empty", "合约不能为空")
	ErrorCurrencyIsEmpty = NewCodeError(100100, "currency is empty", "币种不能为空")

	// user 登录相关  000 - 100
	ErrorAreaCodeIsEmpty         = NewCodeError(200000, "area code is empty", "手机区号不能为空")
	ErrorMobileIsEmpty           = NewCodeError(200001, "mobile is empty", "手机号不能为空")
	ErrorEmailIsEmpty            = NewCodeError(200002, "email is empty", "邮箱不能为空")
	ErrorMobileDuplicateRegister = NewCodeError(200003, "mobile is already register", "手机号已经被注册")
	ErrorEmailDuplicateRegister  = NewCodeError(200004, "email is already register", "邮箱已经被注册")
	ErrorMobileUnRegister        = NewCodeError(200005, "mobile is unregister", "手机号未注册")
	ErrorEmailUnRegister         = NewCodeError(200006, "email is unregister", "邮箱未注册")
	ErrorUserNotFound            = NewCodeError(200007, "user not found", "用户不存在")
	ErrorPasswordIsEmpty         = NewCodeError(200008, "password is empty", "密码不能为空")
	ErrorTwoPasswordNotMatch     = NewCodeError(200009, "two password is not match", "两次密码不一致")
	ErrorPassword                = NewCodeError(200012, "password error", "密码不正确")
	ErrorTokenGenerator          = NewCodeError(200013, "token generator error", "token生成错误")

	//user 信息相关 100 - 199
	ErrorUserIdIsEmpty = NewCodeError(200100, "userID is empty", "用户ID不能为空")
	ErrorUserNotExist  = NewCodeError(200101, "user is not exist", "用户不存在")

	//user 账户安全相关 200-299
	ErrorBindMobileError     = NewCodeError(200200, "bind mobile error", "手机绑定失败")
	ErrorBindEmailError      = NewCodeError(200201, "bind email error", "邮箱绑定失败")
	ErrorAuthCodeIsEmpty     = NewCodeError(200202, "auth code is empty", "验证码不能为空")
	ErrorMobileAuthCodeError = NewCodeError(200203, "mobile auth code error", "手机验证码错误")
	ErrorEmailAuthCodeError  = NewCodeError(200204, "email auth code error", "邮箱验证码错误")
	ErrorEmailAlreadyExist   = NewCodeError(200205, "email is already exist", "邮箱已经被绑定")
	ErrorMobileAlreadyExist  = NewCodeError(200206, "mobile is already exist", "手机已经被绑定")
	ErrorGoogleAlreadyBind   = NewCodeError(200207, "google auth is already bind", "谷歌验证重复绑定")
	ErrorGoogleNotBind       = NewCodeError(200208, "google auth is not bind", "未绑定谷歌验证")
	ErrorGoogleAuthCode      = NewCodeError(200209, "google auth code error", "谷歌验证码错误")

	//upload api
	ErrorFileParse   = NewCodeError(200300, "file parse error", "上传文件解析错误")
	ErrorFileToLarge = NewCodeError(200301, "file too large", "上传文件太大")

	//upload rpc

	//nacos rpc

	//xPert Business 相关
	ErrorBusinessClient = NewCodeError(200400, "business client info error", "业务获取失败")

	// admin offlinetrade相关 000~099
	ErrorOperatorNotFound           = NewCodeError(300000, "Operator does not exists", "用户不存在")
	ErrorDirectionInvalid           = NewCodeError(300001, "Direction is invalid", "交易方向不合法")
	ErrorSystem                     = NewCodeError(300002, "System error", "系统错误")
	ErrorQuantityInvalid            = NewCodeError(300003, "Quantity is invalid", "交易数量不合法")
	ErrorAmountInvalid              = NewCodeError(300004, "Amount is invalid", "交易金额不合法")
	ErrorFeeInvalid                 = NewCodeError(300005, "Fee is invalid", "交易手续费不合法")
	ErrorFeeGtAmount                = NewCodeError(300006, "Fee is greate than amount", "交易手续费不能大于交易金额")
	ErrorAlreadyEnd                 = NewCodeError(300007, "The offline trade is already end state", "交易已结束")
	ErrorOperatorSame               = NewCodeError(300008, "Operator can not be the same", "交易操作员不能相同")
	ErrorExchangeCoinNotExist       = NewCodeError(300009, "Exchange coin is not found", "交易币不存在")
	ErrorBaseCoinNotExist           = NewCodeError(300010, "Base coin is not found", "结算币不存在")
	ErrorExchangeCoinPriceNotExist  = NewCodeError(300011, "Exchange coin can not get the price of market", "不能获取交易币市价")
	ErrorBaseCoinPriceNotExist      = NewCodeError(300012, "Base coin can not not get the price of market", "不能获取结算币市价")
	ErrorOfflineTradeCreateFail     = NewCodeError(300013, "Offline trade create failed", "新建线下交易失败")
	ErrorUidAlreadyUpdate           = NewCodeError(300014, "The id of member can not update", "不能修改客户id")
	ErrorSymbolAlreadyUpdate        = NewCodeError(300015, "The symbol can not update", "不能修改交易币品种")
	ErrorQuantityAlreadyUpdate      = NewCodeError(300016, "The quantity can not update", "不能修改交易数量")
	ErrorAmountAlreadyUpdate        = NewCodeError(300017, "The amount can not update", "不能修改交易金额")
	ErrorFeeAlreadyUpdate           = NewCodeError(300018, "The fee can not update", "不能修改交易费")
	ErrorOfflineTradeUpdateFail     = NewCodeError(300019, "Offline trade update failed", "更新线下交易失败")
	ErrorTradeBusinessNoNotFound    = NewCodeError(300020, "TradeBusinessNo does not exists", "交易号不存在")
	ErrorOfflineTradeNotUpdate      = NewCodeError(300021, "TradeBusinessNo can not be updated", "不能更新交易")
	ErrorOfflineTradeAlreadyExist   = NewCodeError(300022, "The id of offline trade is duplicated", "线下交易号已存在")
	ErrorNotingChange               = NewCodeError(300023, "Nothing is changed", "没有修改,不用更新")
	ErrorOfflineTradeNotFound       = NewCodeError(300024, "The offline trade is not existed", "此交易不存在")
	ErrorTotalBalanceIsEmpty        = NewCodeError(300025, "The totol balance of member is empty", "账户余额为零")
	ErrorMemberNotFound             = NewCodeError(300026, "The member is not found", "账户不存在")
	ErrorOfflineTradeNotDelete      = NewCodeError(300027, "The offline trade can not be deleted", "此交易不能被删除")
	ErrorOfflineTradeNotAllowDelete = NewCodeError(300028, "The offline trade is not allowed to be deleted", "此交易不允许被删除")
	ErrorOfflineTradeIdIsEmpty      = NewCodeError(300029, "The id of offline trade is empty", "交易号不允许为空")
	ErrorRbacUserIdNotFound         = NewCodeError(300030, "The user is not found", "用户不存在")
	ErrorRbacUserOpIdNotFound       = NewCodeError(300031, "The operation of user is not found", "用户权限异常")

	// admin profit相关000~099
	ErrorProfitNotFound = NewCodeError(301000, "The record of profit does not exists", "收益记录不存在")

	// admin symbol相关000~099
	ErrorSymbolNotFound   = NewCodeError(302000, "The symbol does not exists", "交易合约不存在")
	ErrorCurrencyNotFound = NewCodeError(302001, "The currency does not exists", "交易币种不存在")

	// notify 相关000~100
	ErrNotifySendFail          = NewCodeError(400000, "send notify failed", "通知发送失败")
	ErrNotifyRespParseFail     = NewCodeError(400001, "Failed to parse server response", "响应数据解析失败")
	ErrNotifyTooManyRecipients = NewCodeError(400002, "Too many notification recipients", "通知接收者过多")

	//dataService相关
	ErrorDSSymbol   = NewCodeError(500000, "data service get symbol error", "获取合约信息失败")
	ErrorDSCurrency = NewCodeError(500001, "data service get currency error", "获取币种信息失败")

	//order
	ErrorOrderPriceExpired           = NewCodeError(305000, "otc order price is expired", "订单价格已过期")
	ErrorOrderAmountInvalid          = NewCodeError(305001, "Invalid money amount", "下单金额不合法")                                        //原client code:1208
	ErrorOrderVolumeInvalid          = NewCodeError(305002, "The amount do not meet the precision requirement", "下单数量不符合精度要求")        //原client code:1210
	ErrorOrderVolumeMinUnitInvalid   = NewCodeError(305003, "Invalid amount", "下单数量不符合最小交易单位的要求")                                     //原client code:1214
	ErrorOrderAmountPrecisionInvalid = NewCodeError(305004, "The sum of money do not meet the precision requirement", "下单总金额不符合精度要求") //原client code:1212
	ErrorUserAccountNotEnough        = NewCodeError(305005, "member have no enough balance", "用户余额不足")                                //原管理端 code:1236
	ErrorOTCOrderPrice               = NewCodeError(305006, "otc order price error", "otc订单价格错误")
	ErrorQuoteError                  = NewCodeError(305007, "quote error", "询价错误")
)
