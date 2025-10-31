package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
)

func initTagRouter(v1Group *huma.Group) {
	tagHandler := handler.NewTagHandler()

	tagGroup := huma.NewGroup(v1Group, "/tag")
	tagGroup.UseMiddleware(middleware.JwtMiddlewareForHuma())

	huma.Register(tagGroup, huma.Operation{
		OperationID: "listTags",
		Method:      http.MethodGet,
		Path:        "/list",
		Summary:     "ListTags",
		Description: "List all tags with pagination",
		Tags:        []string{"tag"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, tagHandler.HandleListTags)

	huma.Register(tagGroup, huma.Operation{
		OperationID: "getTagInfo",
		Method:      http.MethodGet,
		Path:        "/{tagID}",
		Summary:     "GetTagInfo",
		Description: "Get tag detail by ID",
		Tags:        []string{"tag"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, tagHandler.HandleGetTagInfo)

	securedGroup := huma.NewGroup(tagGroup, "")
	securedGroup.UseMiddleware(middleware.LimitUserPermissionMiddlewareForHuma("tagService", model.PermissionCreator))

	huma.Register(securedGroup, huma.Operation{
		OperationID: "createTag",
		Method:      http.MethodPost,
		Path:        "/",
		Summary:     "CreateTag",
		Description: "Create a new tag",
		Tags:        []string{"tag"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, tagHandler.HandleCreateTag)

	huma.Register(securedGroup, huma.Operation{
		OperationID: "updateTag",
		Method:      http.MethodPatch,
		Path:        "/{tagID}",
		Summary:     "UpdateTag",
		Description: "Update tag information",
		Tags:        []string{"tag"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, tagHandler.HandleUpdateTag)

	huma.Register(securedGroup, huma.Operation{
		OperationID: "deleteTag",
		Method:      http.MethodDelete,
		Path:        "/{tagID}",
		Summary:     "DeleteTag",
		Description: "Delete tag by ID",
		Tags:        []string{"tag"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, tagHandler.HandleDeleteTag)
}
