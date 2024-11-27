package prompt

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/hcd233/Aris-blog/internal/ai/message"
)

type OneTurnPrompt struct {
	role     string
	template string
}

func NewOneTurnPrompt(role, template string) *OneTurnPrompt {
	return &OneTurnPrompt{
		role:     role,
		template: template,
	}
}

func (sp *OneTurnPrompt) Format(params map[string]interface{}) (messages []message.Message, err error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	tmpl, err := template.New("simplePrompt").Option("missingkey=error").Parse(sp.template)
	if err != nil {
		err = fmt.Errorf("failed to parse template: %w", err)
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, params)
	if err != nil {
		if strings.Contains(err.Error(), "map has no entry") {
			err = fmt.Errorf("missing required parameter: %w", err)
		} else {
			err = fmt.Errorf("failed to execute template: %w", err)
		}
		return nil, err
	}

	messages = []message.Message{
		{
			Role:    sp.role,
			Content: buf.String(),
		},
	}
	return
}
