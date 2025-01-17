package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/aris-blog-api/internal/handler"
)

func initOauth2Router(r *gin.RouterGroup) {
	githubOauth2Handler := handler.NewGithubOauth2Handler()
	oauth2Group := r.Group("/oauth2")
	{
		githubRouter := oauth2Group.Group("/github")
		{
			githubRouter.GET("login", githubOauth2Handler.HandleLogin)
			githubRouter.GET("callback", githubOauth2Handler.HandleCallback)
		}
	}
}
