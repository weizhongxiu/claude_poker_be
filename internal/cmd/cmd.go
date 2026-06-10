package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"claude-test/internal/controller/club"
	gameCtrl "claude-test/internal/controller/game"
	"claude-test/internal/controller/stats"
	"claude-test/internal/controller/table"
	"claude-test/internal/controller/user"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()

			// REST API group
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(
					user.NewV1(),
					table.NewV1(),
					stats.NewV1(),
					club.NewV1(),
				)
			})

			// WebSocket — separate group, no JSON response middleware
			wsController := gameCtrl.NewV1()
			s.Group("/ws", func(group *ghttp.RouterGroup) {
				group.GET("/table/{table_id}", wsController.WsTable)
			})

			s.Run()
			return nil
		},
	}
)
