package router

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
	"github.com/hcd233/aris-blog-api/internal/protocol"
)

func initOperationRouter(r fiber.Router) {
	operationHandler := handler.NewOperationHandler()

	operationRouter := r.Group("/operation", middleware.JwtMiddleware())
	{
		userLikeRouter := operationRouter.Group("/like")
		{
			userLikeRouter.Post(
				"/article",
				middleware.RateLimiterMiddleware("likeArticle", constant.CtxKeyUserID, 10*time.Second, 2),
				middleware.ValidateBodyMiddleware(&protocol.LikeArticleBody{}),
				operationHandler.HandleUserLikeArticle,
			)
			userLikeRouter.Post(
				"/comment",
				middleware.RateLimiterMiddleware("likeComment", constant.CtxKeyUserID, 2*time.Second, 2),
				middleware.ValidateBodyMiddleware(&protocol.LikeCommentBody{}),
				operationHandler.HandleUserLikeComment,
			)
			userLikeRouter.Post(
				"/tag",
				middleware.RateLimiterMiddleware("likeTag", constant.CtxKeyUserID, 10*time.Second, 2),
				middleware.ValidateBodyMiddleware(&protocol.LikeTagBody{}),
				operationHandler.HandleUserLikeTag,
			)
		}
		viewRouter := operationRouter.Group("/view")
		{
			viewRouter.Post(
				"/article",
				middleware.RateLimiterMiddleware("logUserViewArticle", constant.CtxKeyUserID, 10*time.Second, 2),
				middleware.ValidateBodyMiddleware(&protocol.LogUserViewArticleBody{}),
				operationHandler.HandleLogUserViewArticle,
			)
		}
	}
}
