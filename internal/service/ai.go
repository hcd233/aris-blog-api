package service

import (
	"errors"
	"strings"
	"time"

	chat_model "github.com/hcd233/aris-blog-api/internal/ai/chat_model"
	"github.com/hcd233/aris-blog-api/internal/ai/prompt"
	"github.com/hcd233/aris-blog-api/internal/logger"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/resource/database"
	"github.com/hcd233/aris-blog-api/internal/resource/database/dao"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
	"github.com/hcd233/aris-blog-api/internal/resource/llm"
	"github.com/hcd233/aris-blog-api/internal/util"
	"github.com/samber/lo"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AIService AI服务
//
//	author centonhuang
//	update 2025-01-05 17:57:43
type AIService interface {
	GetPrompt(req *protocol.GetPromptRequest) (rsp *protocol.GetPromptResponse, err error)
	GetLatestPrompt(req *protocol.GetLatestPromptRequest) (rsp *protocol.GetLatestPromptResponse, err error)
	ListPrompt(req *protocol.ListPromptRequest) (rsp *protocol.ListPromptResponse, err error)
	CreatePrompt(req *protocol.CreatePromptRequest) (rsp *protocol.CreatePromptResponse, err error)
	GenerateContentCompletion(req *protocol.GenerateContentCompletionRequest) (rsp *protocol.GenerateContentCompletionResponse, err error)
	GenerateArticleSummary(req *protocol.GenerateArticleSummaryRequest) (rsp *protocol.GenerateArticleSummaryResponse, err error)
	GenerateArticleTranslation(req *protocol.GenerateArticleTranslationRequest) (rsp *protocol.GenerateArticleTranslationResponse, err error)
	GenerateArticleQA(req *protocol.GenerateArticleQARequest) (rsp *protocol.GenerateArticleQAResponse, err error)
	GenerateTermExplaination(req *protocol.GenerateTermExplainationRequest) (rsp *protocol.GenerateTermExplainationResponse, err error)
}

// NewAIService 创建AI服务
//
//	return AIService
//	author centonhuang
//	update 2025-01-05 17:57:43
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

type aiService struct {
	db                *gorm.DB
	userDAO           *dao.UserDAO
	articleDAO        *dao.ArticleDAO
	articleVersionDAO *dao.ArticleVersionDAO
	promptDAO         *dao.PromptDAO
	openAI            *openai.Client
}

// GetPrompt 获取提示词
//
//	receiver s *aiService
//	param req *protocol.GetPromptRequest
//	return rsp *protocol.GetPromptResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 18:02:44
func (s *aiService) GetPrompt(req *protocol.GetPromptRequest) (rsp *protocol.GetPromptResponse, err error) {
	rsp = &protocol.GetPromptResponse{}
	prompt, err := s.promptDAO.GetPromptByTaskAndVersion(s.db, model.Task(req.TaskName), req.Version, []string{"id", "created_at", "task", "version", "templates", "variables"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[AIService] prompt not found", zap.String("taskName", req.TaskName), zap.Uint("version", req.Version))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[AIService] failed to get prompt", zap.String("taskName", req.TaskName), zap.Uint("version", req.Version), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Prompt = &protocol.Prompt{
		ID:        prompt.ID,
		CreatedAt: prompt.CreatedAt.Format(time.DateTime),
		Task:      string(prompt.Task),
		Version:   prompt.Version,
		Templates: lo.Map(prompt.Templates, func(t model.Template, _ int) protocol.Template {
			return protocol.Template{
				Role:    t.Role,
				Content: t.Content,
			}
		}),
		Variables: prompt.Variables,
	}

	return rsp, nil
}

// GetLatestPrompt 获取最新提示词
//
//	receiver s *aiService
//	param req *protocol.GetLatestPromptRequest
//	return rsp *protocol.GetLatestPromptResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 18:02:55
func (s *aiService) GetLatestPrompt(req *protocol.GetLatestPromptRequest) (rsp *protocol.GetLatestPromptResponse, err error) {
	rsp = &protocol.GetLatestPromptResponse{}

	prompt, err := s.promptDAO.GetLatestPromptByTask(s.db, model.Task(req.TaskName), []string{"id", "created_at", "task", "version", "templates", "variables"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[AIService] prompt not found", zap.String("taskName", req.TaskName))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[AIService] failed to get latest prompt", zap.String("taskName", req.TaskName), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Prompt = &protocol.Prompt{
		ID:        prompt.ID,
		CreatedAt: prompt.CreatedAt.Format(time.DateTime),
		Task:      string(prompt.Task),
		Version:   prompt.Version,
		Templates: lo.Map(prompt.Templates, func(t model.Template, _ int) protocol.Template {
			return protocol.Template{
				Role:    t.Role,
				Content: t.Content,
			}
		}),
		Variables: prompt.Variables,
	}

	return rsp, nil
}

// ListPrompt 列出提示词
//
//	receiver s *aiService
//	param req *protocol.ListPromptRequest
//	return rsp *protocol.ListPromptResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 18:02:55
func (s *aiService) ListPrompt(req *protocol.ListPromptRequest) (rsp *protocol.ListPromptResponse, err error) {
	rsp = &protocol.ListPromptResponse{}

	prompts, pageInfo, err := s.promptDAO.PaginateByTask(s.db, model.Task(req.TaskName),
		[]string{"id", "created_at", "task", "version", "templates", "variables"},
		[]string{},
		req.PageParam.Page, req.PageParam.PageSize,
	)
	if err != nil {
		logger.Logger.Error("[AIService] failed to list prompt", zap.String("taskName", req.TaskName), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Prompts = lo.Map(prompts, func(p *model.Prompt, _ int) *protocol.Prompt {
		return &protocol.Prompt{
			ID:        p.ID,
			CreatedAt: p.CreatedAt.Format(time.DateTime),
			Task:      string(p.Task),
			Version:   p.Version,
			Templates: lo.Map(p.Templates, func(t model.Template, _ int) protocol.Template {
				return protocol.Template{
					Role:    t.Role,
					Content: t.Content,
				}
			}),
			Variables: p.Variables,
		}
	})

	rsp.PageInfo = &protocol.PageInfo{
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
		Total:    pageInfo.Total,
	}

	return rsp, nil
}

// CreatePrompt 创建提示词
//
//	receiver s *aiService
//	param req *protocol.CreatePromptRequest
//	return rsp *protocol.CreatePromptResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 18:03:07
func (s *aiService) CreatePrompt(req *protocol.CreatePromptRequest) (rsp *protocol.CreatePromptResponse, err error) {
	contents := lo.Map(req.Templates, func(tmplate protocol.Template, _ int) string {
		return tmplate.Content
	})

	content := strings.Join(contents, "\n")

	variables := util.ExtractVariablesFromContent(content)

	if len(variables) == 0 {
		return nil, protocol.ErrBadRequest
	}

	prompt, err := s.promptDAO.GetLatestPromptByTask(s.db, model.Task(req.TaskName), []string{"id", "templates", "variables", "version"}, []string{})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Logger.Error("[AIService] failed to get latest prompt", zap.String("taskName", req.TaskName), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	contents = lo.Map(prompt.Templates, func(tmplate model.Template, _ int) string {
		return tmplate.Content
	})

	if latestContent := strings.Join(contents, "\n"); latestContent == content {
		logger.Logger.Info("[AIService] the content of the new version is the same as the latest version", zap.String("taskName", req.TaskName), zap.Any("templates", req.Templates))
		return nil, protocol.ErrBadRequest
	}

	if l, r := lo.Difference(prompt.Variables, variables); prompt.ID != 0 && len(l)+len(r) > 0 {
		logger.Logger.Info("[AIService] the variables of the latest prompt and the new prompt are mismatch",
			zap.String("taskName", req.TaskName), zap.Strings("latestVariables", prompt.Variables), zap.Strings("newVariables", variables))
		return nil, protocol.ErrBadRequest
	}

	prompt = &model.Prompt{
		Task: model.Task(req.TaskName),
		Templates: lo.Map(req.Templates, func(tmplate protocol.Template, _ int) model.Template {
			return model.Template{
				Role:    tmplate.Role,
				Content: tmplate.Content,
			}
		}),
		Variables: variables,
		Version:   prompt.Version + 1,
	}

	if err = s.promptDAO.Create(s.db, prompt); err != nil {
		logger.Logger.Error("[AIService] failed to create prompt", zap.String("taskName", req.TaskName), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	return rsp, nil
}

// GenerateContentCompletion 生成内容补全
//
//	receiver s *aiService
//	param req *protocol.GenerateContentCompletionRequest
//	return rsp *protocol.GenerateContentCompletionResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 18:03:15
func (s *aiService) GenerateContentCompletion(req *protocol.GenerateContentCompletionRequest) (rsp *protocol.GenerateContentCompletionResponse, err error) {
	rsp = &protocol.GenerateContentCompletionResponse{}

	user := lo.Must1(s.userDAO.GetByID(s.db, req.CurUserID, []string{"id", "llm_quota"}, []string{}))
	if user.LLMQuota <= 0 {
		logger.Logger.Info("[AIService] insufficient LLM quota", zap.Uint("userID", req.CurUserID), zap.Int("quota", int(user.LLMQuota)))
		return nil, protocol.ErrInsufficientQuota
	}

	latestPrompt, err := s.promptDAO.GetLatestPromptByTask(s.db, model.TaskContentCompletion, []string{"id", "templates"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[AIService] latest prompt not found", zap.String("taskName", string(model.TaskContentCompletion)))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[AIService] failed to get latest prompt", zap.String("taskName", string(model.TaskContentCompletion)), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	oneTurnPrompts := lo.Map(latestPrompt.Templates, func(template model.Template, _ int) prompt.Prompt {
		return prompt.NewOneTurnPrompt(template.Role, template.Content)
	})

	promptTemplate := prompt.NewMultiTurnPrompt(oneTurnPrompts)
	chatOpenAI := chat_model.NewChatOpenAI(chat_model.ZhipuGlm4Flash, req.Temperature)

	params := map[string]interface{}{
		"context":     req.Context,
		"instruction": req.Instruction,
		"reference":   req.Reference,
	}

	tokenChan, errChan, err := chatOpenAI.Stream(lo.Must1(promptTemplate.Format(params)))
	if err != nil {
		logger.Logger.Error("[AIService] failed to stream", zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	lo.Must0(s.userDAO.Update(s.db, user, map[string]interface{}{"llm_quota": user.LLMQuota - 1}))

	rsp.TokenChan = tokenChan
	rsp.ErrChan = errChan

	return rsp, nil
}

// GenerateArticleSummary 生成文章总结
//
//	receiver s *aiService
//	param req *protocol.GenerateArticleSummaryRequest
//	return rsp *protocol.GenerateArticleSummaryResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 18:03:21
func (s *aiService) GenerateArticleSummary(req *protocol.GenerateArticleSummaryRequest) (rsp *protocol.GenerateArticleSummaryResponse, err error) {
	rsp = &protocol.GenerateArticleSummaryResponse{}

	user := lo.Must1(s.userDAO.GetByID(s.db, req.CurUserID, []string{"id", "llm_quota"}, []string{}))
	if user.LLMQuota <= 0 {
		logger.Logger.Info("[AIService] insufficient LLM quota", zap.Uint("userID", req.CurUserID), zap.Int("quota", int(user.LLMQuota)))
		return nil, protocol.ErrInsufficientQuota
	}

	article, err := s.articleDAO.GetBySlugAndUserID(s.db, req.ArticleSlug, req.CurUserID, []string{"id", "title"}, []string{})
	if err != nil {
		logger.Logger.Error("[AIService] failed to get article", zap.String("articleSlug", req.ArticleSlug), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	latestVersion, err := s.articleVersionDAO.GetLatestByArticleID(s.db, article.ID, []string{"id", "content"}, []string{})
	if err != nil {
		logger.Logger.Error("[AIService] failed to get article version", zap.Uint("articleID", article.ID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	latestPrompt, err := s.promptDAO.GetLatestPromptByTask(s.db, model.TaskContentCompletion, []string{"id", "templates"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[AIService] latest prompt not found", zap.String("taskName", string(model.TaskContentCompletion)))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[AIService] failed to get latest prompt", zap.String("taskName", string(model.TaskContentCompletion)), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	oneTurnPrompts := lo.Map(latestPrompt.Templates, func(template model.Template, _ int) prompt.Prompt {
		return prompt.NewOneTurnPrompt(template.Role, template.Content)
	})

	promptTemplate := prompt.NewMultiTurnPrompt(oneTurnPrompts)
	chatOpenAI := chat_model.NewChatOpenAI(chat_model.ZhipuGlm4Flash, req.Temperature)

	params := map[string]interface{}{
		"title":       article.Title,
		"content":     latestVersion.Content,
		"instruction": req.Instruction,
	}

	tokenChan, errChan, err := chatOpenAI.Stream(lo.Must1(promptTemplate.Format(params)))
	if err != nil {
		logger.Logger.Error("[AIService] failed to stream", zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	lo.Must0(s.userDAO.Update(s.db, user, map[string]interface{}{"llm_quota": user.LLMQuota - 1}))

	rsp.TokenChan = tokenChan
	rsp.ErrChan = errChan

	return rsp, nil
}

// GenerateArticleTranslation 生成文章翻译
//
//	receiver s *aiService
//	param req *protocol.GenerateArticleTranslationRequest
//	return rsp *protocol.GenerateArticleTranslationResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 18:03:27
func (s *aiService) GenerateArticleTranslation(_ *protocol.GenerateArticleTranslationRequest) (rsp *protocol.GenerateArticleTranslationResponse, err error) {
	// TODO: 实现
	return nil, protocol.ErrInternalError
}

// GenerateArticleQA 生成文章问答
//
//	receiver s *aiService
//	param req *protocol.GenerateArticleQARequest
//	return rsp *protocol.GenerateArticleQAResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 18:03:44
func (s *aiService) GenerateArticleQA(req *protocol.GenerateArticleQARequest) (rsp *protocol.GenerateArticleQAResponse, err error) {
	rsp = &protocol.GenerateArticleQAResponse{}
	user := lo.Must1(s.userDAO.GetByID(s.db, req.CurUserID, []string{"id", "llm_quota"}, []string{}))
	if user.LLMQuota <= 0 {
		logger.Logger.Info("[AIService] insufficient LLM quota", zap.Uint("userID", req.CurUserID), zap.Int("quota", int(user.LLMQuota)))
		return nil, protocol.ErrInsufficientQuota
	}

	article, err := s.articleDAO.GetBySlugAndUserID(s.db, req.ArticleSlug, req.CurUserID, []string{"id", "title"}, []string{})
	if err != nil {
		logger.Logger.Error("[AIService] failed to get article", zap.String("articleSlug", req.ArticleSlug), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	latestVersion, err := s.articleVersionDAO.GetLatestByArticleID(s.db, article.ID, []string{"id", "content"}, []string{})
	if err != nil {
		logger.Logger.Error("[AIService] failed to get article version", zap.Uint("articleID", article.ID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	latestPrompt, err := s.promptDAO.GetLatestPromptByTask(s.db, model.TaskContentCompletion, []string{"id", "templates"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[AIService] latest prompt not found", zap.String("taskName", string(model.TaskContentCompletion)))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[AIService] failed to get latest prompt", zap.String("taskName", string(model.TaskContentCompletion)), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	oneTurnPrompts := lo.Map(latestPrompt.Templates, func(template model.Template, _ int) prompt.Prompt {
		return prompt.NewOneTurnPrompt(template.Role, template.Content)
	})

	promptTemplate := prompt.NewMultiTurnPrompt(oneTurnPrompts)
	chatOpenAI := chat_model.NewChatOpenAI(chat_model.ZhipuGlm4Flash, req.Temperature)

	params := map[string]interface{}{
		"title":    article.Title,
		"content":  latestVersion.Content,
		"question": req.Question,
	}

	tokenChan, errChan, err := chatOpenAI.Stream(lo.Must1(promptTemplate.Format(params)))
	if err != nil {
		logger.Logger.Error("[AIService] failed to stream", zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	lo.Must0(s.userDAO.Update(s.db, user, map[string]interface{}{"llm_quota": user.LLMQuota - 1}))

	rsp.TokenChan = tokenChan
	rsp.ErrChan = errChan

	return rsp, nil
}

// GenerateTermExplaination 生成术语解释
//
//	receiver s *aiService
//	param req *protocol.GenerateTermExplainationRequest
//	return rsp *protocol.GenerateTermExplainationResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 18:03:48
func (s *aiService) GenerateTermExplaination(req *protocol.GenerateTermExplainationRequest) (rsp *protocol.GenerateTermExplainationResponse, err error) {
	rsp = &protocol.GenerateTermExplainationResponse{}

	user := lo.Must1(s.userDAO.GetByID(s.db, req.CurUserID, []string{"id", "llm_quota"}, []string{}))
	if user.LLMQuota <= 0 {
		logger.Logger.Info("[AIService] insufficient LLM quota", zap.Uint("userID", req.CurUserID), zap.Int("quota", int(user.LLMQuota)))
		return nil, protocol.ErrInsufficientQuota
	}

	article, err := s.articleDAO.GetBySlugAndUserID(s.db, req.ArticleSlug, req.CurUserID, []string{"id", "title"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[AIService] article not found", zap.String("articleSlug", req.ArticleSlug))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[AIService] failed to get article", zap.String("articleSlug", req.ArticleSlug), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	latestVersion, err := s.articleVersionDAO.GetLatestByArticleID(s.db, article.ID, []string{"id", "content"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[AIService] article version not found", zap.Uint("articleID", article.ID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[AIService] failed to get article version", zap.Uint("articleID", article.ID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	latestPrompt, err := s.promptDAO.GetLatestPromptByTask(s.db, model.TaskTermExplaination, []string{"id", "templates"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[AIService] latest prompt not found", zap.String("taskName", string(model.TaskTermExplaination)))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[AIService] failed to get latest prompt", zap.String("taskName", string(model.TaskTermExplaination)), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	oneTurnPrompts := lo.Map(latestPrompt.Templates, func(template model.Template, _ int) prompt.Prompt {
		return prompt.NewOneTurnPrompt(template.Role, template.Content)
	})

	promptTemplate := prompt.NewMultiTurnPrompt(oneTurnPrompts)
	chatOpenAI := chat_model.NewChatOpenAI(chat_model.ZhipuGlm4Flash, req.Temperature)

	contextWindowLen := 200

	var left, right int
	if int(req.Position) < contextWindowLen/2 {
		left = 0
		right = contextWindowLen
	} else if int(req.Position) > len(latestVersion.Content)-contextWindowLen/2 {
		left = len(latestVersion.Content) - contextWindowLen
		right = len(latestVersion.Content)
	} else {
		left = int(req.Position) - contextWindowLen/2
		right = int(req.Position) + contextWindowLen/2
	}

	params := map[string]interface{}{
		"title":   article.Title,
		"content": latestVersion.Content,
		"context": latestVersion.Content[left:right],
		"term":    req.Term,
	}

	tokenChan, errChan, err := chatOpenAI.Stream(lo.Must1(promptTemplate.Format(params)))
	if err != nil {
		logger.Logger.Error("[AIService] failed to stream", zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	lo.Must0(s.userDAO.Update(s.db, user, map[string]interface{}{"llm_quota": user.LLMQuota - 1}))

	rsp.TokenChan = tokenChan
	rsp.ErrChan = errChan

	return rsp, nil
}
