package protocol

// ArticleURI 文章路径参数
//
//	@author centonhuang
//	@update 2024-09-21 06:13:15
type ArticleURI struct {
	UserURI
	ArticleSlug string `uri:"articleSlug" binding:"required"`
}

// CreateArticleBody 创建文章请求体
//
//	@author centonhuang
//	@update 2024-09-21 09:59:55
type CreateArticleBody struct {
	Title string `json:"title" binding:"required"`
	Slug  string `json:"slug"`
	// TODO support tags
}

// UpdateArticleBody 更新文章请求体
//
//	@author centonhuang
//	@update 2024-09-22 03:56:09
type UpdateArticleBody struct {
	Title string `json:"title"`
	Slug  string `json:"slug"`
}
