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
					"bearerAuth": {
						Type:         "http",
						Scheme:       "bearer",
						BearerFormat: "JWT",
						Description:  "JWT认证，请在Authorization头中传递Bearer Token。格式: Bearer <token>",
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

	v1Router := app.Group("/v1")
	{
		initCategoryRouter(v1Router)
		initTagRouter(v1Router)
		initArticleRouter(v1Router)
		initCommentRouter(v1Router)

		initAssetRouter(v1Router)
		initOperationRouter(v1Router)

		initAIRouter(v1Router)
	}

	v1Group := huma.NewGroup(api, "/v1")
	userGroup := huma.NewGroup(v1Group, "/user")
	initUserRouter(userGroup)

	tokenGroup := huma.NewGroup(v1Group, "/token")
	initTokenRouter(tokenGroup)

	oauth2Group := huma.NewGroup(v1Group, "/oauth2")
	initOauth2Router(oauth2Group)

	huma.Register(api, huma.Operation{
		OperationID: "ping",
		Method:      http.MethodGet,
		Path:        "/",
		Summary:     "Ping",
		Description: "Check service if available.",
		Tags:        []string{"ping"},
	}, pingService.HandlePing)
}
