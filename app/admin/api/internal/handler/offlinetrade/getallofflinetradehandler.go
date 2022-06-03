package offlinetrade

import (
	"market_server/app/admin/api/internal/logic/offlinetrade"
	"market_server/app/admin/api/internal/svc"
	"market_server/app/admin/api/internal/types"
	"net/http"

	"market_server/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetAllOfflineTradeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetAllOfflineTradeReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := offlinetrade.NewGetAllOfflineTradeLogic(r.Context(), svcCtx)
		resp, err := l.GetAllOfflineTrade(&req)
		result.HttpResult(r, w, resp, err)
	}
}
