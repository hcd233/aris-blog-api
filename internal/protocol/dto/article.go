package dto

import "github.com/hcd233/aris-blog-api/internal/resource/database/model"

// Article 文章信息
//
//	author centonhuang
//	update 2025-10-31 05:36:00
type Article struct {
	ArticleID   uint      `json:"articleID" doc:"文章 ID"`
	Title       string    `json:"title" doc:"文章标题"`
	Slug        string    `json:"slug" doc:"文章别名"`
	Status      string    `json:"status" doc:"文章状态"`
	User        *User     `json:"user" doc:"作者信息"`
	Category    *Category `json:"category" doc:"分类信息"`
	CreatedAt   string    `json:"createdAt" doc:"创建时间"`
	UpdatedAt   string    `json:"updatedAt" doc:"更新时间"`
	PublishedAt string    `json:"publishedAt" doc:"发布时间"`
	Likes       uint      `json:"likes" doc:"点赞数量"`
	Views       uint      `json:"views" doc:"浏览量"`
	Tags        []*Tag    `json:"tags" doc:"标签列表"`
	Comments    int       `json:"comments" doc:"评论数量"`
}

// ArticleCreateRequestBody 创建文章请求体
type ArticleCreateRequestBody struct {
	Title      string   `json:"title" doc:"文章标题"`
	Slug       string   `json:"slug" doc:"文章别名"`
	CategoryID uint     `json:"categoryID" doc:"分类 ID"`
	Tags       []string `json:"tags" doc:"标签别名列表"`
}

// ArticleCreateRequest 创建文章请求
type ArticleCreateRequest struct {
	UserID uint                      `json:"-"`
	Body   *ArticleCreateRequestBody `json:"body" doc:"创建文章字段"`
}

// ArticleCreateResponse 创建文章响应
type ArticleCreateResponse struct {
	Article *Article `json:"article" doc:"文章详情"`
}

// ArticlePathParam 文章路径参数
type ArticlePathParam struct {
	ArticleID uint `path:"articleID" doc:"文章 ID"`
}

// ArticleSlugPathParam 文章别名路径参数
type ArticleSlugPathParam struct {
	AuthorName  string `path:"authorName" doc:"作者名称"`
	ArticleSlug string `path:"articleSlug" doc:"文章别名"`
}

// ArticleGetRequest 获取文章详情请求
type ArticleGetRequest struct {
	ArticlePathParam
	UserID uint `json:"-"`
}

// ArticleGetResponse 获取文章详情响应
type ArticleGetResponse struct {
	Article *Article `json:"article" doc:"文章详情"`
}

// ArticleGetBySlugRequest 通过别名获取文章请求
type ArticleGetBySlugRequest struct {
	ArticleSlugPathParam
	UserID uint `json:"-"`
}

// ArticleGetBySlugResponse 通过别名获取文章响应
type ArticleGetBySlugResponse struct {
	Article *Article `json:"article" doc:"文章详情"`
}

// ArticleUpdateRequestBody 更新文章请求体
type ArticleUpdateRequestBody struct {
	Title      string `json:"title" doc:"新的标题"`
	Slug       string `json:"slug" doc:"新的别名"`
	CategoryID uint   `json:"categoryID" doc:"新的分类 ID"`
}

// ArticleUpdateRequest 更新文章请求
type ArticleUpdateRequest struct {
	ArticlePathParam
	UserID uint                      `json:"-"`
	Body   *ArticleUpdateRequestBody `json:"body" doc:"可更新的文章字段"`
}

// ArticleUpdateResponse 更新文章响应
type ArticleUpdateResponse struct{}

// ArticleUpdateStatusRequestBody 更新文章状态请求体
type ArticleUpdateStatusRequestBody struct {
	Status model.ArticleStatus `json:"status" doc:"文章状态"`
}

// ArticleUpdateStatusRequest 更新文章状态请求
type ArticleUpdateStatusRequest struct {
	ArticlePathParam
	UserID uint                            `json:"-"`
	Body   *ArticleUpdateStatusRequestBody `json:"body" doc:"状态字段"`
}

// ArticleUpdateStatusResponse 更新文章状态响应
type ArticleUpdateStatusResponse struct{}

// ArticleDeleteRequest 删除文章请求
type ArticleDeleteRequest struct {
	ArticlePathParam
	UserID uint `json:"-"`
}

// ArticleDeleteResponse 删除文章响应
type ArticleDeleteResponse struct{}

// ArticleListRequest 列出文章请求
type ArticleListRequest struct {
	PaginationQuery
}

// ArticleListResponse 列出文章响应
type ArticleListResponse struct {
	Articles []*Article `json:"articles" doc:"文章列表"`
	PageInfo *PageInfo  `json:"pageInfo" doc:"分页信息"`
}
