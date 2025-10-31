package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
)

func initArticleVersionRouter(articleGroup *huma.Group) {
	articleVersionHandler := handler.NewArticleVersionHandler()

	articleGroup.UseMiddleware(middleware.JwtMiddlewareForHuma())

	// 创建文章版本
	huma.Register(articleGroup, huma.Operation{
		OperationID: "createArticleVersion",
		Method:      http.MethodPost,
		Path:        "/{articleID}/version",
		Summary:     "CreateArticleVersion",
		Description: "Create a new version of an article",
		Tags:        []string{"articleVersion"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, articleVersionHandler.HandleCreateArticleVersion)

	// 列出文章版本
	huma.Register(articleGroup, huma.Operation{
		OperationID: "listArticleVersions",
		Method:      http.MethodGet,
		Path:        "/{articleID}/version/list",
		Summary:     "ListArticleVersions",
		Description: "Get a paginated list of article versions",
		Tags:        []string{"articleVersion"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, articleVersionHandler.HandleListArticleVersions)

	// 获取最新文章版本
	huma.Register(articleGroup, huma.Operation{
		OperationID: "getLatestArticleVersion",
		Method:      http.MethodGet,
		Path:        "/{articleID}/version/latest",
		Summary:     "GetLatestArticleVersion",
		Description: "Get the latest version of an article",
		Tags:        []string{"articleVersion"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, articleVersionHandler.HandleGetLatestArticleVersionInfo)

	// 获取指定版本
	huma.Register(articleGroup, huma.Operation{
		OperationID: "getArticleVersion",
		Method:      http.MethodGet,
		Path:        "/{articleID}/version/v{version}",
		Summary:     "GetArticleVersion",
		Description: "Get a specific version of an article",
		Tags:        []string{"articleVersion"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, articleVersionHandler.HandleGetArticleVersionInfo)
}
