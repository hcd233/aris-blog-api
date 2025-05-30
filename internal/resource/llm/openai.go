// Package llm OpenAI客户端模块
//
//	update 2024-12-09 16:00:18
package llm

import (
	"github.com/hcd233/aris-blog-api/internal/config"
	"github.com/hcd233/aris-blog-api/internal/logger"
	openai "github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

var client *openai.Client

// InitOpenAIClient 初始化OpenAI客户端
//
//	author centonhuang
//	update 2024-12-09 16:00:26
func InitOpenAIClient() {
	clientConfig := openai.DefaultConfig(config.OpenAIAPIKey)
	clientConfig.BaseURL = config.OpenAIBaseURL
	client = openai.NewClientWithConfig(clientConfig)

	logger.Logger().Info("[OpenAI] Connected to OpenAI API", zap.String("baseURL", config.OpenAIBaseURL))
}

// GetOpenAIClient 获取OpenAI客户端
//
//	return *openai.Client
//	author centonhuang
//	update 2024-12-09 16:00:31
func GetOpenAIClient() *openai.Client {
	return client
}
