package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
)

func initCommentRouter(commentGroup *huma.Group) {
	commentHandler := handler.NewCommentHandler()

	commentGroup.UseMiddleware(middleware.JwtMiddlewareForHuma())

	// 创建文章评论
	huma.Register(commentGroup, huma.Operation{
		OperationID: "createArticleComment",
		Method:      http.MethodPost,
		Path:        "/",
		Summary:     "CreateArticleComment",
		Description: "Create a comment on an article",
		Tags:        []string{"comment"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, commentHandler.HandleCreateArticleComment)

	// 删除评论
	huma.Register(commentGroup, huma.Operation{
		OperationID: "deleteComment",
		Method:      http.MethodDelete,
		Path:        "/{commentID}",
		Summary:     "DeleteComment",
		Description: "Delete a comment",
		Tags:        []string{"comment"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, commentHandler.HandleDeleteComment)

	// 列出文章评论
	huma.Register(commentGroup, huma.Operation{
		OperationID: "listArticleComments",
		Method:      http.MethodGet,
		Path:        "/article/{articleID}/list",
		Summary:     "ListArticleComments",
		Description: "Get a paginated list of top-level comments for an article",
		Tags:        []string{"comment"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, commentHandler.HandleListArticleComments)

	// 列出子评论
	huma.Register(commentGroup, huma.Operation{
		OperationID: "listChildrenComments",
		Method:      http.MethodGet,
		Path:        "/{commentID}/subComments",
		Summary:     "ListChildrenComments",
		Description: "Get a paginated list of child comments for a specific comment",
		Tags:        []string{"comment"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, commentHandler.HandleListChildrenComments)
}
