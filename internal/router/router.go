// Package router 路由
package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-blog-api/internal/api"
	"github.com/hcd233/aris-blog-api/internal/handler"
)

// RegisterDocsRouter 注册文档路由
//
//	@return *fiber.App
//	@author centonhuang
//	@update 2025-11-10 18:29:32
func RegisterDocsRouter() {
	app := api.GetFiberApp()
	app.Get("/docs", func(c *fiber.Ctx) error {
		html := `<!doctype html>
<html>
  <head>
    <title>API Reference</title>
    <meta charset="utf-8" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <script
      id="api-reference"
      data-url="/openapi.json"></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>`
		return c.Type("html").SendString(html)
	})
}

// RegisterAPIRouter 注册API路由
//
//	param app *fiber.App
//	author centonhuang
//	update 2025-01-04 15:32:40
func RegisterAPIRouter() {
	pingService := handler.NewPingHandler()

	api := api.GetHumaAPI()

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
