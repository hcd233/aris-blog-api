// Package dto OAuth2 DTO
package dto

// LoginRequest represents a request to initiate OAuth2 login flow
//
//	author centonhuang
//	update 2025-01-05 21:00:00
type LoginRequest struct {
	Provider string `json:"provider" path:"provider" enum:"github,google" doc:"OAuth2 provider name (github or google)"`
}

// LoginResponse represents the response containing the OAuth2 authorization URL
//
//	author centonhuang
//	update 2025-01-05 21:00:00
type LoginResponse struct {
	RedirectURL string `json:"redirectURL" doc:"URL to redirect the user to for OAuth2 authorization"`
}

// CallbackRequest represents a request to handle OAuth2 callback with authorization code
//
//	author centonhuang
//	update 2025-01-05 21:00:00
type CallbackRequest struct {
	Provider string `json:"provider" path:"provider" enum:"github,google" doc:"OAuth2 provider name (github or google)"`
	Code     string `json:"code" query:"code" doc:"Authorization code returned by the OAuth2 provider"`
	State    string `json:"state" query:"state" doc:"State parameter for CSRF protection, must match the initial state"`
}

// CallbackResponse represents the response containing access and refresh tokens after successful OAuth2 authentication
//
//	author centonhuang
//	update 2025-01-05 21:00:00
type CallbackResponse struct {
	AccessToken  string `json:"accessToken" doc:"JWT access token for API authentication"`
	RefreshToken string `json:"refreshToken" doc:"JWT refresh token for obtaining future access tokens"`
}
