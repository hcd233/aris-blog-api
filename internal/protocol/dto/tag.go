package dto

// TagPathParam 标签路径参数
type TagPathParam struct {
	TagID uint `path:"tagID" doc:"Tag ID"`
}

// CreateTagRequestBody 创建标签请求体
type CreateTagRequestBody struct {
	Name        string `json:"name" doc:"Tag name"`
	Slug        string `json:"slug" doc:"Tag slug"`
	Description string `json:"description" doc:"Tag description"`
}

// CreateTagRequest 创建标签请求
type CreateTagRequest struct {
	Body *CreateTagRequestBody `json:"body" doc:"Fields for creating tag"`
}

// CreateTagResponse 创建标签响应
type CreateTagResponse struct {
	Tag *Tag `json:"tag" doc:"Successfully created tag"`
}

// GetTagRequest 获取标签请求
type GetTagRequest struct {
	TagPathParam
}

// GetTagResponse 获取标签响应
type GetTagResponse struct {
	Tag *Tag `json:"tag" doc:"Tag details"`
}

// UpdateTagRequestBody 更新标签请求体
type UpdateTagRequestBody struct {
	Name        string `json:"name" doc:"Tag name"`
	Slug        string `json:"slug" doc:"Tag slug"`
	Description string `json:"description" doc:"Tag description"`
}

// UpdateTagRequest 更新标签请求
type UpdateTagRequest struct {
	TagPathParam
	Body *UpdateTagRequestBody `json:"body" doc:"Updatable tag fields"`
}

// DeleteTagRequest 删除标签请求
type DeleteTagRequest struct {
	TagPathParam
}

// ListTagRequest 标签列表请求
type ListTagRequest struct {
	CommonParam
}

// ListTagResponse 标签列表响应
type ListTagResponse struct {
	Tags     []*Tag    `json:"tags" doc:"List of tags"`
	PageInfo *PageInfo `json:"pageInfo" doc:"Pagination information"`
}
