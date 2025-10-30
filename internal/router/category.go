package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
)

func initCategoryRouter(categoryGroup *huma.Group) {
	categoryHandler := handler.NewCategoryHandler()

	categoryGroup.UseMiddleware(middleware.JwtMiddlewareForHuma())

	// 创建分类
	huma.Register(categoryGroup, huma.Operation{
		OperationID: "createCategory",
		Method:      http.MethodPost,
		Path:        "/",
		Summary:     "CreateCategory",
		Description: "Create a new category",
		Tags:        []string{"category"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, categoryHandler.HandleCreateCategory)

	// 获取根分类
	huma.Register(categoryGroup, huma.Operation{
		OperationID: "getRootCategory",
		Method:      http.MethodGet,
		Path:        "/root",
		Summary:     "GetRootCategory",
		Description: "Get the root category for the current user",
		Tags:        []string{"category"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, categoryHandler.HandleGetRootCategories)

	// 获取分类信息
	huma.Register(categoryGroup, huma.Operation{
		OperationID: "getCategoryInfo",
		Method:      http.MethodGet,
		Path:        "/{categoryID}",
		Summary:     "GetCategoryInfo",
		Description: "Get detailed information about a specific category by ID",
		Tags:        []string{"category"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, categoryHandler.HandleGetCategoryInfo)

	// 更新分类
	huma.Register(categoryGroup, huma.Operation{
		OperationID: "updateCategory",
		Method:      http.MethodPatch,
		Path:        "/{categoryID}",
		Summary:     "UpdateCategory",
		Description: "Update an existing category",
		Tags:        []string{"category"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, categoryHandler.HandleUpdateCategoryInfo)

	// 删除分类
	huma.Register(categoryGroup, huma.Operation{
		OperationID: "deleteCategory",
		Method:      http.MethodDelete,
		Path:        "/{categoryID}",
		Summary:     "DeleteCategory",
		Description: "Delete a category",
		Tags:        []string{"category"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, categoryHandler.HandleDeleteCategory)

	// 列出子分类
	huma.Register(categoryGroup, huma.Operation{
		OperationID: "listChildrenCategories",
		Method:      http.MethodGet,
		Path:        "/{categoryID}/subCategories",
		Summary:     "ListChildrenCategories",
		Description: "Get a paginated list of subcategories for a specific category",
		Tags:        []string{"category"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, categoryHandler.HandleListChildrenCategories)

	// 列出子文章
	huma.Register(categoryGroup, huma.Operation{
		OperationID: "listChildrenArticles",
		Method:      http.MethodGet,
		Path:        "/{categoryID}/subArticles",
		Summary:     "ListChildrenArticles",
		Description: "Get a paginated list of articles in a specific category",
		Tags:        []string{"category"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, categoryHandler.HandleListChildrenArticles)
}
