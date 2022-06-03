// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	hedge "bcts/app/hedge/cmd/api/internal/handler/hedge"
	member "bcts/app/hedge/cmd/api/internal/handler/member"
	withdraw "bcts/app/hedge/cmd/api/internal/handler/withdraw"
	"bcts/app/hedge/cmd/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/risk/monitor",
				Handler: hedge.GetRiskMonitorHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.JwtAuth.AccessSecret),
		rest.WithPrefix("/api/hedge"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/member/:id",
				Handler: member.GetMemberInfoHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.JwtAuth.AccessSecret),
		rest.WithPrefix("/api/admin"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/withdraw/:id",
				Handler: withdraw.GetWithdrawInfoHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.JwtAuth.AccessSecret),
		rest.WithPrefix("/api/withdraw"),
	)
}
