package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
)

func initArticleRouter(r *gin.RouterGroup) {
	articleHandler := handler.NewArticleHandler()

	articleRouter := r.Group("/article", middleware.JwtMiddleware())
	{
		articleRouter.GET("list", middleware.ValidateParamMiddleware(&protocol.PageParam{}), articleHandler.HandleListArticles)

		articleRouter.GET("/slug/:authorName/:articleSlug",
			middleware.ValidateURIMiddleware(&protocol.ArticleSlugURI{}),
			articleHandler.HandleGetArticleInfoBySlug)

		articleRouter.POST(
			"",
			middleware.LimitUserPermissionMiddleware("articleService", model.PermissionCreator),
			middleware.ValidateBodyMiddleware(&protocol.CreateArticleBody{}),
			articleHandler.HandleCreateArticle,
		)
		articleIDRouter := articleRouter.Group("/:articleID", middleware.ValidateURIMiddleware(&protocol.ArticleURI{}))
		{
			articleIDRouter.GET("", articleHandler.HandleGetArticleInfo)
			articleIDRouter.PATCH(
				"",
				middleware.LimitUserPermissionMiddleware("articleService", model.PermissionCreator),
				middleware.ValidateBodyMiddleware(&protocol.UpdateArticleBody{}),
				articleHandler.HandleUpdateArticle,
			)
			articleIDRouter.DELETE(
				"",
				middleware.LimitUserPermissionMiddleware("articleService", model.PermissionCreator),
				articleHandler.HandleDeleteArticle,
			)
			articleIDRouter.PUT(
				"status",
				middleware.LimitUserPermissionMiddleware("articleService", model.PermissionCreator),
				middleware.ValidateBodyMiddleware(&protocol.UpdateArticleStatusBody{}),
				articleHandler.HandleUpdateArticleStatus,
			)

			initArticleVersionRouter(articleIDRouter)
		}
	}
}
