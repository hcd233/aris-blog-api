package humahandler

import (
	"context"

	humadto "github.com/hcd233/aris-blog-api/internal/huma"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
)

type AIHandlers struct{ svc service.AIService }

func NewAIHandlers() *AIHandlers { return &AIHandlers{svc: service.NewAIService()} }

type (
	taskPathInput          struct{ humadto.TaskPath }
	promptVersionPathInput struct{ humadto.PromptVersionPath }
	createPromptInput      struct {
		taskPathInput
		humadto.CreatePromptInput
	}
	listPromptInput struct {
		taskPathInput
		humadto.PaginateParam
	}
	contentCompletionInput struct {
		authHeader
		humadto.GenerateContentCompletionInput
	}
	articleSummaryInput struct {
		authHeader
		humadto.GenerateArticleSummaryInput
	}
	articleQAInput struct {
		authHeader
		humadto.GenerateArticleQAInput
	}
	termExplainationInput struct {
		authHeader
		humadto.GenerateTermExplainationInput
	}
)

func (h *AIHandlers) HandleGetPrompt(ctx context.Context, input *promptVersionPathInput) (*humadto.Output[protocol.GetPromptResponse], error) {
	req := &protocol.GetPromptRequest{TaskName: input.TaskName, Version: input.Version}
	rsp, err := h.svc.GetPrompt(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.GetPromptResponse]{Body: *rsp}, nil
}

func (h *AIHandlers) HandleGetLatestPrompt(ctx context.Context, input *taskPathInput) (*humadto.Output[protocol.GetLatestPromptResponse], error) {
	req := &protocol.GetLatestPromptRequest{TaskName: input.TaskName}
	rsp, err := h.svc.GetLatestPrompt(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.GetLatestPromptResponse]{Body: *rsp}, nil
}

func (h *AIHandlers) HandleListPrompt(ctx context.Context, input *listPromptInput) (*humadto.Output[protocol.ListPromptResponse], error) {
	p := &protocol.PaginateParam{PageParam: &protocol.PageParam{}, QueryParam: &protocol.QueryParam{}}
	if input.PageParam != nil {
		p.PageParam = &protocol.PageParam{Page: input.Page, PageSize: input.PageSize}
	}
	if input.QueryParam != nil {
		p.QueryParam = &protocol.QueryParam{Query: input.Query}
	}
	req := &protocol.ListPromptRequest{TaskName: input.TaskName, PaginateParam: p}
	rsp, err := h.svc.ListPrompt(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.ListPromptResponse]{Body: *rsp}, nil
}

func (h *AIHandlers) HandleCreatePrompt(ctx context.Context, input *createPromptInput) (*humadto.Output[protocol.CreatePromptResponse], error) {
	// 类型转换 humadto.Template -> protocol.Template
	var templates []protocol.Template
	if len(input.Body.Templates) > 0 {
		templates = make([]protocol.Template, len(input.Body.Templates))
		for i, t := range input.Body.Templates {
			templates[i] = protocol.Template{Role: t.Role, Content: t.Content}
		}
	}
	req := &protocol.CreatePromptRequest{TaskName: input.TaskName, Templates: templates}
	rsp, err := h.svc.CreatePrompt(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.CreatePromptResponse]{Body: *rsp}, nil
}

func (h *AIHandlers) HandleGenerateContentCompletion(ctx context.Context, input *contentCompletionInput) (*protocol.GenerateContentCompletionResponse, error) {
	req := &protocol.GenerateContentCompletionRequest{UserID: input.UserID, Context: input.Body.Context, Instruction: input.Body.Instruction, Reference: input.Body.Reference, Temperature: input.Body.Temperature}
	return h.svc.GenerateContentCompletion(ctx, req)
}

func (h *AIHandlers) HandleGenerateArticleSummary(ctx context.Context, input *articleSummaryInput) (*protocol.GenerateArticleSummaryResponse, error) {
	req := &protocol.GenerateArticleSummaryRequest{UserID: input.UserID, ArticleID: input.Body.ArticleID, Instruction: input.Body.Instruction, Temperature: input.Body.Temperature}
	return h.svc.GenerateArticleSummary(ctx, req)
}

func (h *AIHandlers) HandleGenerateArticleTranslation(ctx context.Context, _ *authHeader) (*protocol.GenerateArticleTranslationResponse, error) {
	req := &protocol.GenerateArticleTranslationRequest{}
	return h.svc.GenerateArticleTranslation(ctx, req)
}

func (h *AIHandlers) HandleGenerateArticleQA(ctx context.Context, input *articleQAInput) (*protocol.GenerateArticleQAResponse, error) {
	req := &protocol.GenerateArticleQARequest{UserID: input.UserID, ArticleID: input.Body.ArticleID, Question: input.Body.Question, Temperature: input.Body.Temperature}
	return h.svc.GenerateArticleQA(ctx, req)
}

func (h *AIHandlers) HandleGenerateTermExplaination(ctx context.Context, input *termExplainationInput) (*protocol.GenerateTermExplainationResponse, error) {
	req := &protocol.GenerateTermExplainationRequest{UserID: input.UserID, ArticleID: input.Body.ArticleID, Term: input.Body.Term, Position: input.Body.Position, Temperature: input.Body.Temperature}
	return h.svc.GenerateTermExplaination(ctx, req)
}
