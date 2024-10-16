package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/middleware"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/router/v1/article"
	article_version "github.com/hcd233/Aris-blog/internal/router/v1/article_version"
	"github.com/hcd233/Aris-blog/internal/router/v1/category"
	"github.com/hcd233/Aris-blog/internal/router/v1/oauth2"
	"github.com/hcd233/Aris-blog/internal/router/v1/tag"
	"github.com/hcd233/Aris-blog/internal/router/v1/user"
)

// InitRouter initializes the router.
func InitRouter(r *gin.Engine) {
	rootGroup := r.Group("/")
	{
		rootGroup.GET("/", RootHandler)
	}

	v1Router := r.Group("/v1")
	{
		initOauth2Router(v1Router)
		initTagRouter(v1Router)
		initUserRouter(v1Router)
	}
}

func initOauth2Router(r *gin.RouterGroup) {
	oauth2Group := r.Group("/oauth2")
	{
		githubRouter := oauth2Group.Group("/github")
		{
			githubRouter.GET("/login", oauth2.GithubLoginHandler)
			githubRouter.GET("/callback", oauth2.GithubCallbackHandler)
		}
	}
}

func initTagRouter(r *gin.RouterGroup) {
	tagRouter := r.Group("/tag", middleware.JwtMiddleware())
	{
		tagRouter.GET("", middleware.ValidateParamMiddleware(&protocol.QueryParam{}), tag.SearchTagHandler)
		tagRouter.POST("", middleware.ValidateBodyMiddleware(&protocol.CreateTagBody{}), tag.CreateTagHandler)
		tagRouter.GET("list", middleware.ValidateParamMiddleware(&protocol.PageParam{}), tag.ListTagsHandler)

		tagSlugRouter := tagRouter.Group("/:tagSlug", middleware.ValidateURIMiddleware(&protocol.TagURI{}))
		{
			tagSlugRouter.GET("", tag.GetTagInfoHandler)
			tagSlugRouter.PUT("", middleware.ValidateBodyMiddleware(&protocol.UpdateTagBody{}), tag.UpdateTagHandler)
			tagSlugRouter.DELETE("", tag.DeleteTagHandler)
		}
	}
}

func initUserRouter(r *gin.RouterGroup) {
	userRouter := r.Group("/user", middleware.JwtMiddleware())
	{
		userRouter.GET("/", middleware.ValidateParamMiddleware(&protocol.QueryParam{}), user.QueryUserHandler)

		userNameRouter := userRouter.Group("/:userName", middleware.ValidateURIMiddleware(&protocol.UserURI{}))
		{
			userNameRouter.GET("", user.GetUserInfoHandler)
			userNameRouter.PUT("", middleware.ValidateBodyMiddleware(&protocol.UpdateUserBody{}), user.UpdateInfoHandler)

			initArticleRouter(userNameRouter)
			initCategoryRouter(userNameRouter)
		}

	}
}

func initArticleRouter(r *gin.RouterGroup) {
	articleRouter := r.Group("/article")
	{
		articleRouter.GET("/list", middleware.ValidateParamMiddleware(&protocol.PageParam{}), article.ListArticlesHandler)
		articleRouter.POST("", middleware.ValidateBodyMiddleware(&protocol.CreateArticleBody{}), article.CreateArticleHandler)
	}

	articleSlugRouter := articleRouter.Group("/:articleSlug", middleware.ValidateURIMiddleware(&protocol.ArticleURI{}))
	{
		articleSlugRouter.GET("", article.GetArticleInfoHandler)
		articleSlugRouter.PUT("", middleware.ValidateBodyMiddleware(&protocol.UpdateArticleBody{}), article.UpdateArticleHandler)
		articleSlugRouter.DELETE("", article.DeleteArticleHandler)

		initArticleVersionRouter(articleSlugRouter)
	}
}

func initCategoryRouter(r *gin.RouterGroup) {
	categoryRouter := r.Group("/category")
	{
		categoryRouter.GET("list", middleware.ValidateParamMiddleware(&protocol.PageParam{}), category.ListRootCategoriesHandler)
		categoryRouter.POST("", middleware.ValidateBodyMiddleware(&protocol.CreateCategoryBody{}), category.CreateCategoryHandler)
	}

	categoryIDRouter := categoryRouter.Group("/:categoryID", middleware.ValidateURIMiddleware(&protocol.CategoryURI{}))
	{
		categoryIDRouter.GET("", category.GetCategoryInfoHandler)
		categoryIDRouter.DELETE("", category.DeleteCategoryHandler)
		categoryIDRouter.PUT("", middleware.ValidateBodyMiddleware(&protocol.UpdateCategoryBody{}), category.UpdateCategoryInfoHandler)
		categoryIDRouter.GET("subCategory", middleware.ValidateParamMiddleware(&protocol.PageParam{}), category.ListChildrenCategoriesHandler)
		categoryIDRouter.GET("subArticle", middleware.ValidateParamMiddleware(&protocol.PageParam{}), category.ListChildrenArticlesHandler)
	}
}

func initArticleVersionRouter(r *gin.RouterGroup) {
	articleVersionRouter := r.Group("/version")
	{
		articleVersionRouter.GET("list", middleware.ValidateParamMiddleware(&protocol.PageParam{}), article_version.ListArticleVersionsHandler)
		// articleVersionRouter.POST("", middleware.ValidateBodyMiddleware(&protocol.CreateArticleVersionBody{}), article.CreateArticleVersionHandler)
		// articleVersionRouter.GET("v:version", middleware.ValidateParamMiddleware(&protocol.QueryArticleVersionParam{}), article.GetArticleVersionHandler)
	}
}
