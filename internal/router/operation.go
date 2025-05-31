package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
	"github.com/hcd233/aris-blog-api/internal/protocol"
)

func initOperationRouter(r *gin.RouterGroup) {
	operationHandler := handler.NewOperationHandler()

	operationRouter := r.Group("/operation", middleware.JwtMiddleware())
	{
		userLikeRouter := operationRouter.Group("/like")
		{
			userLikeRouter.POST(
				"article",
				middleware.RateLimiterMiddleware("likeArticle", constant.CtxKeyUserID, 10*time.Second, 2),
				middleware.ValidateBodyMiddleware(&protocol.LikeArticleBody{}),
				operationHandler.HandleUserLikeArticle,
			)
			userLikeRouter.POST(
				"comment",
				middleware.RateLimiterMiddleware("likeComment", constant.CtxKeyUserID, 2*time.Second, 2),
				middleware.ValidateBodyMiddleware(&protocol.LikeCommentBody{}),
				operationHandler.HandleUserLikeComment,
			)
			userLikeRouter.POST(
				"tag",
				middleware.RateLimiterMiddleware("likeTag", constant.CtxKeyUserID, 10*time.Second, 2),
				middleware.ValidateBodyMiddleware(&protocol.LikeTagBody{}),
				operationHandler.HandleUserLikeTag,
			)
		}
		viewRouter := operationRouter.Group("/view")
		{
			viewRouter.POST(
				"article",
				middleware.RateLimiterMiddleware("logUserViewArticle", constant.CtxKeyUserID, 10*time.Second, 2),
				middleware.ValidateBodyMiddleware(&protocol.LogUserViewArticleBody{}),
				operationHandler.HandleLogUserViewArticle,
			)
		}
	}
}
