package prompt

import (
	"github.com/hcd233/Aris-blog/internal/ai/message"
)

type Prompt interface {
	Format(params map[string]interface{}) (messages []message.Message, err error)
}
