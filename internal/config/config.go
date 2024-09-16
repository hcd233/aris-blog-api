// Package config provides the configuration
package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

var (

	// ReadTimeout time is the read timeout key.
	//	@update 2024-06-22 08:59:40
	ReadTimeout time.Duration

	// WriteTimeout time is the write timeout key.
	//	@update 2024-06-22 08:59:37
	WriteTimeout time.Duration

	// MaxHeaderBytes int is the max header bytes key.
	//	@update 2024-06-22 08:59:34
	MaxHeaderBytes int

	// LogLevel string is the log level key.
	//	@update 2024-06-22 08:59:29
	LogLevel string

	// LogDirPath string is the log directory key.
	//	@update 2024-06-22 08:59:26
	LogDirPath string

	// Oauth2GithubClientID string is the Github client ID key.
	//	@update 2024-06-22 08:59:22
	Oauth2GithubClientID string

	// Oauth2GithubClientSecret string
	//	@update 2024-06-22 08:59:17
	Oauth2GithubClientSecret string

	// Oauth2StateString string is the OAuth2 state string key.
	//	@update 2024-06-22 08:59:11
	Oauth2StateString string

	// Oauth2GithubRedirectURL string
	//	@update 2024-06-22 08:59:07
	Oauth2GithubRedirectURL string

	// MysqlUser string
	//	@update 2024-06-22 09:00:30
	MysqlUser string

	// MysqlPassword string
	//	@update 2024-06-22 09:00:45
	MysqlPassword string

	// MysqlHost string
	//	@update 2024-06-22 09:01:02
	MysqlHost string

	// MysqlPort string
	//	@update 2024-06-22 09:01:18
	MysqlPort string

	// MysqlDatabase string
	//	@update 2024-06-22 09:01:34
	MysqlDatabase string

	// JwtTokenExpired int
	//	@update 2024-06-22 11:09:19
	JwtTokenExpired int

	// JwtTokenSecret string
	//	@update 2024-06-22 11:15:55
	JwtTokenSecret string
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

	ReadTimeout = time.Duration(config.GetInt("read.timeout")) * time.Second
	WriteTimeout = time.Duration(config.GetInt("write.timeout")) * time.Second
	MaxHeaderBytes = config.GetInt("max.header.bytes")
	LogLevel = config.GetString("log.level")
	LogDirPath = config.GetString("log.dir")

	Oauth2GithubClientID = config.GetString("oauth2.github.client.id")
	Oauth2GithubClientSecret = config.GetString("oauth2.github.client.secret")
	Oauth2StateString = config.GetString("oauth2.state.string")
	Oauth2GithubRedirectURL = config.GetString("oauth2.github.redirect.url")

	MysqlUser = config.GetString("mysql.user")
	MysqlPassword = config.GetString("mysql.password")
	MysqlHost = config.GetString("mysql.host")
	MysqlPort = config.GetString("mysql.port")
	MysqlDatabase = config.GetString("mysql.database")

	JwtTokenExpired = config.GetInt("jwt.token.expired")
	JwtTokenSecret = config.GetString("jwt.token.secret")

	if Oauth2GithubClientID == "" {
		panic("oauth2.github.client.id is required")
	}

	if Oauth2GithubClientSecret == "" {
		panic("oauth2.github.client.secret is required")
	}
}
