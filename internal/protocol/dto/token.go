// Package dto 令牌DTO
package dto

// RefreshTokenRequest 刷新令牌请求
//
//	author centonhuang
//	update 2025-01-05 21:00:00
type RefreshTokenRequest struct {
	Body *RefreshTokenBody `json:"body"`
}

// RefreshTokenBody 刷新令牌请求体
//
//	author centonhuang
//	update 2025-01-05 21:00:00
type RefreshTokenBody struct {
	RefreshToken string `json:"refreshToken" doc:"刷新令牌"`
}

// RefreshTokenResponse 刷新令牌响应
//
//	author centonhuang
//	update 2025-01-05 21:00:00
type RefreshTokenResponse struct {
	AccessToken  string `json:"accessToken" doc:"访问令牌"`
	RefreshToken string `json:"refreshToken" doc:"刷新令牌"`
}
