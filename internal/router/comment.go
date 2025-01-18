package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
	"github.com/hcd233/aris-blog-api/internal/protocol"
)

func initCommentRouter(r *gin.RouterGroup) {
	commentHandler := handler.NewCommentHandler()

	r.GET("comments", middleware.ValidateParamMiddleware(&protocol.PageParam{}), commentHandler.HandleListArticleComments)
	commentRouter := r.Group("/comment")
	{
		commentRouter.POST(
			"",
			middleware.RateLimiterMiddleware("createComment", "userID", 10*time.Second, 1),
			middleware.ValidateBodyMiddleware(&protocol.CreateArticleCommentBody{}),
			commentHandler.HandleCreateArticleComment,
		)
		commentIDRouter := commentRouter.Group(":commentID", middleware.ValidateURIMiddleware(&protocol.CommentURI{}))
		{
			commentIDRouter.DELETE("", commentHandler.HandleDeleteComment)
			commentIDRouter.GET("subComments", middleware.ValidateParamMiddleware(&protocol.PageParam{}), commentHandler.HandleListChildrenComments)
		}
	}
}
