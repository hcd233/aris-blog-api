package llm

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewOpenAILLM(t *testing.T) {
	llm := NewOpenAILLM(OpenAIGPT4oMini, 0.7)
	assert.NotNil(t, llm, "LLM instance should not be nil")
}

func TestOpenAILLM_Invoke(t *testing.T) {
	tests := []struct {
		name        string
		prompt      string
		wantErr     bool
		errContains string
	}{
		{
			name:        "empty prompt",
			prompt:      "",
			wantErr:     true,
			errContains: "empty prompt",
		},
		{
			name:    "valid prompt",
			prompt:  "自我介绍一下",
			wantErr: false,
		},
	}

	llm := NewOpenAILLM(OpenAIGPT4oMini, 0.7)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := llm.Invoke(tt.prompt)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, resp)
			}
		})
	}
}

func TestOpenAILLM_Stream(t *testing.T) {
	tests := []struct {
		name        string
		prompt      string
		wantErr     bool
		errContains string
	}{
		{
			name:        "empty prompt",
			prompt:      "",
			wantErr:     true,
			errContains: "empty prompt",
		},
		{
			name:    "valid prompt",
			prompt:  "自我介绍一下",
			wantErr: false,
		},
	}

	llm := NewOpenAILLM(OpenAIGPT4oMini, 0.7)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenChan, errChan, err := llm.Stream(tt.prompt)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, tokenChan)
				assert.Nil(t, errChan)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tokenChan)
				assert.NotNil(t, errChan)
				// Test reading from the stream
				timeout := time.After(5 * time.Second)
				received := false

				for {
					select {
					case msg, ok := <-tokenChan:
						if !ok {
							// Channel closed
							return
						}
						assert.NotEmpty(t, msg)
						received = true
					case err, ok := <-errChan:
						if !ok {
							// Channel closed
							return
						}
						assert.NoError(t, err)
					case <-timeout:
						if !received {
							t.Error("timeout waiting for stream response")
						}
						return
					}
				}
			}
		})
	}
}
