package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
)

func initArticleRouter(r fiber.Router) {
	articleHandler := handler.NewArticleHandler()

	articleRouter := r.Group("/article", middleware.JwtMiddleware())
	{
		articleRouter.Get("/list", middleware.ValidateParamMiddleware(&protocol.PageParam{}), articleHandler.HandleListArticles)

		articleRouter.Get("/slug/:authorName/:articleSlug",
			middleware.ValidateURIMiddleware(&protocol.ArticleSlugURI{}),
			articleHandler.HandleGetArticleInfoBySlug)

		articleRouter.Post(
			"/",
			middleware.LimitUserPermissionMiddleware("articleService", model.PermissionCreator),
			middleware.ValidateBodyMiddleware(&protocol.CreateArticleBody{}),
			articleHandler.HandleCreateArticle,
		)
		articleIDRouter := articleRouter.Group("/:articleID", middleware.ValidateURIMiddleware(&protocol.ArticleURI{}))
		{
			articleIDRouter.Get("/", articleHandler.HandleGetArticleInfo)
			articleIDRouter.Patch(
				"/",
				middleware.LimitUserPermissionMiddleware("articleService", model.PermissionCreator),
				middleware.ValidateBodyMiddleware(&protocol.UpdateArticleBody{}),
				articleHandler.HandleUpdateArticle,
			)
			articleIDRouter.Delete(
				"/",
				middleware.LimitUserPermissionMiddleware("articleService", model.PermissionCreator),
				articleHandler.HandleDeleteArticle,
			)
			articleIDRouter.Put(
				"/status",
				middleware.LimitUserPermissionMiddleware("articleService", model.PermissionCreator),
				middleware.ValidateBodyMiddleware(&protocol.UpdateArticleStatusBody{}),
				articleHandler.HandleUpdateArticleStatus,
			)

			initArticleVersionRouter(articleIDRouter)
		}
	}
}
