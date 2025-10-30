// Package dto 数据传输对象
package dto

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
