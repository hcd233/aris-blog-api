package router

import (
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
)

func initOperationRouter(operationGroup *huma.Group) {
	operationHandler := handler.NewOperationHandler()

	operationGroup.UseMiddleware(middleware.JwtMiddleware())

	likeArticleGroup := huma.NewGroup(operationGroup, "/like")
	likeArticleGroup.UseMiddleware(middleware.RateLimiterMiddleware("likeArticle", constant.CtxKeyUserID, 10*time.Second, 2))

	huma.Register(likeArticleGroup, huma.Operation{
		OperationID: "likeArticle",
		Method:      http.MethodPost,
		Path:        "/article",
		Summary:     "LikeArticle",
		Description: "Like or unlike an article",
		Tags:        []string{"operation"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, operationHandler.HandleUserLikeArticle)

	likeCommentGroup := huma.NewGroup(operationGroup, "/like")
	likeCommentGroup.UseMiddleware(middleware.RateLimiterMiddleware("likeComment", constant.CtxKeyUserID, 2*time.Second, 2))

	huma.Register(likeCommentGroup, huma.Operation{
		OperationID: "likeComment",
		Method:      http.MethodPost,
		Path:        "/comment",
		Summary:     "LikeComment",
		Description: "Like or unlike a comment",
		Tags:        []string{"operation"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, operationHandler.HandleUserLikeComment)

	likeTagGroup := huma.NewGroup(operationGroup, "/like")
	likeTagGroup.UseMiddleware(middleware.RateLimiterMiddleware("likeTag", constant.CtxKeyUserID, 10*time.Second, 2))

	huma.Register(likeTagGroup, huma.Operation{
		OperationID: "likeTag",
		Method:      http.MethodPost,
		Path:        "/tag",
		Summary:     "LikeTag",
		Description: "Like or unlike a tag",
		Tags:        []string{"operation"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, operationHandler.HandleUserLikeTag)

	viewGroup := huma.NewGroup(operationGroup, "/view")
	viewGroup.UseMiddleware(middleware.RateLimiterMiddleware("logUserViewArticle", constant.CtxKeyUserID, 10*time.Second, 2))

	huma.Register(viewGroup, huma.Operation{
		OperationID: "logArticleView",
		Method:      http.MethodPost,
		Path:        "/article",
		Summary:     "LogArticleView",
		Description: "Log article view progress",
		Tags:        []string{"operation"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, operationHandler.HandleLogUserViewArticle)
}
