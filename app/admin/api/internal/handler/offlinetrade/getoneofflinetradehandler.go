package offlinetrade

import (
	"bcts/app/admin/api/internal/logic/offlinetrade"
	"bcts/app/admin/api/internal/svc"
	"bcts/app/admin/api/internal/types"
	"net/http"

	"bcts/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetOneOfflineTradeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OfflineTradeReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := offlinetrade.NewGetOneOfflineTradeLogic(r.Context(), svcCtx)
		resp, err := l.GetOneOfflineTrade(&req)
		result.HttpResult(r, w, resp, err)
	}
}
