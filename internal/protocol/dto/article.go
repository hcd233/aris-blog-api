// Package dto 文章DTO
package dto

// CreateArticleRequest 创建文章请求
//
//	author centonhuang
//	update 2025-10-30
type CreateArticleRequest struct {
	Body *CreateArticleBody `json:"body" doc:"Request body containing article information"`
}

// CreateArticleBody 创建文章请求体
//
//	author centonhuang
//	update 2025-10-30
type CreateArticleBody struct {
	Title      string   `json:"title" doc:"Article title"`
	Slug       string   `json:"slug" doc:"URL-friendly article identifier"`
	Tags       []string `json:"tags" doc:"List of tag names for the article"`
	CategoryID uint     `json:"categoryID" doc:"ID of the category this article belongs to"`
}

// CreateArticleResponse 创建文章响应
//
//	author centonhuang
//	update 2025-10-30
type CreateArticleResponse struct {
	Article *Article `json:"article" doc:"Created article information"`
}

// GetArticleInfoRequest 获取文章信息请求
//
//	author centonhuang
//	update 2025-10-30
type GetArticleInfoRequest struct {
	ArticleID uint `path:"articleID" doc:"Unique identifier of the article to retrieve"`
}

// GetArticleInfoResponse 获取文章信息响应
//
//	author centonhuang
//	update 2025-10-30
type GetArticleInfoResponse struct {
	Article *Article `json:"article" doc:"Article information"`
}

// GetArticleInfoBySlugRequest 通过Slug获取文章信息请求
//
//	author centonhuang
//	update 2025-10-30
type GetArticleInfoBySlugRequest struct {
	AuthorName  string `path:"authorName" doc:"Author's username"`
	ArticleSlug string `path:"articleSlug" doc:"URL-friendly article identifier"`
}

// GetArticleInfoBySlugResponse 通过Slug获取文章信息响应
//
//	author centonhuang
//	update 2025-10-30
type GetArticleInfoBySlugResponse struct {
	Article *Article `json:"article" doc:"Article information"`
}

// UpdateArticleRequest 更新文章请求
//
//	author centonhuang
//	update 2025-10-30
type UpdateArticleRequest struct {
	ArticleID uint               `path:"articleID" doc:"Unique identifier of the article to update"`
	Body      *UpdateArticleBody `json:"body" doc:"Request body containing fields to update"`
}

// UpdateArticleBody 更新文章请求体
//
//	author centonhuang
//	update 2025-10-30
type UpdateArticleBody struct {
	Title      string `json:"title,omitempty" doc:"New article title"`
	Slug       string `json:"slug,omitempty" doc:"New URL-friendly article identifier"`
	CategoryID uint   `json:"categoryID,omitempty" doc:"New category ID"`
}

// UpdateArticleResponse 更新文章响应
//
//	author centonhuang
//	update 2025-10-30
type UpdateArticleResponse struct{}

// UpdateArticleStatusRequest 更新文章状态请求
//
//	author centonhuang
//	update 2025-10-30
type UpdateArticleStatusRequest struct {
	ArticleID uint                     `path:"articleID" doc:"Unique identifier of the article"`
	Body      *UpdateArticleStatusBody `json:"body" doc:"Request body containing new status"`
}

// UpdateArticleStatusBody 更新文章状态请求体
//
//	author centonhuang
//	update 2025-10-30
type UpdateArticleStatusBody struct {
	Status string `json:"status" doc:"New article status (draft/publish)"`
}

// UpdateArticleStatusResponse 更新文章状态响应
//
//	author centonhuang
//	update 2025-10-30
type UpdateArticleStatusResponse struct{}

// DeleteArticleRequest 删除文章请求
//
//	author centonhuang
//	update 2025-10-30
type DeleteArticleRequest struct {
	ArticleID uint `path:"articleID" doc:"Unique identifier of the article to delete"`
}

// DeleteArticleResponse 删除文章响应
//
//	author centonhuang
//	update 2025-10-30
type DeleteArticleResponse struct{}

// ListArticlesRequest 列出文章请求
//
//	author centonhuang
//	update 2025-10-30
type ListArticlesRequest struct {
	Page     *int `query:"page" minimum:"1" doc:"Page number (starting from 1)"`
	PageSize *int `query:"pageSize" minimum:"1" maximum:"50" doc:"Number of items per page"`
}

// ListArticlesResponse 列出文章响应
//
//	author centonhuang
//	update 2025-10-30
type ListArticlesResponse struct {
	Articles []*Article `json:"articles" doc:"List of articles"`
	PageInfo *PageInfo  `json:"pageInfo" doc:"Pagination information"`
}

// ArticleVersion 文章版本
//
//	author centonhuang
//	update 2025-10-30
type ArticleVersion struct {
	ArticleVersionID uint   `json:"versionID" doc:"Unique identifier for the version"`
	ArticleID        uint   `json:"articleID" doc:"Article ID"`
	VersionID        uint   `json:"version" doc:"Version number"`
	Content          string `json:"content" doc:"Article content"`
	CreatedAt        string `json:"createdAt" doc:"Creation timestamp"`
	UpdatedAt        string `json:"updatedAt" doc:"Last update timestamp"`
}

// CreateArticleVersionRequest 创建文章版本请求
//
//	author centonhuang
//	update 2025-10-30
type CreateArticleVersionRequest struct {
	ArticleID uint                       `path:"articleID" doc:"Article ID"`
	Body      *CreateArticleVersionBody `json:"body" doc:"Request body containing version content"`
}

// CreateArticleVersionBody 创建文章版本请求体
//
//	author centonhuang
//	update 2025-10-30
type CreateArticleVersionBody struct {
	Content string `json:"content" minLength:"100" maxLength:"20000" doc:"Article content (100-20000 characters)"`
}

// CreateArticleVersionResponse 创建文章版本响应
//
//	author centonhuang
//	update 2025-10-30
type CreateArticleVersionResponse struct {
	ArticleVersion *ArticleVersion `json:"articleVersion" doc:"Created article version"`
}

// GetArticleVersionInfoRequest 获取文章版本信息请求
//
//	author centonhuang
//	update 2025-10-30
type GetArticleVersionInfoRequest struct {
	ArticleID uint `path:"articleID" doc:"Article ID"`
	Version   uint `path:"version" minimum:"1" doc:"Version number"`
}

// GetArticleVersionInfoResponse 获取文章版本信息响应
//
//	author centonhuang
//	update 2025-10-30
type GetArticleVersionInfoResponse struct {
	Version *ArticleVersion `json:"version" doc:"Article version information"`
}

// GetLatestArticleVersionInfoRequest 获取最新文章版本信息请求
//
//	author centonhuang
//	update 2025-10-30
type GetLatestArticleVersionInfoRequest struct {
	ArticleID uint `path:"articleID" doc:"Article ID"`
}

// GetLatestArticleVersionInfoResponse 获取最新文章版本信息响应
//
//	author centonhuang
//	update 2025-10-30
type GetLatestArticleVersionInfoResponse struct {
	Version *ArticleVersion `json:"version" doc:"Latest article version information"`
}

// ListArticleVersionsRequest 列出文章版本请求
//
//	author centonhuang
//	update 2025-10-30
type ListArticleVersionsRequest struct {
	ArticleID uint `path:"articleID" doc:"Article ID"`
	Page      *int `query:"page" minimum:"1" doc:"Page number (starting from 1)"`
	PageSize  *int `query:"pageSize" minimum:"1" maximum:"50" doc:"Number of items per page"`
}

// ListArticleVersionsResponse 列出文章版本响应
//
//	author centonhuang
//	update 2025-10-30
type ListArticleVersionsResponse struct {
	Versions []*ArticleVersion `json:"versions" doc:"List of article versions"`
	PageInfo *PageInfo         `json:"pageInfo" doc:"Pagination information"`
}
