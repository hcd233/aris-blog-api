// Package config provides the configuration
package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

var (
	// Port is the port key.
	Port string
	// ReadTimeout is the read timeout key.
	ReadTimeout time.Duration
	// WriteTimeout is the write timeout key.
	WriteTimeout time.Duration
	// MaxHeaderBytes is the max header bytes key.
	MaxHeaderBytes int
	// LogLevel is the log level key.
	LogLevel string
	// LogDirPath is the log directory key.
	LogDirPath string
	// Oauth2GithubClientID is the Github client ID key.
	Oauth2GithubClientID string
	// Oauth2GithubClientSecret is the Github client secret key.
	Oauth2GithubClientSecret string
	// Oauth2StateString is the OAuth2 state string key.
	Oauth2StateString string
	// Oauth2GithubRedirectURL is the Github redirect URL key.
	Oauth2GithubRedirectURL string
)

func init() {
	initEnvironment()
}

func initEnvironment() {
	config := viper.New()
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	config.SetDefault("port", "8080")
	config.SetDefault("read.timeout", 10)
	config.SetDefault("write.timeout", 10)
	config.SetDefault("max.header.bytes", 1<<20)
	config.SetDefault("log.level", "info")
	config.SetDefault("log.dir", "./logs")

	config.AutomaticEnv()

	Port = config.GetString("port")
	ReadTimeout = time.Duration(config.GetInt("read.timeout")) * time.Second
	WriteTimeout = time.Duration(config.GetInt("write.timeout")) * time.Second
	MaxHeaderBytes = config.GetInt("max.header.bytes")
	LogLevel = config.GetString("log.level")
	LogDirPath = config.GetString("log.dir")

	Oauth2GithubClientID = config.GetString("oauth2.github.client.id")
	Oauth2GithubClientSecret = config.GetString("oauth2.github.client.secret")
	Oauth2StateString = config.GetString("oauth2.state.string")
	Oauth2GithubRedirectURL = config.GetString("oauth2.github.redirect.url")

	if Oauth2GithubClientID == "" {
		panic("oauth2.github.client.id is required")
	}

	if Oauth2GithubClientSecret == "" {
		panic("oauth2.github.client.secret is required")
	}
}
