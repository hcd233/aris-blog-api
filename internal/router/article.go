package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
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

	// 通过slug获取文章信息
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

	// 创建文章
	createArticleGroup := huma.NewGroup(articleGroup, "")
	createArticleGroup.UseMiddleware(middleware.PermissionMiddlewareForHuma("articleService", model.PermissionCreator))
	huma.Register(createArticleGroup, huma.Operation{
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

	// 文章ID相关的路由
	articleIDGroup := huma.NewGroup(articleGroup, "/{articleID}")
	{
		// 获取文章信息
		huma.Register(articleIDGroup, huma.Operation{
			OperationID: "getArticleInfo",
			Method:      http.MethodGet,
			Path:        "/",
			Summary:     "GetArticleInfo",
			Description: "Get article information by article ID",
			Tags:        []string{"article"},
			Security: []map[string][]string{
				{"jwtAuth": {}},
			},
		}, articleHandler.HandleGetArticleInfo)

		// 更新文章
		updateArticleGroup := huma.NewGroup(articleIDGroup, "")
		updateArticleGroup.UseMiddleware(middleware.PermissionMiddlewareForHuma("articleService", model.PermissionCreator))
		huma.Register(updateArticleGroup, huma.Operation{
			OperationID: "updateArticle",
			Method:      http.MethodPatch,
			Path:        "/",
			Summary:     "UpdateArticle",
			Description: "Update article information",
			Tags:        []string{"article"},
			Security: []map[string][]string{
				{"jwtAuth": {}},
			},
		}, articleHandler.HandleUpdateArticle)

		// 删除文章
		deleteArticleGroup := huma.NewGroup(articleIDGroup, "")
		deleteArticleGroup.UseMiddleware(middleware.PermissionMiddlewareForHuma("articleService", model.PermissionCreator))
		huma.Register(deleteArticleGroup, huma.Operation{
			OperationID: "deleteArticle",
			Method:      http.MethodDelete,
			Path:        "/",
			Summary:     "DeleteArticle",
			Description: "Delete an article",
			Tags:        []string{"article"},
			Security: []map[string][]string{
				{"jwtAuth": {}},
			},
		}, articleHandler.HandleDeleteArticle)

		// 更新文章状态
		updateStatusGroup := huma.NewGroup(articleIDGroup, "")
		updateStatusGroup.UseMiddleware(middleware.PermissionMiddlewareForHuma("articleService", model.PermissionCreator))
		huma.Register(updateStatusGroup, huma.Operation{
			OperationID: "updateArticleStatus",
			Method:      http.MethodPut,
			Path:        "/status",
			Summary:     "UpdateArticleStatus",
			Description: "Update article publication status",
			Tags:        []string{"article"},
			Security: []map[string][]string{
				{"jwtAuth": {}},
			},
		}, articleHandler.HandleUpdateArticleStatus)

		// 文章版本路由
		initArticleVersionRouter(articleIDGroup)
	}
}
