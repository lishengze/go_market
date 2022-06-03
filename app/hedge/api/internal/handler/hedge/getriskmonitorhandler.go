package hedge

import (
	"net/http"

	"bcts/app/hedge/cmd/api/internal/logic/hedge"
	"bcts/app/hedge/cmd/api/internal/svc"
	"bcts/app/hedge/cmd/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetRiskMonitorHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ReqQueryHedgeRisk
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := hedge.NewGetRiskMonitorLogic(r.Context(), svcCtx)
		resp, err := l.GetRiskMonitor(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
