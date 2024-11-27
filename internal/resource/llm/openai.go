package llm

import (
	"github.com/hcd233/Aris-blog/internal/config"
	openai "github.com/sashabaranov/go-openai"
)

var (
	client *openai.Client
)

func InitOpenAIClient() {
	clientConfig := openai.DefaultConfig(config.OpenAIAPIKey)
	clientConfig.BaseURL = config.OpenAIBaseURL
	client = openai.NewClientWithConfig(clientConfig)
}

func GetOpenAIClient() *openai.Client {
	return client
}
