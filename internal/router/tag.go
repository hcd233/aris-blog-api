package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/middleware"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"github.com/hcd233/Aris-blog/internal/service"
)

func initTagRouter(r *gin.RouterGroup) {
	tagService := service.NewTagService()

	r.GET("tags", middleware.ValidateParamMiddleware(&protocol.PageParam{}), tagService.ListTagsHandler)
	tagRouter := r.Group("/tag", middleware.JwtMiddleware())
	{
		tagRouter.GET("", middleware.ValidateParamMiddleware(&protocol.QueryParam{}), tagService.QueryTagHandler)
		tagRouter.POST(
			"",
			middleware.LimitUserPermissionMiddleware(model.PermissionCreator),
			middleware.ValidateBodyMiddleware(&protocol.CreateTagBody{}),
			tagService.CreateTagHandler,
		)
		tagSlugRouter := tagRouter.Group("/:tagSlug", middleware.ValidateURIMiddleware(&protocol.TagURI{}))
		{
			tagSlugRouter.GET("", tagService.GetTagInfoHandler)
			tagSlugRouter.PUT(
				"",
				middleware.LimitUserPermissionMiddleware(model.PermissionCreator),
				middleware.ValidateBodyMiddleware(&protocol.UpdateTagBody{}),
				tagService.UpdateTagHandler,
			)
			tagSlugRouter.DELETE(
				"",
				middleware.LimitUserPermissionMiddleware(model.PermissionCreator),
				tagService.DeleteTagHandler,
			)
		}
	}
}

func initUserTagRouter(r *gin.RouterGroup) {
	tagService := service.NewTagService()

	r.GET("tags", middleware.ValidateParamMiddleware(&protocol.PageParam{}), tagService.ListUserTagsHandler)
	tagRouter := r.Group("/tag")
	{
		tagRouter.GET("", middleware.ValidateParamMiddleware(&protocol.QueryParam{}), tagService.QueryUserTagHandler)
	}
}
