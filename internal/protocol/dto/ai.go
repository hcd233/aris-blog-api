package dto

// Template 提示词模板
type Template struct {
    Role    string `json:"role" doc:"消息角色，如 system/user"`
    Content string `json:"content" doc:"提示词内容"`
}

// Prompt 提示词实体
type Prompt struct {
    PromptID  uint       `json:"promptID" doc:"提示词 ID"`
    CreatedAt string     `json:"createdAt" doc:"创建时间"`
    Task      string     `json:"task" doc:"任务名称"`
    Version   uint       `json:"version" doc:"版本号"`
    Templates []Template `json:"templates" doc:"提示词模板列表"`
    Variables []string   `json:"variables" doc:"提示词变量集合"`
}

// TaskPathParam 任务路径参数
type TaskPathParam struct {
    TaskName string `path:"taskName" doc:"任务名称"`
}

// PromptVersionPathParam 提示词版本路径参数
type PromptVersionPathParam struct {
    TaskPathParam
    Version uint `path:"version" minimum:"1" doc:"提示词版本号"`
}

// GetPromptRequest 获取指定版本提示词请求
type GetPromptRequest struct {
    PromptVersionPathParam
}

// GetPromptResponse 获取指定版本提示词响应
type GetPromptResponse struct {
    Prompt *Prompt `json:"prompt" doc:"提示词详情"`
}

// GetLatestPromptRequest 获取最新提示词请求
type GetLatestPromptRequest struct {
    TaskPathParam
}

// GetLatestPromptResponse 获取最新提示词响应
type GetLatestPromptResponse struct {
    Prompt *Prompt `json:"prompt" doc:"提示词详情"`
}

// ListPromptRequest 列出提示词请求
type ListPromptRequest struct {
    TaskPathParam
    CommonParam
}

// ListPromptResponse 列出提示词响应
type ListPromptResponse struct {
    Prompts []*Prompt `json:"prompts" doc:"提示词列表"`
    PageInfo *PageInfo `json:"pageInfo" doc:"分页信息"`
}

// CreatePromptRequestBody 创建提示词请求体
type CreatePromptRequestBody struct {
    Templates []Template `json:"templates" doc:"提示词模板集合"`
}

// CreatePromptRequest 创建提示词请求
type CreatePromptRequest struct {
    TaskPathParam
    Body *CreatePromptRequestBody `json:"body" doc:"创建提示词的请求体"`
}

// AIAppRequestBody AI 应用通用请求体
type AIAppRequestBody struct {
    Temperature float32 `json:"temperature" minimum:"0" maximum:"1" doc:"采样温度，范围 0-1"`
}

// GenerateContentCompletionRequestBody 生成内容补全请求体
type GenerateContentCompletionRequestBody struct {
    AIAppRequestBody
    Context     string `json:"context" doc:"补全上下文"`
    Instruction string `json:"instruction" doc:"补全指令"`
    Reference   string `json:"reference" doc:"参考内容"`
}

// GenerateContentCompletionRequest 生成内容补全请求
type GenerateContentCompletionRequest struct {
    Body *GenerateContentCompletionRequestBody `json:"body" doc:"生成内容补全的请求体"`
}

// GenerateContentCompletionResponse 生成内容补全响应
type GenerateContentCompletionResponse struct {
    TokenChan <-chan string `json:"-"`
    ErrChan   <-chan error  `json:"-"`
}

// GenerateArticleSummaryRequestBody 生成文章摘要请求体
type GenerateArticleSummaryRequestBody struct {
    AIAppRequestBody
    ArticleID   uint   `json:"articleID" doc:"文章 ID"`
    Instruction string `json:"instruction" doc:"摘要指令"`
}

// GenerateArticleSummaryRequest 生成文章摘要请求
type GenerateArticleSummaryRequest struct {
    Body *GenerateArticleSummaryRequestBody `json:"body" doc:"生成文章摘要的请求体"`
}

// GenerateArticleSummaryResponse 生成文章摘要响应
type GenerateArticleSummaryResponse struct {
    TokenChan <-chan string `json:"-"`
    ErrChan   <-chan error  `json:"-"`
}

// GenerateArticleQARequestBody 生成文章问答请求体
type GenerateArticleQARequestBody struct {
    AIAppRequestBody
    ArticleID   uint   `json:"articleID" doc:"文章 ID"`
    Question    string `json:"question" doc:"提问内容"`
}

// GenerateArticleQARequest 生成文章问答请求
type GenerateArticleQARequest struct {
    Body *GenerateArticleQARequestBody `json:"body" doc:"生成文章问答的请求体"`
}

// GenerateArticleQAResponse 生成文章问答响应
type GenerateArticleQAResponse struct {
    TokenChan <-chan string `json:"-"`
    ErrChan   <-chan error  `json:"-"`
}

// GenerateTermExplainationRequestBody 生成术语解释请求体
type GenerateTermExplainationRequestBody struct {
    AIAppRequestBody
    ArticleID uint   `json:"articleID" doc:"文章 ID"`
    Term      string `json:"term" doc:"术语"`
    Position  uint   `json:"position" doc:"术语在正文中的位置"`
}

// GenerateTermExplainationRequest 生成术语解释请求
type GenerateTermExplainationRequest struct {
    Body *GenerateTermExplainationRequestBody `json:"body" doc:"生成术语解释的请求体"`
}

// GenerateTermExplainationResponse 生成术语解释响应
type GenerateTermExplainationResponse struct {
    TokenChan <-chan string `json:"-"`
    ErrChan   <-chan error  `json:"-"`
}

// GenerateArticleTranslationRequest 生成文章翻译请求
type GenerateArticleTranslationRequest struct{}

// GenerateArticleTranslationResponse 生成文章翻译响应
type GenerateArticleTranslationResponse struct {
    TokenChan <-chan string `json:"-"`
    ErrChan   <-chan error  `json:"-"`
}

