package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
)

func initArticleVersionRouter(r *gin.RouterGroup) {
	articleVersionHandler := handler.NewArticleVersionHandler()

	r.GET("latest", articleVersionHandler.HandleGetLatestArticleVersionInfo)
	r.GET("versions", middleware.ValidateParamMiddleware(&protocol.PageParam{}), articleVersionHandler.HandleListArticleVersions)
	articleVersionRouter := r.Group("/version", middleware.LimitUserPermissionMiddleware("articleVersionService", model.PermissionCreator))
	{
		articleVersionRouter.POST(
			"",
			middleware.RateLimiterMiddleware("createArticleVersion", "userID", 10*time.Second, 1),
			middleware.ValidateBodyMiddleware(&protocol.CreateArticleVersionBody{}),
			articleVersionHandler.HandleCreateArticleVersion,
		)
		articleVersionNumberRouter := articleVersionRouter.Group("v:version", middleware.ValidateURIMiddleware(&protocol.ArticleVersionURI{}))
		{
			articleVersionNumberRouter.GET("", articleVersionHandler.HandleGetArticleVersionInfo)
		}
	}
}
