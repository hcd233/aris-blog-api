package protocol

import (
	"context"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
)

// PaginationParams 分页参数（用于 huma）
type PaginationParams struct {
	Page     int `query:"page" doc:"Page number" default:"1" minimum:"1"`
	PageSize int `query:"pageSize" doc:"Page size" default:"10" minimum:"1" maximum:"50"`
}

// SearchParams 搜索参数（用于 huma）
type SearchParams struct {
	Query string `query:"query" doc:"Search query"`
}

// PaginatedSearchParams 分页搜索参数（用于 huma）
type PaginatedSearchParams struct {
	PaginationParams
	SearchParams
}

// UserInput 用户输入（用于 huma）
type UserInput struct {
	UserID uint `path:"userID" doc:"User ID"`
}

// UpdateUserInput 更新用户输入（用于 huma）
type UpdateUserInput struct {
	UserName string `json:"userName" doc:"User name" required:"true"`
}

// UserOutput 用户输出（用于 huma）
type UserOutput struct {
	Body User `json:"body"`
}

// CurUserOutput 当前用户输出（用于 huma）
type CurUserOutput struct {
	Body CurUser `json:"body"`
}

// UserListOutput 用户列表输出（用于 huma）
type UserListOutput struct {
	Body struct {
		Users    []User   `json:"users"`
		PageInfo PageInfo `json:"pageInfo"`
	} `json:"body"`
}

// TagInput 标签输入（用于 huma）
type TagInput struct {
	TagID uint `path:"tagID" doc:"Tag ID"`
}

// CreateTagInput 创建标签输入（用于 huma）
type CreateTagInput struct {
	Body struct {
		Name        string `json:"name" doc:"Tag name" required:"true"`
		Slug        string `json:"slug" doc:"Tag slug" required:"true"`
		Description string `json:"description" doc:"Tag description"`
	} `json:"body"`
}

// UpdateTagInput 更新标签输入（用于 huma）
type UpdateTagInput struct {
	TagID uint `path:"tagID" doc:"Tag ID"`
	Body  struct {
		Name        string `json:"name" doc:"Tag name" required:"true"`
		Slug        string `json:"slug" doc:"Tag slug" required:"true"`
		Description string `json:"description" doc:"Tag description"`
	} `json:"body"`
}

// TagOutput 标签输出（用于 huma）
type TagOutput struct {
	Body Tag `json:"body"`
}

// TagListOutput 标签列表输出（单用于 huma）
type TagListOutput struct {
	Body struct {
		Tags     []Tag    `json:"tags"`
		PageInfo PageInfo `json:"pageInfo"`
	} `json:"body"`
}

// RefreshTokenInput 刷新 token 输入（用于 huma）
type RefreshTokenInput struct {
	Body struct {
		RefreshToken string `json:"refreshToken" doc:"Refresh token" required:"true"`
	} `json:"body"`
}

// RefreshTokenOutput 刷新 token 输出（用于 huma）
type RefreshTokenOutput struct {
	Body RefreshTokenResponse `json:"body"`
}

// ErrorOutput 错误输出（用于 huma）
type ErrorOutput struct {
	huma.ErrorModel
}

// HTTPError 创建 HTTP 错误响应
func HTTPError(ctx context.Context, err error) error {
	// 将现有的错误映射到 huma 错误
	switch err {
	case ErrUnauthorized:
		return huma.Error401Unauthorized("Unauthorized")
	case ErrNoPermission:
		return huma.Error403Forbidden("No permission")
	case ErrDataNotExists:
		return huma.Error404NotFound("Data not found")
	case ErrDataExists:
		return huma.Error409Conflict("Data already exists")
	case ErrTooManyRequests:
		return huma.Error429TooManyRequests("Too many requests")
	case ErrBadRequest:
		return huma.Error400BadRequest("Bad request")
	case ErrInsufficientQuota:
		return huma.Error403Forbidden("Insufficient quota")
	case ErrNoImplement:
		return huma.Error501NotImplemented("Not implemented")
	default:
		return huma.Error500InternalServerError("Internal server error")
	}
}

// EmptyResponse 空响应（用于 huma）
type EmptyResponse struct {
	Body struct{}
}

// StandardHTTPResponse 标准 HTTP 响应
type StandardHTTPResponse[T any] struct {
	Data T `json:"data"`
}

// OperationResponse 操作响应（用于 huma）
type OperationResponse struct {
	StatusCode int
}

// SuccessResponse 成功响应（用于 huma）
type SuccessResponse struct {
	Body struct {
		Message string `json:"message"`
	} `json:"body"`
}

// GetTime 获取时间（用于 huma）
func GetTime() time.Time {
	return time.Now()
}

// ArticleInput 文章输入（用于 huma）
type ArticleInput struct {
	ArticleID uint `path:"articleID" doc:"Article ID"`
}

// ArticleSlugInput 文章别名输入（用于 huma）
type ArticleSlugInput struct {
	AuthorName  string `path:"authorName" doc:"Author name"`
	ArticleSlug string `path:"articleSlug" doc:"Article slug"`
}

// CreateArticleInput 创建文章输入（用于 huma）
type CreateArticleInput struct {
	Body struct {
		Title      string   `json:"title" doc:"Article title" required:"true"`
		Slug       string   `json:"slug" doc:"Article slug" required:"true"`
		Tags       []string `json:"tags" doc:"Article tags" required:"true"`
		CategoryID uint     `json:"categoryID" doc:"Category ID" required:"true"`
	} `json:"body"`
}

// UpdateArticleInput 更新文章输入（用于 huma）
type UpdateArticleInput struct {
	ArticleID uint `path:"articleID" doc:"Article ID"`
	Body      struct {
		Title      string `json:"title" doc:"Article title"`
		Slug       string `json:"slug" doc:"Article slug"`
		CategoryID uint   `json:"categoryID" doc:"Category ID"`
	} `json:"body"`
}

// UpdateArticleStatusInput 更新文章状态输入（用于 huma）
type UpdateArticleStatusInput struct {
	ArticleID uint `path:"articleID" doc:"Article ID"`
	Body      struct {
		Status string `json:"status" doc:"Article status" required:"true" enum:"draft,publish"`
	} `json:"body"`
}

// ArticleOutput 文章输出（用于 huma）
type ArticleOutput struct {
	Body Article `json:"body"`
}

// ArticleListOutput 文章列表输出（用于 huma）
type ArticleListOutput struct {
	Body struct {
		Articles []Article `json:"articles"`
		PageInfo PageInfo  `json:"pageInfo"`
	} `json:"body"`
}

// CategoryInput 分类输入（用于 huma）
type CategoryInput struct {
	CategoryID uint `path:"categoryID" doc:"Category ID"`
}

// CreateCategoryInput 创建分类输入（用于 huma）
type CreateCategoryInput struct {
	Body struct {
		ParentID uint   `json:"parentID" doc:"Parent category ID"`
		Name     string `json:"name" doc:"Category name" required:"true"`
	} `json:"body"`
}

// UpdateCategoryInput 更新分类输入（用于 huma）
type UpdateCategoryInput struct {
	CategoryID uint `path:"categoryID" doc:"Category ID"`
	Body       struct {
		Name     string `json:"name" doc:"Category name"`
		ParentID uint   `json:"parentID" doc:"Parent category ID"`
	} `json:"body"`
}

// CategoryOutput 分类输出（用于 huma）
type CategoryOutput struct {
	Body Category `json:"body"`
}

// CategoryListOutput 分类列表输出（用于 huma）
type CategoryListOutput struct {
	Body struct {
		Categories []Category `json:"categories"`
		PageInfo   PageInfo   `json:"pageInfo"`
	} `json:"body"`
}

// CommentInput 评论输入（用于 huma）
type CommentInput struct {
	CommentID uint `path:"commentID" doc:"Comment ID"`
}

// ArticleCommentInput 文章评论输入（用于 huma）
type ArticleCommentInput struct {
	ArticleID uint `path:"articleID" doc:"Article ID"`
}

// CreateCommentInput 创建评论输入（用于 huma）
type CreateCommentInput struct {
	Body struct {
		ArticleID uint   `json:"articleID" doc:"Article ID" required:"true"`
		ReplyTo   uint   `json:"replyTo" doc:"Reply to comment ID"`
		Content   string `json:"content" doc:"Comment content" required:"true" minLength:"1" maxLength:"300"`
	} `json:"body"`
}

// CommentOutput 评论输出（用于 huma）
type CommentOutput struct {
	Body Comment `json:"body"`
}

// CommentListOutput 评论列表输出（用于 huma）
type CommentListOutput struct {
	Body struct {
		Comments []Comment `json:"comments"`
		PageInfo PageInfo  `json:"pageInfo"`
	} `json:"body"`
}

// EmptyOutput 空输出（用于 huma）
type EmptyOutput struct {
	Body struct{}
}

// GetArticleStatusFromString 从字符串获取 ArticleStatus
func GetArticleStatusFromString(s string) model.ArticleStatus {
	switch s {
	case "publish":
		return model.ArticleStatusPublish
	default:
		return model.ArticleStatusDraft
	}
}
