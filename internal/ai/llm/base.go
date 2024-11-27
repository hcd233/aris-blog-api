package llm

import "github.com/hcd233/Aris-blog/internal/ai/message"

type LLM interface {
	Invoke(messages []message.Message) (string, error)
	Stream(messages []message.Message) (chan string, chan error, error)
}
