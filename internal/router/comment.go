package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/handler"
	"github.com/hcd233/Aris-blog/internal/middleware"
	"github.com/hcd233/Aris-blog/internal/protocol"
)

func initArticleCommentRouter(r *gin.RouterGroup) {
	commentService := handler.NewCommentService()

	r.GET("comments", middleware.ValidateParamMiddleware(&protocol.PageParam{}), commentService.ListArticleCommentsHandler)
	commentRouter := r.Group("/comment")
	{
		commentRouter.POST(
			"",
			middleware.RateLimiterMiddleware(10*time.Second, 1, "createComment", "userID", protocol.CodeCreateCommentRateLimitError),
			middleware.ValidateBodyMiddleware(&protocol.CreateArticleCommentBody{}),
			commentService.CreateArticleCommentHandler,
		)
		commentIDRouter := commentRouter.Group(":commentID", middleware.ValidateURIMiddleware(&protocol.CommentURI{}))
		{
			commentIDRouter.GET("", commentService.GetCommentInfoHandler)
			commentIDRouter.DELETE("", commentService.DeleteCommentHandler)
			commentIDRouter.GET("subComments", middleware.ValidateParamMiddleware(&protocol.PageParam{}), commentService.ListChildrenCommentsHandler)
		}
	}
}
