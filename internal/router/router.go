package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/middleware"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/router/v1/article"
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
		oauth2Group := v1Router.Group("/oauth2")
		{
			githubRouter := oauth2Group.Group("/github")
			{
				githubRouter.GET("/login", oauth2.GithubLoginHandler)
				githubRouter.GET("/callback", oauth2.GithubCallbackHandler)
			}
		}

		tagRouter := v1Router.Group("/tag", middleware.JwtMiddleware())
		{
			tagRouter.GET("", middleware.ValidateParamMiddleware(&protocol.QueryParams{}), tag.SearchTagHandler)
			tagRouter.POST("", middleware.ValidateBodyMiddleware(&protocol.CreateTagBody{}), tag.CreateTagHandler)
			tagRouter.GET("list", middleware.ValidateParamMiddleware(&protocol.PageParams{}), tag.ListTagHandler)

			tagSlugRouter := tagRouter.Group("/:tagSlug", middleware.ValidateURIMiddleware(&protocol.TagURI{}))
			{
				tagSlugRouter.GET("", tag.GetTagInfoHandler)
				tagSlugRouter.PUT("", middleware.ValidateBodyMiddleware(&protocol.UpdateTagBody{}), tag.UpdateTagHandler)
				tagSlugRouter.DELETE("", tag.DeleteTagHandler)
			}
		}

		userRouter := v1Router.Group("/user", middleware.JwtMiddleware())
		{
			userRouter.GET("/", middleware.ValidateParamMiddleware(&protocol.QueryParams{}), user.QueryUserHandler)

			userNameRouter := userRouter.Group("/:userName", middleware.ValidateURIMiddleware(&protocol.UserURI{}))
			{
				userNameRouter.GET("", user.GetUserInfoHandler)
				userNameRouter.PUT("", middleware.ValidateBodyMiddleware(&protocol.UpdateUserBody{}), user.UpdateInfoHandler)

				articleRouter := userNameRouter.Group("/article")
				{
					articleRouter.GET("/list", middleware.ValidateParamMiddleware(&protocol.PageParams{}), article.ListArticleHandler)
					articleRouter.POST("", middleware.ValidateBodyMiddleware(&protocol.CreateArticleBody{}), article.CreateArticleHandler)
				}

				articleSlugRouter := articleRouter.Group("/:articleSlug", middleware.ValidateURIMiddleware(&protocol.ArticleURI{}))
				{
					articleSlugRouter.GET("", article.GetArticleInfoHandler)
					articleSlugRouter.PUT("", middleware.ValidateBodyMiddleware(&protocol.UpdateArticleBody{}), article.UpdateArticleHandler)
					articleSlugRouter.DELETE("", article.DeleteArticleHandler)
				}
			}

		}
	}
}
