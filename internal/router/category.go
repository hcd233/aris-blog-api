package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
)

func initCategoryRouter(r fiber.Router) {
	categoryHandler := handler.NewCategoryHandler()

	categoryRouter := r.Group("/category",
		middleware.JwtMiddleware(),
		middleware.LimitUserPermissionMiddleware("categoryService", model.PermissionCreator))
	{
		categoryRouter.Get("/root", categoryHandler.HandleGetRootCategories)
		categoryRouter.Post("/", middleware.ValidateBodyMiddleware(&protocol.CreateCategoryBody{}), categoryHandler.HandleCreateCategory)

		categoryIDRouter := categoryRouter.Group("/:categoryID", middleware.ValidateURIMiddleware(&protocol.CategoryURI{}))
		{
			categoryIDRouter.Get("/", categoryHandler.HandleGetCategoryInfo)
			categoryIDRouter.Delete("/", categoryHandler.HandleDeleteCategory)
			categoryIDRouter.Patch("/", middleware.ValidateBodyMiddleware(&protocol.UpdateCategoryBody{}), categoryHandler.HandleUpdateCategoryInfo)
			categoryIDRouter.Get("/subCategories", middleware.ValidateParamMiddleware(&protocol.PageParam{}), categoryHandler.HandleListChildrenCategories)
			categoryIDRouter.Get("/subArticles", middleware.ValidateParamMiddleware(&protocol.PageParam{}), categoryHandler.HandleListChildrenArticles)
		}
	}
}
