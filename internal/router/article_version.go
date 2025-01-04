package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/handler"
	"github.com/hcd233/Aris-blog/internal/middleware"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
)

func initArticleVersionRouter(r *gin.RouterGroup) {
	articleVersionHandler := handler.NewArticleVersionHandler()

	r.GET("versions", middleware.ValidateParamMiddleware(&protocol.PageParam{}), articleVersionHandler.HandleListArticleVersions)
	articleVersionRouter := r.Group("/version", middleware.LimitUserPermissionMiddleware(model.PermissionCreator))
	{
		articleVersionRouter.POST(
			"",
			middleware.RateLimiterMiddleware(10*time.Second, 1, "createArticleVersion", "userID", protocol.CodeCreateArticleVersionRateLimitError),
			middleware.ValidateBodyMiddleware(&protocol.CreateArticleVersionBody{}),
			articleVersionHandler.HandleCreateArticleVersion,
		)
		articleVersionRouter.GET("latest", articleVersionHandler.HandleGetLatestArticleVersionInfo)
		articleVersionNumberRouter := articleVersionRouter.Group("v:version", middleware.ValidateURIMiddleware(&protocol.ArticleVersionURI{}))
		{
			articleVersionNumberRouter.GET("", articleVersionHandler.HandleGetArticleVersionInfo)
		}
	}
}
