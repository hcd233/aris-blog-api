package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/handler"
	"github.com/hcd233/Aris-blog/internal/middleware"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
)

func initArticleRouter(r *gin.RouterGroup) {
	articleHandler := handler.NewArticleHandler()

	r.GET("articles", middleware.ValidateParamMiddleware(&protocol.PageParam{}), articleHandler.HandleListArticles)
	articleRouter := r.Group("/article", middleware.JwtMiddleware())
	{
		articleRouter.GET("", middleware.ValidateParamMiddleware(&protocol.QueryParam{}), articleHandler.HandleQueryArticle)
	}
}

func initUserArticleRouter(r *gin.RouterGroup) {
	articleHandler := handler.NewArticleHandler()
	r.GET("articles", middleware.ValidateParamMiddleware(&protocol.PageParam{}), articleHandler.HandleListUserArticles)
	articleRouter := r.Group("/article")
	{
		articleRouter.GET("", middleware.ValidateParamMiddleware(&protocol.QueryParam{}), articleHandler.HandleQueryUserArticle)
		articleRouter.POST(
			"",
			middleware.LimitUserPermissionMiddleware(model.PermissionCreator),
			middleware.ValidateBodyMiddleware(&protocol.CreateArticleBody{}),
			articleHandler.HandleCreateArticle,
		)
	}

	articleSlugRouter := articleRouter.Group("/:articleSlug", middleware.ValidateURIMiddleware(&protocol.ArticleSlugURI{}))
	{
		articleSlugRouter.GET("", articleHandler.HandleGetArticleInfo)
		articleSlugRouter.PUT(
			"",
			middleware.LimitUserPermissionMiddleware(model.PermissionCreator),
			middleware.ValidateBodyMiddleware(&protocol.UpdateArticleBody{}),
			articleHandler.HandleUpdateArticle,
		)
		articleSlugRouter.DELETE(
			"",
			middleware.LimitUserPermissionMiddleware(model.PermissionCreator),
			articleHandler.HandleDeleteArticle,
		)
		articleSlugRouter.PUT(
			"status",
			middleware.LimitUserPermissionMiddleware(model.PermissionCreator),
			middleware.ValidateBodyMiddleware(&protocol.UpdateArticleStatusBody{}),
			articleHandler.HandleUpdateArticleStatus,
		)

		initArticleVersionRouter(articleSlugRouter)
		initArticleCommentRouter(articleSlugRouter)
	}
}
