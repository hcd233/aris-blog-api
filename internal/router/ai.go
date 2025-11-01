package router

import (
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
)

func initAIRouter(aiGroup *huma.Group) {
	aiHandler := handler.NewAIHandler()

	aiGroup.UseMiddleware(middleware.JwtMiddlewareForHuma())

	// Prompt???? - ??Admin??
	promptGroup := huma.NewGroup(aiGroup, "/prompt")
	promptGroup.UseMiddleware(middleware.LimitUserPermissionMiddlewareForHuma("promptService", model.PermissionAdmin))

	// /:taskName/v:version - ???????Prompt
	huma.Register(promptGroup, huma.Operation{
		OperationID: "getPrompt",
		Method:      http.MethodGet,
		Path:        "/{taskName}/v{version}",
		Summary:     "GetPrompt",
		Description: "Get a specific version of a prompt by task name and version",
		Tags:        []string{"ai"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, aiHandler.HandleGetPrompt)

	// /:taskName/latest - ????Prompt
	huma.Register(promptGroup, huma.Operation{
		OperationID: "getLatestPrompt",
		Method:      http.MethodGet,
		Path:        "/{taskName}/latest",
		Summary:     "GetLatestPrompt",
		Description: "Get the latest version of a prompt by task name",
		Tags:        []string{"ai"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, aiHandler.HandleGetLatestPrompt)

	// /:taskName - ??Prompt
	huma.Register(promptGroup, huma.Operation{
		OperationID: "listPrompts",
		Method:      http.MethodGet,
		Path:        "/{taskName}",
		Summary:     "ListPrompts",
		Description: "List all versions of a prompt by task name",
		Tags:        []string{"ai"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, aiHandler.HandleListPrompt)

	// /:taskName - ??Prompt
	huma.Register(promptGroup, huma.Operation{
		OperationID: "createPrompt",
		Method:      http.MethodPost,
		Path:        "/{taskName}",
		Summary:     "CreatePrompt",
		Description: "Create a new version of a prompt",
		Tags:        []string{"ai"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, aiHandler.HandleCreatePrompt)

	// AI????
	appGroup := huma.NewGroup(aiGroup, "/app")

	// Creator??
	creatorGroup := huma.NewGroup(appGroup, "/creator")

	contentCompletionGroup := huma.NewGroup(creatorGroup, "")
	contentCompletionGroup.UseMiddleware(middleware.RedisLockMiddlewareForHuma("contentCompletion", constant.CtxKeyUserID, 30*time.Second))

	huma.Register(contentCompletionGroup, huma.Operation{
		OperationID: "generateContentCompletion",
		Method:      http.MethodPost,
		Path:        "/contentCompletion",
		Summary:     "GenerateContentCompletion",
		Description: "Generate content completion using AI",
		Tags:        []string{"ai"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, aiHandler.HandleGenerateContentCompletion)

	articleSummaryGroup := huma.NewGroup(creatorGroup, "")
	articleSummaryGroup.UseMiddleware(middleware.RedisLockMiddlewareForHuma("articleSummary", constant.CtxKeyUserID, 30*time.Second))

	huma.Register(articleSummaryGroup, huma.Operation{
		OperationID: "generateArticleSummary",
		Method:      http.MethodPost,
		Path:        "/articleSummary",
		Summary:     "GenerateArticleSummary",
		Description: "Generate article summary using AI",
		Tags:        []string{"ai"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, aiHandler.HandleGenerateArticleSummary)

	// Reader??
	readerGroup := huma.NewGroup(appGroup, "/reader")

	articleQAGroup := huma.NewGroup(readerGroup, "")
	articleQAGroup.UseMiddleware(middleware.RedisLockMiddlewareForHuma("articleQA", constant.CtxKeyUserID, 30*time.Second))

	huma.Register(articleQAGroup, huma.Operation{
		OperationID: "generateArticleQA",
		Method:      http.MethodPost,
		Path:        "/articleQA",
		Summary:     "GenerateArticleQA",
		Description: "Generate article Q&A using AI",
		Tags:        []string{"ai"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, aiHandler.HandleGenerateArticleQA)

	// ??????????
	// termExplainationGroup := huma.NewGroup(readerGroup, "")
	// huma.Register(termExplainationGroup, huma.Operation{
	// 	OperationID: "generateTermExplaination",
	// 	Method:      http.MethodPost,
	// 	Path:        "/termExplaination",
	// 	Summary:     "GenerateTermExplaination",
	// 	Description: "Generate term explanation using AI",
	// 	Tags:        []string{"ai"},
	// 	Security:    []map[string][]string{{"jwtAuth": {}}},
	// }, aiHandler.HandleGenerateTermExplaination)
}
