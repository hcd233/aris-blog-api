package protocol

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

// HumaPingRequest 健康检查请求
type HumaPingRequest struct{}

// HumaPingResponse 健康检查响应
type HumaPingResponse struct {
	Body struct {
		Status string `json:"status" example:"ok" doc:"服务状态"`
	}
}

// HumaGetCurUserInfoRequest 获取当前用户信息请求
type HumaGetCurUserInfoRequest struct {
	// 从 JWT 中间件获取，不需要显式参数
}

// HumaGetCurUserInfoResponse 获取当前用户信息响应
type HumaGetCurUserInfoResponse struct {
	Body struct {
		User *CurUser `json:"user" doc:"当前用户信息"`
	}
}

// HumaGetUserInfoRequest 获取用户信息请求
type HumaGetUserInfoRequest struct {
	UserID uint `path:"userID" minimum:"1" doc:"用户ID"`
}

// HumaGetUserInfoResponse 获取用户信息响应
type HumaGetUserInfoResponse struct {
	Body struct {
		User *User `json:"user" doc:"用户信息"`
	}
}

// HumaUpdateUserInfoRequest 更新用户信息请求
type HumaUpdateUserInfoRequest struct {
	Body struct {
		UserName string `json:"userName" minLength:"1" maxLength:"50" doc:"用户名"`
	}
}

// HumaUpdateUserInfoResponse 更新用户信息响应
type HumaUpdateUserInfoResponse struct {
	Body struct {
		Success bool `json:"success" example:"true" doc:"更新是否成功"`
	}
}

// HumaCreateTagRequest 创建标签请求
type HumaCreateTagRequest struct {
	Body struct {
		Name        string `json:"name" minLength:"1" maxLength:"100" doc:"标签名称"`
		Slug        string `json:"slug" minLength:"1" maxLength:"100" doc:"标签别名"`
		Description string `json:"description" maxLength:"500" doc:"标签描述"`
	}
}

// HumaCreateTagResponse 创建标签响应
type HumaCreateTagResponse struct {
	Body struct {
		Tag *Tag `json:"tag" doc:"创建的标签信息"`
	}
}

// HumaGetTagInfoRequest 获取标签信息请求
type HumaGetTagInfoRequest struct {
	TagID uint `path:"tagID" minimum:"1" doc:"标签ID"`
}

// HumaGetTagInfoResponse 获取标签信息响应
type HumaGetTagInfoResponse struct {
	Body struct {
		Tag *Tag `json:"tag" doc:"标签信息"`
	}
}

// HumaUpdateTagRequest 更新标签请求
type HumaUpdateTagRequest struct {
	TagID uint `path:"tagID" minimum:"1" doc:"标签ID"`
	Body  struct {
		Name        string `json:"name" minLength:"1" maxLength:"100" doc:"标签名称"`
		Slug        string `json:"slug" minLength:"1" maxLength:"100" doc:"标签别名"`
		Description string `json:"description" maxLength:"500" doc:"标签描述"`
	}
}

// HumaUpdateTagResponse 更新标签响应
type HumaUpdateTagResponse struct {
	Body struct {
		Success bool `json:"success" example:"true" doc:"更新是否成功"`
	}
}

// HumaDeleteTagRequest 删除标签请求
type HumaDeleteTagRequest struct {
	TagID uint `path:"tagID" minimum:"1" doc:"标签ID"`
}

// HumaDeleteTagResponse 删除标签响应
type HumaDeleteTagResponse struct {
	Body struct {
		Success bool `json:"success" example:"true" doc:"删除是否成功"`
	}
}

// HumaListTagsRequest 列出标签请求
type HumaListTagsRequest struct {
	Page     int    `query:"page" minimum:"1" default:"1" doc:"页码"`
	PageSize int    `query:"pageSize" minimum:"1" maximum:"50" default:"10" doc:"每页大小"`
	Query    string `query:"query" maxLength:"100" doc:"搜索查询"`
}

// HumaListTagsResponse 列出标签响应
type HumaListTagsResponse struct {
	Body struct {
		Tags     []*Tag    `json:"tags" doc:"标签列表"`
		PageInfo *PageInfo `json:"pageInfo" doc:"分页信息"`
	}
}

// HumaRefreshTokenRequest 刷新令牌请求
type HumaRefreshTokenRequest struct {
	Body struct {
		RefreshToken string `json:"refreshToken" minLength:"1" doc:"刷新令牌"`
	}
}

// HumaRefreshTokenResponse 刷新令牌响应
type HumaRefreshTokenResponse struct {
	Body struct {
		AccessToken  string `json:"accessToken" doc:"访问令牌"`
		RefreshToken string `json:"refreshToken" doc:"新的刷新令牌"`
	}
}

// HumaErrorResponse 错误响应
type HumaErrorResponse struct {
	Body struct {
		Error string `json:"error" doc:"错误信息"`
	}
}

// ResolveUserID 从上下文中解析用户ID的辅助函数
func ResolveUserID(ctx context.Context) (uint, error) {
	// 从上下文中获取用户ID
	// 这个值由 Huma JWT 中间件设置
	if userID := ctx.Value("userID"); userID != nil {
		if id, ok := userID.(uint); ok {
			return id, nil
		}
	}
	return 0, huma.Error401Unauthorized("用户未认证")
}