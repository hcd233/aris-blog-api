package dto

import "io"

// Image 表示对象存储中的图片信息
type Image struct {
    Name      string `json:"name" doc:"对象名称"`
    Size      int64  `json:"size" doc:"文件大小（字节）"`
    CreatedAt string `json:"createdAt" doc:"文件最后修改时间"`
}

// ImagePathParam 图片对象路径参数
type ImagePathParam struct {
    ImageName string `path:"objectName" doc:"对象存储中的文件名"`
}

// ImageQualityQueryParam 图片质量查询参数
type ImageQualityQueryParam struct {
    Quality string `query:"quality" enum:"raw,thumb" doc:"图片质量，raw 为原图，thumb 为缩略图"`
}

// GetImageRequest 获取图片请求
type GetImageRequest struct {
    ImagePathParam
    ImageQualityQueryParam
}

// GetImageResponse 获取图片响应
type GetImageResponse struct {
    PresignedURL string `json:"presignedURL" doc:"图片的预签名访问地址"`
}

// DeleteImageRequest 删除图片请求
type DeleteImageRequest struct {
    ImagePathParam
}

// UploadImageRequest 上传图片请求
type UploadImageRequest struct {
    FileName    string            `json:"-" doc:"上传文件名"`
    Size        int64             `json:"-" doc:"文件大小（字节）"`
    ContentType string            `json:"-" doc:"文件内容类型"`
    File        io.ReadSeekCloser `json:"-" doc:"文件读取句柄"`
}

// ListImagesRequest 列出用户图片请求
type ListImagesRequest struct{}

// ListImagesResponse 列出用户图片响应
type ListImagesResponse struct {
    Images []*Image `json:"images" doc:"图片列表"`
}

// ListUserLikeArticlesRequest 列出用户喜欢的文章请求
type ListUserLikeArticlesRequest struct {
    CommonParam
}

// ListUserLikeArticlesResponse 列出用户喜欢的文章响应
type ListUserLikeArticlesResponse struct {
    Articles []*Article `json:"articles" doc:"文章列表"`
    PageInfo *PageInfo  `json:"pageInfo" doc:"分页信息"`
}

// ListUserLikeCommentsRequest 列出用户喜欢的评论请求
type ListUserLikeCommentsRequest struct {
    CommonParam
}

// ListUserLikeCommentsResponse 列出用户喜欢的评论响应
type ListUserLikeCommentsResponse struct {
    Comments []*Comment `json:"comments" doc:"评论列表"`
    PageInfo *PageInfo  `json:"pageInfo" doc:"分页信息"`
}

// ListUserLikeTagsRequest 列出用户喜欢的标签请求
type ListUserLikeTagsRequest struct {
    CommonParam
}

// ListUserLikeTagsResponse 列出用户喜欢的标签响应
type ListUserLikeTagsResponse struct {
    Tags     []*Tag   `json:"tags" doc:"标签列表"`
    PageInfo *PageInfo `json:"pageInfo" doc:"分页信息"`
}

// UserView 用户浏览记录
type UserView struct {
    ViewID       uint   `json:"viewID" doc:"浏览记录 ID"`
    Progress     int8   `json:"progress" doc:"阅读进度（0-100）"`
    LastViewedAt string `json:"lastViewedAt" doc:"最近浏览时间"`
    ArticleID    uint   `json:"articleID" doc:"文章 ID"`
}

// ListUserViewArticlesRequest 列出用户浏览记录请求
type ListUserViewArticlesRequest struct {
    CommonParam
}

// ListUserViewArticlesResponse 列出用户浏览记录响应
type ListUserViewArticlesResponse struct {
    UserViews []*UserView `json:"userViews" doc:"浏览记录列表"`
    PageInfo  *PageInfo   `json:"pageInfo" doc:"分页信息"`
}

// ViewPathParam 浏览记录路径参数
type ViewPathParam struct {
    ViewID uint `path:"viewID" doc:"浏览记录 ID"`
}

// DeleteUserViewRequest 删除浏览记录请求
type DeleteUserViewRequest struct {
    ViewPathParam
}

