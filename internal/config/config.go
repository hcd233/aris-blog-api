// Package config provides the configuration
package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

var (

	// ReadTimeout time Gin读取超时时间
	//	update 2024-06-22 08:59:40
	ReadTimeout time.Duration

	// WriteTimeout time Gin写入超时时间
	//	update 2024-06-22 08:59:37
	WriteTimeout time.Duration

	// MaxHeaderBytes int Gin最大头部字节数
	//	update 2024-06-22 08:59:34
	MaxHeaderBytes int

	// LogLevel string 日志级别
	//	update 2024-06-22 08:59:29
	LogLevel string

	// LogDirPath string 日志目录路径
	//	update 2024-06-22 08:59:26
	LogDirPath string

	// Oauth2GithubClientID string Github OAuth2 Client ID
	//	update 2024-06-22 08:59:22
	Oauth2GithubClientID string

	// Oauth2GithubClientSecret string Github OAuth2 Client Secret
	//	update 2024-06-22 08:59:17
	Oauth2GithubClientSecret string

	// Oauth2StateString string Github OAuth2 State String
	//	update 2024-06-22 08:59:11
	Oauth2StateString string

	// Oauth2GithubRedirectURL string Github OAuth2 Redirect URL
	//	update 2024-06-22 08:59:07
	Oauth2GithubRedirectURL string

	// PostgresUser string Postgres用户名
	//	update 2024-06-22 09:00:30
	PostgresUser string

	// PostgresPassword string Postgres密码
	//	update 2024-06-22 09:00:45
	PostgresPassword string

	// PostgresHost string Postgres主机
	//	update 2024-06-22 09:01:02
	PostgresHost string

	// PostgresPort string Postgres端口
	//	update 2024-06-22 09:01:18
	PostgresPort string

	// PostgresDatabase string Postgres数据库
	//	update 2024-06-22 09:01:34
	PostgresDatabase string

	// PostgresSSLMode string Postgres SSL模式
	//	update 2024-06-22 09:01:50
	PostgresSSLMode string

	// RedisHost string Redis主机
	RedisHost string

	// RedisPort string Redis端口
	RedisPort string

	// RedisPassword string Redis密码
	RedisPassword string

	// MinioEndpoint string Minio Endpoint
	MinioEndpoint string

	// MinioTLS bool Minio TLS
	MinioTLS bool

	// MinioRegion string Minio Region
	MinioRegion string

	// MinioBucketName string Minio Bucket Name
	MinioBucketName string

	// MinioAccessID string Minio Access ID
	MinioAccessID string

	// MinioAccessKey string Minio Access Key
	MinioAccessKey string

	// CosRegion string Cos Region
	CosRegion string

	// CosSecretID string Cos Access ID
	CosSecretID string

	// CosSecretKey string Cos Secret Key
	CosSecretKey string

	// CosBucketName string Cos Bucket Name
	CosBucketName string

	// CosAppID string Cos App ID
	CosAppID string

	// OpenAIModel string OpenAI Model
	OpenAIModel string

	// OpenAIAPIKey string OpenAI API Key
	OpenAIAPIKey string

	// OpenAIBaseURL string OpenAI Base URL
	OpenAIBaseURL string

	// LangfuseHost string Langfuse Host
	LangfuseHost string

	// LangfusePublicKey string Langfuse Public Key
	LangfusePublicKey string

	// LangfuseSecretKey string Langfuse Secret Key
	LangfuseSecretKey string

	// JwtAccessTokenExpired time.Duration Access Jwt Token过期时间
	//	update 2024-06-22 11:09:19
	JwtAccessTokenExpired time.Duration

	// JwtAccessTokenSecret string Jwt Access Token密钥
	//	update 2024-06-22 11:15:55
	JwtAccessTokenSecret string

	// JwtRefreshTokenExpired time.Duration Refresh Jwt Token过期时间
	//	update 2024-06-22 11:09:19
	JwtRefreshTokenExpired time.Duration

	// JwtRefreshTokenSecret string Jwt Refresh Token密钥
	//	update 2024-06-22 11:15:55
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

	config.SetDefault("postgres.sslmode", "disable")
	config.SetDefault("postgres.timezone", "UTC")

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

	PostgresUser = config.GetString("postgres.user")
	PostgresPassword = config.GetString("postgres.password")
	PostgresHost = config.GetString("postgres.host")
	PostgresPort = config.GetString("postgres.port")
	PostgresDatabase = config.GetString("postgres.database")
	PostgresSSLMode = config.GetString("postgres.sslmode")

	RedisHost = config.GetString("redis.host")
	RedisPort = config.GetString("redis.port")
	RedisPassword = config.GetString("redis.password")

	MinioEndpoint = config.GetString("minio.endpoint")
	MinioTLS = config.GetBool("minio.tls")
	MinioRegion = config.GetString("minio.region")
	MinioBucketName = config.GetString("minio.bucket.name")
	MinioAccessID = config.GetString("minio.access.id")
	MinioAccessKey = config.GetString("minio.access.key")

	CosBucketName = config.GetString("cos.bucket.name")
	CosAppID = config.GetString("cos.app.id")
	CosRegion = config.GetString("cos.region")
	CosSecretID = config.GetString("cos.secret.id")
	CosSecretKey = config.GetString("cos.secret.key")

	OpenAIModel = config.GetString("openai.model")
	OpenAIAPIKey = config.GetString("openai.api.key")
	OpenAIBaseURL = config.GetString("openai.base.url")

	LangfuseHost = config.GetString("langfuse.host")
	LangfusePublicKey = config.GetString("langfuse.public.key")
	LangfuseSecretKey = config.GetString("langfuse.secret.key")

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
