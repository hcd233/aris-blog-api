package preprocessor

import (
	"strings"
	"testing"

	"github.com/hcd233/Aris-blog/internal/ai/document"
	"github.com/stretchr/testify/assert"
)

func TestNewTokenPreprocessor(t *testing.T) {
	tests := []struct {
		name              string
		chunkSize         uint
		chunkOverlap      uint
		allowedSpecial    []string
		disallowedSpecial []string
		expectError       bool
	}{
		{
			name:              "正常创建",
			chunkSize:         100,
			chunkOverlap:      20,
			allowedSpecial:    []string{},
			disallowedSpecial: []string{},
			expectError:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewTokenPreprocessor(tt.chunkSize, tt.chunkOverlap, tt.allowedSpecial, tt.disallowedSpecial)
			if tt.expectError {
				assert.Nil(t, p)
			} else {
				assert.NotNil(t, p)
				assert.Equal(t, tt.chunkSize, p.chunkSize)
				assert.Equal(t, tt.chunkOverlap, p.chunkOverlap)
			}
		})
	}
}

func TestTokenPreprocessor_Process(t *testing.T) {
	tests := []struct {
		name          string
		source        string
		chunkSize     uint
		chunkOverlap  uint
		expectedLen   int
		checkMetadata bool
	}{
		{
			name:         "空文本",
			source:       "",
			chunkSize:    100,
			chunkOverlap: 20,
			expectedLen:  0,
		},
		{
			name:          "短文本-不分块",
			source:        "这是一个测试文本",
			chunkSize:     100,
			chunkOverlap:  20,
			expectedLen:   1,
			checkMetadata: true,
		},
		{
			name:          "长文本-需要分块",
			source:        "这是一个很长的测试文本。" + strings.Repeat("需要被分成多个块。", 50),
			chunkSize:     50,
			chunkOverlap:  10,
			expectedLen:   13,
			checkMetadata: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewTokenPreprocessor(tt.chunkSize, tt.chunkOverlap, nil, nil)

			chunks, err := p.Process(tt.source)
			assert.NoError(t, err)
			assert.Len(t, chunks, tt.expectedLen)

			if tt.checkMetadata {
				for _, chunk := range chunks {
					assert.Contains(t, chunk.Metadata, "startToken")
					assert.Contains(t, chunk.Metadata, "endToken")
					start := chunk.Metadata["startToken"].(int)
					end := chunk.Metadata["endToken"].(int)
					assert.Less(t, start, end)
				}
			}
		})
	}
}

func TestTokenPreprocessor_ProcessDocument(t *testing.T) {
	tests := []struct {
		name         string
		doc          *document.TextDocument
		chunkSize    uint
		chunkOverlap uint
	}{
		{
			name: "处理带元数据的文档",
			doc: &document.TextDocument{
				ID:      "test-doc-1",
				Content: "这是一个测试文档，包含一些元数据。",
				Metadata: map[string]interface{}{
					"source": "test",
					"type":   "article",
				},
			},
			chunkSize:    50,
			chunkOverlap: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewTokenPreprocessor(tt.chunkSize, tt.chunkOverlap, nil, nil)

			chunks, err := p.ProcessDocument(tt.doc)
			assert.NoError(t, err)
			assert.NotEmpty(t, chunks)

			// 检查每个chunk
			for _, chunk := range chunks {
				// 检查原始元数据是否保留
				assert.Equal(t, tt.doc.Metadata["source"], chunk.Metadata["source"])
				assert.Equal(t, tt.doc.Metadata["type"], chunk.Metadata["type"])

				// 检查是否添加了必要的token信息
				assert.Contains(t, chunk.Metadata, "parentID")
				assert.Equal(t, tt.doc.ID, chunk.Metadata["parentID"])
				assert.Contains(t, chunk.Metadata, "startToken")
				assert.Contains(t, chunk.Metadata, "endToken")
			}
		})
	}
}

func TestTokenPreprocessor_BatchProcess(t *testing.T) {
	sources := []string{
		"第一个测试文本",
		"第二个测试文本",
		"第三个测试文本",
	}

	p := NewTokenPreprocessor(100, 20, nil, nil)

	chunks, err := p.BatchProcess(sources)
	assert.NoError(t, err)
	assert.Equal(t, len(sources), len(chunks))
}

func TestTokenPreprocessor_BatchProcessDocument(t *testing.T) {
	docs := []*document.TextDocument{
		{
			ID:      "doc-1",
			Content: "第一个文档",
			Metadata: map[string]interface{}{
				"source": "test-1",
			},
		},
		{
			ID:      "doc-2",
			Content: "第二个文档",
			Metadata: map[string]interface{}{
				"source": "test-2",
			},
		},
	}

	p := NewTokenPreprocessor(100, 20, nil, nil)

	chunks, err := p.BatchProcessDocument(docs)
	assert.NoError(t, err)
	assert.NotEmpty(t, chunks)

	// 检查每个chunk都有正确的父文档ID
	for _, chunk := range chunks {
		assert.Contains(t, chunk.Metadata, "parentID")
		parentID := chunk.Metadata["parentID"].(string)
		assert.Contains(t, []string{"doc-1", "doc-2"}, parentID)
	}
}
