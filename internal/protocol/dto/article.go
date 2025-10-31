package dto

import "github.com/hcd233/aris-blog-api/internal/resource/database/model"



// ArticleCreateRequestBody 创建文章请求体
type ArticleCreateRequestBody struct {
	Title      string   `json:"title" doc:"Article title"`
	Slug       string   `json:"slug" doc:"Article slug"`
	CategoryID uint     `json:"categoryID" doc:"Category ID"`
	Tags       []string `json:"tags" doc:"List of tag slugs"`
}

// ArticleCreateRequest 创建文章请求
type ArticleCreateRequest struct {
	Body *ArticleCreateRequestBody `json:"body" doc:"Fields for creating article"`
}

// ArticleCreateResponse 创建文章响应
type ArticleCreateResponse struct {
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

// ArticleGetRequest 获取文章详情请求
type ArticleGetRequest struct {
	ArticlePathParam
}

// ArticleGetResponse 获取文章详情响应
type ArticleGetResponse struct {
	Article *Article `json:"article" doc:"Article details"`
}

// ArticleGetBySlugRequest 通过别名获取文章请求
type ArticleGetBySlugRequest struct {
	ArticleSlugPathParam
}

// ArticleGetBySlugResponse 通过别名获取文章响应
type ArticleGetBySlugResponse struct {
	Article *Article `json:"article" doc:"Article details"`
}

// ArticleUpdateRequestBody 更新文章请求体
type ArticleUpdateRequestBody struct {
	Title      string `json:"title" doc:"New title"`
	Slug       string `json:"slug" doc:"New slug"`
	CategoryID uint   `json:"categoryID" doc:"New category ID"`
}

// ArticleUpdateRequest 更新文章请求
type ArticleUpdateRequest struct {
	ArticlePathParam
	Body *ArticleUpdateRequestBody `json:"body" doc:"Updatable article fields"`
}

// ArticleUpdateStatusRequestBody 更新文章状态请求体
type ArticleUpdateStatusRequestBody struct {
	Status model.ArticleStatus `json:"status" doc:"Article status"`
}

// ArticleUpdateStatusRequest 更新文章状态请求
type ArticleUpdateStatusRequest struct {
	ArticlePathParam
	Body *ArticleUpdateStatusRequestBody `json:"body" doc:"Status field"`
}

// ArticleDeleteRequest 删除文章请求
type ArticleDeleteRequest struct {
	ArticlePathParam
}

// ArticleListRequest 列出文章请求
type ArticleListRequest struct {
	CommonParam
}

// ArticleListResponse 列出文章响应
type ArticleListResponse struct {
	Articles []*Article `json:"articles" doc:"List of articles"`
	PageInfo *PageInfo  `json:"pageInfo" doc:"Pagination information"`
}
