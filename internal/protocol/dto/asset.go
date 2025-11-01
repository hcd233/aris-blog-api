package dto

import "mime/multipart"

// Image 图片信息
type Image struct {
	Name      string `json:"name" doc:"Image name"`
	Size      int64  `json:"size" doc:"Image size in bytes"`
	CreatedAt string `json:"createdAt" doc:"Creation timestamp"`
}

// UserView 用户浏览
type UserView struct {
	ViewID       uint   `json:"viewID" doc:"View record ID"`
	Progress     int8   `json:"progress" doc:"Reading progress percentage (0-100)"`
	LastViewedAt string `json:"lastViewedAt" doc:"Last viewed timestamp"`
	UserID       uint   `json:"userID" doc:"User ID"`
	ArticleID    uint   `json:"articleID" doc:"Article ID"`
}

// ViewPathParam 浏览记录路径参数
type ViewPathParam struct {
	ViewID uint `path:"viewID" doc:"View record ID"`
}

// ObjectPathParam 对象路径参数
type ObjectPathParam struct {
	ObjectName string `path:"objectName" doc:"Object name"`
}

// ImageQueryParam 图片查询参数
type ImageQueryParam struct {
	Quality string `query:"quality" doc:"Image quality (low, medium, high)" enum:"low,medium,high" default:"medium"`
}

// ListUserLikeArticlesRequest 列出用户喜欢的文章请求
type ListUserLikeArticlesRequest struct {
	CommonParam
}

// ListUserLikeArticlesResponse 列出用户喜欢的文章响应
type ListUserLikeArticlesResponse struct {
	Articles []*Article `json:"articles" doc:"List of liked articles"`
	PageInfo *PageInfo  `json:"pageInfo" doc:"Pagination information"`
}

// ListUserLikeCommentsRequest 列出用户喜欢的评论请求
type ListUserLikeCommentsRequest struct {
	CommonParam
}

// ListUserLikeCommentsResponse 列出用户喜欢的评论响应
type ListUserLikeCommentsResponse struct {
	Comments []*Comment `json:"comments" doc:"List of liked comments"`
	PageInfo *PageInfo  `json:"pageInfo" doc:"Pagination information"`
}

// ListUserLikeTagsRequest 列出用户喜欢的标签请求
type ListUserLikeTagsRequest struct {
	CommonParam
}

// ListUserLikeTagsResponse 列出用户喜欢的标签响应
type ListUserLikeTagsResponse struct {
	Tags     []*Tag    `json:"tags" doc:"List of liked tags"`
	PageInfo *PageInfo `json:"pageInfo" doc:"Pagination information"`
}

// ListUserViewArticlesRequest 列出用户浏览的文章请求
type ListUserViewArticlesRequest struct {
	CommonParam
}

// ListUserViewArticlesResponse 列出用户浏览的文章响应
type ListUserViewArticlesResponse struct {
	UserViews []*UserView `json:"userViews" doc:"List of view records"`
	PageInfo  *PageInfo   `json:"pageInfo" doc:"Pagination information"`
}

// DeleteUserViewRequest 删除用户浏览记录请求
type DeleteUserViewRequest struct {
	ViewPathParam
}

// ListImagesRequest 列出图片请求
type ListImagesRequest struct {
	EmptyRequest
}

// ListImagesResponse 列出图片响应
type ListImagesResponse struct {
	Images []*Image `json:"images" doc:"List of images"`
}

// UploadImageRequest 上传图片请求
type UploadImageRequest struct {
	RawBody multipart.FileHeader
}

// UploadImageResponse 上传图片响应
type UploadImageResponse struct {
	ImageName string `json:"imageName" doc:"Uploaded image name"`
}

// GetImageRequest 获取图片请求
type GetImageRequest struct {
	ObjectPathParam
	ImageQueryParam
}

// DeleteImageRequest 删除图片请求
type DeleteImageRequest struct {
	ObjectPathParam
}
