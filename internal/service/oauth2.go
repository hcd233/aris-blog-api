package service

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/hcd233/Aris-blog/internal/auth"
	"github.com/hcd233/Aris-blog/internal/config"
	"github.com/hcd233/Aris-blog/internal/logger"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"

	"github.com/hcd233/Aris-blog/internal/util"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"gorm.io/gorm"
)

const (
	githubUserURL      = "https://api.github.com/user"
	githubUserEmailURL = "https://api.github.com/user/emails"
)

var githubUserScopes = []string{"user:email", "repo", "read:org"}

// Oauth2Service OAuth2服务
//
//	@author centonhuang
//	@update 2025-01-05 13:43:22
type Oauth2Service interface {
	Login(req *protocol.LoginRequest) (rsp *protocol.LoginResponse, err error)
	Callback(req *protocol.CallbackRequest) (rsp *protocol.CallbackResponse, err error)
}

type githubOauth2Service struct {
	oauth2Config       *oauth2.Config
	db                 *gorm.DB
	userDAO            *dao.UserDAO
	accessTokenSigner  auth.JwtTokenSigner
	refreshTokenSigner auth.JwtTokenSigner
}

// NewGithubOauth2Service 创建OAuth2服务
//
//	@return Oauth2Service
//	@author centonhuang
//	@update 2025-01-05 13:43:24
func NewGithubOauth2Service() Oauth2Service {
	return &githubOauth2Service{
		db:      database.GetDBInstance(),
		userDAO: dao.GetUserDAO(),
		oauth2Config: &oauth2.Config{
			Endpoint:     github.Endpoint,
			Scopes:       githubUserScopes,
			ClientID:     config.Oauth2GithubClientID,
			ClientSecret: config.Oauth2GithubClientSecret,
			RedirectURL:  config.Oauth2GithubRedirectURL,
		},
		accessTokenSigner:  auth.GetJwtAccessTokenSigner(),
		refreshTokenSigner: auth.GetJwtRefreshTokenSigner(),
	}
}

// Login 登录
//
//	@receiver s *oauth2Service
//	@return rsp *protocol.LoginResponse
//	@return err error
//	@author centonhuang
//	@update 2025-01-05 14:23:26
func (s *githubOauth2Service) Login(*protocol.LoginRequest) (rsp *protocol.LoginResponse, err error) {
	rsp = &protocol.LoginResponse{}

	url := s.oauth2Config.AuthCodeURL(config.Oauth2StateString, oauth2.AccessTypeOffline)
	rsp.RedirectURL = url

	logger.Logger.Info("[Oauth2Service] login", zap.String("redirectURL", url))

	return rsp, nil
}

// Callback 回调
//
//	@receiver s *oauth2Service
//	@param req *protocol.CallbackRequest
//	@return rsp *protocol.CallbackResponse
//	@return err error
//	@author centonhuang
//	@update 2025-01-05 14:23:26
func (s *githubOauth2Service) Callback(req *protocol.CallbackRequest) (rsp *protocol.CallbackResponse, err error) {
	rsp = &protocol.CallbackResponse{}

	if req.State != config.Oauth2StateString {
		logger.Logger.Error("[Oauth2Service] invalid state",
			zap.String("state", req.State),
			zap.String("expectedState", config.Oauth2StateString))
		return nil, protocol.ErrUnauthorized
	}

	token, err := s.oauth2Config.Exchange(context.Background(), req.Code)
	if err != nil {
		logger.Logger.Error("[Oauth2Service] failed to exchange token", zap.Error(err))
		return nil, protocol.ErrUnauthorized
	}

	data, err := s.getGithubUserInfo(token)
	if err != nil {
		logger.Logger.Error("[Oauth2Service] failed to get github user info", zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	githubID := strconv.FormatFloat(data["id"].(float64), 'f', -1, 64)
	userName, email, avatar := data["login"].(string), data["email"].(string), data["avatar_url"].(string)

	user, err := s.userDAO.GetByEmail(s.db, email, []string{"id", "name", "avatar"}, []string{})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Logger.Error("[Oauth2Service] failed to get user by email",
			zap.String("email", email),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if user.ID != 0 {
		// 更新已存在用户
		if err := s.userDAO.Update(s.db, user, map[string]interface{}{
			"last_login": time.Now(),
		}); err != nil {
			logger.Logger.Error("[Oauth2Service] failed to update user login time",
				zap.Uint("userID", user.ID),
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

		if err := s.userDAO.Create(s.db, user); err != nil {
			logger.Logger.Error("[Oauth2Service] failed to create user",
				zap.String("userName", userName),
				zap.Error(err))
			return nil, protocol.ErrInternalError
		}
	}

	if user.GithubBindID == "" {
		if err := s.userDAO.Update(s.db, user, map[string]interface{}{
			"github_bind_id": githubID,
		}); err != nil {
			logger.Logger.Error("[Oauth2Service] failed to update github bind id",
				zap.Uint("userID", user.ID),
				zap.String("githubID", githubID),
				zap.Error(err))
			return nil, protocol.ErrInternalError
		}
	}

	accessToken, err := s.accessTokenSigner.EncodeToken(user.ID)
	if err != nil {
		logger.Logger.Error("[Oauth2Service] failed to encode access token",
			zap.Uint("userID", user.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	refreshToken, err := s.refreshTokenSigner.EncodeToken(user.ID)
	if err != nil {
		logger.Logger.Error("[Oauth2Service] failed to encode refresh token",
			zap.Uint("userID", user.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.AccessToken = accessToken
	rsp.RefreshToken = refreshToken

	return rsp, nil
}

// getGithubUserInfo 获取Github用户信息
//
//	@receiver s *oauth2Service
//	@param token *oauth2.Token
//	@return map[string]interface{}
//	@return error
//	@author centonhuang
//	@update 2025-01-05 14:23:26
func (s *githubOauth2Service) getGithubUserInfo(token *oauth2.Token) (map[string]interface{}, error) {
	client := s.oauth2Config.Client(context.Background(), token)

	// 获取用户基本信息
	resp, err := client.Get(githubUserURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	// 获取用户邮箱信息
	emailResp, err := client.Get(githubUserEmailURL)
	if err != nil {
		return nil, err
	}
	defer emailResp.Body.Close()

	var emails []map[string]interface{}
	if err := json.NewDecoder(emailResp.Body).Decode(&emails); err != nil {
		return nil, err
	}

	// 选择主邮箱
	for _, email := range emails {
		if email["primary"].(bool) {
			data["email"] = email["email"]
			break
		}
	}

	return data, nil
}
