// Package dto 分类DTO
package dto

// Category 分类
//
//	author centonhuang
//	update 2025-10-30
type Category struct {
	CategoryID uint   `json:"categoryID" doc:"Unique identifier for the category"`
	Name       string `json:"name" doc:"Category name"`
	ParentID   uint   `json:"parentID,omitempty" doc:"ID of the parent category"`
	CreatedAt  string `json:"createdAt,omitempty" doc:"Timestamp when the category was created"`
	UpdatedAt  string `json:"updatedAt,omitempty" doc:"Timestamp when the category was last updated"`
}

// CreateCategoryRequest 创建分类请求
//
//	author centonhuang
//	update 2025-10-30
type CreateCategoryRequest struct {
	Body *CreateCategoryBody `json:"body" doc:"Request body containing category information"`
}

// CreateCategoryBody 创建分类请求体
//
//	author centonhuang
//	update 2025-10-30
type CreateCategoryBody struct {
	ParentID uint   `json:"parentID,omitempty" doc:"ID of the parent category"`
	Name     string `json:"name" doc:"Category name"`
}

// CreateCategoryResponse 创建分类响应
//
//	author centonhuang
//	update 2025-10-30
type CreateCategoryResponse struct {
	Category *Category `json:"category" doc:"Created category information"`
}

// GetCategoryInfoRequest 获取分类信息请求
//
//	author centonhuang
//	update 2025-10-30
type GetCategoryInfoRequest struct {
	CategoryID uint `path:"categoryID" doc:"Unique identifier of the category to retrieve"`
}

// GetCategoryInfoResponse 获取分类信息响应
//
//	author centonhuang
//	update 2025-10-30
type GetCategoryInfoResponse struct {
	Category *Category `json:"category" doc:"Category information"`
}

// GetRootCategoryRequest 获取根分类请求
//
//	author centonhuang
//	update 2025-10-30
type GetRootCategoryRequest struct{}

// GetRootCategoryResponse 获取根分类响应
//
//	author centonhuang
//	update 2025-10-30
type GetRootCategoryResponse struct {
	Category *Category `json:"category" doc:"Root category information"`
}

// UpdateCategoryRequest 更新分类请求
//
//	author centonhuang
//	update 2025-10-30
type UpdateCategoryRequest struct {
	CategoryID uint                 `path:"categoryID" doc:"Unique identifier of the category to update"`
	Body       *UpdateCategoryBody `json:"body" doc:"Request body containing fields to update"`
}

// UpdateCategoryBody 更新分类请求体
//
//	author centonhuang
//	update 2025-10-30
type UpdateCategoryBody struct {
	Name     string `json:"name,omitempty" doc:"New category name"`
	ParentID uint   `json:"parentID,omitempty" doc:"New parent category ID"`
}

// UpdateCategoryResponse 更新分类响应
//
//	author centonhuang
//	update 2025-10-30
type UpdateCategoryResponse struct {
	Category *Category `json:"category" doc:"Updated category information"`
}

// DeleteCategoryRequest 删除分类请求
//
//	author centonhuang
//	update 2025-10-30
type DeleteCategoryRequest struct {
	CategoryID uint `path:"categoryID" doc:"Unique identifier of the category to delete"`
}

// DeleteCategoryResponse 删除分类响应
//
//	author centonhuang
//	update 2025-10-30
type DeleteCategoryResponse struct{}

// ListChildrenCategoriesRequest 列出子分类请求
//
//	author centonhuang
//	update 2025-10-30
type ListChildrenCategoriesRequest struct {
	CategoryID uint `path:"categoryID" doc:"Parent category ID"`
	Page       *int `query:"page" minimum:"1" doc:"Page number (starting from 1)"`
	PageSize   *int `query:"pageSize" minimum:"1" maximum:"50" doc:"Number of items per page"`
}

// ListChildrenCategoriesResponse 列出子分类响应
//
//	author centonhuang
//	update 2025-10-30
type ListChildrenCategoriesResponse struct {
	Categories []*Category `json:"categories" doc:"List of child categories"`
	PageInfo   *PageInfo   `json:"pageInfo" doc:"Pagination information"`
}

// Article 文章（简略）
//
//	author centonhuang
//	update 2025-10-30
type Article struct {
	ArticleID   uint      `json:"articleID" doc:"Unique identifier for the article"`
	Title       string    `json:"title" doc:"Article title"`
	Slug        string    `json:"slug" doc:"URL-friendly article identifier"`
	Status      string    `json:"status" doc:"Article status (draft/publish)"`
	User        *User     `json:"user" doc:"Article author"`
	Category    *Category `json:"category" doc:"Article category"`
	CreatedAt   string    `json:"createdAt" doc:"Creation timestamp"`
	UpdatedAt   string    `json:"updatedAt" doc:"Last update timestamp"`
	PublishedAt string    `json:"publishedAt" doc:"Publication timestamp"`
	Likes       uint      `json:"likes" doc:"Number of likes"`
	Views       uint      `json:"views" doc:"Number of views"`
	Tags        []*Tag    `json:"tags" doc:"Associated tags"`
	Comments    int       `json:"comments" doc:"Number of comments"`
}

// ListChildrenArticlesRequest 列出子文章请求
//
//	author centonhuang
//	update 2025-10-30
type ListChildrenArticlesRequest struct {
	CategoryID uint `path:"categoryID" doc:"Category ID"`
	Page       *int `query:"page" minimum:"1" doc:"Page number (starting from 1)"`
	PageSize   *int `query:"pageSize" minimum:"1" maximum:"50" doc:"Number of items per page"`
}

// ListChildrenArticlesResponse 列出子文章响应
//
//	author centonhuang
//	update 2025-10-30
type ListChildrenArticlesResponse struct {
	Articles []*Article `json:"articles" doc:"List of articles in the category"`
	PageInfo *PageInfo  `json:"pageInfo" doc:"Pagination information"`
}
