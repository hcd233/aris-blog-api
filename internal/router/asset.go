package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
)

func initUserAssetRouter(r *gin.RouterGroup) {
	assetHandler := handler.NewAssetHandler()

	assetRouter := r.Group("/asset")
	{
		likeRouter := assetRouter.Group("/like")
		{
			likeRouter.GET("articles", middleware.ValidateParamMiddleware(&protocol.PageParam{}), assetHandler.HandleListUserLikeArticles)
			likeRouter.GET("comments", middleware.ValidateParamMiddleware(&protocol.PageParam{}), assetHandler.HandleListUserLikeComments)
			likeRouter.GET("tags", middleware.ValidateParamMiddleware(&protocol.PageParam{}), assetHandler.HandleListUserLikeTags)
		}
		viewRouter := assetRouter.Group("/view")
		{
			viewRouter.GET("articles", middleware.ValidateParamMiddleware(&protocol.PageParam{}), assetHandler.HandleListUserViewArticles)
			viewRouter.DELETE(":viewID", middleware.ValidateURIMiddleware(&protocol.ViewURI{}), assetHandler.HandleDeleteUserView)
		}
		objectRouter := assetRouter.Group("/object")
		{
			objectRouter.POST(
				"bucket",
				middleware.LimitUserPermissionMiddleware("objectService", model.PermissionCreator),
				assetHandler.HandleCreateBucket,
			)

			objectRouter.GET(
				"images",
				middleware.LimitUserPermissionMiddleware("objectService", model.PermissionCreator),
				assetHandler.HandleListImages,
			)
			imageRouter := objectRouter.Group("/image")
			{
				imageRouter.POST(
					"",
					middleware.LimitUserPermissionMiddleware("objectService", model.PermissionCreator),
					middleware.RateLimiterMiddleware("uploadImage", "userID", 10*time.Second, 1),
					assetHandler.HandleUploadImage,
				)
				imageIDRouter := imageRouter.Group("/:objectName", middleware.ValidateURIMiddleware(&protocol.ObjectURI{}))
				{
					imageIDRouter.GET("", middleware.ValidateParamMiddleware(&protocol.ImageParam{}), assetHandler.HandleGetImage)
					imageIDRouter.DELETE("", middleware.LimitUserPermissionMiddleware("objectService", model.PermissionCreator), assetHandler.HandleDeleteImage)
				}
			}
		}
	}
}
