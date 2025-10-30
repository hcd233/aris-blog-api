package protocol

import (
	"io"

	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
)

// HumaResponse
//
//	author centonhuang
//	update 2025-10-31 01:18:05
//	param BodyT any
type HumaResponse[BodyT any] struct {
	Status int
	Body BodyT `json:"data"`
}

// PageInfo 分页信息
//
//	author centonhuang
//	update 2025-01-05 12:26:07
type PageInfo struct {
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
	Total    int64 `json:"total"`
}

// PingResponse 健康检查响应
//
//	author centonhuang
//	update 2025-01-04 20:47:11
type PingResponse struct {
	Status string `json:"status"`
}

// RefreshTokenRequest 刷新令牌请求
//
//	author centonhuang
//	update 2025-01-04 17:16:09
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// RefreshTokenResponse 刷新令牌响应
//
//	author centonhuang
//	update 2025-01-04 17:16:12
type RefreshTokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// User 用户
//
//	author centonhuang
//	update 2025-01-05 11:37:01
type User struct {
	UserID    uint   `json:"userID"`
	Name      string `json:"name"`
	Email     string `json:"email,omitempty"`
	Avatar    string `json:"avatar"`
	CreatedAt string `json:"createdAt,omitempty"`
	LastLogin string `json:"lastLogin,omitempty"`
}

// CurUser 当前用户
//
//	author centonhuang
//	update 2025-01-05 11:37:32
type CurUser struct {
	User
	Permission string `json:"permission"`
}

// GetCurUserInfoRequest 获取当前用户信息请求
//
//	author centonhuang
//	update 2025-01-04 21:00:54
type GetCurUserInfoRequest struct {
	UserID uint `json:"userID"`
}

// GetCurUserInfoResponse 获取当前用户信息响应
//
//	author centonhuang
//	update 2025-01-04 21:00:59
type GetCurUserInfoResponse struct {
	User *CurUser `json:"user"`
}

// GetUserInfoRequest 获取用户信息请求
//
//	author centonhuang
//	update 2025-01-04 21:19:41
type GetUserInfoRequest struct {
	UserID uint `json:"userID"`
}

// GetUserInfoResponse 获取用户信息响应
//
//	author centonhuang
//	update 2025-01-04 21:19:44
type GetUserInfoResponse struct {
	User *User `json:"user"`
}

// UpdateUserInfoRequest 更新用户信息请求
//
//	author centonhuang
//	update 2025-01-04 21:19:47
type UpdateUserInfoRequest struct {
	UserID          uint   `json:"userID"`
	UpdatedUserName string `json:"updatedUserName"`
}

// UpdateUserInfoResponse 更新用户信息响应
//
//	author centonhuang
//	update 2025-01-05 11:35:18
type UpdateUserInfoResponse struct{}

// Tag 标签
//
//	author centonhuang
//	update 2025-01-05 12:05:42
type Tag struct {
	TagID       uint   `json:"tagID"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description,omitempty"`
	UserID      uint   `json:"userID,omitempty"`
	CreatedAt   string `json:"createdAt,omitempty"`
	UpdatedAt   string `json:"updatedAt,omitempty"`
	Likes       uint   `json:"likes,omitempty"`
}

// CreateTagRequest 创建标签请求
//
//	author centonhuang
//	update 2025-01-05 11:48:36
type CreateTagRequest struct {
	UserID      uint   `json:"userID"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

// CreateTagResponse 创建标签响应
//
//	author centonhuang
//	update 2025-01-05 11:48:36
type CreateTagResponse struct {
	Tag *Tag `json:"tag"`
}

// GetTagInfoRequest 获取标签信息请求
//
//	author centonhuang
//	update 2025-01-05 12:03:33
type GetTagInfoRequest struct {
	TagID uint `json:"tagID"`
}

// GetTagInfoResponse 获取标签信息响应
//
//	author centonhuang
//	update 2025-01-05 11:48:36
type GetTagInfoResponse struct {
	Tag *Tag `json:"tag"`
}

// UpdateTagRequest 更新标签请求
//
//	author centonhuang
//	update 2025-01-05 12:07:42
type UpdateTagRequest struct {
	UserID      uint   `json:"userID"`
	TagID       uint   `json:"tagID"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

// UpdateTagResponse 更新标签响应
//
//	author centonhuang
//	update 2025-01-05 12:07:42
type UpdateTagResponse struct{}

// DeleteTagRequest 删除标签请求
//
//	author centonhuang
//	update 2025-01-05 11:48:36
type DeleteTagRequest struct {
	UserID uint `json:"userID"`
	TagID  uint `json:"tagID"`
}

// DeleteTagResponse 删除标签响应
//
//	author centonhuang
//	update 2025-01-05 11:48:36
type DeleteTagResponse struct{}

// ListTagsRequest 列出标签请求
//
//	author centonhuang
//	update 2025-01-05 11:48:36
type ListTagsRequest struct {
	PaginateParam *PaginateParam
}

// ListTagsResponse 列出标签响应
//
//	author centonhuang
//	update 2025-01-05 11:48:36
type ListTagsResponse struct {
	Tags     []*Tag    `json:"tags"`
	PageInfo *PageInfo `json:"pageInfo"`
}

// ListUserTagsRequest 列出用户标签请求
//
//	author centonhuang
//	update 2025-01-05 13:23:35
type ListUserTagsRequest struct {
	UserName      string
	PaginateParam *PaginateParam
}

// ListUserTagsResponse 列出用户标签响应
//
//	author centonhuang
//	update 2025-01-05 13:23:37
type ListUserTagsResponse struct {
	Tags     []*Tag    `json:"tags"`
	PageInfo *PageInfo `json:"pageInfo"`
}

// Category 分类
//
//	author centonhuang
//	update 2025-01-05 13:22:49
type Category struct {
	CategoryID uint   `json:"categoryID"`
	Name       string `json:"name"`
	ParentID   uint   `json:"parentID,omitempty"`
	CreatedAt  string `json:"createdAt,omitempty"`
	UpdatedAt  string `json:"updatedAt,omitempty"`
}

// Article 文章
//
//	author centonhuang
//	update 2025-01-05 15:23:26
type Article struct {
	ArticleID   uint      `json:"articleID"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Status      string    `json:"status"`
	User        *User     `json:"userID"`
	Category    *Category `json:"category"`
	CreatedAt   string    `json:"createdAt"`
	UpdatedAt   string    `json:"updatedAt"`
	PublishedAt string    `json:"publishedAt"`
	Likes       uint      `json:"likes"`
	Views       uint      `json:"views"`
	Tags        []*Tag    `json:"tags"`
	Comments    int       `json:"comments"`
}

// CreateCategoryRequest 创建分类请求
//
//	author centonhuang
//	update 2025-01-05 13:22:59
type CreateCategoryRequest struct {
	UserID   uint   `json:"userID"`
	Name     string `json:"name"`
	ParentID uint   `json:"parentID"`
}

// CreateCategoryResponse 创建分类响应
//
//	author centonhuang
//	update 2025-01-05 13:23:01
type CreateCategoryResponse struct {
	Category *Category `json:"category"`
}

// GetCategoryInfoRequest 获取分类信息请求
//
//	author centonhuang
//	update 2025-01-05 13:23:03
type GetCategoryInfoRequest struct {
	UserID     uint `json:"userID"`
	CategoryID uint `json:"categoryID"`
}

// GetCategoryInfoResponse 获取分类信息响应
//
//	author centonhuang
//	update 2025-01-05 13:23:06
type GetCategoryInfoResponse struct {
	Category *Category `json:"category"`
}

// GetRootCategoryRequest 获取根分类请求
//
//	author centonhuang
//	update 2025-01-05 13:23:08
type GetRootCategoryRequest struct {
	UserID uint `json:"userID"`
}

// GetRootCategoryResponse 获取根分类响应
//
//	author centonhuang
//	update 2025-01-05 13:23:13
type GetRootCategoryResponse struct {
	Category *Category `json:"category"`
}

// UpdateCategoryRequest 更新分类请求
//
//	author centonhuang
//	update 2025-01-05 13:23:14
type UpdateCategoryRequest struct {
	UserID     uint   `json:"userID"`
	CategoryID uint   `json:"categoryID"`
	Name       string `json:"name"`
	ParentID   uint   `json:"parentID"`
}

// UpdateCategoryResponse 更新分类响应
//
//	author centonhuang
//	update 2025-01-05 13:23:16
type UpdateCategoryResponse struct {
	Category *Category `json:"category"`
}

// DeleteCategoryRequest 删除分类请求
//
//	author centonhuang
//	update 2025-01-05 13:23:18
type DeleteCategoryRequest struct {
	UserID     uint `json:"userID"`
	CategoryID uint `json:"categoryID"`
}

// DeleteCategoryResponse 删除分类响应
//
//	author centonhuang
//	update 2025-01-05 13:23:20
type DeleteCategoryResponse struct{}

// ListChildrenCategoriesRequest 列出子分类请求
//
//	author centonhuang
//	update 2025-01-05 13:23:21
type ListChildrenCategoriesRequest struct {
	UserID        uint           `json:"userID"`
	CategoryID    uint           `json:"categoryID"`
	PaginateParam *PaginateParam `json:"paginateParam"`
}

// ListChildrenCategoriesResponse 列出子分类响应
//
//	author centonhuang
//	update 2025-01-05 13:23:23
type ListChildrenCategoriesResponse struct {
	Categories []*Category `json:"categories"`
	PageInfo   *PageInfo   `json:"pageInfo"`
}

// ListChildrenArticlesRequest 列出子文章请求
//
//	author centonhuang
//	update 2025-01-05 13:23:25
type ListChildrenArticlesRequest struct {
	UserID        uint           `json:"userID"`
	CategoryID    uint           `json:"categoryID"`
	PaginateParam *PaginateParam `json:"paginateParam"`
}

// ListChildrenArticlesResponse 列出子文章响应
//
//	author centonhuang
//	update 2025-01-05 13:23:26
type ListChildrenArticlesResponse struct {
	Articles []*Article `json:"articles"`
	PageInfo *PageInfo  `json:"pageInfo"`
}

// LoginRequest OAuth2登录请求
//
//	author centonhuang
//	update 2025-01-05 14:23:26
type LoginRequest struct{}

// LoginResponse OAuth2登录响应
//
//	author centonhuang
//	update 2025-01-05 14:23:26
type LoginResponse struct {
	RedirectURL string `json:"redirectURL"`
}

// CallbackRequest OAuth2回调请求
//
//	author centonhuang
//	update 2025-01-05 14:23:26
type CallbackRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

// CallbackResponse OAuth2回调响应
//
//	author centonhuang
//	update 2025-01-05 14:23:26
type CallbackResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// CreateArticleRequest 创建文章请求
//
//	author centonhuang
//	update 2025-01-05 15:23:26
type CreateArticleRequest struct {
	UserID     uint     `json:"userID"`
	Title      string   `json:"title"`
	Slug       string   `json:"slug"`
	CategoryID uint     `json:"categoryID"`
	Tags       []string `json:"tags"`
}

// CreateArticleResponse 创建文章响应
//
//	author centonhuang
//	update 2025-01-05 15:23:26
type CreateArticleResponse struct {
	Article *Article `json:"article"`
}

// GetArticleInfoRequest 获取文章信息请求
//
//	author centonhuang
//	update 2025-01-05 15:23:26
type GetArticleInfoRequest struct {
	UserID    uint `json:"userID"`
	ArticleID uint `json:"articleID"`
}

// GetArticleInfoResponse 获取文章信息响应
//
//	author centonhuang
//	update 2025-01-05 15:23:26
type GetArticleInfoResponse struct {
	Article *Article `json:"article"`
}

// GetArticleInfoBySlugRequest 获取文章信息请求
//
//	author centonhuang
//	update 2025-01-19 15:23:26
type GetArticleInfoBySlugRequest struct {
	UserID      uint   `json:"userID"`
	AuthorName  string `json:"authorName"`
	ArticleSlug string `json:"articleSlug"`
}

// GetArticleInfoBySlugResponse 获取文章信息响应
//
//	author centonhuang
//	update 2025-01-19 15:23:26
type GetArticleInfoBySlugResponse struct {
	Article *Article `json:"article"`
}

// UpdateArticleRequest 更新文章请求
//
//	author centonhuang
//	update 2025-01-05 15:23:26
type UpdateArticleRequest struct {
	UserID            uint   `json:"userID"`
	ArticleID         uint   `json:"articleID"`
	UpdatedTitle      string `json:"title"`
	UpdatedSlug       string `json:"slug"`
	UpdatedCategoryID uint   `json:"categoryID"`
}

// UpdateArticleResponse 更新文章响应
//
//	author centonhuang
//	update 2025-01-05 15:23:26
type UpdateArticleResponse struct{}

// UpdateArticleStatusRequest 更新文章状态请求
//
//	author centonhuang
//	update 2025-01-05 15:23:26
type UpdateArticleStatusRequest struct {
	UserID    uint                `json:"userID"`
	ArticleID uint                `json:"articleID"`
	Status    model.ArticleStatus `json:"status"`
}

// UpdateArticleStatusResponse 更新文章状态响应
//
//	author centonhuang
//	update 2025-01-05 15:23:26
type UpdateArticleStatusResponse struct{}

// DeleteArticleRequest 删除文章请求
//
//	author centonhuang
//	update 2025-01-05 15:23:26
type DeleteArticleRequest struct {
	UserID    uint `json:"userID"`
	ArticleID uint `json:"articleID"`
}

// DeleteArticleResponse 删除文章响应
//
//	author centonhuang
//	update 2025-01-05 15:23:26
type DeleteArticleResponse struct{}

// ListArticlesRequest 列出文章请求
//
//	author centonhuang
//	update 2025-01-05 15:23:26
type ListArticlesRequest struct {
	PaginateParam *PaginateParam `json:"paginateParam"`
}

// ListArticlesResponse 列出文章响应
//
//	author centonhuang
//	update 2025-01-05 15:23:26
type ListArticlesResponse struct {
	Articles []*Article `json:"articles"`
	PageInfo *PageInfo  `json:"pageInfo"`
}

// ArticleVersion 文章版本
type ArticleVersion struct {
	ArticleVersionID uint   `json:"versionID"`
	ArticleID        uint   `json:"articleID"`
	VersionID        uint   `json:"version"`
	Content          string `json:"content"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
}

// CreateArticleVersionRequest 创建文章版本请求
//
//	author centonhuang
//	update 2025-01-05 16:42:48
type CreateArticleVersionRequest struct {
	UserID    uint   `json:"userID"`
	ArticleID uint   `json:"articleID"`
	Content   string `json:"content"`
}

// CreateArticleVersionResponse 创建文章版本响应
//
//	author centonhuang
//	update 2025-01-05 16:42:48
type CreateArticleVersionResponse struct {
	ArticleVersion *ArticleVersion `json:"articleVersion"`
}

// GetArticleVersionInfoRequest 获取文章版本信息请求
//
//	author centonhuang
//	update 2025-01-05 16:42:48
type GetArticleVersionInfoRequest struct {
	UserID    uint `json:"userID"`
	ArticleID uint `json:"articleID"`
	VersionID uint `json:"versionID"`
}

// GetArticleVersionInfoResponse 获取文章版本信息响应
//
//	author centonhuang
//	update 2025-01-05 16:42:48
type GetArticleVersionInfoResponse struct {
	Version *ArticleVersion `json:"version"`
}

// GetLatestArticleVersionInfoRequest 获取最新文章版本信息请求
//
//	author centonhuang
//	update 2025-01-05 16:42:48
type GetLatestArticleVersionInfoRequest struct {
	UserID    uint `json:"userID"`
	ArticleID uint `json:"articleID"`
}

// GetLatestArticleVersionInfoResponse 获取最新文章版本信息响应
//
//	author centonhuang
//	update 2025-01-05 16:42:48
type GetLatestArticleVersionInfoResponse struct {
	Version *ArticleVersion `json:"version"`
}

// ListArticleVersionsRequest 列出文章版本请求
//
//	author centonhuang
//	update 2025-01-05 16:42:48
type ListArticleVersionsRequest struct {
	UserID        uint           `json:"userID"`
	ArticleID     uint           `json:"articleID"`
	PaginateParam *PaginateParam `json:"paginateParam"`
}

// ListArticleVersionsResponse 列出文章版本响应
//
//	author centonhuang
//	update 2025-01-05 16:42:48
type ListArticleVersionsResponse struct {
	Versions []*ArticleVersion `json:"versions"`
	PageInfo *PageInfo         `json:"pageInfo"`
}

// LikeArticleRequest 点赞文章请求
//
//	author centonhuang
//	update 2025-01-05 16:42:48
type LikeArticleRequest struct {
	UserID    uint `json:"userID"`
	ArticleID uint `json:"articleID"`
	Undo      bool `json:"undo"`
}

// LikeArticleResponse 点赞文章响应
//
//	author centonhuang
//	update 2025-01-05 16:43:09
type LikeArticleResponse struct{}

// LikeCommentRequest 点赞评论请求
//
//	author centonhuang
//	update 2025-01-05 16:43:20
type LikeCommentRequest struct {
	UserID    uint `json:"userID"`
	CommentID uint `json:"commentID"`
	Undo      bool `json:"undo"`
}

// LikeCommentResponse 点赞评论响应
//
//	author centonhuang
//	update 2025-01-05 16:43:21
type LikeCommentResponse struct{}

// LikeTagRequest 点赞标签请求
//
//	author centonhuang
//	update 2025-01-05 16:43:23
type LikeTagRequest struct {
	UserID uint `json:"userID"`
	TagID  uint `json:"tagID"`
	Undo   bool `json:"undo"`
}

// LikeTagResponse 点赞标签响应
//
//	author centonhuang
//	update 2025-01-05 16:43:25
type LikeTagResponse struct{}

// LogArticleViewRequest 记录文章浏览请求
//
//	author centonhuang
//	update 2025-01-05 16:43:26
type LogArticleViewRequest struct {
	UserID    uint `json:"userID"`
	ArticleID uint `json:"articleID"`
	Progress  int8 `json:"progress"`
}

// LogArticleViewResponse 记录文章浏览响应
//
//	author centonhuang
//	update 2025-01-05 16:43:28
type LogArticleViewResponse struct{}

// Comment 评论
//
//	author centonhuang
//	update 2025-01-05 16:43:29
type Comment struct {
	CommentID uint   `json:"commentID"`
	Content   string `json:"content"`
	UserID    uint   `json:"userID"`
	ReplyTo   uint   `json:"replyTo"`
	CreatedAt string `json:"createdAt"`
	Likes     uint   `json:"likes"`
}

// CreateArticleCommentRequest 创建文章评论请求
//
//	author centonhuang
//	update 2025-01-05 16:43:31
type CreateArticleCommentRequest struct {
	UserID    uint   `json:"userID"`
	ArticleID uint   `json:"articleID"`
	Content   string `json:"content"`
	ReplyTo   uint   `json:"replyTo"`
}

// CreateArticleCommentResponse 创建文章评论响应
//
//	author centonhuang
//	update 2025-01-05 16:43:33
type CreateArticleCommentResponse struct {
	Comment *Comment `json:"comment"`
}

// GetCommentInfoRequest 获取评论信息请求
//
//	author centonhuang
//	update 2025-01-05 16:43:34
type GetCommentInfoRequest struct {
	UserID    uint `json:"userID"`
	ArticleID uint `json:"articleID"`
	CommentID uint `json:"commentID"`
}

// GetCommentInfoResponse 获取评论信息响应
//
//	author centonhuang
//	update 2025-01-05 16:43:36
type GetCommentInfoResponse struct {
	Comment *Comment `json:"comment"`
}

// DeleteCommentRequest 删除评论请求
//
//	author centonhuang
//	update 2025-01-05 16:43:38
type DeleteCommentRequest struct {
	UserID    uint `json:"userID"`
	CommentID uint `json:"commentID"`
}

// DeleteCommentResponse 删除评论响应
//
//	author centonhuang
//	update 2025-01-05 16:43:39
type DeleteCommentResponse struct{}

// ListArticleCommentsRequest 列出文章评论请求
//
//	author centonhuang
//	update 2025-01-05 16:43:41
type ListArticleCommentsRequest struct {
	UserID        uint           `json:"userID"`
	ArticleID     uint           `json:"articleID"`
	PaginateParam *PaginateParam `json:"paginateParam"`
}

// ListArticleCommentsResponse 列出文章评论响应
//
//	author centonhuang
//	update 2025-01-05 16:43:43
type ListArticleCommentsResponse struct {
	Comments []*Comment `json:"comments"`
	PageInfo *PageInfo  `json:"pageInfo"`
}

// ListChildrenCommentsRequest 列出子评论请求
//
//	author centonhuang
//	update 2025-01-05 16:43:44
type ListChildrenCommentsRequest struct {
	UserID        uint           `json:"userID"`
	CommentID     uint           `json:"commentID"`
	PaginateParam *PaginateParam `json:"paginateParam"`
}

// ListChildrenCommentsResponse 列出子评论响应
//
//	author centonhuang
//	update 2025-01-05 16:43:46
type ListChildrenCommentsResponse struct {
	Comments []*Comment `json:"comments"`
	PageInfo *PageInfo  `json:"pageInfo"`
}

// ListUserLikeArticlesRequest 列出用户喜欢的文章请求
//
//	author centonhuang
//	update 2025-01-05 16:43:48
type ListUserLikeArticlesRequest struct {
	UserID        uint           `json:"userID"`
	PaginateParam *PaginateParam `json:"paginateParam"`
}

// ListUserLikeArticlesResponse 列出用户喜欢的文章响应
//
//	author centonhuang
//	update 2025-01-05 16:43:50
type ListUserLikeArticlesResponse struct {
	Articles []*Article `json:"articles"`
	PageInfo *PageInfo  `json:"pageInfo"`
}

// ListUserLikeCommentsRequest 列出用户喜欢的评论请求
//
//	author centonhuang
//	update 2025-01-05 16:43:52
type ListUserLikeCommentsRequest struct {
	UserID        uint           `json:"userID"`
	PaginateParam *PaginateParam `json:"paginateParam"`
}

// ListUserLikeCommentsResponse 列出用户喜欢的评论响应
//
//	author centonhuang
//	update 2025-01-05 16:43:54
type ListUserLikeCommentsResponse struct {
	Comments []*Comment `json:"comments"`
	PageInfo *PageInfo  `json:"pageInfo"`
}

// ListUserLikeTagsRequest 列出用户喜欢的标签请求
//
//	author centonhuang
//	update 2025-01-05 16:43:56
type ListUserLikeTagsRequest struct {
	UserID        uint           `json:"userID"`
	PaginateParam *PaginateParam `json:"paginateParam"`
}

// ListUserLikeTagsResponse 列出用户喜欢的标签响应
//
//	author centonhuang
//	update 2025-01-05 16:43:58
type ListUserLikeTagsResponse struct {
	Tags     []*Tag    `json:"tags"`
	PageInfo *PageInfo `json:"pageInfo"`
}

// CreateBucketRequest 创建桶请求
//
//	author centonhuang
//	update 2025-01-05 17:03:19
type CreateBucketRequest struct {
	UserID uint `json:"userID"`
}

// CreateBucketResponse 创建桶响应
//
//	author centonhuang
//	update 2025-01-05 17:03:21
type CreateBucketResponse struct{}

// Image 图片
//
//	author centonhuang
//	update 2025-01-05 17:17:53
type Image struct {
	Name      string `json:"name"`
	Size      int64  `json:"size"`
	CreatedAt string `json:"createdAt"`
}

// ListImagesRequest 列出图片请求
//
//	author centonhuang
//	update 2025-01-05 17:03:22
type ListImagesRequest struct {
	UserID uint `json:"userID"`
}

// ListImagesResponse 列出图片响应
//
//	author centonhuang
//	update 2025-01-05 17:03:24
type ListImagesResponse struct {
	Images []*Image `json:"images"`
}

// UploadImageRequest 上传图片请求
//
//	author centonhuang
//	update 2025-01-05 17:03:25
type UploadImageRequest struct {
	UserID      uint   `json:"userID"`
	FileName    string `json:"fileName"`
	Size        int64  `json:"size"`
	ContentType string `json:"contentType"`
	ReadSeeker  io.ReadSeeker
}

// UploadImageResponse 上传图片响应
//
//	author centonhuang
//	update 2025-01-05 17:03:27
type UploadImageResponse struct{}

// GetImageRequest 获取图片请求
//
//	author centonhuang
//	update 2025-01-05 17:03:28
type GetImageRequest struct {
	UserID    uint   `json:"userID"`
	ImageName string `json:"imageName"`
	Quality   string `json:"quality"`
}

// GetImageResponse 获取图片响应
//
//	author centonhuang
//	update 2025-01-05 17:03:30
type GetImageResponse struct {
	PresignedURL string `json:"presignedURL"`
}

// DeleteImageRequest 删除图片请求
//
//	author centonhuang
//	update 2025-01-05 17:03:31
type DeleteImageRequest struct {
	UserID    uint   `json:"userID"`
	ImageName string `json:"imageName"`
}

// DeleteImageResponse 删除图片响应
//
//	author centonhuang
//	update 2025-01-05 17:03:33
type DeleteImageResponse struct{}

// UserView 用户浏览
//
//	author centonhuang
//	update 2025-01-05 17:53:09
type UserView struct {
	ViewID       uint   `json:"viewID"`
	Progress     int8   `json:"progress"`
	LastViewedAt string `json:"lastViewedAt"`
	UserID       uint   `json:"userID"`
	ArticleID    uint   `json:"articleID"`
}

// ListUserViewArticlesRequest 列出用户浏览的文章请求
//
//	author centonhuang
//	update 2025-01-05 17:03:38
type ListUserViewArticlesRequest struct {
	UserID        uint           `json:"userID"`
	PaginateParam *PaginateParam `json:"paginateParam"`
}

// ListUserViewArticlesResponse 列出用户浏览的文章响应
//
//	author centonhuang
//	update 2025-01-05 17:03:40
type ListUserViewArticlesResponse struct {
	UserViews []*UserView `json:"userViews"`
	PageInfo  *PageInfo   `json:"pageInfo"`
}

// DeleteUserViewRequest 删除用户浏览的文章请求
//
//	author centonhuang
//	update 2025-01-05 17:03:41
type DeleteUserViewRequest struct {
	UserID uint `json:"userID"`
	ViewID uint `json:"viewID"`
}

// DeleteUserViewResponse 删除用户浏览的文章响应
//
//	author centonhuang
//	update 2025-01-05 17:03:43
type DeleteUserViewResponse struct{}

// Prompt 提示词
//
//	author centonhuang
//	update 2025-01-05 18:07:58
type Prompt struct {
	PromptID  uint       `json:"promptID"`
	CreatedAt string     `json:"createdAt"`
	Task      string     `json:"task"`
	Version   uint       `json:"version"`
	Templates []Template `json:"templates"`
	Variables []string   `json:"variables"`
}

// GetPromptRequest 获取提示词请求
//
//	author centonhuang
//	update 2025-01-05 18:01:30
type GetPromptRequest struct {
	TaskName string `json:"taskName"`
	Version  uint   `json:"version"`
}

// GetPromptResponse 获取提示词响应
//
//	author centonhuang
//	update 2025-01-05 18:01:30
type GetPromptResponse struct {
	Prompt *Prompt `json:"prompt"`
}

// GetLatestPromptRequest 获取最新提示词请求
//
//	author centonhuang
//	update 2025-01-05 18:10:46
type GetLatestPromptRequest struct {
	TaskName string `json:"taskName"`
}

// GetLatestPromptResponse 获取最新提示词响应
//
//	author centonhuang
//	update 2025-01-05 18:10:51
type GetLatestPromptResponse struct {
	Prompt *Prompt `json:"prompt"`
}

// ListPromptRequest 列出提示词请求
//
//	author centonhuang
//	update 2025-01-05 18:13:20
type ListPromptRequest struct {
	TaskName      string         `json:"taskName"`
	PaginateParam *PaginateParam `json:"paginateParam"`
}

// ListPromptResponse 列出提示词响应
//
//	author centonhuang
//	update 2025-01-05 18:13:22
type ListPromptResponse struct {
	Prompts  []*Prompt `json:"prompts"`
	PageInfo *PageInfo `json:"pageInfo"`
}

// CreatePromptRequest 创建提示词请求
//
//	author centonhuang
//	update 2025-01-05 18:16:13
type CreatePromptRequest struct {
	TaskName  string     `json:"taskName"`
	Templates []Template `json:"templates"`
}

// CreatePromptResponse 创建提示词响应
//
//	author centonhuang
//	update 2025-01-05 18:16:13
type CreatePromptResponse struct{}

// GenerateContentCompletionRequest 生成内容完成请求
//
//	author centonhuang
//	update 2025-01-05 18:41:32
type GenerateContentCompletionRequest struct {
	UserID      uint    `json:"userID"`
	Context     string  `json:"context"`
	Instruction string  `json:"instruction"`
	Reference   string  `json:"reference"`
	Temperature float32 `json:"temperature"`
}

// GenerateContentCompletionResponse 生成内容完成响应
//
//	author centonhuang
//	update 2025-01-05 18:41:32
type GenerateContentCompletionResponse struct {
	TokenChan <-chan string
	ErrChan   <-chan error
}

// GenerateArticleSummaryRequest 生成文章总结请求
//
//	author centonhuang
//	update 2025-01-05 18:41:32
type GenerateArticleSummaryRequest struct {
	UserID      uint    `json:"userID"`
	ArticleID   uint    `json:"articleID"`
	Instruction string  `json:"instruction"`
	Temperature float32 `json:"temperature"`
}

// GenerateArticleSummaryResponse 生成文章总结响应
//
//	author centonhuang
//	update 2025-01-05 18:41:32
type GenerateArticleSummaryResponse struct {
	TokenChan <-chan string
	ErrChan   <-chan error
}

// GenerateArticleTranslationRequest 生成文章翻译请求
//
//	author centonhuang
//	update 2025-01-05 20:36:40
type GenerateArticleTranslationRequest struct{}

// GenerateArticleTranslationResponse 生成文章翻译响应
//
//	author centonhuang
//	update 2025-01-05 20:36:43
type GenerateArticleTranslationResponse struct {
	TokenChan <-chan string
	ErrChan   <-chan error
}

// GenerateArticleQARequest 生成文章问答请求
//
//	author centonhuang
//	update 2025-01-05 18:41:32
type GenerateArticleQARequest struct {
	UserID      uint    `json:"userID"`
	ArticleID   uint    `json:"articleID"`
	Question    string  `json:"question"`
	Temperature float32 `json:"temperature"`
}

// GenerateArticleQAResponse 生成文章问答响应
//
//	author centonhuang
//	update 2025-01-05 18:41:32
type GenerateArticleQAResponse struct {
	TokenChan <-chan string
	ErrChan   <-chan error
}

// GenerateTermExplainationRequest 生成术语解释请求
//
//	author centonhuang
//	update 2025-01-05 18:41:32
type GenerateTermExplainationRequest struct {
	UserID      uint    `json:"userID"`
	ArticleID   uint    `json:"articleID"`
	Term        string  `json:"term"`
	Position    uint    `json:"position"`
	Temperature float32 `json:"temperature"`
}

// GenerateTermExplainationResponse 生成术语解释响应
//
//	author centonhuang
//	update 2025-01-05 18:41:32
type GenerateTermExplainationResponse struct {
	TokenChan <-chan string
	ErrChan   <-chan error
}
