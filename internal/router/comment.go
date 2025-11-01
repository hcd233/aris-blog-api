package router

import (
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
)

func initCommentRouter(commentGroup *huma.Group) {
	commentHandler := handler.NewCommentHandler()

	commentGroup.UseMiddleware(middleware.JwtMiddleware())

	listGroup := huma.NewGroup(commentGroup, "")

	huma.Register(listGroup, huma.Operation{
		OperationID: "listArticleComments",
		Method:      http.MethodGet,
		Path:        "/article/{articleID}/list",
		Summary:     "ListArticleComments",
		Description: "List comments for an article",
		Tags:        []string{"comment"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, commentHandler.HandleListArticleComments)

	huma.Register(listGroup, huma.Operation{
		OperationID: "listChildrenComments",
		Method:      http.MethodGet,
		Path:        "/{commentID}/subComments",
		Summary:     "ListChildrenComments",
		Description: "List child comments of a comment",
		Tags:        []string{"comment"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, commentHandler.HandleListChildrenComments)

	createGroup := huma.NewGroup(commentGroup, "")
	createGroup.UseMiddleware(middleware.RateLimiterMiddleware("createComment", constant.CtxKeyUserID, 10*time.Second, 1))

	huma.Register(createGroup, huma.Operation{
		OperationID: "createComment",
		Method:      http.MethodPost,
		Path:        "/",
		Summary:     "CreateComment",
		Description: "Create a comment on an article",
		Tags:        []string{"comment"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, commentHandler.HandleCreateArticleComment)

	huma.Register(commentGroup, huma.Operation{
		OperationID: "deleteComment",
		Method:      http.MethodDelete,
		Path:        "/{commentID}",
		Summary:     "DeleteComment",
		Description: "Delete a comment",
		Tags:        []string{"comment"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, commentHandler.HandleDeleteComment)
}
