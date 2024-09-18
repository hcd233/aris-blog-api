package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-AI-go/internal/middleware"
	"github.com/hcd233/Aris-AI-go/internal/router/v1/oauth2"
	"github.com/hcd233/Aris-AI-go/internal/router/v1/user"
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

		userRouter := v1Router.Group("/user", middleware.JwtMiddleware())
		{
			userRouter.GET("/", user.QueryUserHandler)
			userRouter.GET("/:userName", user.GetInfoHandler)
			userRouter.POST("/:userName", user.UpdateInfoHandler)
		}
	}
}
