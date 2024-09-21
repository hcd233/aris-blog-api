package protocol

// UserURI 用户路径参数
//
//	@author centonhuang
//	@update 2024-09-18 02:50:19
type UserURI struct {
	UserName string `uri:"userName" binding:"required"`
}

// ArticleURI 文章路径参数
//
//	@author centonhuang
//	@update 2024-09-21 06:13:15
type ArticleURI struct {
	UserURI
	ArticleSlug string `uri:"articleSlug" binding:"required"`
}

// PageParams 列表参数
//
//	@author centonhuang
//	@update 2024-09-21 09:00:57
type PageParams struct {
	Limit  int `form:"limit" binding:"required,min=1,max=50"`
	Offset int `form:"offset" binding:"gte=0"`
}

// QueryParams 查询参数
//
//	@author centonhuang
//	@update 2024-09-18 02:56:39
type QueryParams struct {
	PageParams
	Query string `form:"query" binding:"required"`
}

// UpdateUserBody 更新用户请求体
//
//	@author centonhuang
//	@update 2024-09-18 02:39:31
type UpdateUserBody struct {
	UserName string `json:"userName" binding:"required"`
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
