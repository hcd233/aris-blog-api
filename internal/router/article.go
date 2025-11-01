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

func initArticleRouter(articleGroup *huma.Group) {
	articleHandler := handler.NewArticleHandler()
	articleGroup.UseMiddleware(middleware.JwtMiddleware())

	huma.Register(articleGroup, huma.Operation{
		OperationID: "listArticles",
		Method:      http.MethodGet,
		Path:        "/list",
		Summary:     "ListArticles",
		Description: "List articles with pagination",
		Tags:        []string{"article"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, articleHandler.HandleListArticles)

	huma.Register(articleGroup, huma.Operation{
		OperationID: "getArticleBySlug",
		Method:      http.MethodGet,
		Path:        "/slug/{authorName}/{articleSlug}",
		Summary:     "GetArticleInfoBySlug",
		Description: "Get article information by author name and slug",
		Tags:        []string{"article"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, articleHandler.HandleGetArticleInfoBySlug)

	huma.Register(articleGroup, huma.Operation{
		OperationID: "getArticleInfo",
		Method:      http.MethodGet,
		Path:        "/{articleID}",
		Summary:     "GetArticleInfo",
		Description: "Get article detail by ID",
		Tags:        []string{"article"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, articleHandler.HandleGetArticleInfo)

	creatorArticleGroup := huma.NewGroup(articleGroup, "")
	creatorArticleGroup.UseMiddleware(middleware.LimitUserPermissionMiddleware("articleService", model.PermissionCreator))

	huma.Register(creatorArticleGroup, huma.Operation{
		OperationID: "createArticle",
		Method:      http.MethodPost,
		Path:        "/",
		Summary:     "CreateArticle",
		Description: "Create a new article",
		Tags:        []string{"article"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, articleHandler.HandleCreateArticle)

	huma.Register(creatorArticleGroup, huma.Operation{
		OperationID: "updateArticle",
		Method:      http.MethodPatch,
		Path:        "/{articleID}",
		Summary:     "UpdateArticle",
		Description: "Update article information",
		Tags:        []string{"article"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, articleHandler.HandleUpdateArticle)

	huma.Register(creatorArticleGroup, huma.Operation{
		OperationID: "deleteArticle",
		Method:      http.MethodDelete,
		Path:        "/{articleID}",
		Summary:     "DeleteArticle",
		Description: "Delete article",
		Tags:        []string{"article"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, articleHandler.HandleDeleteArticle)

	huma.Register(creatorArticleGroup, huma.Operation{
		OperationID: "updateArticleStatus",
		Method:      http.MethodPut,
		Path:        "/{articleID}/status",
		Summary:     "UpdateArticleStatus",
		Description: "Update article publish status",
		Tags:        []string{"article"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, articleHandler.HandleUpdateArticleStatus)

	articleVersionGroup := huma.NewGroup(articleGroup, "/{articleID}/version")
	initArticleVersionRouter(articleVersionGroup)
}

func initArticleVersionRouter(articleVersionGroup *huma.Group) {
	versionHandler := handler.NewArticleVersionHandler()

	huma.Register(articleVersionGroup, huma.Operation{
		OperationID: "getLatestArticleVersion",
		Method:      http.MethodGet,
		Path:        "/latest",
		Summary:     "GetLatestArticleVersion",
		Description: "Get latest article version content",
		Tags:        []string{"articleVersion"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, versionHandler.HandleGetLatestArticleVersionInfo)

	creatorArticleVersionGroup := huma.NewGroup(articleVersionGroup, "")
	creatorArticleVersionGroup.UseMiddleware(middleware.LimitUserPermissionMiddleware("articleVersionService", model.PermissionCreator))

	huma.Register(creatorArticleVersionGroup, huma.Operation{
		OperationID: "listArticleVersions",
		Method:      http.MethodGet,
		Path:        "/list",
		Summary:     "ListArticleVersions",
		Description: "List article versions",
		Tags:        []string{"articleVersion"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, versionHandler.HandleListArticleVersions)

	huma.Register(creatorArticleVersionGroup, huma.Operation{
		OperationID: "createArticleVersion",
		Method:      http.MethodPost,
		Path:        "/",
		Summary:     "CreateArticleVersion",
		Description: "Create a new article version",
		Tags:        []string{"articleVersion"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
		Middlewares: huma.Middlewares{middleware.RateLimiterMiddleware("createArticleVersion", constant.CtxKeyUserID, 10*time.Second, 1)},
	}, versionHandler.HandleCreateArticleVersion)

	huma.Register(creatorArticleVersionGroup, huma.Operation{
		OperationID: "getArticleVersionInfo",
		Method:      http.MethodGet,
		Path:        "/v{version}",
		Summary:     "GetArticleVersionInfo",
		Description: "Get article version info",
		Tags:        []string{"articleVersion"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, versionHandler.HandleGetArticleVersionInfo)
}
