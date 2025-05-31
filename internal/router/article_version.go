package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
)

func initArticleVersionRouter(r *gin.RouterGroup) {
	articleVersionHandler := handler.NewArticleVersionHandler()

	r.GET("version/latest", articleVersionHandler.HandleGetLatestArticleVersionInfo)
	articleVersionRouter := r.Group("/version", middleware.LimitUserPermissionMiddleware("articleVersionService", model.PermissionCreator))
	{
		articleVersionRouter.GET("list", middleware.ValidateParamMiddleware(&protocol.PageParam{}), articleVersionHandler.HandleListArticleVersions)

		articleVersionRouter.POST(
			"",
			middleware.RateLimiterMiddleware("createArticleVersion", constant.CtxKeyUserID, 10*time.Second, 1),
			middleware.ValidateBodyMiddleware(&protocol.CreateArticleVersionBody{}),
			articleVersionHandler.HandleCreateArticleVersion,
		)
		articleVersionRouter.GET("v:version", middleware.ValidateURIMiddleware(&protocol.ArticleVersionURI{}), articleVersionHandler.HandleGetArticleVersionInfo)
	}
}
