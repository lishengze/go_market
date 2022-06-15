package {{.PkgName}}

import (
	"net/http"

    "common/httpresult"

	{{if .After1_1_10}}"github.com/zeromicro/go-zero/rest/httpx"{{end}}
	{{.ImportPackages}}
)


func {{.HandlerName}}(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		{{- if .HasRequest}}var req types.{{.RequestType}}
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}{{- end}}

		l := {{.LogicName}}.New{{.LogicType}}(r.Context(), ctx)
		{{- if .HasResp}}
		resp, err := l.{{.Call}}({{if .HasRequest}}req{{end}})
		if err!=nil{
		    {{- if .HasRequest}}l.Logger.Errorf("handler:%s, req:%+v, err:%v", "{{.HandlerName}}" ,req ,err){{- end}}
		    httpresult.Error(r, w, err)
		}else{
			httpx.OkJson(w, resp)
		}
		{{- else}}
		err := l.{{.Call}}({{if .HasRequest}}req{{end}})
		if err!=nil{
		    {{- if .HasRequest}}l.Logger.Errorf("handler:%s, req:%+v, err:%v", "{{.HandlerName}}" ,req ,err){{- end}}
		    l.Logger.Errorf("handler:%s, req:")
            httpresult.Error(r, w, err)
        }else{
        	httpx.Ok(w)
        }
        {{- end}}
	}
}
