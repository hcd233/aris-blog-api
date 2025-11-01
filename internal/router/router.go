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

	v1Group := huma.NewGroup(api, "/v1")
	userGroup := huma.NewGroup(v1Group, "/user")
	initUserRouter(userGroup)

	tagGroup := huma.NewGroup(v1Group, "/tag")
	initTagRouter(tagGroup)

	articleGroup := huma.NewGroup(v1Group, "/article")
	initArticleRouter(articleGroup)

	commentGroup := huma.NewGroup(v1Group, "/comment")
	initCommentRouter(commentGroup)

	tokenGroup := huma.NewGroup(v1Group, "/token")
	initTokenRouter(tokenGroup)

	oauth2Group := huma.NewGroup(v1Group, "/oauth2")
	initOauth2Router(oauth2Group)

	operationGroup := huma.NewGroup(v1Group, "/operation")
	initOperationRouter(operationGroup)

	assetGroup := huma.NewGroup(v1Group, "/asset")
	initAssetRouter(assetGroup)

	aiGroup := huma.NewGroup(v1Group, "/ai")
	initAIRouter(aiGroup)

	huma.Register(api, huma.Operation{
		OperationID: "ping",
		Method:      http.MethodGet,
		Path:        "/",
		Summary:     "Ping",
		Description: "Check service if available.",
		Tags:        []string{"ping"},
	}, pingService.HandlePing)
}
