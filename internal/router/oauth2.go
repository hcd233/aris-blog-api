package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/service"
)

func initOauth2Router(r *gin.RouterGroup) {
	githubOauth2Service := service.NewGithubOauth2Service()
	oauth2Group := r.Group("/oauth2")
	{
		githubRouter := oauth2Group.Group("/github")
		{
			githubRouter.GET("login", githubOauth2Service.LoginHandler)
			githubRouter.GET("callback", githubOauth2Service.CallbackHandler)
		}
	}
}
