// Package handler provides handlers for the API.
//
//	author centonhuang
//	update 2025-11-02 04:14:56
package handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2/sse"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	dto "github.com/hcd233/aris-blog-api/internal/protocol/dto"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// AIHandler AI处理器
type AIHandler interface {
	HandleGetPrompt(ctx context.Context, req *dto.GetPromptRequest) (*protocol.HTTPResponse[*dto.GetPromptResponse], error)
	HandleGetLatestPrompt(ctx context.Context, req *dto.GetLatestPromptRequest) (*protocol.HTTPResponse[*dto.GetLatestPromptResponse], error)
	HandleListPrompt(ctx context.Context, req *dto.ListPromptRequest) (*protocol.HTTPResponse[*dto.ListPromptResponse], error)
	HandleCreatePrompt(ctx context.Context, req *dto.CreatePromptRequest) (*protocol.HTTPResponse[*dto.EmptyResponse], error)
	// SSE streaming methods - will return special responses
	HandleGenerateContentCompletion(ctx context.Context, req *dto.GenerateContentCompletionRequest, sender sse.Sender)
	HandleGenerateArticleSummary(ctx context.Context, req *dto.GenerateArticleSummaryRequest, sender sse.Sender)
	HandleGenerateArticleQA(ctx context.Context, req *dto.GenerateArticleQARequest, sender sse.Sender)
	HandleGenerateTermExplaination(ctx context.Context, req *dto.GenerateTermExplainationRequest, sender sse.Sender)
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

func (h *aiHandler) HandleGetPrompt(ctx context.Context, req *dto.GetPromptRequest) (*protocol.HTTPResponse[*dto.GetPromptResponse], error) {
	return util.WrapHTTPResponse(h.svc.GetPrompt(ctx, req))
}

func (h *aiHandler) HandleGetLatestPrompt(ctx context.Context, req *dto.GetLatestPromptRequest) (*protocol.HTTPResponse[*dto.GetLatestPromptResponse], error) {
	return util.WrapHTTPResponse(h.svc.GetLatestPrompt(ctx, req))
}

func (h *aiHandler) HandleListPrompt(ctx context.Context, req *dto.ListPromptRequest) (*protocol.HTTPResponse[*dto.ListPromptResponse], error) {
	return util.WrapHTTPResponse(h.svc.ListPrompt(ctx, req))
}

func (h *aiHandler) HandleCreatePrompt(ctx context.Context, req *dto.CreatePromptRequest) (*protocol.HTTPResponse[*dto.EmptyResponse], error) {
	return util.WrapHTTPResponse(h.svc.CreatePrompt(ctx, req))
}

// SSE streaming handlers - TODO: Implement SSE response handling
func (h *aiHandler) HandleGenerateContentCompletion(ctx context.Context, req *dto.GenerateContentCompletionRequest, sender sse.Sender) {
	tokenChan, errChan := h.svc.GenerateContentCompletion(ctx, req)
	util.SendStreamEventResponses(sender, tokenChan, errChan)
}

func (h *aiHandler) HandleGenerateArticleSummary(ctx context.Context, req *dto.GenerateArticleSummaryRequest, sender sse.Sender) {
	tokenChan, errChan := h.svc.GenerateArticleSummary(ctx, req)
	util.SendStreamEventResponses(sender, tokenChan, errChan)
}

func (h *aiHandler) HandleGenerateArticleQA(ctx context.Context, req *dto.GenerateArticleQARequest, sender sse.Sender) {
	tokenChan, errChan := h.svc.GenerateArticleQA(ctx, req)
	util.SendStreamEventResponses(sender, tokenChan, errChan)
}

func (h *aiHandler) HandleGenerateTermExplaination(ctx context.Context, req *dto.GenerateTermExplainationRequest, sender sse.Sender) {
	tokenChan, errChan := h.svc.GenerateTermExplaination(ctx, req)
	util.SendStreamEventResponses(sender, tokenChan, errChan)
}
