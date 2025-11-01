// Package router 路由
package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-blog-api/internal/api"
	"github.com/hcd233/aris-blog-api/internal/handler"
)

// RegisterRouter 注册路由
//
//	param app *fiber.App
//	author centonhuang
//	update 2025-01-04 15:32:40
func RegisterRouter() {
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
