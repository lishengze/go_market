package hedge

import (
	"bcts/app/admin/api/internal/logic/hedge"
	"bcts/app/admin/api/internal/svc"
	"bcts/app/admin/api/internal/types"
	"net/http"

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
