package dto

// Tag 标签信息
//
//	author centonhuang
//	update 2025-10-31 05:32:00
type Tag struct {
	TagID       uint   `json:"tagID" doc:"Tag ID"`
	Name        string `json:"name" doc:"Tag name"`
	Slug        string `json:"slug" doc:"Tag slug"`
	Description string `json:"description,omitempty" doc:"Tag description"`
	CreatedAt   string `json:"createdAt,omitempty" doc:"Creation timestamp"`
	UpdatedAt   string `json:"updatedAt,omitempty" doc:"Update timestamp"`
	Likes       uint   `json:"likes,omitempty" doc:"Number of likes"`
}

// TagPathParam 标签路径参数
type TagPathParam struct {
	TagID uint `path:"tagID" doc:"Tag ID"`
}

// TagCreateRequestBody 创建标签请求体
type TagCreateRequestBody struct {
	Name        string `json:"name" doc:"Tag name"`
	Slug        string `json:"slug" doc:"Tag slug"`
	Description string `json:"description" doc:"Tag description"`
}

// TagCreateRequest 创建标签请求
type TagCreateRequest struct {
	Body *TagCreateRequestBody `json:"body" doc:"Fields for creating tag"`
}

// TagCreateResponse 创建标签响应
type TagCreateResponse struct {
	Tag *Tag `json:"tag" doc:"Successfully created tag"`
}

// TagGetRequest 获取标签请求
type TagGetRequest struct {
	TagPathParam
}

// TagGetResponse 获取标签响应
type TagGetResponse struct {
	Tag *Tag `json:"tag" doc:"Tag details"`
}

// TagUpdateRequestBody 更新标签请求体
type TagUpdateRequestBody struct {
	Name        string `json:"name" doc:"Tag name"`
	Slug        string `json:"slug" doc:"Tag slug"`
	Description string `json:"description" doc:"Tag description"`
}

// TagUpdateRequest 更新标签请求
type TagUpdateRequest struct {
	TagPathParam
	Body *TagUpdateRequestBody `json:"body" doc:"Updatable tag fields"`
}

// TagUpdateResponse 更新标签响应
type TagUpdateResponse struct{}

// TagDeleteRequest 删除标签请求
type TagDeleteRequest struct {
	TagPathParam
}

// TagDeleteResponse 删除标签响应
type TagDeleteResponse struct{}

// TagListRequest 标签列表请求
type TagListRequest struct {
	PaginationQuery
}

// TagListResponse 标签列表响应
type TagListResponse struct {
	Tags     []*Tag    `json:"tags" doc:"List of tags"`
	PageInfo *PageInfo `json:"pageInfo" doc:"Pagination information"`
}
