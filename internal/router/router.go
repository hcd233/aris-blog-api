// Package router 路由
package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/hcd233/aris-blog-api/internal/handler"
)

// RegisterRouter 注册路由
//
//	param app *fiber.App
//	author centonhuang
//	update 2025-01-04 15:32:40
func RegisterRouter(app *fiber.App) {
	pingService := handler.NewPingHandler()

	rootRouter := app.Group("")

	v1Router := app.Group("/v1")
	{
		initTokenRouter(v1Router)
		initOauth2Router(v1Router)

		initUserRouter(v1Router)

		initCategoryRouter(v1Router)
		initTagRouter(v1Router)
		initArticleRouter(v1Router)
		initCommentRouter(v1Router)

		initAssetRouter(v1Router)
		initOperationRouter(v1Router)

		initAIRouter(v1Router)
	}

	api := humafiber.NewWithGroup(app, rootRouter, huma.DefaultConfig("Aris-blog", "1.0"))

	huma.Register(api, huma.Operation{
		OperationID: "ping",
		Method:      http.MethodGet,
		Path:        "/",
		Summary:     "Ping Pong!",
		Description: "Check service if available.",
		Tags:        []string{"ping"},
	}, pingService.HandlePing)
}
