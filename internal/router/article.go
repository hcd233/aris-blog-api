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

	r.GET("articles", middleware.ValidateParamMiddleware(&protocol.PageParam{}), articleHandler.HandleListArticles)
	articleRouter := r.Group("/article", middleware.JwtMiddleware())
	{
		articleRouter.POST(
			"",
			middleware.LimitUserPermissionMiddleware("articleService", model.PermissionCreator),
			middleware.ValidateBodyMiddleware(&protocol.CreateArticleBody{}),
			articleHandler.HandleCreateArticle,
		)
		articleSlugRouter := articleRouter.Group("/:articleSlug", middleware.ValidateURIMiddleware(&protocol.ArticleSlugURI{}))
		{
			articleSlugRouter.GET("", articleHandler.HandleGetArticleInfo)
			articleSlugRouter.PUT(
				"",
				middleware.LimitUserPermissionMiddleware("articleService", model.PermissionCreator),
				middleware.ValidateBodyMiddleware(&protocol.UpdateArticleBody{}),
				articleHandler.HandleUpdateArticle,
			)
			articleSlugRouter.DELETE(
				"",
				middleware.LimitUserPermissionMiddleware("articleService", model.PermissionCreator),
				articleHandler.HandleDeleteArticle,
			)
			articleSlugRouter.PUT(
				"status",
				middleware.LimitUserPermissionMiddleware("articleService", model.PermissionCreator),
				middleware.ValidateBodyMiddleware(&protocol.UpdateArticleStatusBody{}),
				articleHandler.HandleUpdateArticleStatus,
			)

			initArticleVersionRouter(articleSlugRouter)
			initArticleCommentRouter(articleSlugRouter)
		}
	}
}
