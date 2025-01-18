package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
)

func initTagRouter(r *gin.RouterGroup) {
	tagHandler := handler.NewTagHandler()

	r.GET("tags", middleware.ValidateParamMiddleware(&protocol.PageParam{}), tagHandler.HandleListTags)
	tagRouter := r.Group("/tag", middleware.JwtMiddleware())
	{
		tagRouter.POST(
			"",
			middleware.LimitUserPermissionMiddleware("createTag", model.PermissionCreator),
			middleware.ValidateBodyMiddleware(&protocol.CreateTagBody{}),
			tagHandler.HandleCreateTag,
		)
		tagSlugRouter := tagRouter.Group("/:tagID", middleware.ValidateURIMiddleware(&protocol.TagURI{}))
		{
			tagSlugRouter.GET("", tagHandler.HandleGetTagInfo)
			tagSlugRouter.PUT(
				"",
				middleware.LimitUserPermissionMiddleware("updateTag", model.PermissionCreator),
				middleware.ValidateBodyMiddleware(&protocol.UpdateTagBody{}),
				tagHandler.HandleUpdateTag,
			)
			tagSlugRouter.DELETE(
				"",
				middleware.LimitUserPermissionMiddleware("deleteTag", model.PermissionCreator),
				tagHandler.HandleDeleteTag,
			)
		}
	}
}

