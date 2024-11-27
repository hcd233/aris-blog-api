package prompt

import "github.com/hcd233/Aris-blog/internal/ai/message"

type EmptyMemory struct {
	messages []message.Message
}

func NewEmptyMemory() Memory {
	return &EmptyMemory{}
}

func (em *EmptyMemory) AddMessage(message message.Message) {
	em.messages = append(em.messages, message)
}

func (em *EmptyMemory) AddMessages(messages []message.Message) {
	em.messages = append(em.messages, messages...)
}

func (em *EmptyMemory) AddSystemContent(content string) {
	em.messages = append(em.messages, message.Message{
		Role:    "system",
		Content: content,
	})
}

func (em *EmptyMemory) AddAIContent(content string) {
	em.messages = append(em.messages, message.Message{
		Role:    "assistant",
		Content: content,
	})
}

func (em *EmptyMemory) AddUserContent(content string) {
	em.messages = append(em.messages, message.Message{
		Role:    "user",
		Content: content,
	})
}

func (em *EmptyMemory) ToMessage() []message.Message {
	return em.messages
}
