package dto

import "mime/multipart"



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

// ListImagesResponse 列出图片响应
type ListImagesResponse struct {
	Images []*Image `json:"images" doc:"List of images"`
}

// UploadImageRequest 上传图片请求
type UploadImageRequest struct {
	RawBody *multipart.FileHeader
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
