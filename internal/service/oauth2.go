package service

import (
	"context"
	"errors"
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
	"gorm.io/gorm"
)

const (
	githubUserURL      = "https://api.github.com/user"
	githubUserEmailURL = "https://api.github.com/user/emails"
)

var githubUserScopes = []string{"user:email", "repo", "read:org"}

// GithubUserInfo Github用户信息结构体
type GithubUserInfo struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// GithubEmail Github邮箱信息结构体
type GithubEmail struct {
	Email   string `json:"email"`
	Primary bool   `json:"primary"`
}

// Oauth2Service OAuth2服务
//
//	author centonhuang
//	update 2025-01-05 13:43:22
type Oauth2Service interface {
	Login(ctx context.Context, req *protocol.LoginRequest) (rsp *protocol.LoginResponse, err error)
	Callback(ctx context.Context, req *protocol.CallbackRequest) (rsp *protocol.CallbackResponse, err error)
}

type githubOauth2Service struct {
	oauth2Config       *oauth2.Config
	userDAO            *dao.UserDAO
	imageObjDAO        objdao.ObjDAO
	thumbnailObjDAO    objdao.ObjDAO
	accessTokenSigner  auth.JwtTokenSigner
	refreshTokenSigner auth.JwtTokenSigner
}

// NewGithubOauth2Service 创建OAuth2服务
//
//	return Oauth2Service
//	author centonhuang
//	update 2025-01-05 13:43:24
func NewGithubOauth2Service() Oauth2Service {
	return &githubOauth2Service{
		userDAO:         dao.GetUserDAO(),
		imageObjDAO:     objdao.GetImageObjDAO(),
		thumbnailObjDAO: objdao.GetThumbnailObjDAO(),
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
//	receiver s *oauth2Service
//	return rsp *protocol.LoginResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 14:23:26
func (s *githubOauth2Service) Login(ctx context.Context, req *protocol.LoginRequest) (rsp *protocol.LoginResponse, err error) {
	rsp = &protocol.LoginResponse{}

	logger := logger.LoggerWithContext(ctx)

	url := s.oauth2Config.AuthCodeURL(config.Oauth2StateString, oauth2.AccessTypeOffline)
	rsp.RedirectURL = url

	logger.Info("[Oauth2Service] login", zap.String("redirectURL", url))

	return rsp, nil
}

// Callback 回调
//
//	receiver s *oauth2Service
//	param req *protocol.CallbackRequest
//	return rsp *protocol.CallbackResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 14:23:26
func (s *githubOauth2Service) Callback(ctx context.Context, req *protocol.CallbackRequest) (rsp *protocol.CallbackResponse, err error) {
	rsp = &protocol.CallbackResponse{}

	logger := logger.LoggerWithContext(ctx)
	db := database.GetDBInstance(ctx)

	if req.State != config.Oauth2StateString {
		logger.Error("[Oauth2Service] invalid state",
			zap.String("state", req.State),
			zap.String("expectedState", config.Oauth2StateString))
		return nil, protocol.ErrUnauthorized
	}

	token, err := s.oauth2Config.Exchange(context.Background(), req.Code)
	if err != nil {
		logger.Error("[Oauth2Service] failed to exchange token", zap.Error(err))
		return nil, protocol.ErrUnauthorized
	}

	userInfo, err := s.getGithubUserInfo(token)
	if err != nil {
		logger.Error("[Oauth2Service] failed to get github user info", zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	githubID := strconv.FormatInt(userInfo.ID, 10)
	userName, email, avatar := userInfo.Login, userInfo.Email, userInfo.AvatarURL

	user, err := s.userDAO.GetByEmail(db, email, []string{"id", "name", "avatar"}, []string{})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("[Oauth2Service] failed to get user by email",
			zap.String("email", email),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if user.ID != 0 {
		// 更新已存在用户
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

		_, err = s.imageObjDAO.CreateDir(user.ID)
		if err != nil {
			logger.Error("[Oauth2Service] failed to create image dir",
				zap.Error(err))
			return nil, protocol.ErrInternalError
		}
		logger.Info("[Oauth2Service] image dir created")
		_, err = s.thumbnailObjDAO.CreateDir(user.ID)
		if err != nil {
			logger.Error("[Oauth2Service] failed to create thumbnail dir",
				zap.Error(err))
			return nil, protocol.ErrInternalError
		}
		logger.Info("[Oauth2Service] thumbnail dir created")
	}

	if user.GithubBindID == "" {
		if err := s.userDAO.Update(db, user, map[string]interface{}{
			"github_bind_id": githubID,
		}); err != nil {
			logger.Error("[Oauth2Service] failed to update github bind id",
				zap.String("githubID", githubID),
				zap.Error(err))
			return nil, protocol.ErrInternalError
		}
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

// getGithubUserInfo 获取Github用户信息
//
//	receiver s *oauth2Service
//	param token *oauth2.Token
//	return *GithubUserInfo
//	return error
//	author centonhuang
//	update 2025-01-05 14:23:26
func (s *githubOauth2Service) getGithubUserInfo(token *oauth2.Token) (*GithubUserInfo, error) {
	client := s.oauth2Config.Client(context.Background(), token)

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
