package memory

import (
	"github.com/hcd233/aris-blog-api/internal/ai/message"
)

type Memory interface {
	ToMessage() []message.Message
}
