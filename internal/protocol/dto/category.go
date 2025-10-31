package dto

// Category 分类信息
//
//	author centonhuang
//	update 2025-10-31 05:34:00
type Category struct {
	CategoryID uint   `json:"categoryID" doc:"分类 ID"`
	Name       string `json:"name" doc:"分类名称"`
	ParentID   uint   `json:"parentID,omitempty" doc:"父级分类 ID"`
	CreatedAt  string `json:"createdAt,omitempty" doc:"创建时间"`
	UpdatedAt  string `json:"updatedAt,omitempty" doc:"更新时间"`
}

// CategoryPathParam 分类路径参数
type CategoryPathParam struct {
	CategoryID uint `path:"categoryID" doc:"分类 ID"`
}

// CategoryCreateRequestBody 创建分类请求体
type CategoryCreateRequestBody struct {
	Name     string `json:"name" doc:"分类名称"`
	ParentID uint   `json:"parentID" doc:"父级分类 ID，没有则为 0"`
}

// CategoryCreateRequest 创建分类请求
type CategoryCreateRequest struct {
	UserID uint                       `json:"-"`
	Body   *CategoryCreateRequestBody `json:"body" doc:"创建分类字段"`
}

// CategoryCreateResponse 创建分类响应
type CategoryCreateResponse struct {
	Category *Category `json:"category" doc:"分类详情"`
}

// CategoryGetRequest 获取分类详情请求
type CategoryGetRequest struct {
	CategoryPathParam
	UserID uint `json:"-"`
}

// CategoryGetResponse 获取分类详情响应
type CategoryGetResponse struct {
	Category *Category `json:"category" doc:"分类详情"`
}

// CategoryGetRootRequest 获取根分类请求
type CategoryGetRootRequest struct {
	UserID uint `json:"-"`
}

// CategoryGetRootResponse 获取根分类响应
type CategoryGetRootResponse struct {
	Category *Category `json:"category" doc:"根分类"`
}

// CategoryUpdateRequestBody 更新分类请求体
type CategoryUpdateRequestBody struct {
	Name     string `json:"name" doc:"新的分类名称"`
	ParentID uint   `json:"parentID" doc:"新的父级分类 ID"`
}

// CategoryUpdateRequest 更新分类请求
type CategoryUpdateRequest struct {
	CategoryPathParam
	UserID uint                       `json:"-"`
	Body   *CategoryUpdateRequestBody `json:"body" doc:"可更新的分类字段"`
}

// CategoryUpdateResponse 更新分类响应
type CategoryUpdateResponse struct {
	Category *Category `json:"category" doc:"更新后的分类"`
}

// CategoryDeleteRequest 删除分类请求
type CategoryDeleteRequest struct {
	CategoryPathParam
	UserID uint `json:"-"`
}

// CategoryDeleteResponse 删除分类响应
type CategoryDeleteResponse struct{}

// CategoryListChildrenCategoriesRequest 列出子分类请求
type CategoryListChildrenCategoriesRequest struct {
	CategoryPathParam
	UserID uint `json:"-"`
	PaginationQuery
}

// CategoryListChildrenCategoriesResponse 列出子分类响应
type CategoryListChildrenCategoriesResponse struct {
	Categories []*Category `json:"categories" doc:"子分类列表"`
	PageInfo   *PageInfo   `json:"pageInfo" doc:"分页信息"`
}

// CategoryListChildrenArticlesRequest 列出分类下文章请求
type CategoryListChildrenArticlesRequest struct {
	CategoryPathParam
	UserID uint `json:"-"`
	PaginationQuery
}

// CategoryListChildrenArticlesResponse 列出分类下文章响应
type CategoryListChildrenArticlesResponse struct {
	Articles []*Article `json:"articles" doc:"分类下文章列表"`
	PageInfo *PageInfo  `json:"pageInfo" doc:"分页信息"`
}
