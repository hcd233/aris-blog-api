// Package dto 用户DTO
package dto

// GetCurUserInfoRequest represents a request to get the current authenticated user's information
//
//	author centonhuang
//	update 2025-01-04 21:00:54
type GetCurUserInfoRequest struct {
	UserID uint `json:"userID" doc:"User ID extracted from the JWT token"`
}

// GetCurUserInfoResponse represents the response containing the current user's detailed information
//
//	author centonhuang
//	update 2025-01-04 21:00:59
type GetCurUserInfoResponse struct {
	User *User `json:"user" doc:"Complete user information including permissions"`
}

// GetUserInfoRequest represents a request to get a specific user's public information
//
//	author centonhuang
//	update 2025-01-04 21:19:41
type GetUserInfoRequest struct {
	UserID uint `json:"userID" path:"userID" doc:"Unique identifier of the user to retrieve"`
}

// GetUserResponse represents the response containing a user's public information
//
//	author centonhuang
//	update 2025-01-04 21:19:44
type GetUserResponse struct {
	User *User `json:"user" doc:"Public user information"`
}

// UpdateUserRequest represents a request to update the current user's information
//
//	author centonhuang
//	update 2025-01-04 21:19:47
type UpdateUserRequest struct {
	Body *UpdateUserRequestBody `json:"body" doc:"Request body containing fields to update"`
}

// UpdateUserRequestBody contains the fields that can be updated for a user
//
//	author centonhuang
//	update 2025-10-31 02:33:48
type UpdateUserRequestBody struct {
	UserName string `json:"userName" doc:"New display name for the user"`
}
