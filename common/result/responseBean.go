package result

type ResponseSuccessBean struct {
	Code int         `json:"code"`
	EN   string      `json:"en"`
	CH   string      `json:"ch"`
	Data interface{} `json:"data"`
}
type NullJson struct{}

func Success(data interface{}) *ResponseSuccessBean {
	return &ResponseSuccessBean{0, "success", "成功", data}
}

type ResponseErrorBean struct {
	Code int    `json:"code"`
	EN   string `json:"en"`
	CH   string `json:"ch"`
}

func Error(errCode int, errMsg, errMsgCHS string) *ResponseErrorBean {
	return &ResponseErrorBean{errCode, errMsg, errMsgCHS}
}
