// Package handler handler层
//
//	@update 2024-12-08 16:59:38
package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	chat_model "github.com/hcd233/Aris-blog/internal/ai/chat_model"
	prompt "github.com/hcd233/Aris-blog/internal/ai/prompt"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"github.com/hcd233/Aris-blog/internal/resource/llm"
	"github.com/hcd233/Aris-blog/internal/util"
	"github.com/samber/lo"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
)

// AIService AI服务
//
//	@author centonhuang
//	@update 2024-12-08 16:45:29
type AIService interface {
	GetPromptHandler(c *gin.Context)
	GetLatestPromptHandler(c *gin.Context)
	ListPromptHandler(c *gin.Context)
	CreatePromptHandler(c *gin.Context)
	GenerateContentCompletionHandler(c *gin.Context)
	GenerateArticleSummaryHandler(c *gin.Context)
	GenerateArticleTranslationHandler(c *gin.Context)
	GenerateArticleQAHandler(c *gin.Context)
	GenerateTermExplainationHandler(c *gin.Context)
}

type aiService struct {
	db                *gorm.DB
	userDAO           *dao.UserDAO
	articleDAO        *dao.ArticleDAO
	articleVersionDAO *dao.ArticleVersionDAO
	promptDAO         *dao.PromptDAO
	openAI            *openai.Client
}

// NewAIService 创建AI服务
//
//	@return AIService
//	@author centonhuang
//	@update 2024-12-08 16:45:37
func NewAIService() AIService {
	return &aiService{
		db:                database.GetDBInstance(),
		userDAO:           dao.GetUserDAO(),
		articleDAO:        dao.GetArticleDAO(),
		articleVersionDAO: dao.GetArticleVersionDAO(),
		promptDAO:         dao.GetPromptDAO(),
		openAI:            llm.GetOpenAIClient(),
	}
}

func (s *aiService) GetPromptHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.PromptVersionURI)

	prompt, err := s.promptDAO.GetPromptByTaskAndVersion(s.db, model.Task(uri.TaskName), uri.Version, []string{"id", "created_at", "task", "version", "templates", "variables"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetPromptError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Data: map[string]interface{}{
			"prompt": prompt.GetDetailedInfo(),
		},
	})
}

func (s *aiService) GetLatestPromptHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.TaskURI)

	prompt, err := s.promptDAO.GetLatestPromptByTask(s.db, model.Task(uri.TaskName), []string{"id", "created_at", "task", "version", "templates", "variables"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetPromptError,
			Message: err.Error(),
		})
	}

	c.JSON(http.StatusOK, protocol.Response{
		Data: map[string]interface{}{
			"prompt": prompt.GetDetailedInfo(),
		},
	})
}

func (s *aiService) ListPromptHandler(c *gin.Context) {
	param := c.MustGet("param").(*protocol.PageParam)
	uri := c.MustGet("uri").(*protocol.TaskURI)

	prompts, pageInfo, err := s.promptDAO.PaginateByTask(s.db, model.Task(uri.TaskName), []string{"id", "created_at", "task", "version"}, []string{}, param.Page, param.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetPromptError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Data: map[string]interface{}{
			"prompts": lo.Map(prompts, func(prompt *model.Prompt, _ int) map[string]interface{} {
				return prompt.GetBasicInfo()
			}),
			"pageInfo": pageInfo,
		},
	})
}

func (s *aiService) CreatePromptHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.TaskURI)
	body := c.MustGet("body").(*protocol.CreatePromptBody)

	contents := lo.Map(body.Templates, func(tmplate protocol.Template, idx int) string {
		return tmplate.Content
	})

	content := strings.Join(contents, "\n")

	variables := util.ExtractVariablesFromContent(content)

	if len(variables) == 0 {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeCreatePromptError,
			Message: "提示词中未找到变量",
		})
		return
	}

	prompt, err := s.promptDAO.GetLatestPromptByTask(s.db, model.Task(uri.TaskName), []string{"id", "templates", "variables", "version"}, []string{})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetPromptError,
			Message: err.Error(),
		})
	}

	contents = lo.Map(prompt.Templates, func(tmplate model.Template, idx int) string {
		return tmplate.Content
	})

	if latestContent := strings.Join(contents, "\n"); latestContent == content {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeCreatePromptError,
			Message: "The content of the new version is the same as the latest version",
		})
		return
	}

	if l, r := lo.Difference(prompt.Variables, variables); prompt.ID != 0 && len(l)+len(r) > 0 {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeCreatePromptError,
			Message: fmt.Sprintf("the variables of the latest prompt and the new prompt are mismatch, latest: %v, new: %v", prompt.Variables, variables),
		})
		return
	}

	prompt = &model.Prompt{
		Task: model.Task(uri.TaskName),
		Templates: lo.Map(body.Templates, func(tmplate protocol.Template, idx int) model.Template {
			return model.Template{
				Role:    tmplate.Role,
				Content: tmplate.Content,
			}
		}),
		Variables: variables,
		Version:   prompt.Version + 1,
	}

	if err = s.promptDAO.Create(s.db, prompt); err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeCreatePromptError,
			Message: err.Error(),
		})
		return
	}

	prompt = lo.Must1(s.promptDAO.GetLatestPromptByTask(s.db, model.Task(uri.TaskName), []string{"id", "created_at", "task", "version"}, []string{}))

	c.JSON(http.StatusOK, protocol.Response{
		Data: map[string]interface{}{
			"prompt": prompt.GetBasicInfo(),
		},
	})
}

func (s *aiService) GenerateContentCompletionHandler(c *gin.Context) {
	userID := c.GetUint("userID")
	body := c.MustGet("body").(*protocol.GenerateContentCompletionBody)

	user := lo.Must1(s.userDAO.GetByID(s.db, userID, []string{"id", "llm_quota"}, []string{}))
	if user.LLMQuota <= 0 {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeInsufficientQuota,
			Message: fmt.Sprintf("Insufficient LLM quota: %d", user.LLMQuota),
		})
		return
	}

	latestPrompt, err := s.promptDAO.GetLatestPromptByTask(s.db, model.TaskContentCompletion, []string{"id", "templates"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetPromptError,
			Message: err.Error(),
		})
		return
	}

	oneTurnPrompts := lo.Map(latestPrompt.Templates, func(template model.Template, idx int) prompt.Prompt {
		return prompt.NewOneTurnPrompt(template.Role, template.Content)
	})

	promptTemplate := prompt.NewMultiTurnPrompt(oneTurnPrompts)
	chatOpenAI := chat_model.NewChatOpenAI(chat_model.ZhipuGlm4Flash, body.Temperature)

	params := map[string]interface{}{
		"context":     body.Context,
		"instruction": body.Instruction,
		"reference":   body.Reference,
	}

	tokenChan, errChan, err := chatOpenAI.Stream(lo.Must1(promptTemplate.Format(params)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGenerateContentCompletionError,
			Message: err.Error(),
		})
		return
	}

	err = util.SendStreamEventResponses(c, tokenChan, errChan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGenerateContentCompletionError,
			Message: err.Error(),
		})
		return
	}

	lo.Must0(s.userDAO.Update(s.db, user, map[string]interface{}{"llm_quota": user.LLMQuota - 1}))
}

func (s *aiService) GenerateArticleSummaryHandler(c *gin.Context) {
	userID := c.GetUint("userID")
	body := c.MustGet("body").(*protocol.GenerateArticleSummaryBody)

	user := lo.Must1(s.userDAO.GetByID(s.db, userID, []string{"id", "llm_quota"}, []string{}))
	if user.LLMQuota <= 0 {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeInsufficientQuota,
			Message: fmt.Sprintf("Insufficient LLM quota: %d", user.LLMQuota),
		})
		return
	}

	article, err := s.articleDAO.GetBySlugAndUserID(s.db, body.ArticleSlug, userID, []string{"id", "title"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	latestVersion, err := s.articleVersionDAO.GetLatestByArticleID(s.db, article.ID, []string{"id", "content"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetArticleVersionError,
			Message: err.Error(),
		})
		return
	}

	latestPrompt, err := s.promptDAO.GetLatestPromptByTask(s.db, model.TaskArticleSummary, []string{"id", "templates"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetPromptError,
			Message: err.Error(),
		})
		return
	}

	oneTurnPrompts := lo.Map(latestPrompt.Templates, func(template model.Template, idx int) prompt.Prompt {
		return prompt.NewOneTurnPrompt(template.Role, template.Content)
	})

	promptTemplate := prompt.NewMultiTurnPrompt(oneTurnPrompts)
	chatOpenAI := chat_model.NewChatOpenAI(chat_model.ZhipuGlm4Flash, body.Temperature)

	params := map[string]interface{}{
		"title":       article.Title,
		"content":     latestVersion.Content,
		"instruction": body.Instruction,
	}

	tokenChan, errChan, err := chatOpenAI.Stream(lo.Must1(promptTemplate.Format(params)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGenerateContentCompletionError,
			Message: err.Error(),
		})
		return
	}

	err = util.SendStreamEventResponses(c, tokenChan, errChan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGenerateContentCompletionError,
			Message: err.Error(),
		})
		return
	}

	lo.Must0(s.userDAO.Update(s.db, user, map[string]interface{}{"llm_quota": user.LLMQuota - 1}))
}

func (s *aiService) GenerateArticleTranslationHandler(c *gin.Context) {
}

func (s *aiService) GenerateArticleQAHandler(c *gin.Context) {
	userID := c.GetUint("userID")
	body := c.MustGet("body").(*protocol.GenerateArticleQABody)

	user := lo.Must1(s.userDAO.GetByID(s.db, userID, []string{"id", "llm_quota"}, []string{}))
	if user.LLMQuota <= 0 {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeInsufficientQuota,
			Message: fmt.Sprintf("Insufficient LLM quota: %d", user.LLMQuota),
		})
		return
	}

	article, err := s.articleDAO.GetBySlugAndUserID(s.db, body.ArticleSlug, userID, []string{"id", "title"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	latestVersion, err := s.articleVersionDAO.GetLatestByArticleID(s.db, article.ID, []string{"id", "content"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetArticleVersionError,
			Message: err.Error(),
		})
		return
	}

	latestPrompt, err := s.promptDAO.GetLatestPromptByTask(s.db, model.TaskArticleQA, []string{"id", "templates"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetPromptError,
			Message: err.Error(),
		})
		return
	}

	oneTurnPrompts := lo.Map(latestPrompt.Templates, func(template model.Template, idx int) prompt.Prompt {
		return prompt.NewOneTurnPrompt(template.Role, template.Content)
	})

	promptTemplate := prompt.NewMultiTurnPrompt(oneTurnPrompts)
	chatOpenAI := chat_model.NewChatOpenAI(chat_model.ZhipuGlm4Flash, body.Temperature)

	params := map[string]interface{}{
		"title":    article.Title,
		"content":  latestVersion.Content,
		"question": body.Question,
	}

	tokenChan, errChan, err := chatOpenAI.Stream(lo.Must1(promptTemplate.Format(params)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGenerateContentCompletionError,
			Message: err.Error(),
		})
		return
	}

	err = util.SendStreamEventResponses(c, tokenChan, errChan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGenerateContentCompletionError,
			Message: err.Error(),
		})
		return
	}

	lo.Must0(s.userDAO.Update(s.db, user, map[string]interface{}{"llm_quota": user.LLMQuota - 1}))
}

func (s *aiService) GenerateTermExplainationHandler(c *gin.Context) {
	userID := c.GetUint("userID")
	body := c.MustGet("body").(*protocol.GenerateTermExplainationBody)

	user := lo.Must1(s.userDAO.GetByID(s.db, userID, []string{"id", "llm_quota"}, []string{}))
	if user.LLMQuota <= 0 {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeInsufficientQuota,
			Message: fmt.Sprintf("Insufficient LLM quota: %d", user.LLMQuota),
		})
		return
	}

	article, err := s.articleDAO.GetBySlugAndUserID(s.db, body.ArticleSlug, userID, []string{"id", "title"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	latestVersion, err := s.articleVersionDAO.GetLatestByArticleID(s.db, article.ID, []string{"id", "content"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetArticleVersionError,
			Message: err.Error(),
		})
		return
	}

	latestPrompt, err := s.promptDAO.GetLatestPromptByTask(s.db, model.TaskTermExplaination, []string{"id", "templates"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetPromptError,
			Message: err.Error(),
		})
		return
	}

	oneTurnPrompts := lo.Map(latestPrompt.Templates, func(template model.Template, idx int) prompt.Prompt {
		return prompt.NewOneTurnPrompt(template.Role, template.Content)
	})

	promptTemplate := prompt.NewMultiTurnPrompt(oneTurnPrompts)
	chatOpenAI := chat_model.NewChatOpenAI(chat_model.ZhipuGlm4Flash, body.Temperature)

	contextWindowLen := 200

	var left, right int
	if int(body.Position) < contextWindowLen/2 {
		left = 0
		right = contextWindowLen
	} else if int(body.Position) > len(latestVersion.Content)-contextWindowLen/2 {
		left = len(latestVersion.Content) - contextWindowLen
		right = len(latestVersion.Content)
	} else {
		left = int(body.Position) - contextWindowLen/2
		right = int(body.Position) + contextWindowLen/2
	}

	params := map[string]interface{}{
		"title":   article.Title,
		"content": latestVersion.Content,
		"context": latestVersion.Content[left:right],
		"term":    body.Term,
	}

	tokenChan, errChan, err := chatOpenAI.Stream(lo.Must1(promptTemplate.Format(params)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGenerateContentCompletionError,
			Message: err.Error(),
		})
		return
	}

	err = util.SendStreamEventResponses(c, tokenChan, errChan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGenerateContentCompletionError,
			Message: err.Error(),
		})
		return
	}

	lo.Must0(s.userDAO.Update(s.db, user, map[string]interface{}{"llm_quota": user.LLMQuota - 1}))
}
