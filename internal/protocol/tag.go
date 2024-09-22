package protocol

// TagURI 标签路径参数
//
//	@author centonhuang
//	@update 2024-09-22 03:20:00
type TagURI struct {
	TagSlug string `uri:"tagSlug" binding:"required"`
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
