package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-blog-api/internal/handler"
)

func initOauth2Router(r fiber.Router) {
	githubOauth2Handler := handler.NewGithubOauth2Handler()
	// qqOauth2Handler := handler.NewQQOauth2Handler()
	googleOauth2Handler := handler.NewGoogleOauth2Handler()

	oauth2Group := r.Group("/oauth2")
	{
		// GitHub OAuth2路由
		githubRouter := oauth2Group.Group("/github")
		{
			githubRouter.Get("/login", githubOauth2Handler.HandleLogin)
			githubRouter.Get("/callback", githubOauth2Handler.HandleCallback)
		}

		// Google OAuth2路由
		googleRouter := oauth2Group.Group("/google")
		{
			googleRouter.Get("/login", googleOauth2Handler.HandleLogin)
			googleRouter.Get("/callback", googleOauth2Handler.HandleCallback)
		}

		// QQ OAuth2路由
		// qqRouter := oauth2Group.Group("/qq")
		// {
		// 	qqRouter.Get("/login", qqOauth2Handler.HandleLogin)
		// 	qqRouter.Get("/callback", qqOauth2Handler.HandleCallback)
		// }
	}
}
