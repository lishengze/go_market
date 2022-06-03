package main

import (
	"flag"
	"fmt"
	"market_server/app/admin/api/internal/config"
	"market_server/app/admin/api/internal/handler"
	"market_server/app/admin/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

//go:generate goctl model mysql datasource --url "broker:QsTZfnXb%yJw&r@7@tcp(36.255.220.139:3306)/broker" --table admin_user --dir /Users/zhq/workspace/go/gitlab.hub.hashkey.com/bcts/market_server/app/admin/model
//go:generate goctl api go --api app/admin/cmd/api/desc/admin.api --dir app/admin/cmd/api/
var configFile = flag.String("f", "etc/admin.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
