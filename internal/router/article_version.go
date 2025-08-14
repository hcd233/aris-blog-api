package router

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
)

func initArticleVersionRouter(r fiber.Router) {
	articleVersionHandler := handler.NewArticleVersionHandler()

	r.Get("/version/latest", articleVersionHandler.HandleGetLatestArticleVersionInfo)
	articleVersionRouter := r.Group("/version", middleware.LimitUserPermissionMiddleware("articleVersionService", model.PermissionCreator))
	{
		articleVersionRouter.Get("/list", middleware.ValidateParamMiddleware(&protocol.PageParam{}), articleVersionHandler.HandleListArticleVersions)

		articleVersionRouter.Post(
			"/",
			middleware.RateLimiterMiddleware("createArticleVersion", constant.CtxKeyUserID, 10*time.Second, 1),
			middleware.ValidateBodyMiddleware(&protocol.CreateArticleVersionBody{}),
			articleVersionHandler.HandleCreateArticleVersion,
		)
		articleVersionRouter.Get("/v:version", middleware.ValidateURIMiddleware(&protocol.ArticleVersionURI{}), articleVersionHandler.HandleGetArticleVersionInfo)
	}
}
