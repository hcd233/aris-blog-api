package chatmodel

import "github.com/hcd233/Aris-blog/internal/ai/message"

type Invokeable interface {
	Invoke(messages []message.Message) (string, error)
}

type Streamable interface {
	Stream(messages []message.Message) (chan string, chan error, error)
}

type ChatModel interface {
	Invokeable
	Streamable
}
