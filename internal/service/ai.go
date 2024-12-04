package service

import (
	"encoding/json"
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
	db        *gorm.DB
	userDAO   *dao.UserDAO
	promptDAO *dao.PromptDAO
	openAI    *openai.Client
}

func NewAIService() AIService {
	return &aiService{
		db:        database.GetDBInstance(),
		userDAO:   dao.GetUserDAO(),
		promptDAO: dao.GetPromptDAO(),
		openAI:    llm.GetOpenAIClient(),
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
	}

	latestPrompt, err := s.promptDAO.GetLatestPromptByTask(s.db, model.TaskContentCompletion, []string{"id", "templates"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetPromptError,
			Message: err.Error(),
		})
	}

	oneTurnPrompts := lo.Map(latestPrompt.Templates, func(template model.Template, idx int) prompt.Prompt {
		return prompt.NewOneTurnPrompt(template.Role, template.Content)
	})

	promptTemplate := prompt.NewMultiTurnPrompt(oneTurnPrompts)
	chatOpenAI := chat_model.NewChatOpenAI(chat_model.OpenAIGPT4oMini, body.Temperature)

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
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	for {
		select {
		case token, ok := <-tokenChan:
			if !ok {
				c.SSEvent("done", string(lo.Must1(json.Marshal(protocol.AIStreamResponse{
					Delta: "",
					Stop:  true,
					Error: "",
				}))))
				c.Writer.Flush()

				lo.Must0(s.userDAO.Update(s.db, user, map[string]interface{}{"llm_quota": user.LLMQuota - 1}))
				return
			}
			c.SSEvent("stream", string(lo.Must1(json.Marshal(protocol.AIStreamResponse{
				Delta: token,
				Stop:  false,
				Error: "",
			}))))
			c.Writer.Flush()
		case err := <-errChan:
			if err != nil {
				c.SSEvent("error", string(lo.Must1(json.Marshal(protocol.AIStreamResponse{
					Delta: "",
					Stop:  true,
					Error: err.Error(),
				}))))
				c.Writer.Flush()
				return
			}
		}
	}
}

func (s *aiService) GenerateArticleSummaryHandler(c *gin.Context) {

}

func (s *aiService) GenerateArticleTranslationHandler(c *gin.Context) {

}

func (s *aiService) GenerateArticleQAHandler(c *gin.Context) {

}

func (s *aiService) GenerateTermExplainationHandler(c *gin.Context) {

}
