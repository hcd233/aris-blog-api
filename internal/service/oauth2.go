package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/hcd233/aris-blog-api/internal/auth"
	"github.com/hcd233/aris-blog-api/internal/config"
	"github.com/hcd233/aris-blog-api/internal/logger"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/resource/database"
	"github.com/hcd233/aris-blog-api/internal/resource/database/dao"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
	objdao "github.com/hcd233/aris-blog-api/internal/resource/storage/obj_dao"

	"github.com/hcd233/aris-blog-api/internal/util"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

// OAuth2Provider 第三方OAuth2提供商类型
type OAuth2Provider string

const (
	// ProviderGithub GitHub OAuth2提供商
	ProviderGithub OAuth2Provider = "github"
	// ProviderQQ QQ OAuth2提供商
	ProviderQQ OAuth2Provider = "qq"
	// ProviderGoogle Google OAuth2提供商
	ProviderGoogle OAuth2Provider = "google"
)

// 第三方服务常量定义
const (
	// GitHub相关
	githubUserURL      = "https://api.github.com/user"
	githubUserEmailURL = "https://api.github.com/user/emails"

	// QQ相关
	qqAuthorizeURL = "https://graph.qq.com/oauth2.0/authorize"
	qqTokenURL     = "https://graph.qq.com/oauth2.0/token"
	qqUserInfoURL  = "https://graph.qq.com/user/get_user_info"
	qqOpenIDURL    = "https://graph.qq.com/oauth2.0/me"

	// Google相关
	googleUserInfoURL = "https://people.googleapis.com/v1/people/me"
)

var (
	githubUserScopes = []string{"user:email", "repo", "read:org"}
	qqUserScopes     = []string{"get_user_info"}
	googleUserScopes = []string{
		"openid",
		"profile",
		"email",
		"https://www.googleapis.com/auth/userinfo.profile",
		"https://www.googleapis.com/auth/userinfo.email",
	}
)

// OAuth2UserInfo 第三方用户信息接口
type OAuth2UserInfo interface {
	GetID() string
	GetName() string
	GetEmail() string
	GetAvatar() string
}

// GithubUserInfo Github用户信息结构体
type GithubUserInfo struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

func (u *GithubUserInfo) GetID() string {
	return strconv.FormatInt(u.ID, 10)
}

func (u *GithubUserInfo) GetName() string {
	return u.Login
}

func (u *GithubUserInfo) GetEmail() string {
	return u.Email
}

func (u *GithubUserInfo) GetAvatar() string {
	return u.AvatarURL
}

// GithubEmail Github邮箱信息结构体
type GithubEmail struct {
	Email   string `json:"email"`
	Primary bool   `json:"primary"`
}

// QQUserInfo QQ用户信息结构体
type QQUserInfo struct {
	OpenID   string `json:"openid"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"figureurl_qq_1"`
}

func (u *QQUserInfo) GetID() string {
	return u.OpenID
}

func (u *QQUserInfo) GetName() string {
	return u.Nickname
}

func (u *QQUserInfo) GetEmail() string {
	// QQ OAuth2默认不提供邮箱，使用openid@qq.oauth.placeholder格式
	return fmt.Sprintf("%s@qq.oauth.placeholder", u.OpenID)
}

func (u *QQUserInfo) GetAvatar() string {
	return u.Avatar
}

// QQOpenIDResponse QQ OpenID响应结构体
type QQOpenIDResponse struct {
	ClientID string `json:"client_id"`
	OpenID   string `json:"openid"`
}

// GoogleUserInfo Google用户信息结构体
type GoogleUserInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	PhotoURL string `json:"picture"`
}

func (u *GoogleUserInfo) GetID() string {
	return u.ID
}

func (u *GoogleUserInfo) GetName() string {
	return u.Name
}

func (u *GoogleUserInfo) GetEmail() string {
	return u.Email
}

func (u *GoogleUserInfo) GetAvatar() string {
	return u.PhotoURL
}

// OAuth2Provider 第三方OAuth2提供商接口
type OAuth2ProviderInterface interface {
	// GetAuthURL 获取授权URL
	GetAuthURL() string
	// ExchangeToken 通过授权码获取Access Token
	ExchangeToken(ctx context.Context, code string) (*oauth2.Token, error)
	// GetUserInfo 获取用户信息
	GetUserInfo(ctx context.Context, token *oauth2.Token) (OAuth2UserInfo, error)
	// GetBindField 获取绑定字段名
	GetBindField() string
}

// Oauth2Service OAuth2服务接口
type Oauth2Service interface {
	Login(ctx context.Context, req *protocol.LoginRequest) (rsp *protocol.LoginResponse, err error)
	Callback(ctx context.Context, req *protocol.CallbackRequest) (rsp *protocol.CallbackResponse, err error)
}

// oauth2Service OAuth2服务基础实现
type oauth2Service struct {
	provider           OAuth2ProviderInterface
	userDAO            *dao.UserDAO
	imageObjDAO        objdao.ObjDAO
	thumbnailObjDAO    objdao.ObjDAO
	accessTokenSigner  auth.JwtTokenSigner
	refreshTokenSigner auth.JwtTokenSigner
}

// githubProvider GitHub OAuth2提供商实现
type githubProvider struct {
	oauth2Config *oauth2.Config
}

func newGithubProvider() OAuth2ProviderInterface {
	return &githubProvider{
		oauth2Config: &oauth2.Config{
			Endpoint:     github.Endpoint,
			Scopes:       githubUserScopes,
			ClientID:     config.Oauth2GithubClientID,
			ClientSecret: config.Oauth2GithubClientSecret,
			RedirectURL:  config.Oauth2GithubRedirectURL,
		},
	}
}

func (p *githubProvider) GetAuthURL() string {
	return p.oauth2Config.AuthCodeURL(config.Oauth2StateString, oauth2.AccessTypeOffline)
}

func (p *githubProvider) ExchangeToken(ctx context.Context, code string) (*oauth2.Token, error) {
	return p.oauth2Config.Exchange(ctx, code)
}

func (p *githubProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (OAuth2UserInfo, error) {
	client := p.oauth2Config.Client(ctx, token)

	// 获取用户基本信息
	resp, err := client.Get(githubUserURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo GithubUserInfo
	if err := sonic.ConfigDefault.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	// 获取用户邮箱信息
	emailResp, err := client.Get(githubUserEmailURL)
	if err != nil {
		return nil, err
	}
	defer emailResp.Body.Close()

	var emails []GithubEmail
	if err := sonic.ConfigDefault.NewDecoder(emailResp.Body).Decode(&emails); err != nil {
		return nil, err
	}

	// 选择主邮箱
	for _, email := range emails {
		if email.Primary {
			userInfo.Email = email.Email
			break
		}
	}

	return &userInfo, nil
}

func (p *githubProvider) GetBindField() string {
	return "github_bind_id"
}

// qqProvider QQ OAuth2提供商实现
type qqProvider struct {
	oauth2Config *oauth2.Config
}

func newQQProvider() OAuth2ProviderInterface {
	return &qqProvider{
		oauth2Config: &oauth2.Config{
			Endpoint: oauth2.Endpoint{
				AuthURL:  qqAuthorizeURL,
				TokenURL: qqTokenURL,
			},
			Scopes:       qqUserScopes,
			ClientID:     config.Oauth2QQClientID,
			ClientSecret: config.Oauth2QQClientSecret,
			RedirectURL:  config.Oauth2QQRedirectURL,
		},
	}
}

func (p *qqProvider) GetAuthURL() string {
	return p.oauth2Config.AuthCodeURL(config.Oauth2StateString, oauth2.AccessTypeOffline)
}

func (p *qqProvider) ExchangeToken(ctx context.Context, code string) (*oauth2.Token, error) {
	return p.oauth2Config.Exchange(ctx, code)
}

func (p *qqProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (OAuth2UserInfo, error) {
	client := p.oauth2Config.Client(ctx, token)

	// 获取OpenID
	openIDResp, err := client.Get(fmt.Sprintf("%s?access_token=%s&fmt=json", qqOpenIDURL, token.AccessToken))
	if err != nil {
		return nil, err
	}
	defer openIDResp.Body.Close()

	var openIDData QQOpenIDResponse
	if err := sonic.ConfigDefault.NewDecoder(openIDResp.Body).Decode(&openIDData); err != nil {
		return nil, err
	}

	// 获取用户信息
	userInfoURL := fmt.Sprintf("%s?access_token=%s&oauth_consumer_key=%s&openid=%s",
		qqUserInfoURL, token.AccessToken, p.oauth2Config.ClientID, openIDData.OpenID)

	userResp, err := client.Get(userInfoURL)
	if err != nil {
		return nil, err
	}
	defer userResp.Body.Close()

	var userInfo QQUserInfo
	if err := sonic.ConfigDefault.NewDecoder(userResp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	userInfo.OpenID = openIDData.OpenID
	return &userInfo, nil
}

func (p *qqProvider) GetBindField() string {
	return "qq_bind_id"
}

// googleProvider Google OAuth2提供商实现
type googleProvider struct {
	oauth2Config *oauth2.Config
}

func newGoogleProvider() OAuth2ProviderInterface {
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

func (p *googleProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (OAuth2UserInfo, error) {
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

// NewGithubOauth2Service 创建Github OAuth2服务
func NewGithubOauth2Service() Oauth2Service {
	return &oauth2Service{
		provider:           newGithubProvider(),
		userDAO:            dao.GetUserDAO(),
		imageObjDAO:        objdao.GetImageObjDAO(),
		thumbnailObjDAO:    objdao.GetThumbnailObjDAO(),
		accessTokenSigner:  auth.GetJwtAccessTokenSigner(),
		refreshTokenSigner: auth.GetJwtRefreshTokenSigner(),
	}
}

// NewQQOauth2Service 创建QQ OAuth2服务
func NewQQOauth2Service() Oauth2Service {
	return &oauth2Service{
		provider:           newQQProvider(),
		userDAO:            dao.GetUserDAO(),
		imageObjDAO:        objdao.GetImageObjDAO(),
		thumbnailObjDAO:    objdao.GetThumbnailObjDAO(),
		accessTokenSigner:  auth.GetJwtAccessTokenSigner(),
		refreshTokenSigner: auth.GetJwtRefreshTokenSigner(),
	}
}

// NewGoogleOauth2Service 创建Google OAuth2服务
func NewGoogleOauth2Service() Oauth2Service {
	return &oauth2Service{
		provider:           newGoogleProvider(),
		userDAO:            dao.GetUserDAO(),
		imageObjDAO:        objdao.GetImageObjDAO(),
		thumbnailObjDAO:    objdao.GetThumbnailObjDAO(),
		accessTokenSigner:  auth.GetJwtAccessTokenSigner(),
		refreshTokenSigner: auth.GetJwtRefreshTokenSigner(),
	}
}

// Login 登录
func (s *oauth2Service) Login(ctx context.Context, req *protocol.LoginRequest) (rsp *protocol.LoginResponse, err error) {
	rsp = &protocol.LoginResponse{}

	logger := logger.WithCtx(ctx)

	url := s.provider.GetAuthURL()
	rsp.RedirectURL = url

	logger.Info("[Oauth2Service] login", zap.String("redirectURL", url))

	return rsp, nil
}

// Callback 回调
func (s *oauth2Service) Callback(ctx context.Context, req *protocol.CallbackRequest) (rsp *protocol.CallbackResponse, err error) {
	rsp = &protocol.CallbackResponse{}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	if req.State != config.Oauth2StateString {
		logger.Error("[Oauth2Service] invalid state",
			zap.String("state", req.State),
			zap.String("expectedState", config.Oauth2StateString))
		return nil, protocol.ErrUnauthorized
	}

	logger.Info("[Oauth2Service] exchanging token",
		zap.String("code", req.Code),
		zap.String("state", req.State))

	token, err := s.provider.ExchangeToken(ctx, req.Code)
	if err != nil {
		logger.Error("[Oauth2Service] failed to exchange token",
			zap.String("code", req.Code),
			zap.Error(err))
		return nil, protocol.ErrUnauthorized
	}

	logger.Info("[Oauth2Service] token exchange successful",
		zap.String("tokenType", token.TokenType),
		zap.Bool("valid", token.Valid()))

	userInfo, err := s.provider.GetUserInfo(ctx, token)
	if err != nil {
		logger.Error("[Oauth2Service] failed to get user info", zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	thirdPartyID := userInfo.GetID()
	userName, email, avatar := userInfo.GetName(), userInfo.GetEmail(), userInfo.GetAvatar()

	user, err := s.userDAO.GetByEmail(db, email, []string{"id", "name", "avatar"}, []string{})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("[Oauth2Service] failed to get user by email",
			zap.String("email", email),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if user.ID != 0 {
		// 更新已存在用户的登录时间
		if err := s.userDAO.Update(db, user, map[string]interface{}{
			"last_login": time.Now(),
		}); err != nil {
			logger.Error("[Oauth2Service] failed to update user login time",
				zap.Error(err))
			return nil, protocol.ErrInternalError
		}
	} else {
		// 创建新用户
		if validateErr := util.ValidateUserName(userName); validateErr != nil {
			userName = "ArisUser" + strconv.FormatInt(time.Now().Unix(), 10)
		}
		defaultCategory := &model.Category{Name: userName}

		user = &model.User{
			Name:       userName,
			Email:      email,
			Avatar:     avatar,
			Permission: model.PermissionReader,
			LastLogin:  time.Now(),
			LLMQuota:   model.QuotaReader,
			Categories: []model.Category{*defaultCategory},
		}

		if err := s.userDAO.Create(db, user); err != nil {
			logger.Error("[Oauth2Service] failed to create user",
				zap.String("userName", userName),
				zap.Error(err))
			return nil, protocol.ErrInternalError
		}

		_, err = s.imageObjDAO.CreateDir(ctx, user.ID)
		if err != nil {
			logger.Error("[Oauth2Service] failed to create image dir",
				zap.Error(err))
			return nil, protocol.ErrInternalError
		}
		logger.Info("[Oauth2Service] image dir created")

		_, err = s.thumbnailObjDAO.CreateDir(ctx, user.ID)
		if err != nil {
			logger.Error("[Oauth2Service] failed to create thumbnail dir",
				zap.Error(err))
			return nil, protocol.ErrInternalError
		}
		logger.Info("[Oauth2Service] thumbnail dir created")
	}

	// 更新第三方平台绑定ID
	bindField := s.provider.GetBindField()
	updateData := map[string]interface{}{
		bindField: thirdPartyID,
	}

	if err := s.userDAO.Update(db, user, updateData); err != nil {
		logger.Error("[Oauth2Service] failed to update third party bind id",
			zap.String("bindField", bindField),
			zap.String("thirdPartyID", thirdPartyID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	accessToken, err := s.accessTokenSigner.EncodeToken(user.ID)
	if err != nil {
		logger.Error("[Oauth2Service] failed to encode access token",
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	refreshToken, err := s.refreshTokenSigner.EncodeToken(user.ID)
	if err != nil {
		logger.Error("[Oauth2Service] failed to encode refresh token",
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.AccessToken = accessToken
	rsp.RefreshToken = refreshToken

	return rsp, nil
}
