package memory

import (
	"github.com/hcd233/Aris-blog/internal/ai/message"
)

type Memory interface {
	ToMessage() []message.Message
}
