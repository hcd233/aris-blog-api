package prompt

import (
	"github.com/hcd233/aris-blog-api/internal/ai/message"
)

type Prompt interface {
	Format(params map[string]interface{}) (messages []message.Message, err error)
}
