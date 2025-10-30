// Package dto 文章版本DTO
package dto

// ArticleVersion 文章版本
type ArticleVersion struct {
	ArticleVersionID uint   `json:"versionID" doc:"Unique identifier for the article version"`
	ArticleID        uint   `json:"articleID" doc:"ID of the article this version belongs to"`
	VersionID        uint   `json:"version" doc:"Version number"`
	Content          string `json:"content" doc:"Content of the article version"`
	CreatedAt        string `json:"createdAt" doc:"Timestamp when the version was created"`
	UpdatedAt        string `json:"updatedAt" doc:"Timestamp when the version was last updated"`
}

// CreateArticleVersionRequest 创建文章版本请求
//
//	author centonhuang
//	update 2025-01-05 16:42:48
type CreateArticleVersionRequest struct {
	ArticleID uint                    `json:"articleID" path:"articleID" doc:"ID of the article"`
	Body      *CreateArticleVersionBody `json:"body" doc:"Request body containing version content"`
}

// CreateArticleVersionBody 创建文章版本请求体
//
//	author centonhuang
//	update 2025-01-05 16:42:48
type CreateArticleVersionBody struct {
	Content string `json:"content" doc:"Content of the article version" example:"This is the article content..."`
}

// CreateArticleVersionResponse 创建文章版本响应
//
//	author centonhuang
//	update 2025-01-05 16:42:48
type CreateArticleVersionResponse struct {
	ArticleVersion *ArticleVersion `json:"articleVersion" doc:"Created article version information"`
}

// GetArticleVersionInfoRequest 获取文章版本信息请求
//
//	author centonhuang
//	update 2025-01-05 16:42:48
type GetArticleVersionInfoRequest struct {
	ArticleID uint `json:"articleID" path:"articleID" doc:"ID of the article"`
	Version   uint `json:"version" path:"version" doc:"Version number"`
}

// GetArticleVersionInfoResponse 获取文章版本信息响应
//
//	author centonhuang
//	update 2025-01-05 16:42:48
type GetArticleVersionInfoResponse struct {
	Version *ArticleVersion `json:"version" doc:"Article version information"`
}

// GetLatestArticleVersionInfoRequest 获取最新文章版本信息请求
//
//	author centonhuang
//	update 2025-01-05 16:42:48
type GetLatestArticleVersionInfoRequest struct {
	ArticleID uint `json:"articleID" path:"articleID" doc:"ID of the article"`
}

// GetLatestArticleVersionInfoResponse 获取最新文章版本信息响应
//
//	author centonhuang
//	update 2025-01-05 16:42:48
type GetLatestArticleVersionInfoResponse struct {
	Version *ArticleVersion `json:"version" doc:"Latest article version information"`
}

// ListArticleVersionsRequest 列出文章版本请求
//
//	author centonhuang
//	update 2025-01-05 16:42:48
type ListArticleVersionsRequest struct {
	ArticleID uint `json:"articleID" path:"articleID" doc:"ID of the article"`
	Page      int  `query:"page" doc:"Page number (starts from 1)" example:"1"`
	PageSize  int  `query:"pageSize" doc:"Number of items per page" example:"10"`
}

// ListArticleVersionsResponse 列出文章版本响应
//
//	author centonhuang
//	update 2025-01-05 16:42:48
type ListArticleVersionsResponse struct {
	Versions []*ArticleVersion `json:"versions" doc:"List of article versions"`
	PageInfo *PageInfo          `json:"pageInfo" doc:"Pagination information"`
}
