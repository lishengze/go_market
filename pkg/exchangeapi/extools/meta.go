package extools

const (
	ReqTypeWait  ReqType = iota + 1 // 遇到频率限制触发等待
	ReqTypeAllow                    // 遇到频率限制,请求被丢弃,触发error
	//ReqTypeWs                       // websocket
)

type (
	Meta interface {
		Url() string
		HttpMethod() string
		Weight() int
		ReqType() ReqType
		NeedSign() bool
	}

	// ReqType 请求类型
	ReqType int

	// meta 每个接口/方法 的元信息
	meta struct {
		httpMethod string
		url        string
		weight     int
		reqType    ReqType
		needSign   bool
	}
)

func (o ReqType) String() string {
	switch o {
	case ReqTypeAllow:
		return "ALLOW"
	case ReqTypeWait:
		return "WAIT"
	//case ReqTypeWs:
	//	return "Websocket"
	default:
		return "UNKNOWN"
	}
}

func NewMetaWithGetWeightFn(httpMethod, url_ string, ReqType ReqType, needSign bool, getWeightFn func() int) Meta {
	meta := &meta{
		httpMethod: httpMethod,
		url:        url_,
		weight:     getWeightFn(),
		reqType:    ReqType,
		needSign:   needSign,
	}
	return meta
}

func NewMeta(httpMethod, url_ string, weight int, ReqType ReqType, needSign bool) Meta {
	meta := &meta{
		httpMethod: httpMethod,
		url:        url_,
		weight:     weight,
		reqType:    ReqType,
		needSign:   needSign,
	}
	return meta
}

func NewMetaWithOneWeight(httpMethod, url_ string, ReqType ReqType, needSign bool) Meta {
	meta := &meta{
		httpMethod: httpMethod,
		url:        url_,
		weight:     1, // 默认1
		reqType:    ReqType,
		needSign:   needSign,
	}
	return meta
}

func (o *meta) Url() string {
	return o.url
}

func (o *meta) HttpMethod() string {
	return o.httpMethod
}

func (o *meta) ReqType() ReqType {
	return o.reqType
}

func (o *meta) Weight() int {
	return o.weight
}

func (o *meta) NeedSign() bool {
	return o.needSign
}
