// Package dto 文章DTO
package dto

import "github.com/hcd233/aris-blog-api/internal/resource/database/model"

// Article 文章
//
//	author centonhuang
//	update 2025-01-05 15:23:26
type Article struct {
	ArticleID   uint   `json:"articleID" doc:"Unique identifier for the article"`
	Title       string `json:"title" doc:"Title of the article"`
	Slug        string `json:"slug" doc:"URL-friendly identifier for the article"`
	Status      string `json:"status" doc:"Publication status of the article (draft/publish)"`
	User        *User  `json:"user" doc:"Author information"`
	Category    *Category `json:"category" doc:"Category information"`
	CreatedAt   string `json:"createdAt" doc:"Timestamp when the article was created"`
	UpdatedAt   string `json:"updatedAt" doc:"Timestamp when the article was last updated"`
	PublishedAt string `json:"publishedAt,omitempty" doc:"Timestamp when the article was published"`
	Likes       uint   `json:"likes" doc:"Number of likes the article has received"`
	Views       uint   `json:"views" doc:"Number of views the article has received"`
	Tags        []*Tag `json:"tags" doc:"Tags associated with the article"`
	Comments    int    `json:"comments" doc:"Number of comments on the article"`
}

// Category 分类
//
//	author centonhuang
//	update 2025-01-05 13:22:49
type Category struct {
	CategoryID uint   `json:"categoryID" doc:"Unique identifier for the category"`
	Name       string `json:"name" doc:"Name of the category"`
	ParentID   uint   `json:"parentID,omitempty" doc:"Parent category ID if this is a subcategory"`
	CreatedAt  string `json:"createdAt,omitempty" doc:"Timestamp when the category was created"`
	UpdatedAt  string `json:"updatedAt,omitempty" doc:"Timestamp when the category was last updated"`
}

// Tag 标签
//
//	author centonhuang
//	update 2025-01-05 12:05:42
type Tag struct {
	TagID       uint   `json:"tagID" doc:"Unique identifier for the tag"`
	Name        string `json:"name" doc:"Name of the tag"`
	Slug        string `json:"slug" doc:"URL-friendly identifier for the tag"`
	Description string `json:"description,omitempty" doc:"Description of the tag"`
	UserID      uint   `json:"userID,omitempty" doc:"ID of the user who created the tag"`
	CreatedAt   string `json:"createdAt,omitempty" doc:"Timestamp when the tag was created"`
	UpdatedAt   string `json:"updatedAt,omitempty" doc:"Timestamp when the tag was last updated"`
	Likes       uint   `json:"likes,omitempty" doc:"Number of likes the tag has received"`
}

// CreateArticleRequest 创建文章请求
//
//	author centonhuang
//	update 2025-01-05 15:23:26
type CreateArticleRequest struct {
	Body *CreateArticleBody `json:"body" doc:"Request body containing article details"`
}

// CreateArticleBody 创建文章请求体
//
//	author centonhuang
//	update 2025-01-05 15:23:26
type CreateArticleBody struct {
	Title      string   `json:"title" doc:"Title of the article" example:"My First Article"`
	Slug       string   `json:"slug" doc:"URL-friendly identifier for the article" example:"my-first-article"`
	Tags       []string `json:"tags" doc:"List of tag names or slugs" example:"go,programming"`
	CategoryID uint     `json:"categoryID" doc:"ID of the category this article belongs to" example:"1"`
}

// CreateArticleResponse 创建文章响应
//
//	author centonhuang
//	update 2025-01-05 15:23:26
type CreateArticleResponse struct {
	Article *Article `json:"article" doc:"Created article information"`
}

// GetArticleInfoRequest 获取文章信息请求
//
//	author centonhuang
//	update 2025-01-05 15:23:26
type GetArticleInfoRequest struct {
	ArticleID uint `json:"articleID" path:"articleID" doc:"Unique identifier of the article to retrieve"`
}

// GetArticleInfoResponse 获取文章信息响应
//
//	author centonhuang
//	update 2025-01-05 15:23:26
type GetArticleInfoResponse struct {
	Article *Article `json:"article" doc:"Article information"`
}

// GetArticleInfoBySlugRequest 通过slug获取文章信息请求
//
//	author centonhuang
//	update 2025-01-19 15:23:26
type GetArticleInfoBySlugRequest struct {
	AuthorName  string `json:"authorName" path:"authorName" doc:"Name of the article author"`
	ArticleSlug string `json:"articleSlug" path:"articleSlug" doc:"URL-friendly identifier of the article"`
}

// GetArticleInfoBySlugResponse 通过slug获取文章信息响应
//
//	author centonhuang
//	update 2025-01-19 15:23:26
type GetArticleInfoBySlugResponse struct {
	Article *Article `json:"article" doc:"Article information"`
}

// UpdateArticleRequest 更新文章请求
//
//	author centonhuang
//	update 2025-01-05 15:23:26
type UpdateArticleRequest struct {
	ArticleID uint                 `json:"articleID" path:"articleID" doc:"Unique identifier of the article to update"`
	Body      *UpdateArticleBody   `json:"body" doc:"Request body containing fields to update"`
}

// UpdateArticleBody 更新文章请求体
//
//	author centonhuang
//	update 2025-01-05 15:23:26
type UpdateArticleBody struct {
	Title      string `json:"title,omitempty" doc:"New title for the article"`
	Slug       string `json:"slug,omitempty" doc:"New URL-friendly identifier for the article"`
	CategoryID uint   `json:"categoryID,omitempty" doc:"New category ID for the article"`
}

// UpdateArticleResponse 更新文章响应
//
//	author centonhuang
//	update 2025-01-05 15:23:26
type UpdateArticleResponse struct{}

// UpdateArticleStatusRequest 更新文章状态请求
//
//	author centonhuang
//	update 2025-01-05 15:23:26
type UpdateArticleStatusRequest struct {
	ArticleID uint                `json:"articleID" path:"articleID" doc:"Unique identifier of the article"`
	Body      *UpdateArticleStatusBody `json:"body" doc:"Request body containing the new status"`
}

// UpdateArticleStatusBody 更新文章状态请求体
//
//	author centonhuang
//	update 2025-01-05 15:23:26
type UpdateArticleStatusBody struct {
	Status model.ArticleStatus `json:"status" doc:"New status for the article (draft/publish)" example:"publish"`
}

// UpdateArticleStatusResponse 更新文章状态响应
//
//	author centonhuang
//	update 2025-01-05 15:23:26
type UpdateArticleStatusResponse struct{}

// DeleteArticleRequest 删除文章请求
//
//	author centonhuang
//	update 2025-01-05 15:23:26
type DeleteArticleRequest struct {
	ArticleID uint `json:"articleID" path:"articleID" doc:"Unique identifier of the article to delete"`
}

// DeleteArticleResponse 删除文章响应
//
//	author centonhuang
//	update 2025-01-05 15:23:26
type DeleteArticleResponse struct{}

// ListArticlesRequest 列出文章请求
//
//	author centonhuang
//	update 2025-01-05 15:23:26
type ListArticlesRequest struct {
	Page     int `query:"page" doc:"Page number (starts from 1)" example:"1"`
	PageSize int `query:"pageSize" doc:"Number of items per page" example:"10"`
}

// ListArticlesResponse 列出文章响应
//
//	author centonhuang
//	update 2025-01-05 15:23:26
type ListArticlesResponse struct {
	Articles []*Article `json:"articles" doc:"List of articles"`
	PageInfo *PageInfo  `json:"pageInfo" doc:"Pagination information"`
}

// PageInfo 分页信息
//
//	author centonhuang
//	update 2025-01-05 12:26:07
type PageInfo struct {
	Page     int   `json:"page" doc:"Current page number"`
	PageSize int   `json:"pageSize" doc:"Number of items per page"`
	Total    int64 `json:"total" doc:"Total number of items"`
}
