package prompt_test

import (
	"testing"

	"github.com/hcd233/aris-blog-api/internal/ai/prompt"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

const oneTurnTemplate = "Hello. I am {{.name}}"

var multiTurnTemplates = []string{
	"Hello. I am {{.name}}",
	"I am {{.age}} years old",
}

func TestOneTurnPrompt(t *testing.T) {
	t.Run("successful format", func(t *testing.T) {
		sp := prompt.NewOneTurnPrompt("user", oneTurnTemplate)
		messages, err := sp.Format(map[string]interface{}{
			"name": "Aris",
		})
		assert.NoError(t, err)
		assert.Equal(t, "user: Hello. I am Aris", messages[0].String())
	})

	t.Run("invalid template", func(t *testing.T) {
		sp := prompt.NewOneTurnPrompt("user", "Hello {{.name")
		messages, err := sp.Format(map[string]interface{}{
			"name": "Aris",
		})
		assert.Error(t, err)
		assert.Nil(t, messages)
	})

	t.Run("missing parameter", func(t *testing.T) {
		sp := prompt.NewOneTurnPrompt("user", oneTurnTemplate)
		messages, err := sp.Format(map[string]interface{}{})
		assert.Error(t, err)
		assert.Nil(t, messages)
	})

	t.Run("invalid parameter type", func(t *testing.T) {
		sp := prompt.NewOneTurnPrompt("user", "I am {{.age}} years old")
		messages, err := sp.Format(map[string]interface{}{
			"age": "not a number",
		})
		assert.NoError(t, err)
		assert.Equal(t, "user: I am not a number years old", messages[0].String())
	})

	t.Run("nil parameters", func(t *testing.T) {
		sp := prompt.NewOneTurnPrompt("user", oneTurnTemplate)
		messages, err := sp.Format(nil)
		assert.Error(t, err)
		assert.Nil(t, messages)
	})
}

func TestMultiTurnPrompt(t *testing.T) {
	t.Run("successful format", func(t *testing.T) {
		prompts := lo.Map(multiTurnTemplates, func(template string, _ int) prompt.Prompt {
			return prompt.NewOneTurnPrompt("user", template)
		})
		mp := prompt.NewMultiTurnPrompt(prompts)
		messages, err := mp.Format(map[string]interface{}{
			"name": "Aris",
			"age":  18,
		})
		assert.NoError(t, err)
		assert.Equal(t, "user: Hello. I am Aris", messages[0].String())
		assert.Equal(t, "user: I am 18 years old", messages[1].String())
	})

	t.Run("empty prompts", func(t *testing.T) {
		mp := prompt.NewMultiTurnPrompt([]prompt.Prompt{})
		messages, err := mp.Format(map[string]interface{}{})
		assert.NoError(t, err)
		assert.Empty(t, messages)
	})

	t.Run("partial parameter missing", func(t *testing.T) {
		prompts := lo.Map(multiTurnTemplates, func(template string, _ int) prompt.Prompt {
			return prompt.NewOneTurnPrompt("user", template)
		})
		mp := prompt.NewMultiTurnPrompt(prompts)
		messages, err := mp.Format(map[string]interface{}{
			"name": "Aris",
			// age is missing
		})
		assert.Error(t, err)
		assert.Nil(t, messages)
	})

	t.Run("nil parameters", func(t *testing.T) {
		prompts := lo.Map(multiTurnTemplates, func(template string, _ int) prompt.Prompt {
			return prompt.NewOneTurnPrompt("user", template)
		})
		mp := prompt.NewMultiTurnPrompt(prompts)
		messages, err := mp.Format(nil)
		assert.Error(t, err)
		assert.Nil(t, messages)
	})

	t.Run("invalid template in sequence", func(t *testing.T) {
		prompts := []prompt.Prompt{
			prompt.NewOneTurnPrompt("user", "Valid {{.name}}"),
			prompt.NewOneTurnPrompt("user", "Invalid {{.age"),
		}
		mp := prompt.NewMultiTurnPrompt(prompts)
		messages, err := mp.Format(map[string]interface{}{
			"name": "Aris",
			"age":  18,
		})
		assert.Error(t, err)
		assert.Nil(t, messages)
	})
}
