package router

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
	"github.com/hcd233/aris-blog-api/internal/protocol"
)

func initCommentRouter(r fiber.Router) {
	commentHandler := handler.NewCommentHandler()

	commentRouter := r.Group("/comment", middleware.JwtMiddleware())
	{
		commentRouter.Get("/article/:articleID/list",
			middleware.ValidateParamMiddleware(&protocol.PageParam{}),
			middleware.ValidateURIMiddleware(&protocol.ArticleURI{}),
			commentHandler.HandleListArticleComments,
		)
		commentRouter.Post(
			"/",
			middleware.RateLimiterMiddleware("createComment", constant.CtxKeyUserID, 10*time.Second, 1),
			middleware.ValidateBodyMiddleware(&protocol.CreateArticleCommentBody{}),
			commentHandler.HandleCreateArticleComment,
		)
		commentIDRouter := commentRouter.Group("/:commentID", middleware.ValidateURIMiddleware(&protocol.CommentURI{}))
		{
			commentIDRouter.Delete("/", commentHandler.HandleDeleteComment)
			commentIDRouter.Get("/subComments", middleware.ValidateParamMiddleware(&protocol.PageParam{}), commentHandler.HandleListChildrenComments)
		}
	}
}
