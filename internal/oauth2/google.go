package oauth2

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/hcd233/aris-blog-api/internal/config"
	"github.com/hcd233/aris-blog-api/internal/logger"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleUserScopes = []string{
	"openid",
	"profile",
	"email",
	"https://www.googleapis.com/auth/userinfo.profile",
	"https://www.googleapis.com/auth/userinfo.email",
}

// GoogleUserInfo Google用户信息结构体
type GoogleUserInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	PhotoURL string `json:"picture"`
}

// GetID 获取Google用户ID
//
//	@receiver u *GoogleUserInfo
//	@return string
//	@author centonhuang
//	@update 2025-10-31 14:48:46
func (u *GoogleUserInfo) GetID() string {
	return u.ID
}

// GetName 获取Google用户名
//
//	@receiver u *GoogleUserInfo
//	@return string
//	@author centonhuang
//	@update 2025-10-31 14:48:48
func (u *GoogleUserInfo) GetName() string {
	return u.Name
}

// GetEmail 获取Google用户邮箱
//
//	@receiver u *GoogleUserInfo
//	@return string
//	@author centonhuang
//	@update 2025-10-31 14:48:50
func (u *GoogleUserInfo) GetEmail() string {
	return u.Email
}

// GetAvatar 获取Google用户头像
//
//	@receiver u *GoogleUserInfo
//	@return string
//	@author centonhuang
//	@update 2025-10-31 14:48:52
func (u *GoogleUserInfo) GetAvatar() string {
	return u.PhotoURL
}

// googleProvider Google OAuth2提供商实现
type googleProvider struct {
	oauth2Config *oauth2.Config
}

// NewGoogleProvider Google提供商
//
//	@return Provider
//	@author centonhuang
//	@update 2025-10-31 14:57:11
func NewGoogleProvider() Provider {
	return &googleProvider{
		oauth2Config: &oauth2.Config{
			Endpoint:     google.Endpoint,
			Scopes:       googleUserScopes,
			ClientID:     config.Oauth2GoogleClientID,
			ClientSecret: config.Oauth2GoogleClientSecret,
			RedirectURL:  config.Oauth2GoogleRedirectURL,
		},
	}
}

func (p *googleProvider) GetAuthURL() string {
	return p.oauth2Config.AuthCodeURL(config.Oauth2StateString, oauth2.AccessTypeOffline)
}

func (p *googleProvider) ExchangeToken(ctx context.Context, code string) (*oauth2.Token, error) {
	logger := logger.WithCtx(ctx)

	logger.Info("[GoogleOauth2] exchanging code for token",
		zap.String("clientID", p.oauth2Config.ClientID),
		zap.String("redirectURL", p.oauth2Config.RedirectURL),
		zap.Strings("scopes", p.oauth2Config.Scopes))

	token, err := p.oauth2Config.Exchange(ctx, code)
	if err != nil {
		logger.Error("[GoogleOauth2] token exchange failed", zap.Error(err))
		return nil, err
	}

	logger.Info("[GoogleOauth2] token exchange successful")
	return token, nil
}

func (p *googleProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (UserInfo, error) {
	logger := logger.WithCtx(ctx)

	// 使用HTTP客户端直接调用Google OAuth2 UserInfo API
	client := p.oauth2Config.Client(ctx, token)

	logger.Info("[GoogleOauth2] calling Google UserInfo API")

	// 调用Google UserInfo API
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		logger.Error("[GoogleOauth2] failed to call userinfo API", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	logger.Info("[GoogleOauth2] userinfo API response",
		zap.Int("statusCode", resp.StatusCode))

	var userInfoResp struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Email   string `json:"email"`
		Picture string `json:"picture"`
	}

	if err := sonic.ConfigDefault.NewDecoder(resp.Body).Decode(&userInfoResp); err != nil {
		logger.Error("[GoogleOauth2] failed to decode userinfo response", zap.Error(err))
		return nil, err
	}

	logger.Info("[GoogleOauth2] successfully decoded user info",
		zap.String("userID", userInfoResp.ID),
		zap.String("userName", userInfoResp.Name),
		zap.String("userEmail", userInfoResp.Email))

	userInfo := &GoogleUserInfo{
		ID:       userInfoResp.ID,
		Name:     userInfoResp.Name,
		Email:    userInfoResp.Email,
		PhotoURL: userInfoResp.Picture,
	}

	return userInfo, nil
}

func (p *googleProvider) GetBindField() string {
	return "google_bind_id"
}
