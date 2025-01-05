package protocol

import (
	"io"

	"github.com/hcd233/Aris-blog/internal/resource/database/model"
)

// PageInfo 分页信息
//
//	@author centonhuang
//	@update 2025-01-05 12:26:07
type PageInfo struct {
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
	Total    int64 `json:"total"`
}

// QueryInfo 查询信息
//
//	@author centonhuang
//	@update 2025-01-05 12:32:25
type QueryInfo struct {
	PageInfo
	Query  string   `json:"query"`
	Filter []string `json:"filter"`
}

// PingResponse 健康检查响应
//
//	@author centonhuang
//	@update 2025-01-04 20:47:11
type PingResponse struct {
	Status string `json:"status"`
}

// RefreshTokenRequest 刷新令牌请求
//
//	@author centonhuang
//	@update 2025-01-04 17:16:09
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// RefreshTokenResponse 刷新令牌响应
//
//	@author centonhuang
//	@update 2025-01-04 17:16:12
type RefreshTokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// User 用户
//
//	@author centonhuang
//	@update 2025-01-05 11:37:01
type User struct {
	UserID    uint   `json:"userID"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
	CreatedAt string `json:"createdAt"`
	LastLogin string `json:"lastLogin"`
}

// CurUser 当前用户
//
//	@author centonhuang
//	@update 2025-01-05 11:37:32
type CurUser struct {
	User
	Permission string `json:"permission"`
}

// GetCurUserInfoRequest 获取当前用户信息请求
//
//	@author centonhuang
//	@update 2025-01-04 21:00:54
type GetCurUserInfoRequest struct {
	CurUserID uint `json:"curUserID"`
}

// GetCurUserInfoResponse 获取当前用户信息响应
//
//	@author centonhuang
//	@update 2025-01-04 21:00:59
type GetCurUserInfoResponse struct {
	User *CurUser `json:"user"`
}

// GetUserInfoRequest 获取用户信息请求
//
//	@author centonhuang
//	@update 2025-01-04 21:19:41
type GetUserInfoRequest struct {
	UserName string `json:"userName"`
}

// GetUserInfoResponse 获取用户信息响应
//
//	@author centonhuang
//	@update 2025-01-04 21:19:44
type GetUserInfoResponse struct {
	User *User `json:"user"`
}

// UpdateUserInfoRequest 更新用户信息请求
//
//	@author centonhuang
//	@update 2025-01-04 21:19:47
type UpdateUserInfoRequest struct {
	CurUserName     string `json:"curUserName"`
	UserName        string `json:"userName"`
	UpdatedUserName string `json:"updatedUserName"`
}

// UpdateUserInfoResponse 更新用户信息响应
//
//	@author centonhuang
//	@update 2025-01-05 11:35:18
type UpdateUserInfoResponse struct{}

// QueryUserRequest 查询用户请求
//
//	@author centonhuang
//	@update 2025-01-05 11:35:17
type QueryUserRequest struct {
	QueryParam *QueryParam
}

// QueryUserResponse 查询用户响应
//
//	@author centonhuang
//	@update 2025-01-05 11:35:23
type QueryUserResponse struct {
	Users     []*User     `json:"users"`
	QueryInfo *QueryParam `json:"queryInfo"`
}

// Tag 标签
//
//	@author centonhuang
//	@update 2025-01-05 12:05:42
type Tag struct {
	TagID       uint   `json:"tagID"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	UserID      uint   `json:"userID"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	Likes       uint   `json:"likes"`
}

// CreateTagRequest 创建标签请求
//
//	@author centonhuang
//	@update 2025-01-05 11:48:36
type CreateTagRequest struct {
	CurUserID   uint   `json:"curUserID"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

// CreateTagResponse 创建标签响应
//
//	@author centonhuang
//	@update 2025-01-05 11:48:36
type CreateTagResponse struct{}

// GetTagInfoRequest 获取标签信息请求
//
//	@author centonhuang
//	@update 2025-01-05 12:03:33
type GetTagInfoRequest struct {
	TagSlug string `json:"tagSlug"`
}

// GetTagInfoResponse 获取标签信息响应
//
//	@author centonhuang
//	@update 2025-01-05 11:48:36
type GetTagInfoResponse struct {
	Tag *Tag `json:"tag"`
}

// UpdateTagRequest 更新标签请求
//
//	@author centonhuang
//	@update 2025-01-05 12:07:42
type UpdateTagRequest struct {
	CurUserID   uint   `json:"curUserID"`
	TagSlug     string `json:"tagSlug"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

// UpdateTagResponse 更新标签响应
//
//	@author centonhuang
//	@update 2025-01-05 12:07:42
type UpdateTagResponse struct{}

// DeleteTagRequest 删除标签请求
//
//	@author centonhuang
//	@update 2025-01-05 11:48:36
type DeleteTagRequest struct {
	CurUserID uint   `json:"curUserID"`
	TagName   string `json:"tagName"`
}

// DeleteTagResponse 删除标签响应
//
//	@author centonhuang
//	@update 2025-01-05 11:48:36
type DeleteTagResponse struct{}

// ListTagsRequest 列出标签请求
//
//	@author centonhuang
//	@update 2025-01-05 11:48:36
type ListTagsRequest struct {
	PageParam *PageParam
}

// ListTagsResponse 列出标签响应
//
//	@author centonhuang
//	@update 2025-01-05 11:48:36
type ListTagsResponse struct {
	Tags     []*Tag    `json:"tags"`
	PageInfo *PageInfo `json:"pageInfo"`
}

// ListUserTagsRequest 列出用户标签请求
//
//	@author centonhuang
//	@update 2025-01-05 13:23:35
type ListUserTagsRequest struct {
	UserName  string
	PageParam *PageParam
}

// ListUserTagsResponse 列出用户标签响应
//
//	@author centonhuang
//	@update 2025-01-05 13:23:37
type ListUserTagsResponse struct {
	Tags     []*Tag    `json:"tags"`
	PageInfo *PageInfo `json:"pageInfo"`
}

// QueryTagRequest 查询标签请求
//
//	@author centonhuang
//	@update 2025-01-05 12:37:54
type QueryTagRequest struct {
	QueryParam *QueryParam
}

// QueryTagResponse 查询标签响应
//
//	@author centonhuang
//	@update 2025-01-05 12:37:54
type QueryTagResponse struct {
	Tags      []*Tag      `json:"tags"`
	QueryInfo *QueryParam `json:"queryInfo"`
}

// QueryUserTagRequest 查询用户标签请求
//
//	@author centonhuang
//	@update 2025-01-05 12:37:54
type QueryUserTagRequest struct {
	UserName   string
	QueryParam *QueryParam
}

// QueryUserTagResponse 查询用户标签响应
//
//	@author centonhuang
//	@update 2025-01-05 12:37:54
type QueryUserTagResponse struct {
	Tags      []*Tag      `json:"tags"`
	QueryInfo *QueryParam `json:"queryInfo"`
}

// Category 分类
//
//	@author centonhuang
//	@update 2025-01-05 13:22:49
type Category struct {
	CategoryID uint   `json:"categoryID"`
	Name       string `json:"name"`
	ParentID   uint   `json:"parentID,omitempty"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

// Article 文章
//
//	@author centonhuang
//	@update 2025-01-05 15:23:26
type Article struct {
	ArticleID   uint     `json:"articleID"`
	Title       string   `json:"title"`
	Slug        string   `json:"slug"`
	Status      string   `json:"status"`
	Author      string   `json:"author"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
	PublishedAt string   `json:"publishedAt"`
	Likes       uint     `json:"likes"`
	Views       uint     `json:"views"`
	Tags        []string `json:"tags"`
	Comments    int      `json:"comments"`
}

// CreateCategoryRequest 创建分类请求
//
//	@author centonhuang
//	@update 2025-01-05 13:22:59
type CreateCategoryRequest struct {
	CurUserName string `json:"curUserName"`
	UserName    string `json:"userName"`
	Name        string `json:"name"`
	ParentID    uint   `json:"parentID"`
}

// CreateCategoryResponse 创建分类响应
//
//	@author centonhuang
//	@update 2025-01-05 13:23:01
type CreateCategoryResponse struct {
	Category *Category `json:"category"`
}

// GetCategoryInfoRequest 获取分类信息请求
//
//	@author centonhuang
//	@update 2025-01-05 13:23:03
type GetCategoryInfoRequest struct {
	CurUserName string `json:"curUserName"`
	UserName    string `json:"userName"`
	CategoryID  uint   `json:"categoryID"`
}

// GetCategoryInfoResponse 获取分类信息响应
//
//	@author centonhuang
//	@update 2025-01-05 13:23:06
type GetCategoryInfoResponse struct {
	Category *Category `json:"category"`
}

// GetRootCategoryRequest 获取根分类请求
//
//	@author centonhuang
//	@update 2025-01-05 13:23:08
type GetRootCategoryRequest struct {
	CurUserName string `json:"curUserName"`
	UserName    string `json:"userName"`
}

// GetRootCategoryResponse 获取根分类响应
//
//	@author centonhuang
//	@update 2025-01-05 13:23:13
type GetRootCategoryResponse struct {
	Category *Category `json:"category"`
}

// UpdateCategoryRequest 更新分类请求
//
//	@author centonhuang
//	@update 2025-01-05 13:23:14
type UpdateCategoryRequest struct {
	CurUserName string `json:"curUserName"`
	UserName    string `json:"userName"`
	CategoryID  uint   `json:"categoryID"`
	Name        string `json:"name"`
	ParentID    uint   `json:"parentID"`
}

// UpdateCategoryResponse 更新分类响应
//
//	@author centonhuang
//	@update 2025-01-05 13:23:16
type UpdateCategoryResponse struct {
	Category *Category `json:"category"`
}

// DeleteCategoryRequest 删除分类请求
//
//	@author centonhuang
//	@update 2025-01-05 13:23:18
type DeleteCategoryRequest struct {
	CurUserName string `json:"curUserName"`
	UserName    string `json:"userName"`
	CategoryID  uint   `json:"categoryID"`
}

// DeleteCategoryResponse 删除分类响应
//
//	@author centonhuang
//	@update 2025-01-05 13:23:20
type DeleteCategoryResponse struct{}

// ListChildrenCategoriesRequest 列出子分类请求
//
//	@author centonhuang
//	@update 2025-01-05 13:23:21
type ListChildrenCategoriesRequest struct {
	CurUserName string     `json:"curUserName"`
	UserName    string     `json:"userName"`
	CategoryID  uint       `json:"categoryID"`
	PageParam   *PageParam `json:"pageParam"`
}

// ListChildrenCategoriesResponse 列出子分类响应
//
//	@author centonhuang
//	@update 2025-01-05 13:23:23
type ListChildrenCategoriesResponse struct {
	Categories []*Category `json:"categories"`
	PageInfo   *PageInfo   `json:"pageInfo"`
}

// ListChildrenArticlesRequest 列出子文章请求
//
//	@author centonhuang
//	@update 2025-01-05 13:23:25
type ListChildrenArticlesRequest struct {
	CurUserName string     `json:"curUserName"`
	UserName    string     `json:"userName"`
	CategoryID  uint       `json:"categoryID"`
	PageParam   *PageParam `json:"pageParam"`
}

// ListChildrenArticlesResponse 列出子文章响应
//
//	@author centonhuang
//	@update 2025-01-05 13:23:26
type ListChildrenArticlesResponse struct {
	Articles []*Article `json:"articles"`
	PageInfo *PageInfo  `json:"pageInfo"`
}

// LoginRequest OAuth2登录请求
//
//	@author centonhuang
//	@update 2025-01-05 14:23:26
type LoginRequest struct{}

// LoginResponse OAuth2登录响应
//
//	@author centonhuang
//	@update 2025-01-05 14:23:26
type LoginResponse struct {
	RedirectURL string `json:"redirectURL"`
}

// CallbackRequest OAuth2回调请求
//
//	@author centonhuang
//	@update 2025-01-05 14:23:26
type CallbackRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

// CallbackResponse OAuth2回调响应
//
//	@author centonhuang
//	@update 2025-01-05 14:23:26
type CallbackResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// CreateArticleRequest 创建文章请求
//
//	@author centonhuang
//	@update 2025-01-05 15:23:26
type CreateArticleRequest struct {
	CurUserName string   `json:"curUserName"`
	UserName    string   `json:"userName"`
	UserID      uint     `json:"userID"`
	Title       string   `json:"title"`
	Slug        string   `json:"slug"`
	CategoryID  uint     `json:"categoryID"`
	Tags        []string `json:"tags"`
}

// CreateArticleResponse 创建文章响应
//
//	@author centonhuang
//	@update 2025-01-05 15:23:26
type CreateArticleResponse struct{}

// GetArticleInfoRequest 获取文章信息请求
//
//	@author centonhuang
//	@update 2025-01-05 15:23:26
type GetArticleInfoRequest struct {
	UserName    string `json:"userName"`
	ArticleSlug string `json:"articleSlug"`
}

// GetArticleInfoResponse 获取文章信息响应
//
//	@author centonhuang
//	@update 2025-01-05 15:23:26
type GetArticleInfoResponse struct {
	Article *Article `json:"article"`
}

// UpdateArticleRequest 更新文章请求
//
//	@author centonhuang
//	@update 2025-01-05 15:23:26
type UpdateArticleRequest struct {
	CurUserName       string `json:"curUserName"`
	UserName          string `json:"userName"`
	ArticleSlug       string `json:"articleSlug"`
	UpdatedTitle      string `json:"title"`
	UpdatedSlug       string `json:"slug"`
	UpdatedCategoryID uint   `json:"categoryID"`
}

// UpdateArticleResponse 更新文章响应
//
//	@author centonhuang
//	@update 2025-01-05 15:23:26
type UpdateArticleResponse struct{}

// UpdateArticleStatusRequest 更新文章状态请求
//
//	@author centonhuang
//	@update 2025-01-05 15:23:26
type UpdateArticleStatusRequest struct {
	CurUserName string              `json:"curUserName"`
	UserName    string              `json:"userName"`
	ArticleSlug string              `json:"articleSlug"`
	Status      model.ArticleStatus `json:"status"`
}

// UpdateArticleStatusResponse 更新文章状态响应
//
//	@author centonhuang
//	@update 2025-01-05 15:23:26
type UpdateArticleStatusResponse struct{}

// DeleteArticleRequest 删除文章请求
//
//	@author centonhuang
//	@update 2025-01-05 15:23:26
type DeleteArticleRequest struct {
	CurUserName string `json:"curUserName"`
	UserName    string `json:"userName"`
	ArticleSlug string `json:"articleSlug"`
}

// DeleteArticleResponse 删除文章响应
//
//	@author centonhuang
//	@update 2025-01-05 15:23:26
type DeleteArticleResponse struct{}

// ListArticlesRequest 列出文章请求
//
//	@author centonhuang
//	@update 2025-01-05 15:23:26
type ListArticlesRequest struct {
	PageParam *PageParam `json:"pageParam"`
}

// ListArticlesResponse 列出文章响应
//
//	@author centonhuang
//	@update 2025-01-05 15:23:26
type ListArticlesResponse struct {
	Articles []*Article `json:"articles"`
	PageInfo *PageInfo  `json:"pageInfo"`
}

// ListUserArticlesRequest 列出用户文章请求
//
//	@author centonhuang
//	@update 2025-01-05 15:23:26
type ListUserArticlesRequest struct {
	UserName  string     `json:"userName"`
	PageParam *PageParam `json:"pageParam"`
}

// ListUserArticlesResponse 列出用户文章响应
//
//	@author centonhuang
//	@update 2025-01-05 15:23:26
type ListUserArticlesResponse struct {
	Articles []*Article `json:"articles"`
	PageInfo *PageInfo  `json:"pageInfo"`
}

// QueryArticleRequest 查询文章请求
//
//	@author centonhuang
//	@update 2025-01-05 15:23:26
type QueryArticleRequest struct {
	QueryParam *QueryParam `json:"queryParam"`
}

// QueryArticleResponse 查询文章响应
//
//	@author centonhuang
//	@update 2025-01-05 15:23:26
type QueryArticleResponse struct {
	QueryInfo *QueryParam `json:"queryInfo"`
}

// QueryUserArticleRequest 查询用户文章请求
//
//	@author centonhuang
//	@update 2025-01-05 15:23:26
type QueryUserArticleRequest struct {
	UserName   string      `json:"userName"`
	QueryParam *QueryParam `json:"queryParam"`
}

// QueryUserArticleResponse 查询用户文章响应
//
//	@author centonhuang
//	@update 2025-01-05 15:23:26
type QueryUserArticleResponse struct {
	QueryInfo *QueryParam `json:"queryInfo"`
}

// ArticleVersion 文章版本
type ArticleVersion struct {
	ArticleVersionID uint   `json:"versionID"`
	ArticleID        uint   `json:"articleID"`
	Version          uint   `json:"version"`
	Content          string `json:"content"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
}

// CreateArticleVersionRequest 创建文章版本请求
//
//	@author centonhuang
//	@update 2025-01-05 16:42:48
type CreateArticleVersionRequest struct {
	CurUserName string `json:"curUserName"`
	UserName    string `json:"userName"`
	ArticleSlug string `json:"articleSlug"`
	Content     string `json:"content"`
}

// CreateArticleVersionResponse 创建文章版本响应
//
//	@author centonhuang
//	@update 2025-01-05 16:42:48
type CreateArticleVersionResponse struct{}

// GetArticleVersionInfoRequest 获取文章版本信息请求
//
//	@author centonhuang
//	@update 2025-01-05 16:42:48
type GetArticleVersionInfoRequest struct {
	CurUserName string `json:"curUserName"`
	UserName    string `json:"userName"`
	ArticleSlug string `json:"articleSlug"`
	Version     uint   `json:"version"`
}

// GetArticleVersionInfoResponse 获取文章版本信息响应
//
//	@author centonhuang
//	@update 2025-01-05 16:42:48
type GetArticleVersionInfoResponse struct {
	Version *ArticleVersion `json:"version"`
}

// GetLatestArticleVersionInfoRequest 获取最新文章版本信息请求
//
//	@author centonhuang
//	@update 2025-01-05 16:42:48
type GetLatestArticleVersionInfoRequest struct {
	CurUserName string `json:"curUserName"`
	UserName    string `json:"userName"`
	ArticleSlug string `json:"articleSlug"`
}

// GetLatestArticleVersionInfoResponse 获取最新文章版本信息响应
//
//	@author centonhuang
//	@update 2025-01-05 16:42:48
type GetLatestArticleVersionInfoResponse struct {
	Version *ArticleVersion `json:"version"`
}

// ListArticleVersionsRequest 列出文章版本请求
//
//	@author centonhuang
//	@update 2025-01-05 16:42:48
type ListArticleVersionsRequest struct {
	CurUserName string     `json:"curUserName"`
	UserName    string     `json:"userName"`
	ArticleSlug string     `json:"articleSlug"`
	PageParam   *PageParam `json:"pageParam"`
}

// ListArticleVersionsResponse 列出文章版本响应
//
//	@author centonhuang
//	@update 2025-01-05 16:42:48
type ListArticleVersionsResponse struct {
	Versions []*ArticleVersion `json:"versions"`
	PageInfo *PageInfo         `json:"pageInfo"`
}

// LikeArticleRequest 点赞文章请求
//
//	@author centonhuang
//	@update 2025-01-05 16:42:48
type LikeArticleRequest struct {
	CurUserID   uint   `json:"curUserID"`
	Author      string `json:"author"`
	ArticleSlug string `json:"articleSlug"`
	Undo        bool   `json:"undo"`
}

// LikeArticleResponse 点赞文章响应
//
//	@author centonhuang
//	@update 2025-01-05 16:43:09
type LikeArticleResponse struct{}

// LikeCommentRequest 点赞评论请求
//
//	@author centonhuang
//	@update 2025-01-05 16:43:20
type LikeCommentRequest struct {
	CurUserID uint `json:"curUserID"`
	CommentID uint `json:"commentID"`
	Undo      bool `json:"undo"`
}

// LikeCommentResponse 点赞评论响应
//
//	@author centonhuang
//	@update 2025-01-05 16:43:21
type LikeCommentResponse struct{}

// LikeTagRequest 点赞标签请求
//
//	@author centonhuang
//	@update 2025-01-05 16:43:23
type LikeTagRequest struct {
	CurUserID uint   `json:"curUserID"`
	TagSlug   string `json:"tagSlug"`
	Undo      bool   `json:"undo"`
}

// LikeTagResponse 点赞标签响应
//
//	@author centonhuang
//	@update 2025-01-05 16:43:25
type LikeTagResponse struct{}

// LogArticleViewRequest 记录文章浏览请求
//
//	@author centonhuang
//	@update 2025-01-05 16:43:26
type LogArticleViewRequest struct {
	CurUserID   uint   `json:"curUserID"`
	Author      string `json:"author"`
	ArticleSlug string `json:"articleSlug"`
	Progress    int8   `json:"progress"`
}

// LogArticleViewResponse 记录文章浏览响应
//
//	@author centonhuang
//	@update 2025-01-05 16:43:28
type LogArticleViewResponse struct{}

// Comment 评论
//
//	@author centonhuang
//	@update 2025-01-05 16:43:29
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
//	@author centonhuang
//	@update 2025-01-05 16:43:31
type CreateArticleCommentRequest struct {
	CurUserID   uint   `json:"curUserID"`
	Author      string `json:"author"`
	ArticleSlug string `json:"articleSlug"`
	Content     string `json:"content"`
	ReplyTo     uint   `json:"replyTo"`
}

// CreateArticleCommentResponse 创建文章评论响应
//
//	@author centonhuang
//	@update 2025-01-05 16:43:33
type CreateArticleCommentResponse struct {
	Comment *Comment `json:"comment"`
}

// GetCommentInfoRequest 获取评论信息请求
//
//	@author centonhuang
//	@update 2025-01-05 16:43:34
type GetCommentInfoRequest struct {
	UserName    string `json:"userName"`
	ArticleSlug string `json:"articleSlug"`
	CommentID   uint   `json:"commentID"`
}

// GetCommentInfoResponse 获取评论信息响应
//
//	@author centonhuang
//	@update 2025-01-05 16:43:36
type GetCommentInfoResponse struct {
	Comment *Comment `json:"comment"`
}

// DeleteCommentRequest 删除评论请求
//
//	@author centonhuang
//	@update 2025-01-05 16:43:38
type DeleteCommentRequest struct {
	CurUserName string `json:"curUserName"`
	UserName    string `json:"userName"`
	ArticleSlug string `json:"articleSlug"`
	CommentID   uint   `json:"commentID"`
}

// DeleteCommentResponse 删除评论响应
//
//	@author centonhuang
//	@update 2025-01-05 16:43:39
type DeleteCommentResponse struct{}

// ListArticleCommentsRequest 列出文章评论请求
//
//	@author centonhuang
//	@update 2025-01-05 16:43:41
type ListArticleCommentsRequest struct {
	CurUserName string     `json:"curUserName"`
	UserName    string     `json:"userName"`
	ArticleSlug string     `json:"articleSlug"`
	PageParam   *PageParam `json:"pageParam"`
}

// ListArticleCommentsResponse 列出文章评论响应
//
//	@author centonhuang
//	@update 2025-01-05 16:43:43
type ListArticleCommentsResponse struct {
	Comments []*Comment `json:"comments"`
	PageInfo *PageInfo  `json:"pageInfo"`
}

// ListChildrenCommentsRequest 列出子评论请求
//
//	@author centonhuang
//	@update 2025-01-05 16:43:44
type ListChildrenCommentsRequest struct {
	UserName    string     `json:"userName"`
	ArticleSlug string     `json:"articleSlug"`
	CommentID   uint       `json:"commentID"`
	PageParam   *PageParam `json:"pageParam"`
}

// ListChildrenCommentsResponse 列出子评论响应
//
//	@author centonhuang
//	@update 2025-01-05 16:43:46
type ListChildrenCommentsResponse struct {
	Comments []*Comment `json:"comments"`
	PageInfo *PageInfo  `json:"pageInfo"`
}

// ListUserLikeArticlesRequest 列出用户喜欢的文章请求
//
//	@author centonhuang
//	@update 2025-01-05 16:43:48
type ListUserLikeArticlesRequest struct {
	CurUserID uint       `json:"curUserID"`
	PageParam *PageParam `json:"pageParam"`
}

// ListUserLikeArticlesResponse 列出用户喜欢的文章响应
//
//	@author centonhuang
//	@update 2025-01-05 16:43:50
type ListUserLikeArticlesResponse struct {
	Articles []*Article `json:"articles"`
	PageInfo *PageInfo  `json:"pageInfo"`
}

// ListUserLikeCommentsRequest 列出用户喜欢的评论请求
//
//	@author centonhuang
//	@update 2025-01-05 16:43:52
type ListUserLikeCommentsRequest struct {
	CurUserID uint       `json:"curUserID"`
	PageParam *PageParam `json:"pageParam"`
}

// ListUserLikeCommentsResponse 列出用户喜欢的评论响应
//
//	@author centonhuang
//	@update 2025-01-05 16:43:54
type ListUserLikeCommentsResponse struct {
	Comments []*Comment `json:"comments"`
	PageInfo *PageInfo  `json:"pageInfo"`
}

// ListUserLikeTagsRequest 列出用户喜欢的标签请求
//
//	@author centonhuang
//	@update 2025-01-05 16:43:56
type ListUserLikeTagsRequest struct {
	CurUserID uint       `json:"curUserID"`
	PageParam *PageParam `json:"pageParam"`
}

// ListUserLikeTagsResponse 列出用户喜欢的标签响应
//
//	@author centonhuang
//	@update 2025-01-05 16:43:58
type ListUserLikeTagsResponse struct {
	Tags     []*Tag    `json:"tags"`
	PageInfo *PageInfo `json:"pageInfo"`
}

// CreateBucketRequest 创建桶请求
//
//	@author centonhuang
//	@update 2025-01-05 17:03:19
type CreateBucketRequest struct {
	CurUserID uint `json:"curUserID"`
}

// CreateBucketResponse 创建桶响应
//
//	@author centonhuang
//	@update 2025-01-05 17:03:21
type CreateBucketResponse struct{}

// Image 图片
//
//	@author centonhuang
//	@update 2025-01-05 17:17:53
type Image struct {
	Name      string `json:"name"`
	Size      int64  `json:"size"`
	CreatedAt string `json:"createdAt"`
}

// ListImagesRequest 列出图片请求
//
//	@author centonhuang
//	@update 2025-01-05 17:03:22
type ListImagesRequest struct {
	CurUserID uint `json:"curUserID"`
}

// ListImagesResponse 列出图片响应
//
//	@author centonhuang
//	@update 2025-01-05 17:03:24
type ListImagesResponse struct {
	Images []*Image `json:"images"`
}

// UploadImageRequest 上传图片请求
//
//	@author centonhuang
//	@update 2025-01-05 17:03:25
type UploadImageRequest struct {
	CurUserID   uint   `json:"curUserID"`
	FileName    string `json:"fileName"`
	Size        int64  `json:"size"`
	ContentType string `json:"contentType"`
	ReadSeeker  io.ReadSeeker
}

// UploadImageResponse 上传图片响应
//
//	@author centonhuang
//	@update 2025-01-05 17:03:27
type UploadImageResponse struct{}

// GetImageRequest 获取图片请求
//
//	@author centonhuang
//	@update 2025-01-05 17:03:28
type GetImageRequest struct {
	CurUserID uint   `json:"curUserID"`
	ImageName string `json:"imageName"`
	Quality   string `json:"quality"`
}

// GetImageResponse 获取图片响应
//
//	@author centonhuang
//	@update 2025-01-05 17:03:30
type GetImageResponse struct {
	PresignedURL string `json:"presignedURL"`
}

// DeleteImageRequest 删除图片请求
//
//	@author centonhuang
//	@update 2025-01-05 17:03:31
type DeleteImageRequest struct {
	CurUserID uint   `json:"curUserID"`
	ImageName string `json:"imageName"`
}

// DeleteImageResponse 删除图片响应
//
//	@author centonhuang
//	@update 2025-01-05 17:03:33
type DeleteImageResponse struct{}

// UserView 用户浏览
//
//	@author centonhuang
//	@update 2025-01-05 17:53:09
type UserView struct {
	ID           uint   `json:"id"`
	Progress     int8   `json:"progress"`
	LastViewedAt string `json:"lastViewedAt"`
	UserID       uint   `json:"userID"`
	ArticleID    uint   `json:"articleID"`
}

// ListUserViewArticlesRequest 列出用户浏览的文章请求
//
//	@author centonhuang
//	@update 2025-01-05 17:03:38
type ListUserViewArticlesRequest struct {
	CurUserID uint       `json:"curUserID"`
	PageParam *PageParam `json:"pageParam"`
}

// ListUserViewArticlesResponse 列出用户浏览的文章响应
//
//	@author centonhuang
//	@update 2025-01-05 17:03:40
type ListUserViewArticlesResponse struct {
	UserViews []*UserView `json:"userViews"`
	PageInfo  *PageInfo   `json:"pageInfo"`
}

// DeleteUserViewRequest 删除用户浏览的文章请求
//
//	@author centonhuang
//	@update 2025-01-05 17:03:41
type DeleteUserViewRequest struct {
	CurUserID uint `json:"curUserID"`
	ViewID    uint `json:"viewID"`
}

// DeleteUserViewResponse 删除用户浏览的文章响应
//
//	@author centonhuang
//	@update 2025-01-05 17:03:43
type DeleteUserViewResponse struct{}

// Prompt 提示词
//
//	@author centonhuang
//	@update 2025-01-05 18:07:58
type Prompt struct {
	ID        uint       `json:"id"`
	CreatedAt string     `json:"createdAt"`
	Task      string     `json:"task"`
	Version   uint       `json:"version"`
	Templates []Template `json:"templates"`
	Variables []string   `json:"variables"`
}

// GetPromptRequest 获取提示词请求
//
//	@author centonhuang
//	@update 2025-01-05 18:01:30
type GetPromptRequest struct {
	TaskName string `json:"taskName"`
	Version  uint   `json:"version"`
}

// GetPromptResponse 获取提示词响应
//
//	@author centonhuang
//	@update 2025-01-05 18:01:30
type GetPromptResponse struct {
	Prompt *Prompt `json:"prompt"`
}

// GetLatestPromptRequest 获取最新提示词请求
//
//	@author centonhuang
//	@update 2025-01-05 18:10:46
type GetLatestPromptRequest struct {
	TaskName string `json:"taskName"`
}

// GetLatestPromptResponse 获取最新提示词响应
//
//	@author centonhuang
//	@update 2025-01-05 18:10:51
type GetLatestPromptResponse struct {
	Prompt *Prompt `json:"prompt"`
}

// ListPromptRequest 列出提示词请求
//
//	@author centonhuang
//	@update 2025-01-05 18:13:20
type ListPromptRequest struct {
	TaskName  string     `json:"taskName"`
	PageParam *PageParam `json:"pageParam"`
}

// ListPromptResponse 列出提示词响应
//
//	@author centonhuang
//	@update 2025-01-05 18:13:22
type ListPromptResponse struct {
	Prompts  []*Prompt `json:"prompts"`
	PageInfo *PageInfo `json:"pageInfo"`
}

// CreatePromptRequest 创建提示词请求
//
//	@author centonhuang
//	@update 2025-01-05 18:16:13
type CreatePromptRequest struct {
	TaskName  string     `json:"taskName"`
	Templates []Template `json:"templates"`
}

// CreatePromptResponse 创建提示词响应
//
//	@author centonhuang
//	@update 2025-01-05 18:16:13
type CreatePromptResponse struct{}

// GenerateContentCompletionRequest 生成内容完成请求
//
//	@author centonhuang
//	@update 2025-01-05 18:41:32
type GenerateContentCompletionRequest struct {
	CurUserID   uint    `json:"curUserID"`
	Context     string  `json:"context"`
	Instruction string  `json:"instruction"`
	Reference   string  `json:"reference"`
	Temperature float64 `json:"temperature"`
}

// GenerateContentCompletionResponse 生成内容完成响应
//
//	@author centonhuang
//	@update 2025-01-05 18:41:32
type GenerateContentCompletionResponse struct {
	TokenChan <-chan string
	ErrChan   <-chan error
}

// GenerateArticleSummaryRequest 生成文章总结请求
//
//	@author centonhuang
//	@update 2025-01-05 18:41:32
type GenerateArticleSummaryRequest struct {
	CurUserID   uint    `json:"curUserID"`
	ArticleSlug string  `json:"articleSlug"`
	Instruction string  `json:"instruction"`
	Temperature float64 `json:"temperature"`
}

// GenerateArticleSummaryResponse 生成文章总结响应
//
//	@author centonhuang
//	@update 2025-01-05 18:41:32
type GenerateArticleSummaryResponse struct {
	TokenChan <-chan string
	ErrChan   <-chan error
}

// GenerateArticleTranslationRequest 生成文章翻译请求
//
//	@author centonhuang
//	@update 2025-01-05 20:36:40
type GenerateArticleTranslationRequest struct{}

// GenerateArticleTranslationResponse 生成文章翻译响应
//
//	@author centonhuang
//	@update 2025-01-05 20:36:43
type GenerateArticleTranslationResponse struct{}

// GenerateArticleQARequest 生成文章问答请求
//
//	@author centonhuang
//	@update 2025-01-05 18:41:32
type GenerateArticleQARequest struct {
	CurUserID   uint    `json:"curUserID"`
	ArticleSlug string  `json:"articleSlug"`
	Question    string  `json:"question"`
	Temperature float64 `json:"temperature"`
}

// GenerateArticleQAResponse 生成文章问答响应
//
//	@author centonhuang
//	@update 2025-01-05 18:41:32
type GenerateArticleQAResponse struct {
	TokenChan <-chan string
	ErrChan   <-chan error
}

// GenerateTermExplainationRequest 生成术语解释请求
//
//	@author centonhuang
//	@update 2025-01-05 18:41:32
type GenerateTermExplainationRequest struct {
	CurUserID   uint    `json:"curUserID"`
	ArticleSlug string  `json:"articleSlug"`
	Term        string  `json:"term"`
	Position    uint    `json:"position"`
	Temperature float64 `json:"temperature"`
}

// GenerateTermExplainationResponse 生成术语解释响应
//
//	@author centonhuang
//	@update 2025-01-05 18:41:32
type GenerateTermExplainationResponse struct {
	TokenChan <-chan string
	ErrChan   <-chan error
}
