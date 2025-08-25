package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
)

func initTagRouter(r fiber.Router) {
	tagHandler := handler.NewTagHandler()

	tagRouter := r.Group("/tag", middleware.JwtMiddleware())
	{
		tagRouter.Get("/list", middleware.ValidateParamMiddleware(&protocol.PaginateParam{}), tagHandler.HandleListTags)
		tagRouter.Post(
			"/",
			middleware.LimitUserPermissionMiddleware("createTag", model.PermissionCreator),
			middleware.ValidateBodyMiddleware(&protocol.CreateTagBody{}),
			tagHandler.HandleCreateTag,
		)
		tagSlugRouter := tagRouter.Group("/:tagID", middleware.ValidateURIMiddleware(&protocol.TagURI{}))
		{
			tagSlugRouter.Get("/", tagHandler.HandleGetTagInfo)
			tagSlugRouter.Patch(
				"/",
				middleware.LimitUserPermissionMiddleware("updateTag", model.PermissionCreator),
				middleware.ValidateBodyMiddleware(&protocol.UpdateTagBody{}),
				tagHandler.HandleUpdateTag,
			)
			tagSlugRouter.Delete(
				"/",
				middleware.LimitUserPermissionMiddleware("deleteTag", model.PermissionCreator),
				tagHandler.HandleDeleteTag,
			)
		}
	}
}
