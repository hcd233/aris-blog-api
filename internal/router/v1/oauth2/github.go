// Package oauth2 login and callback handlers.
package oauth2

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/auth"
	"github.com/hcd233/Aris-blog/internal/config"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"github.com/hcd233/Aris-blog/internal/resource/search"
	"github.com/samber/lo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

const (
	githubUserURL      = "https://api.github.com/user"
	githubUserEmailURL = "https://api.github.com/user/emails"
)

var githubOauthConfig = &oauth2.Config{
	ClientID:     config.Oauth2GithubClientID,
	ClientSecret: config.Oauth2GithubClientSecret,
	Endpoint:     github.Endpoint,
	RedirectURL:  config.Oauth2GithubRedirectURL,
	Scopes:       []string{"user:email", "repo", "read:org"},
}

// GithubLoginHandler Github登录
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-09-16 01:56:09
func GithubLoginHandler(c *gin.Context) {
	url := githubOauthConfig.AuthCodeURL(config.Oauth2StateString, oauth2.AccessTypeOffline)
	c.JSON(200, protocol.Response{
		Code:    protocol.CodeOk,
		Message: "Redirect to Github login page",
		Data:    map[string]interface{}{"redirect_url": url},
	})
}

// GithubCallbackHandler Github登录回调
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-09-16 01:56:03
func GithubCallbackHandler(c *gin.Context) {
	params := protocol.GithubCallbackParams{}
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

	token, err := githubOauthConfig.Exchange(c, params.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code: protocol.CodeTokenError,
		})
		return
	}

	data, err := getGithubUserInfo(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	githubID := strconv.FormatFloat(data["id"].(float64), 'f', -1, 64)
	userName, email, avatar := data["login"].(string), data["email"].(string), data["avatar_url"].(string)
	user, err := model.QueryUserByEmail(email, []string{"id", "name", "avatar"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeQueryUserError,
			Message: err.Error(),
		})
		return
	}

	if user != nil {
		// 如果已有用户，刷新信息
		// atomic
		user = lo.Must1(model.UpdateUserInfoByID(user.ID, map[string]interface{}{
			"last_login": time.Now(),
		}))
		lo.Must0(search.UpdateUserIndex(user.GetUserBasicInfo()))
	} else {
		// 新用户，保存信息
		permission := model.PermissionGeneral
		// atomic
		user = lo.Must(model.CreateUserByBasicInfo(userName, email, avatar, permission))
		lo.Must0(search.AddUserIndex(user.GetUserBasicInfo()))
	}

	if user.GithubBindID == "" {
		user.BindGithubID(githubID)
	}
	tokenString := lo.Must(auth.EncodeToken(user.ID))
	c.JSON(http.StatusOK, protocol.Response{
		Code:    protocol.CodeOk,
		Message: "Login success",
		Data: map[string]interface{}{
			"token": tokenString,
		},
	})
}

func getGithubUserInfo(token *oauth2.Token) (map[string]interface{}, error) {
	client := githubOauthConfig.Client(context.Background(), token)

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
