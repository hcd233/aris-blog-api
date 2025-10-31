package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
)

func initCategoryRouter(v1Group *huma.Group) {
	categoryHandler := handler.NewCategoryHandler()

	categoryGroup := huma.NewGroup(v1Group, "/category")
	categoryGroup.UseMiddleware(middleware.JwtMiddlewareForHuma())
	categoryGroup.UseMiddleware(middleware.LimitUserPermissionMiddlewareForHuma("categoryService", model.PermissionCreator))

	huma.Register(categoryGroup, huma.Operation{
		OperationID: "listCategories",
		Method:      http.MethodGet,
		Path:        "/{categoryID}/subCategories",
		Summary:     "ListChildrenCategories",
		Description: "List sub categories under a specific category",
		Tags:        []string{"category"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, categoryHandler.HandleListChildrenCategories)

	huma.Register(categoryGroup, huma.Operation{
		OperationID: "listCategoryArticles",
		Method:      http.MethodGet,
		Path:        "/{categoryID}/subArticles",
		Summary:     "ListChildrenArticles",
		Description: "List articles under a specific category",
		Tags:        []string{"category"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, categoryHandler.HandleListChildrenArticles)

	huma.Register(categoryGroup, huma.Operation{
		OperationID: "createCategory",
		Method:      http.MethodPost,
		Path:        "/",
		Summary:     "CreateCategory",
		Description: "Create a new category",
		Tags:        []string{"category"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, categoryHandler.HandleCreateCategory)

	huma.Register(categoryGroup, huma.Operation{
		OperationID: "getCategoryInfo",
		Method:      http.MethodGet,
		Path:        "/{categoryID}",
		Summary:     "GetCategoryInfo",
		Description: "Get category detail by ID",
		Tags:        []string{"category"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, categoryHandler.HandleGetCategoryInfo)

	huma.Register(categoryGroup, huma.Operation{
		OperationID: "updateCategoryInfo",
		Method:      http.MethodPatch,
		Path:        "/{categoryID}",
		Summary:     "UpdateCategoryInfo",
		Description: "Update category info",
		Tags:        []string{"category"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, categoryHandler.HandleUpdateCategoryInfo)

	huma.Register(categoryGroup, huma.Operation{
		OperationID: "deleteCategory",
		Method:      http.MethodDelete,
		Path:        "/{categoryID}",
		Summary:     "DeleteCategory",
		Description: "Delete category",
		Tags:        []string{"category"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, categoryHandler.HandleDeleteCategory)

	huma.Register(categoryGroup, huma.Operation{
		OperationID: "getRootCategory",
		Method:      http.MethodGet,
		Path:        "/root",
		Summary:     "GetRootCategory",
		Description: "Get root category for current user",
		Tags:        []string{"category"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, categoryHandler.HandleGetRootCategories)
}
