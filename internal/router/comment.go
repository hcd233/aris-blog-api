package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
	"github.com/hcd233/aris-blog-api/internal/protocol"
)

func initCommentRouter(r *gin.RouterGroup) {
	commentHandler := handler.NewCommentHandler()

	commentRouter := r.Group("/comment", middleware.JwtMiddleware())
	{
		commentRouter.GET("article/:articleID/list",
			middleware.ValidateParamMiddleware(&protocol.PageParam{}),
			middleware.ValidateURIMiddleware(&protocol.ArticleURI{}),
			commentHandler.HandleListArticleComments,
		)
		commentRouter.POST(
			"",
			middleware.RateLimiterMiddleware("createComment", constant.CtxKeyUserID, 10*time.Second, 1),
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
