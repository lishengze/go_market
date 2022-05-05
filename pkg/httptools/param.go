package httptools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

const (
	tagKeyForm  = "form"
	tagKeyParam = "param"
	tagKeyJson  = "json"
)

const (
	tagOptOmitempty = "omitempty"
	tagOptRequired  = "required"
	tagOptDefault   = "default"
	tagOptOptions   = "options"
)

var EmptyReq = struct{}{}

/*
示例
type Req struct {
	Name1  string `form:"Name,optional"`
	Name2 string `json:"Name2,omitempty"`
	Name3 string `json:"Name3,default=3"`
	Name4 string `param:"Name4,default=3"`
	Name5 string `param:"Name4,options=5|6"`
}
*/

type (
	// IntegralParam 包含一次http 请求所有 参数
	IntegralParam struct {
		HttpMethod string
		Url        *url.URL
		Header     http.Header
		Param      url.Values // 路径参数
		Form       url.Values
		JsonBody
	}

	tagValue struct {
		field     string
		omitempty bool
		required  bool
		default_  string
		options   map[string]interface{}
	}

	JsonBody map[string]interface{}
)

func NewIntegralParam() *IntegralParam {
	return &IntegralParam{
		Url:      &url.URL{},
		Header:   http.Header{},
		Param:    url.Values{},
		Form:     url.Values{},
		JsonBody: JsonBody{},
	}
}

// TrimmedString 如果是空 去除 {}
func (o JsonBody) TrimmedString() (string, error) {
	if len(o) == 0 {
		return "", nil
	}

	bodyBytes, err := json.Marshal(o)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), nil
}

func (o *tagValue) verifyValue(v string) error {
	if len(o.options) > 0 {
		_, ok := o.options[v]
		if !ok {
			fmt.Println(v, "-----------")
			return fmt.Errorf("field:%s vaule:%s not in options:%v", o.field, v, o.options)
		}
	}
	return nil
}

func ParseReqParam(req interface{}) (*IntegralParam, error) {
	ip := NewIntegralParam()
	if req == EmptyReq {
		return ip, nil
	}
	t := reflect.TypeOf(req)
	v := reflect.ValueOf(req)
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("not support req type:%s", t.Kind())
	}

	for i := 0; i < t.NumField(); i++ {
		key := t.Field(i)
		value := v.Field(i)
		valueStr := fmt.Sprint(value)
		if !value.IsValid() {
			return nil, fmt.Errorf("the value of field:%s is invalid", key.Name)
		}

		if v, ok := key.Tag.Lookup(tagKeyForm); ok {
			tv, err := parseTagValue(v)
			if err != nil {
				return nil, err
			}

			if value.IsZero() {
				if tv.required {
					ip.Form.Set(tv.field, valueStr)
					continue
				}
				if tv.omitempty { // ignore
					continue
				}
				ip.Form.Set(tv.field, tv.default_)
			} else {
				if err := tv.verifyValue(valueStr); err != nil {
					return nil, err
				}
				ip.Form.Set(tv.field, valueStr)
			}

		} else if v, ok = key.Tag.Lookup(tagKeyParam); ok {
			tv, err := parseTagValue(v)
			if err != nil {
				return nil, err
			}

			if value.IsZero() {
				if tv.required {
					ip.Param.Set(tv.field, valueStr)
					continue
				}

				if tv.omitempty { // ignore
					continue
				}
				ip.Param.Set(tv.field, tv.default_)
			} else {
				if err := tv.verifyValue(valueStr); err != nil {
					return nil, err
				}
				ip.Param.Set(tv.field, valueStr)
			}
		} else if v, ok = key.Tag.Lookup(tagKeyJson); ok {
			tv, err := parseTagValue(v)
			if err != nil {
				return nil, err
			}

			if value.IsZero() {
				if tv.omitempty { // ignore
					continue
				}
				ip.JsonBody[tv.field] = tv.default_
			} else {
				if err := tv.verifyValue(valueStr); err != nil {
					return nil, err
				}
				ip.JsonBody[tv.field] = value.Interface()
			}
		}
	}
	return ip, nil
}

func parseTagValue(v string) (*tagValue, error) {
	tv := &tagValue{
		field:     "",
		omitempty: false,
		default_:  "",
		options:   map[string]interface{}{},
	}
	l := strings.Split(v, ",")
	if len(l) == 0 {
		return nil, fmt.Errorf("Invalid tag:%s ", v)
	}
	for i, item := range l {
		if i == 0 {
			tv.field = item
			continue
		}
		if strings.Contains(item, tagOptDefault) {
			l2 := strings.Split(item, "=")
			if len(l2) != 2 {
				return nil, fmt.Errorf("Invalid default tag:%s ", item)
			}
			tv.default_ = l2[1]
			continue
		}

		if strings.Contains(item, tagOptOptions) {
			l2 := strings.Split(item, "=")
			if len(l2) != 2 {
				return nil, fmt.Errorf("Invalid options tag:%s ", item)
			}
			l3 := strings.Split(l2[1], "|")
			if len(l3) == 0 {
				return nil, fmt.Errorf("Invalid options tag:%s ", item)
			}
			for _, opt := range l3 {
				tv.options[opt] = struct{}{}
			}
			continue
		}

		if strings.Contains(item, tagOptOmitempty) {
			tv.omitempty = true
			continue
		}

		if strings.Contains(item, tagOptRequired) {
			tv.required = true
			continue
		}

	}

	return tv, nil
}
