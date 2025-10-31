package dto

import "github.com/hcd233/aris-blog-api/internal/resource/database/model"

// CreateArticleRequestBody 创建文章请求体
type CreateArticleRequestBody struct {
	Title      string   `json:"title" doc:"Article title"`
	Slug       string   `json:"slug" doc:"Article slug"`
	CategoryID uint     `json:"categoryID" doc:"Category ID"`
	Tags       []string `json:"tags" doc:"List of tag slugs"`
}

// CreateArticleRequest 创建文章请求
type CreateArticleRequest struct {
	Body *CreateArticleRequestBody `json:"body" doc:"Fields for creating article"`
}

// CreateArticleResponse 创建文章响应
type CreateArticleResponse struct {
	Article *Article `json:"article" doc:"Article details"`
}

// ArticlePathParam 文章路径参数
type ArticlePathParam struct {
	ArticleID uint `path:"articleID" doc:"Article ID"`
}

// ArticleSlugPathParam 文章别名路径参数
type ArticleSlugPathParam struct {
	AuthorName  string `path:"authorName" doc:"Author name"`
	ArticleSlug string `path:"articleSlug" doc:"Article slug"`
}

// GetArticleRequest 获取文章详情请求
type GetArticleRequest struct {
	ArticlePathParam
}

// GetArticleResponse 获取文章详情响应
type GetArticleResponse struct {
	Article *Article `json:"article" doc:"Article details"`
}

// GetArticleBySlugRequest 通过别名获取文章请求
type GetArticleBySlugRequest struct {
	ArticleSlugPathParam
}

// GetArticleBySlugResponse 通过别名获取文章响应
type GetArticleBySlugResponse struct {
	Article *Article `json:"article" doc:"Article details"`
}

// UpdateArticleRequestBody 更新文章请求体
type UpdateArticleRequestBody struct {
	Title      string `json:"title" doc:"New title"`
	Slug       string `json:"slug" doc:"New slug"`
	CategoryID uint   `json:"categoryID" doc:"New category ID"`
}

// UpdateArticleRequest 更新文章请求
type UpdateArticleRequest struct {
	ArticlePathParam
	Body *UpdateArticleRequestBody `json:"body" doc:"Updatable article fields"`
}

// UpdateArticleStatusRequestBody 更新文章状态请求体
type UpdateArticleStatusRequestBody struct {
	Status model.ArticleStatus `json:"status" doc:"Article status"`
}

// UpdateArticleStatusRequest 更新文章状态请求
type UpdateArticleStatusRequest struct {
	ArticlePathParam
	Body *UpdateArticleStatusRequestBody `json:"body" doc:"Status field"`
}

// DeleteArticleRequest 删除文章请求
type DeleteArticleRequest struct {
	ArticlePathParam
}

// ListArticleRequest 列出文章请求
type ListArticleRequest struct {
	CommonParam
}

// ListArticleResponse 列出文章响应
type ListArticleResponse struct {
	Articles []*Article `json:"articles" doc:"List of articles"`
	PageInfo *PageInfo  `json:"pageInfo" doc:"Pagination information"`
}
