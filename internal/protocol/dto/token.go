// Package dto 令牌DTO
package dto

// RefreshTokenRequest represents a request to refresh an access token using a refresh token
//
//	author centonhuang
//	update 2025-01-05 21:00:00
type RefreshTokenRequest struct {
	Body *RefreshTokenBody `json:"body" doc:"Request body containing the refresh token"`
}

// RefreshTokenBody contains the refresh token used to obtain a new access token
//
//	author centonhuang
//	update 2025-01-05 21:00:00
type RefreshTokenBody struct {
	RefreshToken string `json:"refreshToken" doc:"JWT refresh token used to obtain a new access token"`
}

// RefreshTokenResponse represents the response containing new access and refresh tokens
//
//	author centonhuang
//	update 2025-01-05 21:00:00
type RefreshTokenResponse struct {
	AccessToken  string `json:"accessToken" doc:"New JWT access token for API authentication"`
	RefreshToken string `json:"refreshToken" doc:"New JWT refresh token for obtaining future access tokens"`
}
