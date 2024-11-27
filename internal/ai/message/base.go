package message

import "fmt"

type Message struct {
	Role    string
	Content string
}

func (m *Message) String() string {
	return fmt.Sprintf("%s: %s", m.Role, m.Content)
}
