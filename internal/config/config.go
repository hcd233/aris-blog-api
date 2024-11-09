// Package config provides the configuration
package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

var (

	// ReadTimeout time Gin读取超时时间
	//	@update 2024-06-22 08:59:40
	ReadTimeout time.Duration

	// WriteTimeout time Gin写入超时时间
	//	@update 2024-06-22 08:59:37
	WriteTimeout time.Duration

	// MaxHeaderBytes int Gin最大头部字节数
	//	@update 2024-06-22 08:59:34
	MaxHeaderBytes int

	// LogLevel string 日志级别
	//	@update 2024-06-22 08:59:29
	LogLevel string

	// LogDirPath string 日志目录路径
	//	@update 2024-06-22 08:59:26
	LogDirPath string

	// Oauth2GithubClientID string Github OAuth2 Client ID
	//	@update 2024-06-22 08:59:22
	Oauth2GithubClientID string

	// Oauth2GithubClientSecret string Github OAuth2 Client Secret
	//	@update 2024-06-22 08:59:17
	Oauth2GithubClientSecret string

	// Oauth2StateString string Github OAuth2 State String
	//	@update 2024-06-22 08:59:11
	Oauth2StateString string

	// Oauth2GithubRedirectURL string Github OAuth2 Redirect URL
	//	@update 2024-06-22 08:59:07
	Oauth2GithubRedirectURL string

	// MysqlUser string Mysql用户名
	//	@update 2024-06-22 09:00:30
	MysqlUser string

	// MysqlPassword string Mysql密码
	//	@update 2024-06-22 09:00:45
	MysqlPassword string

	// MysqlHost string Mysql主机
	//	@update 2024-06-22 09:01:02
	MysqlHost string

	// MysqlPort string Mysql端口
	//	@update 2024-06-22 09:01:18
	MysqlPort string

	// MysqlDatabase string Mysql数据库
	//	@update 2024-06-22 09:01:34
	MysqlDatabase string

	// MeilisearchHost string Meilisearch主机
	//	@update 2024-09-18 12:09:25
	MeilisearchHost string

	// MeilisearchPort string Meilisearch端口
	//	@update 2024-09-18 12:13:17
	MeilisearchPort string

	// MeilisearchMasterKey string Meilisearch API Key
	//	@update 2024-09-18 12:13:29
	MeilisearchMasterKey string

	// JwtAccessTokenExpired time.Duration Access Jwt Token过期时间
	//	@update 2024-06-22 11:09:19
	JwtAccessTokenExpired time.Duration

	// JwtAccessTokenSecret string Jwt Access Token密钥
	//	@update 2024-06-22 11:15:55
	JwtAccessTokenSecret string

	// JwtRefreshTokenExpired time.Duration Refresh Jwt Token过期时间
	//	@update 2024-06-22 11:09:19
	JwtRefreshTokenExpired time.Duration

	// JwtRefreshTokenSecret string Jwt Refresh Token密钥
	//	@update 2024-06-22 11:15:55
	JwtRefreshTokenSecret string
)

func init() {
	initEnvironment()
}

func initEnvironment() {
	config := viper.New()
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

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

	MeilisearchHost = config.GetString("meilisearch.host")
	MeilisearchPort = config.GetString("meilisearch.port")
	MeilisearchMasterKey = config.GetString("meilisearch.master.key")

	JwtAccessTokenExpired = config.GetDuration("jwt.access.token.expired")
	JwtAccessTokenSecret = config.GetString("jwt.access.token.secret")

	JwtRefreshTokenExpired = config.GetDuration("jwt.refresh.token.expired")
	JwtRefreshTokenSecret = config.GetString("jwt.refresh.token.secret")

	if Oauth2GithubClientID == "" {
		panic("oauth2.github.client.id is required")
	}

	if Oauth2GithubClientSecret == "" {
		panic("oauth2.github.client.secret is required")
	}
}
