package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/handler"
	"github.com/hcd233/Aris-blog/internal/middleware"
	"github.com/hcd233/Aris-blog/internal/protocol"
)

func initArticleCommentRouter(r *gin.RouterGroup) {
	commentHandler := handler.NewCommentHandler()

	r.GET("comments", middleware.ValidateParamMiddleware(&protocol.PageParam{}), commentHandler.HandleListArticleComments)
	commentRouter := r.Group("/comment")
	{
		commentRouter.POST(
			"",
			middleware.RateLimiterMiddleware(10*time.Second, 1, "createComment", "userID", protocol.CodeCreateCommentRateLimitError),
			middleware.ValidateBodyMiddleware(&protocol.CreateArticleCommentBody{}),
			commentHandler.HandleCreateArticleComment,
		)
		commentIDRouter := commentRouter.Group(":commentID", middleware.ValidateURIMiddleware(&protocol.CommentURI{}))
		{
			commentIDRouter.GET("", commentHandler.HandleGetCommentInfo)
			commentIDRouter.DELETE("", commentHandler.HandleDeleteComment)
			commentIDRouter.GET("subComments", middleware.ValidateParamMiddleware(&protocol.PageParam{}), commentHandler.HandleListChildrenComments)
		}
	}
}
