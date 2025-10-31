package router

import (
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
)

func initArticleVersionRouter(articleIDGroup *huma.Group) {
	articleVersionHandler := handler.NewArticleVersionHandler()

	// 获取最新文章版本信息（不需要权限限制）
	huma.Register(articleIDGroup, huma.Operation{
		OperationID: "getLatestArticleVersionInfo",
		Method:      http.MethodGet,
		Path:        "/version/latest",
		Summary:     "GetLatestArticleVersionInfo",
		Description: "Get the latest version information of an article",
		Tags:        []string{"articleVersion"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, articleVersionHandler.HandleGetLatestArticleVersionInfo)

	// 文章版本路由组（需要 Creator 权限）
	versionGroup := huma.NewGroup(articleIDGroup, "/version")
	versionGroup.UseMiddleware(middleware.PermissionMiddlewareForHuma("articleVersionService", model.PermissionCreator))

	// 列出文章版本
	huma.Register(versionGroup, huma.Operation{
		OperationID: "listArticleVersions",
		Method:      http.MethodGet,
		Path:        "/list",
		Summary:     "ListArticleVersions",
		Description: "Get a paginated list of article versions",
		Tags:        []string{"articleVersion"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, articleVersionHandler.HandleListArticleVersions)

	// 创建文章版本（需要限流）
	createVersionGroup := huma.NewGroup(versionGroup, "")
	createVersionGroup.UseMiddleware(middleware.RateLimiterMiddlewareForHuma("createArticleVersion", constant.CtxKeyUserID, 10*time.Second, 1))
	huma.Register(createVersionGroup, huma.Operation{
		OperationID: "createArticleVersion",
		Method:      http.MethodPost,
		Path:        "/",
		Summary:     "CreateArticleVersion",
		Description: "Create a new version for an article",
		Tags:        []string{"articleVersion"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, articleVersionHandler.HandleCreateArticleVersion)

	// 获取指定版本的文章版本信息
	huma.Register(versionGroup, huma.Operation{
		OperationID: "getArticleVersionInfo",
		Method:      http.MethodGet,
		Path:        "/v{version}",
		Summary:     "GetArticleVersionInfo",
		Description: "Get article version information by version number",
		Tags:        []string{"articleVersion"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, articleVersionHandler.HandleGetArticleVersionInfo)
}
