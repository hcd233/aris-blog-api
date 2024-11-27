package prompt

import (
	"fmt"

	"github.com/hcd233/Aris-blog/internal/ai/message"
)

type MultiTurnPrompt struct {
	prompts []Prompt
}

func NewMultiTurnPrompt(prompts []Prompt) Prompt {
	return &MultiTurnPrompt{
		prompts: prompts,
	}
}

func (mp *MultiTurnPrompt) Format(params map[string]interface{}) (messages []message.Message, err error) {
	var msgs []message.Message
	for _, prompt := range mp.prompts {
		msgs, err = prompt.Format(params)
		if err != nil {
			err = fmt.Errorf("failed to format prompt: %w", err)
			return
		}
		messages = append(messages, msgs...)
	}
	return
}
