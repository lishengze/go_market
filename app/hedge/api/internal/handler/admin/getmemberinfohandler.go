package admin

import (
	"net/http"

	"market_server/app/hedge/cmd/api/internal/logic/admin"
	"market_server/app/hedge/cmd/api/internal/svc"
	"market_server/app/hedge/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetMemberInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MemberReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := admin.NewGetMemberInfoLogic(r.Context(), svcCtx)
		resp, err := l.GetMemberInfo(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
