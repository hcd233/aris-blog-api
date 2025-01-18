package chatmodel

import "github.com/hcd233/aris-blog-api/internal/ai/message"

// Invokable 可调用对象
//
//	author centonhuang
//	update 2025-01-18 22:13:53
type Invokable interface {
	Invoke(messages []message.Message) (sequence string, err error)
}

// Streamable 可流输对象
//
//	author centonhuang
//	update 2025-01-18 22:13:56
type Streamable interface {
	Stream(messages []message.Message) (tokenChan chan string, errChan chan error, err error)
}

// ChatModel 对话模型
//
//	author centonhuang
//	update 2025-01-18 22:13:57
type ChatModel interface {
	Invokable
	Streamable
}
