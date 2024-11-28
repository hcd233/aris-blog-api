package service

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"github.com/hcd233/Aris-blog/internal/util"
	"github.com/samber/lo"
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
	promptDAO *dao.PromptDAO
}

func NewAIService() AIService {
	return &aiService{
		db:        database.GetDBInstance(),
		promptDAO: dao.GetPromptDAO(),
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

}

func (s *aiService) GenerateArticleSummaryHandler(c *gin.Context) {

}

func (s *aiService) GenerateArticleTranslationHandler(c *gin.Context) {

}

func (s *aiService) GenerateArticleQAHandler(c *gin.Context) {

}

func (s *aiService) GenerateTermExplainationHandler(c *gin.Context) {

}
