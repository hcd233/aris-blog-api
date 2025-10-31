// Package dto 标签DTO
package dto

// Tag 标签
//
//	author centonhuang
//	update 2025-10-30
type Tag struct {
	TagID       uint   `json:"tagID" doc:"Unique identifier for the tag"`
	Name        string `json:"name" doc:"Tag name"`
	Slug        string `json:"slug" doc:"URL-friendly tag identifier"`
	Description string `json:"description,omitempty" doc:"Tag description"`
	UserID      uint   `json:"userID,omitempty" doc:"ID of the user who created the tag"`
	CreatedAt   string `json:"createdAt,omitempty" doc:"Timestamp when the tag was created"`
	UpdatedAt   string `json:"updatedAt,omitempty" doc:"Timestamp when the tag was last updated"`
	Likes       uint   `json:"likes,omitempty" doc:"Number of likes for the tag"`
}

// CreateTagRequest 创建标签请求
//
//	author centonhuang
//	update 2025-10-30
type CreateTagRequest struct {
	Body *CreateTagBody `json:"body" doc:"Request body containing tag information"`
}

// CreateTagBody 创建标签请求体
//
//	author centonhuang
//	update 2025-10-30
type CreateTagBody struct {
	Name        string `json:"name" doc:"Tag name"`
	Slug        string `json:"slug" doc:"URL-friendly tag identifier"`
	Description string `json:"description,omitempty" doc:"Tag description"`
}

// CreateTagResponse 创建标签响应
//
//	author centonhuang
//	update 2025-10-30
type CreateTagResponse struct {
	Tag *Tag `json:"tag" doc:"Created tag information"`
}

// GetTagInfoRequest 获取标签信息请求
//
//	author centonhuang
//	update 2025-10-30
type GetTagInfoRequest struct {
	TagID uint `path:"tagID" doc:"Unique identifier of the tag to retrieve"`
}

// GetTagInfoResponse 获取标签信息响应
//
//	author centonhuang
//	update 2025-10-30
type GetTagInfoResponse struct {
	Tag *Tag `json:"tag" doc:"Tag information"`
}

// UpdateTagRequest 更新标签请求
//
//	author centonhuang
//	update 2025-10-30
type UpdateTagRequest struct {
	TagID uint           `path:"tagID" doc:"Unique identifier of the tag to update"`
	Body  *UpdateTagBody `json:"body" doc:"Request body containing fields to update"`
}

// UpdateTagBody 更新标签请求体
//
//	author centonhuang
//	update 2025-10-30
type UpdateTagBody struct {
	Name        string `json:"name" doc:"New tag name"`
	Slug        string `json:"slug" doc:"New URL-friendly tag identifier"`
	Description string `json:"description,omitempty" doc:"New tag description"`
}

// UpdateTagResponse 更新标签响应
//
//	author centonhuang
//	update 2025-10-30
type UpdateTagResponse struct{}

// DeleteTagRequest 删除标签请求
//
//	author centonhuang
//	update 2025-10-30
type DeleteTagRequest struct {
	TagID uint `path:"tagID" doc:"Unique identifier of the tag to delete"`
}

// DeleteTagResponse 删除标签响应
//
//	author centonhuang
//	update 2025-10-30
type DeleteTagResponse struct{}

// ListTagsRequest 列出标签请求
//
//	author centonhuang
//	update 2025-10-30
type ListTagsRequest struct {
	Page     *int `query:"page" minimum:"1" doc:"Page number (starting from 1)"`
	PageSize *int `query:"pageSize" minimum:"1" maximum:"50" doc:"Number of items per page"`
}

// ListTagsResponse 列出标签响应
//
//	author centonhuang
//	update 2025-10-30
type ListTagsResponse struct {
	Tags     []*Tag    `json:"tags" doc:"List of tags"`
	PageInfo *PageInfo `json:"pageInfo" doc:"Pagination information"`
}

// PageInfo 分页信息
//
//	author centonhuang
//	update 2025-10-30
type PageInfo struct {
	Page     int   `json:"page" doc:"Current page number"`
	PageSize int   `json:"pageSize" doc:"Number of items per page"`
	Total    int64 `json:"total" doc:"Total number of items"`
}
