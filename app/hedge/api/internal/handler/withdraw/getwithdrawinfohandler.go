package withdraw

import (
	"net/http"

	"market_server/app/hedge/cmd/api/internal/logic/withdraw"
	"market_server/app/hedge/cmd/api/internal/svc"
	"market_server/app/hedge/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetWithdrawInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WithdrawReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := withdraw.NewGetWithdrawInfoLogic(r.Context(), svcCtx)
		resp, err := l.GetWithdrawInfo(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
