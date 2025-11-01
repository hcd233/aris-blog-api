package handler

import (
	"context"

	"github.com/hcd233/aris-blog-api/internal/protocol"
	dto "github.com/hcd233/aris-blog-api/internal/protocol/dto"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// AIHandler AI处理器
type AIHandler interface {
	HandleGetPrompt(ctx context.Context, req *dto.GetPromptRequest) (*protocol.HumaHTTPResponse[*dto.GetPromptResponse], error)
	HandleGetLatestPrompt(ctx context.Context, req *dto.GetLatestPromptRequest) (*protocol.HumaHTTPResponse[*dto.GetLatestPromptResponse], error)
	HandleListPrompt(ctx context.Context, req *dto.ListPromptRequest) (*protocol.HumaHTTPResponse[*dto.ListPromptResponse], error)
	HandleCreatePrompt(ctx context.Context, req *dto.CreatePromptRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error)
	// SSE streaming methods - will return special responses
	HandleGenerateContentCompletion(ctx context.Context, req *dto.GenerateContentCompletionRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error)
	HandleGenerateArticleSummary(ctx context.Context, req *dto.GenerateArticleSummaryRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error)
	HandleGenerateArticleQA(ctx context.Context, req *dto.GenerateArticleQARequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error)
	HandleGenerateTermExplaination(ctx context.Context, req *dto.GenerateTermExplainationRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error)
}

type aiHandler struct {
	svc service.AIService
}

// NewAIHandler 创建AI处理器
func NewAIHandler() AIHandler {
	return &aiHandler{
		svc: service.NewAIService(),
	}
}

func (h *aiHandler) HandleGetPrompt(ctx context.Context, req *dto.GetPromptRequest) (*protocol.HumaHTTPResponse[*dto.GetPromptResponse], error) {
	return util.WrapHTTPResponse(h.svc.GetPrompt(ctx, req))
}

func (h *aiHandler) HandleGetLatestPrompt(ctx context.Context, req *dto.GetLatestPromptRequest) (*protocol.HumaHTTPResponse[*dto.GetLatestPromptResponse], error) {
	return util.WrapHTTPResponse(h.svc.GetLatestPrompt(ctx, req))
}

func (h *aiHandler) HandleListPrompt(ctx context.Context, req *dto.ListPromptRequest) (*protocol.HumaHTTPResponse[*dto.ListPromptResponse], error) {
	return util.WrapHTTPResponse(h.svc.ListPrompt(ctx, req))
}

func (h *aiHandler) HandleCreatePrompt(ctx context.Context, req *dto.CreatePromptRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error) {
	return util.WrapHTTPResponse(h.svc.CreatePrompt(ctx, req))
}

// SSE streaming handlers - TODO: Implement SSE response handling
func (h *aiHandler) HandleGenerateContentCompletion(ctx context.Context, req *dto.GenerateContentCompletionRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error) {
	// For now, return not implemented
	// SSE streaming will be implemented separately
	return nil, protocol.ErrNoImplement
}

func (h *aiHandler) HandleGenerateArticleSummary(ctx context.Context, req *dto.GenerateArticleSummaryRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error) {
	return nil, protocol.ErrNoImplement
}

func (h *aiHandler) HandleGenerateArticleQA(ctx context.Context, req *dto.GenerateArticleQARequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error) {
	return nil, protocol.ErrNoImplement
}

func (h *aiHandler) HandleGenerateTermExplaination(ctx context.Context, req *dto.GenerateTermExplainationRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error) {
	return nil, protocol.ErrNoImplement
}
