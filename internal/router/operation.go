package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/middleware"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/service"
)

func initUserOperationRouter(r *gin.RouterGroup) {
	operationService := service.NewOperationService()

	operationRouter := r.Group("/operation")
	{
		userLikeRouter := operationRouter.Group("/like")
		{
			userLikeRouter.POST(
				"article",
				middleware.RateLimiterMiddleware(10*time.Second, 2, "userID", protocol.CodeLikeArticleRateLimitError),
				middleware.ValidateBodyMiddleware(&protocol.LikeArticleBody{}),
				operationService.UserLikeArticleHandler,
			)
			userLikeRouter.POST(
				"comment",
				middleware.RateLimiterMiddleware(2*time.Second, 2, "userID", protocol.CodeLikeCommentRateLimitError),
				middleware.ValidateBodyMiddleware(&protocol.LikeCommentBody{}),
				operationService.UserLikeCommentHandler,
			)
			userLikeRouter.POST(
				"tag",
				middleware.RateLimiterMiddleware(10*time.Second, 2, "userID", protocol.CodeLikeTagRateLimitError),
				middleware.ValidateBodyMiddleware(&protocol.LikeTagBody{}),
				operationService.UserLikeTagHandler,
			)
		}
		viewRouter := operationRouter.Group("/view")
		{
			viewRouter.POST(
				"article",
				middleware.RateLimiterMiddleware(10*time.Second, 2, "userID", protocol.CodeLogUserViewRateLimitError),
				middleware.ValidateBodyMiddleware(&protocol.LogUserViewArticleBody{}),
				operationService.LogUserViewArticleHandler,
			)
		}
	}
}
