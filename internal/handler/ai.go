// Package handler handler层
//
//	update 2024-12-08 16:59:38
package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// AIHandler AI服务
//
//	author centonhuang
//	update 2024-12-08 16:45:29
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
//	return AIService
//	author centonhuang
//	update 2024-12-08 16:45:37
func NewAIHandler() AIHandler {
	return &aiHandler{
		svc: service.NewAIService(),
	}
}

// HandleGetPrompt 获取Prompt
//
//	@Summary		获取Prompt
//	@Description	获取Prompt
//	@Tags			ai
//	@Accept			json
//	@Produce		json
//	@Param			uri		path		protocol.PromptVersionURI	true	"Prompt版本URI"
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.HTTPResponse{data=protocol.DeleteUserViewResponse,error=nil}
//	@Failure		400			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/ai/prompt/{taskName}/{version} [get]
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:46:35
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

// HandleListPrompt 获取Prompt列表
//
//	@Summary		获取Prompt列表
//	@Description	获取Prompt列表
//	@Tags			ai
//	@Accept			json
//	@Produce		json
//	@Param			uri		path		protocol.TaskURI	true	"任务URI"
//	@Param			param	query		protocol.PageParam	true	"分页参数"
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.HTTPResponse{data=protocol.ListPromptResponse,error=nil}
//	@Failure		400			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/ai/prompt [get]
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:46:35
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

// HandleCreatePrompt 创建Prompt
//
//	@Summary		创建Prompt
//	@Description	创建Prompt
//	@Tags			ai
//	@Accept			json
//	@Produce		json
//	@Param			uri		path		protocol.TaskURI	true	"任务URI"
//	@Param			body	body		protocol.CreatePromptBody	true	"创建Prompt请求体"
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.HTTPResponse{data=protocol.CreatePromptResponse,error=nil}
//	@Failure		400			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/ai/prompt [post]
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:46:35
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

// HandleGenerateContentCompletion 生成内容补全
//
//	@Summary		生成内容补全
//	@Description	生成内容补全
//	@Tags			ai
//	@Accept			json
//	@Produce		text/event-stream
//	@Param			body	body		protocol.GenerateContentCompletionBody	true	"生成内容补全请求体"
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.SSEResponse{}
//	@Failure		400			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/ai/content-completion [post]
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:46:35
func (h *aiHandler) HandleGenerateContentCompletion(c *gin.Context) {
	userID := c.GetUint("userID")
	body := c.MustGet("body").(*protocol.GenerateContentCompletionBody)

	req := &protocol.GenerateContentCompletionRequest{
		UserID:      userID,
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

// HandleGenerateArticleSummary 生成文章摘要
//
//	@Summary		生成文章摘要
//	@Description	生成文章摘要
//	@Tags			ai
//	@Accept			json
//	@Produce		text/event-stream
//	@Param			body	body		protocol.GenerateArticleSummaryBody	true	"生成文章摘要请求体"
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.SSEResponse{}
//	@Failure		400			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/ai/article-summary [post]
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:46:35
func (h *aiHandler) HandleGenerateArticleSummary(c *gin.Context) {
	userID := c.GetUint("userID")
	body := c.MustGet("body").(*protocol.GenerateArticleSummaryBody)

	req := &protocol.GenerateArticleSummaryRequest{
		UserID:      userID,
		ArticleID:   body.ArticleID,
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

// HandleGenerateArticleTranslation 生成文章翻译
//
//	@Summary		生成文章翻译
//	@Description	生成文章翻译
//	@Tags			ai
//	@Accept			json
//	@Produce		text/event-stream
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.SSEResponse{}
//	@Failure		400			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/ai/article-translation [post]
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:46:35
func (h *aiHandler) HandleGenerateArticleTranslation(c *gin.Context) {
	// TODO: 实现
	req := &protocol.GenerateArticleTranslationRequest{}

	rsp, err := h.svc.GenerateArticleTranslation(req)
	if err != nil {
		util.SendHTTPResponse(c, nil, err)
		return
	}

	util.SendStreamEventResponses(c, rsp.TokenChan, rsp.ErrChan)
}

// HandleGenerateArticleQA 生成文章问答
//
//	@Summary		生成文章问答
//	@Description	生成文章问答
//	@Tags			ai
//	@Accept			json
//	@Produce		text/event-stream
//	@Param			body	body		protocol.GenerateArticleQABody	true	"生成文章问答请求体"
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.SSEResponse{}
//	@Failure		400			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/ai/article-qa [post]
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:46:35
func (h *aiHandler) HandleGenerateArticleQA(c *gin.Context) {
	userID := c.GetUint("userID")
	body := c.MustGet("body").(*protocol.GenerateArticleQABody)

	req := &protocol.GenerateArticleQARequest{
		UserID:      userID,
		ArticleID:   body.ArticleID,
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

// HandleGenerateTermExplaination 生成术语解释
//
//	@Summary		生成术语解释
//	@Description	生成术语解释
//	@Tags			ai
//	@Accept			json
//	@Produce		text/event-stream
//	@Param			body	body		protocol.GenerateTermExplainationBody	true	"生成术语解释请求体"
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.SSEResponse{}
//	@Failure		400			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/ai/term-explaination [post]
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:46:35
func (h *aiHandler) HandleGenerateTermExplaination(c *gin.Context) {
	userID := c.GetUint("userID")
	body := c.MustGet("body").(*protocol.GenerateTermExplainationBody)

	req := &protocol.GenerateTermExplainationRequest{
		UserID:      userID,
		ArticleID:   body.ArticleID,
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
