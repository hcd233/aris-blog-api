// Package oauth2 login and callback handlers.
package oauth2

import (
	"context"
	"encoding/json"
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

var githubOauthConfig = &oauth2.Config{
	ClientID:     config.Oauth2GithubClientID,
	ClientSecret: config.Oauth2GithubClientSecret,
	Endpoint:     github.Endpoint,
	RedirectURL:  config.Oauth2GithubRedirectURL,
	Scopes:       []string{"user:email", "repo", "read:org"},
}

// HandleGithubLogin handles the Github login.
func handleGithubLogin(c *gin.Context) {
	url := githubOauthConfig.AuthCodeURL(config.Oauth2StateString, oauth2.AccessTypeOffline)
	c.JSON(200, protocol.Response{
		Data:   map[string]interface{}{"url": url},
		Status: protocol.SUCCESS,
	})
}

// HandleGithubCallback handles the Github callback.
func handleGithubCallback(c *gin.Context) {
	state := c.Query("state")
	if state != config.Oauth2StateString {
		c.JSON(400, protocol.Response{
			Message: "Invalid state",
			Status:  protocol.FAILED,
		})
		return
	}

	code := c.Query("code")
	token, err := githubOauthConfig.Exchange(c, code)
	if err != nil {
		c.JSON(400, protocol.Response{
			Message: err.Error(),
			Status:  protocol.FAILED,
		})
		return
	}

	data, err := getGithubUserInfo(token)
	if err != nil {
		c.JSON(500, protocol.Response{
			Message: err.Error(),
			Status:  protocol.FAILED,
		})
		return
	}

	platform := model.PlatformGithub
	bindID := strconv.FormatFloat(data["id"].(float64), 'f', -1, 64)
	user, err := model.QueryUserByPlatformAndID(platform, bindID)
	if err != nil {
		c.JSON(400, protocol.Response{
			Message: err.Error(),
			Status:  protocol.ERROR,
		})
		return
	}
	if user != nil {
		// 如果已有用户，刷新最后登陆时间
		lo.Must0(user.SetLastLoginTime())
	} else {
		// 新用户，保存信息
		userName, avatar := data["login"].(string), data["avatar_url"].(string)
		permission := model.PermissionGeneral
		user = lo.Must(model.AddUserByBasicInfo(userName, avatar, permission, platform, bindID))
	}

	tokenString := lo.Must(auth.EncodeToken(user.ID))
	c.JSON(200, protocol.Response{
		Data: map[string]interface{}{
			"token": tokenString,
		},
	})
}

func getGithubUserInfo(token *oauth2.Token) (map[string]interface{}, error) {
	client := githubOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
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
