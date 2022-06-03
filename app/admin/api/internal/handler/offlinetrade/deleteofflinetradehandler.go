package offlinetrade

import (
	"bcts/app/admin/api/internal/logic/offlinetrade"
	"bcts/app/admin/api/internal/svc"
	"bcts/app/admin/api/internal/types"
	"net/http"

	"bcts/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func DeleteOfflineTradeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteOfflineTradeReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := offlinetrade.NewDeleteOfflineTradeLogic(r.Context(), svcCtx)
		resp, err := l.DeleteOfflineTrade(&req)
		result.HttpResult(r, w, resp, err)
	}
}
