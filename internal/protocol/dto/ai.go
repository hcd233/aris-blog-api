package dto

// Template 提示词模板
type Template struct {
	Role    string `json:"role" doc:"Message role (system/user/assistant)"`
	Content string `json:"content" doc:"Message content"`
}

// Prompt 提示词信息
type Prompt struct {
	PromptID  uint       `json:"promptID" doc:"Prompt ID"`
	CreatedAt string     `json:"createdAt" doc:"Creation timestamp"`
	Task      string     `json:"task" doc:"Task name"`
	Version   uint       `json:"version" doc:"Prompt version"`
	Templates []Template `json:"templates" doc:"Prompt templates"`
	Variables []string   `json:"variables,omitempty" doc:"Template variables"`
}

// TaskPathParam 任务路径参数
type TaskPathParam struct {
	TaskName string `path:"taskName" doc:"Task name" enum:"contentCompletion,articleSummary,articleTranslation,articleQA,termExplaination"`
}

// PromptVersionPathParam 提示词版本路径参数
type PromptVersionPathParam struct {
	Version uint `path:"version" doc:"Prompt version number"`
}

// GetPromptRequest 获取提示词请求
type GetPromptRequest struct {
	TaskPathParam
	PromptVersionPathParam
}

// GetPromptResponse 获取提示词响应
type GetPromptResponse struct {
	Prompt *Prompt `json:"prompt" doc:"Prompt information"`
}

// GetLatestPromptRequest 获取最新提示词请求
type GetLatestPromptRequest struct {
	TaskPathParam
}

// GetLatestPromptResponse 获取最新提示词响应
type GetLatestPromptResponse struct {
	Prompt *Prompt `json:"prompt" doc:"Latest prompt information"`
}

// ListPromptRequest 列出提示词请求
type ListPromptRequest struct {
	TaskPathParam
	CommonParam
}

// ListPromptResponse 列出提示词响应
type ListPromptResponse struct {
	Prompts  []*Prompt `json:"prompts" doc:"List of prompts"`
	PageInfo *PageInfo `json:"pageInfo" doc:"Pagination information"`
}

// CreatePromptRequestBody 创建提示词请求体
type CreatePromptRequestBody struct {
	Templates []Template `json:"templates" doc:"Prompt templates" minItems:"1"`
}

// CreatePromptRequest 创建提示词请求
type CreatePromptRequest struct {
	TaskPathParam
	Body *CreatePromptRequestBody `json:"body" doc:"Fields for creating prompt"`
}

// AIAppRequestBody AI应用基础请求体
type AIAppRequestBody struct {
	Temperature float32 `json:"temperature,omitempty" doc:"Sampling temperature (0-1)" minimum:"0" maximum:"1" default:"0.7"`
}

// GenerateContentCompletionRequestBody 生成内容补全请求体
type GenerateContentCompletionRequestBody struct {
	AIAppRequestBody
	Context     string `json:"context" doc:"Current content context"`
	Instruction string `json:"instruction" doc:"Completion instruction"`
	Reference   string `json:"reference,omitempty" doc:"Reference text"`
}

// GenerateContentCompletionRequest 生成内容补全请求
type GenerateContentCompletionRequest struct {
	Body *GenerateContentCompletionRequestBody `json:"body" doc:"Fields for content completion"`
}

// GenerateArticleSummaryRequestBody 生成文章摘要请求体
type GenerateArticleSummaryRequestBody struct {
	AIAppRequestBody
	ArticleID   uint   `json:"articleID" doc:"Article ID to summarize"`
	Instruction string `json:"instruction" doc:"Summary instruction"`
}

// GenerateArticleSummaryRequest 生成文章摘要请求
type GenerateArticleSummaryRequest struct {
	Body *GenerateArticleSummaryRequestBody `json:"body" doc:"Fields for article summary"`
}

// GenerateArticleQARequestBody 生成文章问答请求体
type GenerateArticleQARequestBody struct {
	AIAppRequestBody
	ArticleID uint   `json:"articleID" doc:"Article ID for Q&A"`
	Question  string `json:"question" doc:"Question about the article"`
}

// GenerateArticleQARequest 生成文章问答请求
type GenerateArticleQARequest struct {
	Body *GenerateArticleQARequestBody `json:"body" doc:"Fields for article Q&A"`
}

// GenerateTermExplainationRequestBody 生成术语解释请求体
type GenerateTermExplainationRequestBody struct {
	AIAppRequestBody
	ArticleID uint   `json:"articleID" doc:"Article ID containing the term"`
	Term      string `json:"term" doc:"Term to explain"`
	Position  int    `json:"position" doc:"Position of term in article"`
}

// GenerateTermExplainationRequest 生成术语解释请求
type GenerateTermExplainationRequest struct {
	Body *GenerateTermExplainationRequestBody `json:"body" doc:"Fields for term explanation"`
}

// SSEResponse SSE流式响应
type SSEResponse struct {
	Delta string `json:"delta" doc:"Incremental token"`
	Stop  bool   `json:"stop" doc:"Whether stream has ended"`
	Error string `json:"error,omitempty" doc:"Error message if any"`
}
