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
	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/logger"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/protocol/dto"
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
	GetPrompt(ctx context.Context, req *dto.GetPromptRequest) (rsp *dto.GetPromptResponse, err error)
	GetLatestPrompt(ctx context.Context, req *dto.GetLatestPromptRequest) (rsp *dto.GetLatestPromptResponse, err error)
	ListPrompt(ctx context.Context, req *dto.ListPromptRequest) (rsp *dto.ListPromptResponse, err error)
	CreatePrompt(ctx context.Context, req *dto.CreatePromptRequest) (rsp *dto.EmptyResponse, err error)
	GenerateContentCompletion(ctx context.Context, req *dto.GenerateContentCompletionRequest) (tokenChan <-chan string, errChan <-chan error)
	GenerateArticleSummary(ctx context.Context, req *dto.GenerateArticleSummaryRequest) (tokenChan <-chan string, errChan <-chan error)
	GenerateArticleTranslation(ctx context.Context, req *dto.GenerateArticleQARequest) (tokenChan <-chan string, errChan <-chan error)
	GenerateArticleQA(ctx context.Context, req *dto.GenerateArticleQARequest) (tokenChan <-chan string, errChan <-chan error)
	GenerateTermExplaination(ctx context.Context, req *dto.GenerateTermExplainationRequest) (tokenChan <-chan string, errChan <-chan error)
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
//	param req *dto.GetPromptRequest
//	return rsp *dto.GetPromptResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 18:02:44
func (s *aiService) GetPrompt(ctx context.Context, req *dto.GetPromptRequest) (rsp *dto.GetPromptResponse, err error) {
	rsp = &dto.GetPromptResponse{}

	logger := logger.WithCtx(ctx)
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

	rsp.Prompt = &dto.Prompt{
		PromptID:  prompt.ID,
		CreatedAt: prompt.CreatedAt.Format(time.DateTime),
		Task:      string(prompt.Task),
		Version:   prompt.Version,
		Templates: lo.Map(prompt.Templates, func(t model.Template, _ int) dto.Template {
			return dto.Template{
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
//	param req *dto.GetLatestPromptRequest
//	return rsp *dto.GetLatestPromptResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 18:02:55
func (s *aiService) GetLatestPrompt(ctx context.Context, req *dto.GetLatestPromptRequest) (rsp *dto.GetLatestPromptResponse, err error) {
	rsp = &dto.GetLatestPromptResponse{}

	logger := logger.WithCtx(ctx)
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

	rsp.Prompt = &dto.Prompt{
		PromptID:  prompt.ID,
		CreatedAt: prompt.CreatedAt.Format(time.DateTime),
		Task:      string(prompt.Task),
		Version:   prompt.Version,
		Templates: lo.Map(prompt.Templates, func(t model.Template, _ int) dto.Template {
			return dto.Template{
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
//	param req *dto.ListPromptRequest
//	return rsp *dto.ListPromptResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 18:02:55
func (s *aiService) ListPrompt(ctx context.Context, req *dto.ListPromptRequest) (rsp *dto.ListPromptResponse, err error) {
	rsp = &dto.ListPromptResponse{}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	param := &dao.CommonParam{
		PageParam: &dao.PageParam{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		QueryParam: &dao.QueryParam{
			Query:       req.Query,
			QueryFields: []string{"task", "version"},
		},
	}
	prompts, pageInfo, err := s.promptDAO.PaginateByTask(db, model.Task(req.TaskName),
		[]string{"id", "created_at", "task", "version", "templates", "variables"},
		[]string{},
		param,
	)
	if err != nil {
		logger.Error("[AIService] failed to list prompt", zap.String("taskName", req.TaskName), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Prompts = lo.Map(prompts, func(p *model.Prompt, _ int) *dto.Prompt {
		return &dto.Prompt{
			PromptID:  p.ID,
			CreatedAt: p.CreatedAt.Format(time.DateTime),
			Task:      string(p.Task),
			Version:   p.Version,
			Templates: lo.Map(p.Templates, func(t model.Template, _ int) dto.Template {
				return dto.Template{
					Role:    t.Role,
					Content: t.Content,
				}
			}),
			Variables: p.Variables,
		}
	})

	rsp.PageInfo = &dto.PageInfo{
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
		Total:    pageInfo.Total,
	}

	return rsp, nil
}

// CreatePrompt 创建提示词
//
//	receiver s *aiService
//	param req *dto.CreatePromptRequest
//	return rsp *dto.EmptyResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 18:03:07
func (s *aiService) CreatePrompt(ctx context.Context, req *dto.CreatePromptRequest) (rsp *dto.EmptyResponse, err error) {
	if req == nil || req.Body == nil {
		return nil, protocol.ErrBadRequest
	}

	rsp = &dto.EmptyResponse{}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	contents := lo.Map(req.Body.Templates, func(tmplate dto.Template, _ int) string {
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
		logger.Info("[AIService] the content of the new version is the same as the latest version", zap.String("taskName", req.TaskName), zap.Any("templates", req.Body.Templates))
		return nil, protocol.ErrBadRequest
	}

	if l, r := lo.Difference(prompt.Variables, variables); prompt.ID != 0 && len(l)+len(r) > 0 {
		logger.Info("[AIService] the variables of the latest prompt and the new prompt are mismatch",
			zap.String("taskName", req.TaskName), zap.Strings("latestVariables", prompt.Variables), zap.Strings("newVariables", variables))
		return nil, protocol.ErrBadRequest
	}

	prompt = &model.Prompt{
		Task: model.Task(req.TaskName),
		Templates: lo.Map(req.Body.Templates, func(tmplate dto.Template, _ int) model.Template {
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
//	param req *dto.GenerateContentCompletionRequest
//	return tokenChan <-chan string
//	return errChan <-chan error
//	author centonhuang
//	update 2025-11-01 18:30:00
func (s *aiService) GenerateContentCompletion(ctx context.Context, req *dto.GenerateContentCompletionRequest) (tokenChan <-chan string, errChan <-chan error) {
	errCh := make(chan error, 1)

	if req == nil || req.Body == nil {
		errCh <- protocol.ErrBadRequest
		return nil, errCh
	}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	user := lo.Must1(s.userDAO.GetByID(db, userID, []string{"id", "name", "llm_quota"}, []string{}))
	if user.LLMQuota <= 0 {
		logger.Info("[AIService] insufficient LLM quota", zap.Int("quota", int(user.LLMQuota)))
		errCh <- protocol.ErrInsufficientQuota
		close(errCh)
		return nil, errCh
	}

	latestPrompt, err := s.promptDAO.GetLatestPromptByTask(db, model.TaskContentCompletion, []string{"id", "task", "templates"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[AIService] latest prompt not found", zap.String("taskName", string(latestPrompt.Task)))
			errCh <- protocol.ErrDataNotExists
			close(errCh)
			return nil, errCh
		}
		logger.Error("[AIService] failed to get latest prompt", zap.String("taskName", string(latestPrompt.Task)), zap.Error(err))
		errCh <- protocol.ErrInternalError
		close(errCh)
		return nil, errCh
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
		Temperature: &req.Body.Temperature,
	})
	if err != nil {
		logger.Error("[AIService] failed to create chat openai", zap.Error(err))
		errCh <- protocol.ErrInternalError
		close(errCh)
		return nil, errCh
	}

	chain := compose.NewChain[map[string]any, *schema.Message]()
	_ = chain.AppendChatTemplate(promptTemplate)
	_ = chain.AppendChatModel(chatOpenAI)
	runnable, err := chain.Compile(ctx)
	if err != nil {
		logger.Error("[AIService] failed to compile chain", zap.Error(err))
		errCh <- protocol.ErrInternalError
		close(errCh)
		return nil, errCh
	}

	input := map[string]interface{}{
		"context":     req.Body.Context,
		"instruction": req.Body.Instruction,
		"reference":   req.Body.Reference,
	}

	userUniqueID := fmt.Sprintf("%s-%d", user.Name, userID)

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

	tokenCh := make(chan string)
	go func() {
		defer close(tokenCh)
		defer close(errCh)

		sr, err := runnable.Stream(ctx, input, compose.WithCallbacks(callbackHandlers...))
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			logger.Error("[AIService] failed to stream", zap.Error(err))
			errCh <- err
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
				errCh <- err
				return
			}

			tokenCh <- chunk.Content
		}
	}()

	lo.Must0(s.userDAO.Update(db, user, map[string]interface{}{"llm_quota": user.LLMQuota - 1}))

	return tokenCh, errCh
}

// GenerateArticleSummary 生成文章总结
//
//	receiver s *aiService
//	param req *dto.GenerateArticleSummaryRequest
//	return tokenChan <-chan string
//	return errChan <-chan error
//	author centonhuang
//	update 2025-11-01 18:30:00
func (s *aiService) GenerateArticleSummary(ctx context.Context, req *dto.GenerateArticleSummaryRequest) (tokenChan <-chan string, errChan <-chan error) {
	errCh := make(chan error, 1)

	if req == nil || req.Body == nil {
		errCh <- protocol.ErrBadRequest
		return nil, errCh
	}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	user := lo.Must1(s.userDAO.GetByID(db, userID, []string{"id", "name", "llm_quota"}, []string{}))
	if user.LLMQuota <= 0 {
		logger.Info("[AIService] insufficient LLM quota", zap.Int("quota", int(user.LLMQuota)))
		errCh <- protocol.ErrInsufficientQuota
		close(errCh)
		return nil, errCh
	}

	article, err := s.articleDAO.GetByID(db, req.Body.ArticleID, []string{"id", "title"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[AIService] article not found",
				zap.Uint("articleID", req.Body.ArticleID))
			errCh <- protocol.ErrDataNotExists
			close(errCh)
			return nil, errCh
		}
		logger.Error("[AIService] failed to get article",
			zap.Uint("articleID", req.Body.ArticleID),
			zap.Error(err))
		errCh <- protocol.ErrInternalError
		close(errCh)
		return nil, errCh
	}

	latestVersion, err := s.articleVersionDAO.GetLatestByArticleID(db, article.ID, []string{"id", "content"}, []string{})
	if err != nil {
		logger.Error("[AIService] failed to get article version", zap.Uint("articleID", article.ID), zap.Error(err))
		errCh <- protocol.ErrInternalError
		close(errCh)
		return nil, errCh
	}

	latestPrompt, err := s.promptDAO.GetLatestPromptByTask(db, model.TaskArticleSummary, []string{"id", "task", "templates"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[AIService] latest prompt not found", zap.String("taskName", string(model.TaskArticleSummary)))
			errCh <- protocol.ErrDataNotExists
			close(errCh)
			return nil, errCh
		}
		logger.Error("[AIService] failed to get latest prompt", zap.String("taskName", string(model.TaskArticleSummary)), zap.Error(err))
		errCh <- protocol.ErrInternalError
		close(errCh)
		return nil, errCh
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
		Temperature: &req.Body.Temperature,
	})
	if err != nil {
		logger.Error("[AIService] failed to create chat openai", zap.Error(err))
		errCh <- protocol.ErrInternalError
		close(errCh)
		return nil, errCh
	}

	chain := compose.NewChain[map[string]any, *schema.Message]()
	_ = chain.AppendChatTemplate(promptTemplate)
	_ = chain.AppendChatModel(chatOpenAI)
	runnable, err := chain.Compile(ctx)
	if err != nil {
		logger.Error("[AIService] failed to compile chain", zap.Error(err))
		errCh <- protocol.ErrInternalError
		close(errCh)
		return nil, errCh
	}

	input := map[string]interface{}{
		"title":       article.Title,
		"content":     latestVersion.Content,
		"instruction": req.Body.Instruction,
	}

	userUniqueID := fmt.Sprintf("%s-%d", user.Name, userID)

	langfuseCallbackHandler, _ := langfuse.NewLangfuseHandler(&langfuse.Config{
		Host:      config.LangfuseHost,
		PublicKey: config.LangfusePublicKey,
		SecretKey: config.LangfuseSecretKey,
		UserID:    userUniqueID,
		Name:      fmt.Sprintf("%s-trace", string(latestPrompt.Task)),
		Tags: []string{
			fmt.Sprintf("%d", req.Body.ArticleID),
			string(latestPrompt.Task),
		},
	})
	callbackHandlers := []callbacks.Handler{
		langfuseCallbackHandler,
		callback.NewLogCallbackHandler(),
	}

	tokenCh := make(chan string)
	go func() {
		defer close(tokenCh)
		defer close(errCh)

		sr, err := runnable.Stream(ctx, input, compose.WithCallbacks(callbackHandlers...))
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			logger.Error("[AIService] failed to stream", zap.Error(err))
			errCh <- err
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
				errCh <- err
				return
			}

			tokenCh <- chunk.Content
		}
	}()

	lo.Must0(s.userDAO.Update(db, user, map[string]interface{}{"llm_quota": user.LLMQuota - 1}))

	return tokenCh, errCh
}

// GenerateArticleTranslation 生成文章翻译
//
//	receiver s *aiService
//	param req *dto.GenerateArticleQARequest
//	return tokenChan <-chan string
//	return errChan <-chan error
//	author centonhuang
//	update 2025-11-01 18:30:00
func (s *aiService) GenerateArticleTranslation(ctx context.Context, req *dto.GenerateArticleQARequest) (tokenChan <-chan string, errChan <-chan error) {
	errCh := make(chan error, 1)

	if req == nil || req.Body == nil {
		errCh <- protocol.ErrBadRequest
		return nil, errCh
	}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	user := lo.Must1(s.userDAO.GetByID(db, userID, []string{"id", "name", "llm_quota"}, []string{}))
	if user.LLMQuota <= 0 {
		logger.Info("[AIService] insufficient LLM quota", zap.Int("quota", int(user.LLMQuota)))
		errCh <- protocol.ErrInsufficientQuota
		close(errCh)
		return nil, errCh
	}

	article, err := s.articleDAO.GetByID(db, req.Body.ArticleID, []string{"id", "title"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[AIService] article not found",
				zap.Uint("articleID", req.Body.ArticleID))
			errCh <- protocol.ErrDataNotExists
			close(errCh)
			return nil, errCh
		}
		logger.Error("[AIService] failed to get article",
			zap.Uint("articleID", req.Body.ArticleID),
			zap.Error(err))
		errCh <- protocol.ErrInternalError
		close(errCh)
		return nil, errCh
	}

	latestVersion, err := s.articleVersionDAO.GetLatestByArticleID(db, article.ID, []string{"id", "content"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[AIService] article version not found",
				zap.Uint("articleID", article.ID))
			errCh <- protocol.ErrDataNotExists
			close(errCh)
			return nil, errCh
		}
		logger.Error("[AIService] failed to get article version",
			zap.Uint("articleID", article.ID),
			zap.Error(err))
		errCh <- protocol.ErrInternalError
		close(errCh)
		return nil, errCh
	}

	latestPrompt, err := s.promptDAO.GetLatestPromptByTask(db, model.TaskArticleTranslation, []string{"id", "task", "templates"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[AIService] latest prompt not found", zap.String("taskName", string(model.TaskArticleTranslation)))
			errCh <- protocol.ErrDataNotExists
			close(errCh)
			return nil, errCh
		}
		logger.Error("[AIService] failed to get latest prompt", zap.String("taskName", string(model.TaskArticleTranslation)), zap.Error(err))
		errCh <- protocol.ErrInternalError
		close(errCh)
		return nil, errCh
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
		Temperature: &req.Body.Temperature,
	})
	if err != nil {
		logger.Error("[AIService] failed to create chat openai", zap.Error(err))
		errCh <- protocol.ErrInternalError
		close(errCh)
		return nil, errCh
	}

	chain := compose.NewChain[map[string]any, *schema.Message]()
	_ = chain.AppendChatTemplate(promptTemplate)
	_ = chain.AppendChatModel(chatOpenAI)
	runnable, err := chain.Compile(ctx)
	if err != nil {
		logger.Error("[AIService] failed to compile chain", zap.Error(err))
		errCh <- protocol.ErrInternalError
		close(errCh)
		return nil, errCh
	}

	input := map[string]interface{}{
		"title":   article.Title,
		"content": latestVersion.Content,
	}

	userUniqueID := fmt.Sprintf("%s-%d", user.Name, userID)

	langfuseCallbackHandler, _ := langfuse.NewLangfuseHandler(&langfuse.Config{
		Host:      config.LangfuseHost,
		PublicKey: config.LangfusePublicKey,
		SecretKey: config.LangfuseSecretKey,
		UserID:    userUniqueID,
		Name:      fmt.Sprintf("%s-trace", string(latestPrompt.Task)),
		Tags: []string{
			fmt.Sprintf("%d", req.Body.ArticleID),
			string(latestPrompt.Task),
		},
	})
	callbackHandlers := []callbacks.Handler{
		langfuseCallbackHandler,
		callback.NewLogCallbackHandler(),
	}

	tokenCh := make(chan string)
	go func() {
		defer close(tokenCh)
		defer close(errCh)

		sr, err := runnable.Stream(ctx, input, compose.WithCallbacks(callbackHandlers...))
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			logger.Error("[AIService] failed to stream", zap.Error(err))
			errCh <- err
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
				errCh <- err
				return
			}

			tokenCh <- chunk.Content
		}
	}()

	lo.Must0(s.userDAO.Update(db, user, map[string]interface{}{"llm_quota": user.LLMQuota - 1}))

	return tokenCh, errCh
}

// GenerateArticleQA 生成文章问答
//
//	receiver s *aiService
//	param req *dto.GenerateArticleQARequest
//	return tokenChan <-chan string
//	return errChan <-chan error
//	author centonhuang
//	update 2025-01-05 18:03:44
func (s *aiService) GenerateArticleQA(ctx context.Context, req *dto.GenerateArticleQARequest) (tokenChan <-chan string, errChan <-chan error) {
	errCh := make(chan error, 1)

	if req == nil || req.Body == nil {
		errCh <- protocol.ErrBadRequest
		return nil, errCh
	}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	user := lo.Must1(s.userDAO.GetByID(db, userID, []string{"id", "name", "llm_quota"}, []string{}))
	if user.LLMQuota <= 0 {
		logger.Info("[AIService] insufficient LLM quota", zap.Int("quota", int(user.LLMQuota)))
		errCh <- protocol.ErrInsufficientQuota
		close(errCh)
		return nil, errCh
	}

	article, err := s.articleDAO.GetByID(db, req.Body.ArticleID, []string{"id", "title"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[AIService] article not found",
				zap.Uint("articleID", req.Body.ArticleID))
			errCh <- protocol.ErrDataNotExists
			close(errCh)
			return nil, errCh
		}
		logger.Error("[AIService] failed to get article",
			zap.Uint("articleID", req.Body.ArticleID),
			zap.Error(err))
		errCh <- protocol.ErrInternalError
		close(errCh)
		return nil, errCh
	}

	latestVersion, err := s.articleVersionDAO.GetLatestByArticleID(db, article.ID, []string{"id", "content"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[AIService] article version not found",
				zap.Uint("articleID", article.ID))
			errCh <- protocol.ErrDataNotExists
			close(errCh)
			return nil, errCh
		}
		logger.Error("[AIService] failed to get article version",
			zap.Uint("articleID", article.ID),
			zap.Error(err))
		errCh <- protocol.ErrInternalError
		close(errCh)
		return nil, errCh
	}

	latestPrompt, err := s.promptDAO.GetLatestPromptByTask(db, model.TaskArticleQA, []string{"id", "task", "templates"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[AIService] latest prompt not found",
				zap.String("taskName", string(model.TaskArticleQA)))
			errCh <- protocol.ErrDataNotExists
			close(errCh)
			return nil, errCh
		}
		logger.Error("[AIService] failed to get latest prompt",
			zap.String("taskName", string(model.TaskArticleQA)),
			zap.Error(err))
		errCh <- protocol.ErrInternalError
		close(errCh)
		return nil, errCh
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
		Temperature: &req.Body.Temperature,
	})
	if err != nil {
		logger.Error("[AIService] failed to create chat openai", zap.Error(err))
		errCh <- protocol.ErrInternalError
		close(errCh)
		return nil, errCh
	}

	chain := compose.NewChain[map[string]any, *schema.Message]()
	_ = chain.AppendChatTemplate(promptTemplate)
	_ = chain.AppendChatModel(chatOpenAI)
	runnable, err := chain.Compile(ctx)
	if err != nil {
		logger.Error("[AIService] failed to compile chain", zap.Error(err))
		errCh <- protocol.ErrInternalError
		close(errCh)
		return nil, errCh
	}

	input := map[string]interface{}{
		"title":    article.Title,
		"content":  latestVersion.Content,
		"question": req.Body.Question,
	}

	userUniqueID := fmt.Sprintf("%s-%d", user.Name, userID)

	langfuseCallbackHandler, _ := langfuse.NewLangfuseHandler(&langfuse.Config{
		Host:      config.LangfuseHost,
		PublicKey: config.LangfusePublicKey,
		SecretKey: config.LangfuseSecretKey,
		UserID:    userUniqueID,
		Name:      fmt.Sprintf("%s-trace", string(latestPrompt.Task)),
		Tags: []string{
			fmt.Sprintf("%d", req.Body.ArticleID),
			string(latestPrompt.Task),
		},
	})
	callbackHandlers := []callbacks.Handler{
		langfuseCallbackHandler,
		callback.NewLogCallbackHandler(),
	}

	tokenCh := make(chan string)
	go func() {
		defer close(tokenCh)
		defer close(errCh)

		sr, err := runnable.Stream(ctx, input, compose.WithCallbacks(callbackHandlers...))
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			logger.Error("[AIService] failed to stream", zap.Error(err))
			errCh <- err
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
				errCh <- err
				return
			}

			tokenCh <- chunk.Content
		}
	}()

	lo.Must0(s.userDAO.Update(db, user, map[string]interface{}{"llm_quota": user.LLMQuota - 1}))

	return tokenCh, errCh
}

// GenerateTermExplaination 生成术语解释
//
//	receiver s *aiService
//	param req *protocol.GenerateTermExplainationRequest
//	return rsp *protocol.GenerateTermExplainationResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 18:03:48
func (s *aiService) GenerateTermExplaination(ctx context.Context, req *dto.GenerateTermExplainationRequest) (tokenChan <-chan string, errChan <-chan error) {
	errCh := make(chan error, 1)

	if req == nil || req.Body == nil {
		errCh <- protocol.ErrBadRequest
		return nil, errCh
	}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	user := lo.Must1(s.userDAO.GetByID(db, userID, []string{"id", "name", "llm_quota"}, []string{}))
	if user.LLMQuota <= 0 {
		logger.Info("[AIService] insufficient LLM quota", zap.Int("quota", int(user.LLMQuota)))
		errCh <- protocol.ErrInsufficientQuota
		close(errCh)
		return nil, errCh
	}

	article, err := s.articleDAO.GetByID(db, req.Body.ArticleID, []string{"id", "title"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[AIService] article not found",
				zap.Uint("articleID", req.Body.ArticleID))
			errCh <- protocol.ErrDataNotExists
			close(errCh)
			return nil, errCh
		}
		logger.Error("[AIService] failed to get article",
			zap.Uint("articleID", req.Body.ArticleID),
			zap.Error(err))
		errCh <- protocol.ErrInternalError
		close(errCh)
		return nil, errCh
	}

	latestVersion, err := s.articleVersionDAO.GetLatestByArticleID(db, article.ID, []string{"id", "content"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[AIService] article version not found", zap.Uint("articleID", article.ID))
			errCh <- protocol.ErrDataNotExists
			close(errCh)
			return nil, errCh
		}
		logger.Error("[AIService] failed to get article version", zap.Uint("articleID", article.ID), zap.Error(err))
		errCh <- protocol.ErrInternalError
		close(errCh)
		return nil, errCh
	}

	latestPrompt, err := s.promptDAO.GetLatestPromptByTask(db, model.TaskTermExplaination, []string{"id", "task", "templates"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[AIService] latest prompt not found", zap.String("taskName", string(model.TaskTermExplaination)))
			errCh <- protocol.ErrDataNotExists
			close(errCh)
			return nil, errCh
		}
		logger.Error("[AIService] failed to get latest prompt", zap.String("taskName", string(model.TaskTermExplaination)), zap.Error(err))
		errCh <- protocol.ErrInternalError
		close(errCh)
		return nil, errCh
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
		Temperature: &req.Body.Temperature,
	})
	if err != nil {
		logger.Error("[AIService] failed to create chat openai", zap.Error(err))
		errCh <- protocol.ErrInternalError
		close(errCh)
		return nil, errCh
	}

	chain := compose.NewChain[map[string]any, *schema.Message]()
	_ = chain.AppendChatTemplate(promptTemplate)
	_ = chain.AppendChatModel(chatOpenAI)
	runnable, err := chain.Compile(ctx)
	if err != nil {
		logger.Error("[AIService] failed to compile chain", zap.Error(err))
		errCh <- protocol.ErrInternalError
		close(errCh)
		return nil, errCh
	}

	contextWindowLen := 200

	var left, right int
	if req.Body.Position < contextWindowLen/2 {
		left = 0
		right = contextWindowLen
	} else if req.Body.Position > len(latestVersion.Content)-contextWindowLen/2 {
		left = len(latestVersion.Content) - contextWindowLen
		right = len(latestVersion.Content)
	} else {
		left = req.Body.Position - contextWindowLen/2
		right = req.Body.Position + contextWindowLen/2
	}

	input := map[string]interface{}{
		"title":   article.Title,
		"content": latestVersion.Content,
		"context": latestVersion.Content[left:right],
		"term":    req.Body.Term,
	}

	userUniqueID := fmt.Sprintf("%s-%d", user.Name, userID)

	langfuseCallbackHandler, _ := langfuse.NewLangfuseHandler(&langfuse.Config{
		Host:      config.LangfuseHost,
		PublicKey: config.LangfusePublicKey,
		SecretKey: config.LangfuseSecretKey,
		UserID:    userUniqueID,
		Name:      fmt.Sprintf("%s-trace", string(latestPrompt.Task)),
		Tags: []string{
			fmt.Sprintf("%d", req.Body.ArticleID),
			string(latestPrompt.Task),
		},
	})
	callbackHandlers := []callbacks.Handler{
		langfuseCallbackHandler,
		callback.NewLogCallbackHandler(),
	}

	tokenCh := make(chan string)
	go func() {
		defer close(tokenCh)
		defer close(errCh)

		sr, err := runnable.Stream(ctx, input, compose.WithCallbacks(callbackHandlers...))
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			logger.Error("[AIService] failed to stream", zap.Error(err))
			errCh <- err
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
				errCh <- err
				return
			}

			tokenCh <- chunk.Content
		}
	}()

	lo.Must0(s.userDAO.Update(db, user, map[string]interface{}{"llm_quota": user.LLMQuota - 1}))

	return tokenCh, errCh
}
