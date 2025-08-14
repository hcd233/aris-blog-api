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

func initAssetRouter(r fiber.Router) {
	assetHandler := handler.NewAssetHandler()

	assetRouter := r.Group("/asset", middleware.JwtMiddleware())
	{
		likeRouter := assetRouter.Group("/like")
		{
			likeRouter.Get("/articles", middleware.ValidateParamMiddleware(&protocol.PageParam{}), assetHandler.HandleListUserLikeArticles)
			likeRouter.Get("/comments", middleware.ValidateParamMiddleware(&protocol.PageParam{}), assetHandler.HandleListUserLikeComments)
			likeRouter.Get("/tags", middleware.ValidateParamMiddleware(&protocol.PageParam{}), assetHandler.HandleListUserLikeTags)
		}
		viewRouter := assetRouter.Group("/view")
		{
			viewRouter.Get("/articles", middleware.ValidateParamMiddleware(&protocol.PageParam{}), assetHandler.HandleListUserViewArticles)
			viewRouter.Delete("/:viewID", middleware.ValidateURIMiddleware(&protocol.ViewURI{}), assetHandler.HandleDeleteUserView)
		}
		objectRouter := assetRouter.Group("/object")
		{
			objectRouter.Get(
				"/images",
				middleware.LimitUserPermissionMiddleware("objectService", model.PermissionCreator),
				assetHandler.HandleListImages,
			)
			imageRouter := objectRouter.Group("/image")
			{
				imageRouter.Post(
					"/",
					middleware.LimitUserPermissionMiddleware("objectService", model.PermissionCreator),
					middleware.RateLimiterMiddleware("uploadImage", constant.CtxKeyUserID, 10*time.Second, 1),
					assetHandler.HandleUploadImage,
				)
				imageIDRouter := imageRouter.Group("/:objectName", middleware.ValidateURIMiddleware(&protocol.ObjectURI{}))
				{
					imageIDRouter.Get("/", middleware.ValidateParamMiddleware(&protocol.ImageParam{}), assetHandler.HandleGetImage)
					imageIDRouter.Delete("/", middleware.LimitUserPermissionMiddleware("objectService", model.PermissionCreator), assetHandler.HandleDeleteImage)
				}
			}
		}
	}
}
