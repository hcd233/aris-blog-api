package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/handler"
	"github.com/hcd233/Aris-blog/internal/middleware"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
)

func initUserCategoryRouter(r *gin.RouterGroup) {
	categoryHandler := handler.NewCategoryHandler()

	r.GET("rootCategory", middleware.LimitUserPermissionMiddleware("categoryService", model.PermissionCreator), categoryHandler.HandleGetRootCategories)
	categoryRouter := r.Group("/category", middleware.LimitUserPermissionMiddleware("categoryService", model.PermissionCreator))
	{
		categoryRouter.POST("", middleware.ValidateBodyMiddleware(&protocol.CreateCategoryBody{}), categoryHandler.HandleCreateCategory)
	}

	categoryIDRouter := categoryRouter.Group(":categoryID", middleware.ValidateURIMiddleware(&protocol.CategoryURI{}))
	{
		categoryIDRouter.GET("", categoryHandler.HandleGetCategoryInfo)
		categoryIDRouter.DELETE("", categoryHandler.HandleDeleteCategory)
		categoryIDRouter.PUT("", middleware.ValidateBodyMiddleware(&protocol.UpdateCategoryBody{}), categoryHandler.HandleUpdateCategoryInfo)
		categoryIDRouter.GET("subCategories", middleware.ValidateParamMiddleware(&protocol.PageParam{}), categoryHandler.HandleListChildrenCategories)
		categoryIDRouter.GET("subArticles", middleware.ValidateParamMiddleware(&protocol.PageParam{}), categoryHandler.HandleListChildrenArticles)
	}
}
