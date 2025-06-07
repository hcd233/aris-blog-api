package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/aris-blog-api/internal/handler"
)

func initOauth2Router(r *gin.RouterGroup) {
	githubOauth2Handler := handler.NewGithubOauth2Handler()
	// qqOauth2Handler := handler.NewQQOauth2Handler()
	googleOauth2Handler := handler.NewGoogleOauth2Handler()

	oauth2Group := r.Group("/oauth2")
	{
		// GitHub OAuth2路由
		githubRouter := oauth2Group.Group("/github")
		{
			githubRouter.GET("login", githubOauth2Handler.HandleLogin)
			githubRouter.GET("callback", githubOauth2Handler.HandleCallback)
		}

		// Google OAuth2路由
		googleRouter := oauth2Group.Group("/google")
		{
			googleRouter.GET("login", googleOauth2Handler.HandleLogin)
			googleRouter.GET("callback", googleOauth2Handler.HandleCallback)
		}

		// QQ OAuth2路由
		// qqRouter := oauth2Group.Group("/qq")
		// {
		// 	qqRouter.GET("login", qqOauth2Handler.HandleLogin)
		// 	qqRouter.GET("callback", qqOauth2Handler.HandleCallback)
		// }
	}
}
