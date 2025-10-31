package dto

// Tag 标签信息
//
//	author centonhuang
//	update 2025-10-31 05:32:00
type Tag struct {
	TagID       uint   `json:"tagID" doc:"标签 ID"`
	Name        string `json:"name" doc:"标签名称"`
	Slug        string `json:"slug" doc:"标签别名"`
	Description string `json:"description,omitempty" doc:"标签描述"`
	UserID      uint   `json:"userID,omitempty" doc:"标签拥有者用户 ID"`
	CreatedAt   string `json:"createdAt,omitempty" doc:"创建时间"`
	UpdatedAt   string `json:"updatedAt,omitempty" doc:"更新时间"`
	Likes       uint   `json:"likes,omitempty" doc:"点赞数量"`
}

// TagPathParam 标签路径参数
type TagPathParam struct {
	TagID uint `path:"tagID" doc:"标签 ID"`
}

// TagCreateRequestBody 创建标签请求体
type TagCreateRequestBody struct {
	Name        string `json:"name" doc:"标签名称"`
	Slug        string `json:"slug" doc:"标签别名"`
	Description string `json:"description" doc:"标签描述"`
}

// TagCreateRequest 创建标签请求
type TagCreateRequest struct {
	UserID uint                  `json:"-"`
	Body   *TagCreateRequestBody `json:"body" doc:"创建标签字段"`
}

// TagCreateResponse 创建标签响应
type TagCreateResponse struct {
	Tag *Tag `json:"tag" doc:"新建成功的标签"`
}

// TagGetRequest 获取标签请求
type TagGetRequest struct {
	TagPathParam
}

// TagGetResponse 获取标签响应
type TagGetResponse struct {
	Tag *Tag `json:"tag" doc:"标签详情"`
}

// TagUpdateRequestBody 更新标签请求体
type TagUpdateRequestBody struct {
	Name        string `json:"name" doc:"标签名称"`
	Slug        string `json:"slug" doc:"标签别名"`
	Description string `json:"description" doc:"标签描述"`
}

// TagUpdateRequest 更新标签请求
type TagUpdateRequest struct {
	TagPathParam
	UserID uint                  `json:"-"`
	Body   *TagUpdateRequestBody `json:"body" doc:"可更新的标签字段"`
}

// TagUpdateResponse 更新标签响应
type TagUpdateResponse struct{}

// TagDeleteRequest 删除标签请求
type TagDeleteRequest struct {
	TagPathParam
	UserID uint `json:"-"`
}

// TagDeleteResponse 删除标签响应
type TagDeleteResponse struct{}

// TagListRequest 标签列表请求
type TagListRequest struct {
	PaginationQuery
}

// TagListResponse 标签列表响应
type TagListResponse struct {
	Tags     []*Tag    `json:"tags" doc:"标签列表"`
	PageInfo *PageInfo `json:"pageInfo" doc:"分页信息"`
}
