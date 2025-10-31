package dto

// ArticleVersion 文章版本信息
//
//	author centonhuang
//	update 2025-10-31 05:38:00
type ArticleVersion struct {
	ArticleVersionID uint   `json:"versionID" doc:"Version ID"`
	ArticleID        uint   `json:"articleID" doc:"Article ID"`
	VersionID        uint   `json:"version" doc:"Version number"`
	Content          string `json:"content" doc:"Version content"`
	CreatedAt        string `json:"createdAt" doc:"Creation timestamp"`
	UpdatedAt        string `json:"updatedAt" doc:"Update timestamp"`
}

// ArticleVersionArticlePathParam 文章路径参数
type ArticleVersionArticlePathParam struct {
	ArticleID uint `path:"articleID" doc:"Article ID"`
}

// ArticleVersionPathParam 文章版本路径参数
type ArticleVersionPathParam struct {
	ArticleVersionArticlePathParam
	Version uint `path:"version" doc:"Version number"`
}

// ArticleVersionCreateRequestBody 创建文章版本请求体
type ArticleVersionCreateRequestBody struct {
	Content string `json:"content" doc:"Version content"`
}

// ArticleVersionCreateRequest 创建文章版本请求
type ArticleVersionCreateRequest struct {
	ArticleVersionArticlePathParam
	Body *ArticleVersionCreateRequestBody `json:"body" doc:"Fields for creating article version"`
}

// ArticleVersionCreateResponse 创建文章版本响应
type ArticleVersionCreateResponse struct {
	ArticleVersion *ArticleVersion `json:"articleVersion" doc:"Article version details"`
}

// ArticleVersionGetRequest 获取文章版本请求
type ArticleVersionGetRequest struct {
	ArticleVersionPathParam
}

// ArticleVersionGetResponse 获取文章版本响应
type ArticleVersionGetResponse struct {
	Version *ArticleVersion `json:"version" doc:"Article version details"`
}

// ArticleVersionGetLatestRequest 获取最新文章版本请求
type ArticleVersionGetLatestRequest struct {
	ArticleVersionArticlePathParam
}

// ArticleVersionGetLatestResponse 获取最新文章版本响应
type ArticleVersionGetLatestResponse struct {
	Version *ArticleVersion `json:"version" doc:"Latest article version"`
}

// ArticleVersionListRequest 列出文章版本请求
type ArticleVersionListRequest struct {
	ArticleVersionArticlePathParam
	PaginationQuery
}

// ArticleVersionListResponse 列出文章版本响应
type ArticleVersionListResponse struct {
	Versions []*ArticleVersion `json:"versions" doc:"List of article versions"`
	PageInfo *PageInfo         `json:"pageInfo" doc:"Pagination information"`
}
