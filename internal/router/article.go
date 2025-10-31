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

func initArticleRouter(v1Group *huma.Group) {
	articleHandler := handler.NewArticleHandler()
	articleGroup := huma.NewGroup(v1Group, "/article")
	articleGroup.UseMiddleware(middleware.JwtMiddlewareForHuma())

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

	securedArticle := huma.NewGroup(articleGroup, "")
	securedArticle.UseMiddleware(middleware.LimitUserPermissionMiddlewareForHuma("articleService", model.PermissionCreator))

	huma.Register(securedArticle, huma.Operation{
		OperationID: "createArticle",
		Method:      http.MethodPost,
		Path:        "/",
		Summary:     "CreateArticle",
		Description: "Create a new article",
		Tags:        []string{"article"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, articleHandler.HandleCreateArticle)

	huma.Register(securedArticle, huma.Operation{
		OperationID: "updateArticle",
		Method:      http.MethodPatch,
		Path:        "/{articleID}",
		Summary:     "UpdateArticle",
		Description: "Update article information",
		Tags:        []string{"article"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, articleHandler.HandleUpdateArticle)

	huma.Register(securedArticle, huma.Operation{
		OperationID: "deleteArticle",
		Method:      http.MethodDelete,
		Path:        "/{articleID}",
		Summary:     "DeleteArticle",
		Description: "Delete article",
		Tags:        []string{"article"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, articleHandler.HandleDeleteArticle)

	huma.Register(securedArticle, huma.Operation{
		OperationID: "updateArticleStatus",
		Method:      http.MethodPut,
		Path:        "/{articleID}/status",
		Summary:     "UpdateArticleStatus",
		Description: "Update article publish status",
		Tags:        []string{"article"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, articleHandler.HandleUpdateArticleStatus)

	initArticleVersionRouter(articleGroup)
}

func initArticleVersionRouter(articleGroup *huma.Group) {
	versionHandler := handler.NewArticleVersionHandler()

	versionBase := huma.NewGroup(articleGroup, "/{articleID}/version")

	huma.Register(versionBase, huma.Operation{
		OperationID: "getLatestArticleVersion",
		Method:      http.MethodGet,
		Path:        "/latest",
		Summary:     "GetLatestArticleVersion",
		Description: "Get latest article version content",
		Tags:        []string{"articleVersion"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, versionHandler.HandleGetLatestArticleVersionInfo)

	securedVersion := huma.NewGroup(versionBase, "")
	securedVersion.UseMiddleware(middleware.LimitUserPermissionMiddlewareForHuma("articleVersionService", model.PermissionCreator))

	huma.Register(securedVersion, huma.Operation{
		OperationID: "listArticleVersions",
		Method:      http.MethodGet,
		Path:        "/list",
		Summary:     "ListArticleVersions",
		Description: "List article versions",
		Tags:        []string{"articleVersion"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, versionHandler.HandleListArticleVersions)

	versionCreate := huma.NewGroup(securedVersion, "")
	versionCreate.UseMiddleware(middleware.RateLimiterMiddlewareForHuma("createArticleVersion", constant.CtxKeyUserID, 10*time.Second, 1))

	huma.Register(versionCreate, huma.Operation{
		OperationID: "createArticleVersion",
		Method:      http.MethodPost,
		Path:        "/",
		Summary:     "CreateArticleVersion",
		Description: "Create a new article version",
		Tags:        []string{"articleVersion"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, versionHandler.HandleCreateArticleVersion)

	huma.Register(securedVersion, huma.Operation{
		OperationID: "getArticleVersionInfo",
		Method:      http.MethodGet,
		Path:        "/v{version}",
		Summary:     "GetArticleVersionInfo",
		Description: "Get article version info",
		Tags:        []string{"articleVersion"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, versionHandler.HandleGetArticleVersionInfo)
}
