package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
)

func initArticleRouter(articleGroup *huma.Group) {
	articleHandler := handler.NewArticleHandler()

	articleGroup.UseMiddleware(middleware.JwtMiddlewareForHuma())

	// 列出文章
	huma.Register(articleGroup, huma.Operation{
		OperationID: "listArticles",
		Method:      http.MethodGet,
		Path:        "/list",
		Summary:     "ListArticles",
		Description: "Get a paginated list of articles",
		Tags:        []string{"article"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, articleHandler.HandleListArticles)

	// 创建文章
	huma.Register(articleGroup, huma.Operation{
		OperationID: "createArticle",
		Method:      http.MethodPost,
		Path:        "/",
		Summary:     "CreateArticle",
		Description: "Create a new article",
		Tags:        []string{"article"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, articleHandler.HandleCreateArticle)

	// 通过Slug获取文章信息
	huma.Register(articleGroup, huma.Operation{
		OperationID: "getArticleInfoBySlug",
		Method:      http.MethodGet,
		Path:        "/slug/{authorName}/{articleSlug}",
		Summary:     "GetArticleInfoBySlug",
		Description: "Get article information by author name and article slug",
		Tags:        []string{"article"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, articleHandler.HandleGetArticleInfoBySlug)

	// 获取文章信息
	huma.Register(articleGroup, huma.Operation{
		OperationID: "getArticleInfo",
		Method:      http.MethodGet,
		Path:        "/{articleID}",
		Summary:     "GetArticleInfo",
		Description: "Get detailed information about a specific article by ID",
		Tags:        []string{"article"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, articleHandler.HandleGetArticleInfo)

	// 更新文章
	huma.Register(articleGroup, huma.Operation{
		OperationID: "updateArticle",
		Method:      http.MethodPatch,
		Path:        "/{articleID}",
		Summary:     "UpdateArticle",
		Description: "Update an existing article",
		Tags:        []string{"article"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, articleHandler.HandleUpdateArticle)

	// 删除文章
	huma.Register(articleGroup, huma.Operation{
		OperationID: "deleteArticle",
		Method:      http.MethodDelete,
		Path:        "/{articleID}",
		Summary:     "DeleteArticle",
		Description: "Delete an article",
		Tags:        []string{"article"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, articleHandler.HandleDeleteArticle)

	// 更新文章状态
	huma.Register(articleGroup, huma.Operation{
		OperationID: "updateArticleStatus",
		Method:      http.MethodPut,
		Path:        "/{articleID}/status",
		Summary:     "UpdateArticleStatus",
		Description: "Update article status (draft/publish)",
		Tags:        []string{"article"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, articleHandler.HandleUpdateArticleStatus)
}
