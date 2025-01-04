package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/handler"
	"github.com/hcd233/Aris-blog/internal/middleware"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
)

func initUserAssetRouter(r *gin.RouterGroup) {
	assetService := handler.NewAssetService()

	assetRouter := r.Group("/asset")
	{
		likeRouter := assetRouter.Group("/like")
		{
			likeRouter.GET("articles", middleware.ValidateParamMiddleware(&protocol.PageParam{}), assetService.ListUserLikeArticlesHandler)
			likeRouter.GET("comments", middleware.ValidateParamMiddleware(&protocol.PageParam{}), assetService.ListUserLikeCommentsHandler)
			likeRouter.GET("tags", middleware.ValidateParamMiddleware(&protocol.PageParam{}), assetService.ListUserLikeTagsHandler)
		}
		viewRouter := assetRouter.Group("/view")
		{
			viewRouter.GET("article", middleware.ValidateParamMiddleware(&protocol.ArticleParam{}), assetService.GetUserViewArticleHandler)
			viewRouter.GET("articles", middleware.ValidateParamMiddleware(&protocol.PageParam{}), assetService.ListUserViewArticlesHandler)
			viewRouter.DELETE(":viewID", middleware.ValidateURIMiddleware(&protocol.ViewURI{}), assetService.DeleteUserViewHandler)
		}
		objectRouter := assetRouter.Group("/object")
		{
			objectRouter.POST(
				"bucket",
				middleware.LimitUserPermissionMiddleware(model.PermissionCreator),
				assetService.CreateBucketHandler,
			)

			objectRouter.GET(
				"images",
				middleware.LimitUserPermissionMiddleware(model.PermissionCreator),
				assetService.ListImagesHandler,
			)
			imageRouter := objectRouter.Group("/image")
			{
				imageRouter.POST(
					"",
					middleware.LimitUserPermissionMiddleware(model.PermissionCreator),
					middleware.RateLimiterMiddleware(10*time.Second, 1, "uploadImage", "userID", protocol.CodeUploadImageRateLimitError),
					assetService.UploadImageHandler,
				)
				imageIDRouter := imageRouter.Group("/:objectName", middleware.ValidateURIMiddleware(&protocol.ObjectURI{}))
				{
					imageIDRouter.GET("", middleware.ValidateParamMiddleware(&protocol.ImageParam{}), assetService.GetImageHandler)
					imageIDRouter.DELETE("", middleware.LimitUserPermissionMiddleware(model.PermissionCreator), assetService.DeleteImageHandler)
				}
			}
		}
	}
}
