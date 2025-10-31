package dto

// ArticleVersion 文章版本信息
//
//	author centonhuang
//	update 2025-10-31 05:38:00
type ArticleVersion struct {
	ArticleVersionID uint   `json:"versionID" doc:"版本 ID"`
	ArticleID        uint   `json:"articleID" doc:"文章 ID"`
	VersionID        uint   `json:"version" doc:"版本号"`
	Content          string `json:"content" doc:"版本内容"`
	CreatedAt        string `json:"createdAt" doc:"创建时间"`
	UpdatedAt        string `json:"updatedAt" doc:"更新时间"`
}

// ArticleVersionArticlePathParam 文章路径参数
type ArticleVersionArticlePathParam struct {
	ArticleID uint `path:"articleID" doc:"文章 ID"`
}

// ArticleVersionPathParam 文章版本路径参数
type ArticleVersionPathParam struct {
	ArticleVersionArticlePathParam
	Version uint `path:"version" doc:"版本号"`
}

// ArticleVersionCreateRequestBody 创建文章版本请求体
type ArticleVersionCreateRequestBody struct {
	Content string `json:"content" doc:"版本内容"`
}

// ArticleVersionCreateRequest 创建文章版本请求
type ArticleVersionCreateRequest struct {
	ArticleVersionArticlePathParam
	UserID uint                             `json:"-"`
	Body   *ArticleVersionCreateRequestBody `json:"body" doc:"创建版本字段"`
}

// ArticleVersionCreateResponse 创建文章版本响应
type ArticleVersionCreateResponse struct {
	ArticleVersion *ArticleVersion `json:"articleVersion" doc:"文章版本详情"`
}

// ArticleVersionGetRequest 获取文章版本请求
type ArticleVersionGetRequest struct {
	ArticleVersionPathParam
	UserID uint `json:"-"`
}

// ArticleVersionGetResponse 获取文章版本响应
type ArticleVersionGetResponse struct {
	Version *ArticleVersion `json:"version" doc:"文章版本详情"`
}

// ArticleVersionGetLatestRequest 获取最新文章版本请求
type ArticleVersionGetLatestRequest struct {
	ArticleVersionArticlePathParam
	UserID uint `json:"-"`
}

// ArticleVersionGetLatestResponse 获取最新文章版本响应
type ArticleVersionGetLatestResponse struct {
	Version *ArticleVersion `json:"version" doc:"文章最新版本"`
}

// ArticleVersionListRequest 列出文章版本请求
type ArticleVersionListRequest struct {
	ArticleVersionArticlePathParam
	UserID uint `json:"-"`
	PaginationQuery
}

// ArticleVersionListResponse 列出文章版本响应
type ArticleVersionListResponse struct {
	Versions []*ArticleVersion `json:"versions" doc:"文章版本列表"`
	PageInfo *PageInfo         `json:"pageInfo" doc:"分页信息"`
}
