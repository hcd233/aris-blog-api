// Package dto 资产DTO
package dto

// Image 图片
//
//	author centonhuang
//	update 2025-10-30
type Image struct {
	Name      string `json:"name" doc:"Image filename"`
	Size      int64  `json:"size" doc:"Image size in bytes"`
	CreatedAt string `json:"createdAt" doc:"Upload timestamp"`
}

// ListImagesRequest 列出图片请求
//
//	author centonhuang
//	update 2025-10-30
type ListImagesRequest struct{}

// ListImagesResponse 列出图片响应
//
//	author centonhuang
//	update 2025-10-30
type ListImagesResponse struct {
	Images []*Image `json:"images" doc:"List of images"`
}

// GetImageRequest 获取图片请求
//
//	author centonhuang
//	update 2025-10-30
type GetImageRequest struct {
	ObjectName string `path:"objectName" doc:"Image filename"`
	Quality    string `query:"quality" enum:"raw,thumb" doc:"Image quality (raw or thumb)"`
}

// GetImageResponse 获取图片响应 (returns redirect)
//
//	author centonhuang
//	update 2025-10-30
type GetImageResponse struct {
	PresignedURL string `json:"presignedURL" doc:"Presigned URL for the image"`
}

// DeleteImageRequest 删除图片请求
//
//	author centonhuang
//	update 2025-10-30
type DeleteImageRequest struct {
	ObjectName string `path:"objectName" doc:"Image filename to delete"`
}

// DeleteImageResponse 删除图片响应
//
//	author centonhuang
//	update 2025-10-30
type DeleteImageResponse struct{}

// UserView 用户浏览记录
//
//	author centonhuang
//	update 2025-10-30
type UserView struct {
	ViewID       uint   `json:"viewID" doc:"Unique identifier for the view record"`
	Progress     int8   `json:"progress" doc:"Reading progress percentage"`
	LastViewedAt string `json:"lastViewedAt" doc:"Last viewed timestamp"`
	UserID       uint   `json:"userID" doc:"User ID"`
	ArticleID    uint   `json:"articleID" doc:"Article ID"`
}

// ListUserViewArticlesRequest 列出用户浏览的文章请求
//
//	author centonhuang
//	update 2025-10-30
type ListUserViewArticlesRequest struct {
	Page     *int `query:"page" minimum:"1" doc:"Page number (starting from 1)"`
	PageSize *int `query:"pageSize" minimum:"1" maximum:"50" doc:"Number of items per page"`
}

// ListUserViewArticlesResponse 列出用户浏览的文章响应
//
//	author centonhuang
//	update 2025-10-30
type ListUserViewArticlesResponse struct {
	UserViews []*UserView `json:"userViews" doc:"List of user view records"`
	PageInfo  *PageInfo   `json:"pageInfo" doc:"Pagination information"`
}

// DeleteUserViewRequest 删除用户浏览记录请求
//
//	author centonhuang
//	update 2025-10-30
type DeleteUserViewRequest struct {
	ViewID uint `path:"viewID" doc:"View record ID to delete"`
}

// DeleteUserViewResponse 删除用户浏览记录响应
//
//	author centonhuang
//	update 2025-10-30
type DeleteUserViewResponse struct{}

// ListUserLikeArticlesRequest 列出用户喜欢的文章请求
//
//	author centonhuang
//	update 2025-10-30
type ListUserLikeArticlesRequest struct {
	Page     *int `query:"page" minimum:"1" doc:"Page number (starting from 1)"`
	PageSize *int `query:"pageSize" minimum:"1" maximum:"50" doc:"Number of items per page"`
}

// ListUserLikeArticlesResponse 列出用户喜欢的文章响应
//
//	author centonhuang
//	update 2025-10-30
type ListUserLikeArticlesResponse struct {
	Articles []*Article `json:"articles" doc:"List of liked articles"`
	PageInfo *PageInfo  `json:"pageInfo" doc:"Pagination information"`
}

// ListUserLikeCommentsRequest 列出用户喜欢的评论请求
//
//	author centonhuang
//	update 2025-10-30
type ListUserLikeCommentsRequest struct {
	Page     *int `query:"page" minimum:"1" doc:"Page number (starting from 1)"`
	PageSize *int `query:"pageSize" minimum:"1" maximum:"50" doc:"Number of items per page"`
}

// ListUserLikeCommentsResponse 列出用户喜欢的评论响应
//
//	author centonhuang
//	update 2025-10-30
type ListUserLikeCommentsResponse struct {
	Comments []*Comment `json:"comments" doc:"List of liked comments"`
	PageInfo *PageInfo  `json:"pageInfo" doc:"Pagination information"`
}

// ListUserLikeTagsRequest 列出用户喜欢的标签请求
//
//	author centonhuang
//	update 2025-10-30
type ListUserLikeTagsRequest struct {
	Page     *int `query:"page" minimum:"1" doc:"Page number (starting from 1)"`
	PageSize *int `query:"pageSize" minimum:"1" maximum:"50" doc:"Number of items per page"`
}

// ListUserLikeTagsResponse 列出用户喜欢的标签响应
//
//	author centonhuang
//	update 2025-10-30
type ListUserLikeTagsResponse struct {
	Tags     []*Tag    `json:"tags" doc:"List of liked tags"`
	PageInfo *PageInfo `json:"pageInfo" doc:"Pagination information"`
}
