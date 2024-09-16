// Package oauth2 login and callback handlers.
package oauth2

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-AI-go/internal/auth"
	"github.com/hcd233/Aris-AI-go/internal/config"
	"github.com/hcd233/Aris-AI-go/internal/protocol"
	"github.com/hcd233/Aris-AI-go/internal/resource/database/model"
	"github.com/samber/lo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

const githubUserURL = "https://api.github.com/user"

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
	state := c.Query("state")
	if state != config.Oauth2StateString {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code: protocol.CodeStateError,
		})
		return
	}

	code := c.Query("code")
	token, err := githubOauthConfig.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code: protocol.CodeTokenError,
		})
		return
	}

	data, err := getGithubUserInfo(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code: protocol.CodeGetUserError,
		})
		return
	}

	platform := model.PlatformGithub
	bindID := strconv.FormatFloat(data["id"].(float64), 'f', -1, 64)
	user := lo.Must1(model.QueryUserByPlatformAndID(platform, bindID))

	if user != nil {
		// 如果已有用户，刷新信息
		lo.Must0(user.UpdateUserInfo())
	} else {
		// 新用户，保存信息
		userName, avatar := data["login"].(string), data["avatar_url"].(string)
		permission := model.PermissionGeneral
		user = lo.Must(model.AddUserByBasicInfo(userName, avatar, permission, platform, bindID))
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
	resp, err := client.Get(githubUserURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}
