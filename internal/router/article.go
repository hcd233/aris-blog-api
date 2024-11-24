package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/middleware"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"github.com/hcd233/Aris-blog/internal/service"
)

func initArticleRouter(r *gin.RouterGroup) {
	articleService := service.NewArticleService()

	r.GET("articles", middleware.ValidateParamMiddleware(&protocol.PageParam{}), articleService.ListArticlesHandler)
	articleRouter := r.Group("/article", middleware.JwtMiddleware())
	{
		articleRouter.GET("", middleware.ValidateParamMiddleware(&protocol.QueryParam{}), articleService.QueryArticleHandler)
	}
}

func initUserArticleRouter(r *gin.RouterGroup) {
	articleService := service.NewArticleService()
	r.GET("articles", middleware.ValidateParamMiddleware(&protocol.PageParam{}), articleService.ListUserArticlesHandler)
	articleRouter := r.Group("/article")
	{
		articleRouter.GET("", middleware.ValidateParamMiddleware(&protocol.QueryParam{}), articleService.QueryUserArticleHandler)
		articleRouter.POST(
			"",
			middleware.LimitUserPermissionMiddleware(model.PermissionCreator),
			middleware.ValidateBodyMiddleware(&protocol.CreateArticleBody{}),
			articleService.CreateArticleHandler,
		)
	}

	articleSlugRouter := articleRouter.Group("/:articleSlug", middleware.ValidateURIMiddleware(&protocol.ArticleSlugURI{}))
	{
		articleSlugRouter.GET("", articleService.GetArticleInfoHandler)
		articleSlugRouter.PUT(
			"",
			middleware.LimitUserPermissionMiddleware(model.PermissionCreator),
			middleware.ValidateBodyMiddleware(&protocol.UpdateArticleBody{}),
			articleService.UpdateArticleHandler,
		)
		articleSlugRouter.DELETE(
			"",
			middleware.LimitUserPermissionMiddleware(model.PermissionCreator),
			articleService.DeleteArticleHandler,
		)
		articleSlugRouter.PUT(
			"status",
			middleware.LimitUserPermissionMiddleware(model.PermissionCreator),
			middleware.ValidateBodyMiddleware(&protocol.UpdateArticleStatusBody{}),
			articleService.UpdateArticleStatusHandler,
		)

		initArticleVersionRouter(articleSlugRouter)
		initArticleCommentRouter(articleSlugRouter)
	}
}
