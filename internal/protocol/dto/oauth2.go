// Package dto OAuth2 DTO
package dto

// LoginRequest OAuth2登录请求
//
//	author centonhuang
//	update 2025-01-05 21:00:00
type LoginRequest struct {
	Provider string `json:"provider" path:"provider" enum:"github,google" doc:"OAuth2提供商 (github/google)"`
}

// LoginResponse OAuth2登录响应
//
//	author centonhuang
//	update 2025-01-05 21:00:00
type LoginResponse struct {
	RedirectURL string `json:"redirectURL" doc:"重定向URL"`
}

// CallbackRequest OAuth2回调请求
//
//	author centonhuang
//	update 2025-01-05 21:00:00
type CallbackRequest struct {
	Provider string `json:"provider" path:"provider" doc:"OAuth2提供商 (github/google)"`
	Code     string `json:"code" query:"code" doc:"授权码"`
	State    string `json:"state" query:"state" doc:"状态码"`
}

// CallbackResponse OAuth2回调响应
//
//	author centonhuang
//	update 2025-01-05 21:00:00
type CallbackResponse struct {
	AccessToken  string `json:"accessToken" doc:"访问令牌"`
	RefreshToken string `json:"refreshToken" doc:"刷新令牌"`
}
