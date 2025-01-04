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
	articleVersionService := handler.NewArticleVersionService()

	r.GET("versions", middleware.ValidateParamMiddleware(&protocol.PageParam{}), articleVersionService.ListArticleVersionsHandler)
	articleVersionRouter := r.Group("/version", middleware.LimitUserPermissionMiddleware(model.PermissionCreator))
	{
		articleVersionRouter.POST(
			"",
			middleware.RateLimiterMiddleware(10*time.Second, 1, "createArticleVersion", "userID", protocol.CodeCreateArticleVersionRateLimitError),
			middleware.ValidateBodyMiddleware(&protocol.CreateArticleVersionBody{}),
			articleVersionService.CreateArticleVersionHandler,
		)
		articleVersionRouter.GET("latest", articleVersionService.GetLatestArticleVersionInfoHandler)
		articleVersionNumberRouter := articleVersionRouter.Group("v:version", middleware.ValidateURIMiddleware(&protocol.ArticleVersionURI{}))
		{
			articleVersionNumberRouter.GET("", articleVersionService.GetArticleVersionInfoHandler)
		}
	}
}
