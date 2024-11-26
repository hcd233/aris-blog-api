package llm

import (
	"context"
	"fmt"
	"io"

	"github.com/hcd233/Aris-blog/internal/config"
	openai "github.com/sashabaranov/go-openai"
)

type OpenAIModel string

const (
	OpenAIGPT4oMini OpenAIModel = "gpt-4o-mini"
)

type OpenAILLM struct {
	client      *openai.Client
	model       OpenAIModel
	temperature float64
}

func NewOpenAILLM(model OpenAIModel, temperature float64) LLM {
	clientConfig := openai.DefaultConfig(config.OpenAIAPIKey)
	clientConfig.BaseURL = config.OpenAIBaseURL
	client := openai.NewClientWithConfig(clientConfig)
	return &OpenAILLM{
		client:      client,
		model:       model,
		temperature: temperature,
	}
}

func (o *OpenAILLM) Invoke(prompt string) (string, error) {
	if prompt == "" {
		return "", fmt.Errorf("empty prompt")
	}

	resp, err := o.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: string(o.model),
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: float32(o.temperature),
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to create chat completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no completion choices returned")
	}

	return resp.Choices[0].Message.Content, nil
}

func (o *OpenAILLM) Stream(prompt string) (chan string, chan error, error) {
	if prompt == "" {
		return nil, nil, fmt.Errorf("empty prompt")
	}

	stream, err := o.client.CreateChatCompletionStream(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: string(o.model),
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: float32(o.temperature),
			Stream:      true,
		},
	)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to create chat completion stream: %w", err)
	}

	tokenChan := make(chan string)
	errChan := make(chan error)

	go func() {
		defer close(tokenChan)
		defer close(errChan)
		defer stream.Close()

		for {
			response, err := stream.Recv()
			if err != nil {
				if err != io.EOF {
					errChan <- fmt.Errorf("stream error: %w", err)
				}
				return
			}

			if len(response.Choices) > 0 {
				content := response.Choices[0].Delta.Content
				if content != "" {
					tokenChan <- content
				}
			}
		}
	}()

	return tokenChan, errChan, nil
}
