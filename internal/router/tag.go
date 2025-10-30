package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
)

func initTagRouter(tagGroup *huma.Group) {
	tagHandler := handler.NewTagHandler()

	tagGroup.UseMiddleware(middleware.JwtMiddlewareForHuma())

	// 列出标签
	huma.Register(tagGroup, huma.Operation{
		OperationID: "listTags",
		Method:      http.MethodGet,
		Path:        "/list",
		Summary:     "ListTags",
		Description: "Get a paginated list of tags",
		Tags:        []string{"tag"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, tagHandler.HandleListTags)

	// 创建标签
	huma.Register(tagGroup, huma.Operation{
		OperationID: "createTag",
		Method:      http.MethodPost,
		Path:        "/",
		Summary:     "CreateTag",
		Description: "Create a new tag (requires creator permission)",
		Tags:        []string{"tag"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, tagHandler.HandleCreateTag)

	// 获取标签信息
	huma.Register(tagGroup, huma.Operation{
		OperationID: "getTagInfo",
		Method:      http.MethodGet,
		Path:        "/{tagID}",
		Summary:     "GetTagInfo",
		Description: "Get detailed information about a specific tag by ID",
		Tags:        []string{"tag"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, tagHandler.HandleGetTagInfo)

	// 更新标签
	huma.Register(tagGroup, huma.Operation{
		OperationID: "updateTag",
		Method:      http.MethodPatch,
		Path:        "/{tagID}",
		Summary:     "UpdateTag",
		Description: "Update an existing tag (requires creator permission)",
		Tags:        []string{"tag"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, tagHandler.HandleUpdateTag)

	// 删除标签
	huma.Register(tagGroup, huma.Operation{
		OperationID: "deleteTag",
		Method:      http.MethodDelete,
		Path:        "/{tagID}",
		Summary:     "DeleteTag",
		Description: "Delete a tag (requires creator permission)",
		Tags:        []string{"tag"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, tagHandler.HandleDeleteTag)
}
