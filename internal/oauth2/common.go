// Package oauth2 Oauth2
package oauth2

import (
	"context"

	"golang.org/x/oauth2"
)

// ProviderType 第三方OAuth2提供商类型
type ProviderType string

const (
	// ProviderTypeGithub GitHub OAuth2提供商
	ProviderTypeGithub ProviderType = "github"
	// ProviderTypeGoogle Google OAuth2提供商
	ProviderTypeGoogle ProviderType = "google"
)

// UserInfo 用户信息
type UserInfo interface {
	GetID() string
	GetName() string
	GetEmail() string
	GetAvatar() string
}

// Provider OAuth2提供商接口
type Provider interface {
	// GetAuthURL 获取授权URL
	GetAuthURL() string
	// ExchangeToken 通过授权码获取Access Token
	ExchangeToken(ctx context.Context, code string) (*oauth2.Token, error)
	// GetUserInfo 获取用户信息
	GetUserInfo(ctx context.Context, token *oauth2.Token) (UserInfo, error)
	// GetBindField 获取绑定字段名
	GetBindField() string
}
