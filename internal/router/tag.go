package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/handler"
	"github.com/hcd233/Aris-blog/internal/middleware"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
)

func initTagRouter(r *gin.RouterGroup) {
	tagHandler := handler.NewTagHandler()

	r.GET("tags", middleware.ValidateParamMiddleware(&protocol.PageParam{}), tagHandler.HandleListTags)
	tagRouter := r.Group("/tag", middleware.JwtMiddleware())
	{
		tagRouter.GET("", middleware.ValidateParamMiddleware(&protocol.QueryParam{}), tagHandler.HandleQueryTag)
		tagRouter.POST(
			"",
			middleware.LimitUserPermissionMiddleware("createTag", model.PermissionCreator),
			middleware.ValidateBodyMiddleware(&protocol.CreateTagBody{}),
			tagHandler.HandleCreateTag,
		)
		tagSlugRouter := tagRouter.Group("/:tagSlug", middleware.ValidateURIMiddleware(&protocol.TagURI{}))
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

func initUserTagRouter(r *gin.RouterGroup) {
	tagHandler := handler.NewTagHandler()

	r.GET("tags", middleware.ValidateParamMiddleware(&protocol.PageParam{}), tagHandler.HandleListUserTags)
	tagRouter := r.Group("/tag")
	{
		tagRouter.GET("", middleware.ValidateParamMiddleware(&protocol.QueryParam{}), tagHandler.HandleQueryUserTag)
	}
}
