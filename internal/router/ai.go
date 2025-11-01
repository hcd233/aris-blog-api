package router

import (
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/sse"
	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
)

func initAIRouter(aiGroup *huma.Group) {
	aiHandler := handler.NewAIHandler()

	aiGroup.UseMiddleware(middleware.JwtMiddleware())

	promptGroup := huma.NewGroup(aiGroup, "/prompt")
	promptGroup.UseMiddleware(middleware.LimitUserPermissionMiddleware("promptService", model.PermissionAdmin))

	huma.Register(promptGroup, huma.Operation{
		OperationID: "getPrompt",
		Method:      http.MethodGet,
		Path:        "/{taskName}/v{version}",
		Summary:     "GetPrompt",
		Description: "Get a specific version of a prompt by task name and version",
		Tags:        []string{"ai"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, aiHandler.HandleGetPrompt)

	huma.Register(promptGroup, huma.Operation{
		OperationID: "getLatestPrompt",
		Method:      http.MethodGet,
		Path:        "/{taskName}/latest",
		Summary:     "GetLatestPrompt",
		Description: "Get the latest version of a prompt by task name",
		Tags:        []string{"ai"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, aiHandler.HandleGetLatestPrompt)

	huma.Register(promptGroup, huma.Operation{
		OperationID: "listPrompts",
		Method:      http.MethodGet,
		Path:        "/{taskName}",
		Summary:     "ListPrompts",
		Description: "List all versions of a prompt by task name",
		Tags:        []string{"ai"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, aiHandler.HandleListPrompt)

	huma.Register(promptGroup, huma.Operation{
		OperationID: "createPrompt",
		Method:      http.MethodPost,
		Path:        "/{taskName}",
		Summary:     "CreatePrompt",
		Description: "Create a new version of a prompt",
		Tags:        []string{"ai"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, aiHandler.HandleCreatePrompt)

	appGroup := huma.NewGroup(aiGroup, "/app")

	creatorGroup := huma.NewGroup(appGroup, "/creator")

	contentCompletionGroup := huma.NewGroup(creatorGroup, "")
	contentCompletionGroup.UseMiddleware(middleware.RedisLockMiddleware("contentCompletion", constant.CtxKeyUserID, 30*time.Second))

	sse.Register(contentCompletionGroup, huma.Operation{
		OperationID: "generateContentCompletion",
		Method:      http.MethodPost,
		Path:        "/contentCompletion",
		Summary:     "GenerateContentCompletion",
		Description: "Generate content completion using AI",
		Tags:        []string{"ai"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	},
		map[string]any{
			"SSEResponse": protocol.SSEResponse{},
		},
		aiHandler.HandleGenerateContentCompletion)

	articleSummaryGroup := huma.NewGroup(creatorGroup, "")
	articleSummaryGroup.UseMiddleware(middleware.RedisLockMiddleware("articleSummary", constant.CtxKeyUserID, 30*time.Second))

	sse.Register(articleSummaryGroup, huma.Operation{
		OperationID: "generateArticleSummary",
		Method:      http.MethodPost,
		Path:        "/articleSummary",
		Summary:     "GenerateArticleSummary",
		Description: "Generate article summary using AI",
		Tags:        []string{"ai"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	},
		map[string]any{
			"SSEResponse": protocol.SSEResponse{},
		}, aiHandler.HandleGenerateArticleSummary)

	readerGroup := huma.NewGroup(appGroup, "/reader")

	articleQAGroup := huma.NewGroup(readerGroup, "")
	articleQAGroup.UseMiddleware(middleware.RedisLockMiddleware("articleQA", constant.CtxKeyUserID, 30*time.Second))

	sse.Register(articleQAGroup, huma.Operation{
		OperationID: "generateArticleQA",
		Method:      http.MethodPost,
		Path:        "/articleQA",
		Summary:     "GenerateArticleQA",
		Description: "Generate article Q&A using AI",
		Tags:        []string{"ai"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	},
		map[string]any{
			"SSEResponse": protocol.SSEResponse{},
		}, aiHandler.HandleGenerateArticleQA)

	//
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
