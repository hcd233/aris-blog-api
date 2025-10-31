package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
)

func initOperationRouter(operationGroup *huma.Group) {
	operationHandler := handler.NewOperationHandler()

	operationGroup.UseMiddleware(middleware.JwtMiddlewareForHuma())

	// 点赞文章
	huma.Register(operationGroup, huma.Operation{
		OperationID: "likeArticle",
		Method:      http.MethodPost,
		Path:        "/like/article",
		Summary:     "LikeArticle",
		Description: "Like or unlike an article",
		Tags:        []string{"operation"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, operationHandler.HandleUserLikeArticle)

	// 点赞评论
	huma.Register(operationGroup, huma.Operation{
		OperationID: "likeComment",
		Method:      http.MethodPost,
		Path:        "/like/comment",
		Summary:     "LikeComment",
		Description: "Like or unlike a comment",
		Tags:        []string{"operation"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, operationHandler.HandleUserLikeComment)

	// 点赞标签
	huma.Register(operationGroup, huma.Operation{
		OperationID: "likeTag",
		Method:      http.MethodPost,
		Path:        "/like/tag",
		Summary:     "LikeTag",
		Description: "Like or unlike a tag",
		Tags:        []string{"operation"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, operationHandler.HandleUserLikeTag)

	// 记录文章浏览
	huma.Register(operationGroup, huma.Operation{
		OperationID: "logArticleView",
		Method:      http.MethodPost,
		Path:        "/view/article",
		Summary:     "LogArticleView",
		Description: "Log article view with reading progress",
		Tags:        []string{"operation"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, operationHandler.HandleLogUserViewArticle)
}
