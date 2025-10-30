// Package dto 数据传输对象
package dto

// User represents a user entity
//
//	author centonhuang
//	update 2025-01-05 11:37:01
type User struct {
	UserID    uint   `json:"userID" doc:"Unique identifier for the user"`
	Name      string `json:"name" doc:"Display name of the user"`
	Email     string `json:"email,omitempty" doc:"Email address of the user"`
	Avatar    string `json:"avatar" doc:"URL or path to the user's avatar image"`
	CreatedAt string `json:"createdAt,omitempty" doc:"Timestamp when the user account was created"`
	LastLogin string `json:"lastLogin,omitempty" doc:"Timestamp of the user's last login"`
}

// CurUser represents the current authenticated user with additional permission information
//
//	author centonhuang
//	update 2025-01-05 11:37:32
type CurUser struct {
	User
	Permission string `json:"permission" doc:"Permission level of the user"`
}
