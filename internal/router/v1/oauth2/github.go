// Package oauth2 login and callback handlers.
package oauth2

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-AI-go/internal/config"
	"github.com/hcd233/Aris-AI-go/internal/protocol"
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
	c.Redirect(302, url)
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

	data := map[string]interface{}{
		"token": token,
	}

	c.JSON(200, protocol.Response{
		Data:   data,
		Status: protocol.SUCCESS,
	})
}
