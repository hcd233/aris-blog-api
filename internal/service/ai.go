package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/cloudwego/eino-ext/callbacks/langfuse"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/hcd233/aris-blog-api/internal/ai/callback"
	"github.com/hcd233/aris-blog-api/internal/config"
	"github.com/hcd233/aris-blog-api/internal/logger"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/resource/database"
	"github.com/hcd233/aris-blog-api/internal/resource/database/dao"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
	"github.com/hcd233/aris-blog-api/internal/util"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AIService AI服务
//
//	author centonhuang
//	update 2025-01-05 17:57:43
type AIService interface {
	GetPrompt(ctx context.Context, req *protocol.GetPromptRequest) (rsp *protocol.GetPromptResponse, err error)
	GetLatestPrompt(ctx context.Context, req *protocol.GetLatestPromptRequest) (rsp *protocol.GetLatestPromptResponse, err error)
	ListPrompt(ctx context.Context, req *protocol.ListPromptRequest) (rsp *protocol.ListPromptResponse, err error)
	CreatePrompt(ctx context.Context, req *protocol.CreatePromptRequest) (rsp *protocol.CreatePromptResponse, err error)
	GenerateContentCompletion(ctx context.Context, req *protocol.GenerateContentCompletionRequest) (rsp *protocol.GenerateContentCompletionResponse, err error)
	GenerateArticleSummary(ctx context.Context, req *protocol.GenerateArticleSummaryRequest) (rsp *protocol.GenerateArticleSummaryResponse, err error)
	GenerateArticleTranslation(ctx context.Context, req *protocol.GenerateArticleTranslationRequest) (rsp *protocol.GenerateArticleTranslationResponse, err error)
	GenerateArticleQA(ctx context.Context, req *protocol.GenerateArticleQARequest) (rsp *protocol.GenerateArticleQAResponse, err error)
	GenerateTermExplaination(ctx context.Context, req *protocol.GenerateTermExplainationRequest) (rsp *protocol.GenerateTermExplainationResponse, err error)
}

// NewAIService 创建AI服务
//
//	return AIService
//	author centonhuang
//	update 2025-01-05 17:57:43
func NewAIService() AIService {
	return &aiService{
		userDAO:           dao.GetUserDAO(),
		articleDAO:        dao.GetArticleDAO(),
		articleVersionDAO: dao.GetArticleVersionDAO(),
		promptDAO:         dao.GetPromptDAO(),
	}
}

type aiService struct {
	userDAO           *dao.UserDAO
	articleDAO        *dao.ArticleDAO
	articleVersionDAO *dao.ArticleVersionDAO
	promptDAO         *dao.PromptDAO
}

// GetPrompt 获取提示词
//
//	receiver s *aiService
//	param req *protocol.GetPromptRequest
//	return rsp *protocol.GetPromptResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 18:02:44
func (s *aiService) GetPrompt(ctx context.Context, req *protocol.GetPromptRequest) (rsp *protocol.GetPromptResponse, err error) {
	rsp = &protocol.GetPromptResponse{}

	logger := logger.LoggerWithContext(ctx)
	db := database.GetDBInstance(ctx)

	prompt, err := s.promptDAO.GetPromptByTaskAndVersion(db, model.Task(req.TaskName), req.Version, []string{"id", "created_at", "task", "version", "templates", "variables"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[AIService] prompt not found", zap.String("taskName", req.TaskName), zap.Uint("version", req.Version))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[AIService] failed to get prompt", zap.String("taskName", req.TaskName), zap.Uint("version", req.Version), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Prompt = &protocol.Prompt{
		PromptID:  prompt.ID,
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
func (s *aiService) GetLatestPrompt(ctx context.Context, req *protocol.GetLatestPromptRequest) (rsp *protocol.GetLatestPromptResponse, err error) {
	rsp = &protocol.GetLatestPromptResponse{}

	logger := logger.LoggerWithContext(ctx)
	db := database.GetDBInstance(ctx)

	prompt, err := s.promptDAO.GetLatestPromptByTask(db, model.Task(req.TaskName), []string{"id", "created_at", "task", "version", "templates", "variables"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[AIService] prompt not found", zap.String("taskName", req.TaskName))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[AIService] failed to get latest prompt", zap.String("taskName", req.TaskName), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Prompt = &protocol.Prompt{
		PromptID:  prompt.ID,
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
func (s *aiService) ListPrompt(ctx context.Context, req *protocol.ListPromptRequest) (rsp *protocol.ListPromptResponse, err error) {
	rsp = &protocol.ListPromptResponse{}

	logger := logger.LoggerWithContext(ctx)
	db := database.GetDBInstance(ctx)

	prompts, pageInfo, err := s.promptDAO.PaginateByTask(db, model.Task(req.TaskName),
		[]string{"id", "created_at", "task", "version", "templates", "variables"},
		[]string{},
		req.PageParam.Page, req.PageParam.PageSize,
	)
	if err != nil {
		logger.Error("[AIService] failed to list prompt", zap.String("taskName", req.TaskName), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Prompts = lo.Map(prompts, func(p *model.Prompt, _ int) *protocol.Prompt {
		return &protocol.Prompt{
			PromptID:  p.ID,
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
func (s *aiService) CreatePrompt(ctx context.Context, req *protocol.CreatePromptRequest) (rsp *protocol.CreatePromptResponse, err error) {
	rsp = &protocol.CreatePromptResponse{}

	logger := logger.LoggerWithContext(ctx)
	db := database.GetDBInstance(ctx)

	contents := lo.Map(req.Templates, func(tmplate protocol.Template, _ int) string {
		return tmplate.Content
	})

	content := strings.Join(contents, "\n")

	variables := util.ExtractVariablesFromContent(content)

	if len(variables) == 0 {
		logger.Error("[AIService] no variables found in the content", zap.String("taskName", req.TaskName), zap.String("content", content))
		return nil, protocol.ErrBadRequest
	}

	prompt, err := s.promptDAO.GetLatestPromptByTask(db, model.Task(req.TaskName), []string{"id", "templates", "variables", "version"}, []string{})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("[AIService] failed to get latest prompt", zap.String("taskName", req.TaskName), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	contents = lo.Map(prompt.Templates, func(tmplate model.Template, _ int) string {
		return tmplate.Content
	})

	if latestContent := strings.Join(contents, "\n"); latestContent == content {
		logger.Info("[AIService] the content of the new version is the same as the latest version", zap.String("taskName", req.TaskName), zap.Any("templates", req.Templates))
		return nil, protocol.ErrBadRequest
	}

	if l, r := lo.Difference(prompt.Variables, variables); prompt.ID != 0 && len(l)+len(r) > 0 {
		logger.Info("[AIService] the variables of the latest prompt and the new prompt are mismatch",
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

	if err = s.promptDAO.Create(db, prompt); err != nil {
		logger.Error("[AIService] failed to create prompt", zap.String("taskName", req.TaskName), zap.Error(err))
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
func (s *aiService) GenerateContentCompletion(ctx context.Context, req *protocol.GenerateContentCompletionRequest) (rsp *protocol.GenerateContentCompletionResponse, err error) {
	rsp = &protocol.GenerateContentCompletionResponse{}

	logger := logger.LoggerWithContext(ctx)
	db := database.GetDBInstance(ctx)

	user := lo.Must1(s.userDAO.GetByID(db, req.UserID, []string{"id", "name", "llm_quota"}, []string{}))
	if user.LLMQuota <= 0 {
		logger.Info("[AIService] insufficient LLM quota", zap.Uint("userID", req.UserID), zap.Int("quota", int(user.LLMQuota)))
		return nil, protocol.ErrInsufficientQuota
	}

	latestPrompt, err := s.promptDAO.GetLatestPromptByTask(db, model.TaskContentCompletion, []string{"id", "task", "templates"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[AIService] latest prompt not found", zap.String("taskName", string(latestPrompt.Task)))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[AIService] failed to get latest prompt", zap.String("taskName", string(latestPrompt.Task)), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	messages := lo.Map(latestPrompt.Templates, func(template model.Template, _ int) schema.MessagesTemplate {
		return &schema.Message{
			Name:    string(latestPrompt.Task),
			Role:    schema.RoleType(template.Role),
			Content: template.Content,
		}
	})

	promptTemplate := prompt.FromMessages(schema.GoTemplate, messages...)

	chatOpenAI, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		Model:       config.OpenAIModel,
		APIKey:      config.OpenAIAPIKey,
		BaseURL:     config.OpenAIBaseURL,
		Temperature: &req.Temperature,
	})
	if err != nil {
		logger.Error("[AIService] failed to create chat openai", zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	chain := compose.NewChain[map[string]any, *schema.Message]()
	_ = chain.AppendChatTemplate(promptTemplate)
	_ = chain.AppendChatModel(chatOpenAI)
	runnable, err := chain.Compile(ctx)
	if err != nil {
		logger.Error("[AIService] failed to compile chain", zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	input := map[string]interface{}{
		"context":     req.Context,
		"instruction": req.Instruction,
		"reference":   req.Reference,
	}

	userUniqueID := fmt.Sprintf("%s-%d", user.Name, req.UserID)

	langfuseCallbackHandler, _ := langfuse.NewLangfuseHandler(&langfuse.Config{
		Host:      config.LangfuseHost,
		PublicKey: config.LangfusePublicKey,
		SecretKey: config.LangfuseSecretKey,
		UserID:    userUniqueID,
		Name:      fmt.Sprintf("%s-trace", string(latestPrompt.Task)),
		Tags: []string{
			string(latestPrompt.Task),
		},
	})
	callbackHandlers := []callbacks.Handler{
		langfuseCallbackHandler,
		callback.NewLogCallbackHandler(),
	}

	tokenChan, errChan := make(chan string), make(chan error)
	go func() {
		defer close(tokenChan)
		defer close(errChan)

		sr, err := runnable.Stream(ctx, input, compose.WithCallbacks(callbackHandlers...))
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			logger.Error("[AIService] failed to stream", zap.Error(err))
			errChan <- err
			return
		}
		defer sr.Close()

		for {
			chunk, err := sr.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return
				}
				logger.Error("[AIService] failed to receive stream", zap.Error(err))
				errChan <- err
				return
			}

			tokenChan <- chunk.Content
		}
	}()

	lo.Must0(s.userDAO.Update(db, user, map[string]interface{}{"llm_quota": user.LLMQuota - 1}))

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
func (s *aiService) GenerateArticleSummary(ctx context.Context, req *protocol.GenerateArticleSummaryRequest) (rsp *protocol.GenerateArticleSummaryResponse, err error) {
	rsp = &protocol.GenerateArticleSummaryResponse{}

	logger := logger.LoggerWithContext(ctx)
	db := database.GetDBInstance(ctx)

	user := lo.Must1(s.userDAO.GetByID(db, req.UserID, []string{"id", "name", "llm_quota"}, []string{}))
	if user.LLMQuota <= 0 {
		logger.Info("[AIService] insufficient LLM quota", zap.Uint("userID", req.UserID), zap.Int("quota", int(user.LLMQuota)))
		return nil, protocol.ErrInsufficientQuota
	}

	article, err := s.articleDAO.GetByID(db, req.ArticleID, []string{"id", "title"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[AIService] article not found",
				zap.Uint("articleID", req.ArticleID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[AIService] failed to get article",
			zap.Uint("articleID", req.ArticleID),
			zap.Uint("userID", req.UserID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	latestVersion, err := s.articleVersionDAO.GetLatestByArticleID(db, article.ID, []string{"id", "content"}, []string{})
	if err != nil {
		logger.Error("[AIService] failed to get article version", zap.Uint("articleID", article.ID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	latestPrompt, err := s.promptDAO.GetLatestPromptByTask(db, model.TaskArticleSummary, []string{"id", "task", "templates"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[AIService] latest prompt not found", zap.String("taskName", string(latestPrompt.Task)))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[AIService] failed to get latest prompt", zap.String("taskName", string(latestPrompt.Task)), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	messages := lo.Map(latestPrompt.Templates, func(template model.Template, _ int) schema.MessagesTemplate {
		return &schema.Message{
			Name:    string(latestPrompt.Task),
			Role:    schema.RoleType(template.Role),
			Content: template.Content,
		}
	})

	promptTemplate := prompt.FromMessages(schema.GoTemplate, messages...)

	chatOpenAI, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		Model:       config.OpenAIModel,
		APIKey:      config.OpenAIAPIKey,
		BaseURL:     config.OpenAIBaseURL,
		Temperature: &req.Temperature,
	})
	if err != nil {
		logger.Error("[AIService] failed to create chat openai", zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	chain := compose.NewChain[map[string]any, *schema.Message]()
	_ = chain.AppendChatTemplate(promptTemplate)
	_ = chain.AppendChatModel(chatOpenAI)
	runnable, err := chain.Compile(ctx)
	if err != nil {
		logger.Error("[AIService] failed to compile chain", zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	input := map[string]interface{}{
		"title":       article.Title,
		"content":     latestVersion.Content,
		"instruction": req.Instruction,
	}

	userUniqueID := fmt.Sprintf("%s-%d", user.Name, req.UserID)

	langfuseCallbackHandler, _ := langfuse.NewLangfuseHandler(&langfuse.Config{
		Host:      config.LangfuseHost,
		PublicKey: config.LangfusePublicKey,
		SecretKey: config.LangfuseSecretKey,
		UserID:    userUniqueID,
		Name:      fmt.Sprintf("%s-trace", string(latestPrompt.Task)),
		Tags: []string{
			fmt.Sprintf("%d", req.ArticleID),
			string(latestPrompt.Task),
		},
	})
	callbackHandlers := []callbacks.Handler{
		langfuseCallbackHandler,
		callback.NewLogCallbackHandler(),
	}

	tokenChan, errChan := make(chan string), make(chan error)
	go func() {
		defer close(tokenChan)
		defer close(errChan)

		sr, err := runnable.Stream(ctx, input, compose.WithCallbacks(callbackHandlers...))
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			logger.Error("[AIService] failed to stream", zap.Error(err))
			errChan <- err
			return
		}
		defer sr.Close()

		for {
			chunk, err := sr.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return
				}
				logger.Error("[AIService] failed to receive stream", zap.Error(err))
				errChan <- err
				return
			}

			tokenChan <- chunk.Content
		}
	}()

	lo.Must0(s.userDAO.Update(db, user, map[string]interface{}{"llm_quota": user.LLMQuota - 1}))

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
func (s *aiService) GenerateArticleTranslation(_ context.Context, _ *protocol.GenerateArticleTranslationRequest) (rsp *protocol.GenerateArticleTranslationResponse, err error) {
	// TODO: 实现
	return nil, protocol.ErrNoImplement
}

// GenerateArticleQA 生成文章问答
//
//	receiver s *aiService
//	param req *protocol.GenerateArticleQARequest
//	return rsp *protocol.GenerateArticleQAResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 18:03:44
func (s *aiService) GenerateArticleQA(ctx context.Context, req *protocol.GenerateArticleQARequest) (rsp *protocol.GenerateArticleQAResponse, err error) {
	rsp = &protocol.GenerateArticleQAResponse{}

	logger := logger.LoggerWithContext(ctx)
	db := database.GetDBInstance(ctx)

	user := lo.Must1(s.userDAO.GetByID(db, req.UserID, []string{"id", "name", "llm_quota"}, []string{}))
	if user.LLMQuota <= 0 {
		logger.Info("[AIService] insufficient LLM quota", zap.Uint("userID", req.UserID), zap.Int("quota", int(user.LLMQuota)))
		return nil, protocol.ErrInsufficientQuota
	}

	article, err := s.articleDAO.GetByID(db, req.ArticleID, []string{"id", "title"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[AIService] article not found",
				zap.Uint("articleID", req.ArticleID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[AIService] failed to get article",
			zap.Uint("articleID", req.ArticleID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	latestVersion, err := s.articleVersionDAO.GetLatestByArticleID(db, article.ID, []string{"id", "content"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[AIService] article version not found",
				zap.Uint("articleID", article.ID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[AIService] failed to get article version",
			zap.Uint("articleID", article.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	latestPrompt, err := s.promptDAO.GetLatestPromptByTask(db, model.TaskArticleQA, []string{"id", "task", "templates"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[AIService] latest prompt not found",
				zap.String("taskName", string(latestPrompt.Task)))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[AIService] failed to get latest prompt",
			zap.String("taskName", string(latestPrompt.Task)),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	messages := lo.Map(latestPrompt.Templates, func(template model.Template, _ int) schema.MessagesTemplate {
		return &schema.Message{
			Name:    string(latestPrompt.Task),
			Role:    schema.RoleType(template.Role),
			Content: template.Content,
		}
	})

	promptTemplate := prompt.FromMessages(schema.GoTemplate, messages...)

	chatOpenAI, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		Model:       config.OpenAIModel,
		APIKey:      config.OpenAIAPIKey,
		BaseURL:     config.OpenAIBaseURL,
		Temperature: &req.Temperature,
	})
	if err != nil {
		logger.Error("[AIService] failed to create chat openai", zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	chain := compose.NewChain[map[string]any, *schema.Message]()
	_ = chain.AppendChatTemplate(promptTemplate)
	_ = chain.AppendChatModel(chatOpenAI)
	runnable, err := chain.Compile(ctx)
	if err != nil {
		logger.Error("[AIService] failed to compile chain", zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	input := map[string]interface{}{
		"title":    article.Title,
		"content":  latestVersion.Content,
		"question": req.Question,
	}

	userUniqueID := fmt.Sprintf("%s-%d", user.Name, req.UserID)

	langfuseCallbackHandler, _ := langfuse.NewLangfuseHandler(&langfuse.Config{
		Host:      config.LangfuseHost,
		PublicKey: config.LangfusePublicKey,
		SecretKey: config.LangfuseSecretKey,
		UserID:    userUniqueID,
		Name:      fmt.Sprintf("%s-trace", string(latestPrompt.Task)),
		Tags: []string{
			fmt.Sprintf("%d", req.ArticleID),
			string(latestPrompt.Task),
		},
	})
	callbackHandlers := []callbacks.Handler{
		langfuseCallbackHandler,
		callback.NewLogCallbackHandler(),
	}

	tokenChan, errChan := make(chan string), make(chan error)
	go func() {
		defer close(tokenChan)
		defer close(errChan)

		sr, err := runnable.Stream(ctx, input, compose.WithCallbacks(callbackHandlers...))
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			logger.Error("[AIService] failed to stream", zap.Error(err))
			errChan <- err
			return
		}
		defer sr.Close()

		for {
			chunk, err := sr.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return
				}
				logger.Error("[AIService] failed to receive stream", zap.Error(err))
				errChan <- err
				return
			}

			tokenChan <- chunk.Content
		}
	}()

	lo.Must0(s.userDAO.Update(db, user, map[string]interface{}{"llm_quota": user.LLMQuota - 1}))

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
func (s *aiService) GenerateTermExplaination(ctx context.Context, req *protocol.GenerateTermExplainationRequest) (rsp *protocol.GenerateTermExplainationResponse, err error) {
	rsp = &protocol.GenerateTermExplainationResponse{}

	logger := logger.LoggerWithContext(ctx)
	db := database.GetDBInstance(ctx)

	user := lo.Must1(s.userDAO.GetByID(db, req.UserID, []string{"id", "name", "llm_quota"}, []string{}))
	if user.LLMQuota <= 0 {
		logger.Info("[AIService] insufficient LLM quota", zap.Uint("userID", req.UserID), zap.Int("quota", int(user.LLMQuota)))
		return nil, protocol.ErrInsufficientQuota
	}

	article, err := s.articleDAO.GetByID(db, req.ArticleID, []string{"id", "title"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[AIService] article not found",
				zap.Uint("articleID", req.ArticleID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[AIService] failed to get article",
			zap.Uint("articleID", req.ArticleID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	latestVersion, err := s.articleVersionDAO.GetLatestByArticleID(db, article.ID, []string{"id", "content"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[AIService] article version not found", zap.Uint("articleID", article.ID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[AIService] failed to get article version", zap.Uint("articleID", article.ID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	latestPrompt, err := s.promptDAO.GetLatestPromptByTask(db, model.TaskTermExplaination, []string{"id", "task", "templates"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[AIService] latest prompt not found", zap.String("taskName", string(latestPrompt.Task)))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[AIService] failed to get latest prompt", zap.String("taskName", string(latestPrompt.Task)), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	messages := lo.Map(latestPrompt.Templates, func(template model.Template, _ int) schema.MessagesTemplate {
		return &schema.Message{
			Name:    string(latestPrompt.Task),
			Role:    schema.RoleType(template.Role),
			Content: template.Content,
		}
	})

	promptTemplate := prompt.FromMessages(schema.GoTemplate, messages...)

	chatOpenAI, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		Model:       config.OpenAIModel,
		APIKey:      config.OpenAIAPIKey,
		BaseURL:     config.OpenAIBaseURL,
		Temperature: &req.Temperature,
	})
	if err != nil {
		logger.Error("[AIService] failed to create chat openai", zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	chain := compose.NewChain[map[string]any, *schema.Message]()
	_ = chain.AppendChatTemplate(promptTemplate)
	_ = chain.AppendChatModel(chatOpenAI)
	runnable, err := chain.Compile(ctx)
	if err != nil {
		logger.Error("[AIService] failed to compile chain", zap.Error(err))
		return nil, protocol.ErrInternalError
	}

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

	input := map[string]interface{}{
		"title":   article.Title,
		"content": latestVersion.Content,
		"context": latestVersion.Content[left:right],
		"term":    req.Term,
	}

	userUniqueID := fmt.Sprintf("%s-%d", user.Name, req.UserID)

	langfuseCallbackHandler, _ := langfuse.NewLangfuseHandler(&langfuse.Config{
		Host:      config.LangfuseHost,
		PublicKey: config.LangfusePublicKey,
		SecretKey: config.LangfuseSecretKey,
		UserID:    userUniqueID,
		Name:      fmt.Sprintf("%s-trace", string(latestPrompt.Task)),
		Tags: []string{
			fmt.Sprintf("%d", req.ArticleID),
			string(latestPrompt.Task),
		},
	})
	callbackHandlers := []callbacks.Handler{
		langfuseCallbackHandler,
		callback.NewLogCallbackHandler(),
	}

	tokenChan, errChan := make(chan string), make(chan error)
	go func() {
		defer close(tokenChan)
		defer close(errChan)

		sr, err := runnable.Stream(ctx, input, compose.WithCallbacks(callbackHandlers...))
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			logger.Error("[AIService] failed to stream", zap.Error(err))
			errChan <- err
			return
		}
		defer sr.Close()
		for {
			chunk, err := sr.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return
				}
				logger.Error("[AIService] failed to receive stream", zap.Error(err))
				errChan <- err
				return
			}

			tokenChan <- chunk.Content
		}
	}()

	lo.Must0(s.userDAO.Update(db, user, map[string]interface{}{"llm_quota": user.LLMQuota - 1}))

	rsp.TokenChan = tokenChan
	rsp.ErrChan = errChan

	return rsp, nil
}
