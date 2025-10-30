// Package router 路由
package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
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

	api := humafiber.NewWithGroup(app, rootRouter, huma.Config{
		OpenAPI: &huma.OpenAPI{
			OpenAPI: "3.1.0",
			Info: &huma.Info{
				Title:   "Aris-blog",
				Version: "1.0",
			},
			Components: &huma.Components{
				Schemas: huma.NewMapRegistry("#/components/schemas/", huma.DefaultSchemaNamer),
				SecuritySchemes: map[string]*huma.SecurityScheme{
					"jwtAuth": {
						Type:        "apiKey",
						Name:        "Authorization",
						In:          "header",
						Description: "JWT Authentication，Please pass the JWT token in the Authorization header.",
					},
				},
			},
		},
		OpenAPIPath:   "/openapi",
		DocsPath:      "/docs",
		SchemasPath:   "/schemas",
		Formats:       huma.DefaultFormats,
		DefaultFormat: "application/json",
	})

	// 保留需要特殊处理的服务（SSE流式响应、文件上传等）
	v1Router := app.Group("/v1")
	{
		initAssetRouter(v1Router)  // 文件上传需要特殊处理
		initAIRouter(v1Router)     // SSE流式响应需要特殊处理
	}

	// Huma路由组 - 已重构的服务
	v1Group := huma.NewGroup(api, "/v1")
	
	// 用户相关
	userGroup := huma.NewGroup(v1Group, "/user")
	initUserRouter(userGroup)

	// 令牌相关
	tokenGroup := huma.NewGroup(v1Group, "/token")
	initTokenRouter(tokenGroup)

	// OAuth2相关
	oauth2Group := huma.NewGroup(v1Group, "/oauth2")
	initOauth2Router(oauth2Group)

	// 标签相关
	tagGroup := huma.NewGroup(v1Group, "/tag")
	initTagRouter(tagGroup)

	// 分类相关
	categoryGroup := huma.NewGroup(v1Group, "/category")
	initCategoryRouter(categoryGroup)

	// 文章相关
	articleGroup := huma.NewGroup(v1Group, "/article")
	initArticleRouter(articleGroup)
	initArticleVersionRouter(articleGroup)  // 文章版本路由注册在文章组下

	// 评论相关
	commentGroup := huma.NewGroup(v1Group, "/comment")
	initCommentRouter(commentGroup)

	// 操作相关
	operationGroup := huma.NewGroup(v1Group, "/operation")
	initOperationRouter(operationGroup)

	huma.Register(api, huma.Operation{
		OperationID: "ping",
		Method:      http.MethodGet,
		Path:        "/",
		Summary:     "Ping",
		Description: "Check service if available.",
		Tags:        []string{"ping"},
	}, pingService.HandlePing)
}
