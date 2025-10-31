// Package dto AI DTO
package dto

// Template 提示词模板
//
//	author centonhuang
//	update 2025-10-30
type Template struct {
	Role    string `json:"role" doc:"Message role (system/user/assistant)"`
	Content string `json:"content" doc:"Message content"`
}

// Prompt 提示词
//
//	author centonhuang
//	update 2025-10-30
type Prompt struct {
	PromptID  uint       `json:"promptID" doc:"Unique identifier for the prompt"`
	CreatedAt string     `json:"createdAt" doc:"Creation timestamp"`
	Task      string     `json:"task" doc:"Task name"`
	Version   uint       `json:"version" doc:"Version number"`
	Templates []Template `json:"templates" doc:"Prompt templates"`
	Variables []string   `json:"variables" doc:"Template variables"`
}

// GetPromptRequest 获取提示词请求
//
//	author centonhuang
//	update 2025-10-30
type GetPromptRequest struct {
	TaskName string `path:"taskName" enum:"contentCompletion,articleSummary,articleTranslation,articleQA,termExplaination" doc:"Task name"`
	Version  uint   `path:"version" minimum:"1" doc:"Version number"`
}

// GetPromptResponse 获取提示词响应
//
//	author centonhuang
//	update 2025-10-30
type GetPromptResponse struct {
	Prompt *Prompt `json:"prompt" doc:"Prompt information"`
}

// GetLatestPromptRequest 获取最新提示词请求
//
//	author centonhuang
//	update 2025-10-30
type GetLatestPromptRequest struct {
	TaskName string `path:"taskName" enum:"contentCompletion,articleSummary,articleTranslation,articleQA,termExplaination" doc:"Task name"`
}

// GetLatestPromptResponse 获取最新提示词响应
//
//	author centonhuang
//	update 2025-10-30
type GetLatestPromptResponse struct {
	Prompt *Prompt `json:"prompt" doc:"Latest prompt information"`
}

// ListPromptRequest 列出提示词请求
//
//	author centonhuang
//	update 2025-10-30
type ListPromptRequest struct {
	TaskName string `path:"taskName" enum:"contentCompletion,articleSummary,articleTranslation,articleQA,termExplaination" doc:"Task name"`
	Page     *int   `query:"page" minimum:"1" doc:"Page number (starting from 1)"`
	PageSize *int   `query:"pageSize" minimum:"1" maximum:"50" doc:"Number of items per page"`
}

// ListPromptResponse 列出提示词响应
//
//	author centonhuang
//	update 2025-10-30
type ListPromptResponse struct {
	Prompts  []*Prompt `json:"prompts" doc:"List of prompts"`
	PageInfo *PageInfo `json:"pageInfo" doc:"Pagination information"`
}

// CreatePromptRequest 创建提示词请求
//
//	author centonhuang
//	update 2025-10-30
type CreatePromptRequest struct {
	TaskName string              `path:"taskName" enum:"contentCompletion,articleSummary,articleTranslation,articleQA,termExplaination" doc:"Task name"`
	Body     *CreatePromptBody `json:"body" doc:"Request body containing prompt templates"`
}

// CreatePromptBody 创建提示词请求体
//
//	author centonhuang
//	update 2025-10-30
type CreatePromptBody struct {
	Templates []Template `json:"templates" minItems:"1" doc:"List of prompt templates (at least one required)"`
}

// CreatePromptResponse 创建提示词响应
//
//	author centonhuang
//	update 2025-10-30
type CreatePromptResponse struct{}

// GenerateContentCompletionRequest 生成内容补全请求
//
//	author centonhuang
//	update 2025-10-30
type GenerateContentCompletionRequest struct {
	Body *GenerateContentCompletionBody `json:"body" doc:"Request body containing generation parameters"`
}

// GenerateContentCompletionBody 生成内容补全请求体
//
//	author centonhuang
//	update 2025-10-30
type GenerateContentCompletionBody struct {
	Context     string  `json:"context" doc:"Context for content generation"`
	Instruction string  `json:"instruction" doc:"Generation instruction"`
	Reference   string  `json:"reference,omitempty" doc:"Reference text (optional)"`
	Temperature float32 `json:"temperature,omitempty" minimum:"0" maximum:"1" doc:"Generation temperature (0-1)"`
}

// GenerateArticleSummaryRequest 生成文章摘要请求
//
//	author centonhuang
//	update 2025-10-30
type GenerateArticleSummaryRequest struct {
	Body *GenerateArticleSummaryBody `json:"body" doc:"Request body containing generation parameters"`
}

// GenerateArticleSummaryBody 生成文章摘要请求体
//
//	author centonhuang
//	update 2025-10-30
type GenerateArticleSummaryBody struct {
	ArticleID   uint    `json:"articleID" doc:"ID of the article to summarize"`
	Instruction string  `json:"instruction" doc:"Summary instruction"`
	Temperature float32 `json:"temperature,omitempty" minimum:"0" maximum:"1" doc:"Generation temperature (0-1)"`
}

// GenerateArticleQARequest 生成文章问答请求
//
//	author centonhuang
//	update 2025-10-30
type GenerateArticleQARequest struct {
	Body *GenerateArticleQABody `json:"body" doc:"Request body containing generation parameters"`
}

// GenerateArticleQABody 生成文章问答请求体
//
//	author centonhuang
//	update 2025-10-30
type GenerateArticleQABody struct {
	ArticleID   uint    `json:"articleID" doc:"ID of the article"`
	Question    string  `json:"question" doc:"Question about the article"`
	Temperature float32 `json:"temperature,omitempty" minimum:"0" maximum:"1" doc:"Generation temperature (0-1)"`
}

// GenerateTermExplainationRequest 生成术语解释请求
//
//	author centonhuang
//	update 2025-10-30
type GenerateTermExplainationRequest struct {
	Body *GenerateTermExplainationBody `json:"body" doc:"Request body containing generation parameters"`
}

// GenerateTermExplainationBody 生成术语解释请求体
//
//	author centonhuang
//	update 2025-10-30
type GenerateTermExplainationBody struct {
	ArticleID   uint    `json:"articleID" doc:"ID of the article"`
	Term        string  `json:"term" doc:"Term to explain"`
	Position    uint    `json:"position" doc:"Position of the term in the article"`
	Temperature float32 `json:"temperature,omitempty" minimum:"0" maximum:"1" doc:"Generation temperature (0-1)"`
}
