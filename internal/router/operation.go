package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
	"github.com/hcd233/aris-blog-api/internal/protocol"
)

func initOperationRouter(r *gin.RouterGroup) {
	operationHandler := handler.NewOperationHandler()

	operationRouter := r.Group("/operation")
	{
		userLikeRouter := operationRouter.Group("/like")
		{
			userLikeRouter.POST(
				"article",
				middleware.RateLimiterMiddleware("likeArticle", "userID", 10*time.Second, 2),
				middleware.ValidateBodyMiddleware(&protocol.LikeArticleBody{}),
				operationHandler.HandleUserLikeArticle,
			)
			userLikeRouter.POST(
				"comment",
				middleware.RateLimiterMiddleware("likeComment", "userID", 2*time.Second, 2),
				middleware.ValidateBodyMiddleware(&protocol.LikeCommentBody{}),
				operationHandler.HandleUserLikeComment,
			)
			userLikeRouter.POST(
				"tag",
				middleware.RateLimiterMiddleware("likeTag", "userID", 10*time.Second, 2),
				middleware.ValidateBodyMiddleware(&protocol.LikeTagBody{}),
				operationHandler.HandleUserLikeTag,
			)
		}
		viewRouter := operationRouter.Group("/view")
		{
			viewRouter.POST(
				"article",
				middleware.RateLimiterMiddleware("logUserViewArticle", "userID", 10*time.Second, 2),
				middleware.ValidateBodyMiddleware(&protocol.LogUserViewArticleBody{}),
				operationHandler.HandleLogUserViewArticle,
			)
		}
	}
}
