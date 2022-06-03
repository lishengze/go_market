package xerror

var message = make(map[int]*CodeError, 0)

//message := make(map[int]*CodeError)

type (
	CodeError struct {
		Code int    `json:"code"`
		EN   string `json:"en"`
		CH   string `json:"ch"`
	}
)

func NewCodeError(code int, en, ch string) *CodeError {
	cr := &CodeError{Code: code, EN: en, CH: ch}
	message[code] = cr
	return cr
}

func GetCodeError(code int) *CodeError {
	if cr, ok := message[code]; ok {
		return cr
	} else {
		return nil
	}
}

func (e *CodeError) ErrCode() int {
	return e.Code
}

func (e *CodeError) Error() string {
	return e.EN
}

func (e *CodeError) ErrorMsg() string {
	return e.EN
}

func (e *CodeError) ErrorMsgCHS() string {
	return e.CH
}

func (e *CodeError) WithMessage(msg, msgChs string) *CodeError {
	if len(msg) > 0 {
		e.EN = msg
	}

	if len(msgChs) > 0 {
		e.CH = msgChs
	}
	return e
}

//func (e *CodeError) Data() *reponse.ResponseBean {
//	return &reponse.ResponseBean{
//		Code:   e.Code,
//		Msg:    e.Msg,
//		MsgCHS: e.MsgCHS,
//	}
//}
