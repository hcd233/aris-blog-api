package chatmodel

import (
	"context"
	"fmt"
	"io"

	"github.com/hcd233/Aris-blog/internal/ai/message"
	"github.com/hcd233/Aris-blog/internal/resource/llm"
	"github.com/samber/lo"
	openai "github.com/sashabaranov/go-openai"
)

type OpenAIModel string

const (
	OpenAIGPT4oMini OpenAIModel = "gpt-4o-mini"
)

type ChatOpenAI struct {
	client      *openai.Client
	model       OpenAIModel
	temperature float64
}

func NewChatOpenAI(model OpenAIModel, temperature float64) ChatModel {
	return &ChatOpenAI{
		client:      llm.GetOpenAIClient(),
		model:       model,
		temperature: temperature,
	}
}

func (o *ChatOpenAI) Invoke(messages []message.Message) (sequence string, err error) {
	if len(messages) == 0 {
		return "", fmt.Errorf("empty messages")
	}

	resp, err := o.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: string(o.model),
			Messages: lo.Map(messages, func(message message.Message, idx int) openai.ChatCompletionMessage {
				return openai.ChatCompletionMessage{
					Role:    message.Role,
					Content: message.Content,
				}
			}),
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

func (o *ChatOpenAI) Stream(messages []message.Message) (tokenChan chan string, errChan chan error, err error) {
	if len(messages) == 0 {
		return nil, nil, fmt.Errorf("empty messages")
	}

	stream, err := o.client.CreateChatCompletionStream(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: string(o.model),
			Messages: lo.Map(messages, func(message message.Message, idx int) openai.ChatCompletionMessage {
				return openai.ChatCompletionMessage{
					Role:    message.Role,
					Content: message.Content,
				}
			}),
			Temperature: float32(o.temperature),
			Stream:      true,
		},
	)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to create chat completion stream: %w", err)
	}

	tokenChan = make(chan string)
	errChan = make(chan error)

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
