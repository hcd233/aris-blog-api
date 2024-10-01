package protocol

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

// UpdateArticleBody 更新文章请求体
//
//	@author centonhuang
//	@update 2024-09-22 03:56:09
type UpdateArticleBody struct {
	Title string `json:"title"`
	Slug  string `json:"slug"`
}

// UpdateTagBody 更新标签请求体
//
//	@author centonhuang
//	@update 2024-09-22 03:20:00
type UpdateTagBody struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Description string `json:"description"`
}

// CreateTagBody 创建标签请求体
//
//	@author centonhuang
//	@update 2024-09-22 03:20:00
type CreateTagBody struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Description string `json:"description"`
}

// CreateCategoryBody 创建分类请求体
//
//	@author centonhuang
//	@update 2024-09-28 07:02:11
type CreateCategoryBody struct {
	ParentID uint   `json:"parentID" binding:"omitempty"`
	Name     string `json:"name" binding:"required"`
}
