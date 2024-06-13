// Package oauth2 login and callback handlers.
package oauth2

import "github.com/gin-gonic/gin"

// InitOauth2Router initializes the OAuth2 router.
func InitOauth2Router(r *gin.RouterGroup) {
	oauth2Group := r.Group("/oauth2")
	{
		githubRouter := oauth2Group.Group("/github")
		{
			githubRouter.GET("/login", handleGithubLogin)
			githubRouter.GET("/callback", handleGithubCallback)
		}
	}
}
