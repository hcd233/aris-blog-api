// Package router 路由
package router

import (
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
	// 注册 Huma 路由 (新的 OpenAPI 集成)
	RegisterHumaRouter(app)

	// 保留原有的 Swagger 路由用于向后兼容
	app.Get("/swagger/*", swagger.HandlerDefault)

	// 原有的健康检查路由
	pingService := handler.NewPingHandler()
	app.Get("/ping", pingService.HandlePing) // 改为 /ping 避免与 Huma 的根路径冲突

	// 原有的 API 路由 - 逐步迁移到 Huma
	v1Router := app.Group("/v1/legacy") // 添加 legacy 前缀来区分
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
}
