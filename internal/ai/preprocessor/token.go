package preprocessor

import (
	"github.com/hcd233/Aris-blog/internal/ai/document"
	"github.com/pkoukk/tiktoken-go"
	"github.com/samber/lo"
)

// TokenPreprocessor 基于token的预处理器
//
//	author centonhuang
//	update 2024-12-08 15:20:59
type TokenPreprocessor struct {
	Preprocessor[string, *document.TextDocument]
	chunkSize               uint
	chunkOverlap            uint
	tokenizer               *tiktoken.Tiktoken
	allowedSpecialTokens    []string
	disallowedSpecialTokens []string
}

// NewTokenPreprocessor 创建新的TokenPreprocessor
//
//	param chunkSize token块大小
//	param chunkOverlap token块重叠大小
//	return *TokenPreprocessor
//	return error
func NewTokenPreprocessor(chunkSize, chunkOverlap uint, allowedSpecialTokens, disallowedSpecialTokens []string) *TokenPreprocessor {
	tokenizer := lo.Must1(tiktoken.GetEncoding("cl100k_base"))

	return &TokenPreprocessor{
		chunkSize:               chunkSize,
		chunkOverlap:            chunkOverlap,
		tokenizer:               tokenizer,
		allowedSpecialTokens:    allowedSpecialTokens,
		disallowedSpecialTokens: disallowedSpecialTokens,
	}
}

// Process 处理文本
//
//	receiver p *TokenPreprocessor
//	param source string
//	return []document.TextDocument
//	return error
func (p *TokenPreprocessor) Process(source string) ([]*document.TextDocument, error) {
	tokenIDs := p.tokenizer.Encode(source, p.allowedSpecialTokens, p.disallowedSpecialTokens)
	var chunks []*document.TextDocument

	for i := 0; i < len(tokenIDs); i += int(p.chunkSize - p.chunkOverlap) {
		end := i + int(p.chunkSize)
		if end > len(tokenIDs) {
			end = len(tokenIDs)
		}

		chunk := p.tokenizer.Decode(tokenIDs[i:end])

		doc := document.NewTextDocument(chunk, map[string]interface{}{
			"startToken": i,
			"endToken":   end,
		})
		chunks = append(chunks, doc)

		if end == len(tokenIDs) {
			break
		}
	}

	return chunks, nil
}

// BatchProcess 批量处理文本
//
//	receiver p *TokenPreprocessor
//	param sources []string
//	return []documentT
//	return error
func (p *TokenPreprocessor) BatchProcess(sources []string) ([]*document.TextDocument, error) {
	var allChunks []*document.TextDocument

	for _, source := range sources {
		chunks, err := p.Process(source)
		if err != nil {
			return nil, err
		}
		allChunks = append(allChunks, chunks...)
	}

	return allChunks, nil
}

// ProcessDocument 处理文档
//
//	receiver p *TokenPreprocessor
//	param rawDocument *document.TextDocument
//	return []documentT
//	return error
func (p *TokenPreprocessor) ProcessDocument(rawDocument *document.TextDocument) ([]*document.TextDocument, error) {
	chunks, err := p.Process(rawDocument.Content)
	if err != nil {
		return nil, err
	}

	// 将原始文档的元数据合并到每个chunk中
	for _, chunk := range chunks {
		// 保留chunk自己的token位置信息
		tokenInfo := map[string]interface{}{
			"parentID":   rawDocument.ID,
			"startToken": chunk.Metadata["startToken"],
			"endToken":   chunk.Metadata["endToken"],
		}
		chunk.Metadata = lo.Assign(rawDocument.Metadata, tokenInfo)
	}

	return chunks, nil
}

// BatchProcessDocument 批量处理文档
//
//	receiver p *TokenPreprocessor
//	param rawDocuments []*document.TextDocument
//	return []documentT
//	return error
func (p *TokenPreprocessor) BatchProcessDocument(rawDocuments []*document.TextDocument) ([]*document.TextDocument, error) {
	var allChunks []*document.TextDocument

	for _, doc := range rawDocuments {
		chunks, err := p.ProcessDocument(doc)
		if err != nil {
			return nil, err
		}
		allChunks = append(allChunks, chunks...)
	}

	return allChunks, nil
}
