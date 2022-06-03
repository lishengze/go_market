package {{.PkgName}}

import (
	"market_server/common/crypto"
    "market_server/common/result"
    "encoding/json"
	"net/http"
	"encoding/json"
    "io/ioutil"

	{{if .After1_1_10}}"github.com/zeromicro/go-zero/rest/httpx"{{end}}
	{{if .After1_1_10}}"github.com/zeromicro/go-zero/core/logx"{{end}}
	{{.ImportPackages}}
)


func {{.HandlerName}}(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
		logx.Infof("header:%+v", r.Header)
			if r.Method == "POST" {
        			iv := r.Header.Get("X-Encrypt-Iv")
        			//1. 可以认为是直接调用
        			if len(iv) == 0 {
        				{{if .HasRequest}}var req types.{{.RequestType}}
        				if err := httpx.Parse(r, &req); err != nil {
        					result.ParamErrorResult(r, w, err, svcCtx.PemFileBase.AESKey, false)
        					return
        				}

        				{{end}}l := {{.LogicName}}.New{{.LogicType}}(r, r.Context(), svcCtx)
		                {{if .HasResp}}resp, {{end}}err := l.{{.Call}}({{if .HasRequest}}&req{{end}})
		                result.HttpResult(r, w, {{if .HasResp}}resp{{else}}nil{{end}}, err, svcCtx.PemFileBase.AESKey, false)
        				return
        			}

                    //2. 需要做加解密的操作
        			requestBody, _ := ioutil.ReadAll(r.Body)

        			var encrypted crypto.EncryptedReqBody
        			if err := json.Unmarshal(requestBody, &encrypted); err != nil {
        				result.HttpResult(r, w, nil, err, svcCtx.PemFileBase.AESKey, svcCtx.Config.ResponseEncrypted)
        				return
        			}
        			logx.Infof("after unmarshal: %+v", encrypted)

        			requestPlainMsg, err := crypto.DecryptReqBody(&encrypted, svcCtx.PemFileBase.AESKey, iv)
        			if err != nil {
        				logx.Errorf("DecryptReqBody err %+v", err)
        				result.HttpResult(r, w, nil, err, svcCtx.PemFileBase.AESKey, svcCtx.Config.ResponseEncrypted)
        				return
        			}
        			logx.Info(string(requestPlainMsg))

        			{{if .HasRequest}}var req types.{{.RequestType}}
        			if err := json.Unmarshal(requestPlainMsg, &req); err != nil {
        				logx.Errorf("{{.LogicType}} Unmarshal err:%+v", err)
        				result.ParamErrorResult(r, w, err, svcCtx.PemFileBase.AESKey, svcCtx.Config.ResponseEncrypted)
        				return
        			}

        			logx.Infof("{{.LogicType}} req:%+v", req)

        			{{end}}l := {{.LogicName}}.New{{.LogicType}}(r, r.Context(), svcCtx)
        		    resp, err := l.{{.Call}}({{if .HasRequest}}&req{{end}})
        			result.HttpResult(r, w, resp, err, svcCtx.PemFileBase.AESKey, svcCtx.Config.ResponseEncrypted)
        		}

        		if r.Method == "GET" {
        			iv := r.Header.Get("X-Encrypt-Iv")
        			logx.Infof("{{.LogicType}} iv:%s", iv)

        			{{if .HasRequest}}var req types.{{.RequestType}}
        			if err := httpx.Parse(r, &req); err != nil {
        				result.ParamErrorResult(r, w, err, svcCtx.PemFileBase.AESKey, svcCtx.Config.ResponseEncrypted)
        				return
        			}

        			{{end}}l := {{.LogicName}}.New{{.LogicType}}(r, r.Context(), svcCtx)
        			{{if .HasResp}}resp, {{end}}err := l.{{.Call}}({{if .HasRequest}}&req{{end}})
                    result.HttpResult(r, w, resp, err, svcCtx.PemFileBase.AESKey, svcCtx.Config.ResponseEncrypted)
        		}
	}
}