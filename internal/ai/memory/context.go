package memory

import "github.com/hcd233/Aris-blog/internal/ai/message"

type ContextMemory struct {
	messages []message.Message
}

func NewContextMemory() Memory {
	return &ContextMemory{}
}

func (cm *ContextMemory) AddMessage(message message.Message) {
	cm.messages = append(cm.messages, message)
}

func (cm *ContextMemory) AddMessages(messages []message.Message) {
	cm.messages = append(cm.messages, messages...)
}

func (cm *ContextMemory) AddSystemContent(content string) {
	cm.messages = append(cm.messages, message.Message{
		Role:    "system",
		Content: content,
	})
}

func (cm *ContextMemory) AddAIContent(content string) {
	cm.messages = append(cm.messages, message.Message{
		Role:    "assistant",
		Content: content,
	})
}

func (cm *ContextMemory) AddUserContent(content string) {
	cm.messages = append(cm.messages, message.Message{
		Role:    "user",
		Content: content,
	})
}

func (cm *ContextMemory) ToMessage() []message.Message {
	return cm.messages
}
