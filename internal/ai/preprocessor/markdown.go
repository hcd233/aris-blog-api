package preprocessor

import (
	"bytes"
	"strings"

	"github.com/hcd233/Aris-blog/internal/ai/document"
	"github.com/pkoukk/tiktoken-go"
	"github.com/samber/lo"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

const (
	defaultChunkSize      = 1024
	defaultChunkOverlap   = 256
	defaultMinHeaderLevel = 1 // 默认从 h1 开始分割
)

// MarkdownPreprocessor markdown预处理器
//
//	author centonhuang
//	update 2024-12-08 14:54:54
type MarkdownPreprocessor struct {
	Preprocessor[string, *document.TextDocument]
	chunkSize               uint
	chunkOverlap            uint
	tokenizer               *tiktoken.Tiktoken
	allowedSpecialTokens    []string
	disallowedSpecialTokens []string
	minHeaderLevel          uint
}

// NewMarkdownPreprocessor 创建markdown预处理器
//
//	param chunkSize uint 块大小
//	param chunkOverlap uint 块重叠
//	param minHeaderLevel uint 最小标题等级（1-6），0表示使用默认值
//	return *MarkdownPreprocessor
//	author centonhuang
//	update 2024-12-08 14:57:33
func NewMarkdownPreprocessor(chunkSize, chunkOverlap uint, allowedSpecialTokens, disallowedSpecialTokens []string) *MarkdownPreprocessor {
	if chunkSize == 0 {
		chunkSize = defaultChunkSize
	}
	if chunkOverlap == 0 {
		chunkOverlap = defaultChunkOverlap
	}

	tokenizer := lo.Must1(tiktoken.GetEncoding("cl100k_base"))

	return &MarkdownPreprocessor{
		chunkSize:               chunkSize,
		chunkOverlap:            chunkOverlap,
		tokenizer:               tokenizer,
		allowedSpecialTokens:    allowedSpecialTokens,
		disallowedSpecialTokens: disallowedSpecialTokens,
	}
}

// Process 处理markdown文档
//
//	receiver p *MarkdownPreprocessor
//	param source string
//	return documents []*document.TextDocument
//	return err error
//	author centonhuang
//	update 2024-12-08 15:02:45
func (p *MarkdownPreprocessor) Process(source string) (documents []*document.TextDocument, err error) {
	md := goldmark.New()
	reader := text.NewReader([]byte(source))
	node := md.Parser().Parse(reader)

	var chunks []string
	var curChunk strings.Builder
	var currentTokens []int // 用于跟踪当前chunk的tokens
	var currentHeader string
	var headers []string // 用于存储每个chunk对应的header
	var headerLevel uint // 当前标题等级

	// 遍历AST
	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		// 如果是标题节点，检查等级并决定是否创建新的chunk
		if heading, ok := n.(*ast.Heading); ok {
			headerLevel = uint(heading.Level)

			// 只有当标题等级小于等于最小标题等级时才分割
			if headerLevel <= p.minHeaderLevel {
				if len(currentTokens) > 0 {
					chunks = append(chunks, curChunk.String())
					headers = append(headers, currentHeader)
					curChunk.Reset()
					currentTokens = nil
				}

				// 获取标题文本
				var headerBuf bytes.Buffer
				if err := md.Renderer().Render(&headerBuf, []byte(source), n); err != nil {
					return ast.WalkStop, err
				}
				currentHeader = strings.TrimSpace(headerBuf.String())
			}
		}

		// 获取节点的文本内容
		var buf bytes.Buffer
		if err := md.Renderer().Render(&buf, []byte(source), n); err != nil {
			return ast.WalkStop, err
		}
		content := buf.String()

		// 使用tokenizer计算新内容的tokens
		newTokens := p.tokenizer.Encode(content, p.allowedSpecialTokens, p.disallowedSpecialTokens)

		// 如果当前chunk加上新内容超过chunkSize，创建新的chunk
		if len(currentTokens)+len(newTokens) > int(p.chunkSize) {
			if len(currentTokens) > 0 {
				chunks = append(chunks, curChunk.String())
				headers = append(headers, currentHeader)

				// 处理重叠
				lastContent := curChunk.String()
				curChunk.Reset()

				// 计算重叠部分的tokens
				lastTokens := p.tokenizer.Encode(lastContent, p.allowedSpecialTokens, p.disallowedSpecialTokens)
				if len(lastTokens) > int(p.chunkOverlap) {
					overlapText := p.tokenizer.Decode(lastTokens[len(lastTokens)-int(p.chunkOverlap):])
					curChunk.WriteString(overlapText)
					currentTokens = p.tokenizer.Encode(overlapText, p.allowedSpecialTokens, p.disallowedSpecialTokens)
				} else {
					currentTokens = nil
				}
			}
		}

		curChunk.WriteString(content)
		if currentTokens == nil {
			currentTokens = newTokens
		} else {
			currentTokens = append(currentTokens, newTokens...)
		}

		return ast.WalkContinue, nil
	})

	// 添加最后一个chunk
	if len(currentTokens) > 0 {
		chunks = append(chunks, curChunk.String())
		headers = append(headers, currentHeader)
	}

	// 创建文档
	documents = make([]*document.TextDocument, len(chunks))
	for i, chunk := range chunks {
		metadata := map[string]interface{}{
			"index": i,
			"total": len(chunks),
		}

		if headers[i] != "" {
			metadata["header"] = headers[i]
		}

		documents[i] = &document.TextDocument{
			Content:  strings.TrimSpace(chunk),
			Metadata: metadata,
		}
	}

	return documents, nil
}

// BatchProcess 批量处理文本
//
//	receiver p *TokenPreprocessor
//	param sources []string
//	return []documentT
//	return error
func (p *MarkdownPreprocessor) BatchProcess(sources []string) ([]*document.TextDocument, error) {
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
func (p *MarkdownPreprocessor) ProcessDocument(rawDocument *document.TextDocument) ([]*document.TextDocument, error) {
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
func (p *MarkdownPreprocessor) BatchProcessDocument(rawDocuments []*document.TextDocument) ([]*document.TextDocument, error) {
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
