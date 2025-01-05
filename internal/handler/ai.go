// Package handler handler层
//
//	@update 2024-12-08 16:59:38
package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/service"
	"github.com/hcd233/Aris-blog/internal/util"
)

// AIHandler AI服务
//
//	@author centonhuang
//	@update 2024-12-08 16:45:29
type AIHandler interface {
	HandleGetPrompt(c *gin.Context)
	HandleGetLatestPrompt(c *gin.Context)
	HandleListPrompt(c *gin.Context)
	HandleCreatePrompt(c *gin.Context)
	HandleGenerateContentCompletion(c *gin.Context)
	HandleGenerateArticleSummary(c *gin.Context)
	HandleGenerateArticleTranslation(c *gin.Context)
	HandleGenerateArticleQA(c *gin.Context)
	HandleGenerateTermExplaination(c *gin.Context)
}

type aiHandler struct {
	svc service.AIService
}

// NewAIHandler 创建AI服务
//
//	@return AIService
//	@author centonhuang
//	@update 2024-12-08 16:45:37
func NewAIHandler() AIHandler {
	return &aiHandler{
		svc: service.NewAIService(),
	}
}

func (h *aiHandler) HandleGetPrompt(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.PromptVersionURI)

	req := &protocol.GetPromptRequest{
		TaskName: string(uri.TaskName),
		Version:  uri.Version,
	}

	rsp, err := h.svc.GetPrompt(req)

	util.SendHTTPResponse(c, rsp, err)
}

func (h *aiHandler) HandleGetLatestPrompt(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.TaskURI)

	req := &protocol.GetLatestPromptRequest{
		TaskName: string(uri.TaskName),
	}

	rsp, err := h.svc.GetLatestPrompt(req)

	util.SendHTTPResponse(c, rsp, err)
}

func (h *aiHandler) HandleListPrompt(c *gin.Context) {
	param := c.MustGet("param").(*protocol.PageParam)
	uri := c.MustGet("uri").(*protocol.TaskURI)

	req := &protocol.ListPromptRequest{
		TaskName:  string(uri.TaskName),
		PageParam: param,
	}

	rsp, err := h.svc.ListPrompt(req)

	util.SendHTTPResponse(c, rsp, err)
}

func (h *aiHandler) HandleCreatePrompt(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.TaskURI)
	body := c.MustGet("body").(*protocol.CreatePromptBody)

	req := &protocol.CreatePromptRequest{
		TaskName:  string(uri.TaskName),
		Templates: body.Templates,
	}

	rsp, err := h.svc.CreatePrompt(req)

	util.SendHTTPResponse(c, rsp, err)
}

func (h *aiHandler) HandleGenerateContentCompletion(c *gin.Context) {
	userID := c.GetUint("userID")
	body := c.MustGet("body").(*protocol.GenerateContentCompletionBody)

	req := &protocol.GenerateContentCompletionRequest{
		CurUserID:   userID,
		Context:     body.Context,
		Instruction: body.Instruction,
		Reference:   body.Reference,
		Temperature: body.Temperature,
	}

	rsp, err := h.svc.GenerateContentCompletion(req)
	if err != nil {
		util.SendHTTPResponse(c, nil, err)
		return
	}

	util.SendStreamEventResponses(c, rsp.TokenChan, rsp.ErrChan)
}

func (h *aiHandler) HandleGenerateArticleSummary(c *gin.Context) {
	userID := c.GetUint("userID")
	body := c.MustGet("body").(*protocol.GenerateArticleSummaryBody)

	req := &protocol.GenerateArticleSummaryRequest{
		CurUserID:   userID,
		ArticleSlug: body.ArticleSlug,
		Instruction: body.Instruction,
		Temperature: body.Temperature,
	}

	rsp, err := h.svc.GenerateArticleSummary(req)
	if err != nil {
		util.SendHTTPResponse(c, nil, err)
		return
	}

	util.SendStreamEventResponses(c, rsp.TokenChan, rsp.ErrChan)
}

func (h *aiHandler) HandleGenerateArticleTranslation(c *gin.Context) {
	// TODO: 实现
	req := &protocol.GenerateArticleTranslationRequest{}

	_, err := h.svc.GenerateArticleTranslation(req)
	if err != nil {
		util.SendHTTPResponse(c, nil, err)
		return
	}
}

func (h *aiHandler) HandleGenerateArticleQA(c *gin.Context) {
	userID := c.GetUint("userID")
	body := c.MustGet("body").(*protocol.GenerateArticleQABody)

	req := &protocol.GenerateArticleQARequest{
		CurUserID:   userID,
		ArticleSlug: body.ArticleSlug,
		Question:    body.Question,
		Temperature: body.Temperature,
	}

	rsp, err := h.svc.GenerateArticleQA(req)
	if err != nil {
		util.SendHTTPResponse(c, nil, err)
		return
	}

	util.SendStreamEventResponses(c, rsp.TokenChan, rsp.ErrChan)
}

func (h *aiHandler) HandleGenerateTermExplaination(c *gin.Context) {
	userID := c.GetUint("userID")
	body := c.MustGet("body").(*protocol.GenerateTermExplainationBody)

	req := &protocol.GenerateTermExplainationRequest{
		CurUserID:   userID,
		ArticleSlug: body.ArticleSlug,
		Term:        body.Term,
		Position:    body.Position,
		Temperature: body.Temperature,
	}

	rsp, err := h.svc.GenerateTermExplaination(req)
	if err != nil {
		util.SendHTTPResponse(c, nil, err)
		return
	}

	util.SendStreamEventResponses(c, rsp.TokenChan, rsp.ErrChan)
}
