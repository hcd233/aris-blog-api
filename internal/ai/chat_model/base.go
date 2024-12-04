package chatmodel

import "github.com/hcd233/Aris-blog/internal/ai/message"

type Invokeable interface {
	Invoke(messages []message.Message) (sequence string, err error)
}

type Streamable interface {
	Stream(messages []message.Message) (tokenChan chan string, errChan chan error, err error)
}

type ChatModel interface {
	Invokeable
	Streamable
}
