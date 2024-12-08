package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/auth"
	"github.com/hcd233/Aris-blog/internal/config"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	doc_dao "github.com/hcd233/Aris-blog/internal/resource/search/doc_dao"
	"github.com/hcd233/Aris-blog/internal/resource/search/document"
	"github.com/hcd233/Aris-blog/internal/util"
	"github.com/samber/lo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"gorm.io/gorm"
)

const (
	githubUserURL      = "https://api.github.com/user"
	githubUserEmailURL = "https://api.github.com/user/emails"
)

var githubUserScopes = []string{"user:email", "repo", "read:org"}

// Oauth2Service Oauth2服务
//
//	@author centonhuang
//	@update 2024-12-08 16:59:38
type Oauth2Service interface {
	LoginHandler(c *gin.Context)
	CallbackHandler(c *gin.Context)
}

type githubOauth2Service struct {
	oauth2Config       *oauth2.Config
	db                 *gorm.DB
	userDAO            *dao.UserDAO
	userDocDAO         *doc_dao.UserDocDAO
	accessTokenSigner  auth.JwtTokenSigner
	refreshTokenSigner auth.JwtTokenSigner
}

// NewGithubOauth2Service 创建Github Oauth2服务
//
//	@return Oauth2Service
//	@author centonhuang
//	@update 2024-12-08 16:59:38
func NewGithubOauth2Service() Oauth2Service {
	return &githubOauth2Service{
		db:         database.GetDBInstance(),
		userDAO:    dao.GetUserDAO(),
		userDocDAO: doc_dao.GetUserDocDAO(),
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

func (s *githubOauth2Service) LoginHandler(c *gin.Context) {
	url := s.oauth2Config.AuthCodeURL(config.Oauth2StateString, oauth2.AccessTypeOffline)
	c.JSON(200, protocol.Response{
		Code:    protocol.CodeOk,
		Message: "Redirect to Github login page",
		Data: map[string]interface{}{
			"redirectURL": url,
		},
	})
}

func (s *githubOauth2Service) CallbackHandler(c *gin.Context) {
	params := protocol.GithubCallbackParam{}

	if err := c.BindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeURIError,
			Message: err.Error(),
		})
		return
	}

	if params.State != config.Oauth2StateString {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code: protocol.CodeStateError,
		})
		return
	}

	token, err := s.oauth2Config.Exchange(c, params.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code: protocol.CodeTokenError,
		})
		return
	}

	data, err := s.getGithubUserInfo(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	githubID := strconv.FormatFloat(data["id"].(float64), 'f', -1, 64)
	userName, email, avatar := data["login"].(string), data["email"].(string), data["avatar_url"].(string)

	user, err := s.userDAO.GetByEmail(s.db, email, []string{"id", "name", "avatar"}, []string{})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeQueryUserError,
			Message: err.Error(),
		})
		return
	}

	if user.ID != 0 {
		lo.Must0(s.userDAO.Update(s.db, user, map[string]interface{}{
			"last_login": time.Now(),
		}))
		lo.Must0(s.userDocDAO.UpdateDocument(document.TransformUserToDocument(user)))
	} else {
		// 新用户，保存信息
		if validateErr := util.ValidateUserName(userName); validateErr != nil {
			userName = fmt.Sprintf("ArisUser" + strconv.FormatInt(time.Now().Unix(), 10))
		}
		defaultCategory := &model.Category{Name: userName}

		user = &model.User{
			Name:       userName,
			Email:      email,
			Avatar:     avatar,
			Permission: model.PermissionReader,
			LLMQuota:   model.QuotaReader,
			Categories: []model.Category{*defaultCategory},
		}

		// 插入用户信息
		lo.Must0(s.userDAO.Create(s.db, user))
		// 插入用户到搜索引擎
		lo.Must0(s.userDocDAO.AddDocument(document.TransformUserToDocument(user)))
	}

	if user.GithubBindID == "" {
		s.userDAO.Update(s.db, user, map[string]interface{}{
			"github_bind_id": githubID,
		})
	}

	accessToken := lo.Must(s.accessTokenSigner.EncodeToken(user.ID))
	refreshToken := lo.Must(s.refreshTokenSigner.EncodeToken(user.ID))

	c.JSON(http.StatusOK, protocol.Response{
		Code:    protocol.CodeOk,
		Message: "Login success",
		Data: map[string]interface{}{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		},
	})
}

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
