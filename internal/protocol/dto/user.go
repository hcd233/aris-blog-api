// Package dto 用户DTO
package dto

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
	UserID uint `json:"userID" path:"userID" doc:"用户ID"`
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
	Body *UpdateUserInfoBody `json:"body"`
}

// UpdateUserInfoBody 更新用户信息请求体
//
//	author centonhuang
//	update 2025-10-31 02:33:48
type UpdateUserInfoBody struct {
	UserName string `json:"userName" doc:"更新后的用户名"`
}
